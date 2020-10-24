package executor

import (
	"io/ioutil"
	"os"

	//"github.com/33cn/chain33/client"
	"testing"

	apimock "github.com/33cn/chain33/client/mocks"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/mocks"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var chainTestCfg = types.NewChain33Config(types.GetDefaultCfgstring())

func init() {
	Init(jvmTypes.JvmX, chainTestCfg, nil)
}

var (
	privOpener = getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944") //opener
	privPlayer = getprivkey("4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01") //player
	opener     = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	player     = "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
	tx         = &types.Transaction{}
	sdb        *db.GoMemDB
	kvdb       db.KVDB
)

type JvmTestEnv struct {
	kvdb      *mocks.KVDB
	jvm       *JVMExecutor
}

func setupTestEnv() *JvmTestEnv{
	jvmTestEnv := &JvmTestEnv{}
	jvmTestEnv.kvdb = new(mocks.KVDB)
	jvmExecutor := &JVMExecutor{DriverBase: drivers.DriverBase{}}

	_, _, kvdb = util.CreateTestDB()
	jvmExecutor.SetLocalDB(kvdb)
	api := new(apimock.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(chainTestCfg, nil)
	jvmExecutor.SetAPI(api)
	sdb, _ = db.NewGoMemDB("JvmTestDb", "test", 128)
	jvmExecutor.SetStateDB(sdb)
	jvmExecutor.SetEnv(10, 100, 1)
	jvmExecutor.SetIsFree(false)
	jvmExecutor.SetChild(jvmExecutor)
	jvmTestEnv.jvm = jvmExecutor

	Chain33LoaderJarPath = "../../../../build"

	return jvmTestEnv
}

//包含query 合约是否存在的功能
func Test_CreateJvmContract(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	code := readJarFile("Guess")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	in := &jvmTypes.CheckJVMContractNameReq{
		JvmContractName:"user.jvm.Guess",
	}
	msg, err := jvmTestEnv.jvm.Query_CheckContractNameExist(in)
	resp := msg.(*jvmTypes.CheckJVMAddrResp)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.ExistAlready)

	in.JvmContractName = "Guess"
	msg, err = jvmTestEnv.jvm.Query_CheckContractNameExist(in)
	resp = msg.(*jvmTypes.CheckJVMAddrResp)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, resp.ExistAlready)

	in.JvmContractName = "user.jvm.Hello"
	msg, err = jvmTestEnv.jvm.Query_CheckContractNameExist(in)
	resp = msg.(*jvmTypes.CheckJVMAddrResp)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, resp.ExistAlready)


	msg, err = jvmTestEnv.jvm.Query_CheckContractNameExist(nil)
	assert.Equal(t, types.ErrInvalidParam, err)

	in.JvmContractName = ""
	msg, err = jvmTestEnv.jvm.Query_CheckContractNameExist(in)
	assert.Equal(t, jvmTypes.ErrNullContractName, err)
}

func Test_CreateJvmContract_errorbranch(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	code := readJarFile("Guess")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	//使用相同的合约名字，检测冲突
	_, err = jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, jvmTypes.ErrContractAddressCollisionJvm, err)

	var bigcode = []byte("hello")
	for i := 0; i < jvmTypes.MaxCodeSize + 2; i++ {
		bigcode = append(bigcode, []byte("1")...)
	}

	bigSizeContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Bigsize",
		Code:common.ToHex(bigcode),
	}
	_, err = jvmTestEnv.jvm.Exec_CreateJvmContract(bigSizeContract, tx, 0)
	assert.Equal(t, jvmTypes.ErrMaxCodeSizeExceededJvm, err)

	zeroSizeContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Zerosize",
	}
	_, err = jvmTestEnv.jvm.Exec_CreateJvmContract(zeroSizeContract, tx, 0)
	assert.Equal(t, jvmTypes.ErrNUllJvmContract, err)
}

