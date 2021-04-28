package chain33

import (
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
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
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

type X2EthDeployResult struct {
	BridgeRegistry *DeployResult
	BridgeBank     *DeployResult
	EthereumBridge *DeployResult
	Valset         *DeployResult
	Oracle         *DeployResult
}

var chain33txLog = log.New("module", "chain33_txs")

// RelayLockToChain33 : RelayLockToChain33 applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
//func RelayLockBurnToChain33(privateKey chain33Crypto.PrivKey, privateKey_ecdsa *ecdsa.PrivateKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL, oracleAddr string) (string, error) {
//	nonceBytes := big.NewInt(claim.Nonce).Bytes()
//	amountBytes := big.NewInt(claim.Amount).Bytes()
//	claimID := crypto.Keccak256Hash(nonceBytes, []byte(claim.EthereumSender), []byte(claim.Chain33Receiver), []byte(claim.Symbol), amountBytes)
//
//	// Sign the hash using the active validator's private key
//	signature, err := utils.SignClaim4Evm(claimID, privateKey_ecdsa)
//	if nil != err {
//		return "", err
//	}
//	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %s, %s, %s)",
//		claim.ClaimType,
//		claim.EthereumSender,
//		claim.Chain33Receiver,
//		claim.TokenAddr,
//		claim.Symbol,
//		claim.Amount,
//		claimID,
//		signature)
//
//	return relayEvmTx2Chain33(privateKey, claim, parameter, rpcURL, oracleAddr)
//}

func createEvmTx(privateKey chain33Crypto.PrivKey, action proto.Message, execer, to string, fee int64) string {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: fee, To: to}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	tx.ChainID = 0

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)
	return dataStr
}

func relayEvmTx2Chain33(privateKey chain33Crypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, parameter, rpcURL, oracleAddr string) (string, error) {
	note := fmt.Sprintf("relayEvmTx2Chain33 by validator:%s with nonce:%d",
		address.PubKeyToAddr(privateKey.PubKey().Bytes()),
		claim.Nonce)
	_, packData, err := evmAbi.Pack(parameter, generated.OracleABI, false)
	if nil != err {
		chain33txLog.Info("relayEvmTx2Chain33", "Failed to do abi.Pack due to:", err.Error())
		return "", ebrelayerTypes.ErrPack
	}

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Para: packData}

	//TODO: 交易费超大问题需要调查，hezhengjun on 20210420
	feeInt64 := int64(5 * 1e7)
	toAddr := oracleAddr

	wholeEvm := getExecerName(claim.ChainName)
	//name表示发给哪个执行器
	data := createEvmTx(privateKey, &action, wholeEvm, toAddr, feeInt64)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  data,
	}
	var txhash string

	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain33.SendTransaction", params, &txhash)
	_, err = ctx.RunResult()
	return txhash, err
}

