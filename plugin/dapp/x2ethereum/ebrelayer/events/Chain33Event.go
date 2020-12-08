package events

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// DplatformMsg : contains data from MsgBurn and MsgLock events
type DplatformMsg struct {
	ClaimType            Event
	DplatformSender        []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}

// NewDplatformMsg : creates a new DplatformMsg
func NewDplatformMsg(
	claimType Event,
	dplatformSender []byte,
	ethereumReceiver common.Address,
	symbol string,
	amount *big.Int,
	tokenContractAddress common.Address,
) DplatformMsg {
	// Package data into a DplatformMsg
	dplatformMsg := DplatformMsg{
		ClaimType:            claimType,
		DplatformSender:        dplatformSender,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
		TokenContractAddress: tokenContractAddress,
	}

	return dplatformMsg
}
