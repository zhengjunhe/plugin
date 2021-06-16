package offline

import (
	"fmt"

	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4chain33/generated"
	"github.com/33cn/plugin/plugin/dapp/dex/utils"
	"github.com/spf13/cobra"
)

func OfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "create and sign offline tx to deploy and set cross contracts to chain33",
	}
	cmd.AddCommand(
		createValSetCmd(),
		createEthereumBridgeCmd(),
		createBridgeBankCmd(),
		createOracleCmd(),
		setBridgeBank2EthBridgeCmd(),
		setOracle2EthBridgeCmd(),
		createBridgeRegistryCmd(),
	)
	return cmd
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

func createEthereumBridge(cmd *cobra.Command, args []string) {
	operator, _ := cmd.Flags().GetString("operator")
	valset, _ := cmd.Flags().GetString("valset")

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
	createPara := fmt.Sprintf("%s,%s", operator, valset)
	content, err := utils.CreateContractAndSign(info, generated.EthereumBridgeBin, generated.EthereumBridgeABI, createPara, "EthereumBridge")
	if nil == err {
		utils.WriteContractFile("./EthereumBridge", content)
	}
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
	content, err := utils.CallContractAndSign(info, action, ethbridge)
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
