package offline

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeFactory"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeRouter"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"math/big"
	"strings"
	"time"
)

//SignFactoryCmd 构造部署factory 合约的交易，并对其签名输出到文件中
type SignCmd struct {
	From        string
	Nonce       uint64
	GasPrice    uint64
	FactoryAddr string
	TxHash      string
	Fee2Addr    string
	Timestamp   string
	SignedTx    string
	Reward uint64
	StartBlock uint64
}

func (s *SignCmd) signCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign contract", //first step
		Short: "sign pancake and farm contract to ethereum ",
		Run:   s.signContract, //对要部署的factory合约进行签名
	}
	s.addFlags(cmd)
	return cmd
}

func (s *SignCmd) addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("file", "f", "accountinfo.txt", "multi params")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringP("feeaddr", "", "", "fee2stter")
	cmd.MarkFlagRequired("feeaddr")
	cmd.Flags().StringP("priv", "p", "", "private key")
	cmd.Flags().Int64P("reward","",5,"Set the reward for each block")
	cmd.MarkFlagRequired("reward")
	cmd.Flags().Int64P("start","",-1,"Set effective height")
	cmd.MarkFlagRequired("start")


}

func (s *SignCmd) signContract(cmd *cobra.Command, args []string) {
	filePath, _ := cmd.Flags().GetString("file")
	fee2setter, _ := cmd.Flags().GetString("feeaddr")
	key, _ := cmd.Flags().GetString("priv")
	reward,_:=cmd.Flags().GetInt64("reward")
	startBlock,_:=cmd.Flags().GetInt64("start")

	if startBlock<=0{
		panic("startBlock  err")
	}

	priv, addr, err := recoverBinancePrivateKey(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("recover addr", addr)
	//解析文件数据
	err = paraseFile(filePath, s)
	if err != nil {
		return
	}
	s.Reward= uint64(reward)
	s.StartBlock= uint64(startBlock)
	//check is timeout
	t,err:=time.Parse(time.RFC3339,s.Timestamp)
	if err!=nil{
		panic(err)
	}
	if time.Now().After(t.Add(time.Hour)){
		panic("after 60 minutes timeout,the accountinfo.txt invalid,please reQuery")
	}
	if !strings.EqualFold(s.From, addr.String()) {
		panic("deployed address mismatch!!!")
	}
	gasPrice := big.NewInt(int64(s.GasPrice))
	//fmt.Println("nonce:",s.Nonce,"gasprice:",s.GasPrice)
	err = s.signContractTx(fee2setter, priv, gasPrice, s.Nonce)
	if nil != err {
		fmt.Println("Failed to deploy contracts due to:", err.Error())
		return
	}

	fmt.Println("Succeed to signed deploy contracts")
}

func (s *SignCmd) signContractTx(fee2setter string, key *ecdsa.PrivateKey, gasPrice *big.Int, nonce uint64) error {
	fee2setterAddr := common.HexToAddress(fee2setter)
	from := crypto.PubkeyToAddress(key.PublicKey)
	//--------------------
	//sign factory
	//--------------------
	signedTx, txHash, err := s.reWriteDeplopyPancakeFactory(nonce, gasPrice, key, fee2setterAddr)
	if nil != err {
		panic(fmt.Sprintf("Failed to DeployPancakeFactory with err:%s", err.Error()))
	}

	factoryAddr := crypto.CreateAddress(from, nonce)
	var signData = make([]*deploayContract, 0)
	var factData deploayContract
	factData.TxHash = txHash
	factData.SignedRawTx = signedTx
	factData.Nonce = s.Nonce
	factData.ContractAddr = factoryAddr.String()
	factData.ContractName = "factory"
	signData = append(signData, &factData)
	//--------------------
	//sign weth9
	//--------------------
	weth := new(SignWeth9Cmd)
	wsignedTx, hash, err := weth.reWriteDeployWETH9(s.Nonce+1, gasPrice, key)
	if nil != err {
		panic(fmt.Sprintf("Failed to DeployPancakeFactory with err:%s", err.Error()))
	}
	weth9Addr := crypto.CreateAddress(from, s.Nonce+1)
	var weth9Data deploayContract
	weth9Data.Nonce = s.Nonce + 1
	weth9Data.TxHash = hash
	weth9Data.SignedRawTx = wsignedTx
	weth9Data.ContractAddr = weth9Addr.String()
	weth9Data.ContractName = "weth9"
	signData = append(signData, &weth9Data)
	//--------------------
	//sign PanCakeRouter
	//--------------------
	panRouter := new(SignPanCakeRout)
	rSignedTx, hash, err := panRouter.reWriteDeployPanCakeRout(weth9Data.Nonce+1, gasPrice, key, factoryAddr, weth9Addr)
	if nil != err {
		panic(fmt.Sprintf("Failed to reWriteDeployPanCakeRout with err:%s", err.Error()))
	}
	panrouterAddr := crypto.CreateAddress(from, weth9Data.Nonce+1)
	var panData deploayContract
	panData.Nonce = weth9Data.Nonce + 1
	panData.SignedRawTx = rSignedTx
	panData.ContractAddr = panrouterAddr.String()
	panData.TxHash = hash
	panData.ContractName = "pancakerouter"
	signData = append(signData, &panData)

	/**************************pancake contract signed completed*************************/
	//--------------------let's begain Farm contract^_^--------------------
	//Sign Farm Contractor
	//--------------------
	farmNonce:=panData.Nonce+1
	var cakeToken=new(SignCakeToken)
	var cakeData=new(deploayContract)
	cakeSignedtx,hash,err:=cakeToken.reWriteDeployCakeToken(farmNonce,gasPrice,key)
	if nil != err {
		panic(fmt.Sprintf("Failed to reWriteDeployCakeToken with err:%s", err.Error()))
	}
	cakeContractAddr:=crypto.CreateAddress(from,farmNonce)
	cakeData.Nonce=farmNonce
	cakeData.SignedRawTx=cakeSignedtx
	cakeData.TxHash=hash
	cakeData.ContractName="caketoken"
	cakeData.ContractAddr=cakeContractAddr.String()
	signData = append(signData, cakeData)
	//--------------------
	//Sign syrupBar Contractor
	//--------------------
	syrupBarNonce:=farmNonce+1
	var syrupBar=new(signsyrupBar)
	var syrupBarData=new(deploayContract)
	syupSignedTx,hash,err:=	syrupBar.reWriteDeploysyrupBar(syrupBarNonce,gasPrice,key,cakeContractAddr)
	if err!=nil{
		panic(err)
	}
	syupContractAddr:=crypto.CreateAddress(from,syrupBarNonce)
	syrupBarData.Nonce=syrupBarNonce
	syrupBarData.TxHash=hash
	syrupBarData.ContractName="syrupbar"
	syrupBarData.ContractAddr=syupContractAddr.String()
	syrupBarData.SignedRawTx=syupSignedTx
	signData=append(signData,syrupBarData)
	//--------------------
	//Sign masterChef Contractor
	//--------------------
	masterChefNonce:=syrupBarNonce+1
	var mChefData=new(deploayContract)
	var mChef=new(signMasterChef)
	reward:=big.NewInt(int64(s.Reward * 1e18))
	startBlockHeight:=big.NewInt(int64(s.StartBlock))
	mchefSignedTx,hash,err:=mChef.reWriteDeployMasterChef(masterChefNonce,gasPrice,key,cakeContractAddr,syupContractAddr,from,reward,startBlockHeight)
	if err!=nil{
		panic(err)
	}
	mChefContractAddr:=crypto.CreateAddress(from,masterChefNonce)
	mChefData.Nonce=masterChefNonce
	mChefData.TxHash=hash
	mChefData.ContractName="masterchef"
	mChefData.SignedRawTx=mchefSignedTx
	mChefData.ContractAddr=mChefContractAddr.String()
	signData=append(signData,mChefData)

	//write signedtx to spec file
	writeToFile("signed.txt", &signData)
	return nil
}

//构造交易，签名交易
func (s *SignCmd) reWriteDeplopyPancakeFactory(nonce uint64, gasPrice *big.Int, key *ecdsa.PrivateKey, fee2addr common.Address,param...interface{}) (signedTx, hash string, err error) {
	parsed, err := abi.JSON(strings.NewReader(pancakeFactory.PancakeFactoryABI))
	if err != nil {
		return
	}
	input, err := parsed.Pack("", fee2addr)
	if err != nil {
		return
	}
	abiBin := pancakeFactory.PancakeFactoryBin
	data := append(common.FromHex(abiBin), input...)
	var gasLimit uint64 = 150000
	var amount = new(big.Int)
	ntx:=types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	//ntx := types.NewTransaction(nonce, common.Address{}, amount, gasLimit, gasPrice, data)
	return signTx(key, ntx)
}


type SignWeth9Cmd struct {
}

//only sign
func (s *SignWeth9Cmd) reWriteDeployWETH9(nonce uint64, gasPrice *big.Int, key *ecdsa.PrivateKey,param...interface{}) (signedTx, hash string, err error) {
	parsed, err := abi.JSON(strings.NewReader(pancakeRouter.WETH9ABI))
	if err != nil {
		return "", "", err
	}
	input, err := parsed.Pack("", param...)
	abiBin := pancakeRouter.PancakeRouterBin
	data := append(common.FromHex(abiBin), input...)
	var gasLimit uint64 = 150000
	var amount = new(big.Int)
	ntx:=types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	return signTx(key, ntx)
}

type SignPanCakeRout struct {
}

func (s *SignPanCakeRout) reWriteDeployPanCakeRout(nonce uint64, gasPrice *big.Int, key *ecdsa.PrivateKey, factoryAddr, Weth9 common.Address) (signedTx, hash string, err error) {
	parsed, err := abi.JSON(strings.NewReader(pancakeRouter.PancakeRouterABI))
	if err != nil {
		return
	}
	input, err := parsed.Pack("", factoryAddr, Weth9)
	if err != nil {
		return
	}
	abiBin := pancakeRouter.PancakeRouterBin
	data := append(common.FromHex(abiBin), input...)
	var gasLimit uint64 = 150000
	var amount = new(big.Int)
	ntx:=types.NewContractCreation(nonce, amount, gasLimit, gasPrice, data)
	return signTx(key, ntx)

}

func signTx(key *ecdsa.PrivateKey, tx *types.Transaction) (signedTx, hash  string, err error) {
	signer := types.HomesteadSigner{}
	txhash := signer.Hash(tx)
	signature, err := crypto.Sign(txhash.Bytes(), key)
	if err != nil {
		return
	}
	tx, err = tx.WithSignature(signer, signature)
	if err != nil {
		return
	}
	txBinary, err := tx.MarshalBinary()
	if err != nil {
		return
	}
	hash = tx.Hash().String()
	signedTx = common.Bytes2Hex(txBinary[:])
	return
}

func recoverBinancePrivateKey(key string) (priv *ecdsa.PrivateKey, address common.Address, err error) {
	//louyuqi: f726c7c704e57ec5d59815dda23ddd794f71ae15f7e0141f00f73eff35334ac6
	//hzj: 2bcf3e23a17d3f3b190a26a098239ad2d20267a673440e0f57a23f44f94b77b9
	priv, err = crypto.ToECDSA(common.FromHex(key))
	if err != nil {
		panic("Failed to recover private key")
	}
	address = crypto.PubkeyToAddress(priv.PublicKey)
	fmt.Println("The address is:", address.String())
	return
}
