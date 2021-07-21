package offline

import (
	"encoding/json"
	"fmt"
	"os"

	//"github.com/33cn/chain33/rpc/jsonclient"
	//rpctypes "github.com/33cn/chain33/rpc/types"
	erc20 "github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/erc20/generated"

	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/system/crypto/secp256k1"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"

	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4chain33/generated"
	"github.com/33cn/plugin/plugin/dapp/dex/utils"
	"github.com/spf13/cobra"
)

var crossXfileName = "deployCrossX2Chain33.txt"

func Boss4xOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "create and sign offline tx to deploy and set cross contracts to chain33",
	}
	cmd.AddCommand(
		createCrossBridgeCmd(),
		sendSignTxs2Chain33Cmd(),
		createERC20Cmd(),
	)
	return cmd
}

func createCrossBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create and sign all the offline cross to ethereum contracts(inclue valset,ethereumBridge,bridgeBank,oracle,bridgeRegistry,mulSign)",
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
	_ = args
	var txs []*utils.Chain33OfflineTx
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

	fmt.Printf("%d: Going to create MulSign2chain33 \n", i)
	i += 1
	createMulSign2chain33Tx, err := createMulSignAndSign(cmd, from)
	if nil != err {
		fmt.Println("Failed to createMulSign2chain33Tx due to cause:", err.Error())
		return
	}
	txs = append(txs, createMulSign2chain33Tx)

	fmt.Printf("%d: Write all the txs to file:   %s \n", i, crossXfileName)
	utils.WriteToFileInJson(crossXfileName, txs)
}

func getTxInfo(cmd *cobra.Command) *utils.TxCreateInfo {
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

	return info
}

func createBridgeRegistryTxAndSign(cmd *cobra.Command, from common.Address, ethereumBridge, valset, bridgeBank, oracle string) (*utils.Chain33OfflineTx, error) {
	createPara := fmt.Sprintf("%s,%s,%s,%s", ethereumBridge, bridgeBank, oracle, valset)
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.BridgeRegistryBin, generated.BridgeRegistryABI, createPara, "BridgeRegistry")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	bridgeRegistryTx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy BridgeRegistry",
		Interval:      time.Second * 5,
	}
	return bridgeRegistryTx, nil
}

func setOracle2EthBridgeTxAndSign(cmd *cobra.Command, ethbridge, oracle string) (*utils.Chain33OfflineTx, error) {
	//function setOracle(
	//	address _oracle
	//)
	parameter := fmt.Sprintf("setOracle(%s)", oracle)
	_, packData, err := evmAbi.Pack(parameter, generated.EthereumBridgeABI, false)
	if nil != err {
		fmt.Println("setOracle2EthBridge", "Failed to do abi.Pack due to:", err.Error())
		return nil, err
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: "setOracle2EthBridge", Para: packData, ContractAddr: ethbridge}
	content, txHash, err := utils.CallContractAndSign(getTxInfo(cmd), action, ethbridge)
	if nil != err {
		return nil, err
	}

	setOracle2EthBridgeTx := &utils.Chain33OfflineTx{
		ContractAddr:  ethbridge,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "setOracle2EthBridge",
		Interval:      time.Second * 5,
	}
	return setOracle2EthBridgeTx, nil
}

func setBridgeBank2EthBridgeTxAndSign(cmd *cobra.Command, ethbridge, bridgebank string) (*utils.Chain33OfflineTx, error) {
	//function setBridgeBank(
	//	address payable _bridgeBank
	//)
	parameter := fmt.Sprintf("setBridgeBank(%s)", bridgebank)
	_, packData, err := evmAbi.Pack(parameter, generated.EthereumBridgeABI, false)
	if nil != err {
		fmt.Println("setBridgeBank2EthBridge", "Failed to do abi.Pack due to:", err.Error())
		return nil, err
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: "setBridgeBank2EthBridge", Para: packData, ContractAddr: ethbridge}
	content, txHash, err := utils.CallContractAndSign(getTxInfo(cmd), action, ethbridge)
	if nil != err {
		return nil, err
	}

	setBridgeBank2EthBridgeTx := &utils.Chain33OfflineTx{
		ContractAddr:  ethbridge,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "setBridgeBank2EthBridge",
		Interval:      time.Second * 5,
	}
	return setBridgeBank2EthBridgeTx, nil
}

