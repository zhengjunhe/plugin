// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blockchain

//message callback
import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	chain33Types "github.com/33cn/chain33/types"
	dappTypes "github.com/33cn/plugin/plugin/dapp/lns/types"
	"github.com/pkg/errors"
	"gitlab.33.cn/chain33/lns/lightningservice/forwarder"
	lnsTypes "gitlab.33.cn/chain33/lns/lightningservice/types"
)

const (
	Send_OnGoing = int32(1)
	Send_Free    = int32(0)

	YES = int32(1)
	NO  = int32(0)

	SyncBlockChanSize = 10
	BufferSize        = 1000
	BatchSizeTimeout  = 100
	Timeout           = (100 * time.Millisecond)
)

//blockProc: 定期地从区块链节点拉取最新块，并进行处理
func (chain *BlockChainSerivce) blockProc() {
	defer chain.wg.Done()
	for {
		select {
		case currentHeight := <-chain.syncBlockChan:
			chain.onNewBlockProc(currentHeight)
			//case <-chain.quit:
			//	return
		case <-chain.ctx.Done():
			return
		}
	}
}

//sendMsg2App:并发地转发每个用户的消息到app，同时，每个用户的消息的发送又是顺序发送的
//----> addr1:storeID_1
//----> addr2:storeID_1
//----> addr1:storeID_2
//----> addr1:storeID_3
//----> addr3:storeID_1
func (chain *BlockChainSerivce) sendMsg2App() {
	defer chain.wg.Done()
	send2AppConcurrnet := atomic.LoadInt32(&chain.send2AppConcurrnet)
	var resourceChannel = make(chan bool, send2AppConcurrnet)
	timer := time.NewTimer(Timeout)
	for {
		select {
		//TODO:此处需要再进行优化，提高转发能力，消息接收通道能力和并行转发能力如何更加步调一致地进行协作
		case task := <-chain.pushTasksChan:
			chain.rw4url.RLock()
			syncInfo, ok := chain.addr2syncInfo[task.Receiver]
			chain.rw4url.RUnlock()
			if !ok {
				chainlog.Error("sendNextMsg2app", "Can't get sync info for address:", task.Receiver)
				continue
			}

			//目前给同一个用户推送数据时，只支持串行方式
			if Send_OnGoing == atomic.LoadInt32(&syncInfo.OngoingSend) {
				chain.pendingSends.PushBack(task)
				chainlog.Info("sendNextMsg2app", "Send msg to app is ongoing:", task.Receiver)
				continue
			}
			atomic.CompareAndSwapInt32(&syncInfo.OngoingSend, Send_Free, Send_OnGoing)

			resourceChannel <- true
			//通过resourceChannel控制链下消息转发并发数
			go func(chain *BlockChainSerivce, task PushTask) {
				defer atomic.CompareAndSwapInt32(&syncInfo.OngoingSend, Send_OnGoing, Send_Free)
				chain.sendNextMsg2app(task.Receiver)
				<-resourceChannel
			}(chain, task)

		case <-timer.C:
			if chain.pendingSends.Len() != 0 {
				for i := 0; i < BatchSizeTimeout; i++ {
					task := chain.pendingSends.Front()
					if nil == task {
						break
					}
					chain.pushTasksChan <- task.Value.(PushTask)
					chain.pendingSends.Remove(task)
				}
			}
			timer.Reset(Timeout)

		case <-chain.ctx.Done():
			return
		}
	}
}

//用于处理移动端重新推送的请求
func (chain *BlockChainSerivce) repushMsg2App() {
	defer chain.wg.Done()
	duration := 50 * time.Millisecond
	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C:
			chain.rwRepush.Lock()
			for addr, cnt := range chain.addr2repush {
				chain.pushTasksChan <- PushTask{addr}
				if 1 == cnt {
					delete(chain.addr2repush, addr)
					continue
				}
				chain.addr2repush[addr] = cnt - 1
			}
			chain.rwRepush.Unlock()
			timer.Reset(duration)
		case <-chain.ctx.Done():
			return
		}
	}
}

