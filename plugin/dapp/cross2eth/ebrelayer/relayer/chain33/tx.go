package chain33

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/33cn/chain33/common"
	chain33Common "github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	chain33Crypto "github.com/33cn/chain33/common/crypto"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/system/crypto/secp256k1"
	types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/contracts/contracts4chain33/generated"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
)

//DeployPara ...
type DeployPara4Chain33 struct {
	Deployer       address.Address
	Operator       address.Address
	InitValidators []address.Address
	InitPowers     []*big.Int
}

type DeployResult struct {
	Address address.Address
	TxHash  string
}

type X2EthDeployInfo struct {
	BridgeRegistry *DeployResult
	BridgeBank     *DeployResult
	Chain33Bridge  *DeployResult
	Valset         *DeployResult
	Oracle         *DeployResult
}

var chain33txLog = log.New("module", "chain33_txs")

// RelayLockToChain33 : RelayLockToChain33 applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockBurnToChain33(privateKey chain33Crypto.PrivKey, privateKey_ecdsa *ecdsa.PrivateKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL, oracleAddr string) (string, error) {
	nonceBytes := big.NewInt(claim.Nonce).Bytes()
	amountBytes := big.NewInt(claim.Amount).Bytes()
	claimID := crypto.Keccak256Hash(nonceBytes, []byte(claim.EthereumSender), []byte(claim.Chain33Receiver), []byte(claim.Symbol), amountBytes)

	// Sign the hash using the active validator's private key
	signature, err := utils.SignClaim4Evm(claimID, privateKey_ecdsa)
	if nil != err {
		return "", err
	}
	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %s, %s, %s)",
		claim.ClaimType,
		claim.EthereumSender,
		claim.Chain33Receiver,
		claim.TokenAddr,
		claim.Symbol,
		claim.Amount,
		claimID,
		signature)

	return relayEvmTx2Chain33(privateKey, claim, parameter, rpcURL, oracleAddr)
}

func createEvmTx(privateKey chain33Crypto.PrivKey, action proto.Message, execer, to string, fee int64) string {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: fee, To: to}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	tx.ChainID = 33

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)
	return dataStr
}

func relayEvmTx2Chain33(privateKey chain33Crypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, parameter, rpcURL, oracleAddr string) (string, error) {
	note := fmt.Sprintf("RelayLockToChain33 by validator:%s with nonce:%d",
		address.PubKeyToAddr(privateKey.PubKey().Bytes()),
		claim.Nonce)

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Abi: parameter}

	feeInt64 := int64(1e7)
	toAddr := oracleAddr
	wholeEvm := claim.ChainName + ".evm"
	//name表示发给哪个执行器
	data := createEvmTx(privateKey, &action, wholeEvm, toAddr, feeInt64)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  data,
	}
	var txhash string

	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain33.SendTransaction", params, &txhash)
	_, err := ctx.RunResult()
	return txhash, err
}

func queryTxsByHashesRes(arg interface{}) (interface{}, error) {
	var receipt *rpctypes.ReceiptDataResult
	for _, v := range arg.(*rpctypes.TransactionDetails).Txs {
		if v == nil {
			continue
		}
		receipt = v.Receipt
		if nil != receipt {
			return receipt.Ty, nil
		}
	}
	return nil, nil
}

func getTxByHashesRpc(txhex, rpcLaddr string) (string, error) {
	hashesArr := strings.Split(txhex, " ")
	params2 := rpctypes.ReqHashes{
		Hashes: hashesArr,
	}

	var res rpctypes.TransactionDetails
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.GetTxByHashes", params2, &res)
	ctx.SetResultCb(queryTxsByHashesRes)
	result, err := ctx.RunResult()
	if err != nil || result == nil {
		return "", err
	}
	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getContractAddr(caller, txhex string) string {
	return address.GetExecAddress(caller + ethcommon.Bytes2Hex(common.HexToHash(txhex).Bytes())).String()
}

func DeployAndInit2Chain33_(rpcLaddr, paraChainName string, para4deploy *DeployPara4Chain33) (*X2EthDeployInfo, error) {
	deployer := para4deploy.Deployer.String()
	deployInfo := &X2EthDeployInfo{}
	var err error
	constructorPara := ""
	paraLen := len(para4deploy.InitValidators)

	var valsetAddr string
	var ethereumBridgeAddr string
	var oracleAddr string
	var bridgeBankAddr string

	//x2EthContracts.Valset, deployInfo.Valset, err = DeployValset(client, para.DeployPrivateKey, para.Deployer, para.Operator, para.InitValidators, para.InitPowers)
	//constructor(
	//	address _operator,
	//	address[] memory _initValidators,
	//	uint256[] memory _initPowers
	//)
	if 1 == paraLen {
		constructorPara = fmt.Sprintf("constructor(%s, %s, %d)", para4deploy.Operator.String(),
			para4deploy.InitValidators[0].String(),
			para4deploy.InitPowers[0].Int64())
	} else if 4 == paraLen {
		constructorPara = fmt.Sprintf("constructor(%s, %s, %s, %s, %s, %d, %d, %d, %d)", para4deploy.Operator.String(),
			para4deploy.InitValidators[0].String(), para4deploy.InitValidators[1].String(), para4deploy.InitValidators[2].String(), para4deploy.InitValidators[3].String(),
			para4deploy.InitPowers[0].Int64(), para4deploy.InitPowers[1].Int64(), para4deploy.InitPowers[2].Int64(), para4deploy.InitPowers[3].Int64())
	} else {
		panic(fmt.Sprintf("Not support valset with parameter count=%d", paraLen))
	}

	deployValsetHash, err := deploySingleContract(ethcommon.FromHex(generated.ValsetBin), generated.ValsetABI, constructorPara, "valset", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployValset due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("\nDeployValset tx hash:", deployInfo.Valset.TxHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployValset timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployValsetHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy valset tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy valset failed due to" + ", ty = " + data)
				}
				valsetAddr = getContractAddr(deployer, deployValsetHash)
				fmt.Println("Succeed to deploy valset with address =", valsetAddr, "\n")
				goto deployEthereumBridge
			}
		}
	}

