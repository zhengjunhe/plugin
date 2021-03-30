package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/rpc/jsonclient"
	"github.com/33cn/chain33/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/33cn/chain33/common"
	rpctypes "github.com/33cn/chain33/rpc/types"
	commandtypes "github.com/33cn/chain33/system/dapp/commands/types"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

var PancakeFactoryBinFile = "./ci/evm/PancakeFactory.bin"
var PancakeFactoryAbiFile = "./ci/evm/PancakeFactory.abi"
var WETH9BinFile = "./ci/evm/WETH9.bin"
var WETH9AbiFile = "./ci/evm/WETH9.abi"
var PancakeRouterBinFile = "./ci/evm/PancakeRouter.bin"
var PancakeRouterAbiFile = "./ci/evm/PancakeRouter.abi"

func DeployPancake(cmd *cobra.Command) error {
	parameter, _ := cmd.Flags().GetString("parameter")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	{
		txhex, err := deployContract(cmd, PancakeFactoryBinFile, PancakeFactoryAbiFile, parameter, "PancakeFactory")
		if err != nil {
			return errors.New(err.Error())
		}

		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployPancakeFactory timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhex, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for DeployPancakeFactory tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("DeployPancakeFactory failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy pancakeFactory with address =", txhex, "\\n")
				goto deployWeth9
			}
		}
	}

deployWeth9:
	{
		txhex, err := deployContract(cmd, WETH9BinFile, WETH9AbiFile, "", "Weth9")
		if err != nil {
			return errors.New(err.Error())
		}

		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy Weth9 timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhex, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy Weth9 tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy Weth9 failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy Weth9 with address =", txhex, "\\n")
				goto deployPancakeRouter
			}
		}
	}

deployPancakeRouter:
	{
		txhex, err := deployContract(cmd, PancakeRouterBinFile, PancakeRouterAbiFile, "", "PancakeRouter")
		if err != nil {
			return errors.New(err.Error())
		}

		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Deploy PancakeRouter timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(txhex, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy PancakeRouter tx and continue to wait")
					continue
				} else if data != "2" {
					return errors.New("Deploy PancakeRouter failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deploy PancakeRouter with address =", txhex, "\\n")
				return nil
			}
		}
	}

	return nil
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
	if err != nil {
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
		result.Txs = append(result.Txs, &td)
	}
	return receipt.Ty, nil
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
