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
		TxCmd(),      //构造交易
		SignCmd(),    //签名交易
		sendTxsCmd(), //发送交易

	)

	return cmd
}
