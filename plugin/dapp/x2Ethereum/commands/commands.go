/*Package commands implement dapp client commands*/
package commands

import (
	"encoding/json"
	"github.com/33cn/chain33/rpc/jsonclient"
	types2 "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	types3 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/spf13/cobra"
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
		CreateRawEthBridgeClaimTxCmd(),
		CreateRawBurnTxCmd(),
		CreateRawLockTxCmd(),
		CreateRawLogInTxCmd(),
		CreateRawLogOutTxCmd(),
		CreateRawSetConsensusTxCmd(),
	)
	return cmd
}

// ethBridgeClaim
func CreateRawEthBridgeClaimTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a EthBridgeClaim",
		Run:   ethBridgeClaim,
	}

	addEthBridgeClaimFlags(cmd)
	return cmd
}

func addEthBridgeClaimFlags(cmd *cobra.Command) {
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

	cmd.Flags().Uint64("amount", 0, "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")

	cmd.Flags().Int64("claimtype", 0, "the type of this claim")
	_ = cmd.MarkFlagRequired("claimtype")

	cmd.Flags().StringP("esymbol", "h", "", "the symbol of ethereum side")
	_ = cmd.MarkFlagRequired("esymbol")
}

func ethBridgeClaim(cmd *cobra.Command, args []string) {
	ethid, _ := cmd.Flags().GetInt64("ethid")
	bcontract, _ := cmd.Flags().GetString("bcontract")
	nonce, _ := cmd.Flags().GetInt64("nonce")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	tcontract, _ := cmd.Flags().GetString("tcontract")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	validator, _ := cmd.Flags().GetString("validator")
	amount, _ := cmd.Flags().GetUint64("amount")
	claimtype, _ := cmd.Flags().GetInt64("claimtype")
	esymbol, _ := cmd.Flags().GetString("esymbol")

	params := &types3.EthBridgeClaim{
		EthereumChainID:       ethid,
		BridgeContractAddress: bcontract,
		Nonce:                 nonce,
		LocalCoinSymbol:       csymbol,
		LocalCoinExec:         cexec,
		TokenContractAddress:  tcontract,
		EthereumSender:        sender,
		Chain33Receiver:       receiver,
		ValidatorAddress:      validator,
		Amount:                amount,
		ClaimType:             claimtype,
		EthSymbol:             esymbol,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "EthBridgeClaim")
}

// Burn
func CreateRawBurnTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn",
		Short: "Create a burn tx in chain33",
		Run:   burn,
	}

	addBurnFlags(cmd)
	return cmd
}

func addBurnFlags(cmd *cobra.Command) {
	cmd.Flags().Int64("ethid", 0, "the ethereum chain ID which send asset to chain33")
	_ = cmd.MarkFlagRequired("ethid")

	cmd.Flags().StringP("contract", "b", "", "token contract address")
	_ = cmd.MarkFlagRequired("contract")

	cmd.Flags().StringP("csymbol", "t", "", "token symbol in chain33")
	_ = cmd.MarkFlagRequired("csymbol")

	cmd.Flags().StringP("cexec", "e", "", "chain execer in chain33")
	_ = cmd.MarkFlagRequired("cexec")

	cmd.Flags().StringP("sender", "s", "", "chain33 sender address")
	_ = cmd.MarkFlagRequired("sender")

	cmd.Flags().StringP("receiver", "r", "", "ethereum receiver address")
	_ = cmd.MarkFlagRequired("cExec")

	cmd.Flags().Uint64("amount", 0, "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")
}

func burn(cmd *cobra.Command, args []string) {
	ethid, _ := cmd.Flags().GetInt64("ethid")
	contract, _ := cmd.Flags().GetString("contract")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetUint64("amount")

	params := &types3.MsgBurn{
		EthereumChainID:  ethid,
		TokenContract:    contract,
		Chain33Sender:    sender,
		EthereumReceiver: receiver,
		Amount:           amount,
		LocalCoinSymbol:  csymbol,
		LocalCoinExec:    cexec,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "MsgBurn")
}

