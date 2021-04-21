package chain33

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	chain33Crypto "github.com/33cn/chain33/common/crypto"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	syncTx "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/chain33/transceiver/sync"
	ebTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var relayerLog = log.New("module", "chain33_relayer")

//Relayer4Chain33 ...
type Relayer4Chain33 struct {
	syncEvmTxLogs       *syncTx.EVMTxLogs
	rpcLaddr            string //用户向指定的blockchain节点进行rpc调用
	chainName           string //用来区别主链中继还是平行链，主链为空，平行链则是user.p.xxx.
	fetchHeightPeriodMs int64
	db                  dbm.DB
	lastHeight4Tx       int64 //等待被处理的具有相应的交易回执的高度
	matDegree           int32 //成熟度         heightSync2App    matDegress   height

	privateKey4Chain33       chain33Crypto.PrivKey
	privateKey4Chain33_ecdsa *ecdsa.PrivateKey
	ctx                      context.Context
	rwLock                   sync.RWMutex
	unlock                   chan int
	bridgeBankEventLockSig   string
	bridgeBankEventBurnSig   string
	bridgeBankAbi            abi.ABI
	deployInfo               *ebTypes.Deploy
	//新增//
	ethBridgeClaimChan <-chan *ebrelayerTypes.EthBridgeClaim
	chain33MsgChan     chan<- *events.Chain33Msg
	bridgeRegistryAddr string
	oracleAddr         string
	bridgeBankAddr     string
	deployResult       *X2EthDeployResult
}

type Chain33StartPara struct {
	ChainName          string
	Ctx                context.Context
	SyncTxConfig       *ebTypes.SyncTxConfig
	BridgeRegistryAddr string
	DeployInfo         *ebTypes.Deploy
	DBHandle           dbm.DB
	EthBridgeClaimChan <-chan *ebrelayerTypes.EthBridgeClaim
	Chain33MsgChan     chan<- *events.Chain33Msg
}

// StartChain33Relayer : initializes a relayer which witnesses events on the chain33 network and relays them to Ethereum
func StartChain33Relayer(startPara *Chain33StartPara) *Relayer4Chain33 {
	chain33Relayer := &Relayer4Chain33{
		rpcLaddr:            startPara.SyncTxConfig.Chain33Host,
		chainName:           startPara.ChainName,
		fetchHeightPeriodMs: startPara.SyncTxConfig.FetchHeightPeriodMs,
		unlock:              make(chan int),
		db:                  startPara.DBHandle,
		ctx:                 startPara.Ctx,
		deployInfo:          startPara.DeployInfo,
		bridgeRegistryAddr:  startPara.BridgeRegistryAddr,
		ethBridgeClaimChan:  startPara.EthBridgeClaimChan,
		chain33MsgChan:      startPara.Chain33MsgChan,
	}

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		Chain33Host:       startPara.SyncTxConfig.Chain33Host,
		PushHost:          startPara.SyncTxConfig.PushHost,
		PushName:          startPara.SyncTxConfig.PushName,
		PushBind:          startPara.SyncTxConfig.PushBind,
		StartSyncHeight:   startPara.SyncTxConfig.StartSyncHeight,
		StartSyncSequence: startPara.SyncTxConfig.StartSyncSequence,
		StartSyncHash:     startPara.SyncTxConfig.StartSyncHash,
		Contracts:         startPara.SyncTxConfig.Contracts,
	}

	registrAddrInDB, err := chain33Relayer.getBridgeRegistryAddr()
	//如果输入的registry地址非空，且和数据库保存地址不一致，则直接使用输入注册地址
	if chain33Relayer.bridgeRegistryAddr != "" && nil == err && registrAddrInDB != chain33Relayer.bridgeRegistryAddr {
		relayerLog.Error("StartChain33Relayer", "BridgeRegistry is setted already with value", registrAddrInDB,
			"but now setting to", startPara.BridgeRegistryAddr)
		_ = chain33Relayer.setBridgeRegistryAddr(startPara.BridgeRegistryAddr)
	} else if startPara.BridgeRegistryAddr == "" && registrAddrInDB != "" {
		//输入地址为空，且数据库中保存地址不为空，则直接使用数据库中的地址
		chain33Relayer.bridgeRegistryAddr = registrAddrInDB
	}

	go chain33Relayer.syncProc(syncCfg)
	return chain33Relayer
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
	//如果该中继器的bridgeRegistryAddr为空，就说明合约未部署，需要等待部署成功之后再继续
	if "" == chain33Relayer.bridgeRegistryAddr {
		<-chain33Relayer.unlock
	}
	//如果oracleAddr为空，则通过bridgeRegistry合约进行查询
	if "" != chain33Relayer.bridgeRegistryAddr && "" == chain33Relayer.oracleAddr {
		oracleAddr, bridgeBankAddr := recoverContractAddrFromRegistry(chain33Relayer.bridgeRegistryAddr, chain33Relayer.rpcLaddr)
		if "" == oracleAddr || "" == bridgeBankAddr {
			panic("Failed to recoverContractAddrFromRegistry")
		}
		chain33Relayer.oracleAddr = oracleAddr
		chain33Relayer.bridgeBankAddr = bridgeBankAddr
		chain33txLog.Debug("recoverContractAddrFromRegistry", "bridgeRegistryAddr", chain33Relayer.bridgeRegistryAddr,
			"oracleAddr", chain33Relayer.oracleAddr, "bridgeBankAddr", chain33Relayer.bridgeBankAddr)
	}

	if 0 == len(syncCfg.Contracts) {
		syncCfg.Contracts = append(syncCfg.Contracts, chain33Relayer.bridgeBankAddr)
	}

	chain33Relayer.syncEvmTxLogs = syncTx.StartSyncEvmTxLogs(syncCfg, chain33Relayer.db)
	chain33Relayer.lastHeight4Tx = chain33Relayer.loadLastSyncHeight()
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

		case ethBridgeClaim := <-chain33Relayer.ethBridgeClaimChan:
			chain33Relayer.relayLockBurnToChain33(ethBridgeClaim)
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
			if !strings.Contains(string(tx.Execer), "evm") || tx.To != chain33Relayer.bridgeBankAddr {
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
	chain33Msg, err := events.ParseBurnLock4chain33(evmEventType, data, chain33Relayer.bridgeBankAbi, chain33TxHash)
	if nil != err {
		return err
	}

	chain33Relayer.chain33MsgChan <- chain33Msg

	return nil
}

