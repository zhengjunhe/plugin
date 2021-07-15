package offline

import (
	"fmt"
	//"github.com/33cn/plugin/plugin/dapp/cross2eth/boss4x/ethereum"
	"math/big"
	"strings"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4eth/generated"
	eoff "github.com/33cn/plugin/plugin/dapp/dex/boss/deploy/ethereum/offline"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)



type DepolyInfo struct {
	//OperatorAddr       string   `toml:"operatorAddr"`
	DeployerPrivateKey string   `toml:"deployerPrivateKey"`
	ValidatorsAddr     []string `toml:"validatorsAddr"`
	InitPowers         []int64  `toml:"initPowers"`
}

func  SignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign", //first step
		Short: "sign tx",
		Run:   signTx,
	}
	addSignFlag(cmd)
	return cmd
}

func  addSignFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "k", "", "private key ")
	cmd.MarkFlagRequired("key")
	cmd.Flags().StringP("file", "f", "deploytxs.txt", "tx file")


}

func signTx(cmd *cobra.Command, args []string) {
	privatekey, _ := cmd.Flags().GetString("key")
	txFilePath,_:=cmd.Flags().GetString("file")
	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(privatekey))
	if err != nil {
		panic(err)
	}

	var deployTxInfos = make([]DeployInfo, 0)
	err = paraseFile(txFilePath, &deployTxInfos)
	if err != nil {
		fmt.Println("paraseFile,err", err.Error())
		return
	}
	fmt.Println("deployTxInfos size:",len(deployTxInfos))
	for i ,info:=range deployTxInfos{
		var tx types.Transaction
		err= tx.UnmarshalBinary(common.FromHex(info.RawTx))
		if err!=nil{
			panic(err)
		}
		signedTx,txHash,err:=eoff.SignTx(deployPrivateKey,&tx)
		if err!=nil{
			panic(err)
		}
		deployTxInfos[i].RawTx=signedTx
		deployTxInfos[i].TxHash=txHash

	}

	//finsh write to file
	writeToFile("deploysigntxs.txt", deployTxInfos)

}

//deploy contract step 1
func deployValSetPackData(validators []common.Address, powers []*big.Int,deployerAddr common.Address)([]byte,error){
	parsed, err := abi.JSON(strings.NewReader(generated.ValsetABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.ValsetBin)
	packdata, err := parsed.Pack("", deployerAddr, validators, powers)
	if err != nil {
		panic(err)
	}
	return append(bin,packdata...),nil
}
//deploy contract step 2
func deploychain33BridgePackData(deployerAddr ,valSetAddr common.Address)([]byte,error){
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.Chain33BridgeBin)
	input, err := parsed.Pack("", deployerAddr, valSetAddr)
	if err != nil {
		panic(err)
	}

	return append(bin, input...),nil
}
//deploy contract step 3
func deployOraclePackData(deployerAddr,valSetAddr,bridgeAddr common.Address)([]byte,error){
	parsed, err := abi.JSON(strings.NewReader(generated.OracleABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.OracleBin)
	packData, err := parsed.Pack("", deployerAddr, valSetAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	return  append(bin, packData...),nil
}

//deploy contract step 4
func deployBridgeBankPackData(deployerAddr,bridgeAddr ,oracalAddr common.Address)([]byte,error){
	parsed, err := abi.JSON(strings.NewReader(generated.BridgeBankABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeBankBin)
	packData, err := parsed.Pack("", deployerAddr, oracalAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	return append(bin, packData...),nil
}

////deploy contract step 5
func callSetBridgeBank(bridgeBankAddr common.Address)([]byte,error){
	method := "setBridgeBank"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	packData, err := parsed.Pack(method, bridgeBankAddr)
	if err != nil {
		panic(err)
	}

	return packData,nil
}

//deploy contract step 6
func callSetOracal(oracalAddr common.Address)([]byte,error){
	method := "setOracle"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	packData, err := parsed.Pack(method, oracalAddr)
	if err != nil {
		panic(err)
	}
	return packData,nil
}

//deploy contract step 7
func deployBridgeRegistry(chain33BridgeAddr,bridgeBankAddr,oracleAddr,valSetAddr common.Address)([]byte,error){
	parsed, err := abi.JSON(strings.NewReader(generated.BridgeRegistryABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeRegistryBin)
	packData, err := parsed.Pack("", chain33BridgeAddr, bridgeBankAddr, oracleAddr, valSetAddr)
	if err != nil {
		panic(err)
	}
	return  append(bin, packData...),nil


}


