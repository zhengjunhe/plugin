package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	chain33types "github.com/33cn/chain33/types"
	token "github.com/33cn/plugin/plugin/dapp/token/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/ethbridge"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/pkg/errors"
)

// stateDB存储KV:
//		ProphecyKey --> DBProphecy
//
//		EthBridgeClaimKey -- > EthBridgeClaim
//
//		ValidatorMapsKey -- > ValidatorMaps arrays
//
//		LastTotalPowerKey -- > totalPower
//
//		ConsensusNeededKey -- > consensusNeeded

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

	moduleAddress := dapp.ExecAddress(types2.ModuleName)
	addressMap := make(map[string]string)
	addressMap[types2.ModuleName] = moduleAddress
	supplyKeeper := common.NewKeeper(addressMap)

	var consensusNeeded float64
	consensusNeededBytes, err := a.GetStateDB().Get(types2.CalConsensusNeededPrefix())
	if err != nil {
		if err == chain33types.ErrNotFound {
			consensusNeeded = types2.DefaultConsensusNeeded
			cb, _ := json.Marshal(types2.MsgSetConsensusNeeded{
				ConsensusNeed: consensusNeeded,
			})
			_ = a.GetStateDB().Set(types2.CalConsensusNeededPrefix(), cb)
		} else {
			return nil
		}
	} else {
		var mc types2.MsgSetConsensusNeeded
		_ = json.Unmarshal(consensusNeededBytes, &mc)
		consensusNeeded = mc.ConsensusNeed
	}
	oracleKeeper := oracle.NewKeeper(a.GetStateDB(), consensusNeeded)

	elog.Info("newAction", "newAction", "done")
	return &action{a.GetAPI(), a.GetCoinsAccount(), a.GetStateDB(), hash, fromaddr,
		a.GetBlockTime(), a.GetHeight(), index, dapp.ExecAddress(string(tx.Execer)), ethbridge.NewKeeper(supplyKeeper, &oracleKeeper, a.GetStateDB())}
}

