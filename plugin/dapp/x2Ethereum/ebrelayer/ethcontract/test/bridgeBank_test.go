package test

import (
	"context"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/test/setup"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethtxs"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

//"BridgeToken creation (Chain33 assets)"
func TestBrigeTokenCreat(t *testing.T) {
	ctx := context.Background()
	println("TEST:BridgeToken creation (Chain33 assets)")
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
	auth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
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
				logEvent := &events.LogNewBridgeToken{}
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

//测试在chain33上锁定资产,然后在以太坊上铸币
//发行token="BTY"
//NewProphecyClaim
//铸币NewOracleClaim,
//ProcessBridgeProphecy
//铸币成功
//Bridge token minting (for locked chain33 assets)
func TestBrigeTokenMint(t *testing.T) {
	ctx := context.Background()
	println("TEST:BridgeToken creation (Chain33 assets)")
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
	auth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
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

	logEvent := &events.LogNewBridgeToken{}
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
	cosmosBridgeABI := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
	logNewProphecyClaimSig := cosmosBridgeABI.Events[eventName].ID().Hex()
	query = ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.Chain33Bridge.Address},
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
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(authVali, events.CLAIM_TYPE_LOCK, chain33Sender, ethReceiver, logEvent.Token, logEvent.Symbol, big.NewInt(amount))
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
			t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
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

//测试在以太坊上lock数字资产,包括Eth和Erc20
//Bridge deposit locking (deposit erc20/eth assets)
func TestBridgeDepositLock(t *testing.T) {
	ctx := context.Background()
	println("TEST:Bridge deposit locking (Erc20/Eth assets)")
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

	//创建token
	operatorAuth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
	symbol := "USDT"
	bridgeTokenAddr, _, bridgeTokenInstance, err := generated.DeployBridgeToken(operatorAuth, backend, symbol)
	require.Nil(t, err)
	sim.Commit()
	t.Logf("The new creaded symbol:%s, address:%s", symbol, bridgeTokenAddr.String())

	//创建实例
	//为userOne铸币
	//userOne为bridgebank允许allowance设置数额
	userOne := para.InitValidators[0]
	callopts := &bind.CallOpts{
		Pending: true,
		From:    userOne,
		Context: ctx,
	}
	symQuery, err := bridgeTokenInstance.Symbol(callopts)
	require.Equal(t, symQuery, symbol)
	t.Logf("symQuery = %s", symQuery)

	//isMiner, err := bridgeTokenInstance.IsMinter(callopts, x2EthDeployInfo.BridgeBank.Address)
	//require.Nil(t, err)
	//t.Logf("\nIsMinter for addr:%s, result:%v", x2EthDeployInfo.BridgeBank.Address.String(), isMiner)
	//require.Equal(t, isMiner, true)
	isMiner, err := bridgeTokenInstance.IsMinter(callopts, para.Operator)
	require.Nil(t, err)
	require.Equal(t, isMiner, true)

	operatorAuth, err = ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)

	mintAmount := int64(1000)
	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	_, err = bridgeTokenInstance.Mint(operatorAuth, userOne, big.NewInt(mintAmount))
	require.Nil(t, err)
	sim.Commit()
	userOneAuth, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	allowAmount := int64(100)
	_, err = bridgeTokenInstance.Approve(userOneAuth, x2EthDeployInfo.BridgeBank.Address, big.NewInt(allowAmount))
	require.Nil(t, err)
	sim.Commit()

	userOneBalance, err := bridgeTokenInstance.BalanceOf(callopts, userOne)
	require.Nil(t, err)
	t.Logf("userOneBalance:%d", userOneBalance.Int64())
	require.Equal(t, userOneBalance.Int64(), mintAmount)

	//***测试子项目:should allow users to lock ERC20 tokens
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)

	//lock 100
	lockAmount := big.NewInt(100)
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, bridgeTokenAddr, lockAmount)
	require.Nil(t, err)
	sim.Commit()

	//balance减少到900
	userOneBalance, err = bridgeTokenInstance.BalanceOf(callopts, userOne)
	require.Nil(t, err)
	expectAmount := int64(900)
	require.Equal(t, userOneBalance.Int64(), expectAmount)
	t.Logf("userOneBalance changes to:%d", userOneBalance.Int64())

	//bridgebank增加了100
	bridgeBankBalance, err := bridgeTokenInstance.BalanceOf(callopts, x2EthDeployInfo.BridgeBank.Address)
	require.Nil(t, err)
	expectAmount = int64(100)
	require.Equal(t, bridgeBankBalance.Int64(), expectAmount)
	t.Logf("bridgeBankBalance changes to:%d", bridgeBankBalance.Int64())

	//***测试子项目:should allow users to lock Ethereum
	bridgeBankBalance, err = sim.BalanceAt(ctx, x2EthDeployInfo.BridgeBank.Address, nil)
	require.Nil(t, err)
	t.Logf("origin eth bridgeBankBalance is:%d", bridgeBankBalance.Int64())

	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	ethAmount := big.NewInt(50)
	userOneAuth.Value = ethAmount

	//lock 50 eth
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, common.Address{}, ethAmount)
	require.Nil(t, err)
	sim.Commit()

	bridgeBankBalance, err = sim.BalanceAt(ctx, x2EthDeployInfo.BridgeBank.Address, nil)
	require.Nil(t, err)
	require.Equal(t, bridgeBankBalance.Int64(), ethAmount.Int64())
	t.Logf("eth bridgeBankBalance changes to:%d", bridgeBankBalance.Int64())
}