func (chain *BlockChainSerivce) sendNextMsg2app(addr string) {
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	syncInfo, _ := chain.addr2syncInfo[addr]
	if syncInfo.Url == "" {
		chainlog.Info("sendNextMsg2app", "No URL is setted for address:", addr)
		return
	}
	storeID2Send := atomic.LoadUint64(&syncInfo.LastSyncStoreID) + 1
	if storeID2Send > atomic.LoadUint64(&syncInfo.MaxStoreID) {
		chainlog.Error("sendNextMsg2app", "No fresh data need to be sent to app(maybe caused by restart push by app) for addr:", addr)
		return
	}

	if YES == atomic.LoadInt32(&syncInfo.ContinueSendFail) {
		chainlog.Info("sendNextMsg2app", "Not going to push data to user due to failed 3 times continuesly for addr:", addr)
		return
	}

	data2Send, err := chain.getLns2WalletMsg(addr, storeID2Send)
	if err != nil {
		chainlog.Error("sendNextMsg2app", "Failed to getLns2WalletMsg for addr:", addr, "storeID2Send:", storeID2Send)
		return
	}
	if err = forwarder.Push2app(data2Send, syncInfo.Url); err != nil {
		i := 1
		for ; i < 3; i++ {
			time.Sleep(time.Second)
			if nil == forwarder.Push2app(data2Send, syncInfo.Url) {
				break
			}
		}
		//尝试3次推送失败，设置连续推送失败标志，在移动端重启推送标志前，后续将不再推送
		if 3 == i {
			chainlog.Error("sendNextMsg2app", "Failed to push data in 3 times for addr:", syncInfo.Address,
				"URL:", syncInfo.Url)
			atomic.CompareAndSwapInt32(&syncInfo.OngoingSend, Send_OnGoing, Send_Free)
			atomic.CompareAndSwapInt32(&syncInfo.ContinueSendFail, NO, YES)
			chain.storeSyncInfo4User(syncInfo)
			return
		}
	}

	atomic.AddUint64(&syncInfo.LastSyncStoreID, 1)
	chain.storeSyncInfo4User(syncInfo)
	return
}

func (chain *BlockChainSerivce) getUrl(receiver string) (string, error) {
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	syncInfo, ok := chain.addr2syncInfo[receiver]
	if !ok {
		chainlog.Error("getUrl", "Receiver is unregistered for address:", receiver)
		return "", errors.New("Receiver is unregistered")
	}
	return syncInfo.Url, nil
}

func (chain *BlockChainSerivce) onNewBlockProc(currentHeight int64) {

	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	//    heightSync2App    matDegress   currentHeight
	syncHeight := chain.heightSync2App
	chainlog.Debug("onNewBlockProc", "currHeight", currentHeight, "syncHeight", syncHeight)
	for syncHeight+int64(chain.matDegree) <= currentHeight {
		blockDetail, err := chain.syncBlock.GetBlock(syncHeight)
		if err != nil {
			chainlog.Error("SyncBlock dealBlocks", "Failed to GetBlock for height:", syncHeight, "err", err)
			return
		}

		block := blockDetail.Block
		//step1:推送新块高度信息给dapp移动端，用于withdraw等的超时处理
		chain.pushNewHeight2App(block)

		//step2:准备处理和推送移动用户相关的交易tx
		txs := block.Txs
		txsonChain := lnsTypes.TxsOnChain{
			Height:    block.Height,
			BlockHash: block.MainHash,
			TxDetails: make([]*lnsTypes.TransactionReceipt, 0),
		}

		for i, tx := range txs {
			//检查是否为lns的交易(包括平行链：user.p.xxx.lns)，将闪电网络交易进行收集
			if 0 == bytes.Compare(tx.Execer, []byte("lns")) ||
				(len(tx.Execer) > 4 && string(tx.Execer[(len(tx.Execer)-4):]) == ".lns") {

				chainlog.Debug("SyncLnsTx", "exec", string(tx.Execer), "action", tx.ActionName(), "fromAddr", tx.From())
				txDetail := &lnsTypes.TransactionReceipt{
					Tx:          tx,
					ReceiptData: blockDetail.Receipts[i],
					Height:      block.Height,
				}
				txsonChain.TxDetails = append(txsonChain.TxDetails, txDetail)
			}
		}
		chainlog.Debug("SyncLnsTxSummary", "syncHeight", syncHeight, "lnsTxNum", len(txsonChain.TxDetails))
		//若存在闪电网络交易，则将其推送到订阅者用户中
		//TODO：更细致的实现需要考虑推送给部分用户成功了，就不再推送，只有不成功的用户才继续推送
		if len(txsonChain.TxDetails) > 0 {
			chain.filterAndPushTxs4Users(&txsonChain)
		}

		syncHeight++
		chain.setLastSyncHeight(syncHeight)
		chain.heightSync2App = syncHeight
	}
}

