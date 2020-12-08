package executor

import (
	"strconv"

	"github.com/33cn/dplatform/types"
	x2eTy "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * 实现交易相关数据本地执行，数据不上链
 * 非关键数据，本地存储(localDB), 用于辅助查询，效率高
 */

func (x *x2ethereum) ExecLocal_Eth2DplatformLock(payload *x2eTy.Eth2Dplatform, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_Eth2DplatformBurn(payload *x2eTy.Eth2Dplatform, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_DplatformToEthBurn(payload *x2eTy.DplatformToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_DplatformToEthLock(payload *x2eTy.DplatformToEth, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := x.execLocal(receiptData)
	if err != nil {
		return set, err
	}
	return x.addAutoRollBack(tx, set.KV), nil
}

func (x *x2ethereum) ExecLocal_AddValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_RemoveValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_ModifyPower(payload *x2eTy.MsgValidator, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

func (x *x2ethereum) ExecLocal_SetConsensusThreshold(payload *x2eTy.MsgConsensusThreshold, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	//implement code
	return x.addAutoRollBack(tx, dbSet.KV), nil
}

//设置自动回滚
func (x *x2ethereum) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {
	dbSet := &types.LocalDBSet{}
	dbSet.KV = x.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}

func (x *x2ethereum) execLocal(receiptData *types.ReceiptData) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	for _, log := range receiptData.Logs {
		switch log.Ty {
		case x2eTy.TyEth2DplatformLog:
			var receiptEth2Dplatform x2eTy.ReceiptEth2Dplatform
			err := types.Decode(log.Log, &receiptEth2Dplatform)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Dplatform.IssuerDotSymbol, receiptEth2Dplatform.TokenAddress, x2eTy.DirEth2Dplatform, "lock"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}
			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptEth2Dplatform.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptEth2Dplatform.IssuerDotSymbol,
				TxType:      "lock",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   1,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Dplatform.IssuerDotSymbol, receiptEth2Dplatform.TokenAddress, x2eTy.DirEth2Dplatform, "lock"),
				Value: TokenAssetsByTxTypeBytes,
			})

			nb, err = x.GetLocalDB().Get(x2eTy.CalTokenSymbolToTokenAddress(receiptEth2Dplatform.IssuerDotSymbol))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var t x2eTy.ReceiptTokenToTokenAddress
			err = types.Decode(nb, &t)
			if err != nil {
				return nil, err
			}
			var exist bool
			for _, addr := range t.TokenAddress {
				if addr == receiptEth2Dplatform.TokenAddress {
					exist = true
				}
			}
			if !exist {
				t.TokenAddress = append(t.TokenAddress, receiptEth2Dplatform.TokenAddress)
			}
			TokenToTokenAddressBytes := types.Encode(&x2eTy.ReceiptTokenToTokenAddress{
				TokenAddress: t.TokenAddress,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolToTokenAddress(receiptEth2Dplatform.IssuerDotSymbol),
				Value: TokenToTokenAddressBytes,
			})
		case x2eTy.TyWithdrawEthLog:
			var receiptEth2Dplatform x2eTy.ReceiptEth2Dplatform
			err := types.Decode(log.Log, &receiptEth2Dplatform)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Dplatform.IssuerDotSymbol, receiptEth2Dplatform.TokenAddress, x2eTy.DirEth2Dplatform, "withdraw"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptEth2Dplatform.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptEth2Dplatform.IssuerDotSymbol,
				TxType:      "withdraw",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   2,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptEth2Dplatform.IssuerDotSymbol, receiptEth2Dplatform.TokenAddress, x2eTy.DirEth2Dplatform, "withdraw"),
				Value: TokenAssetsByTxTypeBytes,
			})
		case x2eTy.TyDplatformToEthLog:
			var receiptDplatformToEth x2eTy.ReceiptDplatformToEth
			err := types.Decode(log.Log, &receiptDplatformToEth)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptDplatformToEth.IssuerDotSymbol, receiptDplatformToEth.TokenContract, x2eTy.DirDplatformToEth, "lock"))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptDplatformToEth.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptDplatformToEth.IssuerDotSymbol,
				TxType:      "lock",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   1,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptDplatformToEth.IssuerDotSymbol, receiptDplatformToEth.TokenContract, x2eTy.DirDplatformToEth, "lock"),
				Value: TokenAssetsByTxTypeBytes,
			})
		case x2eTy.TyWithdrawDplatformLog:
			var receiptDplatformToEth x2eTy.ReceiptDplatformToEth
			err := types.Decode(log.Log, &receiptDplatformToEth)
			if err != nil {
				return nil, err
			}

			nb, err := x.GetLocalDB().Get(x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptDplatformToEth.IssuerDotSymbol, receiptDplatformToEth.TokenContract, x2eTy.DirDplatformToEth, ""))
			if err != nil && err != types.ErrNotFound {
				return nil, err
			}
			var now x2eTy.ReceiptQuerySymbolAssetsByTxType
			err = types.Decode(nb, &now)
			if err != nil {
				return nil, err
			}

			preAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(now.TotalAmount), 64)
			nowAmount, _ := strconv.ParseFloat(x2eTy.TrimZeroAndDot(receiptDplatformToEth.Amount), 64)
			TokenAssetsByTxTypeBytes := types.Encode(&x2eTy.ReceiptQuerySymbolAssetsByTxType{
				TokenSymbol: receiptDplatformToEth.IssuerDotSymbol,
				TxType:      "withdraw",
				TotalAmount: strconv.FormatFloat(preAmount+nowAmount, 'f', 4, 64),
				Direction:   2,
			})
			dbSet.KV = append(dbSet.KV, &types.KeyValue{
				Key:   x2eTy.CalTokenSymbolTotalLockOrBurnAmount(receiptDplatformToEth.IssuerDotSymbol, receiptDplatformToEth.TokenContract, x2eTy.DirDplatformToEth, "withdraw"),
				Value: TokenAssetsByTxTypeBytes,
			})
		default:
			continue
		}
	}
	return dbSet, nil
}