//测试在以太坊上unlock数字资产,包括Eth和Erc20
//Ethereum/ERC20 token unlocking (for burned chain33 assets)
func TestBridgeBankUnlock(t *testing.T) {
	ctx := context.Background()
	println("TEST:Ethereum/ERC20 token unlocking (for burned chain33 assets)")
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

	//1.lockEth资产
	ethAddr := common.Address{}
	ethToken, err := generated.NewBridgeToken(ethAddr, backend)
	userOneAuth, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	userOneAuth.Value = big.NewInt(300)
	_, err = ethToken.Transfer(userOneAuth, x2EthDeployInfo.BridgeBank.Address, userOneAuth.Value)
	sim.Commit()

	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	ethLockAmount := big.NewInt(150)
	userOneAuth.Value = ethLockAmount

	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	//lock 150 eth
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, common.Address{}, ethLockAmount)
	require.Nil(t, err)
	sim.Commit()

	//2.lockErc20资产
	//创建token
	operatorAuth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
	symbol_usdt := "USDT"
	bridgeTokenAddr, _, bridgeTokenInstance, err := generated.DeployBridgeToken(operatorAuth, backend, symbol_usdt)
	require.Nil(t, err)
	sim.Commit()
	t.Logf("The new creaded symbol_usdt:%s, address:%s", symbol_usdt, bridgeTokenAddr.String())

	//创建实例
	//为userOne铸币
	//userOne为bridgebank允许allowance设置数额
	userOne := para.InitValidators[0]
	callopts := &bind.CallOpts{
		Pending: true,
		From:    userOne,
		Context: ctx,
	}
	symQuery, err := bridgeTokenInstance.Symbol(callopts)
	require.Equal(t, symQuery, symbol_usdt)
	t.Logf("symQuery = %s", symQuery)

	//isMiner, err := bridgeTokenInstance.IsMinter(callopts, x2EthDeployInfo.BridgeBank.Address)
	//require.Nil(t, err)
	//t.Logf("\nIsMinter for addr:%s, result:%v", x2EthDeployInfo.BridgeBank.Address.String(), isMiner)
	//require.Equal(t, isMiner, true)
	isMiner, err := bridgeTokenInstance.IsMinter(callopts, para.Operator)
	require.Nil(t, err)
	require.Equal(t, isMiner, true)

	operatorAuth, err = ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)

	mintAmount := int64(1000)
	_, err = bridgeTokenInstance.Mint(operatorAuth, userOne, big.NewInt(mintAmount))
	require.Nil(t, err)
	sim.Commit()
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	allowAmount := int64(100)
	_, err = bridgeTokenInstance.Approve(userOneAuth, x2EthDeployInfo.BridgeBank.Address, big.NewInt(allowAmount))
	require.Nil(t, err)
	sim.Commit()

	userOneBalance, err := bridgeTokenInstance.BalanceOf(callopts, userOne)
	require.Nil(t, err)
	t.Logf("userOneBalance:%d", userOneBalance.Int64())
	require.Equal(t, userOneBalance.Int64(), mintAmount)

	//***测试子项目:should allow users to lock ERC20 tokens
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)

	//lock 100
	lockAmount := big.NewInt(100)
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, bridgeTokenAddr, lockAmount)
	require.Nil(t, err)
	sim.Commit()
	////////////////////////////////////////
	///////////////准备阶段结束///////////////
	////////////////////////////////////////

	////////////////////////////////////////
	///////////////开始测试//////////////////
	///3.should unlock Ethereum upon the processing of a burn prophecy///
	////////////////////////////////////////
	////////////////////订阅LogNewProphecyClaim/////////////////
	//：订阅事件LogNewProphecyClaim
	eventName := "LogNewProphecyClaim"
	cosmosBridgeABI := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
	logNewProphecyClaimSig := cosmosBridgeABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.Chain33Bridge.Address},
	}
	// We will check logs for new events
	newProphecyClaimlogs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	newProphecyClaimSub, err := sim.SubscribeFilterLogs(ctx, query, newProphecyClaimlogs)
	require.Nil(t, err)

	//提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmount := int64(150)
	ethReceivent := para.InitValidators[2]
	ethSym := string("eth")
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		ethAddr,
		ethSym,
		big.NewInt(newProphecyAmount))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
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
			t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
			require.Equal(t, ethSym, newProphecyClaimEvent.Symbol)
			require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmount)
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

	userEthbalance, _ := sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userEthbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userEthbalanceAfter, _ := sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s balance=%d", ethReceivent.String(), userEthbalanceAfter.Int64())
	require.Equal(t, userEthbalance.Int64()+newProphecyAmount, userEthbalanceAfter.Int64())

	//////////////////////////////////////////////////////////////////
	///////should unlock and transfer ERC20 tokens upon the processing of a burn prophecy//////
	//////////////////////////////////////////////////////////////////
	//提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmount = int64(100)
	ethReceivent = para.InitValidators[2]
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		bridgeTokenAddr,
		symbol_usdt,
		big.NewInt(newProphecyAmount))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
	timer.Reset(5*time.Second)
	for {
		select {
		case <-timer.C:
			t.Fatal("timeout for NewProphecyClaimEvent")
		// Handle any errors
		case err := <-newProphecyClaimSub.Err():
			goto latter
			t.Fatalf("sub error:%s", err.Error())
		// vLog is raw event data
		case vLog := <-newProphecyClaimlogs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logNewProphecyClaimSig {
				t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s", eventName,
					vLog.BlockNumber, vLog.TxHash.Hex())

				err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
				require.Nil(t, err)
				t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
				t.Logf("symbol:%s, amount:%d", newProphecyClaimEvent.Symbol, newProphecyClaimEvent.Amount.Int64())
				require.Equal(t, symbol_usdt, newProphecyClaimEvent.Symbol)
				require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmount)
				goto latter
			}
		}
	}
