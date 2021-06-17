package offline

import (
	"fmt"
	"time"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/system/crypto/secp256k1"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"

	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4chain33/generated"
	"github.com/33cn/plugin/plugin/dapp/dex/utils"
	"github.com/spf13/cobra"
)

type Chain33OfflineTx struct {
	ContractAddr  string
	TxHash        string
	SignedRawTx   string
	OperationName string
	Interval      time.Duration
}

var crossXfileName = "deployCrossX2Chain33.txt"

func OfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "create and sign offline tx to deploy and set cross contracts to chain33",
	}
	cmd.AddCommand(
		createCrossBridgeCmd(),
		sendSignTxs2Chain33Cmd(),
	)
	return cmd
}

func createCrossBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create and sign all the offline cross to ethereum contracts(inclue valset,ethereumBridge, bridgeBank,oracle,bridgeRegistry)",
		Run:   createCrossBridge,
	}
	addCreateCrossBridgeFlags(cmd)
	return cmd
}

func addCreateCrossBridgeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("valset", "r", "", "contruct parameter for valset, as: 'addr, [addr, addr, addr, addr], [25, 25, 25, 25]'")
	_ = cmd.MarkFlagRequired("valset")
}

func createCrossBridge(cmd *cobra.Command, args []string) {
	var txs []*Chain33OfflineTx
	privateKeyStr, _ := cmd.Flags().GetString("key")
	var driver secp256k1.Driver
	privateKeySli := common.FromHex(privateKeyStr)
	privateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		fmt.Println("Failed to do PrivKeyFromBytes")
		return
	}
	fromAddr := address.PubKeyToAddress(privateKey.PubKey().Bytes())
	from := common.Address{
		Addr: fromAddr,
	}
	i := 1
	fmt.Printf("%d: Going to create Valset\n", i)
	i += 1
	valsetTx, err := createValsetTxAndSign(cmd, from)
	if nil != err {
		fmt.Println("Failed to createValsetTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, valsetTx)

	fmt.Printf("%d: Going to create EthereumBridge\n", i)
	i += 1
	ethereumBridgeTx, err := createEthereumBridgeAndSign(cmd, from, valsetTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to createEthereumBridgeAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, ethereumBridgeTx)

	fmt.Printf("%d: Going to create Oracle\n", i)
	i += 1
	oracleTx, err := createOracleTxAndSign(cmd, from, valsetTx.ContractAddr, ethereumBridgeTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to createOracleTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, oracleTx)

	fmt.Printf("%d: Going to create BridgeBank\n", i)
	i += 1
	bridgeBankTx, err := createBridgeBankTxAndSign(cmd, from, valsetTx.ContractAddr, ethereumBridgeTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to createBridgeBankTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, bridgeBankTx)

	fmt.Printf("%d: Going to set BridgeBank to EthBridge \n", i)
	i += 1
	setBridgeBank2EthBridgeTx, err := setBridgeBank2EthBridgeTxAndSign(cmd, ethereumBridgeTx.ContractAddr, bridgeBankTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to setBridgeBank2EthBridgeTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, setBridgeBank2EthBridgeTx)

	fmt.Printf("%d: Going to set Oracle to EthBridge \n", i)
	i += 1
	setOracle2EthBridgeTx, err := setOracle2EthBridgeTxAndSign(cmd, ethereumBridgeTx.ContractAddr, oracleTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to setOracle2EthBridgeTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, setOracle2EthBridgeTx)

	fmt.Printf("%d: Going to create BridgeRegistry \n", i)
	i += 1
	createBridgeRegistryTx, err := createBridgeRegistryTxAndSign(cmd, from, ethereumBridgeTx.ContractAddr, valsetTx.ContractAddr, bridgeBankTx.ContractAddr, oracleTx.ContractAddr)
	if nil != err {
		fmt.Println("Failed to createBridgeRegistryTxAndSign due to cause:", err.Error())
		return
	}
	txs = append(txs, createBridgeRegistryTx)

	fmt.Printf("%d: Write all the txs to file:   %s \n", i, crossXfileName)
	utils.WriteToFileInJson(crossXfileName, txs)
}

func createBridgeRegistryTxAndSign(cmd *cobra.Command, from common.Address, ethereumBridge, valset, bridgeBank, oracle string) (*Chain33OfflineTx, error) {
	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//constructor(
	//	address _ethereumBridge,
	//	address _bridgeBank,
	//	address _oracle,
	//	address _valset
	//)
	//constructor(addr, [addr, addr, addr, addr], [25, 25, 25, 25])
	createPara := fmt.Sprintf("%s,%s,%s,%s", ethereumBridge, bridgeBank, oracle, valset)
	content, txHash, err := utils.CreateContractAndSign(info, generated.BridgeRegistryBin, generated.BridgeRegistryABI, createPara, "BridgeRegistry")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	bridgeRegistryTx := &Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy BridgeRegistry",
		Interval:      time.Second * 5,
	}
	return bridgeRegistryTx, nil
}

