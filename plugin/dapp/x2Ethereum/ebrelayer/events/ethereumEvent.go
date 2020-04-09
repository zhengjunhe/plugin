package events

// -----------------------------------------------------
//    ethereumEvent : Creates LockEvents from new events on the
//			  Ethereum blockchain.
// -----------------------------------------------------

import (
	"fmt"
	"math/big"

	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// LockEvent : struct which represents a LogLock event
type LockEvent struct {
	EthereumChainID       *big.Int
	BridgeContractAddress common.Address
	Id                    [32]byte
	From                  common.Address
	To                    []byte
	Token                 common.Address
	Symbol                string
	Value                 *big.Int
	Nonce                 *big.Int
}

// NewProphecyClaimEvent : struct which represents a LogNewProphecyClaim event
type NewProphecyClaimEvent struct {
	ProphecyID       *big.Int
	ClaimType        uint8
	Chain33Sender    []byte
	EthereumReceiver common.Address
	ValidatorAddress common.Address
	TokenAddress     common.Address
	Symbol           string
	Amount           *big.Int
}

type LogNewBridgeToken struct {
	Token  common.Address
	Symbol string
}

// UnpackLogLock : Handles new LogLock events
func UnpackLogLock(clientChainID *big.Int, contractAddress string, contractAbi abi.ABI, eventName string, eventData []byte) (lockEvent *LockEvent, err error) {
	event := &LockEvent{}

	// Bridge contract address
	if !common.IsHexAddress(contractAddress) {
		eventsLog.Error("UnpackLogLock", "Only Ethereum contracts are currently supported. Invalid address: ", contractAddress)
		return nil, ebrelayerTypes.ErrInvalidEthContractAddress
	}
	event.BridgeContractAddress = common.HexToAddress(contractAddress)

	// Ethereum chain ID
	event.EthereumChainID = clientChainID

	// Parse the event's attributes as Ethereum network variables
	err = contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		eventsLog.Error("UnpackLogLock", "Failed to unpack abi due to:", err.Error())
		return nil, ebrelayerTypes.ErrUnpack
	}
	info2Log := printLockEvent(event)
	eventsLog.Error("UnpackLogLock","detail info of event:", info2Log)

	return event, nil
}

// UnpackLogNewProphecyClaim : Handles new LogNewProphecyClaim events
func UnpackLogNewProphecyClaim(contractAbi abi.ABI, eventName string, eventData []byte) (newProphecyClaimEvent NewProphecyClaimEvent, err error) {
	event := NewProphecyClaimEvent{}

	// Parse the event's attributes as Ethereum network variables
	err = contractAbi.Unpack(&event, eventName, eventData)
	if err != nil {
		eventsLog.Error("UnpackLogNewProphecyClaim", "Failed to UnpackLogNewProphecyClaim due to:", err.Error(),
			"eventName", eventName)
		return
	}

	eventsLog.Info("UnpackLogNewProphecyClaim", "event detailed info:", printProphecyClaimEvent(event))
	return
}

// printLockEvent : prints a LockEvent struct's information
func printLockEvent(event *LockEvent) string {
	// Convert the variables into a printable format
	chainID := event.EthereumChainID
	bridgeContractAddress := event.BridgeContractAddress.Hex()
	value := event.Value
	symbol := event.Symbol
	token := event.Token.Hex()
	sender := event.From.Hex()
	recipient := string(event.To)
	nonce := event.Nonce

	// Print the event's information
	return fmt.Sprintf("\nChain ID: %v\nBridge contract address: %v\nToken symbol: %v\nToken contract address: %v\nSender: %v\nRecipient: %v\nValue: %v\nNonce: %v\n\n",
		chainID, bridgeContractAddress, symbol, token, sender, recipient, value, nonce)
}

// printProphecyClaimEvent : prints a NewProphecyClaimEvent struct's information
func printProphecyClaimEvent(event NewProphecyClaimEvent) string {
	// Convert the variables into a printable format
	id := event.ProphecyID
	claimType := event.ClaimType
	sender := string(event.Chain33Sender)
	recipient := event.EthereumReceiver.Hex()
	symbol := event.Symbol
	token := event.TokenAddress.Hex()
	amount := event.Amount
	validator := event.ValidatorAddress.Hex()

	// Print the event's information
	return fmt.Sprintf("\nProphecy ID: %v\nClaim Type: %v\nSender: %v\nRecipient %v\nSymbol %v\nToken %v\nAmount: %v\nValidator: %v\n\n",
		id, claimType, sender, recipient, symbol, token, amount, validator)
}