latter:

	///////////newOracleClaim///////////////////////////
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	OracleClaim, err = ethtxs.ProphecyClaimToSignedOracleClaim(*newProphecyClaimEvent, para.ValidatorPriKey[0])
	require.Nil(t, err)

	_, err = x2EthContracts.Oracle.NewOracleClaim(authOracle, newProphecyClaimEvent.ProphecyID, OracleClaim.Message, OracleClaim.Signature)
	require.Nil(t, err)

	userUSDTbalance, err := bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userUSDTbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userUSDTbalanceAfter, err := bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s symbol:%s, balance=%d", ethReceivent.String(), newProphecyClaimEvent.Symbol, userUSDTbalanceAfter.Int64())
	require.Equal(t, userUSDTbalance.Int64()+newProphecyAmount, userUSDTbalanceAfter.Int64())
}

//测试在以太坊上多次unlock数字资产,包括Eth和Erc20
//Ethereum/ERC20 token second unlocking (for burned chain33 assets)
func TestBridgeBankSecondUnlockEth(t *testing.T) {
	ctx := context.Background()
	println("TEST:to be unlocked incrementally by successive burn prophecies (for burned chain33 assets)")
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

	//1.lockEth资产
	ethAddr := common.Address{}
	ethToken, err := generated.NewBridgeToken(ethAddr, backend)
	userOneAuth, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	userOneAuth.Value = big.NewInt(300)
	_, err = ethToken.Transfer(userOneAuth, x2EthDeployInfo.BridgeBank.Address, userOneAuth.Value)
	sim.Commit()

	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	ethLockAmount := big.NewInt(150)
	userOneAuth.Value = ethLockAmount

	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	//lock 150 eth
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, common.Address{}, ethLockAmount)
	require.Nil(t, err)
	sim.Commit()

	//2.lockErc20资产
	//创建token
	operatorAuth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
	symbol_usdt := "USDT"
	bridgeTokenAddr, _, bridgeTokenInstance, err := generated.DeployBridgeToken(operatorAuth, backend, symbol_usdt)
	require.Nil(t, err)
	sim.Commit()
	t.Logf("The new creaded symbol_usdt:%s, address:%s", symbol_usdt, bridgeTokenAddr.String())

	//创建实例
	//为userOne铸币
	//userOne为bridgebank允许allowance设置数额
	userOne := para.InitValidators[0]
	callopts := &bind.CallOpts{
		Pending: true,
		From:    userOne,
		Context: ctx,
	}
	symQuery, err := bridgeTokenInstance.Symbol(callopts)
	require.Equal(t, symQuery, symbol_usdt)
	t.Logf("symQuery = %s", symQuery)

	//isMiner, err := bridgeTokenInstance.IsMinter(callopts, x2EthDeployInfo.BridgeBank.Address)
	//require.Nil(t, err)
	//t.Logf("\nIsMinter for addr:%s, result:%v", x2EthDeployInfo.BridgeBank.Address.String(), isMiner)
	//require.Equal(t, isMiner, true)
	isMiner, err := bridgeTokenInstance.IsMinter(callopts, para.Operator)
	require.Nil(t, err)
	require.Equal(t, isMiner, true)

	operatorAuth, err = ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)

	mintAmount := int64(1000)
	_, err = bridgeTokenInstance.Mint(operatorAuth, userOne, big.NewInt(mintAmount))
	require.Nil(t, err)
	sim.Commit()
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	allowAmount := int64(100)
	_, err = bridgeTokenInstance.Approve(userOneAuth, x2EthDeployInfo.BridgeBank.Address, big.NewInt(allowAmount))
	require.Nil(t, err)
	sim.Commit()

	userOneBalance, err := bridgeTokenInstance.BalanceOf(callopts, userOne)
	require.Nil(t, err)
	t.Logf("userOneBalance:%d", userOneBalance.Int64())
	require.Equal(t, userOneBalance.Int64(), mintAmount)

	//***测试子项目:should allow users to lock ERC20 tokens
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)

	//lock 100
	lockAmount := big.NewInt(100)
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, bridgeTokenAddr, lockAmount)
	require.Nil(t, err)
	sim.Commit()
	////////////////////////////////////////
	///////////////准备阶段结束///////////////
	////////////////////////////////////////

	////////////////////////////////////////
	///////////////开始测试//////////////////
	///3.should unlock Ethereum upon the processing of a burn prophecy///
	////////////////////////////////////////
	////////////////////订阅LogNewProphecyClaim/////////////////
	//：订阅事件LogNewProphecyClaim
	eventName := "LogNewProphecyClaim"
	cosmosBridgeABI := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
	logNewProphecyClaimSig := cosmosBridgeABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.Chain33Bridge.Address},
	}
	// We will check logs for new events
	newProphecyClaimlogs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	newProphecyClaimSub, err := sim.SubscribeFilterLogs(ctx, query, newProphecyClaimlogs)
	require.Nil(t, err)

	//提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmount := int64(44)
	ethReceivent := para.InitValidators[2]
	ethSym := string("eth")
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		ethAddr,
		ethSym,
		big.NewInt(newProphecyAmount))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
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
			err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
			require.Nil(t, err)
			t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
			t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s, symbol:%s", eventName,
				vLog.BlockNumber, vLog.TxHash.Hex(), newProphecyClaimEvent.Symbol)
			require.Equal(t, ethSym, newProphecyClaimEvent.Symbol)
			require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmount)
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

	userEthbalance, _ := sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userEthbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userEthbalanceAfter, _ := sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s balance=%d", ethReceivent.String(), userEthbalanceAfter.Int64())
	require.Equal(t, userEthbalance.Int64()+newProphecyAmount, userEthbalanceAfter.Int64())

	//第二次提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmountSecond := int64(33)
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		ethAddr,
		ethSym,
		big.NewInt(newProphecyAmountSecond))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
	for {
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
				err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
				require.Nil(t, err)
				t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s, symbol:%s", eventName,
					vLog.BlockNumber, vLog.TxHash.Hex(), newProphecyClaimEvent.Symbol)
				t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
				require.Equal(t, ethSym, newProphecyClaimEvent.Symbol)
				require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmountSecond)
				goto latter
			}
		}
	}

	latter:
	///////////newOracleClaim///////////////////////////
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	OracleClaim, err = ethtxs.ProphecyClaimToSignedOracleClaim(*newProphecyClaimEvent, para.ValidatorPriKey[0])
	require.Nil(t, err)

	_, err = x2EthContracts.Oracle.NewOracleClaim(authOracle, newProphecyClaimEvent.ProphecyID, OracleClaim.Message, OracleClaim.Signature)
	require.Nil(t, err)

	userEthbalance, _ = sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userEthbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userEthbalanceAfter, _ = sim.BalanceAt(ctx, ethReceivent, nil)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s balance=%d", ethReceivent.String(), userEthbalanceAfter.Int64())
	require.Equal(t, userEthbalance.Int64()+newProphecyAmountSecond, userEthbalanceAfter.Int64())
}

