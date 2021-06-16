package offline

import (
	"crypto/ecdsa"
	"time"

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

type SignCmd struct {
	From         string
	Nonce        uint64
	GasPrice     uint64
	key          *ecdsa.PrivateKey
	deployerAddr common.Address
	Timestamp    string
}

type DepolyInfo struct {
	//OperatorAddr       string   `toml:"operatorAddr"`
	DeployerPrivateKey string   `toml:"deployerPrivateKey"`
	ValidatorsAddr     []string `toml:"validatorsAddr"`
	InitPowers         []int64  `toml:"initPowers"`
}

func (s *SignCmd) signCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign contract", //first step
		Short: " create contract ,and sign",
		Run:   s.SignContractTx, //对要部署的factory合约进行签名
	}
	s.addSignFlag(cmd)
	return cmd
}

func (s *SignCmd) addSignFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("conf", "c", "", "config file")
	cmd.MarkFlagRequired("conf")
	cmd.Flags().Uint64P("nonce","n",-1,"transaction count")
	cmd.MarkFlagRequired("nonce")
	cmd.Flags().Uint64P("gasprice","g",1000000000,"gas price")// 1Gwei=1e9wei
	cmd.MarkFlagRequired("gasprice")

}

func (s *SignCmd) SignContractTx(cmd *cobra.Command, args []string) {
	cfgpath, _ := cmd.Flags().GetString("conf")
	gasprice,_:=cmd.Flags().GetUint64("gasprice")
	nonce,_:=cmd.Flags().GetUint64("nonce")
	var deployCfg  DepolyInfo
	InitCfg(cfgpath, &deployCfg)



	deployPrivateKey, err := crypto.ToECDSA(common.FromHex(deployCfg.DeployerPrivateKey))
	if err != nil {
		//fmt.Printf("privkey",deployCfg.DeployerPrivateKey,"configpath",cfgpath)
		panic(err)
	}

	deployerAddr := crypto.PubkeyToAddress(deployPrivateKey.PublicKey)
	s.deployerAddr=deployerAddr
	s.Nonce=nonce
	s.GasPrice=gasprice
	s.From=deployerAddr.String()
	s.key=deployPrivateKey

	if len(deployCfg.InitPowers) != len(deployCfg.ValidatorsAddr) {
		panic("not same number for validator address and power")
	}

	if len(deployCfg.ValidatorsAddr) < 3 {
		panic("the number of validator must be not less than 3")
	}

	var validators []common.Address
	var initPowers []*big.Int
	var signData = make([]*eoff.DeployContract, 0)
	for i, addr := range deployCfg.ValidatorsAddr {
		validators = append(validators, common.HexToAddress(addr))
		initPowers = append(initPowers, big.NewInt(deployCfg.InitPowers[i]))
	}

	//Step1 Sign ValSet Contractor
	valSet := s.signValSet(validators, initPowers)
	signData = append(signData, valSet)
	//Step2 Sign chain33 bridge
	bridge := s.signChain33Bridge(valSet.Nonce+1, deployerAddr, common.HexToAddress(valSet.ContractAddr))
	signData = append(signData, bridge)
	//step3 sign oracle
	oracle := s.signOracle(bridge.Nonce+1, common.HexToAddress(valSet.ContractAddr), common.HexToAddress(bridge.ContractAddr))
	signData = append(signData, oracle)
	//step4  sign bridgebank
	bank := s.SignBridgeBank(oracle.Nonce+1, common.HexToAddress(bridge.ContractAddr), common.HexToAddress(oracle.ContractAddr))
	signData = append(signData, bank)
	//step5 sign SetBridgeBank
	setbank := s.SignSetBridgeBank(bank.Nonce+1, common.HexToAddress(bank.ContractAddr), common.HexToAddress(bridge.ContractAddr))
	signData = append(signData, setbank)
	//step6 Sign setOracle
	setOracle := s.SignsetOracle(setbank.Nonce+1, common.HexToAddress(oracle.ContractAddr), common.HexToAddress(bridge.ContractAddr))
	signData = append(signData, setOracle)
	//step7 Sign BridgeRegistry
	bridgeRegistry := s.SignBridgeRegistry(setOracle.Nonce+1, common.HexToAddress(bridge.ContractAddr), common.HexToAddress(bank.ContractAddr), common.HexToAddress(oracle.ContractAddr), common.HexToAddress(valSet.ContractAddr))
	signData = append(signData, bridgeRegistry)
	//finsh write to file
	writeToFile("signed_cross2eth.txt", signData)

}

func (s *SignCmd) signValSet(validators []common.Address, powers []*big.Int) *eoff.DeployContract {

	parsed, err := abi.JSON(strings.NewReader(generated.ValsetABI))
	if err != nil {
		panic(err)
	}
	vbin := common.FromHex(generated.ValsetBin)

	input, err := parsed.Pack("", s.deployerAddr, validators, powers)
	if err != nil {
		panic(err)
	}
	data := append(vbin, input...)
	rawTx := types.NewContractCreation(s.Nonce, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), data)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	contractAddress := crypto.CreateAddress(s.deployerAddr, s.Nonce)
	var valSet eoff.DeployContract
	valSet.Nonce = s.Nonce
	valSet.ContractName = "valset"
	valSet.SignedRawTx = signedtx
	valSet.ContractAddr = contractAddress.String()
	valSet.TxHash = hash
	return &valSet

}

