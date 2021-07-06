package chain33

import (
	"context"
	"flag"
	"fmt"
	chain33Common "github.com/33cn/chain33/common"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/contracts4chain33/generated"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/test/setup"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/ethereum/ethinterface"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/ethereum/ethtxs"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/relayer/events"
	ebTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	relayerTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	tml "github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/33cn/chain33/system"

	// 需要显示引用系统插件，以加载系统内置合约
	"github.com/33cn/chain33/client/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	configPath    = flag.String("f", "./../../relayer.toml", "configfile")
	privateKeyStr = "0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"
	accountAddr   = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	passphrase    = "123456hzj"
	test          = "0ac3050aa3020a0a7832657468657265756d126d60671a690a2a3078303030303030303030303030303030303030303030303030303030303030303030303030303030301a2a307830633035626135633233306664616135303362353337303261663139363265303864306336306266220831303030303030302a0365746838121a6e080112210320bbac09528e19c55b0f89cb37ab265e7e856b1a8c388780322dbbfd194b52ba1a473045022100c403d9a6e531292336b44d52e4f4dbb9b8ab1e16335383954583728b909478da022031d8a29efcbcea8df648c4054f3c09ab1ab7a330797cf79fd891a3d9336922e920a08d0628e0f193f60530a1d7ad93e5ebc28e253a22314c7538586d537459765777664e716951336e4e4b33345239466648346b5270425612ce0208021a5e0802125a0a2b10c0d59294bb192222313271796f6361794e46374c7636433971573461767873324537553431664b536676122b10a0c88c94bb192222313271796f6361794e46374c7636433971573461767873324537553431664b5366761a55080f12510a291080ade2042222313271796f6361794e46374c7636433971573461767873324537553431664b53667612242222313271796f6361794e46374c7636433971573461767873324537553431664b5366761a92010867128d010a2a3078303030303030303030303030303030303030303030303030303030303030303030303030303030301222313271796f6361794e46374c7636433971573461767873324537553431664b5366761a2a307830633035626135633233306664616135303362353337303261663139363265303864306336306266220831303030303030302a03657468301220c4092a207a38e1da7de4444f2d34c7488293f3a2e01ce2561e720e9bbef355e83755ad833220e68d8418f69d5f18278a53dca53b101f26f76883337a60a5754d5f6d94e42e3c400148c409"
)

//func TestAll(t *testing.T) {
//	//mock33 := newMock33()
//	//defer mock33.Close()
//	//test_HandleRequest(t)
//	Test_ImportPrivateKey(t)
//	//test_Lockbty(t)
//	Test_RestorePrivateKeys(t)
//}

//func Test_HandleRequest(t *testing.T) {
//	mock33 := newMock33()
//	defer mock33.Close()
//	_, _, _, x2EthDeployInfo, err := deployContracts()
//	require.NoError(t, err)
//	chain33Relayer := newChain33Relayer(x2EthDeployInfo, "127.0.0.1:60002")
//	err = chain33Relayer.ImportPrivateKey(passphrase, privateKeyStr)
//	assert.NoError(t, err)
//
//	body, err := hex.DecodeString(test)
//	assert.NoError(t, err)
//
//	//chain33Relayer.statusCheckedIndex = 1220
//	err = syncTx.HandleRequest(body)
//	assert.NoError(t, err)
//
//	time.Sleep(200 * time.Millisecond)
//
//	ret := chain33Relayer.QueryTxhashRelay2Eth()
//	assert.NotEmpty(t, ret)
//
//	event := getOracleClaimType(events.MsgLock.String())
//	assert.Equal(t, event, events.Event(events.ClaimTypeLock))
//}

// getOracleClaimType : sets the OracleClaim's claim type based upon the witnessed event type
func getOracleClaimType(eventType string) events.Event {
	var claimType events.Event

	switch eventType {
	case events.MsgBurn.String():
		claimType = events.Event(events.ClaimTypeBurn)
	case events.MsgLock.String():
		claimType = events.Event(events.ClaimTypeLock)
	default:
		//panic(errors.New("eventType invalid"))
	}

	return claimType
}

