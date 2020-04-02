package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	chain33types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/ethbridge"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/pkg/errors"
	"strconv"
)

// stateDB存储KV:
//		CalProphecyPrefix --> DBProphecy
//		CalEth2Chain33Prefix -- > Eth2Chain33
//		CalWithdrawEthPrefix -- > Eth2Chain33
//		CalWithdrawChain33Prefix -- > Chain33ToEth
//		CalChain33ToEthPrefix -- > Chain33ToEth
//		CalValidatorMapsPrefix -- > MsgValidator maps
//		CalLastTotalPowerPrefix -- > ReceiptQueryTotalPower
//		CalConsensusThresholdPrefix -- > ReceiptSetConsensusThreshold
//		CalTokenSymbolTotalAmountPrefix -- > ReceiptQuerySymbolAssets

// 当前存在一个问题：
// token的发行需要提前授权，所以账户模型该如何设计？
//
// 解决方案：
// 当eth-->chain33时，采用 mavl-x2ethereum-symbol的账户模型，但是这样的资产是无法提现的，是一个完全虚拟的资产
// 而在chain33-->eth时，采用 mavl-coins-bty 的账户模型（后续可以升级为mavl-token-symbol以支持多个token资产）

// token 合约转币到x2ethereum合约
// 个人账户地址 = mavl-token-symbol-execAddr:aliceAddr
// 不同币种账户地址 = mavl-token-symbol-execAddr

// eth -- > chain33:
// 在 mavl-token-symbol-execAddr 上铸币，然后转到 mavl-token-symbol-execAddr:addr 上
// withdraw 的时候先从mavl-coins-symbol-execAddr:addr 中withdraw到 mavl-token-symbol-execAddr，然后 burn

// chain33 -- > eth:
// 在 mavl-token-symbol-execAddr:addr 上withdraw到 mavl-token-symbol-execAddr 上，然后frozen住
// withdraw 的时候从 mavl-token-symbol-execAddr 上 deposit mavl-token-symbol-execAddr:addr

type action struct {
	api          client.QueueProtocolAPI
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	index        int32
	execaddr     string
	keeper       ethbridge.Keeper
}

func newAction(a *x2ethereum, tx *chain33types.Transaction, index int32) *action {
	hash := tx.Hash()
	fromaddr := tx.From()

	var ConsensusThreshold float64
	consensusNeededBytes, err := a.GetStateDB().Get(types2.CalConsensusThresholdPrefix())
	if err != nil {
		if err == chain33types.ErrNotFound {
			ConsensusThreshold = types2.DefaultConsensusNeeded
			cb, _ := json.Marshal(types2.ReceiptSetConsensusThreshold{
				PreConsensusThreshold: 0,
				NowConsensusThreshold: int64(ConsensusThreshold * 100),
			})
			_ = a.GetStateDB().Set(types2.CalConsensusThresholdPrefix(), cb)
		} else {
			return nil
		}
	} else {
		var mc types2.ReceiptSetConsensusThreshold
		_ = json.Unmarshal(consensusNeededBytes, &mc)
		ConsensusThreshold = float64(mc.NowConsensusThreshold) / 100.0
	}
	oracleKeeper := oracle.NewKeeper(a.GetStateDB(), ConsensusThreshold)

	elog.Info("newAction", "newAction", "done")
	return &action{a.GetAPI(), a.GetCoinsAccount(), a.GetStateDB(), hash, fromaddr,
		a.GetBlockTime(), a.GetHeight(), index, dapp.ExecAddress(string(tx.Execer)), ethbridge.NewKeeper(&oracleKeeper, a.GetStateDB())}
}

