package chain33

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	chain33Types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	syncTx "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/chain33/transceiver/sync"
	relayerTx "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/txs"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var relayerLog = log.New("module", "chain33_relayer")

type Chain33Relayer struct {
	syncTxReceipts      *syncTx.SyncTxReceipts
	web3Provider        string
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	fetchHeightPeriodMs int64
	db                  dbm.DB
	syncTxChan          <-chan int64
	//height              int64 //当前区块高度  +++++++++||++++++++++++||++++++++++||
	//heightSync2App      int64 //已经同步高度           ^             ^           ^
	lastHeight4Tx        int64 //等待被处理的具有相应的交易回执的高度
	matDegree            int32 //成熟度         heightSync2App    matDegress   height
	passphase            string
	privateKey4Ethereum  *ecdsa.PrivateKey
	ethSender            ethCommon.Address
	contractAddress      ethCommon.Address
	totalTx4Chain33ToEth int64
	ctx                  context.Context
	wg                   sync.WaitGroup
	rwLock               sync.RWMutex
	unlock               chan int
	status               int32
}

// StartChain33Relayer : initializes a relayer which witnesses events on the chain33 network and relays them to Ethereum
func StartChain33Relayer(syncTxConfig *ebTypes.SyncTxConfig, db dbm.DB, ctx context.Context) *Chain33Relayer {
	relayer := &Chain33Relayer{
		rpcLaddr:            syncTxConfig.Chain33Host,
		fetchHeightPeriodMs: syncTxConfig.FetchHeightPeriodMs,
		unlock:              make(chan int),
		db:                  db,
		ctx:                 ctx,
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		Chain33Host:       syncTxConfig.Chain33Host,
		PushHost:          syncTxConfig.PushHost,
		PushName:          syncTxConfig.PushName,
		PushBind:          syncTxConfig.PushBind,
		StartSyncHeight:   syncTxConfig.StartSyncHeight,
		StartSyncSequence: syncTxConfig.StartSyncSequence,
		StartSyncHash:     syncTxConfig.StartSyncHash,
	}

	go relayer.syncProc(syncCfg)
	return relayer
}

func (chain33Relayer *Chain33Relayer) SetPassphase(passphase string) {
	chain33Relayer.rwLock.Lock()
	chain33Relayer.passphase = passphase
	chain33Relayer.rwLock.Unlock()
}

func (chain33Relayer *Chain33Relayer) QueryTxhashRelay2Eth() ebTypes.Txhashes {
	txhashs := utils.QueryTxhashes([]byte(chain33ToEthBurnLockTxHashPrefix), chain33Relayer.db)
	return ebTypes.Txhashes{Txhash: txhashs}
}

func (chain33Relayer *Chain33Relayer) GetRunningStatus() (relayerRunStatus *ebTypes.RelayerRunStatus) {
	relayerRunStatus = &ebTypes.RelayerRunStatus{}
	chain33Relayer.rwLock.RLock()
	relayerRunStatus.Status = chain33Relayer.status
	chain33Relayer.rwLock.RUnlock()
	if relayerRunStatus.Status == ebTypes.StatusPending {
		if nil == chain33Relayer.privateKey4Ethereum {
			relayerRunStatus.Details = "Ethereum's private key not imported"
		}
		return
	}
	relayerRunStatus.Details = "Running"
	return
}

