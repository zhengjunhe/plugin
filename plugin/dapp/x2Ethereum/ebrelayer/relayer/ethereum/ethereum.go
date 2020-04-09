package ethereum

// -----------------------------------------------------
//      Ethereum relayer
//
//      Initializes the relayer service, which parses,
//      encodes, and packages named events on an Ethereum
//      Smart Contract for validator's to sign and send
//      to the Chain33 bridge.
// -----------------------------------------------------

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	chain33Crypto "github.com/33cn/chain33/common/crypto"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethtxs"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
)

type EthereumRelayer struct {
	provider           string
	bridgeRegistryAddr common.Address
	//validatorName        string
	db dbm.DB
	//passphase            string
	rwLock               sync.RWMutex
	validatorAddress     []byte
	privateKey4Chain33   chain33Crypto.PrivKey
	privateKey4Ethereum  *ecdsa.PrivateKey
	ethValidator         common.Address
	totalTx4Eth2Chain33  int64
	totalTx4Chain33ToEth int64
	rpcURL2Chain33       string
	unlockchan           chan int
	status               int32
	//bridgeBankAbi        *abi.ABI
	client                *ethclient.Client
	bridgeBankAddr        common.Address
	chain33BridgeAddr     common.Address
	bridgeBankSub         ethereum.Subscription
	chain33BridgeSub      ethereum.Subscription
	bridgeBankLog         chan types.Log
	chain33BridgeLog      chan types.Log
	bridgeBankEventSig    string
	chain33BridgeEventSig string
	bridgeBankAbi         abi.ABI
	chain33BridgeAbi      abi.ABI
	deployInfo            *ebTypes.Deploy
	x2EthDeployInfo       *ethtxs.X2EthDeployInfo
	deployPara            *ethtxs.DeployPara
	x2EthContracts        *ethtxs.X2EthContracts
}

var (
	relayerLog = log.New("module", "ethereum_relayer")
)

func StartEthereumRelayer(rpcURL2Chain33 string, db dbm.DB, provider, registryAddress string, deploy *ebTypes.Deploy) *EthereumRelayer {
	relayer := &EthereumRelayer{
		provider:           provider,
		db:                 db,
		unlockchan:         make(chan int),
		rpcURL2Chain33:     rpcURL2Chain33,
		status:             ebTypes.StatusPending,
		bridgeRegistryAddr: common.HexToAddress(registryAddress),
		deployInfo:         deploy,
	}

	registrAddrInDB, err := relayer.getBridgeRegistryAddr()
	if nil == err && registrAddrInDB != registryAddress {
		relayerLog.Error("StartEthereumRelayer", "BridgeRegistry is setted already with value", registrAddrInDB,
			"but now setting to", registryAddress)
		_ = relayer.setBridgeRegistryAddr(registryAddress)
	}

	go relayer.proc()
	return relayer
}

func (ethRelayer *EthereumRelayer) SetPrivateKey4Ethereum(privateKey4Ethereum *ecdsa.PrivateKey) {
	ethRelayer.rwLock.Lock()
	defer ethRelayer.rwLock.Unlock()
	ethRelayer.privateKey4Ethereum = privateKey4Ethereum
	if ethRelayer.privateKey4Chain33 != nil {
		ethRelayer.unlockchan <- start
	}
}

func (ethRelayer *EthereumRelayer) GetRunningStatus() (relayerRunStatus *ebTypes.RelayerRunStatus) {
	relayerRunStatus = &ebTypes.RelayerRunStatus{}
	ethRelayer.rwLock.RLock()
	relayerRunStatus.Status = ethRelayer.status
	ethRelayer.rwLock.RUnlock()
	if relayerRunStatus.Status == ebTypes.StatusPending {
		if nil == ethRelayer.privateKey4Ethereum {
			relayerRunStatus.Details = "Ethereum's private key not imported"
		}

		if nil == ethRelayer.privateKey4Chain33 {
			relayerRunStatus.Details += "\nChain33's private key not imported"
		}
		return
	}
	relayerRunStatus.Details = "Running"
	return
}

