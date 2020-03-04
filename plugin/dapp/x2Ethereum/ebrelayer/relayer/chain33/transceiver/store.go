package blockchain

import (
	"fmt"
	"sync/atomic"

	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	"github.com/pkg/errors"
	"gitlab.33.cn/chain33/lns/lightningservice/blockchain/syncblock"
	lnsTypes "gitlab.33.cn/chain33/lns/lightningservice/types"
)

var (
	onchainTxs4addrPrefix = []byte("onchainTxs4addr:")
	lastSyncHeightPrefix  = []byte("lastSyncHeight:")
	transferPrefix        = []byte("transfer:")
	transferConfirmPrefix = []byte("transferConfirm:")
	withdrawPrefix        = []byte("withdraw:")
	withdrawConfirmPrefix = []byte("withdrawConfirm:")
	syncInfo4AppPrefix    = []byte("syncInfo4App:")
	lns2walletMsgPrefix   = []byte("lns2walletMsg:")
)

const (
	Invalid_Store_ID = 0
)

func blockKey4Height(addr string, height int64) []byte {
	return append(onchainTxs4addrPrefix, []byte(fmt.Sprintf("%s:%012d", height, addr))...)
}

func key4transfer(receiver, paymentID string) []byte {
	return append(transferPrefix, []byte(fmt.Sprintf("%s:%s", receiver, paymentID))...)
}

func key4transConfirm(receiver, paymentID string) []byte {
	return append(transferConfirmPrefix, []byte(fmt.Sprintf("%s:%s", receiver, paymentID))...)
}

func key4withdraw(receiver, withdrawID string) []byte {
	return append(withdrawPrefix, []byte(fmt.Sprintf("%s:%s", receiver, withdrawID))...)
}

func key4withdrawConfirm(receiver, withdrawID string) []byte {
	return append(withdrawConfirmPrefix, []byte(fmt.Sprintf("%s:%s", receiver, withdrawID))...)
}

func key4SyncInfo(addr string) []byte {
	return append(syncInfo4AppPrefix, []byte(fmt.Sprintf("%s", addr))...)
}

func key4StoreID(addr string, storeID uint64) []byte {
	return append(lns2walletMsgPrefix, []byte(fmt.Sprintf("%s:%024d", addr, storeID))...)
}

func (chain *BlockChainSerivce) GetNextStoreID4User(addr string) (uint64, error) {
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	syncInfo4App, ok := chain.addr2syncInfo[addr]
	if !ok {
		info := fmt.Sprintf("GetNextStoreID4User: Cant't get syncInfo4App for addr:%s", addr)
		panic(info)
	}
	return atomic.AddUint64(&syncInfo4App.MaxStoreID, 1), nil
}

func (chain *BlockChainSerivce) decreaseStoreID4User(addr string) (uint64, error) {
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	syncInfo4App, ok := chain.addr2syncInfo[addr]
	if !ok {
		return Invalid_Store_ID, errors.New("The Receiver can't be reached now,as it's not registered")
	}
	return atomic.AddUint64(&syncInfo4App.MaxStoreID, ^uint64(0)), nil
}

//StoreLns2WalletMsg:生成一个存储ID,并将其填充并保存到数据库中
func (chain *BlockChainSerivce) StoreLns2WalletMsg(addr string, msg *lnsTypes.LnsWalletMsg) error {
	storeID, err := chain.GetNextStoreID4User(addr)
	if err != nil {
		return err
	}
	msg.StoreID = storeID
	data := types.Encode(msg)
	if err := chain.db.Set(key4StoreID(addr, storeID), data); err != nil {
		//这里的回退似乎没有什么意义
		chain.decreaseStoreID4User(addr)
		panic("Failed to StoreLns2WalletMsg due to fail to write db")
		return err
	}
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	syncInfo4App, ok := chain.addr2syncInfo[addr]
	if !ok {
		info := fmt.Sprintf("StoreLns2WalletMsg: Cant't get syncInfo4App for addr:%s", addr)
		panic(info)
	}
	//将信息MaxStoreID更新到数据库
	chain.storeSyncInfo4User(syncInfo4App)
	return nil
}

func (chain *BlockChainSerivce) getLns2WalletMsg(addr string, storeID uint64) (*lnsTypes.LnsWalletMsg, error) {
	data, err := chain.db.Get(key4StoreID(addr, storeID))
	if err != nil {
		return nil, err
	}

	msg := &lnsTypes.LnsWalletMsg{}
	err = types.Decode(data, msg)
	return msg, err
}

func (chain *BlockChainSerivce) storeOnchainTxs4addr(receiver string, TxsOnChain *lnsTypes.TxsOnChain) error {
	msg := lnsTypes.LnsWalletMsg{
		Type: lnsTypes.MsgTxsOnChain,
		Value: &lnsTypes.LnsWalletMsg_TxsOnChain{
			TxsOnChain: TxsOnChain,
		},
		StoreID: Invalid_Store_ID,
		Receiver: receiver,
	}
	return chain.StoreLns2WalletMsg(receiver, &msg)
}

