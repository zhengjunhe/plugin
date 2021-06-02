package offline

import (
	"github.com/spf13/cobra"
)

func Chain33OfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain33",
		Short: "create and sign offline tx to deploy and set dex contracts to chain33",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		createERC20ContractCmd(),
		createFactoryCmd(),
		createWeth9Cmd(),
		createRouterCmd(),
		farmofflineCmd(),
	)
	return cmd
}