//测试在以太坊上多次unlock数字资产Erc20
//Ethereum/ERC20 token unlocking (for burned chain33 assets)
func TestBridgeBankSedondUnlockErc20(t *testing.T) {
	ctx := context.Background()
	println("TEST:ERC20 to be unlocked incrementally by successive burn prophecies (for burned chain33 assets))")
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

	//1.lockEth资产
	ethAddr := common.Address{}
	ethToken, err := generated.NewBridgeToken(ethAddr, backend)
	userOneAuth, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	userOneAuth.Value = big.NewInt(300)
	_, err = ethToken.Transfer(userOneAuth, x2EthDeployInfo.BridgeBank.Address, userOneAuth.Value)
	sim.Commit()

	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	ethLockAmount := big.NewInt(150)
	userOneAuth.Value = ethLockAmount

	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	//lock 150 eth
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, common.Address{}, ethLockAmount)
	require.Nil(t, err)
	sim.Commit()

	//2.lockErc20资产
	//创建token
	operatorAuth, err := ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)
	symbol_usdt := "USDT"
	bridgeTokenAddr, _, bridgeTokenInstance, err := generated.DeployBridgeToken(operatorAuth, backend, symbol_usdt)
	require.Nil(t, err)
	sim.Commit()
	t.Logf("The new creaded symbol_usdt:%s, address:%s", symbol_usdt, bridgeTokenAddr.String())

	//创建实例
	//为userOne铸币
	//userOne为bridgebank允许allowance设置数额
	userOne := para.InitValidators[0]
	callopts := &bind.CallOpts{
		Pending: true,
		From:    userOne,
		Context: ctx,
	}
	symQuery, err := bridgeTokenInstance.Symbol(callopts)
	require.Equal(t, symQuery, symbol_usdt)
	t.Logf("symQuery = %s", symQuery)
	isMiner, err := bridgeTokenInstance.IsMinter(callopts, para.Operator)
	require.Nil(t, err)
	require.Equal(t, isMiner, true)

	operatorAuth, err = ethtxs.PrepareAuth(backend, para.DeployPrivateKey, para.Operator)

	mintAmount := int64(1000)
	_, err = bridgeTokenInstance.Mint(operatorAuth, userOne, big.NewInt(mintAmount))
	require.Nil(t, err)
	sim.Commit()
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	allowAmount := int64(100)
	_, err = bridgeTokenInstance.Approve(userOneAuth, x2EthDeployInfo.BridgeBank.Address, big.NewInt(allowAmount))
	require.Nil(t, err)
	sim.Commit()

	userOneBalance, err := bridgeTokenInstance.BalanceOf(callopts, userOne)
	require.Nil(t, err)
	t.Logf("userOneBalance:%d", userOneBalance.Int64())
	require.Equal(t, userOneBalance.Int64(), mintAmount)

	//***测试子项目:should allow users to lock ERC20 tokens
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)

	//lock 100
	lockAmount := big.NewInt(100)
	_, err = x2EthContracts.BridgeBank.Lock(userOneAuth, chain33Sender, bridgeTokenAddr, lockAmount)
	require.Nil(t, err)
	sim.Commit()
	////////////////////////////////////////
	///////////////准备阶段结束///////////////
	////////////////////////////////////////

	////////////////////////////////////////
	///////////////开始测试//////////////////
	///3.should unlock Ethereum upon the processing of a burn prophecy///
	////////////////////////////////////////
	////////////////////订阅LogNewProphecyClaim/////////////////
	//：订阅事件LogNewProphecyClaim
	eventName := "LogNewProphecyClaim"
	cosmosBridgeABI := ethtxs.LoadABI(ethtxs.Chain33BridgeABI)
	logNewProphecyClaimSig := cosmosBridgeABI.Events[eventName].ID().Hex()
	query := ethereum.FilterQuery{
		Addresses: []common.Address{x2EthDeployInfo.Chain33Bridge.Address},
	}
	// We will check logs for new events
	newProphecyClaimlogs := make(chan types.Log)
	// Filter by contract and event, write results to logs
	newProphecyClaimSub, err := sim.SubscribeFilterLogs(ctx, query, newProphecyClaimlogs)
	require.Nil(t, err)

	//提交NewProphecyClaim

	//////////////////////////////////////////////////////////////////
	///////should unlock and transfer ERC20 tokens upon the processing of a burn prophecy//////
	//////////////////////////////////////////////////////////////////
	//提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmount := int64(33)
	ethReceivent := para.InitValidators[2]
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		bridgeTokenAddr,
		symbol_usdt,
		big.NewInt(newProphecyAmount))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
	newProphecyClaimEvent := &events.NewProphecyClaimEvent{}
	timer := time.NewTimer(5*time.Second)
	for {
		select {
		case <-timer.C:
			t.Fatal("timeout for NewProphecyClaimEvent")
		// Handle any errors
		case err := <-newProphecyClaimSub.Err():
			goto latter
			t.Fatalf("sub error:%s", err.Error())
		// vLog is raw event data
		case vLog := <-newProphecyClaimlogs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logNewProphecyClaimSig {
				t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s", eventName,
					vLog.BlockNumber, vLog.TxHash.Hex())

				err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
				require.Nil(t, err)
				t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
				t.Logf("symbol:%s, amount:%d", newProphecyClaimEvent.Symbol, newProphecyClaimEvent.Amount.Int64())
				require.Equal(t, symbol_usdt, newProphecyClaimEvent.Symbol)
				require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmount)
				goto latter
			}
		}
	}
