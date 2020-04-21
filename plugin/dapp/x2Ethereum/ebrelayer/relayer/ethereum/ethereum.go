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
	rwLock                 sync.RWMutex
	privateKey4Chain33     chain33Crypto.PrivKey
	privateKey4Ethereum    *ecdsa.PrivateKey
	ethValidator           common.Address
	totalTx4Eth2Chain33    int64
	totalTx4Chain33ToEth   int64
	rpcURL2Chain33         string
	unlockchan             chan int
	status                 int32
	client                 *ethclient.Client
	bridgeBankAddr         common.Address
	chain33BridgeAddr      common.Address
	bridgeBankSub          ethereum.Subscription
	chain33BridgeSub       ethereum.Subscription
	bridgeBankLog          chan types.Log
	chain33BridgeLog       chan types.Log
	bridgeBankEventLockSig string
	bridgeBankEventBurnSig string
	chain33BridgeEventSig  string
	bridgeBankAbi          abi.ABI
	chain33BridgeAbi       abi.ABI
	deployInfo             *ebTypes.Deploy
	x2EthDeployInfo        *ethtxs.X2EthDeployInfo
	deployPara             *ethtxs.DeployPara
	operatorInfo           *ethtxs.OperatorInfo
	x2EthContracts         *ethtxs.X2EthContracts
}

var (
	relayerLog = log.New("module", "ethereum_relayer")
)

func StartEthereumRelayer(rpcURL2Chain33 string, db dbm.DB, provider, registryAddress string, deploy *ebTypes.Deploy) *EthereumRelayer {
	relayer := &EthereumRelayer{
		provider:           provider,
		db:                 db,
		unlockchan:         make(chan int, 2),
		rpcURL2Chain33:     rpcURL2Chain33,
		status:             ebTypes.StatusPending,
		bridgeRegistryAddr: common.HexToAddress(registryAddress),
		deployInfo:         deploy,
	}

	registrAddrInDB, err := relayer.getBridgeRegistryAddr()
	//如果输入的registry地址非空，且和数据库保存地址不一致，则直接使用输入注册地址
	if registryAddress != "" && nil == err && registrAddrInDB != registryAddress {
		relayerLog.Error("StartEthereumRelayer", "BridgeRegistry is setted already with value", registrAddrInDB,
			"but now setting to", registryAddress)
		_ = relayer.setBridgeRegistryAddr(registryAddress)
	} else if registryAddress == "" && registrAddrInDB != "" {
		//输入地址为空，且数据库中保存地址不为空，则直接使用数据库中的地址
		relayer.bridgeRegistryAddr = common.HexToAddress(registrAddrInDB)
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

func (ethRelayer *EthereumRelayer) recoverDeployPara() (err error) {
	if nil == ethRelayer.deployInfo {
		return nil
	}
	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(ethRelayer.deployInfo.DeployerPrivateKey))
	if nil != err {
		return err
	}
	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)

	ethRelayer.operatorInfo = &ethtxs.OperatorInfo{
		PrivateKey:deployPrivateKey,
		Address:deployerAddr,
	}

	return nil
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
		ValidatorPriKey:  []*ecdsa.PrivateKey{deployPrivateKey},
		InitPowers:       initPowers,
	}

	for i, power := range para.InitPowers {
		relayerLog.Info("deploy", "the validator address ", para.InitValidators[i].String(),
			                                   "power", power.String())
	}


	x2EthContracts, x2EthDeployInfo, err := ethtxs.DeployAndInit(ethRelayer.client, para)
	if err != nil {
		return bridgeRegistry, err
	}
	ethRelayer.operatorInfo = &ethtxs.OperatorInfo{
		PrivateKey:deployPrivateKey,
		Address:deployerAddr,
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
	return ethtxs.GetBalance(ethRelayer.client, tokenAddr, owner)
}

func (ethRelayer *EthereumRelayer) ShowBridgeBankAddr() (string, error) {
	if nil == ethRelayer.x2EthDeployInfo {
		return "", errors.New("The relayer is not started yes")
	}

	return ethRelayer.x2EthDeployInfo.BridgeBank.Address.String(), nil
}

func (ethRelayer *EthereumRelayer) IsProphecyPending(claimID [32]byte) (bool, error) {
	return ethtxs.IsProphecyPending(claimID, ethRelayer.ethValidator, ethRelayer.x2EthContracts.Chain33Bridge)
}