deployEthereumBridge:
	//x2EthContracts.Chain33Bridge, deployInfo.Chain33Bridge, err = DeployChain33Bridge(client, para.DeployPrivateKey, para.Deployer, para.Operator, deployInfo.Valset.Address)
	//constructor(
	//	address _operator,
	//	address _valset
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s)", para4deploy.Operator.String(), valsetAddr)
	deployEthereumBridgeHash, err := deploySingleContract(ethcommon.FromHex(generated.EthereumBridgeBin), generated.EthereumBridgeABI, constructorPara, "EthereumBridge", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to deployEthereumBridge due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("\nDeploy EthereumBridge Hash tx hash:", deployEthereumBridgeHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployEthereumBridge timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployEthereumBridgeHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy EthereumBridge tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy EthereumBridge failed due to" + ", ty = " + data)
				}
				ethereumBridgeAddr = getContractAddr(deployer, deployEthereumBridgeHash)
				fmt.Println("Succeed to deploy EthereumBridge with address =", ethereumBridgeAddr, "\n")
				goto deployOracle
			}
		}
	}

deployOracle:
	//constructor(
	//	address _operator,
	//	address _valset,
	//	address _ethereumBridge
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s)", para4deploy.Operator.String(), valsetAddr, ethereumBridgeAddr)
	//x2EthContracts.Oracle, deployInfo.Oracle, err = DeployOracle(client, para.DeployPrivateKey, para.Deployer, para.Operator, deployInfo.Valset.Address, deployInfo.Chain33Bridge.Address)
	deployOracleHash, err := deploySingleContract(ethcommon.FromHex(generated.OracleBin), generated.OracleABI, constructorPara, "Oracle", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployOracle due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("DeployOracle tx hash:", deployOracleHash)

		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployOracle timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployOracleHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy Oracle tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy Oracle failed due to" + ", ty = " + data)
				}
				oracleAddr = getContractAddr(deployer, deployOracleHash)
				fmt.Println("Succeed to deploy EthereumBridge with address =", oracleAddr, "\n")
				goto deployBridgeBank
			}
		}
	}
	/////////////////////////////////////
deployBridgeBank:
	//constructor (
	//	address _operatorAddress,
	//	address _oracleAddress,
	//	address _ethereumBridgeAddress
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s)", para4deploy.Operator.String(), oracleAddr, ethereumBridgeAddr)
	deployBridgeBankHash, err := deploySingleContract(ethcommon.FromHex(generated.BridgeBankBin), generated.BridgeBankABI, constructorPara, "BridgeBank", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployBridgeBank due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("deployBridgeBank tx hash:", deployBridgeBankHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployBridgeBank timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy BridgeBank tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy BridgeBank failed due to" + ", ty = " + data)
				}
				bridgeBankAddr = getContractAddr(deployer, deployBridgeBankHash)
				fmt.Println("Succeed to deploy BridgeBank with address =", bridgeBankAddr, "\n")
				goto settingBridgeBank
			}
		}
	}

settingBridgeBank:
	////////////////////////
	//function setBridgeBank(
	//	address payable _bridgeBank
	//)
	callPara := fmt.Sprintf("setBridgeBank(%s)", bridgeBankAddr)
	settingBridgeBankHash, err := sendTx2Evm(callPara, rpcLaddr, ethereumBridgeAddr, paraChainName, deployer)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to settingBridgeBank due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("setBridgeBank tx hash:", settingBridgeBankHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("setBridgeBank timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for setBridgeBank tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("setBridgeBank failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to setBridgeBank ")
				goto setOracle
			}
		}
	}

setOracle:
	//function setOracle(
	//	address _oracle
	//)
	callPara = fmt.Sprintf("setOracle(%s)", oracleAddr)
	setOracleHash, err := sendTx2Evm(callPara, rpcLaddr, ethereumBridgeAddr, paraChainName, deployer)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to setOracle due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("setOracle tx hash:", setOracleHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("setOracle timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for setOracle tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("setOracle failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to setOracle ")
				goto deployBridgeRegistry
			}
		}
	}