//部署以太坊合约
func (ethRelayer *EthereumRelayer) DeployContrcts() (bridgeRegistry string, err error) {
	bridgeRegistry = ""
	if nil == ethRelayer.deployInfo {
		return bridgeRegistry, errors.New("No deploy info configured yet")
	}
	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(ethRelayer.deployInfo.DeployerPrivateKey))
	if nil != err {
		return bridgeRegistry, err
	}
	if len(ethRelayer.deployInfo.ValidatorsAddr) != len(ethRelayer.deployInfo.InitPowers) {
		return bridgeRegistry, errors.New("Not same number for validator address and power")
	}
	if len(ethRelayer.deployInfo.ValidatorsAddr) < 3 {
		return bridgeRegistry, errors.New("The number of validator must be not less than 3")
	}

	nilAddr := common.Address{}

	//已经设置了注册合约地址，说明已经部署了相关的合约，不再重复部署
	if ethRelayer.bridgeRegistryAddr != nilAddr {
		return bridgeRegistry, errors.New("Contract deployed already")
	}

	var validators []common.Address
	var initPowers []*big.Int

	for i, addr := range ethRelayer.deployInfo.ValidatorsAddr {
		validators = append(validators, common.HexToAddress(addr))
		initPowers = append(initPowers, big.NewInt(ethRelayer.deployInfo.InitPowers[i]))
	}
	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)
	para := &ethtxs.DeployPara{
		DeployPrivateKey: deployPrivateKey,
		Deployer:         deployerAddr,
		Operator:         deployerAddr,
		InitValidators:   validators,
		InitPowers:       initPowers,
	}

	x2EthContracts, x2EthDeployInfo, err := ethtxs.DeployAndInit(ethRelayer.client, para)
	if err != nil {
		return bridgeRegistry, err
	}
	ethRelayer.deployPara = para
	ethRelayer.x2EthDeployInfo = x2EthDeployInfo
	ethRelayer.x2EthContracts = x2EthContracts
	bridgeRegistry = x2EthDeployInfo.BridgeRegistry.Address.String()
	_ = ethRelayer.setBridgeRegistryAddr(bridgeRegistry)
	//设置注册合约地址，同时设置启动中继服务的信号
	ethRelayer.bridgeRegistryAddr = x2EthDeployInfo.BridgeRegistry.Address
	ethRelayer.unlockchan <- start
	relayerLog.Info("deploy", "the BridgeRegistry address is", bridgeRegistry)

	return bridgeRegistry, nil
}

//GetBalance：获取某一个币种的余额
func (ethRelayer *EthereumRelayer) GetBalance(tokenAddr, owner string) (int64, error) {
	return ethtxs.GetBalance(ethRelayer.client, common.HexToAddress(tokenAddr), common.HexToAddress(owner))
}

func (ethRelayer *EthereumRelayer) ProcessProphecyClaim(prophecyID int64) (string, error) {
	return ethtxs.ProcessProphecyClaim(ethRelayer.client, ethRelayer.deployPara, ethRelayer.x2EthContracts, prophecyID)
}

func (ethRelayer *EthereumRelayer) MakeNewProphecyClaim(claimType uint8, chain33Sender, tokenAddr, symbol string) (string, error) {
	newProphecyClaimPara := &ethtxs.NewProphecyClaimPara{
		ClaimType:claimType,
		Chain33Sender:common.FromHex(chain33Sender),
		TokenAddr:common.HexToAddress(tokenAddr),
		Symbol:symbol,
	}
	return ethtxs.MakeNewProphecyClaim(newProphecyClaimPara, ethRelayer.client, ethRelayer.deployPara, ethRelayer.x2EthContracts)
}

func (ethRelayer *EthereumRelayer) CreateBridgeToken(symbol string) (string, error) {
	return ethtxs.CreateBridgeToken(symbol, ethRelayer.client, ethRelayer.deployPara, ethRelayer.x2EthDeployInfo, ethRelayer.x2EthContracts)
}

func (ethRelayer *EthereumRelayer) ShowTxReceipt(hash string) (*types.Receipt, error) {
	txhash := common.HexToHash(hash)
	return ethRelayer.client.TransactionReceipt(context.Background(), txhash)
}

func (ethRelayer *EthereumRelayer) proc() {
	// Start client with infura ropsten provider
	relayerLog.Info("EthereumRelayer proc", "Started Ethereum websocket with provider:", ethRelayer.provider,
		"rpcURL2Chain33:", ethRelayer.rpcURL2Chain33)
	client, err := setupWebsocketEthClient(ethRelayer.provider)
	if err != nil {
		panic(err)
	}
	ethRelayer.client = client

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		errinfo := fmt.Sprintf("Failed to get NetworkID due to:%s", err.Error())
		panic(errinfo)
	}

	//等待用户导入
	_, _ = fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Ethereum relayer")
	nilAddr := common.Address{}
	if nilAddr != ethRelayer.bridgeRegistryAddr {
		ethRelayer.unlockchan <- start
	}

	for {
		select {
		case <-ethRelayer.unlockchan:
			if nil != ethRelayer.privateKey4Ethereum && nil != ethRelayer.privateKey4Chain33 && nilAddr != ethRelayer.bridgeRegistryAddr {
				ethRelayer.ethValidator, err = ethtxs.LoadSender(ethRelayer.privateKey4Ethereum)
				if nil != err {
					errinfo := fmt.Sprintf("Failed to load validator for ethereum due to:%s", err.Error())
					panic(errinfo)
				}

				//向chain33Bridge订阅事件
				ethRelayer.subscribeEvent(false)
				//向bridgeBank订阅事件
				ethRelayer.subscribeEvent(true)
				_, _ = fmt.Fprintln(os.Stdout, "Ethereum relayer starts to run...")
				goto latter
			}
		}
	}

