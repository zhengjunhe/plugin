package commands

import (
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/spf13/cobra"
	"os"
)

func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query x2ethereum",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(queryEthProphecyCmd(), queryValidatorsCmd(), queryConsensusCmd(), queryTotalPowerCmd())
	return cmd
}

func queryEthProphecyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prophecy",
		Short: "query prophecy",
		Run:   queryEthProphecy,
	}

	cmd.Flags().StringP("id", "i", "", "prophecy id")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}

func queryEthProphecy(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	id, _ := cmd.Flags().GetString("id")

	get := &types2.QueryEthProphecyParams{
		ID: id,
	}

	payLoad, err := types.PBToJSON(get)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ErrPbToJson:"+err.Error())
		return
	}

	query := rpctypes.Query4Jrpc{
		Execer:   types2.X2ethereumX,
		FuncName: types2.FuncQueryEthProphecy,
		Payload:  payLoad,
	}

	channel := &types2.ReceiptEthProphecy{}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", query, channel)
	ctx.Run()
}

func queryValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "query current validators",
		Run:   queryValidators,
	}
	cmd.Flags().StringP("validator", "v", "", "write if you want to check specific validator")
	//_ = cmd.MarkFlagRequired("validator")
	return cmd
}

func queryValidators(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	validator, _ := cmd.Flags().GetString("validator")

	get := &types2.QueryValidatorsParams{
		Validator: validator,
	}

	payLoad, err := types.PBToJSON(get)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ErrPbToJson:"+err.Error())
		return
	}

	query := rpctypes.Query4Jrpc{
		Execer:   types2.X2ethereumX,
		FuncName: types2.FuncQueryValidators,
		Payload:  payLoad,
	}

	channel := &types2.ReceiptQueryValidator{}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", query, channel)
	ctx.Run()
}

func queryConsensusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consensus",
		Short: "query current consensus need",
		Run:   queryConsensus,
	}
	return cmd
}

func queryConsensus(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	query := rpctypes.Query4Jrpc{
		Execer:   types2.X2ethereumX,
		FuncName: types2.FuncQueryConsensusNeeded,
	}

	channel := &types2.ReceiptQueryConsensusNeeded{}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", query, channel)
	ctx.Run()
}

func queryTotalPowerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "totalpower",
		Short: "query current total power",
		Run:   queryTotalPower,
	}
	return cmd
}

func queryTotalPower(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	query := rpctypes.Query4Jrpc{
		Execer:   types2.X2ethereumX,
		FuncName: types2.FuncQueryTotalPower,
	}

	channel := &types2.ReceiptQueryTotalPower{}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", query, channel)
	ctx.Run()
}
