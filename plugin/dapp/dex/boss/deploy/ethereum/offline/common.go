package offline

import "github.com/spf13/cobra"

func EthOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eth",
		Short: "create and sign offline tx to deploy and set dex contracts to eth",
		Args:  cobra.MinimumNArgs(1),
	}
	var query = new(queryCmd)
	var sign = new(SignCmd)
	var deploy = new(deploayContract)
	cmd.AddCommand(
		query.queryCmd(),   //query fromAccount info such as: nonce,gasprice
		sign.signCmd(),     //sign fatory.weth9,pancakrouter contract
		deploy.deployCmd(), //send singned tx to deploy contract:factory,weth9,pancakerouter.

	)
	return cmd
}
