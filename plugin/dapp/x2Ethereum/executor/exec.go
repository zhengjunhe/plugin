package executor

import (
	"github.com/33cn/chain33/types"
	x2ethereumtypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

/*
 * 实现交易的链上执行接口
 * 关键数据上链（statedb）并生成交易回执（log）
 */

func (x *x2ethereum) Exec_EthBridgeClaim(payload *x2ethereumtypes.EthBridgeClaim, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgEthBridgeClaim(payload)
}

func (x *x2ethereum) Exec_MsgBurn(payload *x2ethereumtypes.MsgBurn, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgBurn(payload)
}

func (x *x2ethereum) Exec_MsgLock(payload *x2ethereumtypes.MsgLock, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgLock(payload)
}

func (x *x2ethereum) Exec_MsgLogInValidator(payload *x2ethereumtypes.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgLogInValidator(payload)
}

func (x *x2ethereum) Exec_MsgLogOutValidator(payload *x2ethereumtypes.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgLogOutValidator(payload)
}

func (x *x2ethereum) Exec_MsgSetConsensusNeeded(payload *x2ethereumtypes.MsgSetConsensusNeeded, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	return action.procMsgSetConsensusNeeded(payload)
}