func Test_Create_CallJvmContract(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	//1st step: create Guess contract
	code := readJarFile("Dice")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Dice",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	////////////////////////////////////
	//2nd step: call Guess contract
	////////////////////////////////////
	callJvmContract := &jvmTypes.CallJvmContract{
		Name:"user.jvm.Dice",
		ActionData:[]string{"startGame"},
	}

	payload2call := types.Encode(callJvmContract)
	tx2call := createTx(jvmTypes.CallJvmContractAction, payload2call, []byte("user.jvm.Dice"))

	receipt2call, err := jvmTestEnv.jvm.Exec_CallJvmContract(callJvmContract, tx2call, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt2call.Ty)
	removeFile("./Dice.jar")
}

func Test_CallJvmContract_errorBranch(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	//1st step: contract not exist
	callJvmContract := &jvmTypes.CallJvmContract{
		Name:"user.jvm.Guess",
		ActionData:[]string{"startGame"},
	}
	payload2call := types.Encode(callJvmContract)
	tx2call := createTx(jvmTypes.CallJvmContractAction, payload2call, []byte("user.jvm.Guess"))

	_, err := jvmTestEnv.jvm.Exec_CallJvmContract(callJvmContract, tx2call, 0)
	assert.Equal(t, jvmTypes.ErrContractNotExist, err)
}

func Test_UpdateJvmContract_errorBranch(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	updateJvmContract := &jvmTypes.UpdateJvmContract{
		Name:"user.jvm.Guess",
	}

	payload2update := types.Encode(updateJvmContract)
	tx2update := createTx(jvmTypes.UpdateJvmContractAction, payload2update, []byte(jvmTypes.JvmX))

	//指定合约名字的合约不存在，更新错误
	_, err := jvmTestEnv.jvm.Exec_UpdateJvmContract(updateJvmContract, tx2update, 0)
	assert.Equal(t, jvmTypes.ErrContractNotExist, err)

	//创建合约
	code := readJarFile("Guess")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	//更新合约
	updateJvmContract2 := &jvmTypes.UpdateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload2update2 := types.Encode(updateJvmContract2)
	tx2update2 := createTx(jvmTypes.UpdateJvmContractAction, payload2update2, []byte(jvmTypes.JvmX))
	tx2update2.Sign(types.SECP256K1, privPlayer)

	_, err = jvmTestEnv.jvm.Exec_UpdateJvmContract(updateJvmContract2, tx2update2, 0)
	assert.Equal(t, jvmTypes.ErrNoPermission, err)

	//合约超大
	var bigcode []byte = []byte("hello")
	for i := 0; i < jvmTypes.MaxCodeSize + 2; i++ {
		bigcode = append(bigcode, []byte("1")...)
	}

	bigSizeContract := &jvmTypes.UpdateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(bigcode),
	}
	payloadBigsize := types.Encode(bigSizeContract)
	txBigsize := createTx(jvmTypes.UpdateJvmContractAction, payloadBigsize, []byte(jvmTypes.JvmX))
	_, err = jvmTestEnv.jvm.Exec_UpdateJvmContract(bigSizeContract, txBigsize, 0)
	assert.Equal(t, jvmTypes.ErrMaxCodeSizeExceededJvm, err)

	//空合约
	zeroSizeContract := &jvmTypes.UpdateJvmContract{
		Name:"user.jvm.Guess",
	}
	payloadZero := types.Encode(zeroSizeContract)
	txZerosize := createTx(jvmTypes.UpdateJvmContractAction, payloadZero, []byte(jvmTypes.JvmX))
	_, err = jvmTestEnv.jvm.Exec_UpdateJvmContract(zeroSizeContract, txZerosize, 0)
	assert.Equal(t, jvmTypes.ErrNUllJvmContract, err)
}

func Test_Create_Update_CallJvmContract(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	//1st step: create Guess contract
	//创建时，使用dice合约，更新时使用guess
	code := readJarFile("Dice")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	////////////////////////////////////
	//2nd step: update Guess contract
	////////////////////////////////////
	updateJvmContract := &jvmTypes.UpdateJvmContract{
		Name:"user.jvm.Guess",
		Code:common.ToHex(code),
	}

	payload2update := types.Encode(updateJvmContract)
	tx2update := createTx(jvmTypes.UpdateJvmContractAction, payload2update, []byte(jvmTypes.JvmX))

	receipt2update, err := jvmTestEnv.jvm.Exec_UpdateJvmContract(updateJvmContract, tx2update, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt2update.Ty)

	////////////////////////////////////
	//3rd step: call the updated Guess contract
	////////////////////////////////////
	callJvmContract := &jvmTypes.CallJvmContract{
		Name:"user.jvm.Guess",
		ActionData:[]string{"startGame"},
	}

	payload2call := types.Encode(callJvmContract)
	tx2call := createTx(jvmTypes.CallJvmContractAction, payload2call, []byte("user.jvm.Guess"))

	receipt2call, err := jvmTestEnv.jvm.Exec_CallJvmContract(callJvmContract, tx2call, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt2call.Ty)

	removeFile("./Guess.jar")
}

