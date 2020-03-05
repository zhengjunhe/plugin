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
	cosmosMsg := Chain33Msg{
		ClaimType:            claimType,
		Chain33Sender:        chain33Sender,
		EthereumReceiver:     ethereumReceiver,
		Symbol:               symbol,
		Amount:               amount,
		TokenContractAddress: tokenContractAddress,
	}

	PrintCosmosMsg(cosmosMsg)

	return cosmosMsg
}

// PrintCosmosMsg : prints a Chain33Msg struct's information
func PrintCosmosMsg(event Chain33Msg) {
	claimType := event.ClaimType.String()
	cosmosSender := string(event.Chain33Sender)
	ethereumReceiver := event.EthereumReceiver.Hex()
	tokenContractAddress := event.TokenContractAddress.Hex()
	symbol := event.Symbol
	amount := event.Amount

	fmt.Printf("\nClaim Type: %v\nCosmos Sender: %v\nEthereum Recipient: %v\nToken Address: %v\nSymbol: %v\nAmount: %v\n",
		claimType, cosmosSender, ethereumReceiver, tokenContractAddress, symbol, amount)
}
