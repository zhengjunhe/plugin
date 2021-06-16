package offline

import (
	eoff "github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum/offline"
	"github.com/spf13/cobra"
)

const gasLimit uint64 = 10000 * 800

func OfflineDeployContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline deploy",
		Short: "deploy the corresponding Ethereum contracts",
		Args:  cobra.MinimumNArgs(1),
	}
	var query = new(queryCmd)
	var sign = new(SignCmd)
	var deploy = new(eoff.DeployContract)
	cmd.AddCommand(
		query.queryCmd(),
		sign.signCmd(),
		deploy.DeployCmd(), //send singned tx
	)

	return cmd
}
