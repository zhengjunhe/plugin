package chain33

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/events"

	"github.com/33cn/chain33/common"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	ethContract "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/contracts/contracts4eth/generated"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/ethinterface"
	ethTx "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/ethtxs"
	syncTx "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/chain33/transceiver/sync"
	ebTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

var relayerLog = log.New("module", "chain33_relayer")

//Relayer4Chain33 ...
type Relayer4Chain33 struct {
	syncEvmTxLogs       *syncTx.EVMTxLogs
	ethClient           ethinterface.EthClientSpec
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	fetchHeightPeriodMs int64
	db                  dbm.DB
	lastHeight4Tx       int64 //等待被处理的具有相应的交易回执的高度
	matDegree           int32 //成熟度         heightSync2App    matDegress   height
	//passphase            string
	privateKey4Ethereum    *ecdsa.PrivateKey
	ethSender              ethCommon.Address
	bridgeRegistryAddr     ethCommon.Address
	oracleInstanceOnEth    *ethContract.Oracle //此处需要使用eth的合约句柄
	totalTx4Chain33ToEth   int64
	statusCheckedIndex     int64
	ctx                    context.Context
	rwLock                 sync.RWMutex
	unlock                 chan int
	chain33EVMAddr         string
	bridgeBankEventLockSig string
	bridgeBankEventBurnSig string
	bridgeBankAbi          abi.ABI
}

