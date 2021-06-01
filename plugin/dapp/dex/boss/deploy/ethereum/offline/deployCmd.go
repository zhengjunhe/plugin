package offline

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeFactory"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
)

//SignFactoryCmd 构造部署factory 合约的交易，并对其签名输出到文件中
type SignFactoryCmd struct{
	From string
	Nonce uint64
	GasPrice uint64
	FactoryAddr string
	TxHash string
	Fee2Addr string
	Timestamp string
	SignedTx string



}
func (s *SignFactoryCmd)signFactoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign factory", //first step
		Short: "deploy pancake router to ethereum ",
		Run:   s.signFactory,//对要部署的factory合约进行签名
	}
	s.addFactoryFlags(cmd)
	return cmd
}


func  (s *SignFactoryCmd)addFactoryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file","f","accountinfo.txt","multi params")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringP("feeaddr", "", "", "fee2stter")
	cmd.MarkFlagRequired("feeaddr")
	cmd.Flags().StringP("priv", "p", "", "private key")

}




func  (s *SignFactoryCmd)signFactory(cmd *cobra.Command, args []string) {
	filePath,_:=cmd.Flags().GetString("file")
	fee2setter, _ := cmd.Flags().GetString("feeaddr")
	key,_:=cmd.Flags().GetString("priv")
	priv,addr,err:=	recoverBinancePrivateKey(key)
	if err!=nil{
		panic(err)
	}
	fmt.Println("recover addr",addr)
	//解析文件数据
	err=paraseFile(filePath,s)
	if err!=nil{
		return
	}


	if !strings.EqualFold(s.From,addr.String()){
		panic("deployed address mismatch!!!")
	}
	gasPrice:=big.NewInt(int64(s.GasPrice))
	fmt.Println("nonce:",s.Nonce,"gasprice:",s.GasPrice)
	err =s.signFactoryTx(fee2setter,priv,gasPrice,s.Nonce)
	if nil != err {
		fmt.Println("Failed to deploy contracts due to:", err.Error())
		return
	}

	fmt.Println("Succeed to signed deploy contracts")
}


func  (s *SignFactoryCmd)signFactoryTx(fee2setter string,key *ecdsa.PrivateKey,gasPrice *big.Int,nonce uint64) error {
	fee2setterAddr := common.HexToAddress(fee2setter)
	signedTx,txHash, err := s.reWriteDeplopyPancakeFactory(nonce,gasPrice,key ,fee2setterAddr)
	if nil != err {
		panic(fmt.Sprintf("Failed to DeployPancakeFactory with err:%s", err.Error()))
	}

	//fmt.Println("signedTx:",signedtx)
	//write to file,把签好名的交易写入文件
	from:= crypto.PubkeyToAddress(key.PublicKey)
	factoryAddr:= crypto.CreateAddress(from, nonce)

	var wData deploayFactory
	wData.TxHash=txHash
	wData.SignedRawTx=signedTx
	wData.Nonce=s.Nonce
	wData.FactoryAddr=factoryAddr.String()
	writeToFile("factorySigned.txt",&wData)
	return nil
}

//构造交易，签名交易
func  (s *SignFactoryCmd)reWriteDeplopyPancakeFactory(nonce uint64,gasPrice *big.Int,key*ecdsa.PrivateKey, fee2addr common.Address)(signedTx ,hash string,err error){
	parsed, err := abi.JSON(strings.NewReader(pancakeFactory.PancakeFactoryABI))
	if err != nil {
		return
	}
	input,err:=parsed.Pack("", fee2addr)
	if err != nil {
		return
	}
	abiBin:=pancakeFactory.PancakeFactoryBin
	data:=append(common.FromHex(abiBin),input...)
	var gasLimit uint64=150000
	var amount =new(big.Int)
	ntx:= types.NewTransaction(nonce, common.Address{}, amount, gasLimit, gasPrice, data)
	signer:=types.HomesteadSigner{}
	txhash:=signer.Hash(ntx)
	signature, err := crypto.Sign(txhash.Bytes(), key)
	if err != nil {
		return
	}
	tx,err:= ntx.WithSignature(signer, signature)
	if err!=nil{
		return
	}
	txBinary,err:=	tx.MarshalBinary()
	if err!=nil{
		return
	}
	hash =tx.Hash().String()
	signedTx=common.Bytes2Hex(txBinary[:])

	return


}


func recoverBinancePrivateKey(key string) (priv *ecdsa.PrivateKey,address common.Address, err error) {
	//louyuqi: f726c7c704e57ec5d59815dda23ddd794f71ae15f7e0141f00f73eff35334ac6
	//hzj: 2bcf3e23a17d3f3b190a26a098239ad2d20267a673440e0f57a23f44f94b77b9
	priv, err = crypto.ToECDSA(common.FromHex(key))
	if err!=nil {
		panic("Failed to recover private key")
	}
	address = crypto.PubkeyToAddress(priv.PublicKey)
	fmt.Println("The address is:", address.String())
	return
}


