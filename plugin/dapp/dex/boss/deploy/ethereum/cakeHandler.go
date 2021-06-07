package ethereum

import (
	"context"
	"crypto/ecdsa"

	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeFactory"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeRouter"

	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//var TestNodeAddr = "https://data-seed-prebsc-1-s1.binance.org:8545/"
var ethClient *ethclient.Client
var privateKey *ecdsa.PrivateKey
var deployerAddr common.Address

//const ...
const (
	// GasLimit : the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit        = uint64(100 * 10000)
	GasLimit4Deploy = uint64(0) //此处需要设置为0,让交易自行估计,否则将会导致部署失败,TODO:其他解决途径后续调研解决
	fee2setter      = "0x0f2e821517D4f64a012a04b668a6b1aa3B262e08"
)

func setupWebsocketEthClient(ethNodeAddr string) {
	var err error
	ethClient, err = SetupWebsocketEthClient(ethNodeAddr)
	if nil != err {
		panic(fmt.Sprintf("Failed to SetupWebsocketEthClient with url:%s", ethNodeAddr))
	}
	fmt.Println("Succeed to establish connection to bsc")
}

// SetupWebsocketEthClient : returns boolean indicating if a URL is valid websocket ethclient
func SetupWebsocketEthClient(ethURL string) (*ethclient.Client, error) {
	if strings.TrimSpace(ethURL) == "" {
		return nil, nil
	}

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		return nil, fmt.Errorf("error dialing websocket client %w", err)
	}

	return client, nil
}

func GetBalance(tokenAddr, owner string) (string, error) {
	//查询ERC20余额
	//if tokenAddr != "" {
	//	bridgeToken, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
	//	if nil != err {
	//		return "", err
	//	}
	//	ownerAddr := common.HexToAddress(owner)
	//	opts := &bind.CallOpts{
	//		Pending: true,
	//		From:    ownerAddr,
	//		Context: context.Background(),
	//	}
	//	balance, err := bridgeToken.BalanceOf(opts, ownerAddr)
	//	if nil != err {
	//		return "", err
	//	}
	//	return balance.String(), nil
	//}

	//查询ETH余额
	balance, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(owner), nil)
	if nil != err {
		return "", err
	}
	return balance.String(), nil
}

func getNonce(sender common.Address) (*big.Int, error) {
	nonce, err := ethClient.PendingNonceAt(context.Background(), sender)
	if nil != err {
		return nil, err
	}
	return big.NewInt(int64(nonce)), nil
}

//PrepareAuth ...
func PrepareAuth(privateKey *ecdsa.PrivateKey, transactor common.Address) (*bind.TransactOpts, error) {
	//ctx := context.Background()
	//gasPrice, err := ethClient.SuggestGasPrice(ctx)
	//if err != nil {
	//	return nil, errors.New("failed to get suggest gas price")
	//}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = GasLimit4Deploy
	//auth.GasPrice = gasPrice

	var err error
	if auth.Nonce, err = getNonce(transactor); err != nil {
		return nil, err
	}

	return auth, nil
}

func DeployPancake() error {
	_ = recoverBinancePrivateKey()
	//1st step to deploy factory
	auth, err := PrepareAuth(privateKey, deployerAddr)
	if nil != err {
		return err
	}

	fee2setterAddr := common.HexToAddress(fee2setter)
	factoryAddr, deployFactoryTx, _, err := pancakeFactory.DeployPancakeFactory(auth, ethClient, fee2setterAddr)
	if nil != err {
		panic(fmt.Sprintf("Failed to DeployPancakeFactory with err:%s", err.Error()))
		return err
	}

	{
		fmt.Println("\nDeployPancakeFactory tx hash:", deployFactoryTx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployPancakeFactory timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), deployFactoryTx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for DeployPancakeFactory tx and continue to wait")
					continue
				} else if err != nil {
					panic("DeployPancakeFactory failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to deploy pancakeFactory with address =", factoryAddr.String())
				goto deployWeth9
			}
		}
	}

deployWeth9:
	//部署合约 PancakeRouter
	auth, err = PrepareAuth(privateKey, deployerAddr)
	if nil != err {
		return err
	}

	weth9Addr, deployWeth9Tx, _, err := pancakeRouter.DeployWETH9(auth, ethClient)
	if err != nil {
		panic(fmt.Sprintf("Failed to DeployWETH9 with err:%s", err.Error()))
		return err
	}
	{
		fmt.Println("\nDeployWETH9 tx hash:", deployWeth9Tx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployWETH9 timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), deployWeth9Tx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for deployWeth9Tx tx and continue to wait")
					continue
				} else if err != nil {
					panic("deployWeth9Tx failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to deploy weth9 with address =", weth9Addr.String())
				goto deployPancakeRouter
			}
		}
	}

