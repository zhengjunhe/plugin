package executor

import (
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"strconv"
)

// MsgLock defines a message for locking coins and triggering a related event
type MsgLock struct {
	EthereumChainID  int                     `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	TokenContract    executor.EthAddress     `json:"token_contract_address" yaml:"token_contract_address"`
	Chain33Sender    executor.Chain33Address `json:"chain33_sender" yaml:"chain33_sender"`
	EthereumReceiver executor.EthAddress     `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	Amount           uint64                  `json:"amount" yaml:"amount"`
}

// NewMsgLock is a constructor function for MsgLock
func NewMsgLock(ethereumChainID int, tokenContract string, cosmosSender string, ethereumReceiver string, amount uint64) MsgLock {
	return MsgLock{
		EthereumChainID:  ethereumChainID,
		TokenContract:    executor.NewEthereumAddress(tokenContract),
		Chain33Sender:    executor.NewChain33Address(cosmosSender),
		EthereumReceiver: executor.NewEthereumAddress(ethereumReceiver),
		Amount:           amount,
	}
}

// Route should return the name of the module
func (msg MsgLock) Route() string { return executor.ModuleName }

// Type should return the action
func (msg MsgLock) Type() string { return "lock" }

// ValidateBasic runs stateless checks on the message
func (msg MsgLock) ValidateBasic() error {
	if strconv.Itoa(msg.EthereumChainID) == "" {
		return types.ErrInvalidChainID
	}

	if msg.TokenContract.String() == "" {
		return types.ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.TokenContract.String()) {
		return types.ErrInvalidEthAddress
	}

	if AddressIsEmpty(msg.Chain33Sender.Enc58str) {
		return types.ErrInvalidAddress
	}

	if msg.EthereumReceiver.String() == "" {
		return types.ErrInvalidEthAddress
	}

	if !gethCommon.IsHexAddress(msg.EthereumReceiver.String()) {
		return types.ErrInvalidEthAddress
	}

	return nil
}
