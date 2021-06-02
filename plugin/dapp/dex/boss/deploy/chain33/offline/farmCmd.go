package offline

import (
	"fmt"

	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-farm/src/cakeToken"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-farm/src/masterChef"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-farm/src/syrupBar"
	"github.com/33cn/plugin/plugin/dapp/dex/utils"
	"github.com/spf13/cobra"
)

func farmofflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farm",
		Short: "create and sign tx to deploy farm and set lp, transfer ownership",
	}
	cmd.AddCommand(
		createCakeTokenCmd(),
		createSyrupBarCmd(),
		createMasterChefCmd(),
		AddPoolCmd(),
		updateAllocPointCmd(),
	)
	return cmd
}

func createCakeTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cakeToken",
		Short: "create cakeToken contract",
		Run:   createCakeToken,
	}
	addCreateCakeTokenFlags(cmd)
	return cmd
}

func addCreateCakeTokenFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller's private key")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")
}

func createCakeToken(cmd *cobra.Command, args []string) {
	privateKey, _ := cmd.Flags().GetString("caller")
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
	createPara := ""
	content, err := utils.CreateContractAndSign(info, cakeToken.CakeTokenBin, cakeToken.CakeTokenABI, createPara, "cakeToken")
	if nil != err {
		utils.WriteContractFile("./cakeToken", content)
	}
}

func createSyrupBarCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "syrupBar",
		Short: "create syrupBar contract",
		Run:   createSyrupBar,
	}
	addCreateSyrupBarFlags(cmd)
	return cmd
}

func addCreateSyrupBarFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller's private key")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("cakeToken", "a", "", "address of cake token")
	cmd.MarkFlagRequired("cakeToken")
}

func createSyrupBar(cmd *cobra.Command, args []string) {
	cakeToken, _ := cmd.Flags().GetString("cakeToken")

	privateKey, _ := cmd.Flags().GetString("caller")
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
	//constructor(CakeToken _cake)
	createPara := cakeToken
	content, err := utils.CreateContractAndSign(info, syrupBar.SyrupBarBin, syrupBar.SyrupBarABI, createPara, "syrupBar")
	if nil != err {
		utils.WriteContractFile("./syrupBar", content)
	}
}

func createMasterChefCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "masterChef",
		Short: "create masterChef contract",
		Run:   createMasterChef,
	}
	addCreateMasterChefFlags(cmd)
	return cmd
}

func addCreateMasterChefFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller's private key")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("cakeToken", "a", "", "address of cake token")
	cmd.MarkFlagRequired("cakeToken")
	cmd.Flags().StringP("syrup", "s", "", "address of syrup")
	cmd.MarkFlagRequired("syrup")
	cmd.Flags().StringP("devaddr", "d", "", "address of develop")
	cmd.MarkFlagRequired("devaddr")
	cmd.Flags().Int64P("cakePerBlock", "m", 0, "cake Per Block, should multiply 1e18")
	cmd.MarkFlagRequired("cakePerBlock")
	cmd.Flags().Int64P("startBlock", "h", 0, "start Block height")
	cmd.MarkFlagRequired("startBlock")
}

func createMasterChef(cmd *cobra.Command, args []string) {
	cakeToken, _ := cmd.Flags().GetString("cakeToken")
	syrup, _ := cmd.Flags().GetString("syrup")
	devaddr, _ := cmd.Flags().GetString("devaddr")
	cakePerBlock, _ := cmd.Flags().GetInt64("cakePerBlock")
	startBlock, _ := cmd.Flags().GetInt64("startBlock")

	privateKey, _ := cmd.Flags().GetString("caller")
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
	//	CakeToken _cake,
	//	SyrupBar _syrup,
	//	address _devaddr,
	//	uint256 _cakePerBlock,
	//	uint256 _startBlock
	//) public {
	createPara := fmt.Sprintf("%s,%s,%s,%d,%d", cakeToken, syrup, devaddr, cakePerBlock, startBlock)
	content, err := utils.CreateContractAndSign(info, masterChef.MasterChefBin, masterChef.MasterChefABI, createPara, "masterChef")
	if nil != err {
		utils.WriteContractFile("./masterChef", content)
	}
}

func AddPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addPool",
		Short: "add lp to pool",
		Run:   addPool,
	}
	addPoolFlags(cmd)
	return cmd
}

func addPoolFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().StringP("lptoken", "l", "", "lp Addr ")
	_ = cmd.MarkFlagRequired("lptoken")

	cmd.Flags().Int64P("alloc", "p", 0, "allocation point ")
	_ = cmd.MarkFlagRequired("alloc")

	cmd.Flags().BoolP("update", "u", true, "with update")
	_ = cmd.MarkFlagRequired("update")

	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")
	cmd.MarkFlagRequired("fee")

	cmd.Flags().StringP("caller", "c", "", "caller address")
	_ = cmd.MarkFlagRequired("caller")
}

func addPool(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	allocPoint, _ := cmd.Flags().GetInt64("alloc")
	lpToken, _ := cmd.Flags().GetString("lptoken")
	update, _ := cmd.Flags().GetBool("update")

	privateKey, _ := cmd.Flags().GetString("caller")
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
	parameter := fmt.Sprintf("add(%d, %s, %v)", allocPoint, lpToken, update)
	_, packData, err := evmAbi.Pack(parameter, masterChef.MasterChefABI, false)
	if nil != err {
		fmt.Println("AddPool2FarmHandle", "Failed to do abi.Pack due to:", err.Error())
		return
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, err := utils.CallContractAndSign(info, action, masterChefAddrStr)
	if nil != err {
		utils.WriteContractFile("./addPool", content)
	}
}

func updateAllocPointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateAllocPoint",
		Short: "update Alloc Point",
		Run:   updateAllocPoint,
	}
	addUpdateAllocPointFlags(cmd)
	return cmd
}

func addUpdateAllocPointFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().Int64P("pid", "d", 0, "id of pool")
	_ = cmd.MarkFlagRequired("pid")

	cmd.Flags().Int64P("alloc", "p", 0, "allocation point ")
	_ = cmd.MarkFlagRequired("alloc")

	cmd.Flags().BoolP("update", "u", true, "with update")
	_ = cmd.MarkFlagRequired("update")

	cmd.Flags().StringP("caller", "c", "", "caller address")
	_ = cmd.MarkFlagRequired("caller")
}

func updateAllocPoint(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	pid, _ := cmd.Flags().GetInt64("pid")
	allocPoint, _ := cmd.Flags().GetInt64("alloc")
	update, _ := cmd.Flags().GetBool("update")

	privateKey, _ := cmd.Flags().GetString("caller")
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
	parameter := fmt.Sprintf("set(%d, %d, %v)", pid, allocPoint, update)
	_, packData, err := evmAbi.Pack(parameter, masterChef.MasterChefABI, false)
	if nil != err {
		fmt.Println("UpdateAllocPoint", "Failed to do abi.Pack due to:", err.Error())
		return
	}
	action := &evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: parameter, Para: packData}
	content, err := utils.CallContractAndSign(info, action, masterChefAddrStr)
	if nil != err {
		utils.WriteContractFile("./updateAllocPoint", content)
	}
}
