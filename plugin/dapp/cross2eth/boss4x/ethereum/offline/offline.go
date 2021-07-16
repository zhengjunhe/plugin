package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type DeployContractRet struct {
	ContractAddr string
	TxHash       string
	ContractName string
}

func DeployOfflineContractsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "deploy the corresponding Ethereum contracts",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		CreateCmd(), //构造交易
		CreateWithFileCmd(),
		SignCmd(),    //签名交易
		SendTxsCmd(), //发送交易
	)

	return cmd
}

func SendTxsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send", //first step
		Short: "send signed raw tx",
		Run:   sendTxs,
	}
	sendTxsFlags(cmd)
	return cmd
}

func sendTxsFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file", "f", "deploysigntxs.txt", "signed tx, default: deploysigntxs.txt")
	//_ = cmd.MarkFlagRequired("file")
}

func sendTxs(cmd *cobra.Command, _ []string) {
	url, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	filePath, _ := cmd.Flags().GetString("file")
	//解析文件数据
	var rdata = make([]*DeployInfo, 0)
	err := paraseFile(filePath, &rdata)
	if err != nil {
		fmt.Println("paraseFile,err", err.Error())
		return
	}
	var respData = make([]*DeployContractRet, 0)
	for _, deployInfo := range rdata {
		tx := new(types.Transaction)
		err = tx.UnmarshalBinary(common.FromHex(deployInfo.RawTx))
		if err != nil {
			panic(err)
		}
		client, err := ethclient.Dial(url)
		if err != nil {
			panic(err)
		}
		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			fmt.Println("err:", err)
			panic(err)
		}
		ret := &DeployContractRet{ContractAddr: deployInfo.ContractorAddr.String(), TxHash: tx.Hash().String(), ContractName: deployInfo.Name}
		respData = append(respData, ret)
		if !checkTxStatus(client, tx.Hash().String(), deployInfo.Name) {
			fmt.Println("FATAL ERROR! DEPLOY CONTRACTOR TERMINATION……:-(")
			break
		}
		//fmt.Println("Step:", i+1, "tx hash", tx.Hash().String(), "op:", deployInfo.Name)
	}

	data, err := json.MarshalIndent(respData, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
	//writeToFile("deploysendtx.txt", data)
	fmt.Println(string(data))
}

func checkTxStatus(client *ethclient.Client, txhash, txName string) bool {
	var checkticket = time.NewTicker(time.Second * 3)
	var timeout = time.NewTicker(time.Second * 300)
	for {
		select {
		case <-timeout.C:
			panic("Deploy timeout")
		case <-checkticket.C:
			receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txhash))
			if err == ethereum.NotFound {
				fmt.Println("\n No receipt received yet for "+txName, " tx and continue to wait")
				continue
			} else if err != nil {
				panic("failed due to" + err.Error())
			}

			//hasBridgeBank
			//	fmt.Println("receip:",receipt)

			if receipt.Status == types.ReceiptStatusSuccessful {
				return true
			}

			if receipt.Status == types.ReceiptStatusFailed {
				fmt.Println("tx status:", types.ReceiptStatusFailed)
				return false
			}
		}
	}

	return false

}
