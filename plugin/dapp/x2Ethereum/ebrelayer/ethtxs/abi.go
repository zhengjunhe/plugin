package ethtxs

import (
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

const (
	BridgeBankABI = "BridgeBankABI"
	CosmosBankABI = "CosmosBankABI"
	CosmosBridgeABI = "CosmosBridgeABI"
	EthereumBankABI = "EthereumBankABI"
)

func LoadABI(contractName string) abi.ABI {
	var abiJson string
	switch contractName {
	case BridgeBankABI:
		abiJson = generated.BridgeBankABI
	case CosmosBankABI:
		abiJson = generated.CosmosBankABI
	case CosmosBridgeABI:
		abiJson = generated.CosmosBridgeABI
	case EthereumBankABI:
		abiJson = generated.EthereumBankABI
	default:
		panic("No abi matched")
	}

	// Convert the raw abi into a usable format
	contractABI, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(err)
	}

	return contractABI
}
