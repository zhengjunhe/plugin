package events

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Chain33Msg : contains data from MsgBurn and MsgLock events
type Chain33Msg struct {
	ClaimType            Event
	Chain33Sender        []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}

// NewChain33Msg : creates a new Chain33Msg
func NewChain33Msg(
	claimType Event,
	chain33Sender []byte,
	ethereumReceiver common.Address,
	symbol string,
	amount *big.Int,
	tokenContractAddress common.Address,
) Chain33Msg {
	// Package data into a Chain33Msg
	chain33Msg := Chain33Msg{
		ClaimType:            claimType,
		Chain33Sender:        chain33Sender,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
		TokenContractAddress: tokenContractAddress,
	}

	PrintChain33Msg(chain33Msg)

	return chain33Msg
}

// PrintChain33Msg : prints a Chain33Msg struct's information
func PrintChain33Msg(event Chain33Msg) {
	claimType := event.ClaimType.String()
	chain33Sender := string(event.Chain33Sender)
	ethereumReceiver := event.EthereumReceiver.Hex()
	tokenContractAddress := event.TokenContractAddress.Hex()
	symbol := event.Symbol
	amount := event.Amount

	fmt.Printf("\nClaim Type: %v\nChain33 Sender: %v\nEthereum Recipient: %v\nToken Address: %v\nSymbol: %v\nAmount: %v\n",
		claimType, chain33Sender, ethereumReceiver, tokenContractAddress, symbol, amount)
}