// Lock
func CreateRawLockTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Create a lock tx in chain33",
		Run:   lock,
	}

	addLockFlags(cmd)
	return cmd
}

func addLockFlags(cmd *cobra.Command) {
	cmd.Flags().Int64("ethid", 0, "the ethereum chain ID which send asset to chain33")
	_ = cmd.MarkFlagRequired("ethid")

	cmd.Flags().StringP("contract", "b", "", "token contract address")
	_ = cmd.MarkFlagRequired("contract")

	cmd.Flags().StringP("csymbol", "t", "", "token symbol in chain33")
	_ = cmd.MarkFlagRequired("csymbol")

	cmd.Flags().StringP("cexec", "e", "", "chain execer in chain33")
	_ = cmd.MarkFlagRequired("cexec")

	cmd.Flags().StringP("sender", "s", "", "chain33 sender address")
	_ = cmd.MarkFlagRequired("sender")

	cmd.Flags().StringP("receiver", "r", "", "ethereum receiver address")
	_ = cmd.MarkFlagRequired("cExec")

	cmd.Flags().Uint64("amount", 0, "the amount of this contract want to lock")
	_ = cmd.MarkFlagRequired("amount")
}

func lock(cmd *cobra.Command, args []string) {
	ethid, _ := cmd.Flags().GetInt64("ethid")
	contract, _ := cmd.Flags().GetString("contract")
	csymbol, _ := cmd.Flags().GetString("csymbol")
	cexec, _ := cmd.Flags().GetString("cexec")
	sender, _ := cmd.Flags().GetString("sender")
	receiver, _ := cmd.Flags().GetString("receiver")
	amount, _ := cmd.Flags().GetUint64("amount")

	params := &types3.MsgLock{
		EthereumChainID:  ethid,
		TokenContract:    contract,
		Chain33Sender:    sender,
		EthereumReceiver: receiver,
		Amount:           amount,
		LocalCoinSymbol:  csymbol,
		LocalCoinExec:    cexec,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "MsgLock")
}

// MsgLogInValidator
func CreateRawLogInTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Create a LogIn tx in chain33",
		Run:   logIn,
	}

	addLogInFlags(cmd)
	return cmd
}

func addLogInFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", "", "the address you want to log in ")
	_ = cmd.MarkFlagRequired("address")

	cmd.Flags().Float64P("power", "p", 0, "validator power set")
	_ = cmd.MarkFlagRequired("power")
}

func logIn(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")
	power, _ := cmd.Flags().GetFloat64("power")

	params := &types3.MsgValidator{
		Address: address,
		Power:   power,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "MsgLogInValidator")
}

// MsgLogOutValidator
func CreateRawLogOutTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Create a LogOut tx in chain33",
		Run:   logOut,
	}

	addLogOutFlags(cmd)
	return cmd
}

func addLogOutFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", "", "the address you want to log out ")
	_ = cmd.MarkFlagRequired("address")

	cmd.Flags().Float64P("power", "p", 0, "validator power set")
	_ = cmd.MarkFlagRequired("power")
}

func logOut(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")
	power, _ := cmd.Flags().GetFloat64("power")

	params := &types3.MsgValidator{
		Address: address,
		Power:   power,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "MsgLogOutValidator")
}

// MsgSetConsensusNeeded
func CreateRawSetConsensusTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setconsensus",
		Short: "Create a set consensus Need tx in chain33",
		Run:   setConsensus,
	}

	addSetConsensusFlags(cmd)
	return cmd
}

func addSetConsensusFlags(cmd *cobra.Command) {
	cmd.Flags().Float64P("power", "p", 0, "the power you want to set consensus need")
	_ = cmd.MarkFlagRequired("power")
}

func setConsensus(cmd *cobra.Command, args []string) {
	power, _ := cmd.Flags().GetFloat64("power")

	params := &types3.MsgSetConsensusNeeded{
		Power: power,
	}

	payLoad, err := json.Marshal(params)
	if err != nil {
		return
	}

	createTx(cmd, payLoad, "MsgSetConsensusNeeded")
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