deployPancakeRouter:
	//部署合约 PancakeRouter
	auth, err = PrepareAuth(privateKey, deployerAddr)
	if nil != err {
		return err
	}
	pancakeRouterAddr, deploypancakeTx, _, err := pancakeRouter.DeployPancakeRouter(auth, ethClient, factoryAddr, weth9Addr)
	if err != nil {
		panic(fmt.Sprintf("Failed to DeployPancakeRouter with err:%s", err.Error()))
		return err
	}

	{
		fmt.Println("\nDeployPancakeRouter tx hash:", deploypancakeTx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployPancakeRouter timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), deploypancakeTx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for DeployPancakeRouter tx and continue to wait")
					continue
				} else if err != nil {
					panic("DeployPancakeRouter failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to deploy pancakeRouter with address =", pancakeRouterAddr.String())
				return nil
			}
		}
	}

	return nil
}

func recoverBinancePrivateKey() (err error) {
	//louyuqi: f726c7c704e57ec5d59815dda23ddd794f71ae15f7e0141f00f73eff35334ac6
	//hzj: 2bcf3e23a17d3f3b190a26a098239ad2d20267a673440e0f57a23f44f94b77b9
	privateKey, err = crypto.ToECDSA(common.FromHex("2bcf3e23a17d3f3b190a26a098239ad2d20267a673440e0f57a23f44f94b77b9"))
	if nil != err {
		panic("Failed to recover private key")
		return err
	}
	deployerAddr = crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Println("The address is:", deployerAddr.String())

	return nil
}


//DeployVaultFactory tx hash: 0x97c0e9b7f42aed3ddf5bbead6d042fcd5609ee1dde9113e31f54b2afae56b551
//Succeed to deploy DeployVaultFactory with address = 0xe534945F98f344d6D5E53e0E747A44704c7C3806
//
//DeployVault tx hash: 0x5b1519a8bed301d349f22ef9cbdde37042a4321578f8e0cbcd1f4325c7e2f32a
//Succeed to deploy DeployVault with address = 0xe534945F98f344d6D5E53e0E747A44704c7C3806
//
//Computered address is 0x0183661e6b9288ebF98De625B4501bCF05c7b4cD
//   last Vault Addr is 0x0183661e6b9288ebF98De625B4501bCF05c7b4cD
//Succeed to deploy contracts

func AddAllowance4LPHandle(lp string, spender string, amount int64) (err error) {
	_ = recoverBinancePrivateKey()
	pairInt, err := pancakeFactory.NewPancakePair(common.HexToAddress(lp), ethClient)
	if nil != err {
		return err
	}
	auth, err := PrepareAuth(privateKey, deployerAddr)
	if nil != err {
		return err
	}
	approveTx, err := pairInt.Approve(auth, common.HexToAddress(spender), big.NewInt(amount))
	if nil != err {
		panic(fmt.Sprintf("Failed to Approve with err:%s", err.Error()))
		return err
	}

	{
		fmt.Println("\nApprove tx hash:", approveTx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("Approve timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), approveTx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for Approve tx and continue to wait")
					continue
				} else if err != nil {
					panic("Approve failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to do the approve operation")
				goto checkAllowance
				return nil
			}
		}
	}

checkAllowance:
	opts := &bind.CallOpts{
		Pending: true,
		Context: context.Background(),
	}
	//allowance[owner][spender] = value;
	result, err := pairInt.Allowance(opts, deployerAddr, common.HexToAddress(spender))
	if nil != err {
		return err
	}
	fmt.Println("\n The allowance recetrived is:", result.Int64())

	return nil
}

func CheckAllowance4LPHandle(lp string, spender string) (err error) {
	_ = recoverBinancePrivateKey()
	pairInt, err := pancakeFactory.NewPancakePair(common.HexToAddress(lp), ethClient)
	if nil != err {
		return err
	}
	//checkAllowance:
	opts := &bind.CallOpts{
		Pending: true,
		Context: context.Background(),
	}
	//allowance[owner][spender] = value;
	result, err := pairInt.Allowance(opts, deployerAddr, common.HexToAddress(spender))
	if nil != err {
		return err
	}
	fmt.Println("\n The allowance recetrived is:", result.String())

	return nil
}