func (chain *BlockChainSerivce) pushNewHeight2App(block *chain33Types.Block) {
	chain.rw4url.RLock()
	defer chain.rw4url.RUnlock()
	//TODO：
	newBlock := &lnsTypes.NewBlock{
		Height:     block.Height,
		BlockHash:  block.MainHash,
		ParentHash: block.ParentHash,
	}
	//高度消息推送给所有用户，所以receiver字段在for循环中每次都要更新
	msg := lnsTypes.LnsWalletMsg{
		Type: lnsTypes.MsgNewBlock,
		Value: &lnsTypes.LnsWalletMsg_NewBlock{
			NewBlock: newBlock,
		},
	}
	concurrnet := atomic.LoadInt32(&chain.newHeightConcurrnet)
	var resourceChannel = make(chan bool, concurrnet)
	for _, syncInfo := range chain.addr2syncInfo {
		url := syncInfo.Url
		if url == "" {
			continue
		}
		//添加msg的receiver
		msg.Receiver = syncInfo.Address
		resourceChannel <- true
		go func(walletMsg lnsTypes.LnsWalletMsg, url string) {
			if err := forwarder.Push2app(&walletMsg, url); nil != err {
				chainlog.Error("pushNewHeight2App", "push error", err)
			}
			<-resourceChannel
		}(msg, url)
	}
}

func (chain *BlockChainSerivce) filterAndPushTxs4Users(txsOnChain *lnsTypes.TxsOnChain) {

	addr2txs := make(map[string]*lnsTypes.TxsOnChain)
	for _, txDetail := range txsOnChain.TxDetails {
		action := dappTypes.LnsAction{}
		err := chain33Types.Decode(txDetail.Tx.GetPayload(), &action)
		if err != nil {
			chainlog.Error("filterAndPushTxs4Users", "decode error", err)
			return
		}
		switch action.Ty {
		case dappTypes.TyOpenAction:
			chain.procOpen(txDetail, addr2txs)

		case dappTypes.TyDepositAction:
			chain.procDeposit(txDetail, addr2txs)

		case dappTypes.TyWithdrawAction:
			chain.procWithdraw(txDetail, addr2txs)

		case dappTypes.TyCloseAction:
			chain.procClose(txDetail, addr2txs)

		case dappTypes.TyUpdateProofAction:
			chain.procUpdate(txDetail, addr2txs)

		case dappTypes.TySettleAction:
			chain.procSettle(txDetail, addr2txs)

		default:
			chainlog.Error("filterAndPushTxs4Users", "decode error", err)
		}

	}
	for addr, txonChain4Push := range addr2txs {
		txonChain4Push.Height = txsOnChain.Height
		txonChain4Push.BlockHash = txsOnChain.BlockHash

		//在将tx推送到钱包侧之前，先进行持久化操作
		chain.storeOnchainTxs4addr(addr, txonChain4Push)
		chain.pushTasksChan <- PushTask{addr}
	}
	return
}