//获取上次同步到app的高度
func (chain *BlockChainSerivce) loadLastSyncHeight() (int64, error) {
	return chain.LoadInt64FromDB(lastSyncHeightPrefix, chain.db)
}

func (chain *BlockChainSerivce) setLastSyncHeight(syncHeight int64) {
	bytes := types.Encode(&types.Int64{Data: syncHeight})
	chain.db.Set(lastSyncHeightPrefix, bytes)
}

//func (chain *BlockChainSerivce) StoreTransfer(in lnsTypes.TransferSignedState) {
//	bytes := types.Encode(&in)
//	chain.db.Set(key4transfer(in.Receiver, in.PaymentID), bytes)
//}

func (chain *BlockChainSerivce) getTransfer(receiver, paymentID string) (tranfer lnsTypes.TransferSignedState, err error) {
	bytes, err := chain.db.Get(key4transfer(receiver, paymentID))
	if nil != err {
		return
	}
	err = types.Decode(bytes, &tranfer)
	return
}

func (chain *BlockChainSerivce) StoreTransferConfirm(in lnsTypes.TransferConfirm) {
	bytes := types.Encode(&in)
	chain.db.Set(key4transConfirm(in.TransferPaymentID.Payer, in.TransferPaymentID.PaymentID), bytes)
}

func (chain *BlockChainSerivce) getTransferConfirm(receiver, paymentID string) (tranferConfirm lnsTypes.TransferConfirm, err error) {
	bytes, err := chain.db.Get(key4transConfirm(receiver, paymentID))
	if nil != err {
		return
	}
	err = types.Decode(bytes, &tranferConfirm)
	return
}

func (chain *BlockChainSerivce) StoreWithdrawAction(in lnsTypes.WithdrawAction) {
	bytes := types.Encode(&in)
	chain.db.Set(key4withdraw(in.Partner, in.WithdrawID), bytes)
}

func (chain *BlockChainSerivce) getWithdrawAction(receiver, withdrawID string) (withdrawAction lnsTypes.WithdrawAction, err error) {
	bytes, err := chain.db.Get(key4withdraw(receiver, withdrawID))
	if nil != err {
		return
	}
	err = types.Decode(bytes, &withdrawAction)
	return
}

func (chain *BlockChainSerivce) StoreWithdrawConfirm(in lnsTypes.WithdrawConfirmReply) {
	bytes := types.Encode(&in)
	chain.db.Set(key4withdrawConfirm(in.WithdrawConfirmProof.Withdrawer, in.WithdrawID), bytes)
}

func (chain *BlockChainSerivce) getWithdrawConfirm(receiver, withdrawID string) (withdrawConfirmReply lnsTypes.WithdrawConfirmReply, err error) {
	bytes, err := chain.db.Get(key4withdrawConfirm(receiver, withdrawID))
	if nil != err {
		return
	}
	err = types.Decode(bytes, &withdrawConfirmReply)
	return
}

func (chain *BlockChainSerivce) storeSyncInfo4User(in *lnsTypes.SyncInfo4App) {
	bytes := types.Encode(in)
	chain.db.Set(key4SyncInfo(in.Address), bytes)
}

func (chain *BlockChainSerivce) getSyncInfo4User(addr string) (in lnsTypes.SyncInfo4App, err error) {
	bytes, err := chain.db.Get(key4SyncInfo(addr))
	if nil != err {
		return
	}
	err = types.Decode(bytes, &in)
	return
}

func (chain *BlockChainSerivce) loadAllSyncInfo4User() (err error) {
	kvdb := dbm.NewKVDB(chain.db)
	//TODO:后续用户量多的情况，考虑批量获取，而不是一次性，同时也可以考虑将addr2syncInfo通过LRU的方式进行总量的管理
	datas, err := kvdb.List(syncInfo4AppPrefix, []byte{}, 0, 0)
	if nil != err {
		//对于首次加载失败的情况，不作处理
		chainlog.Info("Load nothing of SyncInfo for users")
		return nil
	}

	chain.rw4url.Lock()
	defer chain.rw4url.Unlock()
	for _, data := range datas {
		syncInfo4App := lnsTypes.SyncInfo4App{}
		if err := types.Decode(data, &syncInfo4App); err != nil {
			return err
		}
		chain.addr2syncInfo[syncInfo4App.Address] = &syncInfo4App
		if syncInfo4App.MaxStoreID > syncInfo4App.LastSyncStoreID && NO == syncInfo4App.ContinueSendFail {
			chain.rwRepush.Lock()
			chain.addr2repush[syncInfo4App.Address] = syncInfo4App.MaxStoreID - syncInfo4App.LastSyncStoreID
			chain.rwRepush.Unlock()
		}
	}
	return nil
}

//为了保证每个移动端用户数据的完整、及时地推送给用户，为每个用户每项持久化的数据分配一个存储store_id，
// 这样就可以保存当前的用户信息成功地推送到哪个存储项了
//addr-store_id<----->key<----->item
//addr-lastStoreID<------>store_id,获取最近推送成功的storeid
