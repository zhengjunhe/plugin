package offline

import (
	"github.com/spf13/cobra"
)

func OfflineDeployContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "deploy the corresponding Ethereum contracts",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		TxCmd(),
		SignCmd(),
		sendTxsCmd(),

	)

	return cmd
}