func Test_ImportPrivateKey(t *testing.T) {
	mock33 := newMock33()
	defer mock33.Close()
	_, _, _, x2EthDeployInfo, err := setup.DeployContracts()
	require.NoError(t, err)
	chain33Relayer := newChain33Relayer(x2EthDeployInfo, "127.0.0.1:60000")

	err = chain33Relayer.ImportPrivateKey(passphrase, privateKeyStr)
	assert.NoError(t, err)
	//assert.Equal(t, addr, accountAddr)

	time.Sleep(50 * time.Millisecond)

	addr, err := chain33Relayer.GetAccountAddr()
	assert.NoError(t, err)
	assert.Equal(t, addr, accountAddr)

	key, _, _ := chain33Relayer.GetAccount("123")
	assert.NotEqual(t, key, privateKeyStr)

	key, _, _ = chain33Relayer.GetAccount(passphrase)
	assert.Equal(t, key, privateKeyStr)
}

//func Test_Lockbty(t *testing.T) {
//	mock33 := newMock33()
//	defer mock33.Close()
//	para, sim, x2EthContracts, x2EthDeployInfo, err := setup.DeployContracts()
//	require.NoError(t, err)
//	chain33Relayer := newChain33Relayer(x2EthDeployInfo, "127.0.0.1:60001")
//	err = chain33Relayer.ImportPrivateKey(passphrase, privateKeyStr)
//	assert.NoError(t, err)
//
//	ctx := context.Background()
//	//2nd：订阅事件
//	eventName := "LogNewBridgeToken"
//	bridgeBankABI := ethtxs.LoadABI(ethtxs.BridgeBankABI)
//	logNewBridgeTokenSig := bridgeBankABI.Events[eventName].ID.Hex()
//	query := ethereum.FilterQuery{
//		Addresses: []common.Address{x2EthDeployInfo.BridgeBank.Address},
//	}
//	// We will check logs for new events
//	logs := make(chan ethTypes.Log)
//	// Filter by contract and event, write results to logs
//	sub, err := sim.SubscribeFilterLogs(ctx, query, logs)
//	assert.Nil(t, err)
//	//require.Nil(t, err)
//
//	opts := &bind.CallOpts{
//		Pending: true,
//		From:    para.Operator,
//		Context: ctx,
//	}
//
//	tokenCount, err := x2EthContracts.BridgeBank.BridgeTokenCount(opts)
//	require.Nil(t, err)
//	assert.Equal(t, tokenCount.Int64(), int64(0))
//
//	//3rd：创建token
//	symbol := "BTY"
//	auth, err := ethtxs.PrepareAuth(sim, para.DeployPrivateKey, para.Operator)
//	require.Nil(t, err)
//	_, err = x2EthContracts.BridgeBank.BridgeBankTransactor.CreateNewBridgeToken(auth, symbol)
//	require.Nil(t, err)
//	sim.Commit()
//
//	logEvent := &events.LogNewBridgeToken{}
//	select {
//	// Handle any errors
//	case err := <-sub.Err():
//		t.Fatalf("sub error:%s", err.Error())
//	// vLog is raw event data
//	case vLog := <-logs:
//		// Check if the event is a 'LogLock' event
//		if vLog.Topics[0].Hex() == logNewBridgeTokenSig {
//			_, err = bridgeBankABI.Unpack(eventName, vLog.Data)
//			require.Nil(t, err)
//			fmt.Println("symbol", symbol, "logEvent.Symbol", logEvent.Symbol)
//			require.Equal(t, symbol, logEvent.Symbol)
//
//			//tokenCount正确加1
//			tokenCount, err = x2EthContracts.BridgeBank.BridgeTokenCount(opts)
//			require.Nil(t, err)
//			require.Equal(t, tokenCount.Int64(), int64(1))
//			break
//		}
//	}
//
//	///////////newOracleClaim///////////////////////////
//	balance, _ := sim.BalanceAt(ctx, para.InitValidators[0], nil)
//	fmt.Println("InitValidators[0] addr,", para.InitValidators[0].String(), "balance =", balance.String())
//
//	chain33Sender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
//	amount := int64(99)
//	ethReceiver := para.InitValidators[2]
//	claimID := crypto.Keccak256Hash(chain33Sender, ethReceiver.Bytes(), logEvent.Token.Bytes(), big.NewInt(amount).Bytes())
//
//	authOracle, err := ethtxs.PrepareAuth(sim, para.ValidatorPriKey[0], para.InitValidators[0])
//	require.Nil(t, err)
//
//	signature, err := utils.SignClaim4Evm(claimID, para.ValidatorPriKey[0])
//	//signature, err := ethtxs.SignClaim4Eth(claimID, para.ValidatorPriKey[0])
//	require.Nil(t, err)
//
//	bridgeToken, err := generated.NewBridgeToken(logEvent.Token, sim)
//	require.Nil(t, err)
//	opts = &bind.CallOpts{
//		Pending: true,
//		Context: ctx,
//	}
//
//	balance, err = bridgeToken.BalanceOf(opts, ethReceiver)
//	require.Nil(t, err)
//	require.Equal(t, balance.Int64(), int64(0))
//
//	tx, err := x2EthContracts.Oracle.NewOracleClaim(
//		authOracle,
//		uint8(events.ClaimTypeLock),
//		chain33Sender,
//		ethReceiver,
//		logEvent.Token,
//		logEvent.Symbol,
//		big.NewInt(amount),
//		claimID,
//		signature)
//	require.Nil(t, err)
//
//	sim.Commit()
//	balance, err = bridgeToken.BalanceOf(opts, ethReceiver)
//	require.Nil(t, err)
//	require.Equal(t, balance.Int64(), amount)
//	//t.Logf("The minted amount is:%d", balance.Int64())
//
//	txhash := tx.Hash().Hex()
//	fmt.Println(txhash)
//
//	chain33Relayer.rwLock.Lock()
//	//chain33Relayer.statusCheckedIndex = 1
//	chain33Relayer.totalTx4Chain33ToEth = 2
//	chain33Relayer.rwLock.Unlock()
//	//_ = chain33Relayer.setLastestRelay2Chain33TxStatics(relayerTx.EthTxPending.String(), 2, txhash)
//
//	time.Sleep(200 * time.Millisecond)
//
//	chain33Relayer.rwLock.Lock()
//	//chain33Relayer.statusCheckedIndex = 9
//	chain33Relayer.totalTx4Chain33ToEth = 11
//	chain33Relayer.rwLock.Unlock()
//	//_ = chain33Relayer.setLastestRelay2EthTxhash(relayerTx.EthTxPending.String(), "", 11)
//
//	time.Sleep(200 * time.Millisecond)
//}