//ethereum ---> chain33
func (a *action) procMsgEth2Chain33(ethBridgeClaim *types2.Eth2Chain33) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)
	msgEthBridgeClaim := ethbridge.NewMsgCreateEthBridgeClaim(*ethBridgeClaim)
	if err := msgEthBridgeClaim.ValidateBasic(); err != nil {
		return nil, err
	}

	status, err := a.keeper.ProcessClaim(*ethBridgeClaim)
	if err != nil {
		return nil, err
	}

	ID := strconv.Itoa(int(msgEthBridgeClaim.EthereumChainID)) + strconv.Itoa(int(msgEthBridgeClaim.Nonce)) + msgEthBridgeClaim.EthereumSender

	//记录ethProphecy
	bz, err := a.db.Get(types2.CalProphecyPrefix())
	if err != nil {
		return nil, types2.ErrProphecyGet
	}

	var dbProphecys []oracle.DBProphecy
	var dbProphecy oracle.DBProphecy
	err = json.Unmarshal(bz, &dbProphecys)
	if err != nil {
		return nil, chain33types.ErrUnmarshal
	}

	for _, p := range dbProphecys {
		if p.ID == ID {
			dbProphecy = p
			break
		}
	}

	dRes, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return nil, err
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{
		Key:   types2.CalProphecyPrefix(),
		Value: bz,
	})
	receipt.Logs = append(receipt.Logs, &chain33types.ReceiptLog{Ty: types2.TyProphecyLog, Log: chain33types.Encode(&types2.ReceiptEthProphecy{
		ID: dRes.ID,
		Status: &types2.ProphecyStatus{
			Text:       types2.EthBridgeStatus(dRes.Status.Text),
			FinalClaim: dRes.Status.FinalClaim,
		},
		ClaimValidators: dRes.ClaimValidators,
		ValidatorClaims: dRes.ValidatorClaims,
	})})

	if status.Text == oracle.StatusText(types2.EthBridgeStatus_SuccessStatusText) {
		accDB, err := account.NewAccountDB(a.api.GetConfig(), msgEthBridgeClaim.LocalCoinExec, msgEthBridgeClaim.LocalCoinSymbol, a.db)
		if err != nil {
			return nil, errors.Wrapf(err, "relay procMsgEth2Chain33,exec=%s,sym=%s", msgEthBridgeClaim.LocalCoinExec, msgEthBridgeClaim.LocalCoinSymbol)
		}

		r, err := a.keeper.ProcessSuccessfulClaimForLock(status.FinalClaim, a.execaddr, ethBridgeClaim.LocalCoinSymbol, accDB)
		if err != nil {
			return nil, err
		}
		receipt.KV = append(receipt.KV, r.KV...)
		receipt.Logs = append(receipt.Logs, r.Logs...)

		// 记录该token的总量
		var resAmount uint64
		amount, err := a.getTotalAmountByTokenSymbol(msgEthBridgeClaim.LocalCoinSymbol, types2.DirEth2Chain33)
		if err != nil {
			if err != chain33types.ErrNotFound {
				return nil, err
			} else {
				resAmount = msgEthBridgeClaim.Amount
			}
		} else {
			resAmount = amount + msgEthBridgeClaim.Amount
		}
		symbolAssets := types2.ReceiptQuerySymbolAssets{
			TokenSymbol: msgEthBridgeClaim.LocalCoinSymbol,
			TotalAmount: resAmount,
		}
		symbolAssetsBytes, _ := json.Marshal(symbolAssets)
		receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalTokenSymbolTotalAmountPrefix(msgEthBridgeClaim.LocalCoinSymbol, types2.DirEth2Chain33), Value: symbolAssetsBytes})

		assetsLogs := &chain33types.ReceiptLog{
			Ty: types2.TySymbolAssetsLog,
			Log: chain33types.Encode(&types2.ReceiptQuerySymbolAssets{
				TokenSymbol: msgEthBridgeClaim.LocalCoinSymbol,
				TotalAmount: resAmount,
			})}
		receipt.Logs = append(receipt.Logs, assetsLogs)

		//记录成功lock的日志
		msgEthBridgeClaimBytes, err := json.Marshal(msgEthBridgeClaim)
		if err != nil {
			return nil, chain33types.ErrMarshal
		}
		receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalEth2Chain33Prefix(), Value: msgEthBridgeClaimBytes})

		execlog := &chain33types.ReceiptLog{Ty: types2.TyEth2Chain33Log, Log: chain33types.Encode(&types2.ReceiptEth2Chain33{
			EthereumChainID:       msgEthBridgeClaim.EthereumChainID,
			BridgeContractAddress: msgEthBridgeClaim.BridgeContractAddress,
			Nonce:                 msgEthBridgeClaim.Nonce,
			LocalCoinSymbol:       msgEthBridgeClaim.LocalCoinSymbol,
			LocalCoinExec:         msgEthBridgeClaim.LocalCoinExec,
			TokenContractAddress:  msgEthBridgeClaim.TokenContractAddress,
			EthereumSender:        msgEthBridgeClaim.EthereumSender,
			Chain33Receiver:       msgEthBridgeClaim.Chain33Receiver,
			ValidatorAddress:      msgEthBridgeClaim.ValidatorAddress,
			Amount:                msgEthBridgeClaim.Amount,
			ClaimType:             msgEthBridgeClaim.ClaimType,
			EthSymbol:             msgEthBridgeClaim.EthSymbol,
			XTxHash:               a.txhash,
			XHeight:               uint64(a.height),
			ProphecyID:            ID,
		})}
		receipt.Logs = append(receipt.Logs, execlog)

	}

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procWithdrawEth(withdrawEth *types2.Eth2Chain33) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)
	msgWithdrawEth := ethbridge.NewMsgCreateEthBridgeClaim(*withdrawEth)
	if err := msgWithdrawEth.ValidateBasic(); err != nil {
		return nil, err
	}

	status, err := a.keeper.ProcessClaim(*withdrawEth)
	if err != nil {
		return nil, err
	}

	ID := strconv.Itoa(int(msgWithdrawEth.EthereumChainID)) + strconv.Itoa(int(msgWithdrawEth.Nonce)) + msgWithdrawEth.EthereumSender

	//记录ethProphecy
	bz, err := a.db.Get(types2.CalProphecyPrefix())
	if err != nil {
		return nil, types2.ErrProphecyGet
	}

	var dbProphecys []oracle.DBProphecy
	var dbProphecy oracle.DBProphecy
	err = json.Unmarshal(bz, &dbProphecys)
	if err != nil {
		return nil, chain33types.ErrUnmarshal
	}

	for _, p := range dbProphecys {
		if p.ID == ID {
			dbProphecy = p
			break
		}
	}

	dRes, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return nil, err
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{
		Key:   types2.CalProphecyPrefix(),
		Value: bz,
	})
	receipt.Logs = append(receipt.Logs, &chain33types.ReceiptLog{Ty: types2.TyProphecyLog, Log: chain33types.Encode(&types2.ReceiptEthProphecy{
		ID: dRes.ID,
		Status: &types2.ProphecyStatus{
			Text:       types2.EthBridgeStatus(dRes.Status.Text),
			FinalClaim: dRes.Status.FinalClaim,
		},
		ClaimValidators: dRes.ClaimValidators,
		ValidatorClaims: dRes.ValidatorClaims,
	})})

	if status.Text == oracle.StatusText(types2.EthBridgeStatus_WithdrawedStatusText) {
		accDB, err := account.NewAccountDB(a.api.GetConfig(), msgWithdrawEth.LocalCoinExec, msgWithdrawEth.LocalCoinSymbol, a.db)
		if err != nil {
			return nil, errors.Wrapf(err, "relay procWithdrawEth,exec=%s,sym=%s", withdrawEth.LocalCoinExec, withdrawEth.LocalCoinSymbol)
		}

		r, err := a.keeper.ProcessSuccessfulClaimForBurn(status.FinalClaim, a.execaddr, withdrawEth.LocalCoinSymbol, accDB)
		if err != nil {
			return nil, err
		}
		receipt.KV = append(receipt.KV, r.KV...)
		receipt.Logs = append(receipt.Logs, r.Logs...)

		// 记录该token的总量
		var resAmount uint64
		amount, err := a.getTotalAmountByTokenSymbol(msgWithdrawEth.LocalCoinSymbol, types2.DirEth2Chain33)
		if err != nil {
			return nil, err
		} else {
			resAmount = amount - msgWithdrawEth.Amount
		}
		symbolAssets := types2.ReceiptQuerySymbolAssets{
			TokenSymbol: msgWithdrawEth.LocalCoinSymbol,
			TotalAmount: resAmount,
		}
		symbolAssetsBytes, _ := json.Marshal(symbolAssets)
		receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalTokenSymbolTotalAmountPrefix(msgWithdrawEth.LocalCoinSymbol, types2.DirEth2Chain33), Value: symbolAssetsBytes})

		assetsLogs := &chain33types.ReceiptLog{
			Ty: types2.TySymbolAssetsLog,
			Log: chain33types.Encode(&types2.ReceiptQuerySymbolAssets{
				TokenSymbol: msgWithdrawEth.LocalCoinSymbol,
				TotalAmount: resAmount,
			})}
		receipt.Logs = append(receipt.Logs, assetsLogs)

		msgWithdrawEthBytes, err := json.Marshal(msgWithdrawEth)
		if err != nil {
			return nil, chain33types.ErrMarshal
		}
		receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalWithdrawEthPrefix(), Value: msgWithdrawEthBytes})

		execlog := &chain33types.ReceiptLog{Ty: types2.TyWithdrawEthLog, Log: chain33types.Encode(&types2.ReceiptEth2Chain33{
			EthereumChainID:       msgWithdrawEth.EthereumChainID,
			BridgeContractAddress: msgWithdrawEth.BridgeContractAddress,
			Nonce:                 msgWithdrawEth.Nonce,
			LocalCoinSymbol:       msgWithdrawEth.LocalCoinSymbol,
			LocalCoinExec:         msgWithdrawEth.LocalCoinExec,
			TokenContractAddress:  msgWithdrawEth.TokenContractAddress,
			EthereumSender:        msgWithdrawEth.EthereumSender,
			Chain33Receiver:       msgWithdrawEth.Chain33Receiver,
			ValidatorAddress:      msgWithdrawEth.ValidatorAddress,
			Amount:                msgWithdrawEth.Amount,
			ClaimType:             msgWithdrawEth.ClaimType,
			EthSymbol:             msgWithdrawEth.EthSymbol,
			XTxHash:               a.txhash,
			XHeight:               uint64(a.height),
			ProphecyID:            ID,
		})}
		receipt.Logs = append(receipt.Logs, execlog)

	}

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgLock(msgLock *types2.Chain33ToEth) (*chain33types.Receipt, error) {
	accDB := account.NewCoinsAccount(a.api.GetConfig())
	accDB.SetDB(a.db)
	receipt, err := a.keeper.ProcessLock(msgLock.Chain33Sender, a.execaddr, int64(msgLock.Amount), accDB)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyChain33ToEthLog, Log: chain33types.Encode(&types2.ReceiptChain33ToEth{
		TokenContract:    msgLock.TokenContract,
		Chain33Sender:    msgLock.Chain33Sender,
		EthereumReceiver: msgLock.EthereumReceiver,
		Amount:           msgLock.Amount,
		EthSymbol:        msgLock.EthSymbol,
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	msgLockBytes, err := json.Marshal(msgLock)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalChain33ToEthPrefix(), Value: msgLockBytes})

	// 记录锁定总资产
	var resAmount uint64
	amount, err := a.getTotalAmountByTokenSymbol(msgLock.LocalCoinSymbol, types2.DirChain33ToEth)
	if err != nil {
		if err != chain33types.ErrNotFound {
			return nil, err
		} else {
			resAmount = msgLock.Amount
		}
	} else {
		resAmount = amount + msgLock.Amount
	}
	symbolAssets := types2.ReceiptQuerySymbolAssets{
		TokenSymbol: msgLock.LocalCoinSymbol,
		TotalAmount: resAmount,
	}
	symbolAssetsBytes, _ := json.Marshal(symbolAssets)
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalTokenSymbolTotalAmountPrefix(msgLock.LocalCoinSymbol, types2.DirChain33ToEth), Value: symbolAssetsBytes})

	assetsLogs := &chain33types.ReceiptLog{
		Ty: types2.TySymbolAssetsLog,
		Log: chain33types.Encode(&types2.ReceiptQuerySymbolAssets{
			TokenSymbol: msgLock.LocalCoinSymbol,
			TotalAmount: resAmount,
		})}
	receipt.Logs = append(receipt.Logs, assetsLogs)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgBurn(msgBurn *types2.Chain33ToEth) (*chain33types.Receipt, error) {
	accDB := account.NewCoinsAccount(a.api.GetConfig())
	accDB.SetDB(a.db)
	receipt, err := a.keeper.ProcessBurn(msgBurn.Chain33Sender, a.execaddr, int64(msgBurn.Amount), accDB)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyWithdrawChain33Log, Log: chain33types.Encode(&types2.ReceiptChain33ToEth{
		TokenContract:    msgBurn.TokenContract,
		Chain33Sender:    msgBurn.Chain33Sender,
		EthereumReceiver: msgBurn.EthereumReceiver,
		Amount:           msgBurn.Amount,
		EthSymbol:        msgBurn.EthSymbol,
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	msgBurnBytes, err := json.Marshal(msgBurn)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalWithdrawChain33Prefix(), Value: msgBurnBytes})
	// 记录锁定总资产
	// todo
	// 存 0 的时候读出来是null?
	var resAmount uint64
	amount, err := a.getTotalAmountByTokenSymbol(msgBurn.LocalCoinSymbol, types2.DirChain33ToEth)
	if err != nil {
		return nil, err
	} else {
		resAmount = amount - msgBurn.Amount
	}
	symbolAssets := types2.ReceiptQuerySymbolAssets{
		TokenSymbol: msgBurn.LocalCoinSymbol,
		TotalAmount: resAmount,
	}
	symbolAssetsBytes, _ := json.Marshal(symbolAssets)
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalTokenSymbolTotalAmountPrefix(msgBurn.LocalCoinSymbol, types2.DirChain33ToEth), Value: symbolAssetsBytes})

	assetsLogs := &chain33types.ReceiptLog{
		Ty: types2.TySymbolAssetsLog,
		Log: chain33types.Encode(&types2.ReceiptQuerySymbolAssets{
			TokenSymbol: msgBurn.LocalCoinSymbol,
			TotalAmount: resAmount,
		})}
	receipt.Logs = append(receipt.Logs, assetsLogs)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