latter:

	///////////newOracleClaim///////////////////////////
	authOracle, err := ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	OracleClaim, err := ethtxs.ProphecyClaimToSignedOracleClaim(*newProphecyClaimEvent, para.ValidatorPriKey[0])
	require.Nil(t, err)

	_, err = x2EthContracts.Oracle.NewOracleClaim(authOracle, newProphecyClaimEvent.ProphecyID, OracleClaim.Message, OracleClaim.Signature)
	require.Nil(t, err)

	userUSDTbalance, err := bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userUSDTbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userUSDTbalanceAfter, err := bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s symbol:%s, balance=%d", ethReceivent.String(), newProphecyClaimEvent.Symbol, userUSDTbalanceAfter.Int64())
	require.Equal(t, userUSDTbalance.Int64()+newProphecyAmount, userUSDTbalanceAfter.Int64())

	//再次提交NewProphecyClaim
	userOneAuth, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	newProphecyAmountSecond := int64(66)
	_, err = x2EthContracts.Chain33Bridge.NewProphecyClaim(
		userOneAuth,
		events.CLAIM_TYPE_BURN,
		chain33Sender,
		ethReceivent,
		bridgeTokenAddr,
		symbol_usdt,
		big.NewInt(newProphecyAmountSecond))
	sim.Commit()
	require.Nil(t, err)

	//接收并处理NewProphecyClaim
	for {
		select {
		case <-timer.C:
			t.Fatal("timeout for NewProphecyClaimEvent")
		// Handle any errors
		case err := <-newProphecyClaimSub.Err():
			goto latter
			t.Fatalf("sub error:%s", err.Error())
		// vLog is raw event data
		case vLog := <-newProphecyClaimlogs:
			// Check if the event is a 'LogLock' event
			if vLog.Topics[0].Hex() == logNewProphecyClaimSig {
				t.Logf("Witnessed new logNewProphecyClaim event:%s, Block number:%d, Tx hash:%s", eventName,
					vLog.BlockNumber, vLog.TxHash.Hex())

				err = cosmosBridgeABI.Unpack(newProphecyClaimEvent, eventName, vLog.Data)
				require.Nil(t, err)
				t.Logf("chain33Sender:%s", string(newProphecyClaimEvent.Chain33Sender))
				t.Logf("symbol:%s, amount:%d", newProphecyClaimEvent.Symbol, newProphecyClaimEvent.Amount.Int64())
				require.Equal(t, symbol_usdt, newProphecyClaimEvent.Symbol)
				require.Equal(t, newProphecyClaimEvent.Amount.Int64(), newProphecyAmountSecond)
				goto latter2
			}
		}
	}
