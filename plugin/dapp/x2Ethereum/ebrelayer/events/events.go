package events

import (
	log "github.com/33cn/chain33/common/log/log15"
)
// Event : enum containing supported contract events
type Event int

var eventsLog = log.New("module", "ethereum_relayer")

const (
	// Unsupported : unsupported Chain33 or Ethereum event
	Unsupported Event = iota
	// MsgBurn : Chain33 event 'Chain33Msg' type MsgBurn
	MsgBurn
	// MsgLock :  Chain33 event 'Chain33Msg' type MsgLock
	MsgLock
	// LogLock : Ethereum event 'LockEvent'
	LogLock
	// LogNewProphecyClaim : Ethereum event 'NewProphecyClaimEvent'
	LogNewProphecyClaim
)

const (
	CLAIM_TYPE_BURN = uint8(1)
	CLAIM_TYPE_LOCK = uint8(2)
)

// String : returns the event type as a string
func (d Event) String() string {
	return [...]string{"unsupported", "burn", "lock", "LogLock", "LogNewProphecyClaim"}[d]
}

// Chain33MsgAttributeKey : enum containing supported attribute keys
type Chain33MsgAttributeKey int

const (
	// UnsupportedAttributeKey : unsupported attribute key
	UnsupportedAttributeKey Chain33MsgAttributeKey = iota
	// Chain33Sender : sender's address on Chain33 network
	Chain33Sender
	// EthereumReceiver : receiver's address on Ethereum network
	EthereumReceiver
	// Coin : coin type
	Coin
	// TokenContractAddress : coin's corresponding contract address deployed on the Ethereum network
	TokenContractAddress
)

// String : returns the event type as a string
func (d Chain33MsgAttributeKey) String() string {
	return [...]string{"unsupported", "chain33_sender", "ethereum_receiver", "amount", "token_contract_address"}[d]
}
