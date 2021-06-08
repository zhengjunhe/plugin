package ethereum

import (
	"encoding/json"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-farm/src/masterChef"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"

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
		ShowPoolInfosCmd(),
	)
	return cmd
}

func ShowPoolInfosCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool info",
		Short: "query pool length,and pool info by pid",
		Run:   showPools,
	}

	ShowPoolInfoFlags(cmd)
	return cmd
}

func ShowPoolInfoFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("lp", "l", "", "lp address")
	cmd.Flags().Int64P("poolid", "p", 0, "pool id")
	//_ = cmd.MarkFlagRequired("lp")
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

}

func showPools(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	lpAddrStr, _ := cmd.Flags().GetString("lp")
	poolid, _ := cmd.Flags().GetInt64("poolid")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	setupWebsocketEthClient(ethNodeAddr)
	masterChefAddr := common.HexToAddress(masterChefAddrStr)

	masterChefInt, err := masterChef.NewMasterChef(masterChefAddr, ethClient)
	if nil != err {
		return
	}
	var opts bind.CallOpts
	pl, err := masterChefInt.PoolLength(&opts)
	if err != nil {
		fmt.Println("query masterChef PoolLength err", err.Error())
		return
	}
	fmt.Println("++++++++++++++++\ntotal pool num:", pl.Int64(), "\n++++++++++++++++\n")
	//var pid int64 =1
	totalPid := pl.Int64()
	for pid := 1; pid < int(totalPid); pid++ {
		info, err := masterChefInt.PoolInfo(&opts, big.NewInt(int64(pid)))
		if err != nil {
			fmt.Println("query poolinfo err", err.Error(), "pid", pid)
			continue
		}
		jinfo, _ := json.MarshalIndent(info, "", "\t")
		if lpAddrStr != "" {
			//find lpaddr-->pid
			if strings.ToLower(info.LpToken.String()) == strings.ToLower(lpAddrStr) {

				fmt.Println("Find LP PID:", pid, "\nLP info:", string(jinfo))
				return
			}
			continue
		}
		if poolid != 0 {
			//find pid--->lpaddr
			if pid == int(poolid) {
				fmt.Println("Find LP PID:", pid, "\nLP info:", string(jinfo))
				return
			}
			continue
		}
		fmt.Println("LP PID:", pid, "\nLP info:", string(jinfo))

	}

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
	cmd.Flags().Uint64P("gasLimit", "g", 80*10000, "set gas limit")
}

func AddPool2Farm(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	allocPoint, _ := cmd.Flags().GetInt64("alloc")
	lpToken, _ := cmd.Flags().GetString("lptoken")
	update, _ := cmd.Flags().GetBool("update")
	gasLimit, _ := cmd.Flags().GetUint64("gaslimit")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	err := AddPool2FarmHandle(masterChefAddrStr, allocPoint, lpToken, update, gasLimit)
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
