/*Package commands implement dapp client commands*/
package commands

import (
	"github.com/33cn/chain33/rpc/jsonclient"
	types2 "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	types3 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/spf13/cobra"
)

/*
 * 实现合约对应客户端
 */

//TODO
// 在本地维护一张不同token的unit表格

// Cmd x2ethereum client command
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "x2ethereum",
		Short: "x2ethereum command",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		CreateRawEth2Chain33TxCmd(),
		CreateRawWithdrawEthTxCmd(),
		CreateRawWithdrawChain33TxCmd(),
		CreateRawChain33ToEthTxCmd(),
		CreateRawAddValidatorTxCmd(),
		CreateRawRemoveValidatorTxCmd(),
		CreateRawModifyValidatorTxCmd(),
		CreateRawSetConsensusTxCmd(),
		queryCmd(),
	)
	return cmd
}

// Eth2Chain33
func CreateRawEth2Chain33TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Eth2Chain33",
		Run:   Eth2Chain33,
	}

	addEth2Chain33Flags(cmd)
	return cmd
}

func addEth2Chain33Flags(cmd *cobra.Command) {
	cmd.Flags().Int64("ethid", 0, "the ethereum chain ID which send asset to chain33")
	_ = cmd.MarkFlagRequired("ethid")

	cmd.Flags().StringP("bcontract", "b", "", "BridgeContractAddress")
	_ = cmd.MarkFlagRequired("bcontract")

	cmd.Flags().Int64("nonce", 0, "the nonce for this tx in ethereum")
	_ = cmd.MarkFlagRequired("nonce")

	cmd.Flags().StringP("csymbol", "t", "", "token symbol in chain33")
	_ = cmd.MarkFlagRequired("csymbol")

	cmd.Flags().StringP("cexec", "e", "", "chain execer in chain33")
	_ = cmd.MarkFlagRequired("cexec")

	cmd.Flags().StringP("tcontract", "q", "", "token contract address in ethereum")
	_ = cmd.MarkFlagRequired("tcontract")

	cmd.Flags().StringP("sender", "s", "", "ethereum sender address")
	_ = cmd.MarkFlagRequired("sender")

	cmd.Flags().StringP("receiver", "r", "", "chain33 receiver address")
	_ = cmd.MarkFlagRequired("cExec")

	cmd.Flags().StringP("validator", "v", "", "validator address")
	_ = cmd.MarkFlagRequired("validator")

	cmd.Flags().Float64P("amount", "a", float64(0), "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")

	cmd.Flags().Int64("claimtype", 0, "the type of this claim,lock=1,burn=2")
	_ = cmd.MarkFlagRequired("claimtype")

}

func Eth2Chain33(cmd *cobra.Command, args []string) {
	ethid, _ := cmd.Flags().GetInt64("ethid")
	bcontract, _ := cmd.Flags().GetString("bcontract")
	nonce, _ := cmd.Flags().GetInt64("nonce")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	tcontract, _ := cmd.Flags().GetString("tcontract")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	validator, _ := cmd.Flags().GetString("validator")
	amount, _ := cmd.Flags().GetFloat64("amount")
	claimtype, _ := cmd.Flags().GetInt64("claimtype")

	//todo
	// 这边暂时默认是ethereum，即decimals为18
	params := &types3.Eth2Chain33{
		EthereumChainID:       ethid,
		BridgeContractAddress: bcontract,
		Nonce:                 nonce,
		LocalCoinSymbol:       csymbol,
		LocalCoinExec:         cexec,
		TokenContractAddress:  tcontract,
		EthereumSender:        sender,
		Chain33Receiver:       receiver,
		ValidatorAddress:      validator,
		Amount:                uint64(amount * 1e18),
		ClaimType:             claimtype,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameEth2Chain33Action)
}

// WithdrawEth
func CreateRawWithdrawEthTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraweth",
		Short: "withdraw a Eth2Chain33",
		Run:   WithdrawEth,
	}

	addEth2Chain33Flags(cmd)
	return cmd
}

func WithdrawEth(cmd *cobra.Command, args []string) {
	ethid, _ := cmd.Flags().GetInt64("ethid")
	bcontract, _ := cmd.Flags().GetString("bcontract")
	nonce, _ := cmd.Flags().GetInt64("nonce")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	tcontract, _ := cmd.Flags().GetString("tcontract")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	validator, _ := cmd.Flags().GetString("validator")
	amount, _ := cmd.Flags().GetFloat64("amount")
	claimtype, _ := cmd.Flags().GetInt64("claimtype")

	params := &types3.Eth2Chain33{
		EthereumChainID:       ethid,
		BridgeContractAddress: bcontract,
		Nonce:                 nonce,
		LocalCoinSymbol:       csymbol,
		LocalCoinExec:         cexec,
		TokenContractAddress:  tcontract,
		EthereumSender:        sender,
		Chain33Receiver:       receiver,
		ValidatorAddress:      validator,
		Amount:                uint64(amount * 1e8),
		ClaimType:             claimtype,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameWithdrawEthAction)
}

