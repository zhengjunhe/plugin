package ethtxs

import (
	"context"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type NewProphecyClaimPara struct {
	ClaimType uint8
	Chain33Sender []byte
	TokenAddr common.Address
	EthReceiver common.Address
	Symbol string
	Amount int64
}

func CreateBridgeToken(symbol string, client *ethclient.Client, para *DeployPara, x2EthDeployInfo *X2EthDeployInfo, x2EthContracts *X2EthContracts) (string, error) {
	ctx := context.Background()

	//订阅事件
	eventName := "LogNewBridgeToken"
	bridgeBankABI := LoadABI(BridgeBankABI)
	logNewBridgeTokenSig := bridgeBankABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.BridgeBank.Address},
	}
	// We will check logs for new events
	logs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if nil != err {
		txslog.Error("CreateBrigeToken", "failed to SubscribeFilterLogs", err.Error())
		return "", err
	}

	//创建token
	auth, err := PrepareAuth(client, para.DeployPrivateKey, para.Operator)
	if nil != err {
		return "", err
	}

	tx, err := x2EthContracts.BridgeBank.BridgeBankTransactor.CreateNewBridgeToken(auth, symbol)
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "CreateBridgeToken")
	if nil != err {
		return "", err
	}

	logEvent := &events.LogNewBridgeToken{}
	select {
	// Handle any errors
	case err := <-sub.Err():
		return "", err
	// vLog is raw event data
	case vLog := <-logs:
		// Check if the event is a 'LogLock' event
		if vLog.Topics[0].Hex() == logNewBridgeTokenSig {
			txslog.Debug("CreateBrigeToken","Witnessed new event", eventName, "Block number", vLog.BlockNumber)

			err = bridgeBankABI.Unpack(logEvent, eventName, vLog.Data)
			if nil != err {
				return "", err
			}
			if symbol != logEvent.Symbol {
				txslog.Error("CreateBrigeToken","symbol", symbol, "logEvent.Symbol", logEvent.Symbol)
			}
			txslog.Info("CreateBrigeToken","Witnessed new event", eventName, "Block number", vLog.BlockNumber, "token address", logEvent.Token.String())
			break
		}
	}
	return logEvent.Token.String(), nil
}

func CreateERC20Token(symbol string, client *ethclient.Client, para *DeployPara, x2EthDeployInfo *X2EthDeployInfo, x2EthContracts *X2EthContracts) (string, error) {
	auth, err := PrepareAuth(client, para.DeployPrivateKey, para.Operator)
	if nil != err {
		return "", err
	}

	tokenAddr, tx, _, err := generated.DeployBridgeToken(auth, client, symbol)
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "CreateERC20Token")
	if nil != err {
		return "", err
	}

	return tokenAddr.String(), nil
}

func MintERC20Token(tokenAddr, ownerAddr string, amount int64, client *ethclient.Client, para *DeployPara)  (string, error) {
	operatorAuth, err := PrepareAuth(client, para.DeployPrivateKey, para.Operator)
	if nil != err {
		return "", err
	}
	erc20TokenInstance, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
	if nil != err {
		return "", err
	}
	tx, err := erc20TokenInstance.Mint(operatorAuth, common.HexToAddress(ownerAddr), big.NewInt(amount))
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "MintERC20Token")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

func ApproveAllowance(ownerPrivateKeyStr, tokenAddr string, bridgeBank common.Address, amount int64, client *ethclient.Client,)  (string, error) {
	ownerPrivateKey, err := crypto.ToECDSA(common.FromHex(ownerPrivateKeyStr))
	if nil != err {
		return "", err
	}
	ownerAddr := crypto.PubkeyToAddress(ownerPrivateKey.PublicKey)

	auth, err := PrepareAuth(client, ownerPrivateKey, ownerAddr)
	if nil != err {
		return "", err
	}
	erc20TokenInstance, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
	if nil != err {
		return "", err
	}

	tx, err := erc20TokenInstance.Approve(auth, bridgeBank, big.NewInt(amount))
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "ApproveAllowance")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

