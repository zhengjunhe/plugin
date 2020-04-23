package chain33

import (
	"fmt"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	"sync/atomic"
)

var (
	lastSyncHeightPrefix              = []byte("lastSyncHeight:")
	chain33ToEthBurnLockTxHashPrefix  = "chain33ToEthBurnLockTxHash"
	chain33ToEthBurnLockTxTotalAmount = []byte("chain33ToEthBurnLockTxTotalAmount")
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

func (chain33Relayer *Chain33Relayer) setLastestRelay2EthTxhash(txhash string, txIndex int64) error {
	key := calcRelay2EthTxhash(txIndex)
	return chain33Relayer.db.Set(key, []byte(txhash))
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
	chain33Relayer.db.Set(lastSyncHeightPrefix, bytes)
}
