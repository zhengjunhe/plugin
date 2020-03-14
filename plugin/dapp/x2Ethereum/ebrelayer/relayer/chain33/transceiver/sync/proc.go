package sync

import (
	"fmt"
	"math"
	"sync/atomic"

	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	"github.com/33cn/chain33/types"
	"github.com/pkg/errors"
)

// SeqType
const (
	SeqTypeAdd = int32(1)
	SeqTypeDel = int32(2)
)

var (
	syncLastHeight  = []byte("syncLastHeight:")
	txReceiptPrefix = []byte("txReceiptPrefix:")
	lastSequences   = []byte("lastSequences:")
	lastSyncHeight  = []byte("lastSyncHeight:")
	seqOperationType = []string{"SeqTypeAdd", "SeqTypeDel"}
)

var txReceiptCh chan *types.TxReceipts4Subscribe
var resultCh chan error

func init() {
	txReceiptCh = make(chan *types.TxReceipts4Subscribe)
	resultCh = make(chan error)
}

func txReceiptsKey4Height(height int64) []byte {
	return append(txReceiptPrefix, []byte(fmt.Sprintf("%012d", height))...)
}

// pushTxReceipts push block to backend
func pushTxReceipts(txReceipts *types.TxReceipts4Subscribe) error {
	txReceiptCh <- txReceipts
	err := <-resultCh
	return err
}

type SyncTxReceipts struct {
	db              dbm.DB
	seqNum          int64 //当前同步的序列号
	height          int64 //当前区块高度
	syncReceiptChan chan<- int64
	quit            chan struct{}
}

func NewSyncTxReceipts(db dbm.DB, syncReceiptChan chan<- int64) *SyncTxReceipts {
	sync := &SyncTxReceipts{
		db: db,
	}
	sync.seqNum, _ = sync.loadBlockLastSequence()
	sync.height, _ = sync.LoadLastBlockHeight()
	sync.quit = make(chan struct{})
	sync.syncReceiptChan = syncReceiptChan
	sync.initSyncReceiptDataBase()

	return sync
}

//此处添加一个高度为0的空块，只是为了查找下一个比较方便，并不需要使用其信息
func (syncTx *SyncTxReceipts) initSyncReceiptDataBase() {
	txblock0, _ := syncTx.GetTxReceipts(0)
	if nil != txblock0 {
		return
	}
	txsPerBlock := &types.TxReceipts4SubscribePerBlk{
		Height:0,
	}
	syncTx.setTxReceiptsPerBlock(txsPerBlock)
}

func (syncTx *SyncTxReceipts) Stop() {
	close(syncTx.quit)
}

// SaveAndSyncTxs2Relayer save block to db
func (syncTx *SyncTxReceipts) SaveAndSyncTxs2Relayer() {
	for {
		select {
		case txReceipts := <-txReceiptCh:
			if len(txReceipts.TxReceipts) > 0 {
				log.Info("to deal request", "seq", txReceipts.TxReceipts[0].SeqNum, "count", len(txReceipts.TxReceipts))
			}
			syncTx.dealTxReceipts(txReceipts)
		case <-syncTx.quit:
			return
		}
	}
}

// 保存区块步骤
// 1. 记录 seqNumber ->  seq
// 2. 记录 lastseq
// 3. 更新高度
//
// 重启恢复
// 1. 看高度， 对应高度是已经完成的
// 2. 继续重新下一个高度即可。 重复写， 幂等
// 所以不需要恢复过程， 读出高度即可

// 处理输入流程
func (syncTx *SyncTxReceipts) dealTxReceipts(txReceipts *types.TxReceipts4Subscribe) {
	count, start, txReceiptsParsed, err := parseTxReceipts(txReceipts)
	if err != nil {
		resultCh <- err
	}
	// 在app 端保存成功， 但回复ok时，程序挂掉, 记录日志
	if start <= syncTx.seqNum {
		log.Error("dealTxReceipts", "start", start, "current_seq", syncTx.seqNum)
	}
	if start+int64(count-1) < syncTx.seqNum {
		log.Error("dealTxReceipts", "start", start, "current_seq", syncTx.seqNum)
		resultCh <- errors.New("The tx and receipt pushed is old")
		return
	}
	var height int64
	for i := 0; i < count; i++ {
		txsPerBlock := txReceiptsParsed[i]
		if txsPerBlock.AddDelType == SeqTypeAdd {
			syncTx.setTxReceiptsPerBlock(txsPerBlock)
			syncTx.setBlockLastSequence(txsPerBlock.SeqNum)
			syncTx.setBlockHeight(txsPerBlock.Height)
			height = txsPerBlock.Height
		} else {
			//删除分叉区块处理
			syncTx.delTxReceipts(txsPerBlock.Height)
			syncTx.setBlockLastSequence(txsPerBlock.SeqNum)
			height = txsPerBlock.Height - 1
			//删除区块不需要通知新的高度，因为这只会降低未处理区块的成熟度
			syncTx.setBlockHeight(height)
		}
	}
	log.Debug("dealTxReceipts", "seqStart", start, "count", count, "maxBlockHeight", height)
	syncTx.syncReceiptChan <- height
	//发送回复，确认接收成功
	resultCh <- nil
	return
}

