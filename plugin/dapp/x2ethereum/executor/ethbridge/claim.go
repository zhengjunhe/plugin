package ethbridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"strconv"
)

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(ethereumChainID int64, bridgeContract string, nonce int64, symbol string, tokenContact string, ethereumSender string, chain33Receiver string, validator string, amount uint64, claimType int64) types.EthBridgeClaim {
	return types.EthBridgeClaim{
		EthereumChainID:       ethereumChainID,
		BridgeContractAddress: bridgeContract,
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContact,
		EthereumSender:        ethereumSender,
		Chain33Receiver:       chain33Receiver,
		ValidatorAddress:      validator,
		Amount:                amount,
		ClaimType:             claimType,
	}
}

// NewOracleClaimContent is a constructor function for OracleClaim
func NewOracleClaimContent(chain33Receiver string, amount uint64, claimType int64) types.OracleClaimContent {
	return types.OracleClaimContent{
		Chain33Receiver: chain33Receiver,
		Amount:          amount,
		ClaimType:       claimType,
	}
}

// NewClaim returns a new Claim
func NewClaim(id string, validatorAddress string, content string) types.OracleClaim {
	return types.OracleClaim{
		ID:               id,
		ValidatorAddress: validatorAddress,
		Content:          content,
	}
}

// CreateOracleClaimFromEthClaim converts a specific ethereum bridge claim to a general oracle claim to be used by
// the oracle module. The oracle module expects every claim for a particular prophecy to have the same id, so this id
// must be created in a deterministic way that all validators can follow. For this, we use the Nonce an Ethereum Sender provided,
// as all validators will see this same data from the smart contract.
func CreateOracleClaimFromEthClaim(ethClaim types.EthBridgeClaim) (types.OracleClaim, error) {
	oracleID := strconv.Itoa(int(ethClaim.EthereumChainID)) + strconv.Itoa(int(ethClaim.Nonce)) + ethClaim.EthereumSender
	claimContent := NewOracleClaimContent(ethClaim.Chain33Receiver, ethClaim.Amount, ethClaim.ClaimType)
	claimBytes, err := json.Marshal(claimContent)
	if err != nil {
		return types.OracleClaim{}, err
	}
	claimString := string(claimBytes)
	claim := NewClaim(oracleID, ethClaim.ValidatorAddress, claimString)
	return claim, nil
}

// CreateEthClaimFromOracleString converts a string from any generic claim from the oracle module into an ethereum bridge specific claim.
func CreateEthClaimFromOracleString(ethereumChainID int64, bridgeContract string, nonce int64, symbol string, tokenContract string, ethereumAddress string, validator string, oracleClaimString string) (types.EthBridgeClaim, error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return types.EthBridgeClaim{}, err
	}

	return NewEthBridgeClaim(
		ethereumChainID,
		bridgeContract,
		nonce,
		symbol,
		tokenContract,
		ethereumAddress,
		oracleClaim.Chain33Receiver,
		validator,
		oracleClaim.Amount,
		oracleClaim.ClaimType,
	), nil
}

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module. In general, it is
// expected that the oracle module will store claims in this JSON format and so this should be used to convert oracle claims.
func CreateOracleClaimFromOracleString(oracleClaimString string) (types.OracleClaimContent, error) {
	var oracleClaimContent types.OracleClaimContent

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaimContent); err != nil {
		return types.OracleClaimContent{}, errors.New(fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaimContent, nil
}
