// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

import (
	"container/list"
	"context"
	"sync"

	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	relayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
)

var (
	chainlog = log.New("module", "blockchain")
	syncCfg  = &relayerTypes.SyncTxReceiptConfig{}
)

type PushTask struct {
	Receiver string
}

type RePushTask struct {
	Receiver string
	Count    uint64
}

type BlockChainSerivce struct {
	syncBlock          *syncblock.SyncBlock
	rpcLaddr           string //用户向指定的blockchain节点进行rpc调用
	fetchBlockPeriodMs int32
	confirmBlocks      int32 //块的成熟度，只有超过该确认数值的block才会进行解析并使用其结果
	runcount           int32
	db                 dbm.DB
	syncBlockChan      <-chan int64
	//所有发送到移动端的消息都需要通过channel汇总到一个地方，然后进行并发地发送到多个移动端
	//目的：这样就可以顺序地发送每个地址名下的信息(由storeID决定发送先后顺序)，如果发生panic等情况
	//就可以记录下当前成功发送的store item
	pushTasksChan chan PushTask
	//repushTasksChan     chan RePushTask //重新推送任务队列
	//用户地址到url等信息的映射,其中的地址需要事先在闪电网络服务进行注册，否则不会帮他推送相应的消息
	addr2syncInfo       map[string]*lnsTypes.SyncInfo4App
	addr2repush         map[string]uint64
	pendingSends        *list.List
	rw4url              sync.RWMutex
	rwRepush            sync.RWMutex
	height              int64 //当前区块高度  +++++++++||++++++++++++||++++++++++||
	heightSync2App      int64 //已经同步高度           ^             ^           ^
	matDegree           int32 //成熟度         heightSync2App    matDegress   height
	ctx                 context.Context
	wg                  sync.WaitGroup
	newHeightConcurrnet int32 //通知高度goroutine最大并发数
	send2AppConcurrnet  int32
}

func New(wg sync.WaitGroup, ctx context.Context, syncTxConfig *relayerTypes.TxReceiptConfig) *BlockChainSerivce {
	blockChainSerivce := &BlockChainSerivce{}
	//blockChainSerivce.quit = make(chan struct{})
	blockChainSerivce.rpcLaddr = syncTxConfig.Chain33Host
	blockChainSerivce.db = dbm.NewDB("lns_blockchainservice", syncTxConfig.Dbdriver, syncTxConfig.DbPath, syncTxConfig.DbCache)
	blockChainSerivce.rw4url = *new(sync.RWMutex)
	blockChainSerivce.rwRepush = *new(sync.RWMutex)
	blockChainSerivce.ctx = ctx
	blockChainSerivce.wg = wg
	blockChainSerivce.newHeightConcurrnet = syncTxConfig.NewHeightConcurrnet
	blockChainSerivce.send2AppConcurrnet = syncTxConfig.Send2AppConcurrnet
	blockChainSerivce.send2AppConcurrnet = syncTxConfig.Send2AppConcurrnet
	blockChainSerivce.matDegree = syncTxConfig.MaturityDegree

	blockChainSerivce.addr2syncInfo = make(map[string]*lnsTypes.SyncInfo4App)
	blockChainSerivce.addr2repush = make(map[string]uint64)
	blockChainSerivce.pushTasksChan = make(chan PushTask, BufferSize)
	blockChainSerivce.pendingSends = list.New()

	syncCfg.Chain33Host = syncTxConfig.Chain33Host
	syncCfg.PushHost = syncTxConfig.PushHost
	syncCfg.PushName = syncTxConfig.PushName
	syncCfg.PushBind = syncTxConfig.PushBind

	return blockChainSerivce
}

//Close 关闭区块链服务
func (chain *BlockChainSerivce) Close() {
	syncblock.StopSyncBlock()
	//chain.quit <- struct{}{}
}

//SetQueueClient 设置队列
func (chain *BlockChainSerivce) Start() {
	//一次最多缓存10个新高度,正常情况是不需要缓存的，就为了应付分叉之后的删除和增加多个块的场景
	syncChan := make(chan int64, SyncBlockChanSize)
	chain.syncBlockChan = syncChan
	chain.syncBlock = syncblock.StartSyncBlock(syncCfg, syncChan, chain.db)
	chain.loadFromDB()

	chain.wg.Add(3)
	go chain.blockProc()
	go chain.sendMsg2App()
	go chain.repushMsg2App()
	chainlog.Info("BlockChainSerivce is started successfully...")
}

func (chain *BlockChainSerivce) loadFromDB() {
	//如果第一启动，读取不到高度值
	height, err := chain.syncBlock.LoadLastBlockHeight()
	if nil != err {
		chain.height = 0
		chainlog.Info("Failed to LoadLastBlockHeight")
	}
	chain.height = height
	heightSync2App, err := chain.loadLastSyncHeight()
	if nil != err {
		chainlog.Info("Failed to loadLastSyncHeight")
		heightSync2App = 0
	}
	chain.heightSync2App = heightSync2App
	chainlog.Info("LoadInfoFromDB", "blockHeight", height, "syncHeight2App", heightSync2App)
	if nil != chain.loadAllSyncInfo4User() {
		panic("Failed to loadAllSyncInfo4User")
	}
	return
}

func (chain *BlockChainSerivce) GetRpcAddr() string {
	return chain.rpcLaddr
}