func createBridgeBankTxAndSign(cmd *cobra.Command, from common.Address, oracle, ethereumBridge string) (*utils.Chain33OfflineTx, error) {
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s,%s", operator, oracle, ethereumBridge)
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.BridgeBankBin, generated.BridgeBankABI, createPara, "bridgeBank")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	bridgeBankTx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy bridgeBank",
		Interval:      time.Second * 5,
	}
	return bridgeBankTx, nil
}

func createOracleTxAndSign(cmd *cobra.Command, from common.Address, valset, ethereumBridge string) (*utils.Chain33OfflineTx, error) {
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s,%s", operator, valset, ethereumBridge)
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.OracleBin, generated.OracleABI, createPara, "oralce")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	oracleTx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy oracle",
		Interval:      time.Second * 5,
	}
	return oracleTx, nil
}

func createValsetTxAndSign(cmd *cobra.Command, from common.Address) (*utils.Chain33OfflineTx, error) {
	contructParameter, _ := cmd.Flags().GetString("valset")
	createPara := contructParameter
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.ValsetBin, generated.ValsetABI, createPara, "valset")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	valsetTx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy valset",
		Interval:      time.Second * 5,
	}
	return valsetTx, nil
}

func createEthereumBridgeAndSign(cmd *cobra.Command, from common.Address, valset string) (*utils.Chain33OfflineTx, error) {
	operator := from.String()
	createPara := fmt.Sprintf("%s,%s", operator, valset)
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.EthereumBridgeBin, generated.EthereumBridgeABI, createPara, "EthereumBridge")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	ethereumBridgeTx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy ethereumBridge",
		Interval:      time.Second * 5,
	}
	return ethereumBridgeTx, nil
}

func createMulSignAndSign(cmd *cobra.Command, from common.Address) (*utils.Chain33OfflineTx, error) {
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), generated.GnosisSafeBin, generated.GnosisSafeABI, "", "mulSign2chain33")
	if nil != err {
		return nil, err
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	mulSign2chain33Tx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy mulSign2chain33",
		Interval:      time.Second * 5,
	}
	return mulSign2chain33Tx, nil
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
	cmd.Flags().StringP("file", "f", "", "signed tx file")
	_ = cmd.MarkFlagRequired("file")
}

func sendSignTxs2Chain33(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("file")
	url, _ := cmd.Flags().GetString("rpc_laddr")
	utils.SendSignTxs2Chain33(filePath, url)
}

func createERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_erc20",
		Short: "create erc20 contracts and sign, default 3300*1e8 to be minted",
		Run:   CreateERC20,
	}
	CreateERC20Flags(cmd)
	return cmd
}

//CreateERC20Flags ...
func CreateERC20Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")
	//cmd.Flags().StringP("owner", "o", "", "owner address")
	//_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringP("symbol", "s", "", "token symbol")
	_ = cmd.MarkFlagRequired("symbol")
	cmd.Flags().Float64P("amount", "a", 0, "amount to be minted(optional),default to 3300*1e8")
}

func CreateERC20(cmd *cobra.Command, args []string) {
	//rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	symbol, _ := cmd.Flags().GetString("symbol")
	//owner, _ := cmd.Flags().GetString("owner")
	amount, _ := cmd.Flags().GetFloat64("amount")
	amountInt64 := int64(3300 * 1e8)
	if 0 != int64(amount) {
		amountInt64 = int64(amount)
	}

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

	createPara := fmt.Sprintf("%s,%s,%s,%s", symbol, symbol, fmt.Sprintf("%d", amountInt64), fromAddr.String())
	content, txHash, err := utils.CreateContractAndSign(getTxInfo(cmd), erc20.ERC20Bin, erc20.ERC20ABI, createPara, "ERC20:"+symbol)
	if nil != err {
		fmt.Println("CreateContractAndSign erc20 fail")
		return
	}

	newContractAddr := common.NewContractAddress(from, txHash).String()
	Erc20Tx := &utils.Chain33OfflineTx{
		ContractAddr:  newContractAddr,
		TxHash:        common.Bytes2Hex(txHash),
		SignedRawTx:   content,
		OperationName: "deploy ERC20:" + symbol,
	}

	data, err := json.MarshalIndent(Erc20Tx, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(data))

	var txs []*utils.Chain33OfflineTx
	txs = append(txs, Erc20Tx)

	fileName := fmt.Sprintf("deployErc20%sChain33.txt", symbol)
	fmt.Printf("Write all the txs to file:   %s \n", fileName)
	utils.WriteToFileInJson(fileName, txs)
}
