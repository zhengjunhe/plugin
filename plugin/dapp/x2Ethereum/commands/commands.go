/*Package commands implement dapp client commands*/
package commands

import (
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	types2 "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/system/dapp/commands"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebcli/buildflags"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	types3 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

/*
 * 实现合约对应客户端
 */

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
		CreateTransferCmd(),
		queryCmd(),
		queryRelayerBalanceCmd(),
	)

	if buildflags.NodeAddr == "" {
		buildflags.NodeAddr = "http://127.0.0.1:7545"
	}
	cmd.PersistentFlags().String("node_addr", buildflags.NodeAddr, "eth node url")
	return cmd
}

// Eth2Chain33
func CreateRawEth2Chain33TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Eth2Chain33 && lock a erc20 and mint a token in chain33",
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

	cmd.Flags().Int64("decimal", 0, "the decimal of this token")
	_ = cmd.MarkFlagRequired("decimal")

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
	decimal, _ := cmd.Flags().GetInt64("decimal")

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
		Amount:                strconv.FormatFloat(types3.MultiplySpecifyTimes(amount, decimal), 'f', 4, 64),
		ClaimType:             claimtype,
		Decimals:              decimal,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameEth2Chain33Action)
}

// WithdrawEth
func CreateRawWithdrawEthTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraweth",
		Short: "withdraw a tx && burn erc20 back to chain33",
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
	decimal, _ := cmd.Flags().GetInt64("decimal")

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
		Amount:                strconv.FormatFloat(amount*1e8, 'f', 4, 64),
		ClaimType:             claimtype,
		Decimals:              decimal,
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
	cmd.Flags().StringP("contract", "q", "", "token contract address,nil for ETH")

	cmd.Flags().StringP("symbol", "t", "", "token symbol in chain33")
	_ = cmd.MarkFlagRequired("symbol")

	cmd.Flags().StringP("receiver", "r", "", "ethereum receiver address")
	_ = cmd.MarkFlagRequired("cExec")

	cmd.Flags().Float64P("amount", "a", float64(0), "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")

}

func burn(cmd *cobra.Command, args []string) {
	contract, _ := cmd.Flags().GetString("contract")
	csymbol, _ := cmd.Flags().GetString("symbol")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetFloat64("amount")
	nodeAddr, _ := cmd.Flags().GetString("node_addr")

	decimal, err := utils.GetDecimalsFromNode(contract, nodeAddr)
	if err != nil {
		fmt.Println("get decimal error")
		return
	}

	params := &types3.Chain33ToEth{
		TokenContract:    contract,
		EthereumReceiver: receiver,
		Amount:           types3.TrimZeroAndDot(strconv.FormatFloat(types3.MultiplySpecifyTimes(amount, decimal), 'f', 4, 64)),
		LocalCoinSymbol:  csymbol,
		Decimals:         decimal,
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
	contract, _ := cmd.Flags().GetString("contract")
	csymbol, _ := cmd.Flags().GetString("symbol")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetFloat64("amount")
	nodeAddr, _ := cmd.Flags().GetString("node_addr")

	decimal, err := utils.GetDecimalsFromNode(contract, nodeAddr)
	if err != nil {
		fmt.Println("get decimal error")
		return
	}

	params := &types3.Chain33ToEth{
		TokenContract:    contract,
		EthereumReceiver: receiver,
		Amount:           strconv.FormatFloat(amount*1e8, 'f', 4, 64),
		LocalCoinSymbol:  csymbol,
		Decimals:         decimal,
	}

	payLoad := types.MustPBToJSON(params)

	createTx(cmd, payLoad, types3.NameChain33ToEthAction)
}

// Transfer
func CreateTransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create a transfer tx in chain33",
		Run:   transfer,
	}

	addTransferFlags(cmd)
	return cmd
}

func transfer(cmd *cobra.Command, args []string) {
	commands.CreateAssetTransfer(cmd, args, types3.X2ethereumX)
}

func addTransferFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("to", "t", "", "receiver account address")
	_ = cmd.MarkFlagRequired("to")

	cmd.Flags().Float64P("amount", "a", 0, "transaction amount")
	_ = cmd.MarkFlagRequired("amount")

	cmd.Flags().StringP("note", "n", "", "transaction note info,optional")

	cmd.Flags().StringP("symbol", "s", "", "token symbol")
	_ = cmd.MarkFlagRequired("symbol")

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

func queryRelayerBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "query balance of x2ethereum",
		Run:   queryRelayerBalance,
	}

	cmd.Flags().StringP("token", "t", "", "token symbol")
	_ = cmd.MarkFlagRequired("token")

	cmd.Flags().StringP("address", "s", "", "the address you want to query")
	_ = cmd.MarkFlagRequired("address")
	return cmd
}

func queryRelayerBalance(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	token, _ := cmd.Flags().GetString("token")
	address, _ := cmd.Flags().GetString("address")

	get := &types3.QueryRelayerBalance{
		TokenSymbol: token,
		Address:     address,
	}

	payLoad, err := types.PBToJSON(get)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ErrPbToJson:"+err.Error())
		return
	}

	query := types2.Query4Jrpc{
		Execer:   types3.X2ethereumX,
		FuncName: types3.FuncQueryRelayerBalance,
		Payload:  payLoad,
	}

	channel := &types3.ReceiptQueryRelayerBalance{}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", query, channel)
	ctx.Run()
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