package ethereum

import (
	"fmt"

	"github.com/spf13/cobra"
)

func EthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ethereum",
		Short: "ethereumcake command",
	}
	cmd.AddCommand(
		CakeCmd(),
		FarmCmd(),
		GetBalanceCmd(),
	)
	return cmd
}

//GetBalanceCmd ...
func GetBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "get owner's balance for ETH or ERC20",
		Run:   ShowBalance,
	}
	GetBalanceFlags(cmd)
	return cmd
}

//GetBalanceFlags ...
func GetBalanceFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("owner", "o", "", "owner address")
	_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringP("tokenAddr", "t", "", "token address, optional, nil for Eth")
}

//GetBalance ...
func ShowBalance(cmd *cobra.Command, args []string) {
	owner, _ := cmd.Flags().GetString("owner")
	tokenAddr, _ := cmd.Flags().GetString("tokenAddr")
	ethNodeAddr, _ := cmd.Flags().GetString("rpc_laddr_ethereum")

	setupWebsocketEthClient(ethNodeAddr)
	balance, err := GetBalance(tokenAddr, owner)
	if nil != err {
		fmt.Println("err:", err.Error())
	}
	fmt.Println("balance =", balance)
}