deployBridgeRegistry:
	//constructor(
	//	address _ethereumBridge,
	//	address _bridgeBank,
	//	address _oracle,
	//	address _valset
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s, %s)", ethereumBridgeAddr, bridgeBankAddr, oracleAddr, valsetAddr)
	deployBridgeRegistryHash, err := deploySingleContract(ethcommon.FromHex(generated.BridgeRegistryBin), generated.BridgeRegistryABI, constructorPara, "BridgeRegistry", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to deployBridgeRegistry due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("deployBridgeRegistryHash tx hash:", deployBridgeRegistryHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployBridgeRegistry timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for deployBridgeRegistry tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("deployBridgeRegistry failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deployBridgeRegistry")
				goto finished
			}
		}
	}
finished:

	return deployInfo, nil
}

func deploySingleContract(code []byte, abi, constructorPara, contractName, paraChainName, deployer, rpcLaddr string) (string, error) {
	note := "deploy " + contractName

	var action evmtypes.EVMContractAction
	action = evmtypes.EVMContractAction{Amount: 0, Code: code, GasLimit: 0, GasPrice: 0, Note: note, Alias: contractName, Abi: string(abi)}
	if constructorPara != "" {
		packData, err := evmAbi.PackContructorPara(constructorPara, abi)
		if err != nil {
			return "", errors.New(contractName + " " + constructorPara + " Pack Contructor Para error:" + err.Error())
		}
		action.Code = append(action.Code, packData...)
	}
	data, err := createSignedEvmTx(&action, paraChainName+"evm", deployer, rpcLaddr)
	if err != nil {
		return "", errors.New(contractName + " create contract error:" + err.Error())
	}

	txhex, err := sendTransactionRpc(data, rpcLaddr)
	if err != nil {
		return "", errors.New(contractName + " send transaction error:" + err.Error())
	}
	fmt.Println("Deploy", contractName, "tx hash:", txhex)

	return txhex, nil
}

func createSignedEvmTx(action proto.Message, execer, caller, rpcLaddr string) (string, error) {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: int64(1e8)}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	//tx.ChainID = cfg.GetChainID()
	txHex := types.Encode(tx)
	rawTx := hex.EncodeToString(txHex)

	unsignedTx := &types.ReqSignRawTx{
		Addr:  caller,
		TxHex: rawTx,
		Fee:   tx.Fee,
	}

	var res string
	client, err := jsonclient.NewJSONClient(rpcLaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}
	err = client.Call("Chain33.SignRawTx", unsignedTx, &res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}

	return res, nil
}

func sendTransactionRpc(data, rpcLaddr string) (string, error) {
	params := rpctypes.RawParm{
		Data: data,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.SendTransaction", params, nil)
	var txhex string
	rpc, err := jsonclient.NewJSONClient(ctx.Addr)
	if err != nil {
		return "", err
	}

	err = rpc.Call(ctx.Method, ctx.Params, &txhex)
	if err != nil {
		return "", err
	}

	return txhex, nil
}

func sendTx2Evm(parameter, rpcURL, toAddr, chainName, caller string) (string, error) {
	note := fmt.Sprintf("sendTx2Evm by caller:%s", caller)

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Abi: parameter}
	wholeEvm := chainName + "evm"
	data, err := createSignedEvmTx(&action, wholeEvm, caller, rpcURL)
	if err != nil {
		return "", errors.New(toAddr + " createSignedEvmTx error:" + err.Error())
	}

	txhex, err := sendTransactionRpc(data, rpcURL)
	if err != nil {
		return "", errors.New(toAddr + " send transaction error:" + err.Error())
	}
	return txhex, nil
}

func approve(privateKey chain33Crypto.PrivKey, contractAddr, spender, chainName, rpcURL string, amount int64) (string, error) {
	note := fmt.Sprintf("approve for spender:%s ", spender)

	//approve(address spender, uint256 amount)
	parameter := fmt.Sprint("approve(%s, %d)", spender, amount)
	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Abi: parameter}

	feeInt64 := int64(1e7)
	toAddr := contractAddr
	wholeEvm := chainName + "evm"
	//name表示发给哪个执行器
	data := createEvmTx(privateKey, &action, wholeEvm, toAddr, feeInt64)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  data,
	}
	var txhash string

	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain33.SendTransaction", params, &txhash)
	_, err := ctx.RunResult()
	return txhash, err

}

func BurnAsync(ownerPrivateKeyStr, tokenAddrstr, chain33Receiver string, amount int64, bridgeBankAddr string, chainName, rpcURL string) (string, error) {
	var driver secp256k1.Driver
	privateKeySli, err := chain33Common.FromHex(ownerPrivateKeyStr)
	if nil != err {
		return "", err
	}
	ownerPrivateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return "", err
	}

	approveTxHash, err := approve(ownerPrivateKey, tokenAddrstr, bridgeBankAddr, chainName, rpcURL, amount)
	if err != nil {
		chain33txLog.Error("BurnAsync", "failed to send approve tx due to:", err.Error())
		return "", err
	}

	return "", err
}
