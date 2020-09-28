package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

// JvmCmd jvm command
func JvmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jvm",
		Short: "java contracts operation",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		jvmCheckContractNameCmd(),
		jvmCreateContractCmd(),
		jvmUpdateContractCmd(),
		jvmCallContractCmd(),
		jvmQueryContractCmd(),
		jvmDebugCmd(),
	)

	return cmd
}

// 创建jvm合约
func jvmCreateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new jvm contract",
		Run:   jvmCreateContract,
	}
	jvmAddCreateContractFlags(cmd)
	return cmd
}

func jvmAddCreateContractFlags(cmd *cobra.Command) {
	jvmAddCommonFlags(cmd)

	cmd.Flags().StringP("contract", "x", "", "contract name same with the code and abi file")
	_ = cmd.MarkFlagRequired("contract")

	cmd.Flags().StringP("path", "d", "", "path where stores jvm code and abi")
	_ = cmd.MarkFlagRequired("path")
}

func jvmCreateContract(cmd *cobra.Command, args []string) {
	contractName, _ := cmd.Flags().GetString("contract")
	path, _ := cmd.Flags().GetString("path")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	nameReg, _ := regexp.Compile(jvmTypes.NameRegExp)
	if !nameReg.MatchString(contractName) {
		_, _ = fmt.Fprintln(os.Stderr, "Wrong jvm contract name format, which should be a-z and 0-9 ")
		return
	}

	if len(contractName) > 16 || len(contractName) < 4 {
		_, _ = fmt.Fprintln(os.Stderr, "jvm contract name's length should be within range [4-16]")
		return
	}

	feeInt64 := uint64(fee*1e4) * 1e4
	codePath := path + "/" + contractName + ".jar"
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "read code error ", err)
		return
	}

	var createReq = jvmTypes.CreateContractReq{
		Name: contractName,
		Note: note,
		Code: code,
		Fee:  int64(feeInt64),
	}
	var createResp = types.ReplyString{}
	query := sendQuery4jvm(rpcLaddr, jvmTypes.CreateJvmContractStr, &createReq, &createResp)
	if query {
		fmt.Println(createResp.Data)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "get create to transaction error")
		return
	}
}

// 更新jvm合约
func jvmUpdateContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a old jvm contract",
		Run:   jvmUpdateContract,
	}
	jvmAddUpdateContractFlags(cmd)
	return cmd
}

func jvmAddUpdateContractFlags(cmd *cobra.Command) {
	jvmAddCommonFlags(cmd)

	cmd.Flags().StringP("contract", "x", "", "contract name same with the code and abi file")
	_ = cmd.MarkFlagRequired("contract")

	cmd.Flags().StringP("path", "d", "", "path where stores jvm code and abi")
	_ = cmd.MarkFlagRequired("path")
}

func jvmUpdateContract(cmd *cobra.Command, args []string) {
	contractName, _ := cmd.Flags().GetString("contract")
	path, _ := cmd.Flags().GetString("path")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	//paraName, _ := cmd.Flags().GetString("paraName")

	nameReg, _ := regexp.Compile(jvmTypes.NameRegExp)
	if !nameReg.MatchString(contractName) {
		_, _ = fmt.Fprintln(os.Stderr, "Wrong jvm contract name format, which should be a-z and 0-9 ")
		return
	}

	if len(contractName) > 16 || len(contractName) < 4 {
		_, _ = fmt.Fprintln(os.Stderr, "jvm contract name's length should be within range [4-16]")
		return
	}

	feeInt64 := uint64(fee*1e4) * 1e4

	codePath := path + "/" + contractName + ".class"
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "read code error ", err)
		return
	}

	var createReq = jvmTypes.UpdateContractReq{
		Name: contractName,
		Note: note,
		Code: code,
		Fee:  int64(feeInt64),
	}
	var createResp = types.ReplyString{}
	ok := sendQuery4jvm(rpcLaddr, jvmTypes.UpdateJvmContractStr, &createReq, &createResp)
	if ok {
		fmt.Println(createResp.Data)
		return
	}
	_, _ = fmt.Fprintln(os.Stderr, "get update to transaction error")
}

//运行jvm合约的查询请求
func jvmQueryContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "query the jvm contract",
		Run:   jvmQueryContract,
	}
	jvmAddQueryContractFlags(cmd)
	return cmd
}

func jvmQueryContract(cmd *cobra.Command, args []string) {
	contractName, _ := cmd.Flags().GetString("exec")
	paraOneStr, _ := cmd.Flags().GetString("para")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var paraParsed []string
	paraOneStr = strings.TrimSpace(paraOneStr)
	paraParsed = strings.Split(paraOneStr, " ")

	queryReq := jvmTypes.JVMQueryReq{
		Contract: contractName,
		Para:     paraParsed,
	}

	var jvmQueryResponse jvmTypes.JVMQueryResponse
	query := sendQuery4jvm(rpcLaddr, jvmTypes.JvmGetContractTable, &queryReq, &jvmQueryResponse)
	if !query {
		_, _ = fmt.Fprintln(os.Stderr, "get jvm query error")
		return
	}

	if !jvmQueryResponse.Success {
		fmt.Println("Exception occurred")
		return
	}

	for _, info := range jvmQueryResponse.Result {
		fmt.Println(info)
	}
	return
}

func jvmAddQueryContractFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("exec", "e", "", "jvm contract name")
	_ = cmd.MarkFlagRequired("exec")

	cmd.Flags().StringP("para", "r", "", "multiple parameter splitted by space")
	_ =  cmd.MarkFlagRequired("para")
}

// 调用jvm合约
func jvmCallContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call",
		Short: "Call the jvm contract",
		Run:   jvmCallContract,
	}
	jvmAddCallContractFlags(cmd)
	return cmd
}

func jvmCallContract(cmd *cobra.Command, args []string) {
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	contractName, _ := cmd.Flags().GetString("exec")
	actionName, _ := cmd.Flags().GetString("action")
	paraOneStr, _ := cmd.Flags().GetString("para")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	feeInt64 := uint64(fee*1e4) * 1e4

	var paraParsed []string
	paraOneStr = strings.TrimSpace(paraOneStr)
	paraParsed = strings.Split(paraOneStr, " ")

	var para  = []string{actionName}
	para = append(para, paraParsed...)

	var createReq = jvmTypes.CallJvmContract{
		Name:       contractName,
		Note:       note,
		ActionData: para,
		Fee:        int64(feeInt64),
	}
	var createResp = types.ReplyString{}

	query := sendQuery4jvm(rpcLaddr, jvmTypes.CallJvmContractStr, &createReq, &createResp)
	if query {
		fmt.Println(createResp.Data)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "get call jvm to transaction error")
		return
	}
}

func jvmAddCallContractFlags(cmd *cobra.Command) {
	jvmAddCommonFlags(cmd)
	cmd.Flags().StringP("exec", "e", "", "jvm contract name,like user.jvm.xxx")
	_ = cmd.MarkFlagRequired("exec")

	cmd.Flags().StringP("action", "x", "", "contract tx name")
	_ = cmd.MarkFlagRequired("action")

	cmd.Flags().StringP("para", "r", "", "multiple contract execution parameter splitted by space(optional)")
}

func jvmAddCommonFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")

	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")
}

// 检查地址是否为Jvm合约
func jvmCheckContractNameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if jvm contract used has been used already",
		Run:   jvmCheckContractAddr,
	}
	jvmAddCheckContractAddrFlags(cmd)
	return cmd
}

func jvmAddCheckContractAddrFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("exec", "e", "", "jvm contract name, like user.jvm.xxxxx(a-z0-9, within length [4-16])")
	_ = cmd.MarkFlagRequired("exec")
}

func jvmCheckContractAddr(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("exec")
	if bytes.Contains([]byte(name), []byte(jvmTypes.UserJvmX)) {
		name = name[len(jvmTypes.UserJvmX):]
	}

	match, _ := regexp.MatchString(jvmTypes.NameRegExp, name)
	if !match {
		_, _ = fmt.Fprintln(os.Stderr, "Wrong jvm contract name format, which should be a-z and 0-9 ")
		return
	}

	var checkAddrReq = jvmTypes.CheckJVMContractNameReq{JvmContractName: name}
	var checkAddrResp jvmTypes.CheckJVMAddrResp
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	query := sendQuery4jvm(rpcLaddr, jvmTypes.CheckNameExistsFunc, &checkAddrReq, &checkAddrResp)
	if query {
		_, _ = fmt.Fprintln(os.Stdout, checkAddrResp.ExistAlready)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "error")
	}
}

// 查询或设置Jvm调试开关
func jvmDebugCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Query or set external debug status",
	}
	cmd.AddCommand(
		jvmDebugQueryCmd(),
		jvmDebugSetCmd(),
		jvmDebugClearCmd())

	return cmd
}

func jvmDebugQueryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Query external debug status",
		Run:   jvmDebugQuery,
	}
}
func jvmDebugSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set",
		Short: "Set external debug to ON",
		Run:   jvmDebugSet,
	}
}
func jvmDebugClearCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Set external debug to OFF",
		Run:   jvmDebugClear,
	}
}

func jvmDebugQuery(cmd *cobra.Command, args []string) {
	jvmDebugRPC(cmd, 0)
}

func jvmDebugSet(cmd *cobra.Command, args []string) {
	jvmDebugRPC(cmd, 1)
}

func jvmDebugClear(cmd *cobra.Command, args []string) {
	jvmDebugRPC(cmd, -1)
}
func jvmDebugRPC(cmd *cobra.Command, flag int32) {
	var debugReq = jvmTypes.JVMDebugReq{Optype: flag}
	var debugResp jvmTypes.JVMDebugResp
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	query := sendQuery4jvm(rpcLaddr, "jvmDebug", &debugReq, &debugResp)

	if query {
		_ = proto.MarshalText(os.Stdout, &debugResp)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "error")
	}
}

func sendQuery4jvm(rpcAddr, funcName string, request types.Message, result proto.Message) bool {
	params := rpctypes.Query4Jrpc{
		Execer:   jvmTypes.JvmX,
		FuncName: funcName,
		Payload:  types.MustPBToJSON(request),
	}

	jsonrpc, err := jsonclient.NewJSONClient(rpcAddr)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return false
	}

	err = jsonrpc.Call("Chain33.Query", params, result)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return false
	}
	return true
}
