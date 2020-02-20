package executor

import (
	"github.com/33cn/chain33/types"
	x2ethereumtypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * 实现交易的链上执行接口
 * 关键数据上链（statedb）并生成交易回执（log）
 */

func (x *x2ethereum) Exec_EthBridgeClaim(payload *x2ethereumtypes.EthBridgeClaim, tx *types.Transaction, index int) (*types.Receipt, error) {
	var receipt *types.Receipt
	//implement code
	return receipt, nil
}

func (x *x2ethereum) Exec_MsgBurn(payload *x2ethereumtypes.MsgBurn, tx *types.Transaction, index int) (*types.Receipt, error) {
	var receipt *types.Receipt
	//implement code
	return receipt, nil
}

func (x *x2ethereum) Exec_MsgLock(payload *x2ethereumtypes.MsgLock, tx *types.Transaction, index int) (*types.Receipt, error) {
	var receipt *types.Receipt
	//implement code
	return receipt, nil
}