func (chain *BlockChainSerivce) addTx4Users(participants []string, txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, addr := range participants {
		if _, ok := chain.addr2syncInfo[addr]; !ok {
			//当用户未注册时，也将其相应的数据进行保存，但是在未注册其接收数据的URL之前不进行推送
			info := &lnsTypes.SyncInfo4App{
				Url:              "",
				MaxStoreID:       Invalid_Store_ID,
				LastSyncStoreID:  Invalid_Store_ID,
				OngoingSend:      Send_Free,
				Address:          addr,
				ContinueSendFail: NO,
			}
			chain.rw4url.Lock()
			chain.addr2syncInfo[addr] = info
			chain.rw4url.Unlock()
			chain.storeSyncInfo4User(info)
			chainlog.Info("addTx4Users", "Add new user with addr:", addr)
		}
		if addr2txs[addr] == nil {
			addr2txs[addr] = new(lnsTypes.TxsOnChain)
		}
		addr2txs[addr].TxDetails = append(addr2txs[addr].TxDetails, txDetail)
	}
}

func (chain *BlockChainSerivce) procOpen(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyOpenLog == log.Ty {
			open := &dappTypes.ReceiptOpen{}
			err := chain33Types.Decode(log.Log, open)
			if nil != err {
				chainlog.Error("procOpen", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{open.Opener, open.Partner}, txDetail, addr2txs)
			return
		}
	}
	return
}

func (chain *BlockChainSerivce) procDeposit(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyDepositLog == log.Ty {
			deposit := &dappTypes.ReceiptDeposit{}
			err := chain33Types.Decode(log.Log, deposit)
			if nil != err {
				chainlog.Error("procDeposit", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{deposit.Partner, deposit.Depositor}, txDetail, addr2txs)
			return
		}
	}
	return
}

func (chain *BlockChainSerivce) procWithdraw(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyDepositLog == log.Ty {
			withdraw := &dappTypes.ReceiptWithdraw{}
			err := chain33Types.Decode(log.Log, withdraw)
			if nil != err {
				chainlog.Error("procWithdraw", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{withdraw.Withdrawer, withdraw.Partner}, txDetail, addr2txs)
			return
		}
	}
	return
}

func (chain *BlockChainSerivce) procClose(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyCloseLog == log.Ty {
			close := &dappTypes.ReceiptClose{}
			err := chain33Types.Decode(log.Log, close)
			if nil != err {
				chainlog.Error("procClose", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{close.Closer, close.Partner}, txDetail, addr2txs)
			return
		}
	}
	return
}

func (chain *BlockChainSerivce) procUpdate(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyCloseLog == log.Ty {
			update := &dappTypes.ReceiptUpdate{}
			err := chain33Types.Decode(log.Log, update)
			if nil != err {
				chainlog.Error("procUpdate", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{update.Updater, update.Partner}, txDetail, addr2txs)
			return
		}
	}
	return
}

func (chain *BlockChainSerivce) procSettle(txDetail *lnsTypes.TransactionReceipt, addr2txs map[string]*lnsTypes.TxsOnChain) {
	for _, log := range txDetail.ReceiptData.Logs {
		if dappTypes.TyCloseLog == log.Ty {
			settle := &dappTypes.ReceiptSettle{}
			err := chain33Types.Decode(log.Log, settle)
			if nil != err {
				chainlog.Error("procSettle", "wrong lns tx with corrupted payload with txhash",
					common.ToHex(txDetail.Tx.Hash()), "err info", err)
				return
			}
			chain.addTx4Users([]string{settle.Participant1, settle.Participant2}, txDetail, addr2txs)
			return
		}
	}
	return
}

//func (chain *BlockChainSerivce) IsReachable(Receiver string) bool {
//	chain.rw4url.RLock()
//	defer chain.rw4url.Unlock()
//	if _, ok := chain.addr2syncInfo[Receiver]; !ok {
//		return false
//	}
//	return true
//}

func (chain *BlockChainSerivce) GetPushTasksChan() chan<- PushTask {
	return chain.pushTasksChan
}

