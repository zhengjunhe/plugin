package offline

import (
	"context"
	"encoding/json"
	"fmt"
	tml "github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/big"
	"os"
)

//查询deploy 私钥的nonce信息，并输出到文件中

func TxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx", //first step
		Short: "create deploy tx",
		Run:   newTx, //对要部署的factory合约进行签名
	}
	addQueryFlags(cmd)
	return cmd
}

func addQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("conf", "c", "", "config file")
	cmd.MarkFlagRequired("conf")
}

type DeployInfo struct {
	Name           string
	PackData       []byte
	ContractorAddr common.Address
	Nonce          uint64
	To             *common.Address
	RawTx          string
	TxHash         string
	Gas            uint64
}

func newTx(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("rpc_laddr_ethereum")
	cfgpath, _ := cmd.Flags().GetString("conf")
	var deployCfg DepolyInfo
	InitCfg(cfgpath, &deployCfg)
	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(deployCfg.DeployerPrivateKey))
	if err != nil {
		panic(err)
	}

	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)
	if len(deployCfg.InitPowers) != len(deployCfg.ValidatorsAddr) {
		panic("not same number for validator address and power")
	}

	if len(deployCfg.ValidatorsAddr) < 3 {
		panic("the number of validator must be not less than 3")
	}

	var validators []common.Address
	var initPowers []*big.Int
	for i, addr := range deployCfg.ValidatorsAddr {
		validators = append(validators, common.HexToAddress(addr))
		initPowers = append(initPowers, big.NewInt(deployCfg.InitPowers[i]))
	}

	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	price, err := client.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}

	startNonce, err := client.PendingNonceAt(ctx, deployerAddr)
	if nil != err {
		panic(err)
	}

	var infos []*DeployInfo
	//step1 valSet
	packData, err := deployValSetPackData(validators, initPowers, deployerAddr)
	if err != nil {
		panic(err)
	}
	valSetAddr := crypto.CreateAddress(deployerAddr, startNonce)
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: valSetAddr, Name: "valSet", Nonce: startNonce, To: nil})

	//step2 chain33bridge
	packData, err = deploychain33BridgePackData(deployerAddr, valSetAddr)
	if err != nil {
		panic(err)
	}
	chain33BridgeAddr := crypto.CreateAddress(deployerAddr, startNonce+1)
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: chain33BridgeAddr, Name: "chain33Bridge", Nonce: startNonce + 1, To: nil})

	//step3 oracle
	packData, err = deployOraclePackData(deployerAddr, valSetAddr, chain33BridgeAddr)
	if err != nil {
		panic(err)
	}
	oracleAddr := crypto.CreateAddress(deployerAddr, startNonce+2)
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: oracleAddr, Name: "oracle", Nonce: startNonce + 2, To: nil})

	//step4 bridgebank
	packData, err = deployBridgeBankPackData(deployerAddr, chain33BridgeAddr, oracleAddr)
	if err != nil {
		panic(err)
	}
	bridgeBankAddr := crypto.CreateAddress(deployerAddr, startNonce+3)
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: bridgeBankAddr, Name: "bridgebank", Nonce: startNonce + 3, To: nil})

	//step5
	packData, err = callSetBridgeBank(bridgeBankAddr)
	if err != nil {
		panic(err)
	}
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: common.Address{}, Name: "setbridgebank", Nonce: startNonce + 4, To: &chain33BridgeAddr})

	//step6
	packData, err = callSetOracal(oracleAddr)
	if err != nil {
		panic(err)
	}
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: common.Address{}, Name: "setoracle", Nonce: startNonce + 5, To: &chain33BridgeAddr})

	//step7 bridgeRegistry
	packData, err = deployBridgeRegistry(chain33BridgeAddr, bridgeBankAddr, oracleAddr, valSetAddr)
	if err != nil {
		panic(err)
	}
	bridgeRegAddr := crypto.CreateAddress(deployerAddr, startNonce+6)
	infos = append(infos, &DeployInfo{PackData: packData, ContractorAddr: bridgeRegAddr, Name: "bridgeRegistry", Nonce: startNonce + 6, To: nil})

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
			panic(err)
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
			panic(err)
		}
		infos[i].RawTx = common.Bytes2Hex(txBytes)
		infos[i].Gas = gasLimit
	}

	jbytes, err := json.MarshalIndent(&infos, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jbytes))
	writeToFile("deploytxs.txt", &infos)
	return

}

func paraseFile(file string, result interface{}) error {
	_, err := os.Stat(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(b, result)
}

func writeToFile(fileName string, content interface{}) {
	jbytes, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, jbytes, 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
	fmt.Println("tx is written to file: ", fileName)
	//fmt.Println("tx is written to file: ", fileName, "writeContent:", string(jbytes))
}

func InitCfg(filepath string, cfg interface{}) {
	if _, err := tml.DecodeFile(filepath, cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return
}
