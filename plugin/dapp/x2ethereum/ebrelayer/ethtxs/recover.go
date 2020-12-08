package ethtxs

import (
	"errors"

	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethinterface"
	"github.com/ethereum/go-ethereum/common"
)

//RecoverContractHandler ...
func RecoverContractHandler(client ethinterface.EthClientSpec, sender, registry common.Address) (*X2EthContracts, *X2EthDeployInfo, error) {
	bridgeBankAddr, err := GetAddressFromBridgeRegistry(client, sender, registry, BridgeBank)
	if nil != err {
		return nil, nil, errors.New("failed to get addr for bridgeBank from registry")
	}
	bridgeBank, err := generated.NewBridgeBank(*bridgeBankAddr, client)
	if nil != err {
		return nil, nil, errors.New("failed to NewBridgeBank")
	}

	dplatformBridgeAddr, err := GetAddressFromBridgeRegistry(client, sender, registry, DplatformBridge)
	if nil != err {
		return nil, nil, errors.New("failed to get addr for dplatformBridgeAddr from registry")
	}
	dplatformBridge, err := generated.NewDplatformBridge(*dplatformBridgeAddr, client)
	if nil != err {
		return nil, nil, errors.New("failed to NewDplatformBridge")
	}

	oracleAddr, err := GetAddressFromBridgeRegistry(client, sender, registry, Oracle)
	if nil != err {
		return nil, nil, errors.New("failed to get addr for oracleBridgeAddr from registry")
	}
	oracle, err := generated.NewOracle(*oracleAddr, client)
	if nil != err {
		return nil, nil, errors.New("failed to NewOracle")
	}

	valsetAddr, err := GetAddressFromBridgeRegistry(client, sender, registry, Valset)
	if nil != err {
		return nil, nil, errors.New("failed to get addr for valset from registry")
	}
	valset, err := generated.NewValset(*valsetAddr, client)
	if nil != err {
		return nil, nil, errors.New("failed to NewValset")
	}

	registryInstance, _ := generated.NewBridgeRegistry(registry, client)
	x2EthContracts := &X2EthContracts{
		BridgeRegistry: registryInstance,
		BridgeBank:     bridgeBank,
		DplatformBridge:  dplatformBridge,
		Oracle:         oracle,
		Valset:         valset,
	}

	x2EthDeployInfo := &X2EthDeployInfo{
		BridgeRegistry: &DeployResult{Address: registry},
		BridgeBank:     &DeployResult{Address: *bridgeBankAddr},
		DplatformBridge:  &DeployResult{Address: *dplatformBridgeAddr},
		Oracle:         &DeployResult{Address: *oracleAddr},
		Valset:         &DeployResult{Address: *valsetAddr},
	}

	return x2EthContracts, x2EthDeployInfo, nil
}

//RecoverOracleInstance ...
func RecoverOracleInstance(client ethinterface.EthClientSpec, sender, registry common.Address) (*generated.Oracle, error) {
	oracleAddr, err := GetAddressFromBridgeRegistry(client, sender, registry, Oracle)
	if nil != err {
		return nil, errors.New("failed to get addr for oracleBridgeAddr from registry")
	}
	oracle, err := generated.NewOracle(*oracleAddr, client)
	if nil != err {
		return nil, errors.New("failed to NewOracle")
	}

	return oracle, nil
}