//Register:用来为移动用户第一次注册URL或者后续更新URL
func (chain *BlockChainSerivce) Register(req lnsTypes.RegisterUserReq) (bool, error) {
	var RegisterUserData lnsTypes.RegisterUserReq
	RegisterUserData = req
	RegisterUserData.Signature = nil
	hash := common.Sha256(chain33Types.Encode(&RegisterUserData))
	if req.GetSignature() == nil || !chain33Types.CheckSign(hash, "", req.Signature) {
		return false, errors.New("Wrong Signature")
	}
	if req.Address != address.PubKeyToAddr(req.GetSignature().GetPubkey()) {
		return false, errors.New("The address is not consistent with the signature")
	}

	chain.rw4url.Lock()
	defer chain.rw4url.Unlock()
	//如果该地址的信息已经存在，就直接更新url
	if info, ok := chain.addr2syncInfo[req.Address]; ok {
		info.Url = req.Url
		chain.storeSyncInfo4User(info)
		//debug
		fmt.Println("registered addr:", chain.addr2syncInfo)
		//debug
		return true, nil
	}

	info := &lnsTypes.SyncInfo4App{
		Url:              req.Url,
		MaxStoreID:       Invalid_Store_ID,
		LastSyncStoreID:  Invalid_Store_ID,
		OngoingSend:      Send_Free,
		Address:          req.Address,
		ContinueSendFail: NO,
	}
	chain.addr2syncInfo[req.Address] = info
	chain.storeSyncInfo4User(info)
	//debug
	fmt.Println("registered addr:", chain.addr2syncInfo)
	//debug
	return true, nil
}

//RestartPush：响应移动端重新推送消息的请求
func (chain *BlockChainSerivce) RestartPush(req lnsTypes.RestartPushReq) (bool, error) {
	var RestartPushdata lnsTypes.RestartPushReq
	RestartPushdata = req
	RestartPushdata.Signature = nil
	hash := common.Sha256(chain33Types.Encode(&RestartPushdata))
	if !chain33Types.CheckSign(hash, "", req.Signature) {
		return false, errors.New("Wrong Signature")
	}
	if req.Address != address.PubKeyToAddr(req.GetSignature().GetPubkey()) {
		return false, errors.New("The address is not consistent with the signature")
	}
	//TODO 可以封装成函数, 统一控制锁
	chain.rw4url.RLock()
	syncInfo, ok := chain.addr2syncInfo[req.Address]
	chain.rw4url.RUnlock()
	if !ok {
		chainlog.Error("RestartPush", "No syncInfo for addr:", req.Address)
		return false, errors.New("No data saved in lns33 now,pls check it yourself first")
	}

	//如果未出现连续三次发送失败的情况，lns应该一直推送,需要移动端先关闭后重启才能允许
	//TODO:这部分是否允许重新启动推送的判断机制后续再进行讨论修改
	if YES != atomic.LoadInt32(&syncInfo.ContinueSendFail) {
		chainlog.Info("RestartPush", "Pushing is always ongoing for addr:", req.Address)
		return true, nil
	}
	//重置推送失败标志
	atomic.CompareAndSwapInt32(&syncInfo.ContinueSendFail, YES, NO)
	chain.storeSyncInfo4User(syncInfo)

	lastSyncStoreID := atomic.LoadUint64(&syncInfo.LastSyncStoreID)
	if req.StoreID > lastSyncStoreID+1 {
		chainlog.Error("RestartPush", "addr:", req.Address, "req.StoreID:", req.StoreID,
			"but lastSyncStoreID:", lastSyncStoreID)
		return false, errors.New("Request restart store id is greater than last SyncStoreID")
	}

	chain.rwRepush.Lock()
	defer chain.rwRepush.Unlock()
	if _, ok := chain.addr2repush[req.Address]; !ok {
		chainlog.Error("RestartPush", "Last repush is not finished for addr:", req.Address)
		return false, errors.New("Last repush is not finished")
	}
	//---------------------req.StoreID-----------lastSyncStoreID-------MaxStoreID
	pushCnt := atomic.LoadUint64(&syncInfo.MaxStoreID) - req.StoreID + 1
	//将lastSyncStoreID的值直接更新为重启push请求中的值
	atomic.CompareAndSwapUint64(&syncInfo.LastSyncStoreID, lastSyncStoreID, req.StoreID)
	chain.storeSyncInfo4User(syncInfo)
	chain.addr2repush[req.Address] = pushCnt

	return true, nil
}
