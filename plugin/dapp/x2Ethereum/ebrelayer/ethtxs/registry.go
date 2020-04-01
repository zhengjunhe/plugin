package ethtxs

import (
	"context"
	"log"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	bridgeRegistry "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
    ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
)

// ContractRegistry :
type ContractRegistry byte

const (
	// Valset : valset contract
	Valset ContractRegistry = iota + 1
	// Oracle : oracle contract
	Oracle
	// BridgeBank : bridgeBank contract
	BridgeBank
	// CosmosBridge : cosmosBridge contract
	CosmosBridge
)

// String : returns the event type as a string
func (d ContractRegistry) String() string {
	return [...]string{"valset", "oracle", "bridgebank", "cosmosbridge"}[d-1]
}

// GetAddressFromBridgeRegistry : utility method which queries the requested contract address from the BridgeRegistry
func GetAddressFromBridgeRegistry(client *ethclient.Client, sender, registry common.Address, target ContractRegistry) (address *common.Address, err error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		txslog.Error("GetAddressFromBridgeRegistry", "Failed to get HeaderByNumber due to:", err.Error())
		return nil, err
	}

	// Set up CallOpts auth
	auth := bind.CallOpts{
		Pending:     true,
		From:        sender,
		BlockNumber: header.Number,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	registryInstance, err := bridgeRegistry.NewBridgeRegistry(registry, client)
	if err != nil {
		txslog.Error("GetAddressFromBridgeRegistry", "Failed to NewBridgeRegistry to:", err.Error())
		return nil, err
	}

	switch target {
	case Valset:
		valsetAddress, err := registryInstance.Valset(&auth)
		if err != nil {
			log.Fatal(err)
		}
		return &valsetAddress, nil
	case Oracle:
		oracleAddress, err := registryInstance.Oracle(&auth)
		if err != nil {
			log.Fatal(err)
		}
		return &oracleAddress, nil
	case BridgeBank:
		bridgeBankAddress, err := registryInstance.BridgeBank(&auth)
		if err != nil {
			log.Fatal(err)
		}
		return &bridgeBankAddress, nil
	case CosmosBridge:
		cosmosBridgeAddress, err := registryInstance.CosmosBridge(&auth)
		if err != nil {
			log.Fatal(err)
		}
		return &cosmosBridgeAddress, nil
	default:
		txslog.Error("GetAddressFromBridgeRegistry", "invalid target contract type:", target)
		return nil, ebrelayerTypes.ErrInvalidContractAddress
	}
}