//DeployContrcts 部署以太坊合约
func (chain33Relayer *Relayer4Chain33) DeployContracts() (bridgeRegistry string, err error) {
	bridgeRegistry = ""
	if nil == chain33Relayer.deployInfo {
		return bridgeRegistry, errors.New("no deploy info configured yet")
	}
	if len(chain33Relayer.deployInfo.ValidatorsAddr) != len(chain33Relayer.deployInfo.InitPowers) {
		return bridgeRegistry, errors.New("not same number for validator address and power")
	}
	if len(chain33Relayer.deployInfo.ValidatorsAddr) < 3 {
		return bridgeRegistry, errors.New("the number of validator must be not less than 3")
	}

	//已经设置了注册合约地址，说明已经部署了相关的合约，不再重复部署
	if chain33Relayer.bridgeRegistryAddr != "" {
		return bridgeRegistry, errors.New("contract deployed already")
	}

	var validators []address.Address
	var initPowers []*big.Int

	for i, addrStr := range chain33Relayer.deployInfo.ValidatorsAddr {
		addr, err := address.NewAddrFromString(addrStr)
		if nil != err {
			panic(fmt.Sprintf("Failed to NewAddrFromString for:%s", addrStr))
		}
		validators = append(validators, *addr)
		initPowers = append(initPowers, big.NewInt(chain33Relayer.deployInfo.InitPowers[i]))
	}
	deployerAddr, err := address.NewAddrFromString(chain33Relayer.deployInfo.OperatorAddr)
	if nil != err {
		panic(fmt.Sprintf("Failed to NewAddrFromString for:%s", chain33Relayer.deployInfo.OperatorAddr))
	}
	para4deploy := &DeployPara4Chain33{
		Deployer:       *deployerAddr,
		Operator:       *deployerAddr,
		InitValidators: validators,
		InitPowers:     initPowers,
	}

	for i, power := range para4deploy.InitPowers {
		relayerLog.Info("deploy", "the validator address ", para4deploy.InitValidators[i].String(),
			"power", power.String())
	}

	x2EthDeployInfo, err := deployAndInit2Chain33(chain33Relayer.rpcLaddr, chain33Relayer.chainName, para4deploy)
	if err != nil {
		return bridgeRegistry, err
	}
	chain33Relayer.rwLock.Lock()

	chain33Relayer.deployResult = x2EthDeployInfo
	bridgeRegistry = x2EthDeployInfo.BridgeRegistry.Address.String()
	_ = chain33Relayer.setBridgeRegistryAddr(bridgeRegistry)
	//设置注册合约地址，同时设置启动中继服务的信号
	chain33Relayer.bridgeRegistryAddr = bridgeRegistry
	chain33Relayer.oracleAddr = x2EthDeployInfo.Oracle.Address.String()
	chain33Relayer.bridgeBankAddr = x2EthDeployInfo.BridgeBank.Address.String()
	chain33Relayer.rwLock.Unlock()
	chain33Relayer.unlock <- start
	relayerLog.Info("deploy", "the BridgeRegistry address is", bridgeRegistry)

	return bridgeRegistry, nil
}