func Burn(ownerPrivateKeyStr, tokenAddrstr, chain33Receiver string, bridgeBank common.Address, amount int64, bridgeBankIns *generated.BridgeBank, client *ethclient.Client,)  (string, error) {
	ownerPrivateKey, err := crypto.ToECDSA(common.FromHex(ownerPrivateKeyStr))
	if nil != err {
		return "", err
	}
	ownerAddr := crypto.PubkeyToAddress(ownerPrivateKey.PublicKey)
	auth, err := PrepareAuth(client, ownerPrivateKey, ownerAddr)
	if nil != err {
		return "", err
	}
	tokenAddr := common.HexToAddress(tokenAddrstr)
	tokenInstance, err := generated.NewBridgeToken(tokenAddr, client)
	if nil != err {
		return "", err
	}
	//chain33bank 是bridgeBank的基类，所以使用bridgeBank的地址
	tx, err := tokenInstance.Approve(auth, bridgeBank, big.NewInt(amount))
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "Approve")
	if nil != err {
		return "", err
	}
	txslog.Info("Burn","Approve tx with hash", tx.Hash().String())

	auth, err = PrepareAuth(client, ownerPrivateKey, ownerAddr)
	if nil != err {
		return "", err
	}
	tx, err = bridgeBankIns.BurnBridgeTokens(auth, []byte(chain33Receiver), tokenAddr, big.NewInt(amount))
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "Burn")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

func LockEthErc20Asset(ownerPrivateKeyStr, tokenAddrStr, chain33Receiver string, amount int64, client *ethclient.Client, bridgeBank *generated.BridgeBank)  (string, error) {
	ownerPrivateKey, err := crypto.ToECDSA(common.FromHex(ownerPrivateKeyStr))
	if nil != err {
		return "", err
	}
	ownerAddr := crypto.PubkeyToAddress(ownerPrivateKey.PublicKey)

	auth, err := PrepareAuth(client, ownerPrivateKey, ownerAddr)
	if nil != err {
		return "", err
	}
	//ETH转账，空地址，且设置value
	var tokenAddr common.Address
    if "" == tokenAddrStr {
		auth.Value = big.NewInt(amount)
	}

	if "" != tokenAddrStr {
		tokenAddr = common.HexToAddress(tokenAddrStr)
	}
	tx, err := bridgeBank.Lock(auth, []byte(chain33Receiver), tokenAddr, big.NewInt(amount))
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "LockEthErc20Asset")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

////////////////////订阅LogNewProphecyClaim/////////////////
//chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
/////////////////NewProphecyClaim////////////////
func MakeNewProphecyClaim(newProphecyClaimPara *NewProphecyClaimPara, client *ethclient.Client, para *DeployPara, x2EthContracts *X2EthContracts) (string, error) {
	authVali, err := PrepareAuth(client, para.ValidatorPriKey[0], para.InitValidators[0])
	if nil != err {
		return "", err
	}

	amount := newProphecyClaimPara.Amount
	ethReceiver := newProphecyClaimPara.EthReceiver
	tx, err := x2EthContracts.Chain33Bridge.NewProphecyClaim(authVali, newProphecyClaimPara.ClaimType, newProphecyClaimPara.Chain33Sender, ethReceiver, newProphecyClaimPara.TokenAddr, newProphecyClaimPara.Symbol, big.NewInt(amount))
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "MakeNewProphecyClaim")
	if nil != err {
		return "", err
	}
	return tx.Hash().String(), nil
}

func ProcessProphecyClaim(client *ethclient.Client, para *DeployPara, x2EthContracts *X2EthContracts, prophecyID int64) (string, error) {
	authOracle, err := PrepareAuth(client, para.ValidatorPriKey[0], para.InitValidators[0])
	if nil != err {
		return "", err
	}
	tx, err := x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, big.NewInt(prophecyID))
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "ProcessProphecyClaim")
	if nil != err {
		return tx.Hash().String(), err
	}
	return tx.Hash().String(), nil
}