latter:
	for {
		select {
		case err := <-ethRelayer.bridgeBankSub.Err():
			panic("bridgeBankSub" + err.Error())
		case err := <-ethRelayer.chain33BridgeSub.Err():
			panic("chain33BridgeSub" + err.Error())
		// vLog is raw event data
		case vLog := <-ethRelayer.bridgeBankLog:
			_, _ = fmt.Fprintln(os.Stdout, "^_^ ^_^ Received bridgeBankLog")
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == ethRelayer.bridgeBankEventSig {
				eventName := events.LogLock.String()
				relayerLog.Info("EthereumRelayer proc", "Witnessed new event:", eventName,
					"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
				err := ethRelayer.handleLogLockEvent(clientChainID, ethRelayer.bridgeBankAbi, eventName, vLog)
				if err != nil {
					errinfo := fmt.Sprintf("Failed to handleLogLockEvent due to:%s", err.Error())
					panic(errinfo)
				}
			}
		case vLog := <-ethRelayer.chain33BridgeLog:
			_, _ = fmt.Fprintln(os.Stdout, "^_^ ^_^ Received chain33BridgeLog", "total topic number", len(vLog.Topics),
				"topic", vLog.Topics[0].Hex())
			for i := 0; i < len(vLog.Topics); i++ {
				_, _ = fmt.Fprintln(os.Stdout, "topic:", i, "topic=", vLog.Topics[i].Hex())
			}
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == ethRelayer.chain33BridgeEventSig {
				eventName := events.LogNewProphecyClaim.String()
				relayerLog.Info("EthereumRelayer proc", "Witnessed new event:", eventName,
					"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
				err := ethRelayer.handleLogNewProphecyClaimEvent(ethRelayer.chain33BridgeAbi, eventName, vLog)
				if err != nil {
					errinfo := fmt.Sprintf("Failed to handleLogNewProphecyClaimEvent due to:%s", err.Error())
					panic(errinfo)
				}
			}
		}
	}
}

func (ethRelayer *EthereumRelayer) subscribeEvent(makeClaims bool) {
	var eventName string
	var target ethtxs.ContractRegistry

	switch makeClaims {
	case true:
		eventName = events.LogNewProphecyClaim.String()
		contactAbi := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
		// Load unique event signature from the named event contained within the contract's ABI
		ethRelayer.chain33BridgeEventSig = contactAbi.Events[eventName].ID().Hex()
		ethRelayer.chain33BridgeAbi = contactAbi
		target = ethtxs.Chain33Bridge
	case false:
		eventName = events.LogLock.String()
		contactAbi := ethtxs.LoadABI(ethtxs.BridgeBankABI)
		ethRelayer.bridgeBankEventSig = contactAbi.Events[eventName].ID().Hex()
		ethRelayer.bridgeBankAbi = contactAbi
		target = ethtxs.BridgeBank
	}

	client := ethRelayer.client
	sender := ethRelayer.ethValidator
	// Get the specific contract's address (Chain33Bridge or BridgeBank)
	targetAddress, err := ethtxs.GetAddressFromBridgeRegistry(client, sender, ethRelayer.bridgeRegistryAddr, target)
	if err != nil {
		errinfo := fmt.Sprintf("Failed to GetAddressFromBridgeRegistry due to:%s", err.Error())
		panic(errinfo)
	}
	if makeClaims {
		ethRelayer.chain33BridgeAddr = *targetAddress
		_, _ = fmt.Fprintln(os.Stdout, "Succeed to get chain33Bridge addr:", targetAddress.String(),
			"event signature", ethRelayer.chain33BridgeEventSig)
	} else {
		ethRelayer.bridgeBankAddr = *targetAddress
		_, _ = fmt.Fprintln(os.Stdout, "Succeed to get BridgeBank addr:", targetAddress.String(),
			"event signature", ethRelayer.bridgeBankEventSig)
	}

	// We need the target address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{*targetAddress},
	}
	// We will check logs for new events
	logs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		errinfo := fmt.Sprintf("Failed to SubscribeFilterLogs due to:%s", err.Error())
		panic(errinfo)
	}
	relayerLog.Info("EthereumRelayer proc", "Subscribed to contract at address:", targetAddress.Hex())

	if makeClaims {
		ethRelayer.chain33BridgeLog = logs
		ethRelayer.chain33BridgeSub = sub
		return
	}
	ethRelayer.bridgeBankLog = logs
	ethRelayer.bridgeBankSub = sub
}

