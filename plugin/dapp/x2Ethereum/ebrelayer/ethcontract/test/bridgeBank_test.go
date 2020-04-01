package test

import (
	"context"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"math/big"
	"time"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/test/setup"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethtxs"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type LogNewBridgeToken struct {
	Token  common.Address
	Symbol string
}

const (
	CLAIM_TYPE_BURN = uint8(1)
	CLAIM_TYPE_LOCK = uint8(2)
)

//"BridgeToken creation (Cosmos assets)"
func TestBrigeTokenCreat(t *testing.T) {
	ctx := context.Background()
	println("TEST:BridgeToken creation (Cosmos assets)")
	//1st部署相关合约
	backend, para := setup.PrepareTestEnv()
	sim := backend.(*backends.SimulatedBackend)

	balance, _ := sim.BalanceAt(ctx, para.Deployer, nil)
	fmt.Println("deployer addr,", para.Deployer.String(), "balance =", balance.String())

	/////////////////////////EstimateGas///////////////////////////
	callMsg := ethereum.CallMsg{
		From: para.Deployer,
		Data: common.FromHex(generated.BridgeBankBin),
	}

	gas, err := sim.EstimateGas(ctx, callMsg)
	if nil != err {
		panic("failed to estimate gas due to:" + err.Error())
	}
	fmt.Printf("\nThe estimated gas=%d", gas)
	////////////////////////////////////////////////////

	x2EthContracts, x2EthDeployInfo, err := ethtxs.DeployAndInit(backend, para)
	if nil != err {
		t.Fatalf("DeployAndInit failed due to:%s", err.Error())
	}
	sim.Commit()

	//2nd：订阅事件
	eventName := "LogNewBridgeToken"
	bridgeBankABI := ethtxs.LoadABI(ethtxs.BridgeBankABI)
	logNewBridgeTokenSig := bridgeBankABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.BridgeBank.Address},
	}
	// We will check logs for new events
	logs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	sub, err := sim.SubscribeFilterLogs(ctx, query, logs)
	require.Nil(t, err)

	//fmt.Printf("\n*****BridgeBank addr:%s, BridgeBank:%v***\n\n", deployInfo.BridgeBank.Address.String(), x2EthContracts.BridgeBank)
	t.Logf("x2EthDeployInfo.BridgeBank.Address is:%s", x2EthDeployInfo.BridgeBank.Address.String())
	bridgeBank, err := generated.NewBridgeBank(x2EthDeployInfo.BridgeBank.Address, backend)
	require.Nil(t, err)

	opts := &bind.CallOpts{
		Pending: true,
		From:    para.Operator,
		Context: ctx,
	}
	BridgeBankAddr, err := x2EthContracts.BridgeRegistry.BridgeBank(opts)
	require.Nil(t, err)
	t.Logf("BridgeBankAddr is:%s", BridgeBankAddr.String())

	//tokenCount, err := x2EthContracts.BridgeBank.BridgeTokenCount(opts)
	tokenCount, err := bridgeBank.BridgeBankCaller.BridgeTokenCount(opts)
	require.Nil(t, err)
	require.Equal(t, tokenCount.Int64(), int64(0))

	//3rd：创建token
	auth, err := ethtxs.PrepareAuth(backend, para.PrivateKey, para.Operator)
	if nil != err {
		t.Fatalf("PrepareAuth failed due to:%s", err.Error())
	}
	symbol := "BTY"
	_, err = x2EthContracts.BridgeBank.BridgeBankTransactor.CreateNewBridgeToken(auth, symbol)
	if nil != err {
		t.Fatalf("CreateNewBridgeToken failed due to:%s", err.Error())
	}
	sim.Commit()

	timer := time.NewTimer(30 * time.Second)
	for {
		select {
		case <-timer.C:
			t.Fatal("failed due to timeout")
		// Handle any errors
		case err := <-sub.Err():
			t.Fatalf("sub error:%s", err.Error())
		// vLog is raw event data
		case vLog := <-logs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logNewBridgeTokenSig {
				t.Logf("Witnessed new event:%s, Block number:%d, Tx hash:%s", eventName,
					vLog.BlockNumber, vLog.TxHash.Hex())
				logEvent := &LogNewBridgeToken{}
				err = bridgeBankABI.Unpack(logEvent, eventName, vLog.Data)
				require.Nil(t, err)
				t.Logf("token addr:%s, symbol:%s", logEvent.Token.String(), logEvent.Symbol)
				require.Equal(t, symbol, logEvent.Symbol)

				//tokenCount正确加1
				tokenCount, err := x2EthContracts.BridgeBank.BridgeTokenCount(opts)
				require.Nil(t, err)
				require.Equal(t, tokenCount.Int64(), int64(1))

				return
			}
		}
	}
}

