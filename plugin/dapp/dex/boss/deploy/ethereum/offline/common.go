package offline

import "github.com/spf13/cobra"

func EthOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eth",
		Short: "create and sign offline tx to deploy and set dex contracts to eth",
		Args:  cobra.MinimumNArgs(1),
	}
	var query =new(queryCmd)
	var signFac =new(SignFactoryCmd)
	var deployFac=new(deploayFactory)
	cmd.AddCommand(
		query.queryAddrInfoCmd(),
		signFac.signFactoryCmd(),     //step1
		deployFac.deployFactoryCmd(), //send tx
		//signWth9Cmd(),//step2
		//signPancakeRouter(),//step3
	)
	return cmd
}