//需要一笔交易来注册validator
//这里注册的validator的power之和可能不为1，需要在内部进行加权
//返回的回执中，KV包含所有validator的power值，Log中包含本次注册的validator的power值
func (a *action) procAddValidator(msgAddValidator *types2.MsgValidator) (*chain33types.Receipt, error) {
	elog.Info("procAddValidator", "start", msgAddValidator)

	receipt, err := a.keeper.ProcessAddValidator(msgAddValidator.Address, msgAddValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyAddValidatorLog, Log: chain33types.Encode(&types2.ReceiptValidator{
		Address: msgAddValidator.Address,
		Power:   msgAddValidator.Power,
		XTxHash: a.txhash,
		XHeight: uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procRemoveValidator(msgRemoveValidator *types2.MsgValidator) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)

	receipt, err := a.keeper.ProcessRemoveValidator(msgRemoveValidator.Address)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyRemoveValidatorLog, Log: chain33types.Encode(&types2.ReceiptValidator{
		Address: msgRemoveValidator.Address,
		Power:   msgRemoveValidator.Power,
		XTxHash: a.txhash,
		XHeight: uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procModifyValidator(msgModifyValidator *types2.MsgValidator) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)

	receipt, err := a.keeper.ProcessModifyValidator(msgModifyValidator.Address, msgModifyValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyModifyPowerLog, Log: chain33types.Encode(&types2.ReceiptValidator{
		Address: msgModifyValidator.Address,
		Power:   msgModifyValidator.Power,
		XTxHash: a.txhash,
		XHeight: uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgSetConsensusThreshold(msgSetConsensusThreshold *types2.MsgConsensusThreshold) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)

	preConsensusNeeded, nowConsensusNeeded, err := a.keeper.ProcessSetConsensusNeeded(float64(msgSetConsensusThreshold.ConsensusThreshold) / 100.0)
	if err != nil {
		return nil, err
	}

	setConsensusThreshold := &types2.ReceiptSetConsensusThreshold{
		PreConsensusThreshold: int64(preConsensusNeeded * 100),
		NowConsensusThreshold: int64(nowConsensusNeeded * 100),
		XTxHash:               a.txhash,
		XHeight:               uint64(a.height),
	}
	execlog := &chain33types.ReceiptLog{Ty: types2.TySetConsensusThresholdLog, Log: chain33types.Encode(setConsensusThreshold)}
	receipt.Logs = append(receipt.Logs, execlog)

	msgSetConsensusThresholdBytes, err := json.Marshal(setConsensusThreshold)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalConsensusThresholdPrefix(), Value: msgSetConsensusThresholdBytes})

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) getTotalAmountByTokenSymbol(symbol, direction string) (uint64, error) {
	res, err := a.db.Get(types2.CalTokenSymbolTotalAmountPrefix(symbol, direction))
	if err != nil {
		return 0, err
	}
	var tokenAssets types2.ReceiptQuerySymbolAssets
	err = json.Unmarshal(res, &tokenAssets)
	if err != nil {
		return 0, chain33types.ErrUnmarshal
	}
	return tokenAssets.TotalAmount, nil
}