func Test_RestorePrivateKeys(t *testing.T) {
	mock33 := newMock33()
	defer mock33.Close()
	_, _, _, x2EthDeployInfo, err := setup.DeployContracts()
	require.NoError(t, err)
	chain33Relayer := newChain33Relayer(x2EthDeployInfo, "127.0.0.1:60003")
	err = chain33Relayer.ImportPrivateKey(passphrase, privateKeyStr)
	assert.NoError(t, err)

	go func() {
		for range chain33Relayer.unlockChan {
		}
	}()

	err = chain33Relayer.RestorePrivateKeys("123")
	assert.NotEqual(t, chain33Common.ToHex(chain33Relayer.privateKey4Chain33.Bytes()), privateKeyStr)
	fmt.Println("err", err)
	assert.NoError(t, err)

	err = chain33Relayer.RestorePrivateKeys(passphrase)
	assert.Equal(t, chain33Common.ToHex(chain33Relayer.privateKey4Chain33.Bytes()), privateKeyStr)
	assert.NoError(t, err)

	err = chain33Relayer.StoreAccountWithNewPassphase("new123", passphrase)
	assert.NoError(t, err)

	err = chain33Relayer.RestorePrivateKeys("new123")
	assert.Equal(t, chain33Common.ToHex(chain33Relayer.privateKey4Chain33.Bytes()), privateKeyStr)
	assert.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
}

