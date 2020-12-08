package events

import (
	log "github.com/33cn/dplatform/common/log/log15"
)

// Event : enum containing supported contract events
type Event int

var eventsLog = log.New("module", "ethereum_relayer")

const (
	// Unsupported : unsupported Dplatform or Ethereum event
	Unsupported Event = iota
	// MsgBurn : Dplatform event 'DplatformMsg' type MsgBurn
	MsgBurn
	// MsgLock :  Dplatform event 'DplatformMsg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogDplatformTokenBurn : Ethereum event 'LogDplatformTokenBurn' in contract dplatformBank
	LogDplatformTokenBurn
	// LogNewProphecyClaim : Ethereum event 'NewProphecyClaimEvent'
	LogNewProphecyClaim
)

//const
const (
	ClaimTypeBurn = uint8(1)
	ClaimTypeLock = uint8(2)
)

// String : returns the event type as a string
func (d Event) String() string {
	return [...]string{"unknown-x2ethereum", "DplatformToEthBurn", "DplatformToEthLock", "LogLock", "LogDplatformTokenBurn", "LogNewProphecyClaim"}[d]
}

// DplatformMsgAttributeKey : enum containing supported attribute keys
type DplatformMsgAttributeKey int

const (
	// UnsupportedAttributeKey : unsupported attribute key
	UnsupportedAttributeKey DplatformMsgAttributeKey = iota
	// DplatformSender : sender's address on Dplatform network
	DplatformSender
	// EthereumReceiver : receiver's address on Ethereum network
	EthereumReceiver
	// Coin : coin type
	Coin
	// TokenContractAddress : coin's corresponding contract address deployed on the Ethereum network
	TokenContractAddress
)

// String : returns the event type as a string
func (d DplatformMsgAttributeKey) String() string {
	return [...]string{"unsupported", "dplatform_sender", "ethereum_receiver", "amount", "token_contract_address"}[d]
}
