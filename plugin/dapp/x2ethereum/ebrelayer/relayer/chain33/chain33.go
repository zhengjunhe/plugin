package dplatform

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	dbm "github.com/33cn/dplatform/common/db"
	log "github.com/33cn/dplatform/common/log/log15"
	"github.com/33cn/dplatform/rpc/jsonclient"
	rpctypes "github.com/33cn/dplatform/rpc/types"
	dplatformTypes "github.com/33cn/dplatform/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethinterface"
	relayerTx "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethtxs"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/events"
	syncTx "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/relayer/dplatform/transceiver/sync"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/utils"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var relayerLog = log.New("module", "dplatform_relayer")

//Relayer4Dplatform ...
type Relayer4Dplatform struct {
	syncTxReceipts      *syncTx.TxReceipts
	ethClient           ethinterface.EthClientSpec
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	fetchHeightPeriodMs int64
	db                  dbm.DB
	lastHeight4Tx       int64 //等待被处理的具有相应的交易回执的高度
	matDegree           int32 //成熟度         heightSync2App    matDegress   height
	//passphase            string
	privateKey4Ethereum  *ecdsa.PrivateKey
	ethSender            ethCommon.Address
	bridgeRegistryAddr   ethCommon.Address
	oracleInstance       *generated.Oracle
	totalTx4DplatformToEth int64
	statusCheckedIndex   int64
	ctx                  context.Context
	rwLock               sync.RWMutex
	unlock               chan int
}

// StartDplatformRelayer : initializes a relayer which witnesses events on the dplatform network and relays them to Ethereum
func StartDplatformRelayer(ctx context.Context, syncTxConfig *ebTypes.SyncTxConfig, registryAddr, provider string, db dbm.DB) *Relayer4Dplatform {
	chian33Relayer := &Relayer4Dplatform{
		rpcLaddr:            syncTxConfig.DplatformHost,
		fetchHeightPeriodMs: syncTxConfig.FetchHeightPeriodMs,
		unlock:              make(chan int),
		db:                  db,
		ctx:                 ctx,
		bridgeRegistryAddr:  ethCommon.HexToAddress(registryAddr),
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		DplatformHost:       syncTxConfig.DplatformHost,
		PushHost:          syncTxConfig.PushHost,
		PushName:          syncTxConfig.PushName,
		PushBind:          syncTxConfig.PushBind,
		StartSyncHeight:   syncTxConfig.StartSyncHeight,
		StartSyncSequence: syncTxConfig.StartSyncSequence,
		StartSyncHash:     syncTxConfig.StartSyncHash,
	}

	client, err := relayerTx.SetupWebsocketEthClient(provider)
	if err != nil {
		panic(err)
	}
	chian33Relayer.ethClient = client
	chian33Relayer.totalTx4DplatformToEth = chian33Relayer.getTotalTxAmount2Eth()
	chian33Relayer.statusCheckedIndex = chian33Relayer.getStatusCheckedIndex()

	go chian33Relayer.syncProc(syncCfg)
	return chian33Relayer
}

//QueryTxhashRelay2Eth ...
func (dplatformRelayer *Relayer4Dplatform) QueryTxhashRelay2Eth() ebTypes.Txhashes {
	txhashs := utils.QueryTxhashes([]byte(dplatformToEthBurnLockTxHashPrefix), dplatformRelayer.db)
	return ebTypes.Txhashes{Txhash: txhashs}
}

func (dplatformRelayer *Relayer4Dplatform) syncProc(syncCfg *ebTypes.SyncTxReceiptConfig) {
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Dplatform relayer")
	<-dplatformRelayer.unlock
	_, _ = fmt.Fprintln(os.Stdout, "Dplatform relayer starts to run...")

	dplatformRelayer.syncTxReceipts = syncTx.StartSyncTxReceipt(syncCfg, dplatformRelayer.db)
	dplatformRelayer.lastHeight4Tx = dplatformRelayer.loadLastSyncHeight()

	oracleInstance, err := relayerTx.RecoverOracleInstance(dplatformRelayer.ethClient, dplatformRelayer.bridgeRegistryAddr, dplatformRelayer.bridgeRegistryAddr)
	if err != nil {
		panic(err.Error())
	}
	dplatformRelayer.oracleInstance = oracleInstance

	timer := time.NewTicker(time.Duration(dplatformRelayer.fetchHeightPeriodMs) * time.Millisecond)
	for {
		select {
		case <-timer.C:
			height := dplatformRelayer.getCurrentHeight()
			relayerLog.Debug("syncProc", "getCurrentHeight", height)
			dplatformRelayer.onNewHeightProc(height)

		case <-dplatformRelayer.ctx.Done():
			timer.Stop()
			return
		}
	}
}

func (dplatformRelayer *Relayer4Dplatform) getCurrentHeight() int64 {
	var res rpctypes.Header
	ctx := jsonclient.NewRPCCtx(dplatformRelayer.rpcLaddr, "Dplatform.GetLastHeader", nil, &res)
	_, err := ctx.RunResult()
	if nil != err {
		relayerLog.Error("getCurrentHeight", "Failede due to:", err.Error())
	}
	return res.Height
}

