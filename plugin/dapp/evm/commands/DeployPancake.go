package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/rpc/jsonclient"
	"github.com/33cn/chain33/types"
	"github.com/spf13/cobra"

	"github.com/33cn/chain33/common"
	rpctypes "github.com/33cn/chain33/rpc/types"
	commandtypes "github.com/33cn/chain33/system/dapp/commands/types"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"

	//evmcommon "github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
	ethereumcommon "github.com/ethereum/go-ethereum/common"
)

var ERC20BinFile = "./ci/evm/ERC20.bin"
var ERC20AbiFile = "./ci/evm/ERC20.abi"
var PancakeFactoryBinFile = "./ci/evm/PancakeFactory.bin"
var PancakeFactoryAbiFile = "./ci/evm/PancakeFactory.abi"
var WETH9BinFile = "./ci/evm/WETH9.bin"
var WETH9AbiFile = "./ci/evm/WETH9.abi"
var PancakeRouterBinFile = "./ci/evm/PancakeRouter.bin"
var PancakeRouterAbiFile = "./ci/evm/PancakeRouter.abi"
var MulticallBinFile = "./ci/evm/Multicall.bin"
var MulticallAbiFile = "./ci/evm/Multicall.abi"

func DeployMulticall(cmd *cobra.Command) error {
	caller, _ := cmd.Flags().GetString("caller")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	txMulticall, err := deployContract(cmd, MulticallBinFile, MulticallAbiFile, "", "Multicall")
	if err != nil {
		return errors.New(err.Error())
	}

	{
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy ERC20 timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txMulticall, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy Multicall tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy Multicall failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy Multicall with address =", getContractAddr(caller, txMulticall), "\n")
				return nil

			}
		}
	}
}

func DeployPancake(cmd *cobra.Command) error {
	caller, _ := cmd.Flags().GetString("caller")
	parameter, _ := cmd.Flags().GetString("parameter")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	txhexERC20, err := deployContract(cmd, ERC20BinFile, ERC20AbiFile, "ycc, ycc, 3300000000, "+caller, "ERC20")
	if err != nil {
		return errors.New(err.Error())
	}

	{
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy ERC20 timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhexERC20, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy ERC20 tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy ERC20 failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy ERC20 with address =", getContractAddr(caller, txhexERC20), "\n")
				goto deployPancakeFactory

			}
		}
	}

deployPancakeFactory:
	txhexPancakeFactory, err := deployContract(cmd, PancakeFactoryBinFile, PancakeFactoryAbiFile, parameter, "PancakeFactory")
	if err != nil {
		return errors.New(err.Error())
	}

	{
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy PancakeFactory timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhexPancakeFactory, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy PancakeFactory tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy PancakeFactory failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy pancakeFactory with address =", getContractAddr(caller, txhexPancakeFactory), "\n")
				goto deployWeth9
			}
		}
	}

deployWeth9:
	txhexWeth9, err := deployContract(cmd, WETH9BinFile, WETH9AbiFile, "", "Weth9")
	if err != nil {
		return errors.New(err.Error())
	}

	{
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy Weth9 timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhexWeth9, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy Weth9 tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy Weth9 failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy Weth9 with address =", getContractAddr(caller, txhexWeth9), "\n")
				goto deployPancakeRouter
			}
		}
	}

deployPancakeRouter:
	param := getContractAddr(caller, txhexPancakeFactory) + "," + getContractAddr(caller, txhexWeth9)
	txhexPancakeRouter, err := deployContract(cmd, PancakeRouterBinFile, PancakeRouterAbiFile, param, "PancakeRouter")
	if err != nil {
		return errors.New(err.Error())
	}

	{
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy PancakeRouter timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhexPancakeRouter, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy PancakeRouter tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy PancakeRouter failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy PancakeRouter with address =", getContractAddr(caller, txhexPancakeRouter), "\n")
				return nil
			}
		}
	}

	return nil
}

func getContractAddr(caller, txhex string) string {
	return address.GetExecAddress(caller + ethereumcommon.Bytes2Hex(common.HexToHash(txhex).Bytes())).String()
}

