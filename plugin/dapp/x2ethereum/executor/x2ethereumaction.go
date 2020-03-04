package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	"github.com/33cn/chain33/common/address"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	token "github.com/33cn/plugin/plugin/dapp/token/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/ethbridge"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/oracle"
	types2 "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"github.com/pkg/errors"
)

// stateDB存储KV:
//		id --> DBProphecy
//
//
//		ValidatorMapsKey -- > ValidatorMaps arrays
//

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

func newAction(a *x2ethereum, tx *types.Transaction, index int32) *action {
	hash := tx.Hash()
	fromaddr := tx.From()

	moduleAddress, err := address.NewAddrFromString(types2.ModuleName)
	if err != nil {
		return nil
	}
	addressMap := make(map[string]string)
	addressMap[types2.ModuleName] = moduleAddress.String()
	supplyKeeper := common.NewKeeper(addressMap)
	oracleKeeper := oracle.NewKeeper(a.GetStateDB(), types2.DefaultConsensusNeeded)

	return &action{a.GetAPI(), a.GetCoinsAccount(), a.GetStateDB(), hash, fromaddr,
		a.GetBlockTime(), a.GetHeight(), index, dapp.ExecAddress(string(tx.Execer)), ethbridge.NewKeeper(supplyKeeper, oracleKeeper, a.GetStateDB())}
}

//ethereum ---> chain33
func (a *action) procMsgEthBridgeClaim(ethBridgeClaim *types2.EthBridgeClaim) (*types.Receipt, error) {
	var receipt *types.Receipt
	msgEthBridgeClaim := ethbridge.NewMsgCreateEthBridgeClaim(*ethBridgeClaim)
	if err := msgEthBridgeClaim.ValidateBasic(); err != nil {
		return nil, err
	}

	status, err := a.keeper.ProcessClaim(*ethBridgeClaim)
	if err != nil {
		return nil, err
	}

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
		receipt.KV = append(receipt.KV, receipt.KV...)
		receipt.Logs = append(receipt.Logs, r.Logs...)
	}

	msgEthBridgeClaimBytes, err := json.Marshal(msgEthBridgeClaim)
	if err != nil {
		return nil, types.ErrMarshal
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return nil, types.ErrMarshal
	}

	receipt.KV = append(receipt.KV, &types.KeyValue{Key: msgEthBridgeClaimBytes, Value: statusBytes})

	execlog := &types.ReceiptLog{Ty: types2.TyLogMsgEthBridgeClaim, Log: types.Encode(ethBridgeClaim)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk
	return receipt, nil
}

func (a *action) procMsgLock(msgLock *types2.MsgLock) (*types.Receipt, error) {
	accDB, err := a.createAccount(msgLock.LocalCoinExec, msgLock.LocalCoinSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "relay procMsgLock,exec=%s,sym=%s", msgLock.LocalCoinExec, msgLock.LocalCoinSymbol)
	}
	receipt, err := a.keeper.ProcessLock(msgLock.Chain33Sender, a.execaddr, int64(msgLock.Amount), accDB)
	if err != nil {
		return nil, err
	}

	receipt.Ty = types.ExecOk
	return receipt, nil
}

func (a *action) procMsgBurn(msgBurn *types2.MsgBurn) (*types.Receipt, error) {
	accDB, err := a.createAccount(msgBurn.LocalCoinExec, msgBurn.LocalCoinSymbol)
	if err != nil {
		return nil, errors.Wrapf(err, "relay procMsgBurn,exec=%s,sym=%s", msgBurn.LocalCoinExec, msgBurn.LocalCoinSymbol)
	}

	receipt, err := a.keeper.ProcessBurn(msgBurn.Chain33Sender, a.execaddr, int64(msgBurn.Amount), accDB)
	if err != nil {
		return nil, err
	}

	receipt.Ty = types.ExecOk
	return receipt, nil
}

//需要一笔交易来注册validator
//这里注册的validator的power之和可能不为1，需要在内部进行加权
//返回的回执中，KV包含所有validator的power值，Log中包含本次注册的validator的power值
func (a *action) procMsgLogInValidator(msgLogInValidator *types2.MsgLogInValidator) (*types.Receipt, error) {
	receipt := new(types.Receipt)

	receipt, err := a.keeper.ProcessLogInValidator(msgLogInValidator.Address, msgLogInValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &types.ReceiptLog{Ty: types2.TyLogMsgLogInValidator, Log: types.Encode(msgLogInValidator)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk
	return receipt, nil
}

func (a *action) procMsgLogOutValidator(msgLogOutValidator *types2.MsgLogOutValidator) (*types.Receipt, error) {
	receipt := new(types.Receipt)

	receipt, err := a.keeper.ProcessLogOutValidator(msgLogOutValidator.Address, msgLogOutValidator.Power)
	if err != nil {
		return nil, err
	}

	execlog := &types.ReceiptLog{Ty: types2.TyLogMsgLogOutValidator, Log: types.Encode(msgLogOutValidator)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk
	return receipt, nil
}

func (a *action) procMsgSetConsensusNeeded(msgSetConsensusNeeded *types2.MsgSetConsensusNeeded) (*types.Receipt, error) {
	receipt := new(types.Receipt)

	receipt, err := a.keeper.ProcessSetConsensusNeeded(msgSetConsensusNeeded.Power)
	if err != nil {
		return nil, err
	}

	execlog := &types.ReceiptLog{Ty: types2.TyLogMsgSetConsensusNeeded, Log: types.Encode(msgSetConsensusNeeded)}
	receipt.Logs = append(receipt.Logs, execlog)

	receipt.Ty = types.ExecOk
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
