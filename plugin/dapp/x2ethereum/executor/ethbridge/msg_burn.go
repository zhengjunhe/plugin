package executor

import (
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"strconv"
)

type Msg_Burn struct {
	EthereumChainID  int64                   `json:"ethereum_chain_id" yaml:"ethereum_chain_id"`
	TokenContract    executor.EthAddress     `json:"token_contract_address" yaml:"token_contract_address"`
	Chain33Sender    executor.Chain33Address `json:"chain33_sender" yaml:"chain33_sender"`
	EthereumReceiver executor.EthAddress     `json:"ethereum_receiver" yaml:"ethereum_receiver"`
	Amount           uint64                  `json:"amount" yaml:"amount"`
}

// NewMsgBurn is a constructor function for MsgBurn
func NewMsgBurn(ethereumChainID int64, tokenContract string, chain33Sender string, ethereumReceiver string, amount uint64) Msg_Burn {
	return Msg_Burn{
		EthereumChainID:  ethereumChainID,
		TokenContract:    executor.NewEthereumAddress(tokenContract),
		Chain33Sender:    executor.NewChain33Address(chain33Sender),
		EthereumReceiver: executor.NewEthereumAddress(ethereumReceiver),
		Amount:           amount,
	}
}

// Route should return the name of the module
func (msg Msg_Burn) Route() string { return executor.ModuleName }

// Type should return the action
func (msg Msg_Burn) Type() string { return "burn" }

// ValidateBasic runs stateless checks on the message
func (msg Msg_Burn) ValidateBasic() error {
	if strconv.Itoa(int(msg.EthereumChainID)) == "" {
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

func AddressIsEmpty(address string) bool {
	if address == "" {
		return true
	}

	var aa2 string
	return address == aa2
}