func (syncTx *SyncTxReceipts) loadBlockLastSequence() (int64, error) {
	return utils.LoadInt64FromDB(lastSequences, syncTx.db)
}

func (syncTx *SyncTxReceipts) LoadLastBlockHeight() (int64, error) {
	return utils.LoadInt64FromDB(syncLastHeight, syncTx.db)
}

func (syncTx *SyncTxReceipts) setBlockLastSequence(newSequence int64) {
	Sequencebytes := types.Encode(&types.Int64{Data: newSequence})
	syncTx.db.Set(lastSequences, Sequencebytes)
	//同时更新内存中的seq
	syncTx.updateSequence(newSequence)
}

func (syncTx *SyncTxReceipts) setBlockHeight(height int64) {
	bytes := types.Encode(&types.Int64{Data: height})
	syncTx.db.Set(syncLastHeight, bytes)
	atomic.StoreInt64(&syncTx.height, height)
}

func (syncTx *SyncTxReceipts) updateSequence(newSequence int64) {
	atomic.StoreInt64(&syncTx.seqNum, newSequence)
}

func (syncTx *SyncTxReceipts) setTxReceiptsPerBlock(txReceipts *types.TxReceipts4SubscribePerBlk) {
	key := txReceiptsKey4Height(txReceipts.Height)
	value := types.Encode(txReceipts)
	if err := syncTx.db.Set(key, value); nil != err {
		panic("setTxReceiptsPerBlock failed due to:"+err.Error())
	}
}

func (syncTx *SyncTxReceipts) GetTxReceipts(height int64) (*types.TxReceipts4SubscribePerBlk, error) {
	key := txReceiptsKey4Height(height)
	value, err := syncTx.db.Get(key)
	if err != nil {
		return nil, err
	}
	detail := &types.TxReceipts4SubscribePerBlk{}
	err = types.Decode(value, detail)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (syncTx *SyncTxReceipts) GetNextValidTxReceipts(height int64) (*types.TxReceipts4SubscribePerBlk, error) {
	key := txReceiptsKey4Height(height)
	helper := dbm.NewListHelper(syncTx.db)
	TxReceipts := helper.List(txReceiptPrefix, key, 1, dbm.ListASC)
	if nil == TxReceipts {
		return nil, nil
	}
	detail := &types.TxReceipts4SubscribePerBlk{}
	err := types.Decode(TxReceipts[0], detail)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (syncTx *SyncTxReceipts) delTxReceipts(height int64) {
	key := txReceiptsKey4Height(height)
	_ = syncTx.db.Set(key, nil)
}

// 检查输入是否有问题, 并解析输入
func parseTxReceipts(txReceipts *types.TxReceipts4Subscribe) (count int, start int64, txsWithReceipt []*types.TxReceipts4SubscribePerBlk, err error) {
	count = len(txReceipts.TxReceipts)
	txsWithReceipt = make([]*types.TxReceipts4SubscribePerBlk, 0)
	start = math.MaxInt64
	for i := 0; i < count; i++ {
		if txReceipts.TxReceipts[i].AddDelType != SeqTypeAdd && txReceipts.TxReceipts[i].AddDelType != SeqTypeDel {
			log.Error("parseTxReceipts seq op not support", "seq", txReceipts.TxReceipts[i].SeqNum,
				"height", txReceipts.TxReceipts[i].Height, "seqOp", txReceipts.TxReceipts[i].AddDelType)
			continue
		}
		txsWithReceipt = append(txsWithReceipt, txReceipts.TxReceipts[i])
		if txReceipts.TxReceipts[i].SeqNum < start {
			start = txReceipts.TxReceipts[i].SeqNum
		}
		log.Debug("parseTxReceipts get one block's tx with receipts", "seq", txReceipts.TxReceipts[i].SeqNum,
			"height", txReceipts.TxReceipts[i].Height, "seqOpType", seqOperationType[txReceipts.TxReceipts[i].AddDelType - 1])

	}
	if len(txsWithReceipt) != count {
		err = errors.New("duplicate block's tx receipt")
		return
	}
	return
}
