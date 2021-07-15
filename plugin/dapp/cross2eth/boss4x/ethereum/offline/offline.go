package offline

import (
	"context"
	"encoding/json"
	"fmt"
	eoff "github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum/offline"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeployContractRet struct {
	ContractAddr string
	TxHash       string
	ContractName string
}

func Boss4xEthOfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offline",
		Short: "create and sign offline tx to deploy and set cross contracts to ethereum",
	}
	cmd.AddCommand(
		getNonceCmd(),
		createAndSignTxsCmd(),
		sendTxsCmd(),
	)
	return cmd
}

func getNonceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNonce", //first step
		Short: "query gasPrice, nonce from the spec address",
		Run:   getNonce,
	}
	getNonceFlags(cmd)
	return cmd
}

func getNonceFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", "", "query address")
	_ = cmd.MarkFlagRequired("address")
}

func getNonce(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	addr, _ := cmd.Flags().GetString("address")

	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	price, err := client.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(addr))
	if nil != err {
		fmt.Println("err:", err)
	}

	fmt.Println("    Address: ", addr)
	fmt.Println("    GasPrice: ", price.Uint64())
	fmt.Println("    Nonce: ", nonce)
}

func createAndSignTxsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_sign", //first step
		Short: "create contract and sign",
		Run:   createAndSignTxs,
	}
	createAndSignTxsFlag(cmd)
	return cmd
}

func createAndSignTxsFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("validatorsAddrs", "v", "", "validatorsAddrs, as: 'addr,addr,addr,addr'")
	_ = cmd.MarkFlagRequired("validatorsAddrs")
	cmd.Flags().StringP("initPowers", "p", "", "initPowers, as: '25,25,25,25'")
	_ = cmd.MarkFlagRequired("initPowers")
	cmd.Flags().StringP("key", "k", "", "the deployer private key")
	_ = cmd.MarkFlagRequired("key")
	cmd.Flags().Uint64P("nonce", "n", 0, "transaction count")
	cmd.Flags().Uint64P("gasprice", "g", 1000000000, "gas price") // 1Gwei=1e9wei
	cmd.Flags().Uint64P("gaslimit", "l", 21000, "gas limit")
}

func createAndSignTxs(cmd *cobra.Command, args []string) {
	deployerPrivateKey, _ := cmd.Flags().GetString("key")
	gasprice, _ := cmd.Flags().GetUint64("gasprice")
	gaslimit, _ := cmd.Flags().GetUint64("gaslimit")
	nonce, _ := cmd.Flags().GetUint64("nonce")
	validatorsAddrs, _ := cmd.Flags().GetString("validatorsAddrs")
	initPowers, _ := cmd.Flags().GetString("initPowers")

	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(deployerPrivateKey))
	if err != nil {
		panic(err)
	}
	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)

	validatorsAddrsArray := strings.Split(validatorsAddrs, ",")
	initPowersArray := strings.Split(initPowers, ",")

	if len(validatorsAddrsArray) != len(initPowersArray) {
		fmt.Println("input validatorsAddrs initPowers error!")
		return
	}

	var validators []common.Address
	var initpowers []*big.Int
	for _, v := range validatorsAddrsArray {
		validators = append(validators, common.HexToAddress(v))
	}

	for _, v := range initPowersArray {
		vint64, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		initpowers = append(initpowers, big.NewInt(vint64))
	}

	var signData = make([]*eoff.DeployContract, 0)
	valSet := signValSet(validators, initpowers, deployerAddr, nonce, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, valSet)
	//Step2 Sign chain33 bridge
	bridge := signChain33Bridge(deployerAddr, common.HexToAddress(valSet.ContractAddr), valSet.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, bridge)
	//step3 sign oracle
	oracle := signOracle(common.HexToAddress(valSet.ContractAddr), common.HexToAddress(bridge.ContractAddr), deployerAddr, bridge.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, oracle)
	//step4  sign bridgebank
	bank := signBridgeBank(common.HexToAddress(bridge.ContractAddr), common.HexToAddress(oracle.ContractAddr), deployerAddr, oracle.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, bank)
	//step5 sign SetBridgeBank
	setbank := signSetBridgeBank(common.HexToAddress(bank.ContractAddr), common.HexToAddress(bridge.ContractAddr), deployerAddr, bank.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, setbank)
	//step6 Sign setOracle
	setOracle := signsetOracle(common.HexToAddress(oracle.ContractAddr), common.HexToAddress(bridge.ContractAddr), deployerAddr, setbank.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, setOracle)
	//step7 Sign BridgeRegistry
	bridgeRegistry := signBridgeRegistry(common.HexToAddress(bridge.ContractAddr), common.HexToAddress(bank.ContractAddr), common.HexToAddress(oracle.ContractAddr), common.HexToAddress(valSet.ContractAddr), deployerAddr, setOracle.Nonce+1, gaslimit, gasprice, deployPrivateKey)
	signData = append(signData, bridgeRegistry)
	//finsh write to file
	writeToFile("signed_cross2eth.txt", signData)
}

func sendTxsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send", //first step
		Short: "send signed raw tx",
		Run:   sendTxs,
	}
	sendTxsFlags(cmd)
	return cmd
}

func sendTxsFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file", "f", "", "*.txt signed tx")
	_ = cmd.MarkFlagRequired("file")
}

func sendTxs(cmd *cobra.Command, args []string) {
	_ = args
	url, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	filePath, _ := cmd.Flags().GetString("file")
	//解析文件数据
	var rdata = make([]*eoff.DeployContract, 0)
	err := paraseFile(filePath, &rdata)
	if err != nil {
		fmt.Println("paraseFile,err", err.Error())
		return
	}
	var resData = make([]*DeployContractRet, 0)
	for _, deployInfo := range rdata {
		if deployInfo.Interval != 0 {
			time.Sleep(deployInfo.Interval)
		}
		tx := new(types.Transaction)
		err = tx.UnmarshalBinary(common.FromHex(deployInfo.SignedRawTx))
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
		ret := &DeployContractRet{ContractAddr: deployInfo.ContractAddr, TxHash: tx.Hash().String(), ContractName: deployInfo.ContractName}
		resData = append(resData, ret)
		time.Sleep(time.Second)
	}

	data, err := json.MarshalIndent(resData, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	fmt.Println(string(data))
}
