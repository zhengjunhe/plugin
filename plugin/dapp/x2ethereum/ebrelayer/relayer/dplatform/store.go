package dplatform

import (
	"fmt"
	"sync/atomic"

	"github.com/33cn/dplatform/types"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/common"
)

//key ...
var (
	lastSyncHeightPrefix              = []byte("lastSyncHeight:")
	dplatformToEthBurnLockTxHashPrefix  = "dplatformToEthBurnLockTxHash"
	dplatformToEthBurnLockTxTotalAmount = []byte("dplatformToEthBurnLockTxTotalAmount")
	EthTxStatusCheckedIndex           = []byte("EthTxStatusCheckedIndex")
)

func calcRelay2EthTxhash(txindex int64) []byte {
	return []byte(fmt.Sprintf("%s-%012d", dplatformToEthBurnLockTxHashPrefix, txindex))
}

func (dplatformRelayer *Relayer4Dplatform) updateTotalTxAmount2Eth(total int64) error {
	totalTx := &types.Int64{
		Data: atomic.LoadInt64(&dplatformRelayer.totalTx4DplatformToEth),
	}
	//更新成功见证的交易数
	return dplatformRelayer.db.Set(dplatformToEthBurnLockTxTotalAmount, types.Encode(totalTx))
}

func (dplatformRelayer *Relayer4Dplatform) getTotalTxAmount2Eth() int64 {
	totalTx, _ := utils.LoadInt64FromDB(dplatformToEthBurnLockTxTotalAmount, dplatformRelayer.db)
	return totalTx
}

func (dplatformRelayer *Relayer4Dplatform) setLastestRelay2EthTxhash(status, txhash string, txIndex int64) error {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{
		Status: status,
		Txhash: txhash,
	}
	data := types.Encode(ethTxStatus)
	return dplatformRelayer.db.Set(key, data)
}

func (dplatformRelayer *Relayer4Dplatform) getEthTxhash(txIndex int64) (common.Hash, error) {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{}
	data, err := dplatformRelayer.db.Get(key)
	if nil != err {
		return common.Hash{}, err
	}
	err = types.Decode(data, ethTxStatus)
	if nil != err {
		return common.Hash{}, err
	}
	return common.HexToHash(ethTxStatus.Txhash), nil
}

func (dplatformRelayer *Relayer4Dplatform) setStatusCheckedIndex(txIndex int64) error {
	index := &types.Int64{
		Data: txIndex,
	}
	data := types.Encode(index)
	return dplatformRelayer.db.Set(EthTxStatusCheckedIndex, data)
}

func (dplatformRelayer *Relayer4Dplatform) getStatusCheckedIndex() int64 {
	index, _ := utils.LoadInt64FromDB(EthTxStatusCheckedIndex, dplatformRelayer.db)
	return index
}

//获取上次同步到app的高度
func (dplatformRelayer *Relayer4Dplatform) loadLastSyncHeight() int64 {
	height, err := utils.LoadInt64FromDB(lastSyncHeightPrefix, dplatformRelayer.db)
	if nil != err && err != types.ErrHeightNotExist {
		relayerLog.Error("loadLastSyncHeight", "err:", err.Error())
		return 0
	}
	return height
}

func (dplatformRelayer *Relayer4Dplatform) setLastSyncHeight(syncHeight int64) {
	bytes := types.Encode(&types.Int64{Data: syncHeight})
	_ = dplatformRelayer.db.Set(lastSyncHeightPrefix, bytes)
}
