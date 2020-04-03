package executor

import (
	"github.com/33cn/chain33/types"
	x2ethereumtypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

/*
 * 实现交易相关数据本地执行，数据不上链
 * 非关键数据，本地存储(localDB), 用于辅助查询，效率高
 */

// todo
// 将锁定和销毁分开存储

func (s *x2ethereum) ExecLocal_Eth2Chain33(payload *x2ethereumtypes.Eth2Chain33, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_WithdrawEth(payload *x2ethereumtypes.Eth2Chain33, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_WithdrawChain33(payload *x2ethereumtypes.Chain33ToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_Chain33ToEth(payload *x2ethereumtypes.Chain33ToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_AddValidator(payload *x2ethereumtypes.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_RemoveValidator(payload *x2ethereumtypes.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_ModifyPower(payload *x2ethereumtypes.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

func (s *x2ethereum) ExecLocal_SetConsensusThreshold(payload *x2ethereumtypes.MsgConsensusThreshold, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return s.addAutoRollBack(tx, dbSet.KV), nil
}

//设置自动回滚
func (s *x2ethereum) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {

	dbSet := &types.LocalDBSet{}
	dbSet.KV = s.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}
