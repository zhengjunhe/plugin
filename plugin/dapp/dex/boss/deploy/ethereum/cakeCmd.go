package ethereum

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd x2ethereum client command
func CakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cake",
		Short: "cake command",
	}
	cmd.AddCommand(
		GetBalanceCmd(),
		DeployPancakeCmd(),
		AddAllowance4LPCmd(),
		CheckAllowance4LPCmd(),
	)
	return cmd
}

func DeployPancakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy pancake",
		Short: "deploy pancake router to ethereum ",
		Run:   DeployContractsCake,
	}
	return cmd
}

func DeployContractsCake(cmd *cobra.Command, args []string) {
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)
	err := DeployPancake()
	if nil != err {
		fmt.Println("Failed to deploy contracts due to:", err.Error())
		return
	}
	fmt.Println("Succeed to deploy contracts")
}

func AddAllowance4LPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "allowance",
		Short: "approve allowance for add lp to pool",
		Run:   AddAllowance4LP,
	}

	AddAllowance4LPFlags(cmd)

	return cmd
}

func AddAllowance4LPFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().StringP("lptoken", "l", "", "lp Addr ")
	_ = cmd.MarkFlagRequired("lptoken")

	cmd.Flags().Int64P("amount", "p", 0, "amount to approve")
	_ = cmd.MarkFlagRequired("amount")
}

func AddAllowance4LP(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	amount, _ := cmd.Flags().GetInt64("amount")
	lpToken, _ := cmd.Flags().GetString("lptoken")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	//owner string, spender string, amount int64
	err := AddAllowance4LPHandle(lpToken, masterChefAddrStr, amount)
	if nil != err {
		fmt.Println("Failed to AddPool2Farm due to:", err.Error())
		return
	}
	fmt.Println("Succeed to AddPool2Farm")
}

func CheckAllowance4LPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check-allowance",
		Short: "check allowance for add lp to pool",
		Run:   CheckAllowance4LP,
	}

	CheckAllowance4LPFlags(cmd)

	return cmd
}

func CheckAllowance4LPFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("masterchef", "m", "", "master Chef Addr ")
	_ = cmd.MarkFlagRequired("masterchef")

	cmd.Flags().StringP("lptoken", "l", "", "lp Addr ")
	_ = cmd.MarkFlagRequired("lptoken")
}

func CheckAllowance4LP(cmd *cobra.Command, args []string) {
	masterChefAddrStr, _ := cmd.Flags().GetString("masterchef")
	lpToken, _ := cmd.Flags().GetString("lptoken")

	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)

	//owner string, spender string, amount int64
	err := CheckAllowance4LPHandle(lpToken, masterChefAddrStr)
	if nil != err {
		fmt.Println("Failed to CheckAllowance4LP due to:", err.Error())
		return
	}
	fmt.Println("Succeed to CheckAllowance4LP")
}
