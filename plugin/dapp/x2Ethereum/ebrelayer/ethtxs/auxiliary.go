package ethtxs

import (
	"context"
	"crypto/ecdsa"
	"errors"
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
	ClaimType     uint8
	Chain33Sender []byte
	TokenAddr     common.Address
	EthReceiver   common.Address
	Symbol        string
	Amount        *big.Int
	Txhash        []byte
}

func CreateBridgeToken(symbol string, client *ethclient.Client, para *OperatorInfo, x2EthDeployInfo *X2EthDeployInfo, x2EthContracts *X2EthContracts) (string, error) {
	if nil == para {
		return "", errors.New("No operator private key configured")
	}
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
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if nil != err {
		txslog.Error("CreateBrigeToken", "failed to SubscribeFilterLogs", err.Error())
		return "", err
	}

	//创建token
	auth, err := PrepareAuth(client, para.PrivateKey, para.Address)
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
			txslog.Debug("CreateBrigeToken", "Witnessed new event", eventName, "Block number", vLog.BlockNumber)

			err = bridgeBankABI.Unpack(logEvent, eventName, vLog.Data)
			if nil != err {
				return "", err
			}
			if symbol != logEvent.Symbol {
				txslog.Error("CreateBrigeToken", "symbol", symbol, "logEvent.Symbol", logEvent.Symbol)
			}
			txslog.Info("CreateBrigeToken", "Witnessed new event", eventName, "Block number", vLog.BlockNumber, "token address", logEvent.Token.String())
			break
		}
	}
	return logEvent.Token.String(), nil
}

func CreateERC20Token(symbol string, client *ethclient.Client, para *OperatorInfo, x2EthDeployInfo *X2EthDeployInfo, x2EthContracts *X2EthContracts) (string, error) {
	if nil == para {
		return "", errors.New("No operator private key configured")
	}
	auth, err := PrepareAuth(client, para.PrivateKey, para.Address)
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

func MintERC20Token(tokenAddr, ownerAddr string, amount *big.Int, client *ethclient.Client, para *OperatorInfo) (string, error) {
	if nil == para {
		return "", errors.New("No operator private key configured")
	}

	operatorAuth, err := PrepareAuth(client, para.PrivateKey, para.Address)
	if nil != err {
		return "", err
	}
	erc20TokenInstance, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
	if nil != err {
		return "", err
	}
	tx, err := erc20TokenInstance.Mint(operatorAuth, common.HexToAddress(ownerAddr), amount)
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "MintERC20Token")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

func ApproveAllowance(ownerPrivateKeyStr, tokenAddr string, bridgeBank common.Address, amount *big.Int, client *ethclient.Client) (string, error) {
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

	tx, err := erc20TokenInstance.Approve(auth, bridgeBank, amount)
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "ApproveAllowance")
	if nil != err {
		return "", err
	}

	return tx.Hash().String(), nil
}

func Burn(ownerPrivateKeyStr, tokenAddrstr, chain33Receiver string, bridgeBank common.Address, amount int64, bridgeBankIns *generated.BridgeBank, client *ethclient.Client) (string, error) {
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
	txslog.Info("Burn", "Approve tx with hash", tx.Hash().String())

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

func TransferToken(tokenAddr, fromPrivateKeyStr, toAddr string, amount *big.Int, client *ethclient.Client) (string, error) {
	tokenInstance, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
	if nil != err {
		return "", err
	}

	fromPrivateKey, err := crypto.ToECDSA(common.FromHex(fromPrivateKeyStr))
	if nil != err {
		return "", err
	}
	fromAddr := crypto.PubkeyToAddress(fromPrivateKey.PublicKey)
	auth, err := PrepareAuth(client, fromPrivateKey, fromAddr)
	if nil != err {
		return "", err
	}

	tx, err := tokenInstance.Transfer(auth, common.HexToAddress(toAddr), amount)
	if nil != err {
		return "", err
	}

	err = waitEthTxFinished(client, tx.Hash(), "TransferFromToken")
	if nil != err {
		return "", err
	}
	return tx.Hash().String(), nil
}

func LockEthErc20Asset(ownerPrivateKeyStr, tokenAddrStr, chain33Receiver string, amount *big.Int, client *ethclient.Client, bridgeBank *generated.BridgeBank) (string, error) {
	txslog.Info("LockEthErc20Asset", "ownerPrivateKeyStr", ownerPrivateKeyStr, "tokenAddrStr", tokenAddrStr, "chain33Receiver", chain33Receiver, "amount", amount.String())
	ownerPrivateKey, err := crypto.ToECDSA(common.FromHex(ownerPrivateKeyStr))
	if nil != err {
		return "", err
	}
	ownerAddr := crypto.PubkeyToAddress(ownerPrivateKey.PublicKey)

	auth, err := PrepareAuth(client, ownerPrivateKey, ownerAddr)
	if nil != err {
		txslog.Error("LockEthErc20Asset", "PrepareAuth err", err.Error())
		return "", err
	}
	//ETH转账，空地址，且设置value
	var tokenAddr common.Address
	if "" == tokenAddrStr {
		auth.Value = amount
	}

	if "" != tokenAddrStr {
		tokenAddr = common.HexToAddress(tokenAddrStr)
	}
	tx, err := bridgeBank.Lock(auth, []byte(chain33Receiver), tokenAddr, amount)
	if nil != err {
		txslog.Error("LockEthErc20Asset", "lock err", err.Error())
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "LockEthErc20Asset")
	if nil != err {
		txslog.Error("LockEthErc20Asset", "waitEthTxFinished err", err.Error())
		return "", err
	}

	return tx.Hash().String(), nil
}

/////////////////NewProphecyClaim////////////////
func MakeNewProphecyClaim(newProphecyClaimPara *NewProphecyClaimPara, client *ethclient.Client, privateKey *ecdsa.PrivateKey, transactor common.Address, x2EthContracts *X2EthContracts) (string, error) {

	authVali, err := PrepareAuth(client, privateKey, transactor)
	if nil != err {
		return "", err
	}

	amount := newProphecyClaimPara.Amount
	ethReceiver := newProphecyClaimPara.EthReceiver

	// Generate rawHash using ProphecyClaim data
	claimID := crypto.Keccak256Hash(newProphecyClaimPara.Txhash, newProphecyClaimPara.Chain33Sender, newProphecyClaimPara.EthReceiver.Bytes(), newProphecyClaimPara.TokenAddr.Bytes(), amount.Bytes())

	// Sign the hash using the active validator's private key
	signature, err := SignClaim4Eth(claimID, privateKey)
	if nil != err {
		return "", err
	}

	tx, err := x2EthContracts.Oracle.NewOracleClaim(authVali, newProphecyClaimPara.ClaimType, newProphecyClaimPara.Chain33Sender, ethReceiver, newProphecyClaimPara.TokenAddr, newProphecyClaimPara.Symbol, amount, claimID, signature)
	if nil != err {
		return "", err
	}
	err = waitEthTxFinished(client, tx.Hash(), "MakeNewProphecyClaim")
	if nil != err {
		return "", err
	}
	return tx.Hash().String(), nil
}
