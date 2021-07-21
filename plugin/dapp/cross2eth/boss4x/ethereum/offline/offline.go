package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	erc20 "github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/erc20/generated"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	eoff "github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum/offline"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
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
		DeployERC20Cmd(),
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
	cmd.Flags().StringP("file", "f", "", "signed tx file")
	_ = cmd.MarkFlagRequired("file")
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
	}

	data, err := json.MarshalIndent(respData, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}
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

func DeployERC20Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_erc20",
		Short: "create ERC20 contracts and sign",
		Run:   DeployERC20,
	}
	DeployERC20Flags(cmd)
	return cmd
}

func DeployERC20Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "private key ")
	_ = cmd.MarkFlagRequired("key")
	cmd.Flags().StringP("owner", "o", "", "owner address")
	_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringP("symbol", "s", "", "erc20 symbol")
	_ = cmd.MarkFlagRequired("symbol")
	cmd.Flags().StringP("amount", "m", "0", "amount")
	_ = cmd.MarkFlagRequired("amount")
}

func DeployERC20(cmd *cobra.Command, _ []string) {
	url, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	key, _ := cmd.Flags().GetString("key")
	owner, _ := cmd.Flags().GetString("owner")
	symbol, _ := cmd.Flags().GetString("symbol")
	amount, _ := cmd.Flags().GetString("amount")
	bnAmount := big.NewInt(1)
	bnAmount, _ = bnAmount.SetString(utils.TrimZeroAndDot(amount), 10)

	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(key))
	if err != nil {
		fmt.Println("ToECDSA error", err.Error())
		return
	}

	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)
	client, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("ethclient Dial error", err.Error())
		return
	}

	ctx := context.Background()
	price, err := client.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("SuggestGasPrice error", err.Error())
		return
	}

	startNonce, err := client.PendingNonceAt(ctx, deployerAddr)
	if nil != err {
		fmt.Println("PendingNonceAt error", err.Error())
		return
	}

	var infos []*DeployInfo

	parsed, err := abi.JSON(strings.NewReader(erc20.ERC20ABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(erc20.ERC20Bin)
	Erc20OwnerAddr := common.HexToAddress(owner)
	packdata, err := parsed.Pack("", symbol, symbol, bnAmount, Erc20OwnerAddr)
	if err != nil {
		panic(err)
	}
	Erc20Addr := crypto.CreateAddress(deployerAddr, startNonce)
	deployInfo := DeployInfo{
		PackData:       append(bin, packdata...),
		ContractorAddr: Erc20Addr,
		Name:           "Erc20: " + symbol,
		Nonce:          startNonce,
		To:             nil,
	}
	infos = append(infos, &deployInfo)
	//预估gas,批量构造交易
	for i, info := range infos {
		var msg ethereum.CallMsg
		msg.From = deployerAddr
		msg.To = info.To
		msg.Value = big.NewInt(0)
		msg.Data = info.PackData
		//估算gas
		gasLimit, err := client.EstimateGas(ctx, msg)
		if err != nil {
			fmt.Println("EstimateGas error", err.Error())
			return
		}
		if gasLimit < 100*10000 {
			gasLimit = 100 * 10000
		}
		ntx := types.NewTx(&types.LegacyTx{
			Nonce:    info.Nonce,
			Gas:      gasLimit,
			GasPrice: price,
			Data:     info.PackData,
			To:       info.To,
		})

		txBytes, err := ntx.MarshalBinary()
		if err != nil {
			fmt.Println("MarshalBinary error", err.Error())
			return
		}
		infos[i].RawTx = common.Bytes2Hex(txBytes)
		infos[i].Gas = gasLimit

		var tx types.Transaction
		err = tx.UnmarshalBinary(common.FromHex(info.RawTx))
		if err != nil {
			panic(err)
		}
		signedTx, txHash, err := eoff.SignTx(deployPrivateKey, &tx)
		if err != nil {
			panic(err)
		}
		infos[i].RawTx = signedTx
		infos[i].TxHash = txHash
	}

	jbytes, err := json.MarshalIndent(&infos, "", "\t")
	if err != nil {
		fmt.Println("MarshalIndent error", err.Error())
		return
	}

	fmt.Println(string(jbytes))
	fileName := fmt.Sprintf("deployErc20%s.txt", symbol)
	writeToFile(fileName, &infos)
}
