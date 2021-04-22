package chain33

import (
	"errors"
	"fmt"
	"sync/atomic"

	dbm "github.com/33cn/chain33/common/db"
	chain33Types "github.com/33cn/chain33/types"
	types "github.com/33cn/chain33/types"
	ebTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/common"
)

//key ...
var (
	lastSyncHeightPrefix              = []byte("lastSyncHeight:")
	chain33ToEthBurnLockTxHashPrefix  = "chain33ToEthBurnLockTxHash"
	chain33ToEthBurnLockTxTotalAmount = []byte("chain33ToEthBurnLockTxTotalAmount")
	EthTxStatusCheckedIndex           = []byte("EthTxStatusCheckedIndex")
	bridgeRegistryAddrOnChain33       = []byte("x2EthBridgeRegistryAddrOnChain33")
	tokenSymbol2AddrPrefix            = []byte("chain33TokenSymbol2AddrPrefix")
)

func tokenSymbol2AddrKey(symbol string) []byte {
	return append(tokenSymbol2AddrPrefix, []byte(fmt.Sprintf("-symbol-%s", symbol))...)
}

func calcRelay2EthTxhash(txindex int64) []byte {
	return []byte(fmt.Sprintf("%s-%012d", chain33ToEthBurnLockTxHashPrefix, txindex))
}

func (chain33Relayer *Relayer4Chain33) updateTotalTxAmount2Eth(total int64) error {
	totalTx := &chain33Types.Int64{
		Data: atomic.LoadInt64(&chain33Relayer.totalTx4Chain33ToEth),
	}
	//更新成功见证的交易数
	return chain33Relayer.db.Set(chain33ToEthBurnLockTxTotalAmount, types.Encode(totalTx))
}

func (chain33Relayer *Relayer4Chain33) getTotalTxAmount2Eth() int64 {
	totalTx, _ := utils.LoadInt64FromDB(chain33ToEthBurnLockTxTotalAmount, chain33Relayer.db)
	return totalTx
}

func (chain33Relayer *Relayer4Chain33) setLastestRelay2EthTxhash(status, txhash string, txIndex int64) error {
	key := calcRelay2EthTxhash(txIndex)
	ethTxStatus := &ebTypes.EthTxStatus{
		Status: status,
		Txhash: txhash,
	}
	data := types.Encode(ethTxStatus)
	return chain33Relayer.db.Set(key, data)
}

func (chain33Relayer *Relayer4Chain33) getEthTxhash(txIndex int64) (common.Hash, error) {
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

func (chain33Relayer *Relayer4Chain33) setStatusCheckedIndex(txIndex int64) error {
	index := &types.Int64{
		Data: txIndex,
	}
	data := types.Encode(index)
	return chain33Relayer.db.Set(EthTxStatusCheckedIndex, data)
}

func (chain33Relayer *Relayer4Chain33) getStatusCheckedIndex() int64 {
	index, _ := utils.LoadInt64FromDB(EthTxStatusCheckedIndex, chain33Relayer.db)
	return index
}

//获取上次同步到app的高度
func (chain33Relayer *Relayer4Chain33) loadLastSyncHeight() int64 {
	height, err := utils.LoadInt64FromDB(lastSyncHeightPrefix, chain33Relayer.db)
	if nil != err && err != types.ErrHeightNotExist {
		relayerLog.Error("loadLastSyncHeight", "err:", err.Error())
		return 0
	}
	return height
}

func (chain33Relayer *Relayer4Chain33) setLastSyncHeight(syncHeight int64) {
	bytes := types.Encode(&types.Int64{Data: syncHeight})
	_ = chain33Relayer.db.Set(lastSyncHeightPrefix, bytes)
}

func (chain33Relayer *Relayer4Chain33) setBridgeRegistryAddr(bridgeRegistryAddr string) error {
	return chain33Relayer.db.Set(bridgeRegistryAddrOnChain33, []byte(bridgeRegistryAddr))
}

func (chain33Relayer *Relayer4Chain33) getBridgeRegistryAddr() (string, error) {
	addr, err := chain33Relayer.db.Get(bridgeRegistryAddrOnChain33)
	if nil != err {
		return "", err
	}
	return string(addr), nil
}

func (chain33Relayer *Relayer4Chain33) SetTokenAddress(token2set ebTypes.TokenAddress) error {
	bytes := chain33Types.Encode(&token2set)
	chain33Relayer.rwLock.Lock()
	chain33Relayer.symbol2Addr[token2set.Symbol] = token2set.Address
	chain33Relayer.rwLock.Unlock()
	return chain33Relayer.db.Set(tokenSymbol2AddrKey(token2set.Symbol), bytes)
}

func (chain33Relayer *Relayer4Chain33) RestoreTokenAddress() error {
	helper := dbm.NewListHelper(chain33Relayer.db)
	datas := helper.List(tokenSymbol2AddrPrefix, nil, 100, dbm.ListASC)
	if nil == datas {
		return nil
	}

	chain33Relayer.rwLock.Lock()
	for _, data := range datas {

		var token2set ebTypes.TokenAddress
		err := chain33Types.Decode(data, &token2set)
		if nil != err {
			return err
		}
		relayerLog.Info("RestoreTokenAddress", "symbol", token2set.Symbol, "address", token2set.Address)
		chain33Relayer.symbol2Addr[token2set.Symbol] = token2set.Address
	}
	chain33Relayer.rwLock.Unlock()
	return nil
}

func (chain33Relayer *Relayer4Chain33) ShowTokenAddress(token2show ebTypes.TokenAddress) (*ebTypes.TokenAddressArray, error) {
	res := &ebTypes.TokenAddressArray{}

	if len(token2show.Symbol) > 0 {
		data, err := chain33Relayer.db.Get(tokenSymbol2AddrKey(token2show.Symbol))
		if err != nil {
			return nil, err
		}
		var token2set ebTypes.TokenAddress
		err = chain33Types.Decode(data, &token2set)
		if nil != err {
			return nil, err
		}
		res.TokenAddress = append(res.TokenAddress, &token2set)
		return res, nil
	}
	helper := dbm.NewListHelper(chain33Relayer.db)
	datas := helper.List(tokenSymbol2AddrPrefix, nil, 100, dbm.ListASC)
	if nil == datas {
		return nil, errors.New("Not found")
	}

	for _, data := range datas {

		var token2set ebTypes.TokenAddress
		err := chain33Types.Decode(data, &token2set)
		if nil != err {
			return nil, err
		}
		res.TokenAddress = append(res.TokenAddress, &token2set)

	}
	return res, nil
}