func (s *SignCmd) signChain33Bridge(nonce uint64, operater, valSetAddr common.Address) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	bridgebin := common.FromHex(generated.Chain33BridgeBin)
	input, err := parsed.Pack("", operater, valSetAddr)
	if err != nil {
		panic(err)
	}

	data := append(bridgebin, input...)
	rawTx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), data)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	contractAddress := crypto.CreateAddress(operater, nonce)
	var bridge eoff.DeployContract
	bridge.Nonce = nonce
	bridge.ContractName = "chain33bridge"
	bridge.SignedRawTx = signedtx
	bridge.ContractAddr = contractAddress.String()
	bridge.TxHash = hash
	return &bridge
}

func (s *SignCmd) signOracle(nonce uint64, valsetAddr, bridgeAddr common.Address) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.OracleABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.OracleBin)
	input, err := parsed.Pack("", s.deployerAddr, valsetAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	data := append(bin, input...)
	rawTx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), data)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	contractAddress := crypto.CreateAddress(s.deployerAddr, nonce)
	var oracle eoff.DeployContract
	oracle.Nonce = nonce
	oracle.ContractName = "oracle"
	oracle.SignedRawTx = signedtx
	oracle.ContractAddr = contractAddress.String()
	oracle.TxHash = hash
	return &oracle
}

func (s *SignCmd) SignBridgeBank(nonce uint64, bridgeAddr, oracalAddr common.Address) *eoff.DeployContract {

	parsed, err := abi.JSON(strings.NewReader(generated.BridgeBankABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeBankBin)
	input, err := parsed.Pack("", s.deployerAddr, oracalAddr, bridgeAddr)
	if err != nil {
		panic(err)
	}

	data := append(bin, input...)
	rawTx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), data)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	contractAddress := crypto.CreateAddress(s.deployerAddr, nonce)
	var bank eoff.DeployContract
	bank.Nonce = nonce
	bank.ContractName = "bridgeBank"
	bank.SignedRawTx = signedtx
	bank.ContractAddr = contractAddress.String()
	bank.TxHash = hash
	return &bank
}

//SignSetBridgeBank SetBridgeBank
func (s *SignCmd) SignSetBridgeBank(nonce uint64, bridgebank, chain33bridge common.Address) *eoff.DeployContract {
	method := "setBridgeBank"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	input, err := parsed.Pack(method, bridgebank)
	if err != nil {
		panic(err)
	}

	rawTx := types.NewTransaction(nonce, chain33bridge, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), input)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	if err != nil {
		panic(err)
	}
	contractAddress := crypto.CreateAddress(s.deployerAddr, nonce)
	var setbank eoff.DeployContract
	setbank.Interval = time.Second * 20
	setbank.Nonce = nonce
	setbank.ContractName = "setbridgebank"
	setbank.SignedRawTx = signedtx
	setbank.ContractAddr = contractAddress.String()
	setbank.TxHash = hash
	return &setbank

}

func (s *SignCmd) SignsetOracle(nonce uint64, oracalAddr, chain33bridge common.Address) *eoff.DeployContract {
	method := "setOracle"
	parsed, err := abi.JSON(strings.NewReader(generated.Chain33BridgeABI))
	if err != nil {
		panic(err)
	}
	input, err := parsed.Pack(method, oracalAddr)
	if err != nil {
		panic(err)
	}
	rawTx := types.NewTransaction(nonce, chain33bridge, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), input)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	if err != nil {
		panic(err)
	}
	contractAddress := crypto.CreateAddress(s.deployerAddr, nonce)
	var setoracle eoff.DeployContract
	setoracle.Interval = time.Second * 20
	setoracle.Nonce = nonce
	setoracle.ContractName = "setOracle"
	setoracle.SignedRawTx = signedtx
	setoracle.ContractAddr = contractAddress.String()
	setoracle.TxHash = hash
	return &setoracle

}

func (s *SignCmd) SignBridgeRegistry(nonce uint64, chain33Bridge, bridgebank, oracleAddr, valSetAddr common.Address) *eoff.DeployContract {
	parsed, err := abi.JSON(strings.NewReader(generated.BridgeRegistryABI))
	if err != nil {
		panic(err)
	}
	bin := common.FromHex(generated.BridgeRegistryBin)
	input, err := parsed.Pack("", chain33Bridge, bridgebank, oracleAddr, valSetAddr)
	if err != nil {
		panic(err)
	}
	data := append(bin, input...)
	rawTx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, big.NewInt(int64(s.GasPrice)), data)
	//signedtx
	signedtx, hash, err := eoff.SignTx(s.key, rawTx)
	if err != nil {
		panic(err)
	}
	contractAddress := crypto.CreateAddress(s.deployerAddr, nonce)
	var bridgeRegistry eoff.DeployContract
	bridgeRegistry.Nonce = nonce
	bridgeRegistry.ContractName = "bridgeRegistry"
	bridgeRegistry.SignedRawTx = signedtx
	bridgeRegistry.ContractAddr = contractAddress.String()
	bridgeRegistry.TxHash = hash
	return &bridgeRegistry

}