func setOracle2EthBridgeTxAndSign(cmd *cobra.Command, ethbridge, oracle string) (*Chain33OfflineTx, error) {
	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//function setOracle(
	//	address _oracle
	//)
	parameter := fmt.Sprintf("setOracle(%s)", oracle)
	_, packData, err := evmAbi.Pack(parameter, generated.EthereumBridgeABI, false)
	if nil != err {
		fmt.Println("setOracle2EthBridge", "Failed to do abi.Pack due to:", err.Error())
		return nil, err
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, txHash, err := utils.CallContractAndSign(info, action, ethbridge)
	if nil != err {
		return nil, err
	}

	setOracle2EthBridgeTx := &Chain33OfflineTx{
		ContractAddr:  ethbridge,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "setOracle2EthBridge",
		Interval:      time.Second * 5,
	}
	return setOracle2EthBridgeTx, nil
}

func setBridgeBank2EthBridgeTxAndSign(cmd *cobra.Command, ethbridge, bridgebank string) (*Chain33OfflineTx, error) {

	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//function setBridgeBank(
	//	address payable _bridgeBank
	//)
	parameter := fmt.Sprintf("setBridgeBank(%s)", bridgebank)
	_, packData, err := evmAbi.Pack(parameter, generated.EthereumBridgeABI, false)
	if nil != err {
		fmt.Println("setBridgeBank2EthBridge", "Failed to do abi.Pack due to:", err.Error())
		return nil, err
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, txHash, err := utils.CallContractAndSign(info, action, ethbridge)
	if nil != err {
		return nil, err
	}

	setBridgeBank2EthBridgeTx := &Chain33OfflineTx{
		ContractAddr:  ethbridge,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "setBridgeBank2EthBridge",
		Interval:      time.Second * 5,
	}
	return setBridgeBank2EthBridgeTx, nil
}

func createBridgeBankTxAndSign(cmd *cobra.Command, from common.Address, oracle, ethereumBridge string) (*Chain33OfflineTx, error) {
	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//constructor (
	//	address _operatorAddress,
	//	address _oracleAddress,
	//	address _ethereumBridgeAddress
	//)
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s,%s", operator, oracle, ethereumBridge)
	content, txHash, err := utils.CreateContractAndSign(info, generated.BridgeBankBin, generated.BridgeBankABI, createPara, "bridgeBank")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	bridgeBankTx := &Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy bridgeBank",
		Interval:      time.Second * 5,
	}
	return bridgeBankTx, nil
}

func createOracleTxAndSign(cmd *cobra.Command, from common.Address, valset, ethereumBridge string) (*Chain33OfflineTx, error) {
	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//constructor(
	//	address _operator,
	//	address _valset,
	//	address _ethereumBridge
	//)
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s,%s", operator, valset, ethereumBridge)
	content, txHash, err := utils.CreateContractAndSign(info, generated.OracleBin, generated.OracleABI, createPara, "oralce")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	oracleTx := &Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy oracle",
		Interval:      time.Second * 5,
	}
	return oracleTx, nil
}

func createValsetTxAndSign(cmd *cobra.Command, from common.Address) (*Chain33OfflineTx, error) {
	contructParameter, _ := cmd.Flags().GetString("valset")

	privateKeyStr, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKeyStr,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//constructor(
	//	address _operator,
	//	address[] memory _initValidators,
	//	uint256[] memory _initPowers
	//)
	//constructor(addr, [addr, addr, addr, addr], [25, 25, 25, 25])
	createPara := contructParameter
	content, txHash, err := utils.CreateContractAndSign(info, generated.ValsetBin, generated.ValsetABI, createPara, "valset")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	valsetTx := &Chain33OfflineTx{
		ContractAddr: newContractAddr,

		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy valset",
		Interval:      time.Second * 5,
	}
	return valsetTx, nil
}

func createEthereumBridgeAndSign(cmd *cobra.Command, from common.Address, valset string) (*Chain33OfflineTx, error) {
	privateKey, _ := cmd.Flags().GetString("key")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &utils.TxCreateInfo{
		PrivateKey: privateKey,
		Expire:     expire,
		Note:       note,
		Fee:        feeInt64,
		ParaName:   paraName,
		ChainID:    chainID,
	}
	//constructor(
	//	address _operator,
	//	address _valset
	//)
	//constructor(addr, [addr, addr, addr, addr], [25, 25, 25, 25])
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s", operator, valset)
	content, txHash, err := utils.CreateContractAndSign(info, generated.EthereumBridgeBin, generated.EthereumBridgeABI, createPara, "EthereumBridge")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	ethereumBridgeTx := &Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy ethereumBridge",
		Interval:      time.Second * 5,
	}
	return ethereumBridgeTx, nil
}

func sendSignTxs2Chain33Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "send all the crossX txs to chain33 in serial",
		Run:   sendSignTxs2Chain33,
	}
	addSendSignTxs2Chain33Flags(cmd)
	return cmd
}

func addSendSignTxs2Chain33Flags(cmd *cobra.Command) {
	usage := fmt.Sprintf("(optional)path of txs file,default to current directroy ./, and file name is:%s", crossXfileName)
	cmd.Flags().StringP("path", "p", "./", usage)
}

func sendSignTxs2Chain33(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("path")
	url, _ := cmd.Flags().GetString("rpc_laddr")
	//解析文件数据
	filePath += crossXfileName
	var rdata []*Chain33OfflineTx
	err := utils.ParseFileInJson(filePath, &rdata)
	if err != nil {
		fmt.Printf("parse file with error:%s, make sure file:%s exist", err.Error(), filePath)
		return
	}
	for i, deployInfo := range rdata {
		txhash, err := sendTransactionRpc(deployInfo.SignedRawTx, url)
		if nil != err {
			fmt.Printf("Failed to send %s to chain33 due to error:%s", deployInfo.OperationName, err.Error())
			return
		}
		fmt.Printf("   %d:Succeed to send %s to chain33 with tx hash:%s\n", i+1, deployInfo.OperationName, txhash)
		if deployInfo.Interval != 0 {
			time.Sleep(deployInfo.Interval)
		}
	}

	fmt.Println("All txs are sent successfully.^-^ ^-^")
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
