package offline

import (
	"github.com/spf13/cobra"
)

const gasLimit uint64 = 50000 //10000 * 800

func OfflineDeployContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "deploy the corresponding Ethereum contracts",
		Args:  cobra.MinimumNArgs(1),
	}
	var query = new(queryCmd)
	var sign = new(SignCmd)
	//var deploy = new(eoff.DeployContract)
	cmd.AddCommand(
		query.queryCmd(),
		sign.signCmd(),
		sendTxsCmd(),
		//deploy.DeployCmd(), //send singned tx
	)

	return cmd
}