//Bridge token minting (for locked Cosmos assets)
func TestBrigeTokenMint(t *testing.T) {
	ctx := context.Background()
	println("TEST:BridgeToken creation (Cosmos assets)")
	//1st部署相关合约
	backend, para := setup.PrepareTestEnv()
	sim := backend.(*backends.SimulatedBackend)

	balance, _ := sim.BalanceAt(ctx, para.Deployer, nil)
	fmt.Println("deployer addr,", para.Deployer.String(), "balance =", balance.String())

	/////////////////////////EstimateGas///////////////////////////
	callMsg := ethereum.CallMsg{
		From: para.Deployer,
		Data: common.FromHex(generated.BridgeBankBin),
	}

	gas, err := sim.EstimateGas(ctx, callMsg)
	if nil != err {
		panic("failed to estimate gas due to:" + err.Error())
	}
	fmt.Printf("\nThe estimated gas=%d", gas)
	////////////////////////////////////////////////////

	x2EthContracts, x2EthDeployInfo, err := ethtxs.DeployAndInit(backend, para)
	if nil != err {
		t.Fatalf("DeployAndInit failed due to:%s", err.Error())
	}
	sim.Commit()
	auth, err := ethtxs.PrepareAuth(backend, para.PrivateKey, para.Operator)
	if nil != err {
		t.Fatalf("PrepareAuth failed due to:%s", err.Error())
	}

	//2nd：订阅事件
	eventName := "LogNewBridgeToken"
	bridgeBankABI := ethtxs.LoadABI(ethtxs.BridgeBankABI)
	logNewBridgeTokenSig := bridgeBankABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.BridgeBank.Address},
	}
	// We will check logs for new events
	logs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	sub, err := sim.SubscribeFilterLogs(ctx, query, logs)
	require.Nil(t, err)

	opts := &bind.CallOpts{
		Pending: true,
		From:    para.Operator,
		Context: ctx,
	}

	tokenCount, err := x2EthContracts.BridgeBank.BridgeTokenCount(opts)
	require.Equal(t, tokenCount.Int64(), int64(0))

	//3rd：创建token
	symbol := "BTY"
	_, err = x2EthContracts.BridgeBank.BridgeBankTransactor.CreateNewBridgeToken(auth, symbol)
	if nil != err {
		t.Fatalf("CreateNewBridgeToken failed due to:%s", err.Error())
	}
	sim.Commit()

	logEvent := &LogNewBridgeToken{}
	select {
	// Handle any errors
	case err := <-sub.Err():
		t.Fatalf("sub error:%s", err.Error())
	// vLog is raw event data
	case vLog := <-logs:
		// Check if the event is a 'LogLock' event
		if vLog.Topics[0].Hex() == logNewBridgeTokenSig {
			t.Logf("Witnessed new event:%s, Block number:%d, Tx hash:%s", eventName,
				vLog.BlockNumber, vLog.TxHash.Hex())

			err = bridgeBankABI.Unpack(logEvent, eventName, vLog.Data)
			require.Nil(t, err)
			t.Logf("token addr:%s, symbol:%s", logEvent.Token.String(), logEvent.Symbol)
			require.Equal(t, symbol, logEvent.Symbol)

			//tokenCount正确加1
			tokenCount, err = x2EthContracts.BridgeBank.BridgeTokenCount(opts)
			require.Equal(t, tokenCount.Int64(), int64(1))
			break
		}
	}

	////////////////////订阅LogNewProphecyClaim/////////////////
	//：订阅事件LogNewProphecyClaim
	eventName = "LogNewProphecyClaim"
	cosmosBridgeABI := ethtxs.LoadABI(ethtxs.CosmosBridgeABI)
	logNewProphecyClaimSig := cosmosBridgeABI.Events[eventName].ID().Hex()
	query = ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.CosmosBridge.Address},
	}
	// We will check logs for new events
	newProphecyClaimlogs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	newProphecyClaimSub, err := sim.SubscribeFilterLogs(ctx, query, newProphecyClaimlogs)
	require.Nil(t, err)

	/////////////////NewProphecyClaim////////////////
	balance, _ = sim.BalanceAt(ctx, para.InitValidators[0], nil)
	fmt.Println("InitValidators[0] addr,", para.InitValidators[0].String(), "balance =", balance.String())

	authVali, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	amount := int64(99)
	ethReceiver := para.InitValidators[2]
	_, err = x2EthContracts.CosmosBridge.NewProphecyClaim(authVali, CLAIM_TYPE_LOCK, chain33Sender, ethReceiver, logEvent.Token, logEvent.Symbol, big.NewInt(amount))
	sim.Commit()
	require.Nil(t, err)

	newProphecyClaimEvent := &events.NewProphecyClaimEvent{}
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		t.Fatal("timeout for NewProphecyClaimEvent")
	// Handle any errors
	case err := <-newProphecyClaimSub.Err():
		t.Fatalf("sub error:%s", err.Error())
	// vLog is raw event data
	case vLog := <-newProphecyClaimlogs:
		// Check if the event is a 'LogLock' event
		if vLog.Topics[0].Hex() == logNewProphecyClaimSig {
			t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s", eventName,
				vLog.BlockNumber, vLog.TxHash.Hex())

			err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
			require.Nil(t, err)
			t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.CosmosSender))
			require.Equal(t, symbol, newProphecyClaimEvent.Symbol)
			require.Equal(t, newProphecyClaimEvent.Amount.Int64(), amount)
			break
		}
	}

	///////////newOracleClaim///////////////////////////
	authOracle, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	OracleClaim, err := ethtxs.ProphecyClaimToSignedOracleClaim(*newProphecyClaimEvent, para.ValidatorPriKey[0])
	require.Nil(t, err)

	_, err = x2EthContracts.Oracle.NewOracleClaim(authOracle, newProphecyClaimEvent.ProphecyID, OracleClaim.Message, OracleClaim.Signature)
	require.Nil(t, err)

	bridgeToken, err := generated.NewBridgeToken(logEvent.Token, backend)
	require.Nil(t, err)
	opts = &bind.CallOpts{
		Pending: true,
		Context: ctx,
	}
	balance, err = bridgeToken.BalanceOf(opts, ethReceiver)
	require.Nil(t, err)
	require.Equal(t, balance.Int64(), int64(0))

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()

	balance, err = bridgeToken.BalanceOf(opts, ethReceiver)
	require.Nil(t, err)
	require.Equal(t, balance.Int64(), amount)
	t.Logf("The minted amount is:%d", balance.Int64())

	////////////////
	ethMessageHash, err := x2EthContracts.Valset.DebugEthMessageHash(opts, OracleClaim.Message)
	require.Nil(t, err)
	t.Logf("The ethMessageHash is:%s", common.Bytes2Hex(ethMessageHash[:]))

	ethMessagePack, err := x2EthContracts.Valset.DebugPacked(opts, OracleClaim.Message)
	require.Nil(t, err)
	t.Logf("The ethMessagePack is:%s", common.Bytes2Hex(ethMessagePack[:]))

}