func Test_Exec_Order(t *testing.T) {
	jvmTestEnv := setupTestEnv()
	assert.Equal(t, drivers.ExecLocalSameTime, jvmTestEnv.jvm.ExecutorOrder())
}

func Test_Allow(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Guess",
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))
	err := jvmTestEnv.jvm.Allow(tx, 0)
	assert.Equal(t, nil, err)

	callJvmContract := &jvmTypes.CallJvmContract{
		Name:"user.jvm.Dice",
		ActionData:[]string{"startGame"},
	}
	payload2call := types.Encode(callJvmContract)
	tx2call := createTx(jvmTypes.CallJvmContractAction, payload2call, []byte("user.jvm.Dice"))
	err = jvmTestEnv.jvm.Allow(tx2call, 0)
	assert.Equal(t, nil, err)

	callJvmContract2 := &jvmTypes.CallJvmContract{
		Name:"user.wasm.Dice",
		ActionData:[]string{"startGame"},
	}
	payload2call2 := types.Encode(callJvmContract2)
	tx2call2 := createTx(jvmTypes.CallJvmContractAction, payload2call2, []byte("user.wasm.Dice"))
	err = jvmTestEnv.jvm.Allow(tx2call2, 0)
	assert.Equal(t, types.ErrNotAllow, err)
}

func Test_QueryContractRun(t *testing.T) {
	jvmTestEnv := setupTestEnv()

	//1st step: create Guess contract
	code := readJarFile("Dice")
	assert.NotEqual(t,nil, code)
	createJvmContract := &jvmTypes.CreateJvmContract{
		Name:"user.jvm.Dice",
		Code:common.ToHex(code),
	}

	payload := types.Encode(createJvmContract)
	tx := createTx(jvmTypes.CreateJvmContractAction, payload, []byte(jvmTypes.JvmX))

	receipt, err := jvmTestEnv.jvm.Exec_CreateJvmContract(createJvmContract, tx, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt.Ty)

	////////////////////////////////////
	//2nd step: call Guess contract to start game
	////////////////////////////////////
	createJarLib(t)
	callJvmContract := &jvmTypes.CallJvmContract{
		Name:"user.jvm.Dice",
		ActionData:[]string{"startGame"},
	}

	payload2call := types.Encode(callJvmContract)
	tx2call := createTx(jvmTypes.CallJvmContractAction, payload2call, []byte("user.jvm.Dice"))

	receipt2call, err := jvmTestEnv.jvm.Exec_CallJvmContract(callJvmContract, tx2call, 0)
	assert.Equal(t, nil, err)
	assert.Equal(t, int32(types.ExecOk), receipt2call.Ty)


	//////////////////////////////////////
	////3rd step: call Guess contract to play game
	//////////////////////////////////////
	//deposit2contract(t, jvmTestEnv.jvm, "user.jvm.Dice")
	//playGame := &jvmTypes.CallJvmContract{
	//	Name:"user.jvm.Dice",
	//	ActionData:[]string{"playGame", "6", "2"},
	//}
	//
	//payload2play := types.Encode(playGame)
	//tx2play := createTx(jvmTypes.CallJvmContractAction, payload2play, []byte("user.jvm.Dice"))
	//
	//tx2play.Sign(types.SECP256K1, privPlayer)
	//receipt2play, err := jvmTestEnv.jvm.Exec_CallJvmContract(playGame, tx2play, 0)
	//assert.Equal(t, nil, err)
	//assert.Equal(t, int32(types.ExecOk), receipt2play.Ty)
	//
	//jvmTestEnv.jvm.SetEnv(20, 200, 1)
	//
	//closeGame := &jvmTypes.CallJvmContract{
	//	Name:"user.jvm.Dice",
	//	ActionData:[]string{"closeGame"},
	//}
	//
	//payload2close := types.Encode(closeGame)
	//tx2close := createTx(jvmTypes.CallJvmContractAction, payload2close, []byte("user.jvm.Dice"))
	//
	//receipt2close, err := jvmTestEnv.jvm.Exec_CallJvmContract(closeGame, tx2close, 0)
	//assert.Equal(t, nil, err)
	//assert.Equal(t, int32(types.ExecOk), receipt2close.Ty)

	//查询
	removeFile("./Dice.jar")
	jvmQueryReq := &jvmTypes.JVMQueryReq{
		Contract:"user.jvm.Dice",
		Para: []string{"getDiceRecordByRound", "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv", "1"},
	}
	msg, err := jvmTestEnv.jvm.Query_JavaContract(jvmQueryReq)
	assert.Equal(t, nil, err)
	resp := msg.(*jvmTypes.JVMQueryResponse)
	assert.Equal(t, "No info", resp.Result[0])
	assert.Equal(t, true, resp.Success)
	removeFile("./Dice.jar")
	removeJarLib(t)
}

