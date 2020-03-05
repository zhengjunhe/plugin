package ethereum

// -----------------------------------------------------
//      Ethereum relayer
//
//      Initializes the relayer service, which parses,
//      encodes, and packages named events on an Ethereum
//      Smart Contract for validator's to sign and send
//      to the Cosmos bridge.
// -----------------------------------------------------

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	//"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer"
	"math/big"
	"os"
	"sync"
	"sync/atomic"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	chain33Crypto "github.com/33cn/chain33/common/crypto"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/contract"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/txs"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
)

type EthereumRelayer struct {
	provider             string
	contractAddress      common.Address
	//validatorName        string
	db                   dbm.DB
	//passphase            string
	rwLock               sync.RWMutex
	validatorAddress     []byte
	privateKey4Chain33   chain33Crypto.PrivKey
	privateKey4Ethereum  *ecdsa.PrivateKey
	ethSender            common.Address
	totalTx4Eth2Chain33  int64
	totalTx4Chain33ToEth int64
	rpcURL2Chain33       string
	unlockchan           chan int
}

var (
	relayerLog = log.New("module", "ethereum_relayer")
)

func StartEthereumRelayer(rpcURL2Chain33 string, db dbm.DB) *EthereumRelayer {
	relayer := &EthereumRelayer{
		db:         db,
		unlockchan: make(chan int),
		rpcURL2Chain33:rpcURL2Chain33,
	}

	go relayer.proc()
	return relayer
}

func (ethRelayer *EthereumRelayer) proc() {
	// Start client with infura ropsten provider
	client, err := setupWebsocketEthClient(ethRelayer.provider)
	if err != nil {
		panic(err)
	}
	relayerLog.Info("EthereumRelayer proc", "\nStarted Ethereum websocket with provider:", ethRelayer.provider)

	var targetContract txs.ContractRegistry
	var eventName string

	targetContract = txs.BridgeBank
	eventName = events.LogLock.String()

	clientChainID, err := client.NetworkID(context.Background())
	if err != nil {
		errinfo := fmt.Sprintf("Failed to get NetworkID due to:%s", err.Error())
		panic(errinfo)
	}
	// Load our contract's ABI
	contractABI := contract.LoadABI(false)

	// Load unique event signature from the named event contained within the contract's ABI
	eventSignature := contractABI.Events[eventName].Id().Hex()

	fmt.Fprintln(os.Stdout, "Pls unlock or import private key for Ethereum relayer")
	<-ethRelayer.unlockchan
	fmt.Fprintln(os.Stdout, "Ethereum relayer has been unlocked or private key has been imported")
	//等待解锁获取私钥或者导入私钥

	sender, err := txs.LoadSender(ethRelayer.privateKey4Ethereum)
	if nil != err {
		errinfo := fmt.Sprintf("Failed to load sender due to:%s", err.Error())
		panic(errinfo)
	}
	ethRelayer.ethSender = *sender

	// Get the specific contract's address (CosmosBridge or BridgeBank)
	targetAddress, err := txs.GetAddressFromBridgeRegistry(client, *sender, ethRelayer.contractAddress, targetContract)
	if err != nil {
		errinfo := fmt.Sprintf("Failed to GetAddressFromBridgeRegistry due to:%s", err.Error())
		panic(errinfo)
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
	relayerLog.Info("EthereumRelayer proc", "Subscribed to %v contract:", targetContract, "at address:", targetAddress.Hex())

	for {
		select {
		// Handle any errors
		case err := <-sub.Err():
			panic(err)
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == eventSignature {
				relayerLog.Info("EthereumRelayer proc", "Witnessed new event:", eventName,
					"Block number:", vLog.BlockNumber, "Tx hash:", vLog.TxHash.Hex())
				switch eventName {
				// lock事件为用户在eth侧锁定ETH/ERC20，然后向chain33进行跨链转移
				case events.LogLock.String():
					err := ethRelayer.handleLogLockEvent(clientChainID, contractABI, eventName, vLog)
					if err != nil {
						errinfo := fmt.Sprintf("Failed to handleLogLockEvent due to:%s", err.Error())
						panic(errinfo)
					}
				case events.LogNewProphecyClaim.String():
					err := ethRelayer.handleLogNewProphecyClaimEvent(contractABI, eventName, vLog)
					if err != nil {
						errinfo := fmt.Sprintf("Failed to handleLogNewProphecyClaimEvent due to:%s", err.Error())
						panic(errinfo)
					}
				}
			}
		}
	}
}

//func (ethRelayer *EthereumRelayer) SetPassphase(passphase string) {
//	ethRelayer.rwLock.Lock()
//	ethRelayer.passphase = passphase
//	ethRelayer.rwLock.Unlock()
//}

func (ethRelayer *EthereumRelayer) QueryTxhashRelay2Eth() ebrelayerTypes.Txhashes {
	txhashs := ethRelayer.queryTxhashes([]byte(chain33ToEthTxHashPrefix))
	return ebrelayerTypes.Txhashes{Txhash:txhashs}
}

func (ethRelayer *EthereumRelayer) QueryTxhashRelay2Chain33() *ebrelayerTypes.Txhashes {
	txhashs := ethRelayer.queryTxhashes([]byte(eth2chain33TxHashPrefix))
	return &ebrelayerTypes.Txhashes{Txhash:txhashs}
}

// handleLogLockEvent : unpacks a LogLock event, converts it to a ProphecyClaim, and relays a tx to chain33
func (ethRelayer *EthereumRelayer) handleLogLockEvent(clientChainID *big.Int, contractABI abi.ABI, eventName string, log types.Log) error {
	contractAddress := ethRelayer.contractAddress.Hex()
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
	prophecyClaim, err := txs.LogLockToEthBridgeClaim(validatorAddress, event)
	if err != nil {
		return err
	}

	// Initiate the relay
	txhash, err := txs.RelayLockToChain33(ethRelayer.privateKey4Chain33, &prophecyClaim, rpcURL)
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
	oracleClaim, err := txs.ProphecyClaimToSignedOracleClaim(event, ethRelayer.privateKey4Ethereum)
	if nil != err {
		return err
	}

	// Initiate the relay
	txhash, err := txs.RelayOracleClaimToEthereum(ethRelayer.provider, ethRelayer.ethSender, ethRelayer.contractAddress, events.LogNewProphecyClaim, oracleClaim, ethRelayer.privateKey4Ethereum)

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