func (dplatformRelayer *Relayer4Dplatform) onNewHeightProc(currentHeight int64) {
	//检查已经提交的交易结果
	dplatformRelayer.rwLock.Lock()
	for dplatformRelayer.statusCheckedIndex < dplatformRelayer.totalTx4DplatformToEth {
		index := dplatformRelayer.statusCheckedIndex + 1
		txhash, err := dplatformRelayer.getEthTxhash(index)
		if nil != err {
			relayerLog.Error("onNewHeightProc", "getEthTxhash for index ", index, "error", err.Error())
			break
		}
		status := relayerTx.GetEthTxStatus(dplatformRelayer.ethClient, txhash)
		//按照提交交易的先后顺序检查交易，只要出现当前交易还在pending状态，就不再检查后续交易，等到下个区块再从该交易进行检查
		//TODO:可能会由于网络和打包挖矿的原因，使得交易执行顺序和提交顺序有差别，后续完善该检查逻辑
		if status == relayerTx.EthTxPending.String() {
			break
		}
		_ = dplatformRelayer.setLastestRelay2EthTxhash(status, txhash.Hex(), index)
		atomic.AddInt64(&dplatformRelayer.statusCheckedIndex, 1)
		_ = dplatformRelayer.setStatusCheckedIndex(dplatformRelayer.statusCheckedIndex)
	}
	dplatformRelayer.rwLock.Unlock()
	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	// lastHeight4Tx    matDegress   currentHeight
	for dplatformRelayer.lastHeight4Tx+int64(dplatformRelayer.matDegree)+1 <= currentHeight {
		relayerLog.Info("onNewHeightProc", "currHeight", currentHeight, "lastHeight4Tx", dplatformRelayer.lastHeight4Tx)

		lastHeight4Tx := dplatformRelayer.lastHeight4Tx
		TxReceipts, err := dplatformRelayer.syncTxReceipts.GetNextValidTxReceipts(lastHeight4Tx)
		if nil == TxReceipts || nil != err {
			if err != nil {
				relayerLog.Error("onNewHeightProc", "Failed to GetNextValidTxReceipts due to:", err.Error())
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
			_ = dplatformTypes.Decode(tx.Payload, &ss)
			actionName := ss.GetActionName()
			if relayerTx.BurnAction == actionName || relayerTx.LockAction == actionName {
				relayerLog.Debug("^_^ ^_^ Processing dplatform tx receipt", "ActionName", actionName, "fromAddr", tx.From(), "exec", string(tx.Execer))
				actionEvent := getOracleClaimType(actionName)
				if err := dplatformRelayer.handleBurnLockMsg(actionEvent, TxReceipts.ReceiptData[i], tx.Hash()); nil != err {
					errInfo := fmt.Sprintf("Failed to handleBurnLockMsg due to:%s", err.Error())
					panic(errInfo)
				}
			}
		}
		dplatformRelayer.lastHeight4Tx = TxReceipts.Height
		dplatformRelayer.setLastSyncHeight(dplatformRelayer.lastHeight4Tx)
	}
}

// getOracleClaimType : sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) events.Event {
	var claimType events.Event

	switch eventType {
	case events.MsgBurn.String():
		claimType = events.Event(events.ClaimTypeBurn)
	case events.MsgLock.String():
		claimType = events.Event(events.ClaimTypeLock)
	default:
		panic(errors.New("eventType invalid"))
	}

	return claimType
}

// handleBurnLockMsg : parse event data as a DplatformMsg, package it into a ProphecyClaim, then relay tx to the Ethereum Network
func (dplatformRelayer *Relayer4Dplatform) handleBurnLockMsg(claimEvent events.Event, receipt *dplatformTypes.ReceiptData, dplatformTxHash []byte) error {
	relayerLog.Info("handleBurnLockMsg", "Received tx with hash", ethCommon.Bytes2Hex(dplatformTxHash))

	// Parse the witnessed event's data into a new DplatformMsg
	dplatformMsg := relayerTx.ParseBurnLockTxReceipt(claimEvent, receipt)
	if nil == dplatformMsg {
		//收到执行失败的交易，直接跳过
		relayerLog.Error("handleBurnLockMsg", "Received failed tx with hash", ethCommon.Bytes2Hex(dplatformTxHash))
		return nil
	}

	// Parse the DplatformMsg into a ProphecyClaim for relay to Ethereum
	prophecyClaim := relayerTx.DplatformMsgToProphecyClaim(*dplatformMsg)

	// Relay the DplatformMsg to the Ethereum network
	txhash, err := relayerTx.RelayOracleClaimToEthereum(dplatformRelayer.oracleInstance, dplatformRelayer.ethClient, dplatformRelayer.ethSender, claimEvent, prophecyClaim, dplatformRelayer.privateKey4Ethereum, dplatformTxHash)
	if nil != err {
		return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&dplatformRelayer.totalTx4DplatformToEth, 1)
	txIndex := atomic.LoadInt64(&dplatformRelayer.totalTx4DplatformToEth)
	if err = dplatformRelayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToDplatform due to:", err.Error())
		return err
	}
	if err = dplatformRelayer.setLastestRelay2EthTxhash(relayerTx.EthTxPending.String(), txhash, txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToDplatform due to:", err.Error())
		return err
	}
	return nil
}
