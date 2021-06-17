package offline

import (
	"fmt"
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

type Chain33OfflineTx struct {
	ContractAddr  string
	TxHash        string
	SignedRawTx   string
	OperationName string
	Interval      time.Duration
}

func OfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "create and sign offline tx to deploy and set cross contracts to chain33",
	}
	cmd.AddCommand(
		createCrossBridgeCmd(),
		createValSetCmd(),
		createEthereumBridgeCmd(),
		createOracleCmd(),
		createBridgeBankCmd(),
		setBridgeBank2EthBridgeCmd(),
		setOracle2EthBridgeCmd(),
		createBridgeRegistryCmd(),
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

	fileName := "deployCrossX2Chain33.txt"
	fmt.Printf("%d: Write all the txs to file:%s \n", i, fileName)
	utils.WriteToFileInJson(fileName, txs)
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
		OperationName: "bridgeBank",
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
		OperationName: "oracle",
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
		OperationName: "valset",
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
		OperationName: "ethereumBridge",
		Interval:      time.Second * 5,
	}
	return ethereumBridgeTx, nil
}

func createBridgeBankCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridge-bank",
		Short: "create and sign offline bridge bank contract",
		Run:   createBridgeBank,
	}
	addCreateBridgeBankFlags(cmd)
	return cmd
}

func addCreateBridgeBankFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("operator", "o", "", "operator address")
	_ = cmd.MarkFlagRequired("operator")

	cmd.Flags().StringP("eth-bridge", "b", "", "address of eth-bridge contract ")
	_ = cmd.MarkFlagRequired("eth-bridge")

	cmd.Flags().StringP("oracle", "r", "", "address of oracle contract ")
	_ = cmd.MarkFlagRequired("oracle")
}

func createBridgeBank(cmd *cobra.Command, args []string) {
	operator, _ := cmd.Flags().GetString("operator")
	oracle, _ := cmd.Flags().GetString("oracle")
	ethereumBridge, _ := cmd.Flags().GetString("eth-bridge")

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
	createPara := fmt.Sprintf("%s,%s,%s", operator, oracle, ethereumBridge)
	content, err := utils.CreateContractAndSign(info, generated.BridgeBankBin, generated.BridgeBankABI, createPara, "bridgeBank")
	if nil == err {
		utils.WriteContractFile("./bridgeBank", content)
	}
}

func createValSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "valset",
		Short: "create and sign offline valset contract",
		Run:   createValset,
	}
	addCreateValsetFlags(cmd)
	return cmd
}

func addCreateValsetFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("contructParameter", "r", "", "contruct parameter, as: 'addr, [addr, addr, addr, addr], [25, 25, 25, 25]'")
	_ = cmd.MarkFlagRequired("contructParameter")
}

func createValset(cmd *cobra.Command, args []string) {
	contructParameter, _ := cmd.Flags().GetString("contructParameter")

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
	//	address[] memory _initValidators,
	//	uint256[] memory _initPowers
	//)
	//constructor(addr, [addr, addr, addr, addr], [25, 25, 25, 25])
	createPara := contructParameter
	content, err := utils.CreateContractAndSign(info, generated.ValsetBin, generated.ValsetABI, createPara, "valset")
	if nil == err {
		utils.WriteContractFile("./valset", content)
	}
}

func createEthereumBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eth-bridge",
		Short: "create and sign offline ethereum bridge contract",
		Run:   createEthereumBridge,
	}
	addCreateEthereumBridgeFlags(cmd)
	return cmd
}

func addCreateEthereumBridgeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("operator", "o", "", "operator address")
	_ = cmd.MarkFlagRequired("operator")

	cmd.Flags().StringP("valset", "v", "", "address of valset contract ")
	_ = cmd.MarkFlagRequired("valset")
}

func createOracleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "oracle",
		Short: "create and sign offline oracle contract",
		Run:   createOracle,
	}
	addCreateOracleFlags(cmd)
	return cmd
}

func addCreateOracleFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("operator", "o", "", "operator address")
	_ = cmd.MarkFlagRequired("operator")

	cmd.Flags().StringP("eth-bridge", "b", "", "address of eth-bridge contract ")
	_ = cmd.MarkFlagRequired("eth-bridge")

	cmd.Flags().StringP("valset", "v", "", "address of valset contract ")
	_ = cmd.MarkFlagRequired("valset")
}

