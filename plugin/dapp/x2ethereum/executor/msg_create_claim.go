package executor

import (
	"errors"
	"fmt"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"strings"
)

// MsgCreateEthBridgeClaim defines a message for creating claims on the ethereum bridge
type MsgCreateEthBridgeClaim types.EthBridgeClaim

// NewMsgCreateEthBridgeClaim is a constructor function for MsgCreateBridgeClaim
func NewMsgCreateEthBridgeClaim(ethBridgeClaim types.EthBridgeClaim) MsgCreateEthBridgeClaim {
	return MsgCreateEthBridgeClaim(ethBridgeClaim)
}

// Route should return the name of the module
func (msg MsgCreateEthBridgeClaim) Route() string { return ModuleName }

// Type should return the action
func (msg MsgCreateEthBridgeClaim) Type() string { return "create_bridge_claim" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreateEthBridgeClaim) ValidateBasic() error {
	if AddressIsEmpty(msg.Chain33Receiver) {
		return types.ErrInvalidAddress
	}

	if AddressIsEmpty(msg.ValidatorAddress) {
		return types.ErrInvalidAddress
	}

	if msg.Nonce < 0 {
		return types.ErrInvalidEthNonce
	}

	if !gethCommon.IsHexAddress(msg.EthereumSender) {
		return types.ErrInvalidEthAddress
	}
	if !gethCommon.IsHexAddress(msg.BridgeContractAddress) {
		return types.ErrInvalidEthAddress
	}
	if strings.ToLower(msg.Symbol) == "eth" && msg.TokenContractAddress != "0x0000000000000000000000000000000000000000" {
		return types.ErrInvalidEthSymbol
	}
	return nil
}

// MapOracleClaimsToEthBridgeClaims maps a set of generic oracle claim data into EthBridgeClaim objects
func MapOracleClaimsToEthBridgeClaims(ethereumChainID int, bridgeContract string, nonce int, symbol string, tokenContract string, ethereumSender string, oracleValidatorClaims map[string]string, f func(int, string, int, string, string, string, string, string) (types.EthBridgeClaim, error)) ([]types.EthBridgeClaim, error) {
	mappedClaims := make([]types.EthBridgeClaim, len(oracleValidatorClaims))
	i := 0
	for validator, validatorClaim := range oracleValidatorClaims {
		parseErr := address.CheckAddress(validator)
		if parseErr != nil {
			return nil, errors.New(fmt.Sprintf("failed to parse claim: %s", parseErr))
		}
		mappedClaim, err := f(ethereumChainID, bridgeContract, nonce, symbol, tokenContract, ethereumSender, validator, validatorClaim)
		if err != nil {
			return nil, err
		}
		mappedClaims[i] = mappedClaim
		i++
	}
	return mappedClaims, nil
}