func (ethRelayer *EthereumRelayer) MakeNewProphecyClaim(newProphecyClaimPara *ethtxs.NewProphecyClaimPara) (string, error) {
	return ethtxs.MakeNewProphecyClaim(newProphecyClaimPara, ethRelayer.client, ethRelayer.privateKey4Ethereum, ethRelayer.ethValidator, ethRelayer.x2EthContracts)
}

func (ethRelayer *EthereumRelayer) CreateBridgeToken(symbol string) (string, error) {
	return ethtxs.CreateBridgeToken(symbol, ethRelayer.client, ethRelayer.operatorInfo, ethRelayer.x2EthDeployInfo, ethRelayer.x2EthContracts)
}

func (ethRelayer *EthereumRelayer) CreateERC20Token(symbol string) (string, error) {
	return ethtxs.CreateERC20Token(symbol, ethRelayer.client, ethRelayer.operatorInfo, ethRelayer.x2EthDeployInfo, ethRelayer.x2EthContracts)
}

func (ethRelayer *EthereumRelayer) MintERC20Token(tokenAddr, ownerAddr string, amount int64) (string, error) {
	return ethtxs.MintERC20Token(tokenAddr, ownerAddr, amount, ethRelayer.client, ethRelayer.operatorInfo)
}

func (ethRelayer *EthereumRelayer) ApproveAllowance(ownerPrivateKey, tokenAddr string, amount int64) (string, error) {
	return ethtxs.ApproveAllowance(ownerPrivateKey, tokenAddr, ethRelayer.x2EthDeployInfo.BridgeBank.Address, amount, ethRelayer.client)
}

func (ethRelayer *EthereumRelayer) Burn(ownerPrivateKey, tokenAddr, chain33Receiver string, amount int64) (string, error) {
	return ethtxs.Burn(ownerPrivateKey, tokenAddr, chain33Receiver, ethRelayer.x2EthDeployInfo.BridgeBank.Address, amount, ethRelayer.x2EthContracts.BridgeBank, ethRelayer.client)
}

func (ethRelayer *EthereumRelayer) TransferToken(tokenAddr, fromKey, toAddr string, amount int64) (string, error) {
	return ethtxs.TransferToken(tokenAddr, fromKey, toAddr, amount, ethRelayer.client)
}

