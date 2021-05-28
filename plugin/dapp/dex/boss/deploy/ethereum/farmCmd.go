package ethereum

import (
	"fmt"

	"github.com/spf13/cobra"
)

func FarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "farm",
		Short: "farm command",
	}
	cmd.AddCommand(
		DeployFarmCmd(),
		AddPoolCmd(),
		UpdateAllocPointCmd(),
		TransferOwnerShipCmd(),
		ShowCackeBalanceCmd(),
	)
	return cmd
}

func ShowCackeBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "show cake balance within a specified pool",
		Run:   ShowCakeBalance,
	}
	ShowBalanceFlags(cmd)
	return cmd
}

//GetBalanceFlags ...
func ShowBalanceFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("owner", "o", "", "owner address")
	_ = cmd.MarkFlagRequired("owner")

	cmd.Flags().Int64P("pid", "d", 0, "id of pool")
	_ = cmd.MarkFlagRequired("pid")
}

//GetBalance ...
func ShowCakeBalance(cmd *cobra.Command, args []string) {
	owner, _ := cmd.Flags().GetString("owner")
	pid, _ := cmd.Flags().GetInt64("pid")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)
	balance, err := GetCakeBalance(owner, pid)
	if nil != err {
		fmt.Println("err:", err.Error())
	}
	fmt.Println("balance =", balance)
}

func DeployFarmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy farm",
		Short: "deploy farm to bsc ",
		Run:   DeployContracts,
	}
	return cmd
}

func DeployContracts(cmd *cobra.Command, args []string) {
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)
	err := DeployFarm()
	if nil != err {
		fmt.Println("Failed to deploy contracts due to:", err.Error())
		return
	}
	fmt.Println("Succeed to deploy contracts")
}

func AddPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add pool",
		Short: "add pool to farm ",
		Run:   AddPool2Farm,
	}

	addAddPoolCmdFlags(cmd)

	return cmd
}

func addAddPoolCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().StringP("lptoken", "l", "", "lp Addr ")
	_ = cmd.MarkFlagRequired("lptoken")

	cmd.Flags().Int64P("alloc", "p", 0, "allocation point ")
	_ = cmd.MarkFlagRequired("alloc")

	cmd.Flags().BoolP("update", "u", true, "with update")
	_ = cmd.MarkFlagRequired("update")
}

func AddPool2Farm(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	allocPoint, _ := cmd.Flags().GetInt64("alloc")
	lpToken, _ := cmd.Flags().GetString("lptoken")
	update, _ := cmd.Flags().GetBool("update")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	err := AddPool2FarmHandle(masterChefAddrStr, allocPoint, lpToken, update)
	if nil != err {
		fmt.Println("Failed to AddPool2Farm due to:", err.Error())
		return
	}
	fmt.Println("Succeed to AddPool2Farm")
}

func UpdateAllocPointCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update alloc point",
		Short: "Update the given pool's CAKE allocation point",
		Run:   UpdateAllocPoint,
	}

	updateAllocPointCmdFlags(cmd)

	return cmd
}

func updateAllocPointCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().Int64P("pid", "d", 0, "id of pool")
	_ = cmd.MarkFlagRequired("pid")

	cmd.Flags().Int64P("alloc", "p", 0, "allocation point ")
	_ = cmd.MarkFlagRequired("alloc")

	cmd.Flags().BoolP("update", "u", true, "with update")
	_ = cmd.MarkFlagRequired("update")
}

func UpdateAllocPoint(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	pid, _ := cmd.Flags().GetInt64("pid")
	allocPoint, _ := cmd.Flags().GetInt64("alloc")
	update, _ := cmd.Flags().GetBool("update")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	err := UpdateAllocPointHandle(masterChefAddrStr, pid, allocPoint, update)
	if nil != err {
		fmt.Println("Failed to AddPool2Farm due to:", err.Error())
		return
	}
	fmt.Println("Succeed to AddPool2Farm")
}

func TransferOwnerShipCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Transfer OwnerShip",
		Short: "Transfer OwnerShip",
		Run:   TransferOwnerShip,
	}

	TransferOwnerShipFlags(cmd)

	return cmd
}

func TransferOwnerShipFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("new", "n", "", "new owner")
	_ = cmd.MarkFlagRequired("new")

	cmd.Flags().StringP("contract", "c", "", "contract address")
	_ = cmd.MarkFlagRequired("contract")
}

func TransferOwnerShip(cmd *cobra.Command, args []string) {
	newOwner, _ := cmd.Flags().GetString("new")
	contract, _ := cmd.Flags().GetString("contract")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	err := TransferOwnerShipHandle(newOwner, contract)
	if nil != err {
		fmt.Println("Failed to TransferOwnerShip due to:", err.Error())
		return
	}
	fmt.Println("Succeed to TransferOwnerShip")
}
