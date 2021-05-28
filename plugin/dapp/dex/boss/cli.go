package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/33cn/plugin/plugin/dapp/dex/boss/buildFlags"
	"github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/chain33"
	"github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum"
	"github.com/spf13/cobra"
)

func main() {
	if buildFlags.RPCAddr4Chain33 == "" {
		buildFlags.RPCAddr4Chain33 = "http://localhost:8801"
	}
	buildFlags.RPCAddr4Chain33 = testTLS(buildFlags.RPCAddr4Chain33)

	if buildFlags.RPCAddr4Ethereum == "" {
		buildFlags.RPCAddr4Ethereum = "https://data-seed-prebsc-1-s1.binance.org:8545"
	}

	rootCmd := Cmd()
	rootCmd.PersistentFlags().String("rpc_laddr", buildFlags.RPCAddr4Chain33, "http url")
	rootCmd.PersistentFlags().String("rpc_laddr_ethereum", buildFlags.RPCAddr4Ethereum, "http url")
	rootCmd.PersistentFlags().String("paraName", "", "para chain name,Eg:user.p.fzm.")
	rootCmd.PersistentFlags().String("expire", "120s", "transaction expire time (optional)")
	rootCmd.PersistentFlags().Int32("chainID", 0, "chain id, default to 0")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Cmd x2ethereum client command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy to ethereum or chain33",
	}
	cmd.AddCommand(
		ethereum.EthCmd(),
		chain33.Chain33Cmd(),
	)
	return cmd
}

func testTLS(RPCAddr string) string {
	rpcaddr := RPCAddr
	if !strings.HasPrefix(rpcaddr, "http://") {
		return RPCAddr
	}
	// if http://
	if rpcaddr[len(rpcaddr)-1] != '/' {
		rpcaddr += "/"
	}
	rpcaddr += "test"
	/* #nosec */
	resp, err := http.Get(rpcaddr)
	if err != nil {
		return "https://" + RPCAddr[7:]
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return RPCAddr
	}
	return "https://" + RPCAddr[7:]
}