func (ethRelayer *EthereumRelayer) LockEthErc20Asset(ownerPrivateKey, tokenAddr string, amount int64, chain33Receiver string) (string, error) {
	return ethtxs.LockEthErc20Asset(ownerPrivateKey, tokenAddr, chain33Receiver, amount, ethRelayer.client, ethRelayer.x2EthContracts.BridgeBank)
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
	relayerLog.Info("Please unlock or import private key for Ethereum relayer")
	nilAddr := common.Address{}
	if nilAddr != ethRelayer.bridgeRegistryAddr {
		relayerLog.Info("proc", "Going to recover corresponding solidity contract handler with bridgeRegistryAddr", ethRelayer.bridgeRegistryAddr.String())
		ethRelayer.x2EthContracts, ethRelayer.x2EthDeployInfo, err = ethtxs.RecoverContractHandler(client, ethRelayer.bridgeRegistryAddr, ethRelayer.bridgeRegistryAddr)
		if nil != err {
			panic("Failed to recover corresponding solidity contract handler due to:" + err.Error())
		}
		relayerLog.Info("^-^ ^-^ Succeed to recover corresponding solidity contract handler")
		if nil != ethRelayer.recoverDeployPara() {
			panic("Failed to recoverDeployPara")
		}
		ethRelayer.unlockchan <- start
	}

	for {
		select {
		case <-ethRelayer.unlockchan:
			relayerLog.Info("Received ethRelayer.unlockchan")
			if nil != ethRelayer.privateKey4Ethereum && nil != ethRelayer.privateKey4Chain33 && nilAddr != ethRelayer.bridgeRegistryAddr {
				ethRelayer.ethValidator, err = ethtxs.LoadSender(ethRelayer.privateKey4Ethereum)
				if nil != err {
					errinfo := fmt.Sprintf("Failed to load validator for ethereum due to:%s", err.Error())
					panic(errinfo)
				}
				relayerLog.Info("Ethereum relayer starts to run...")
				ethRelayer.prePareSubscribeEvent()
				ethRelayer.filterLogEvents(clientChainID)

				//向chain33Bridge订阅事件
				ethRelayer.subscribeEvent(false)
				//向bridgeBank订阅事件
				ethRelayer.subscribeEvent(true)
				relayerLog.Info("Ethereum relayer starts to process online log event...")
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
		case vLog := <-ethRelayer.bridgeBankLog:
			ethRelayer.procBridgeBankLogs(vLog, clientChainID)
		case vLog := <-ethRelayer.chain33BridgeLog:
			ethRelayer.procChain33BridgeLogs(vLog)
		}
	}
}

func (ethRelayer *EthereumRelayer) procBridgeBankLogs(vLog types.Log, clientChainID *big.Int) {
	if ethRelayer.checkTxProcessed(vLog.TxHash.Bytes()) {
		relayerLog.Info("procBridgeBankLogs", "Tx Processed with hash:", vLog.TxHash.Hex())
		return
	}

	//lock,用于捕捉 (ETH/ERC20----->chain33) 跨链转移
	if vLog.Topics[0].Hex() == ethRelayer.bridgeBankEventLockSig {
		eventName := events.LogLock.String()
		relayerLog.Info("EthereumRelayer proc", "^_^ ^_^ Received bridgeBankLog for event", eventName,
			"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
		err := ethRelayer.handleLogLockEvent(clientChainID, ethRelayer.bridgeBankAbi, eventName, vLog)
		if err != nil {
			errinfo := fmt.Sprintf("Failed to handleLogLockEvent due to:%s", err.Error())
			relayerLog.Info("EthereumRelayer procBridgeBankLogs", "errinfo", errinfo)
			//panic(errinfo)
		}
	} else if vLog.Topics[0].Hex() == ethRelayer.bridgeBankEventBurnSig {
		//burn,用于捕捉 (chain33 token----->chain33) 实现chain33资产withdraw操作，之后在chain33上实现unlock操作
		eventName := events.LogChain33TokenBurn.String()
		relayerLog.Info("EthereumRelayer proc", "^_^ ^_^ Received bridgeBankLog for event", eventName,
			"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
		err := ethRelayer.handleLogBurnEvent(clientChainID, ethRelayer.bridgeBankAbi, eventName, vLog)
		if err != nil {
			errinfo := fmt.Sprintf("Failed to handleLogBurnEvent due to:%s", err.Error())
			relayerLog.Info("EthereumRelayer procBridgeBankLogs", "errinfo", errinfo)
			//panic(errinfo)
		}
	}

	if err := ethRelayer.setTxProcessed(vLog.TxHash.Bytes()); nil != err {
		panic(err.Error())
	}

	if err := ethRelayer.setHeight4BridgeBankLogAt(vLog.BlockNumber); nil != err {
		panic(err.Error())
	}
}

// 捕捉NewProphecyClaim(包括lock和burn)，
// NewProphecyClaim：lock 实现chain33 asset ------> ethereum(跨链转移操作)
// NewProphecyClaim：burn 实现ETH/ERC20 ------> ethereum(资产取回操作)
func (ethRelayer *EthereumRelayer) procChain33BridgeLogs(vLog types.Log) {
	if ethRelayer.checkTxProcessed(vLog.TxHash.Bytes()) {
		relayerLog.Info("procChain33BridgeLogs", "Tx Processed with hash:", vLog.TxHash.Hex())
		return
	}
	if vLog.Topics[0].Hex() == ethRelayer.chain33BridgeEventSig {
		eventName := events.LogNewProphecyClaim.String()
		relayerLog.Info("EthereumRelayer proc", "^_^ ^_^ Received chain33BridgeLog with new event:", eventName,
			"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
		err := ethRelayer.handleLogNewProphecyClaimEvent(ethRelayer.chain33BridgeAbi, eventName, vLog)
		if err != nil {
			errinfo := fmt.Sprintf("Failed to handleLogNewProphecyClaimEvent due to:%s", err.Error())
			//panic(errinfo)
			relayerLog.Info("EthereumRelayer procChain33BridgeLogs", "errinfo", errinfo)
		}
	}

	if err := ethRelayer.setTxProcessed(vLog.TxHash.Bytes()); nil != err {
		panic(err.Error())
	}

	if err := ethRelayer.setHeight4chain33BridgeLogAt(vLog.BlockNumber); nil != err {
		panic(err.Error())
	}
}

func (ethRelayer *EthereumRelayer) filterLogEvents(clientChainID *big.Int) {
	//debug code, just for debug now
	ethRelayer.setHeight4BridgeBankLogAt(0)
	ethRelayer.setHeight4chain33BridgeLogAt(0)

	deployHeight := int64(1)
	height4BridgeBankLogAt := int64(ethRelayer.getHeight4BridgeBankLogAt())
	height4chain33BridgeLogAt := int64(ethRelayer.getHeight4chain33BridgeLogAt())

	if height4BridgeBankLogAt < deployHeight {
		height4BridgeBankLogAt = deployHeight
	}
	if height4chain33BridgeLogAt < deployHeight {
		height4chain33BridgeLogAt = deployHeight
	}

	header, err := ethRelayer.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		errinfo := fmt.Sprintf("Failed to get HeaderByNumbers due to:%s", err.Error())
		panic(errinfo)
	}
	curHeight := int64(header.Number.Uint64())
	relayerLog.Info("filterLogEvents", "curHeight:", curHeight)

	if curHeight < height4BridgeBankLogAt && curHeight < height4chain33BridgeLogAt {
		return
	}

	bridgeBankSig := make(map[string]bool)
	bridgeBankSig[ethRelayer.bridgeBankEventLockSig] = true
	bridgeBankSig[ethRelayer.bridgeBankEventBurnSig] = true
	bridgeBankLog := make(chan types.Log)
	done := make(chan int)
	go ethRelayer.filterLogEventsProc(bridgeBankLog, done, "bridgeBank", curHeight, height4BridgeBankLogAt, ethRelayer.bridgeBankAddr, bridgeBankSig)

	chain33BridgeSig := make(map[string]bool)
	chain33BridgeSig[ethRelayer.chain33BridgeEventSig] = true
	chain33BridgeLog := make(chan types.Log)
	go ethRelayer.filterLogEventsProc(chain33BridgeLog, done, "chain33Bridge", curHeight, height4chain33BridgeLogAt, ethRelayer.chain33BridgeAddr, chain33BridgeSig)

	doneCnt := 0
	for {
		select {
		case vLog := <-bridgeBankLog:
			ethRelayer.procBridgeBankLogs(vLog, clientChainID)
		case vLog := <-chain33BridgeLog:
			ethRelayer.procChain33BridgeLogs(vLog)
		case <-done:
			doneCnt++
			if 2 == doneCnt {
				relayerLog.Info("Finshed offline logs processed")
				return
			}
		}
	}
	relayerLog.Info("Finshed offline logs processed")
}

func (ethRelayer *EthereumRelayer) filterLogEventsProc(logchan chan<- types.Log, done chan<- int, title string, curHeight, heightLogProcAt int64, contractAddr common.Address, eventSig map[string]bool) {
	relayerLog.Info(title, "eventSig", eventSig)

	startHeight := heightLogProcAt
	batchCount := int64(10)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddr},
	}

	for {
		if batchCount < (curHeight - startHeight + 1) {
			stopHeight := startHeight + batchCount - 1
			query.FromBlock = big.NewInt(startHeight)
			query.ToBlock = big.NewInt(stopHeight)
		} else {
			query.FromBlock = big.NewInt(startHeight)
			query.ToBlock = big.NewInt(curHeight)
		}

		// Filter by contract and event, write results to logs
		logs, err := ethRelayer.client.FilterLogs(context.Background(), query)
		if err != nil {
			errinfo := fmt.Sprintf("Failed to filterLogEvents due to:%s", err.Error())
			panic(errinfo)
		}

		relayerLog.Info(title, "received logs", len(logs))
		for _, log := range logs {
			relayerLog.Info(title, "received log with topics", log.Topics[0].Hex(), "BlockNumber", log.BlockNumber)
			if _, exist := eventSig[log.Topics[0].Hex()]; !exist {
				continue
			}
			//if ethRelayer.checkTxProcessed(log.TxHash.Bytes()) {
			//	relayerLog.Info("filterLogEvents", "get processed log with topic:", log.Topics[0].String(),
			//		"BlockNumber", log.BlockNumber)
			//	continue
			//}
			logchan <- log
			relayerLog.Info(title, "get unprocessed log with topic:", log.Topics[0].String(),
				"BlockNumber", log.BlockNumber)
		}

		if query.ToBlock.Int64() == curHeight {
			relayerLog.Info(title,  "Finished FilterLogs to height", curHeight)
			done <- 1
			break
		}
		startHeight = query.ToBlock.Int64() + 1
	}
}

func (ethRelayer *EthereumRelayer) prePareSubscribeEvent() {
	var eventName string
	//chain33Bridge处理
	eventName = events.LogNewProphecyClaim.String()
	contactAbi := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
	ethRelayer.chain33BridgeEventSig = contactAbi.Events[eventName].ID().Hex()
	ethRelayer.chain33BridgeAbi = contactAbi
	ethRelayer.chain33BridgeAddr = ethRelayer.x2EthDeployInfo.Chain33Bridge.Address

	//bridgeBank处理
	contactAbi = ethtxs.LoadABI(ethtxs.BridgeBankABI)
	ethRelayer.bridgeBankAbi = contactAbi
	eventName = events.LogLock.String()
	ethRelayer.bridgeBankEventLockSig = contactAbi.Events[eventName].ID().Hex()
	eventName = events.LogChain33TokenBurn.String()
	ethRelayer.bridgeBankEventBurnSig = contactAbi.Events[eventName].ID().Hex()
    ethRelayer.bridgeBankAddr = ethRelayer.x2EthDeployInfo.BridgeBank.Address
}

func (ethRelayer *EthereumRelayer) subscribeEvent(makeClaims bool) {
	var targetAddress common.Address
	if makeClaims {
		targetAddress = ethRelayer.chain33BridgeAddr

	} else {
		targetAddress = ethRelayer.bridgeBankAddr
	}

	// We need the target address in bytes[] for the query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{targetAddress},
		FromBlock: big.NewInt(int64(1)),
	}
	// We will check logs for new events
	logs := make(chan types.Log, 10)
	// Filter by contract and event, write results to logs
	sub, err := ethRelayer.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		errinfo := fmt.Sprintf("Failed to SubscribeFilterLogs due to:%s", err.Error())
		panic(errinfo)
	}
	relayerLog.Info("subscribeEvent", "Subscribed to contract at address:", targetAddress.Hex())

	if makeClaims {
		ethRelayer.chain33BridgeLog = logs
		ethRelayer.chain33BridgeSub = sub
		return
	}
	ethRelayer.bridgeBankLog = logs
	ethRelayer.bridgeBankSub = sub
	return
}

