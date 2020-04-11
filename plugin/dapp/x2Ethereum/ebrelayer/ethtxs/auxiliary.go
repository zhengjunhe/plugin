package ethtxs

import (
	"context"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type NewProphecyClaimPara struct {
	ClaimType uint8
	Chain33Sender []byte
	TokenAddr common.Address
	Symbol string
}


//发行token="BTY"
//NewProphecyClaim
//铸币NewOracleClaim,
//ProcessBridgeProphecy
//铸币成功
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
	err = waitEthTxFinished(client, tx.Hash())
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

////////////////////订阅LogNewProphecyClaim/////////////////
//chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
/////////////////NewProphecyClaim////////////////
func MakeNewProphecyClaim(newProphecyClaimPara *NewProphecyClaimPara, client *ethclient.Client, para *DeployPara, x2EthContracts *X2EthContracts) (string, error) {
	authVali, err := PrepareAuth(client, para.ValidatorPriKey[0], para.InitValidators[0])
	if nil != err {
		return "", err
	}

	amount := int64(99)
	ethReceiver := para.InitValidators[2]
	tx, err := x2EthContracts.Chain33Bridge.NewProphecyClaim(authVali, newProphecyClaimPara.ClaimType, newProphecyClaimPara.Chain33Sender, ethReceiver, newProphecyClaimPara.TokenAddr, newProphecyClaimPara.Symbol, big.NewInt(amount))
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash())
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

	err = waitEthTxFinished(client, tx.Hash())
	if nil != err {
		return tx.Hash().String(), err
	}
	return tx.Hash().String(), nil
}

func GetBalance(client *ethclient.Client, tokenAddr, owner common.Address) (int64, error) {
	bridgeToken, err := generated.NewBridgeToken(tokenAddr, client)
	if nil != err {
		return 0, err
	}
	opts := &bind.CallOpts{
		Pending: true,
		From:    owner,
		Context: context.Background(),
	}
	balance, err := bridgeToken.BalanceOf(opts, owner)
	if nil != err {
		return 0, err
	}
	return balance.Int64(), nil
}