func removeFile(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil && !os.IsExist(err) {
		return
	}
	_ = os.Remove(filePath)
}

func createJarLib(t *testing.T) {
	err := os.Mkdir("./jarlib", 0775)
	assert.Equal(t, nil, err)

	codePath := "../../../../build/jarlib/Gson.jar"
	code, _ := ioutil.ReadFile(codePath)
	jarfile, err := os.OpenFile("./jarlib/Gson.jar", os.O_WRONLY|os.O_CREATE, 0666)
	assert.Equal(t, nil, err)
	writeLen, err := jarfile.Write(code)
	assert.Equal(t, writeLen, len(code))
	closeErr := jarfile.Close()
	assert.Equal(t, nil, closeErr)
}

func removeJarLib(t *testing.T) {
	err := os.RemoveAll("./jarlib")
	assert.Equal(t, nil, err)
}

func createTx(txType int, payload []byte, execer []byte) *types.Transaction {

	switch txType {
	case jvmTypes.CreateJvmContractAction:
		tx := &types.Transaction{
			//Execer:[]byte(jvmTypes.JvmX),
			Execer:execer,
			Payload:payload,
			To:address.ExecAddress(string(tx.Execer)),
			Nonce:1,
		}
		tx.Sign(types.SECP256K1, privOpener)
		return tx
	case jvmTypes.CallJvmContractAction:
		tx := &types.Transaction{
			//Execer:[]byte("user.jvm.Guess"),
			Execer:execer,
			Payload:payload,
			To:address.ExecAddress(string(tx.Execer)),
			Nonce:1,
		}
		tx.Sign(types.SECP256K1, privOpener)
		return tx
	case jvmTypes.UpdateJvmContractAction:
		tx := &types.Transaction{
			//Execer:[]byte(jvmTypes.JvmX),
			Execer:execer,
			Payload:payload,
			To:address.ExecAddress(string(tx.Execer)),
			Nonce:1,
		}
		tx.Sign(types.SECP256K1, privOpener)
		return tx
	default:
		return nil
    }
}

func readJarFile(jarName string) []byte {
	codePath := "../../../../build/" + jarName + ".jar"
	code, _ := ioutil.ReadFile(codePath)
	return code
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func deposit2contract(t *testing.T, jvm *JVMExecutor,  contractName string) {
	acc := jvm.GetCoinsAccount()

	account := &types.Account{
		Balance: 1000 * 1e8,
		Addr:    player,
	}
	contractAddr := address.ExecAddress(contractName)
	acc.SaveAccount(account)
	account = acc.LoadAccount(player)
	assert.Equal(t, int64(1000*1e8), account.Balance)
	_, err := acc.TransferToExec(player, contractAddr, 200*1e8)
	assert.Nil(t, err)
	account = acc.LoadExecAccount(player, contractAddr)
	assert.Equal(t, int64(200*1e8), account.Balance)
}