func (chain33Relayer *Chain33Relayer) syncProc(syncCfg *ebTypes.SyncTxReceiptConfig) {
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Chain33 relayer")
	<-chain33Relayer.unlock
	_, _ = fmt.Fprintln(os.Stdout, "Chain33 relayer starts to run...")

	syncChan := make(chan int64, 10)
	chain33Relayer.syncTxReceipts = syncTx.StartSyncTxReceipt(syncCfg, syncChan, chain33Relayer.db)
	chain33Relayer.lastHeight4Tx = chain33Relayer.loadLastSyncHeight()

	timer := time.NewTicker(time.Duration(chain33Relayer.fetchHeightPeriodMs) * time.Millisecond)
	for {
		select {
		case <-timer.C:
			height := chain33Relayer.getCurrentHeight()
			relayerLog.Debug("syncProc", "getCurrentHeight", height)
			chain33Relayer.onNewHeightProc(height)

		case <-chain33Relayer.ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (chain33Relayer *Chain33Relayer) getCurrentHeight() int64 {
	var res rpctypes.Header
	ctx := jsonclient.NewRPCCtx(chain33Relayer.rpcLaddr, "Chain33.GetLastHeader", nil, &res)
	_, err := ctx.RunResult()
	if nil != err {
		relayerLog.Error("getCurrentHeight", "Failede due to:", err.Error())
	}
	return res.Height
}

func (chain33Relayer *Chain33Relayer) onNewHeightProc(currentHeight int64) {
	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	// lastHeight4Tx    matDegress   currentHeight
	for chain33Relayer.lastHeight4Tx+int64(chain33Relayer.matDegree)+1 <= currentHeight {
		relayerLog.Debug("onNewHeightProc", "currHeight", currentHeight, "lastHeight4Tx", chain33Relayer.lastHeight4Tx)

		lastHeight4Tx := chain33Relayer.lastHeight4Tx
		TxReceipts, err := chain33Relayer.syncTxReceipts.GetNextValidTxReceipts(lastHeight4Tx)
		if nil == TxReceipts || nil != err {
			if err != nil {
				relayerLog.Error("onNewHeightProc", "Failede to GetNextValidTxReceipts due to:", err.Error())
			}
			break
		}
		relayerLog.Debug("onNewHeightProc", "currHeight", currentHeight, "valid tx receipt with height:", TxReceipts.Height)

		txs := TxReceipts.Tx
		for i, tx := range txs {
			//检查是否为lns的交易(包括平行链：user.p.xxx.lns)，将闪电网络交易进行收集
			if 0 != bytes.Compare(tx.Execer, []byte(relayerTx.X2Eth)) &&
				(len(tx.Execer) > 4 && string(tx.Execer[(len(tx.Execer)-4):]) != "."+relayerTx.X2Eth) {
				relayerLog.Debug("onNewHeightProc, the tx is not x2ethereum", "Execer", string(tx.Execer), "height:", TxReceipts.Height)
				continue
			}
			var ss types.X2EthereumAction
			_ = chain33Types.Decode(tx.Payload, &ss)
			relayerLog.Debug("onNewHeightProc", "exec", string(tx.Execer), "tx", ss.GetActionName(), "action", tx.ActionName(), "fromAddr", tx.From(), "Name", tx)
			actionName := ss.GetActionName()
			if relayerTx.BurnAction == actionName || relayerTx.LockAction == actionName {
				actionEvent := getOracleClaimType(actionName)
				if err := chain33Relayer.handleBurnLockMsg(actionEvent, TxReceipts.ReceiptData[i]); nil != err {
					errInfo := fmt.Sprintf("Failed to handleBurnLockMsg due to:", err.Error())
					panic(errInfo)
				}
			}
		}
		chain33Relayer.lastHeight4Tx = TxReceipts.Height
		chain33Relayer.setLastSyncHeight(chain33Relayer.lastHeight4Tx)
	}
}

// getOracleClaimType : sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) events.Event {
	var claimType events.Event

	switch eventType {
	case events.MsgBurn.String():
		claimType = events.MsgBurn
	case events.MsgLock.String():
		claimType = events.MsgLock
	default:
		claimType = events.Unsupported
	}

	return claimType
}

// handleBurnLockMsg : parse event data as a Chain33Msg, package it into a ProphecyClaim, then relay tx to the Ethereum Network
func (chain33Relayer *Chain33Relayer) handleBurnLockMsg(claimEvent events.Event, receipt *chain33Types.ReceiptData) error {
	// Parse the witnessed event's data into a new Chain33Msg
	chain33Msg := relayerTx.BurnLockTxReceiptToChain33Msg(claimEvent, receipt)

	// Parse the Chain33Msg into a ProphecyClaim for relay to Ethereum
	prophecyClaim := relayerTx.Chain33MsgToProphecyClaim(chain33Msg)

	// TODO: Need some sort of delay on this so validators aren't all submitting at the same time
	// Relay the Chain33Msg to the Ethereum network
	txhash, err := relayerTx.RelayProphecyClaimToEthereum(chain33Relayer.web3Provider, chain33Relayer.ethSender, chain33Relayer.contractAddress, claimEvent, prophecyClaim, chain33Relayer.privateKey4Ethereum)

	//保存交易hash，方便查询
	atomic.AddInt64(&chain33Relayer.totalTx4Chain33ToEth, 1)
	txIndex := atomic.LoadInt64(&chain33Relayer.totalTx4Chain33ToEth)
	if err = chain33Relayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	if err = chain33Relayer.setLastestRelay2EthTxhash(txhash, txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	return nil
}