// Burn
func CreateRawWithdrawChain33TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn",
		Short: "Create a burn tx in chain33,withdraw chain33ToEth",
		Run:   burn,
	}

	addChain33ToEthFlags(cmd)
	return cmd
}

func addChain33ToEthFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("tcontract", "q", "", "token contract address")
	_ = cmd.MarkFlagRequired("tcontract")

	cmd.Flags().StringP("csymbol", "t", "", "token symbol in chain33")
	_ = cmd.MarkFlagRequired("csymbol")

	cmd.Flags().StringP("cexec", "e", "", "chain execer in chain33")
	_ = cmd.MarkFlagRequired("cexec")

	cmd.Flags().StringP("sender", "s", "", "chain33 sender address")
	_ = cmd.MarkFlagRequired("sender")

	cmd.Flags().StringP("receiver", "r", "", "ethereum receiver address")
	_ = cmd.MarkFlagRequired("cExec")

	cmd.Flags().Float64P("amount", "a", float64(0), "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")

}

func burn(cmd *cobra.Command, args []string) {
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetFloat64("amount")
	tcontract, _ := cmd.Flags().GetString("tcontract")

	params := &types3.Chain33ToEth{
		TokenContract:    tcontract,
		Chain33Sender:    sender,
		EthereumReceiver: receiver,
		Amount:           uint64(amount * 1e8),
		LocalCoinSymbol:  csymbol,
		LocalCoinExec:    cexec,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameWithdrawChain33Action)
}

// Lock
func CreateRawChain33ToEthTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Create a lock tx in chain33,create a chain33ToEth tx",
		Run:   lock,
	}

	addChain33ToEthFlags(cmd)
	return cmd
}

func lock(cmd *cobra.Command, args []string) {
	contract, _ := cmd.Flags().GetString("tcontract")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetFloat64("amount")

	params := &types3.Chain33ToEth{
		TokenContract:    contract,
		Chain33Sender:    sender,
		EthereumReceiver: receiver,
		Amount:           uint64(amount * 1e8),
		LocalCoinSymbol:  csymbol,
		LocalCoinExec:    cexec,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameChain33ToEthAction)
}

// AddValidator
func CreateRawAddValidatorTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Create a add validator tx in chain33",
		Run:   addValidator,
	}

	addValidatorFlags(cmd)
	cmd.Flags().Int64P("power", "p", 0, "validator power set")
	_ = cmd.MarkFlagRequired("power")
	return cmd
}

func addValidatorFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", "", "the address you want to add/remove/modify as a validator ")
	_ = cmd.MarkFlagRequired("address")
}

func addValidator(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")
	power, _ := cmd.Flags().GetInt64("power")

	params := &types3.MsgValidator{
		Address: address,
		Power:   power,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameAddValidatorAction)
}

// RemoveValidator
func CreateRawRemoveValidatorTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Create a remove validator tx in chain33",
		Run:   removeValidator,
	}

	addValidatorFlags(cmd)
	return cmd
}

func removeValidator(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")

	params := &types3.MsgValidator{
		Address: address,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameRemoveValidatorAction)
}

// ModifyValidator
func CreateRawModifyValidatorTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "modify",
		Short: "Create a modify validator tx in chain33",
		Run:   modify,
	}

	addValidatorFlags(cmd)

	cmd.Flags().Int64P("power", "p", 0, "validator power set")
	_ = cmd.MarkFlagRequired("power")
	return cmd
}

func modify(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")
	power, _ := cmd.Flags().GetInt64("power")

	params := &types3.MsgValidator{
		Address: address,
		Power:   power,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameModifyPowerAction)
}

// MsgSetConsensusNeeded
func CreateRawSetConsensusTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setconsensus",
		Short: "Create a set consensus threshold tx in chain33",
		Run:   setConsensus,
	}

	addSetConsensusFlags(cmd)
	return cmd
}

func addSetConsensusFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("power", "p", 0, "the power you want to set consensus need")
	_ = cmd.MarkFlagRequired("power")
}

func setConsensus(cmd *cobra.Command, args []string) {
	power, _ := cmd.Flags().GetInt64("power")

	params := &types3.MsgConsensusThreshold{
		ConsensusThreshold: power,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameSetConsensusThresholdAction)
}

func createTx(cmd *cobra.Command, payLoad []byte, action string) {
	title, _ := cmd.Flags().GetString("title")
	cfg := types.GetCliSysParam(title)
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	pm := &types2.CreateTxIn{
		Execer:     cfg.ExecName(types3.X2ethereumX),
		ActionName: action,
		Payload:    payLoad,
	}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.CreateTransaction", pm, &res)
	ctx.RunWithoutMarshal()
}