func (ethRelayer *EthereumRelayer) IsValidatorActive(addr string) (bool, error) {

	return true, nil
}

func (ethRelayer *EthereumRelayer) ShowOperator() (string, error) {
	operator, err := ethtxs.GetOperator(ethRelayer.client, ethRelayer.ethValidator, ethRelayer.bridgeBankAddr)
	if nil != err {
		return "", err
	}
	return operator.String(), nil
}

func (ethRelayer *EthereumRelayer) QueryTxhashRelay2Eth() ebTypes.Txhashes {
	txhashs := ethRelayer.queryTxhashes([]byte(chain33ToEthTxHashPrefix))
	return ebTypes.Txhashes{Txhash: txhashs}
}

func (ethRelayer *EthereumRelayer) QueryTxhashRelay2Chain33() *ebTypes.Txhashes {
	txhashs := ethRelayer.queryTxhashes([]byte(eth2chain33TxHashPrefix))
	return &ebTypes.Txhashes{Txhash: txhashs}
}

// handleLogLockEvent : unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to chain33
func (ethRelayer *EthereumRelayer) handleLogLockEvent(clientChainID *big.Int, contractABI abi.ABI, eventName string, log types.Log) error {
	contractAddress := ethRelayer.bridgeRegistryAddr.Hex()
	validatorAddress := ethRelayer.validatorAddress
	rpcURL := ethRelayer.rpcURL2Chain33

	// Unpack the LogLock event using its unique event signature from the contract's ABI
	event, err := events.UnpackLogLock(clientChainID, contractAddress, contractABI, eventName, log.Data)
	if nil != err {
		return err
	}
	// Add the event to the record
	events.NewEventWrite(log.TxHash.Hex(), *event)

	// Parse the LogLock event's payload into a struct
	prophecyClaim, err := ethtxs.LogLockToEthBridgeClaim(validatorAddress, event)
	if err != nil {
		return err
	}

	// Initiate the relay
	txhash, err := ethtxs.RelayLockToChain33(ethRelayer.privateKey4Chain33, &prophecyClaim, rpcURL)
	if err != nil {
		relayerLog.Error("handleLogLockEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&ethRelayer.totalTx4Eth2Chain33, 1)
	txIndex := atomic.LoadInt64(&ethRelayer.totalTx4Eth2Chain33)
	if err = ethRelayer.updateTotalTxAmount2chain33(txIndex); nil != err {
		relayerLog.Error("handleLogLockEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	if err = ethRelayer.setLastestRelay2Chain33Txhash(txhash, txIndex); nil != err {
		relayerLog.Error("handleLogLockEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}

	return nil
}

// handleLogNewProphecyClaimEvent : unpacks a LogNewProphecyClaim event, converts it to a OracleClaim, and relays a tx to Ethereum
func (ethRelayer *EthereumRelayer) handleLogNewProphecyClaimEvent(contractABI abi.ABI, eventName string, log types.Log) error {
	// Unpack the LogNewProphecyClaim event using its unique event signature from the contract's ABI
	event, err := events.UnpackLogNewProphecyClaim(contractABI, eventName, log.Data)
	if nil != err {
		return err
	}
	// Parse ProphecyClaim's data into an OracleClaim
	oracleClaim, err := ethtxs.ProphecyClaimToSignedOracleClaim(event, ethRelayer.privateKey4Ethereum)
	if nil != err {
		return err
	}

	// Initiate the relay
	txhash, err := ethtxs.RelayOracleClaimToEthereum(ethRelayer.provider, ethRelayer.ethValidator, ethRelayer.bridgeRegistryAddr, events.LogNewProphecyClaim, oracleClaim, ethRelayer.privateKey4Ethereum)
    if "" == txhash {
    	return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&ethRelayer.totalTx4Chain33ToEth, 1)
	txIndex := atomic.LoadInt64(&ethRelayer.totalTx4Chain33ToEth)
	if err = ethRelayer.updateTotalTxAmount2Eth(txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	if err = ethRelayer.setLastestRelay2EthTxhash(txhash, txIndex); nil != err {
		relayerLog.Error("handleLogNewProphecyClaimEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	return nil
}