latter2:

	///////////newOracleClaim///////////////////////////
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	OracleClaim, err = ethtxs.ProphecyClaimToSignedOracleClaim(*newProphecyClaimEvent, para.ValidatorPriKey[0])
	require.Nil(t, err)

	_, err = x2EthContracts.Oracle.NewOracleClaim(authOracle, newProphecyClaimEvent.ProphecyID, OracleClaim.Message, OracleClaim.Signature)
	require.Nil(t, err)

	userUSDTbalance, err = bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance for addr:%s balance=%d", ethReceivent.String(), userUSDTbalance.Int64())

	////////////////should mint bridge tokens upon the successful processing of a burn prophecy claim
	authOracle, err = ethtxs.PrepareAuth(backend, para.ValidatorPriKey[0], para.InitValidators[0])
	require.Nil(t, err)
	_, err = x2EthContracts.Oracle.ProcessBridgeProphecy(authOracle, newProphecyClaimEvent.ProphecyID)
	require.Nil(t, err)
	sim.Commit()
	userUSDTbalanceAfter, err = bridgeTokenInstance.BalanceOf(callopts, ethReceivent)
	t.Logf("userEthbalance after ProcessBridgeProphecy for addr:%s symbol:%s, balance=%d", ethReceivent.String(), newProphecyClaimEvent.Symbol, userUSDTbalanceAfter.Int64())
	require.Equal(t, userUSDTbalance.Int64()+newProphecyAmountSecond, userUSDTbalanceAfter.Int64())
}