func newChain33Relayer(x2EthDeployInfo *ethtxs.X2EthDeployInfo, pushBind string) *Relayer4Chain33 {
	cfg := initCfg(*configPath)
	cfg.SyncTxConfig.Chain33Host = "http://127.0.0.1:8801"
	cfg.BridgeRegistry = x2EthDeployInfo.BridgeRegistry.Address.String()
	cfg.SyncTxConfig.PushBind = pushBind
	cfg.SyncTxConfig.FetchHeightPeriodMs = 50
	cfg.SyncTxConfig.Dbdriver = "memdb"

	db := dbm.NewDB("relayer_db_service", cfg.SyncTxConfig.Dbdriver, cfg.SyncTxConfig.DbPath, cfg.SyncTxConfig.DbCache)
	ctx, cancel := context.WithCancel(context.Background())
	ethBridgeClaimchan := make(chan *relayerTypes.EthBridgeClaim, 100)
	chain33Msgchan := make(chan *events.Chain33Msg, 100)

	var wg sync.WaitGroup

	relayer := &Relayer4Chain33{
		rpcLaddr:             cfg.SyncTxConfig.Chain33Host,
		fetchHeightPeriodMs:  cfg.SyncTxConfig.FetchHeightPeriodMs,
		db:                   db,
		ctx:                  ctx,
		bridgeRegistryAddr:   x2EthDeployInfo.BridgeRegistry.Address.String(),
		chainName:            "",
		chainID:              0,
		unlockChan:           make(chan int),
		deployInfo:           cfg.Deploy,
		ethBridgeClaimChan:   ethBridgeClaimchan,
		chain33MsgChan:       chain33Msgchan,
		totalTx4Chain33ToEth: 0,
		symbol2Addr:          make(map[string]string),
		oracleAddr:           x2EthDeployInfo.Oracle.Address.String(),
		bridgeBankAddr:       x2EthDeployInfo.BridgeBank.Address.String(),
	}

	//err := relayer.setStatusCheckedIndex(1)
	//assert.NoError(t, err)
	//relayer.totalTx4Chain33ToEth = relayer.getTotalTxAmount2Eth()
	//relayer.statusCheckedIndex = relayer.getStatusCheckedIndex()
	//assert.Equal(t, relayer.statusCheckedIndex, int64(1))

	syncCfg := &ebTypes.SyncTxReceiptConfig{
		Chain33Host:       cfg.SyncTxConfig.Chain33Host,
		PushHost:          cfg.SyncTxConfig.PushHost,
		PushName:          cfg.SyncTxConfig.PushName,
		PushBind:          pushBind,
		StartSyncHeight:   cfg.SyncTxConfig.StartSyncHeight,
		StartSyncSequence: cfg.SyncTxConfig.StartSyncSequence,
		StartSyncHash:     cfg.SyncTxConfig.StartSyncHash,
	}
	go relayer.syncProc(syncCfg)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
		wg.Wait()
		os.Exit(0)
	}()
	return relayer
}