func createOracle(cmd *cobra.Command, args []string) {
	operator, _ := cmd.Flags().GetString("operator")
	valset, _ := cmd.Flags().GetString("valset")
	ethereumBridge, _ := cmd.Flags().GetString("eth-bridge")

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
	//constructor(addr, [addr, addr, addr, addr], [25, 25, 25, 25])
	createPara := fmt.Sprintf("%s,%s,%s", operator, valset, ethereumBridge)
	content, err := utils.CreateContractAndSign(info, generated.OracleBin, generated.OracleABI, createPara, "oralce")
	if nil == err {
		utils.WriteContractFile("./oracle", content)
	}
}

func setBridgeBank2EthBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-bridgebank",
		Short: "set-bridgebank to contract:ethereum bank",
		Run:   setBridgeBank2EthBridge,
	}
	addSetBridgeBank2EthBridgeFlags(cmd)
	return cmd
}

func addSetBridgeBank2EthBridgeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("bridgebank", "b", "", "bridgebank address")
	_ = cmd.MarkFlagRequired("bridgebank")

	cmd.Flags().StringP("ethbridge", "e", "", "ethereum bridge address")
	_ = cmd.MarkFlagRequired("ethbridge")
}

func setBridgeBank2EthBridge(cmd *cobra.Command, args []string) {
	bridgebank, _ := cmd.Flags().GetString("bridgebank")
	ethbridge, _ := cmd.Flags().GetString("ethbridge")

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
		return
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, txHash, err := utils.CallContractAndSign(info, action, ethbridge)
	if nil == err {
		utils.WriteContractFile("./addPool", content)
	}
}

func setOracle2EthBridgeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-Oracle",
		Short: "set Oracle to contract:ethereum bank",
		Run:   setOracle2EthBridge,
	}
	addSetOracle2EthBridgeFlags(cmd)
	return cmd
}

func addSetOracle2EthBridgeFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("oracle", "o", "", "oracle address")
	_ = cmd.MarkFlagRequired("oracle")

	cmd.Flags().StringP("ethbridge", "e", "", "ethereum bridge address")
	_ = cmd.MarkFlagRequired("ethbridge")
}

func setOracle2EthBridge(cmd *cobra.Command, args []string) {
	oracle, _ := cmd.Flags().GetString("oracle")
	ethbridge, _ := cmd.Flags().GetString("ethbridge")

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
		return
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, err := utils.CallContractAndSign(info, action, ethbridge)
	if nil == err {
		utils.WriteContractFile("./setOracle", content)
	}
}

func createBridgeRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridgeRegistry",
		Short: "create and sign offline BridgeRegistry contract",
		Run:   createBridgeRegistry,
	}
	addCreateBridgeRegistryFlags(cmd)
	return cmd
}

func addCreateBridgeRegistryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("ethereumBridge", "e", "", "ethereum bridge address")
	_ = cmd.MarkFlagRequired("ethereumBridge")

	cmd.Flags().StringP("bridgeBank", "b", "", "address of bridgeBank contract ")
	_ = cmd.MarkFlagRequired("bridgeBank")

	cmd.Flags().StringP("valset", "v", "", "address of valset contract ")
	_ = cmd.MarkFlagRequired("valset")

	cmd.Flags().StringP("oracle", "o", "", "address of oracle contract ")
	_ = cmd.MarkFlagRequired("oracle")
}

func createBridgeRegistry(cmd *cobra.Command, args []string) {
	ethereumBridge, _ := cmd.Flags().GetString("ethereumBridge")
	valset, _ := cmd.Flags().GetString("valset")
	bridgeBank, _ := cmd.Flags().GetString("bridgeBank")
	oracle, _ := cmd.Flags().GetString("oracle")

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
	content, err := utils.CreateContractAndSign(info, generated.BridgeRegistryBin, generated.BridgeRegistryABI, createPara, "BridgeRegistry")
	if nil == err {
		utils.WriteContractFile("./BridgeRegistry", content)
	}
}