func getExecerName(name string) string {
	var ret string
	names := strings.Split(name, ".")
	for _, v := range names {
		if v != "" {
			ret = ret + v + "."
		}
	}
	ret += "evm"
	return ret
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

func getContractAddr(caller, txhex string) address.Address {
	return *address.GetExecAddress(caller + ethcommon.Bytes2Hex(common.HexToHash(txhex).Bytes()))
}

func deploySingleContract(code []byte, abi, constructorPara, contractName, paraChainName, deployer, rpcLaddr string) (string, error) {
	note := "deploy " + contractName

	var action evmtypes.EVMContractAction
	action = evmtypes.EVMContractAction{Amount: 0, Code: code, GasLimit: 0, GasPrice: 0, Note: note, Alias: contractName}
	if constructorPara != "" {
		packData, err := evmAbi.PackContructorPara(constructorPara, abi)
		if err != nil {
			return "", errors.New(contractName + " " + constructorPara + " Pack Contructor Para error:" + err.Error())
		}
		action.Code = append(action.Code, packData...)
	}

	exector := paraChainName + "evm"
	to := address.ExecAddress(exector)
	data, err := createSignedEvmTx(&action, exector, deployer, rpcLaddr, to)
	if err != nil {
		return "", errors.New(contractName + " create contract error:" + err.Error())
	}

	txhex, err := sendTransactionRpc(data, rpcLaddr)
	if err != nil {
		return "", errors.New(contractName + " send transaction error:" + err.Error())
	}
	chain33txLog.Info("deploySingleContract", "Deploy contract for", contractName, " with tx hash:", txhex)
	return txhex, nil
}

func createSignedEvmTx(action proto.Message, execer, caller, rpcLaddr, to string) (string, error) {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: int64(1e8), To: to}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	//tx.ChainID = cfg.GetChainID()
	txHex := types.Encode(tx)
	rawTx := hex.EncodeToString(txHex)

	unsignedTx := &types.ReqSignRawTx{
		Addr:   caller,
		TxHex:  rawTx,
		Fee:    tx.Fee,
		Expire: "120s",
	}

	var res string
	client, err := jsonclient.NewJSONClient(rpcLaddr)
	if err != nil {
		chain33txLog.Error("createSignedEvmTx", "jsonclient.NewJSONClient", err.Error())
		return "", err
	}
	err = client.Call("Chain33.SignRawTx", unsignedTx, &res)
	if err != nil {
		chain33txLog.Error("createSignedEvmTx", "Chain33.SignRawTx", err.Error())
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

func sendTx2Evm(parameter []byte, rpcURL, toAddr, chainName, caller string) (string, error) {
	note := fmt.Sprintf("sendTx2Evm by caller:%s", caller)

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Para: parameter}
	wholeEvm := chainName + "evm"
	data, err := createSignedEvmTx(&action, wholeEvm, caller, rpcURL, toAddr)
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
	parameter := fmt.Sprintf("approve(%s, %d)", spender, amount)
	_, packData, err := evmAbi.Pack(parameter, generated.BridgeTokenABI, false)
	if nil != err {
		chain33txLog.Info("approve", "Failed to do abi.Pack due to:", err.Error())
		return "", err
	}
	return sendEvmTx(privateKey, contractAddr, chainName, rpcURL, note, packData)
}

func burn(privateKey chain33Crypto.PrivKey, contractAddr, ethereumReceiver, ethereumTokenAddress, chainName, rpcURL string, amount int64) (string, error) {
	//    function burnBridgeTokens(
	//        bytes memory _ethereumReceiver,
	//        address _ethereumTokenAddress,
	//        uint256 _amount
	//    )
	parameter := fmt.Sprintf("burnBridgeTokens(%s, %s, %d)", ethereumReceiver, ethereumTokenAddress, amount)
	note := parameter
	_, packData, err := evmAbi.Pack(parameter, generated.BridgeBankABI, false)
	if nil != err {
		chain33txLog.Info("burn", "Failed to do abi.Pack due to:", err.Error())
		return "", err
	}

	return sendEvmTx(privateKey, contractAddr, chainName, rpcURL, note, packData)
}

func lockBty(privateKey chain33Crypto.PrivKey, contractAddr, ethereumReceiver, chainName, rpcURL string, amount int64) (string, error) {
	//function lock(
	//	bytes memory _recipient,
	//	address _token,
	//	uint256 _amount
	//)
	parameter := fmt.Sprintf("lock(%s, %s, %d)", ethereumReceiver, "1111111111111111111114oLvT2", amount)
	note := parameter
	_, packData, err := evmAbi.Pack(parameter, generated.BridgeBankABI, false)
	if nil != err {
		chain33txLog.Info("setOracle", "Failed to do abi.Pack due to:", err.Error())
		return "", ebrelayerTypes.ErrPack
	}
	return sendEvmTx(privateKey, contractAddr, chainName, rpcURL, note, packData)
}

func sendEvmTx(privateKey chain33Crypto.PrivKey, contractAddr, chainName, rpcURL, note string, parameter []byte) (string, error) {
	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Para: parameter}

	feeInt64 := int64(1e7)
	toAddr := contractAddr
	wholeEvm := getExecerName(chainName)
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

func burnAsync(ownerPrivateKeyStr, tokenAddrstr, ethereumReceiver string, amount int64, bridgeBankAddr string, chainName, rpcURL string) (string, error) {
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
	chain33txLog.Debug("BurnAsync", "approve with tx hash", approveTxHash)

	//privateKey chain33Crypto.PrivKey, contractAddr, ethereumReceiver, ethereumTokenAddress, chainName, rpcURL string, amount int6
	burnTxHash, err := burn(ownerPrivateKey, bridgeBankAddr, ethereumReceiver, tokenAddrstr, chainName, rpcURL, amount)
	if err != nil {
		chain33txLog.Error("BurnAsync", "failed to send burn tx due to:", err.Error())
		return "", err
	}
	chain33txLog.Debug("BurnAsync", "burn with tx hash", burnTxHash)

	return "", err
}

func lockAsync(ownerPrivateKeyStr, ethereumReceiver string, amount int64, bridgeBankAddr string, chainName, rpcURL string) (string, error) {
	var driver secp256k1.Driver
	privateKeySli, err := chain33Common.FromHex(ownerPrivateKeyStr)
	if nil != err {
		return "", err
	}
	ownerPrivateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return "", err
	}

	//privateKey chain33Crypto.PrivKey, contractAddr, ethereumReceiver, ethereumTokenAddress, chainName, rpcURL string, amount int6
	lockBtyTxHash, err := lockBty(ownerPrivateKey, bridgeBankAddr, ethereumReceiver, chainName, rpcURL, amount)
	if err != nil {
		chain33txLog.Error("lockBty", "failed to send approve tx due to:", err.Error())
		return "", err
	}
	chain33txLog.Debug("lockBty", "lockBty with tx hash", lockBtyTxHash)

	return "", err
}