func deployContract(cmd *cobra.Command, binFile, abiFile, parameter, contractName string) (string, error) {
	title, _ := cmd.Flags().GetString("title")
	cfg := types.GetCliSysParam(title)

	caller, _ := cmd.Flags().GetString("caller")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	paraName, _ := cmd.Flags().GetString("paraName")
	feeInt64 := uint64(fee*1e4) * 1e4

	code, err := readContractFile(binFile)
	if err != nil {
		return "", errors.New(contractName + " read file error " + binFile + err.Error())
	}
	abi, err := readContractFile(abiFile)
	if err != nil {
		return "", errors.New(contractName + " read file error " + abiFile + err.Error())
	}

	var action evmtypes.EVMContractAction
	bCode, err := common.FromHex(code)
	if err != nil {
		return "", errors.New(contractName + " parse evm code error " + err.Error())
	}
	action = evmtypes.EVMContractAction{Amount: 0, Code: bCode, GasLimit: 0, GasPrice: 0, Note: note, Alias: "PancakeFactory", Abi: abi}
	if parameter != "" {
		constructorPara := "constructor(" + parameter + ")"
		packData, err := evmAbi.PackContructorPara(constructorPara, abi)
		if err != nil {
			return "", errors.New(contractName + " Pack Contructor Para error:" + err.Error())
		}
		action.Code = append(action.Code, packData...)
	}
	data, err := createEvmTx(cfg, &action, cfg.ExecName(paraName+"evm"), caller, address.ExecAddress(cfg.ExecName(paraName+"evm")), expire, rpcLaddr, feeInt64)
	if err != nil {
		return "", errors.New(contractName + " create contract error:" + err.Error())
	}

	txhex, err := sendTransactionRpc(data, rpcLaddr)
	if err != nil {
		return "", errors.New(contractName + " send transaction error:" + err.Error())
	}
	fmt.Println("Deploy", contractName, "tx hash:", txhex)

	return txhex, nil
}

func getTxByHashesRpc(txhex, rpcLaddr string) (string, error) {
	hashesArr := strings.Split(txhex, " ")
	params2 := rpctypes.ReqHashes{
		Hashes: hashesArr,
	}

	var res rpctypes.TransactionDetails
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.GetTxByHashes", params2, &res)
	ctx.SetResultCb(queryTxsByHashesRes)
	result, err := ctx.RunResult()
	if err != nil || result == nil {
		return "", err
	}
	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func queryTxsByHashesRes(arg interface{}) (interface{}, error) {
	var result commandtypes.TxDetailsResult
	var receipt *rpctypes.ReceiptDataResult
	for _, v := range arg.(*rpctypes.TransactionDetails).Txs {
		if v == nil {
			result.Txs = append(result.Txs, nil)
			continue
		}
		amountResult := strconv.FormatFloat(float64(v.Amount)/float64(types.Coin), 'f', 4, 64)
		td := commandtypes.TxDetailResult{
			Tx:         commandtypes.DecodeTransaction(v.Tx),
			Receipt:    v.Receipt,
			Proofs:     v.Proofs,
			Height:     v.Height,
			Index:      v.Index,
			Blocktime:  v.Blocktime,
			Amount:     amountResult,
			Fromaddr:   v.Fromaddr,
			ActionName: v.ActionName,
			Assets:     v.Assets,
		}
		receipt = v.Receipt
		if nil != receipt {
			return receipt.Ty, nil
		}
		result.Txs = append(result.Txs, &td)
	}
	return nil, nil
}

func sendTransactionRpc(data, rpcLaddr string) (string, error) {
	params := rpctypes.RawParm{
		Data: data,
	}
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.SendTransaction", params, nil)
	var txhex string
	rpc, err := jsonclient.NewJSONClient(ctx.Addr)
	if err != nil {
		return "", err
	}

	err = rpc.Call(ctx.Method, ctx.Params, &txhex)
	if err != nil {
		return "", err
	}

	return txhex, nil
}

func readContractFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return "", err
	}

	fileContent, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}