func (chain33Relayer *Relayer4Chain33) relayLockBurnToChain33(claim *ebrelayerTypes.EthBridgeClaim) {
	relayerLog.Debug("relayLockBurnToChain33", "new EthBridgeClaim received", claim)

	nonceBytes := big.NewInt(claim.Nonce).Bytes()
	amountBytes := big.NewInt(claim.Amount).Bytes()
	claimID := crypto.Keccak256Hash(nonceBytes, []byte(claim.EthereumSender), []byte(claim.Chain33Receiver), []byte(claim.Symbol), amountBytes)

	// Sign the hash using the active validator's private key
	signature, err := utils.SignClaim4Evm(claimID, chain33Relayer.privateKey4Chain33_ecdsa)
	if nil != err {
		panic("SignClaim4Evm due to" + err.Error())
	}
	//function newOracleClaim(
	//	ClaimType _claimType,
	//	bytes memory _ethereumSender,
	//	address payable _chain33Receiver,
	//	address _tokenAddress,
	//	string memory _symbol,
	//	uint256 _amount,
	//	bytes32 _claimID,
	//	bytes memory _signature
	//)
	//addr := chain33EvmCommon.BytesToAddress([]byte{1})
	//precompiledContract := runtime.PrecompiledContractsByzantium[addr]
	//result, err := precompiledContract.Run(signature)
	relayerLog.Debug("relayLockBurnToChain33", "chain33Relayer.privateKey4Chain33.PubKey().Bytes()",
		common.ToHex(chain33Relayer.privateKey4Chain33.PubKey().Bytes()))
	//chain33Relayer.privateKey4Chain33.PubKey().Bytes()=0x02504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a4387
	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %d, %s, %s)",
		claim.ClaimType,
		claim.EthereumSender,
		claim.Chain33Receiver,
		claim.TokenAddr,
		claim.Symbol,
		claim.Amount,
		claimID.String(),
		common.ToHex(signature))

	claim.ChainName = chain33Relayer.chainName
	txhash, err := relayEvmTx2Chain33(chain33Relayer.privateKey4Chain33, claim, parameter, chain33Relayer.rpcLaddr, chain33Relayer.oracleAddr)
	if err != nil {
		relayerLog.Error("relayLockBurnToChain33", "Failed to RelayEvmTx2Chain33 due to:", err.Error())
		return
	}
	relayerLog.Info("relayLockBurnToChain33", "RelayLockToChain33 with hash", txhash)
}

func (chain33Relayer *Relayer4Chain33) BurnAsyncFromChain33(ownerPrivateKey, tokenAddr, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return burnAsync(ownerPrivateKey, tokenAddr, ethereumReceiver, bn.Int64(), chain33Relayer.bridgeBankAddr, chain33Relayer.chainName, chain33Relayer.rpcLaddr)
}

func (chain33Relayer *Relayer4Chain33) LockBTYAssetAsync(ownerPrivateKey, ethereumReceiver, amount string) (string, error) {
	bn := big.NewInt(1)
	bn, _ = bn.SetString(utils.TrimZeroAndDot(amount), 10)
	return lockAsync(ownerPrivateKey, ethereumReceiver, bn.Int64(), chain33Relayer.bridgeBankAddr, chain33Relayer.chainName, chain33Relayer.rpcLaddr)
}

//ShowBridgeRegistryAddr ...
func (chain33Relayer *Relayer4Chain33) ShowBridgeRegistryAddr() (string, error) {
	if "" == chain33Relayer.bridgeRegistryAddr {
		return "", errors.New("the relayer is not started yet")
	}

	return chain33Relayer.bridgeRegistryAddr, nil
}
