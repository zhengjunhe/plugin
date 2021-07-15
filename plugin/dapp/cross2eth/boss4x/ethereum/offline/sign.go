package offline

import (
	"crypto/ecdsa"
	eoff "github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum/offline"
	"math/big"
	"strings"
	"time"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4eth/generated"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func signContractTx(contractName string, data []byte, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	rawTx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: big.NewInt(int64(gasPrice)),
		Data:     data,
	})

	//signedtx
	signedtx, hash, err := eoff.SignTx(key, rawTx)
	if err != nil {
		panic(err)
	}
	contractAddress := crypto.CreateAddress(deployerAddr, nonce)
	var contract eoff.DeployContract
	contract.Nonce = nonce
	contract.ContractName = contractName
	contract.SignedRawTx = signedtx
	contract.ContractAddr = contractAddress.String()
	contract.TxHash = hash
	return &contract
}

func signTransactionTx(contractName string, data []byte, deployerAddr, to common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	rawTx := types.NewTx(&types.LegacyTx{
		//To:       &to,
		Nonce:    nonce,
		Value:    big.NewInt(0),
		Gas:      gasLimit,
		GasPrice: big.NewInt(int64(gasPrice)),
		Data:     data,
	})
	//signedtx
	signedtx, hash, err := eoff.SignTx(key, rawTx)
	if err != nil {
		panic(err)
	}
	contractAddress := crypto.CreateAddress(deployerAddr, nonce)
	var contract eoff.DeployContract
	contract.Interval = time.Second * 20
	contract.Nonce = nonce
	contract.ContractName = contractName
	contract.SignedRawTx = signedtx
	contract.ContractAddr = contractAddress.String()
	contract.TxHash = hash
	return &contract
}

func signValSet(validators []common.Address, powers []*big.Int, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.ValsetABI))
	if err != nil {
		panic(err)
	}

	vbin := common.FromHex(generated.ValsetBin)
	input, err := parsed.Pack("", deployerAddr, validators, powers)
	if err != nil {
		panic(err)
	}

	data := append(vbin, input...)
	return signContractTx("valset", data, deployerAddr, nonce, gasLimit, gasPrice, key)
}

func signChain33Bridge(deployerAddr, valSetAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}

	bridgebin := common.FromHex(generated.Chain33BridgeBin)
	input, err := parsed.Pack("", deployerAddr, valSetAddr)
	if err != nil {
		panic(err)
	}

	data := append(bridgebin, input...)
	return signContractTx("chain33bridge", data, deployerAddr, nonce, gasLimit, gasPrice, key)
}

func signOracle(valsetAddr, bridgeAddr common.Address, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.OracleABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.OracleBin)
	input, err := parsed.Pack("", deployerAddr, valsetAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	data := append(bin, input...)
	return signContractTx("oracle", data, deployerAddr, nonce, gasLimit, gasPrice, key)
}

func signBridgeBank(bridgeAddr, oracalAddr common.Address, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.BridgeBankABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeBankBin)
	input, err := parsed.Pack("", deployerAddr, oracalAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	data := append(bin, input...)
	return signContractTx("bridgeBank", data, deployerAddr, nonce, gasLimit, gasPrice, key)
}

//SignSetBridgeBank SetBridgeBank
func signSetBridgeBank(bridgebank, chain33bridge common.Address, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	method := "setBridgeBank"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	input, err := parsed.Pack(method, bridgebank)
	if err != nil {
		panic(err)
	}
	return signTransactionTx("setbridgebank", input, deployerAddr, chain33bridge, nonce, gasLimit, gasPrice, key)
}

func signsetOracle(oracalAddr, chain33bridge common.Address, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	method := "setOracle"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	input, err := parsed.Pack(method, oracalAddr)
	if err != nil {
		panic(err)
	}
	return signTransactionTx("setOracle", input, deployerAddr, chain33bridge, nonce, gasLimit, gasPrice, key)
}

func signBridgeRegistry(chain33Bridge, bridgebank, oracleAddr, valSetAddr common.Address, deployerAddr common.Address, nonce, gasLimit, gasPrice uint64, key *ecdsa.PrivateKey) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.BridgeRegistryABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeRegistryBin)
	input, err := parsed.Pack("", chain33Bridge, bridgebank, oracleAddr, valSetAddr)
	if err != nil {
		panic(err)
	}
	data := append(bin, input...)
	return signContractTx("bridgeRegistry", data, deployerAddr, nonce, gasLimit, gasPrice, key)
}
