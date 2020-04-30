package chain33

import (
	"fmt"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"sync/atomic"
)

var (
	lastSyncHeightPrefix              = []byte("lastSyncHeight:")
	chain33ToEthBurnLockTxHashPrefix  = "chain33ToEthBurnLockTxHash"
	chain33ToEthBurnLockTxTotalAmount = []byte("chain33ToEthBurnLockTxTotalAmount")
	EthTxStatusCheckedIndex           = []byte("EthTxStatusCheckedIndex")
)

func calcRelay2EthTxhash(txindex int64) []byte {
	return []byte(fmt.Sprintf("%s-%012d", chain33ToEthBurnLockTxHashPrefix, txindex))
}

func (chain33Relayer *Chain33Relayer) updateTotalTxAmount2Eth(total int64) error {
	totalTx := &types.Int64{
		Data: atomic.LoadInt64(&chain33Relayer.totalTx4Chain33ToEth),
	}
	//更新成功见证的交易数
	return chain33Relayer.db.Set(chain33ToEthBurnLockTxTotalAmount, types.Encode(totalTx))
}

func (chain33Relayer *Chain33Relayer) getTotalTxAmount2Eth() int64 {
	totalTx, _ := utils.LoadInt64FromDB(chain33ToEthBurnLockTxTotalAmount, chain33Relayer.db)
	return totalTx
}

func (chain33Relayer *Chain33Relayer) setLastestRelay2EthTxhash(status, txhash string, txIndex int64) error {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{
		Status:status,
		Txhash:txhash,
	}
	data := types.Encode(ethTxStatus)
	return chain33Relayer.db.Set(key, data)
}

func (chain33Relayer *Chain33Relayer) getEthTxhash(txIndex int64) (common.Hash, error) {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{}
	data, err := chain33Relayer.db.Get(key)
	if nil != err {
		return common.Hash{}, err
	}
	err = types.Decode(data, ethTxStatus)
	if nil != err {
		return common.Hash{}, err
	}
	return common.HexToHash(ethTxStatus.Txhash), nil
}

func (chain33Relayer *Chain33Relayer) setStatusCheckedIndex(txIndex int64) error {
	index := &types.Int64{
		Data:txIndex,
	}
	data := types.Encode(index)
	return chain33Relayer.db.Set(EthTxStatusCheckedIndex, data)
}

func (chain33Relayer *Chain33Relayer) getStatusCheckedIndex() int64 {
	index, _ := utils.LoadInt64FromDB(EthTxStatusCheckedIndex, chain33Relayer.db)
	return index
}

//获取上次同步到app的高度
func (chain33Relayer *Chain33Relayer) loadLastSyncHeight() int64 {
	height, err := utils.LoadInt64FromDB(lastSyncHeightPrefix, chain33Relayer.db)
	if nil != err && err != types.ErrHeightNotExist {
		relayerLog.Error("loadLastSyncHeight", "err:", err.Error())
		return 0
	}
	return height
}

func (chain33Relayer *Chain33Relayer) setLastSyncHeight(syncHeight int64) {
	bytes := types.Encode(&types.Int64{Data: syncHeight})
	_ = chain33Relayer.db.Set(lastSyncHeightPrefix, bytes)
}