//ethereum ---> chain33
func (a *action) procMsgEthBridgeClaim(ethBridgeClaim *types2.EthBridgeClaim) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)
	msgEthBridgeClaim := ethbridge.NewMsgCreateEthBridgeClaim(*ethBridgeClaim)
	if err := msgEthBridgeClaim.ValidateBasic(); err != nil {
		return nil, err
	}

	status, err := a.keeper.ProcessClaim(*ethBridgeClaim)
	if err != nil {
		return nil, err
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}

	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalProphecyPrefix(), Value: statusBytes})

	if status.Text == oracle.StatusText(types2.EthBridgeStatus_SuccessStatusText) {
		accDB, err := a.createAccount(ethBridgeClaim.LocalCoinExec, ethBridgeClaim.LocalCoinSymbol)
		if err != nil {
			return nil, errors.Wrapf(err, "relay procMsgEthBridgeClaim,exec=%s,sym=%s", ethBridgeClaim.LocalCoinExec, ethBridgeClaim.LocalCoinSymbol)
		}

		//增发货币到exec地址
		//需要在配置项中配置挖矿
		receipt, err = accDB.ExecIssueCoins(a.execaddr, int64(ethBridgeClaim.Amount))
		if err != nil {
			return nil, err
		}

		r, err := a.keeper.ProcessSuccessfulClaim(status.FinalClaim, a.execaddr, accDB)
		if err != nil {
			return nil, err
		}
		receipt.KV = append(receipt.KV, r.KV...)
		receipt.Logs = append(receipt.Logs, r.Logs...)
	}

	msgEthBridgeClaimBytes, err := json.Marshal(msgEthBridgeClaim)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalEthBridgeClaimPrefix(), Value: msgEthBridgeClaimBytes})

	execlog := &chain33types.ReceiptLog{Ty: types2.TyEthBridgeClaimLog, Log: chain33types.Encode(&types2.ReceiptEthBridgeClaim{
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
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgLock(msgLock *types2.MsgLock) (*chain33types.Receipt, error) {
	accDB, err := a.createAccount(msgLock.LocalCoinExec, msgLock.LocalCoinSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "relay procMsgLock,exec=%s,sym=%s", msgLock.LocalCoinExec, msgLock.LocalCoinSymbol)
	}
	receipt, err := a.keeper.ProcessLock(msgLock.Chain33Sender, a.execaddr, int64(msgLock.Amount), accDB)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyMsgLockLog, Log: chain33types.Encode(&types2.ReceiptLock{
		EthereumChainID:  msgLock.EthereumChainID,
		TokenContract:    msgLock.TokenContract,
		Chain33Sender:    msgLock.Chain33Sender,
		EthereumReceiver: msgLock.EthereumReceiver,
		Amount:           msgLock.Amount,
		LocalCoinSymbol:  msgLock.LocalCoinSymbol,
		LocalCoinExec:    msgLock.LocalCoinExec,
		XTxHash:          a.txhash,
		XHeight:          uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	msgLockBytes, err := json.Marshal(msgLock)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalLockPrefix(), Value: msgLockBytes})

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgBurn(msgBurn *types2.MsgBurn) (*chain33types.Receipt, error) {
	accDB, err := a.createAccount(msgBurn.LocalCoinExec, msgBurn.LocalCoinSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "relay procMsgBurn,exec=%s,sym=%s", msgBurn.LocalCoinExec, msgBurn.LocalCoinSymbol)
	}

	receipt, err := a.keeper.ProcessBurn(msgBurn.Chain33Sender, a.execaddr, int64(msgBurn.Amount), accDB)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyMsgBurnLog, Log: chain33types.Encode(&types2.ReceiptBurn{
		EthereumChainID:  msgBurn.EthereumChainID,
		TokenContract:    msgBurn.TokenContract,
		Chain33Sender:    msgBurn.Chain33Sender,
		EthereumReceiver: msgBurn.EthereumReceiver,
		Amount:           msgBurn.Amount,
		LocalCoinSymbol:  msgBurn.LocalCoinSymbol,
		LocalCoinExec:    msgBurn.LocalCoinExec,
		XTxHash:          a.txhash,
		XHeight:          uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	msgBurnBytes, err := json.Marshal(msgBurn)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalBurnPrefix(), Value: msgBurnBytes})

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

//需要一笔交易来注册validator
//这里注册的validator的power之和可能不为1，需要在内部进行加权
//返回的回执中，KV包含所有validator的power值，Log中包含本次注册的validator的power值
func (a *action) procMsgLogInValidator(msgLogInValidator *types2.MsgValidator) (*chain33types.Receipt, error) {
	elog.Info("procMsgLogInValidator", "start", msgLogInValidator)

	receipt, err := a.keeper.ProcessLogInValidator(msgLogInValidator.Address, msgLogInValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyMsgLogInValidatorLog, Log: chain33types.Encode(&types2.ReceiptLogInOut{
		Address: msgLogInValidator.Address,
		Power:   msgLogInValidator.Power,
		XTxHash: a.txhash,
		XHeight: uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgLogOutValidator(msgLogOutValidator *types2.MsgValidator) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)

	receipt, err := a.keeper.ProcessLogOutValidator(msgLogOutValidator.Address, msgLogOutValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyMsgLogOutValidatorLog, Log: chain33types.Encode(&types2.ReceiptLogInOut{
		Address: msgLogOutValidator.Address,
		Power:   msgLogOutValidator.Power,
		XTxHash: a.txhash,
		XHeight: uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) procMsgSetConsensusNeeded(msgSetConsensusNeeded *types2.MsgSetConsensusNeeded) (*chain33types.Receipt, error) {
	receipt := new(chain33types.Receipt)

	preConsensusNeeded, nowConsensusNeeded, err := a.keeper.ProcessSetConsensusNeeded(msgSetConsensusNeeded.ConsensusNeed)
	if err != nil {
		return nil, err
	}

	execlog := &chain33types.ReceiptLog{Ty: types2.TyMsgSetConsensusNeededLog, Log: chain33types.Encode(&types2.ReceiptSetConsensusNeeded{
		PreConsensusNeeded: preConsensusNeeded,
		NowConsensusNeeded: nowConsensusNeeded,
		XTxHash:            a.txhash,
		XHeight:            uint64(a.height),
	})}
	receipt.Logs = append(receipt.Logs, execlog)

	msgSetConsensusNeededBytes, err := json.Marshal(msgSetConsensusNeeded)
	if err != nil {
		return nil, chain33types.ErrMarshal
	}
	receipt.KV = append(receipt.KV, &chain33types.KeyValue{Key: types2.CalConsensusNeededPrefix(), Value: msgSetConsensusNeededBytes})

	receipt.Ty = chain33types.ExecOk
	return receipt, nil
}

func (a *action) createAccount(exec, symbol string) (*account.DB, error) {
	var accDB *account.DB
	cfg := a.api.GetConfig()

	if symbol == "" {
		accDB = account.NewCoinsAccount(cfg)
		accDB.SetDB(a.db)
		return accDB, nil
	}
	if exec == "" {
		exec = token.TokenX
	}
	return account.NewAccountDB(cfg, exec, symbol, a.db)
}
