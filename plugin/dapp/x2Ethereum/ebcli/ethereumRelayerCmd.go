package main

import (
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	ebTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
)

// EthereumRelayerCmd command func
func EthereumRelayerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ethereum",
		Short: "Ethereum relayer",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		ImportChain33PrivateKeyCmd(),
		ImportEthValidatorPrivateKeyCmd(),
		GenEthPrivateKeyCmd(),
		ShowValidatorsAddrCmd(),
		ShowChain33TxsHashCmd(),
		ShowEthereumTxsHashCmd(),
		ShowEthRelayerStatusCmd(),
		IsValidatorActiveCmd(),
		ShowOperatorCmd(),
		DeployContrctsCmd(),
		ShowTxReceiptCmd(),
		//////auxiliary///////
		CreateBridgeTokenCmd(),
		MakeNewProphecyClaimCmd(),
		ProcessProphecyClaimCmd(),
		GetBalanceCmd(),

	)

	return cmd
}

func ImportChain33PrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import_chain33privatekey",
		Short: "import chain33 private key to sign txs to be submitted to chain33",
		Run:   importChain33Privatekey,
	}
	addImportChain33PrivateKeyFlags(cmd)
	return cmd
}

func addImportChain33PrivateKeyFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "chain33 private key")
	cmd.MarkFlagRequired("key")
}

func importChain33Privatekey(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	privateKey, _ := cmd.Flags().GetString("key")
	params := privateKey

	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ImportChain33PrivateKey4EthRelayer", params, &res)
	ctx.Run()
}

func ImportEthValidatorPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import_ethprivatekey",
		Short: "import ethereum's validator private key ",
		Run:   importEthValidtorPrivatekey,
	}
	addImportPrivateKeyFlags(cmd)
	return cmd
}

func importEthValidtorPrivatekey(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	privateKey, _ := cmd.Flags().GetString("key")
	params := privateKey

	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ImportEthValidatorPrivateKey", params, &res)
	ctx.Run()
}

func GenEthPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_eth_key",
		Short: "create ethereum's private key to sign txs to be submitted to ethereum",
		Run:   generateEthereumPrivateKey,
	}
	return cmd
}

func generateEthereumPrivateKey(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var res ebTypes.Account4Show
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.GenerateEthereumPrivateKey", nil, &res)
	ctx.Run()
}

func ShowValidatorsAddrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show_validators",
		Short: "show me the validators including ethereum and chain33",
		Run:   showValidatorsAddr,
	}
	return cmd
}

func showValidatorsAddr(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	var res ebTypes.ValidatorAddr4EthRelayer
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowEthRelayerValidator", nil, &res)
	ctx.Run()
}

func ShowChain33TxsHashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show_chain33_tx",
		Short: "show me the chain33 tx hashes",
		Run:   showChain33Txs,
	}
	return cmd
}

func showChain33Txs(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var res ebTypes.Txhashes
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowEthRelayer2Chain33Txs", nil, &res)
	if _, err := ctx.RunResult(); nil != err {
		errInfo := err.Error()
		fmt.Println("errinfo:"+errInfo)
		return
	}
	for _, hash := range res.Txhash {
		fmt.Println(hash)
	}
}

func ShowEthereumTxsHashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show_eth_tx",
		Short: "show me the ethereum tx hashes",
		Run:   showEthTxs,
	}
	return cmd
}

func showEthTxs(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var res ebTypes.Txhashes
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowEthRelayer2Chain33Txs", nil, &res)
	if _, err := ctx.RunResult(); nil != err {
		errInfo := err.Error()
		fmt.Println("errinfo:"+errInfo)
		return
	}
	for _, hash := range res.Txhash {
		fmt.Println(hash)
	}
}

func ShowEthRelayerStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "show ethereum-relayer status",
		Run:   showEthRelayerStatus,
	}
	return cmd
}

func showEthRelayerStatus(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var res ebTypes.RelayerRunStatus
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowEthRelayerStatus", nil, &res)
	ctx.Run()
}

func IsValidatorActiveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "active",
		Short: "show whether the validator is active or not",
		Run:   IsValidatorActive,
	}
	IsValidatorActiveFlags(cmd)
	return cmd
}

func IsValidatorActiveFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("addr", "a", "", "validator address")
	_ = cmd.MarkFlagRequired("addr")
}

func IsValidatorActive(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("addr")

	params := addr
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.IsValidatorActive", params, &res)
	ctx.Run()
}

func ShowOperatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator",
		Short: "show me the operator",
		Run:   ShowOperator,
	}
	return cmd
}

func ShowOperator(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowOperator", nil, &res)
	ctx.Run()
}

func DeployContrctsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy the corresponding Ethereum contracts",
		Run:   DeployContrcts,
	}
	return cmd
}

func DeployContrcts(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.DeployContrcts", nil, &res)
	ctx.Run()
}

func ShowTxReceiptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receipt",
		Short: "show me the tx receipt for Ethereum",
		Run:   ShowTxReceipt,
	}
	ShowTxReceiptFlags(cmd)
	return cmd
}

func ShowTxReceiptFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("hash", "s", "", "tx hash")
	_ = cmd.MarkFlagRequired("hash")
}

func ShowTxReceipt(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	txhash, _ := cmd.Flags().GetString("hash")
	para := txhash
	var res ethTypes.Receipt
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ShowTxReceipt", para, &res)
	ctx.Run()
}

func CreateBridgeTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "create new token on Ethereum",
		Run:   CreateBridgeToken,
	}
	CreateBridgeTokenFlags(cmd)
	return cmd
}

func CreateBridgeTokenFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("symbol", "s", "", "token symbol")
	_ = cmd.MarkFlagRequired("symbol")
}

func CreateBridgeToken(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	token, _ := cmd.Flags().GetString("symbol")
	para := token
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.CreateBridgeToken", para, &res)
	ctx.Run()
}

func MakeNewProphecyClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prophecy",
		Short: "Make New Prophecy Claim",
		Run:   MakeNewProphecyClaim,
	}
	MakeNewProphecyClaimFlags(cmd)
	return cmd
}

func MakeNewProphecyClaimFlags(cmd *cobra.Command) {
	cmd.Flags().Uint32P("claim", "c", uint32(1), "claim type, 1 denote burn, and 2 denotes lock")
	_ = cmd.MarkFlagRequired("claim")
	cmd.Flags().StringP("chain33Sender", "t", "", "Chain33Sender")
	_ = cmd.MarkFlagRequired("chain33Sender")
	cmd.Flags().StringP("address", "a", "", "token address")
	_ = cmd.MarkFlagRequired("address")
	cmd.Flags().StringP("symbol", "s", "", "token symbol")
	_ = cmd.MarkFlagRequired("symbol")
}

func MakeNewProphecyClaim(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	claimType, _ := cmd.Flags().GetUint32("claim")
	if claimType != uint32(1) && claimType != uint32(2) {
		fmt.Println("Wrong claim type")
		return
	}
	chain33Sender, _ := cmd.Flags().GetString("chain33Sender")
	tokenAddr, _ := cmd.Flags().GetString("address")
	symbol, _ := cmd.Flags().GetString("symbol")
	para := ebTypes.NewProphecyClaim{
		ClaimType:claimType,
		Chain33Sender:chain33Sender,
		TokenAddr:tokenAddr,
		Symbol:symbol,
	}
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.MakeNewProphecyClaim", para, &res)
	ctx.Run()
}


func ProcessProphecyClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process",
		Short: "process Prophecy Claim",
		Run:   ProcessProphecyClaim,
	}
	ProcessProphecyClaimFlags(cmd)
	return cmd
}

func ProcessProphecyClaimFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("prophecyID", "i", int64(0), "prophecy id to be processed")
	_ = cmd.MarkFlagRequired("prophecyID")
}

func ProcessProphecyClaim(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	prophecyID, _ := cmd.Flags().GetUint32("prophecyID")
	para := prophecyID
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.ProcessProphecyClaim", para, &res)
	ctx.Run()
}

func GetBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "get balance for addr of token",
		Run:   GetBalance,
	}
	GetBalanceFlags(cmd)
	return cmd
}

func GetBalanceFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("owner", "o", "", "owner address")
	_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringP("tokenAddr", "a", "", "token address")
	_ = cmd.MarkFlagRequired("tokenAddr")
}

func GetBalance(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	owner, _ := cmd.Flags().GetString("owner")
	tokenAddr, _ := cmd.Flags().GetString("tokenAddr")
	para := ebTypes.BalanceAddr{
		Owner:                owner,
		TokenAddr:            tokenAddr,
	}
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "RelayerManager.GetBalance", para, &res)
	ctx.Run()
}