func deployContracts() (*ethtxs.DeployPara, *ethinterface.SimExtend, *ethtxs.X2EthContracts, *ethtxs.X2EthDeployInfo, error) {
	// 0x8AFDADFC88a1087c9A1D6c0F5Dd04634b87F303a
	deployerPrivateKey := "8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
	// 0x92C8b16aFD6d423652559C6E266cBE1c29Bfd84f
	ethValidatorAddrKeyA := "3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
	ethValidatorAddrKeyB := "a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
	ethValidatorAddrKeyC := "bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
	ethValidatorAddrKeyD := "c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"

	ethValidatorAddrKeys := make([]string, 0)
	ethValidatorAddrKeys = append(ethValidatorAddrKeys, ethValidatorAddrKeyA)
	ethValidatorAddrKeys = append(ethValidatorAddrKeys, ethValidatorAddrKeyB)
	ethValidatorAddrKeys = append(ethValidatorAddrKeys, ethValidatorAddrKeyC)
	ethValidatorAddrKeys = append(ethValidatorAddrKeys, ethValidatorAddrKeyD)

	ctx := context.Background()
	backend, para := setup.PrepareTestEnvironment(deployerPrivateKey, ethValidatorAddrKeys)
	sim := new(ethinterface.SimExtend)
	sim.SimulatedBackend = backend.(*backends.SimulatedBackend)

	opts, _ := bind.NewKeyedTransactorWithChainID(para.DeployPrivateKey, big.NewInt(1337))
	parsed, _ := abi.JSON(strings.NewReader(generated.BridgeBankBin))
	contractAddr, _, _, _ := bind.DeployContract(opts, parsed, common.FromHex(generated.BridgeBankBin), sim)
	sim.Commit()

	callMsg := ethereum.CallMsg{
		From: para.Deployer,
		To:   &contractAddr,
		Data: common.FromHex(generated.BridgeBankBin),
	}

	_, err := sim.EstimateGas(ctx, callMsg)
	if nil != err {
		panic("failed to estimate gas due to:" + err.Error())
	}
	x2EthContracts, x2EthDeployInfo, err := ethtxs.DeployAndInit(sim, para)
	if nil != err {
		fmt.Println(err.Error())
		return nil, nil, nil, nil, err
	}
	sim.Commit()

	return para, sim, x2EthContracts, x2EthDeployInfo, nil
}

func initCfg(path string) *relayerTypes.RelayerConfig {
	var cfg relayerTypes.RelayerConfig
	if _, err := tml.DecodeFile(path, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return &cfg
}

//func init() {
//	fmt.Println("======================= init =======================")
//	var ret = types.ReplySubscribePush{IsOk: true, Msg: ""}
//	var he = types.Header{Height: 10000}
//
//	mockapi := &mocks.QueueProtocolAPI{}
//	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
//	mockapi.On("Close").Return()
//	mockapi.On("AddPushSubscribe", mock.Anything).Return(&ret, nil)
//	mockapi.On("GetLastHeader", mock.Anything).Return(&he, nil)
//
//	mock33 := testnode.New("", mockapi)
//	defer mock33.Close()
//	rpcCfg := mock33.GetCfg().RPC
//	// 这里必须设置监听端口，默认的是无效值
//	rpcCfg.JrpcBindAddr = "127.0.0.1:8801"
//	mock33.GetRPC().Listen()
//}

func newMock33() *testnode.Chain33Mock {
	var ret = types.ReplySubscribePush{IsOk: true, Msg: ""}
	var he = types.Header{Height: 10000}

	mockapi := &mocks.QueueProtocolAPI{}
	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
	mockapi.On("Close").Return()
	mockapi.On("AddPushSubscribe", mock.Anything).Return(&ret, nil)
	mockapi.On("GetLastHeader", mock.Anything).Return(&he, nil)

	mock33 := testnode.New("", mockapi)
	//defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	// 这里必须设置监听端口，默认的是无效值
	rpcCfg.JrpcBindAddr = "127.0.0.1:8801"
	mock33.GetRPC().Listen()

	return mock33
}

func Test_getExecerName(t *testing.T) {
	assert.Equal(t, getExecerName(""), "evm")
	assert.Equal(t, getExecerName("user.p.para."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p.para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user...p.para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p...para.."), "user.p.para.evm")
	assert.Equal(t, getExecerName("user.p.para"), "user.p.para.evm")
	assert.Equal(t, getExecerName("user"), "user.evm")
}
