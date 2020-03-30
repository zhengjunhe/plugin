package executor

//
//import (
//	"github.com/33cn/chain33/client"
//	"github.com/33cn/chain33/common"
//	"github.com/33cn/chain33/common/crypto"
//	dbm "github.com/33cn/chain33/common/db"
//	"github.com/33cn/chain33/common/log"
//	"github.com/33cn/chain33/queue"
//	"github.com/33cn/chain33/system/dapp"
//	"github.com/33cn/chain33/types"
//	"github.com/33cn/chain33/util"
//	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
//)
//
//var (
//	adminAddr   = "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
//	receiveAddr = "1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"
//
//	ethid                 = 0
//	bridgeContractAddress = "0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB"
//	nonce                 = 0
//	chain33Symbol         = "eth"
//	erhSymbol             = "eth"
//	cexec                 = "x2ethereum"
//	tokenContractAddress  = "0x0000000000000000000000000000000000000000"
//)
//
//func init() {
//	log.SetLogLevel("error")
//	Init(types2.X2ethereumX, testCfg, nil)
//}
//
//type testExecMock struct {
//	dbDir   string
//	localDB dbm.KVDB
//	stateDB dbm.DB
//	exec    dapp.Driver
//	cfg     *types.Chain33Config
//	q       queue.Queue
//	qapi    client.QueueProtocolAPI
//}
//
//type testcase struct {
//	payload            types.Message
//	expectExecErr      error
//	expectCheckErr     error
//	expectExecLocalErr error
//	expectExecDelErr   error
//	priv               string
//	index              int
//	systemCreate       bool
//	testState          int
//	testSign           []byte
//	testFee            int64
//}
//
//// InitEnv init env
//func (mock *testExecMock) InitEnv() {
//
//	mock.cfg = testCfg
//	util.ResetDatadir(mock.cfg.GetModuleConfig(), "$TEMP/")
//	mock.q = queue.New("channel")
//	mock.q.SetConfig(mock.cfg)
//	mock.qapi, _ = client.New(mock.q.Client(), nil)
//	mock.initExec()
//
//}
//
//func (mock *testExecMock) FreeEnv() {
//	util.CloseTestDB(mock.dbDir, mock.stateDB)
//}
//
//func (mock *testExecMock) initExec() {
//	mock.dbDir, mock.stateDB, mock.localDB = util.CreateTestDB()
//	exec := newX2ethereum()
//	exec.SetAPI(mock.qapi)
//	exec.SetStateDB(mock.stateDB)
//	exec.SetLocalDB(mock.localDB)
//	exec.SetEnv(100, 1539918074, 1539918074)
//	mock.exec = exec
//}
//
//func createTx(mock *testExecMock, payload types.Message, priv string, systemCreate bool) (*types.Transaction, error) {
//
//	c, err := crypto.New(crypto.GetName(types.SECP256K1))
//	if err != nil {
//		return nil, err
//	}
//	bytes, err := common.FromHex(priv[:])
//	if err != nil {
//		return nil, err
//	}
//	privKey, err := c.PrivKeyFromBytes(bytes)
//	if err != nil {
//		return nil, err
//	}
//	if systemCreate {
//		action, _ := buildAction(payload)
//		tx, err := types.CreateFormatTx(mock.cfg, mock.cfg.ExecName(types2.PrivacyX), types.Encode(action))
//		if err != nil {
//			return nil, err
//		}
//		tx.Sign(int32(types.SECP256K1), privKey)
//		return tx, nil
//	}
//	req := payload.(*types2.ReqCreatePrivacyTx)
//	if req.GetAssetExec() == "" {
//		req.AssetExec = "coins"
//	}
//	reply, err := mock.wallet.GetAPI().ExecWalletFunc(testPolicyName, "CreateTransaction", payload)
//	if err != nil {
//		return nil, errors.New("createTxErr:" + err.Error())
//	}
//	signTxReq := &types.ReqSignRawTx{
//		TxHex: common.ToHex(types.Encode(reply)),
//	}
//	_, signTx, err := mock.policy.SignTransaction(privKey, signTxReq)
//	if err != nil {
//		return nil, errors.New("signPrivacyTxErr:" + err.Error())
//	}
//	signTxBytes, _ := common.FromHex(signTx)
//	tx := &types.Transaction{}
//	err = types.Decode(signTxBytes, tx)
//	if err != nil {
//		return nil, err
//	}
//	return tx, nil
//}
//
//func buildAction(param types.Message,actionName string) (types.Message, error) {
//
//	action := &types2.X2EthereumAction{
//		Value: nil,
//		Ty:    0,
//	}
//	if val, ok := param.(*types2.Eth2Chain33); ok {
//		if actionName == types2.NameEth2Chain33Action {
//			action.Value = &types2.X2EthereumAction_Eth2Chain33{Eth2Chain33: val}
//			action.Ty = types2.TyEth2Chain33Action
//		} else if actionName == types2.NameWithdrawEthAction {
//			action.Value = &types2.X2EthereumAction_WithdrawEth{WithdrawEth: val}
//			action.Ty = types2.TyWithdrawEthAction
//		}
//	} else if val, ok := param.(*types2.Chain33ToEth); ok {
//		if actionName == types2.NameChain33ToEthAction {
//			action.Value = &types2.X2EthereumAction_Chain33ToEth{Chain33ToEth: val}
//			action.Ty = types2.TyChain33ToEthAction
//		} else if actionName == types2.NameWithdrawChain33Action {
//			action.Value = &types2.X2EthereumAction_WithdrawChain33{WithdrawChain33: val}
//			action.Ty = types2.TyWithdrawChain33Action
//		}
//	} else if val, ok := param.(*types2.MsgValidator); ok {
//		if actionName == types2.NameAddValidatorAction {
//			action.Value = &types2.X2EthereumAction_AddValidator{AddValidator: val}
//			action.Ty = types2.TyAddValidatorAction
//		} else if actionName == types2.NameRemoveValidatorAction {
//			action.Value = &types2.X2EthereumAction_RemoveValidator{RemoveValidator: val}
//			action.Ty = types2.TyRemoveValidatorAction
//		} else if actionName == types2.NameModifyPowerAction {
//			action.Value = &types2.X2EthereumAction_ModifyPower{ModifyPower: val}
//			action.Ty = types2.TyModifyPowerAction
//		}
//	} else if val, ok := param.(*types2.MsgConsensusThreshold); ok {
//		action.Value = &types2.X2EthereumAction_SetConsensusThreshold{SetConsensusThreshold: val}
//		action.Ty = types2.TySetConsensusThresholdAction
//	} else {
//		return nil, types.ErrActionNotSupport
//	}
//	return action, nil
//}