func recoverContractAddrFromRegistry(bridgeRegistry, rpcLaddr string) (oracle, bridgeBank string) {
	parameter := fmt.Sprint("oracle()")

	result := query(bridgeRegistry, parameter, bridgeRegistry, rpcLaddr, generated.BridgeRegistryABI)
	if nil == result {
		return "", ""
	}
	oracle = result.(string)

	parameter = fmt.Sprint("bridgeBank()")
	result = query(bridgeRegistry, parameter, bridgeRegistry, rpcLaddr, generated.BridgeRegistryABI)
	if nil == result {
		return "", ""
	}
	bridgeBank = result.(string)
	return
}

func getToken2address(bridgeBank, symbol, rpcLaddr string) string {
	parameter := fmt.Sprintf("getToken2address(%s)", symbol)

	result := query(bridgeBank, parameter, bridgeBank, rpcLaddr, generated.BridgeBankABI)
	if nil == result {
		return ""
	}
	return result.(string)
}

func query(contractAddr, input, caller, rpcLaddr, abiData string) interface{} {
	methodName, packedinput, err := evmAbi.Pack(input, abiData, true)
	if err != nil {
		chain33txLog.Debug("query", "Failed to do para pack due to", err.Error())
		return nil
	}

	var req = evmtypes.EvmQueryReq{Address: contractAddr, Input: common.ToHex(packedinput), Caller: caller}
	var resp evmtypes.EvmQueryResp
	query := sendQuery(rpcLaddr, "Query", &req, &resp)

	if !query {
		return nil
	}
	_, err = json.MarshalIndent(&resp, "", "  ")
	if err != nil {
		fmt.Println(resp.String())
		return nil
	}

	data, err := common.FromHex(resp.RawData)
	if nil != err {
		fmt.Println("common.FromHex failed due to:", err.Error())
	}

	outputs, err := evmAbi.Unpack(data, methodName, abiData)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "unpack evm return error", err)
	}
	chain33txLog.Debug("query", "outputs", outputs)

	return outputs[0].Value
}

func sendQuery(rpcAddr, funcName string, request types.Message, result proto.Message) bool {
	params := rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: funcName,
		Payload:  types.MustPBToJSON(request),
	}

	jsonrpc, err := jsonclient.NewJSONClient(rpcAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}

	err = jsonrpc.Call("Chain33.Query", params, result)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false
	}
	return true
}