// StartChain33Relayer : initializes a relayer which witnesses events on the chain33 network and relays them to Ethereum
func StartChain33Relayer(ctx context.Context, syncTxConfig *ebTypes.SyncTxConfig, registryAddr, provider string, db dbm.DB) *Relayer4Chain33 {
	chian33Relayer := &Relayer4Chain33{
		rpcLaddr:            syncTxConfig.Chain33Host,
		fetchHeightPeriodMs: syncTxConfig.FetchHeightPeriodMs,
		unlock:              make(chan int),
		db:                  db,
		ctx:                 ctx,
		bridgeRegistryAddr:  ethCommon.HexToAddress(registryAddr),
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

	client, err := ethTx.SetupWebsocketEthClient(provider)
	if err != nil {
		panic(err)
	}
	chian33Relayer.ethClient = client
	chian33Relayer.totalTx4Chain33ToEth = chian33Relayer.getTotalTxAmount2Eth()
	chian33Relayer.statusCheckedIndex = chian33Relayer.getStatusCheckedIndex()

	go chian33Relayer.syncProc(syncCfg)
	return chian33Relayer
}

//QueryTxhashRelay2Eth ...
func (chain33Relayer *Relayer4Chain33) QueryTxhashRelay2Eth() ebTypes.Txhashes {
	txhashs := utils.QueryTxhashes([]byte(chain33ToEthBurnLockTxHashPrefix), chain33Relayer.db)
	return ebTypes.Txhashes{Txhash: txhashs}
}

func (chain33Relayer *Relayer4Chain33) syncProc(syncCfg *ebTypes.SyncTxReceiptConfig) {
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Chain33 relayer")
	<-chain33Relayer.unlock
	_, _ = fmt.Fprintln(os.Stdout, "Chain33 relayer starts to run...")

	chain33Relayer.syncEvmTxLogs = syncTx.StartSyncEvmTxLogs(syncCfg, chain33Relayer.db)
	chain33Relayer.lastHeight4Tx = chain33Relayer.loadLastSyncHeight()

	oracleInstance, err := ethTx.RecoverOracleInstance(chain33Relayer.ethClient, chain33Relayer.bridgeRegistryAddr, chain33Relayer.bridgeRegistryAddr)
	if err != nil {
		panic(err.Error())
	}
	chain33Relayer.oracleInstanceOnEth = oracleInstance
	chain33Relayer.prePareSubscribeEvent()

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

func (chain33Relayer *Relayer4Chain33) getCurrentHeight() int64 {
	var res rpctypes.Header
	ctx := jsonclient.NewRPCCtx(chain33Relayer.rpcLaddr, "Chain33.GetLastHeader", nil, &res)
	_, err := ctx.RunResult()
	if nil != err {
		relayerLog.Error("getCurrentHeight", "Failede due to:", err.Error())
	}
	return res.Height
}

func (chain33Relayer *Relayer4Chain33) onNewHeightProc(currentHeight int64) {
	//检查已经提交的交易结果
	chain33Relayer.rwLock.Lock()
	for chain33Relayer.statusCheckedIndex < chain33Relayer.totalTx4Chain33ToEth {
		index := chain33Relayer.statusCheckedIndex + 1
		txhash, err := chain33Relayer.getEthTxhash(index)
		if nil != err {
			relayerLog.Error("onNewHeightProc", "getEthTxhash for index ", index, "error", err.Error())
			break
		}
		status := ethTx.GetEthTxStatus(chain33Relayer.ethClient, txhash)
		//按照提交交易的先后顺序检查交易，只要出现当前交易还在pending状态，就不再检查后续交易，等到下个区块再从该交易进行检查
		//TODO:可能会由于网络和打包挖矿的原因，使得交易执行顺序和提交顺序有差别，后续完善该检查逻辑
		if status == ethTx.EthTxPending.String() {
			break
		}
		_ = chain33Relayer.setLastestRelay2EthTxhash(status, txhash.Hex(), index)
		atomic.AddInt64(&chain33Relayer.statusCheckedIndex, 1)
		_ = chain33Relayer.setStatusCheckedIndex(chain33Relayer.statusCheckedIndex)
	}
	chain33Relayer.rwLock.Unlock()
	//未达到足够的成熟度，不进行处理
	//  +++++++++||++++++++++++||++++++++++||
	//           ^             ^           ^
	// lastHeight4Tx    matDegress   currentHeight
	for chain33Relayer.lastHeight4Tx+int64(chain33Relayer.matDegree)+1 <= currentHeight {
		relayerLog.Info("onNewHeightProc", "currHeight", currentHeight, "lastHeight4Tx", chain33Relayer.lastHeight4Tx)

		lastHeight4Tx := chain33Relayer.lastHeight4Tx
		txLogs, err := chain33Relayer.syncEvmTxLogs.GetNextValidEvmTxLogs(lastHeight4Tx)
		if nil == txLogs || nil != err {
			if err != nil {
				relayerLog.Error("onNewHeightProc", "Failed to GetNextValidTxReceipts due to:", err.Error())
			}
			break
		}
		relayerLog.Debug("onNewHeightProc", "currHeight", currentHeight, "valid tx receipt with height:", txLogs.Height)

		txAndLogs := txLogs.TxAndLogs
		for _, txAndLog := range txAndLogs {
			tx := txAndLog.Tx

			//确认订阅的evm交易类型和合约地址
			if !strings.Contains(string(tx.Execer), "evm") || tx.To != chain33Relayer.chain33EVMAddr {
				relayerLog.Error("onNewHeightProc received logs not expected", "tx.Execer", string(tx.Execer), "tx.To", tx.To)
				continue
			}

			for _, evmlog := range txAndLog.LogsPerTx.Logs {
				var evmEventType events.Chain33EvmEvent
				if chain33Relayer.bridgeBankEventBurnSig == common.ToHex(evmlog.Topic) {
					evmEventType = events.Chain33EventLogBurn
				} else if chain33Relayer.bridgeBankEventLockSig == common.ToHex(evmlog.Topic) {
					evmEventType = events.Chain33EventLogLock
				} else {
					continue
				}

				if err := chain33Relayer.handleBurnLockEvent(evmEventType, evmlog.Data, tx.Hash()); nil != err {
					errInfo := fmt.Sprintf("Failed to handleBurnLockEvent due to:%s", err.Error())
					panic(errInfo)
				}
			}
		}
		chain33Relayer.lastHeight4Tx = txLogs.Height
		chain33Relayer.setLastSyncHeight(chain33Relayer.lastHeight4Tx)
	}
}

// handleBurnLockMsg : parse event data as a Chain33Msg, package it into a ProphecyClaim, then relay tx to the Ethereum Network
//
func (chain33Relayer *Relayer4Chain33) handleBurnLockEvent(evmEventType events.Chain33EvmEvent, data []byte, chain33TxHash []byte) error {
	relayerLog.Info("handleBurnLockEvent", "Received tx with hash", ethCommon.Bytes2Hex(chain33TxHash))

	// Parse the witnessed event's data into a new Chain33Msg
	chain33Msg, err := events.ParseBurnLock4chain33(evmEventType, data, chain33Relayer.bridgeBankAbi)
	if nil != err {
		return err
	}

	// Parse the Chain33Msg into a ProphecyClaim for relay to Ethereum
	prophecyClaim := ethTx.Chain33MsgToProphecyClaim(*chain33Msg)

	// Relay the Chain33Msg to the Ethereum network
	txhash, err := ethTx.RelayOracleClaimToEthereum(chain33Relayer.oracleInstanceOnEth, chain33Relayer.ethClient, chain33Relayer.ethSender, prophecyClaim, chain33Relayer.privateKey4Ethereum, chain33TxHash)
	if nil != err {
		return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&chain33Relayer.totalTx4Chain33ToEth, 1)
	txIndex := atomic.LoadInt64(&chain33Relayer.totalTx4Chain33ToEth)
	if err = chain33Relayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	if err = chain33Relayer.setLastestRelay2EthTxhash(ethTx.EthTxPending.String(), txhash, txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	return nil
}