func (ethRelayer *EthereumRelayer) IsValidatorActive(addr string) (bool, error) {
	return ethtxs.IsActiveValidator(common.HexToAddress(addr), ethRelayer.x2EthContracts.Valset)
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
	rpcURL := ethRelayer.rpcURL2Chain33

	// Unpack the LogLock event using its unique event signature from the contract's ABI
	event, err := events.UnpackLogLock(contractABI, eventName, log.Data)
	if nil != err {
		return err
	}
	// Add the event to the record
	events.NewEventWrite(log.TxHash.Hex(), *event)

	// Parse the LogLock event's payload into a struct
	prophecyClaim, err := ethtxs.LogLockToEthBridgeClaim(event, clientChainID.Int64(), ethRelayer.bridgeBankAddr.String())
	if err != nil {
		return err
	}

	// Initiate the relay
	txhash, err := ethtxs.RelayLockToChain33(ethRelayer.privateKey4Chain33, prophecyClaim, rpcURL)
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

// handleLogBurnEvent : unpacks a burn event, converts it to a ProphecyClaim, and relays a tx to chain33
func (ethRelayer *EthereumRelayer) handleLogBurnEvent(clientChainID *big.Int, contractABI abi.ABI, eventName string, log types.Log) error {
	rpcURL := ethRelayer.rpcURL2Chain33

	event, err := events.UnpackLogBurn(contractABI, eventName, log.Data)
	if nil != err {
		return err
	}

	// Parse the LogLock event's payload into a struct
	prophecyClaim, err := ethtxs.LogBurnToEthBridgeClaim(event, clientChainID.Int64(), ethRelayer.bridgeBankAddr.String())
	if err != nil {
		return err
	}

	// Initiate the relay
	txhash, err := ethtxs.RelayBurnToChain33(ethRelayer.privateKey4Chain33, prophecyClaim, rpcURL)
	if err != nil {
		relayerLog.Error("handleLogLockEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}

	//保存交易hash，方便查询
	atomic.AddInt64(&ethRelayer.totalTx4Eth2Chain33, 1)
	txIndex := atomic.LoadInt64(&ethRelayer.totalTx4Eth2Chain33)
	if err = ethRelayer.updateTotalTxAmount2chain33(txIndex); nil != err {
		relayerLog.Error("handleLogBurnEvent", "Failed to RelayLockToChain33 due to:", err.Error())
		return err
	}
	if err = ethRelayer.setLastestRelay2Chain33Txhash(txhash, txIndex); nil != err {
		relayerLog.Error("handleLogBurnEvent", "Failed to RelayLockToChain33 due to:", err.Error())
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
