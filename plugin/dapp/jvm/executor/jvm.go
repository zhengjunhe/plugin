package executor

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"unsafe"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor/state"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
	lru "github.com/hashicorp/golang-lru"
)

type subConfig struct {
	ParaRemoteGrpcClient string `json:"paraRemoteGrpcClient"`
}

type exception struct {
	occurred bool
	info error
}

// JVMExecutor 执行器结构
type JVMExecutor struct {
	drivers.DriverBase
	mStateDB  *state.MemoryStateDB
	tx        *types.Transaction
	contract  string
	txIndex   int
	excep     exception
	queryChan chan QueryResult
}

type QueryResult struct {
	exceptionOccurred bool
	info []string
}

var (
	log        = log15.New("module", "execs.jvm")
	jvmsCached *lru.Cache
	jvmsCacheCreated = int32(0)
	jdkPath string
)

func initExecType() {
	ety := types.LoadExecutorType(jvmTypes.JvmX)
	ety.InitFuncList(types.ListMethod(&JVMExecutor{}))
}

// Init register function
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	if sub != nil {
		types.MustDecode(sub, &cfg)
	}
	drivers.Register(cfg, GetName(), newJVM, cfg.GetDappFork(jvmTypes.JvmX, "Enable"))
	initExecType()

	conf := types.ConfSub(cfg, jvmTypes.JvmX)
	jdkPath = conf.GStr("jdkPath")
	if "" == jdkPath {
		panic("JDK path is not configured")
	}
	log.Info("jvm::Init", "JDK path is configured to:", jdkPath)
}

func newJVM() drivers.Driver {
	return newJVMDriver()
}

// GetName get name for execname
func GetName() string {
	return newJVM().GetName()
}

func newJVMDriver() drivers.Driver {
	jvm := NewJVMExecutor()
	return jvm
}

// NewJVMExecutor new a jvm executor
func NewJVMExecutor() *JVMExecutor {
	exec := &JVMExecutor{}
	exec.SetChild(exec)
	exec.SetExecutorType(types.LoadExecutorType(jvmTypes.JvmX))
	atomic.LoadInt32(&jvmsCacheCreated)
	if int32(Bool_TRUE) != atomic.LoadInt32(&jvmsCacheCreated) {
		var err error
		jvmsCached, err = lru.New(1000)
		if nil != err {
			panic("Failed to new lru for caching jvms due to:"+ err.Error())
		}
		atomic.StoreInt32(&jvmsCacheCreated, int32(Bool_TRUE))
	}
	return exec
}

func recordTxJVMEnv(jvm *JVMExecutor, envHandle uintptr ) bool {
	jvmsCached.Add(envHandle, jvm)
	_, ok := jvmsCached.Get(envHandle)
	return ok
}

func getJvmExector(envHandle uintptr) (*JVMExecutor, bool) {
	value, ok := jvmsCached.Get(envHandle)
	if !ok {
		log.Error("getJvmExector", "Failed to get JVMExecutor from lru cache with key", envHandle)
		return nil, false
	}

	jvmExecutor, ok := value.(*JVMExecutor)
	if !ok {
		log.Error("getJvmExector", "Failed to get JVMExecutor for query with key", envHandle)
		return nil, false
	}
	return jvmExecutor, true
}

/////////////////////////LocalDB interface//////////////////////////////////////////
func getValueFromLocal(key []byte, envHandle uintptr) []byte {
	log.Debug("Entering GetValueFromLocal", "key", string(key))
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return nil
	}
	contractAddrgo := jvmExecutor.GetContractAddr()
	value := jvmExecutor.mStateDB.GetValueFromLocal(contractAddrgo, string(key))
	if 0 == len(value) {
		log.Debug("Entering Get GetValueFromLocal", "get null value for key", string(key))
		return nil
	}
	return value
}

func setValue2Local(key, value []byte, envHandle uintptr) bool {
	log.Debug("setValue2Local", "key", string(key), "value in string:", string(value),
		"value in slice:", value)
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	contractAddrgo :=  jvmExecutor.GetContractAddr()
	return jvmExecutor.mStateDB.SetValue2Local(contractAddrgo, string(key), value)
}

func stateDBGetState(key []byte, envHandle uintptr) []byte {
	log.Debug("Entering StateDBGetState", "key", string(key))
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		log.Error("stateDBGetState", "Can't get jvmExecutor for key", string(key))
		return nil
	}
	contractAddrgo := jvmExecutor.GetContractAddr()
	value := jvmExecutor.mStateDB.GetState(contractAddrgo, string(key))
	if 0 == len(value) {
		log.Debug("Entering Get StateDBGetState", "get null value for key", string(key))
		return nil
	}

	log.Debug("StateDBGetState Succeed to get value", "value in string", string(value), "value in slice", value)

	return value
}

func stateDBSetState(key, value []byte, envHandle uintptr) bool {
	log.Debug("StateDBSetStateCallback", "key", string(key), "value in string:", string(value),
		"value in slice:", value)
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	contractAddrgo :=  jvmExecutor.GetContractAddr()
	return jvmExecutor.mStateDB.SetState(contractAddrgo, string(key), value)
}

////////////以下接口用于user.jvm.xxx合约内部转账/////////////////////////////
//必须要使用回传的envhandle获取jvm结构指针，否则存在java合约跨合约操作的安全性问题,
//比如在查询的时候，恶意发起数据库写（其中的账户操作就是）的操作，
func execFrozen(from string, amount int64, envHandle uintptr) bool {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	if nil == jvmExecutor || nil == jvmExecutor.mStateDB {
		log.Error("ExecFrozen failed due to nil handle", "pJvm", jvmExecutor,
			"pJvmMap[uint64(jvmIndex)].mStateDB", jvmExecutor.mStateDB)
		return jvmTypes.AccountOpFail
	}
	return jvmExecutor.mStateDB.ExecFrozen(jvmExecutor.tx, from, amount * jvmTypes.Coin_Precision)
}

// ExecActive 激活user.jvm.xxx合约addr上的部分余额
func execActive(from string, amount int64, envHandle uintptr) bool {
	log.Debug("Enter ExecActive", "from", from, "amount", amount)
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		log.Error("ExecActive", "Failed to getJvmExector")
		return jvmTypes.AccountOpFail
	}
	if nil == jvmExecutor || nil == jvmExecutor.mStateDB {
		log.Error("ExecActive failed due to nil handle", "pJvm", jvmExecutor,
			"pJvmMap[uint64(jvmIndex)].mStateDB", jvmExecutor.mStateDB)
		return jvmTypes.AccountOpFail
	}
	return jvmExecutor.mStateDB.ExecActive(jvmExecutor.tx, from, amount*jvmTypes.Coin_Precision)
}

// ExecTransfer transfer exec
func execTransfer(from, to string, amount int64, envHandle uintptr) bool {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	if nil == jvmExecutor || nil == jvmExecutor.mStateDB {
		log.Error("ExecTransfer failed due to nil handle", "pJvm", jvmExecutor,
			"pJvmMap[uint64(jvmIndex)].mStateDB", jvmExecutor.mStateDB)
		return jvmTypes.AccountOpFail
	}
	return jvmExecutor.mStateDB.ExecTransfer(jvmExecutor.tx, from, to, amount * jvmTypes.Coin_Precision)
}

// ExecTransferFrozen 冻结的转账
func execTransferFrozen(from, to string, amount int64, envHandle uintptr) bool {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	if nil == jvmExecutor || nil == jvmExecutor.mStateDB {
		log.Error("ExecTransferFrozen failed due to nil handle", "pJvm", jvmExecutor,
			"pJvmMap[uint64(jvmIndex)].mStateDB", jvmExecutor.mStateDB)
		return jvmTypes.AccountOpFail
	}
	return jvmExecutor.mStateDB.ExecTransferFrozen(jvmExecutor.tx, from, to, int64(amount)*jvmTypes.Coin_Precision)
}

// GetRandom 为jvm用户自定义合约提供随机数，该随机数是64位hash值,返回值为实际返回的长度
func getRandom(envHandle uintptr) (string, error) {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return "", jvmTypes.ErrGetJvmFailed
	}

	if consensusType != "ticket" {
		return "0x42f4eada40e876c476204dfb0749b2cda90020c68992dcacba6ea5a0fa75a371", nil
	}

	req := &types.ReqRandHash{
		ExecName: "ticket",
		BlockNum: jvmExecutor.GetHeight(),
		Hash:     jvmExecutor.GetLastHash(),
	}
	data, err := jvmExecutor.GetExecutorAPI().GetRandNum(req)
	if nil != err {
		log.Error("GetRandom failed due to:", err.Error())
		return "", err
	}
	return string(data), nil
}

func getFrom(envHandle uintptr) string {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return ""
	}
	if nil == jvmExecutor || nil == jvmExecutor.tx {
		log.Error("GetFrom failed due to nil jvmExecutor or nil tx ", "pJvm", jvmExecutor)
		return ""
	}
	return jvmExecutor.tx.From()
}

func getHeight(envHandle uintptr) int64 {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return 0
	}
	if nil == jvmExecutor {
		log.Error("GetFrom failed due to nil handle", "pJvm", jvmExecutor)
		return 0
	}
	return jvmExecutor.GetHeight()
}

func stopTransWithErrInfo(err string, envHandle uintptr) bool {
	jvmExecutor, ok := getJvmExector(envHandle)
	if !ok {
		return false
	}
	if nil == jvmExecutor {
		log.Error("StopTransWithErrInfo failed due to nil handle", "pJvm", jvmExecutor)
		return false
	}
	jvmExecutor.excep.occurred = true
	jvmExecutor.excep.info = errors.New(err)

	log.Info("StopTransWithErrInfo", "error info", err)

	return true
}

//forward the query result to the corresponding jvm
func ForwardQueryResult(exceptionOccurred bool, info []string, jvmHandle uintptr) bool {
	queryResult := QueryResult{
		exceptionOccurred:exceptionOccurred,
		info:info,
	}
	jvm := (*JVMExecutor)(unsafe.Pointer(jvmHandle))
	jvm.queryChan<-queryResult
	log.Info("ForwardQueryResult get query result and forward it", "queryResult", queryResult)
	return true
}

// GetDriverName 获取driver 名称
func (jvm *JVMExecutor) GetDriverName() string {
	return jvmTypes.JvmX
}

// Allow 允许哪些交易在本命执行器执行
func (jvm *JVMExecutor) Allow(tx *types.Transaction, index int) error {
	err := jvm.DriverBase.Allow(tx, index)
	if err == nil {
		return nil
	}
	//增加新的规则:
	//主链: user.jvm.xxx  执行 jvm用户自定义 合约
	//平行链: user.p.guodun.user.jvm.xxx 执行 jvm用户自定义合约
	exec := jvm.GetAPI().GetConfig().GetParaExec(tx.Execer)
	if jvm.AllowIsUserDot2(exec) {
		return nil
	}
	return types.ErrNotAllow
}

func (jvm *JVMExecutor) prepareExecContext(tx *types.Transaction, index int) {
	paraExector := string(jvm.GetAPI().GetConfig().GetParaExec(tx.Execer))
	if jvm.mStateDB == nil {
		jvm.mStateDB = state.NewMemoryStateDB(paraExector, jvm.GetStateDB(), jvm.GetLocalDB(), jvm.GetCoinsAccount(), jvm.GetHeight())
	}
	// 合约具体分为jvm平台合约和基于平台合约的具体合约，如dice合约
	// 获取字节码时通过jvm合约获取，具体执行时要通过具体的合约如dice
	// 每一个区块中会对执行器对象进行缓存，而jvm和user.jvm.dice是两个不同的执行器，因此会产生两个执行器缓存对象，两个对象的accounts字段是互相独立的
	// 更新合约会调用update接口，该接口会更新合约字节码和abi，并缓存在accounts字段中，但只会更新jvm执行器的缓存，而不会更新dice执行器的缓存，因此需要手动从数据库获取数据更新dice缓存
	if strings.HasPrefix(paraExector, jvmTypes.UserJvmX) {
		//create和update接口的执行器名都是jvm, call的执行器名是user.jvm.XXX
		oldExecName := jvm.mStateDB.ExecutorName
		jvm.mStateDB.SetCurrentExecutorName(jvmTypes.JvmX)
		jvm.mStateDB.UpdateAccounts()
		jvm.mStateDB.SetCurrentExecutorName(oldExecName)
	}

	jvm.tx = tx
	jvm.txIndex = index
}

func (jvm *JVMExecutor) prepareQueryContext(executorName []byte) {
	if jvm.mStateDB == nil {
		jvm.mStateDB = state.NewMemoryStateDB(string(jvm.GetAPI().GetConfig().GetParaExec(executorName)), jvm.GetStateDB(), jvm.GetLocalDB(), jvm.GetCoinsAccount(), jvm.GetHeight())
	}
}

// GenerateExecReceipt generate exec receipt
func (jvm *JVMExecutor) GenerateExecReceipt(snapshot int, execName, caller, contractAddr string, opType jvmTypes.JvmContratOpType) (*types.Receipt, error) {
	curVer := jvm.mStateDB.GetLastSnapshot()

	// 打印合约中生成的日志
	jvm.mStateDB.PrintLogs()

	if curVer == nil {
		return nil, nil
	}
	// 从状态机中获取数据变更和变更日志
	data, logs := jvm.mStateDB.GetChangedData(curVer.GetID(), opType)
	contractReceipt := &jvmTypes.ReceiptJVMContract{Caller: caller, ContractName: execName, ContractAddr: contractAddr}

	runLog := &types.ReceiptLog{
		Ty:  jvmTypes.TyLogCallContractJvm,
		Log: types.Encode(contractReceipt)}
	if opType == jvmTypes.CreateJvmContractAction {
		runLog.Ty = jvmTypes.TyLogCreateUserJvmContract
	} else if opType == jvmTypes.UpdateJvmContractAction {
		runLog.Ty = jvmTypes.TyLogUpdateUserJvmContract
	}

	//jvm子合约的debug信息
	//if len(debugInfo) != 0 {
	//	debugLog := &jvmTypes.JVMLog{LogInfo: debugInfo}
	//	runLog2 := &types.ReceiptLog{
	//		Ty:  jvmTypes.TyLogOutputItemJvm,
	//		Log: types.Encode(debugLog),
	//	}
	//	logs = append(logs, runLog2)
	//}

	logs = append(logs, runLog)
	logs = append(logs, jvm.mStateDB.GetReceiptLogs(contractAddr)...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: data, Logs: logs}

	// 返回之前，把本次交易在区块中生成的合约日志集中打印出来
	if jvm.mStateDB != nil {
		jvm.mStateDB.WritePreimages(jvm.GetHeight())
	}

	jvm.collectJvmTxLog(jvm.tx, contractReceipt, receipt)

	return receipt, nil
}

func (jvm *JVMExecutor) collectJvmTxLog(tx *types.Transaction, cr *jvmTypes.ReceiptJVMContract, receipt *types.Receipt) {
	log.Debug("jvm collect begin")
	log.Debug("Tx info", "txHash", common.ToHex(tx.Hash()), "height", jvm.GetHeight())
	log.Debug("ReceiptJVMContract", "data", fmt.Sprintf("caller=%v, name=%v, addr=%v", cr.Caller, cr.ContractName, cr.ContractAddr))
	log.Debug("receipt data", "type", receipt.Ty)
	for _, kv := range receipt.KV {
		log.Debug("KeyValue", "key", common.ToHex(kv.Key), "value", common.ToHex(kv.Value))
	}
	for _, kv := range receipt.Logs {
		log.Debug("ReceiptLog", "Type", kv.Ty, "log", common.ToHex(kv.Log))
	}
	log.Debug("jvm collect end")
}

// ExecLocal 执行本地的transaction, 并写入localdb
func (jvm *JVMExecutor) ExecLocal(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := jvm.DriverBase.ExecLocal(tx, receipt, index)
	if err != nil {
		return nil, err
	}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	// 需要将Exec中生成的合约状态变更信息写入localdb
	for _, logItem := range receipt.Logs {
		if jvmTypes.TyLogStateChangeItemJvm == logItem.Ty {
			data := logItem.Log
			var changeItem jvmTypes.JVMStateChangeItem
			err = types.Decode(data, &changeItem)
			if err != nil {
				return set, err
			}
			set.KV = append(set.KV, &types.KeyValue{Key: []byte(changeItem.Key), Value: changeItem.CurrentValue})
		}
	}

	return set, err
}

// ExecDelLocal 撤销本地的执行
func (jvm *JVMExecutor) ExecDelLocal(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set, err := jvm.DriverBase.ExecDelLocal(tx, receipt, index)
	if err != nil {
		return nil, err
	}
	if receipt.GetTy() != types.ExecOk {
		return set, nil
	}

	// 需要将Exec中生成的合约状态变更信息从localdb中恢复
	for _, logItem := range receipt.Logs {
		if jvmTypes.TyLogStateChangeItemJvm == logItem.Ty {
			data := logItem.Log
			var changeItem jvmTypes.JVMStateChangeItem
			err = types.Decode(data, &changeItem)
			if err != nil {
				return set, err
			}
			set.KV = append(set.KV, &types.KeyValue{Key: []byte(changeItem.Key), Value: changeItem.PreValue})
		}
	}

	return set, err
}

// 检查合约地址是否存在，此操作不会改变任何状态，所以可以直接从statedb查询
func (jvm *JVMExecutor) checkContractNameExists(req *jvmTypes.CheckJVMContractNameReq) (types.Message, error) {
	contractName := req.JvmContractName
	if len(contractName) == 0 {
		return nil, jvmTypes.ErrAddrNotExists
	}

	if !bytes.Contains([]byte(contractName), []byte(jvmTypes.UserJvmX)) {
		contractName = jvmTypes.UserJvmX + contractName
	}
	exists := jvm.GetMStateDB().Exist(address.ExecAddress(jvm.GetAPI().GetConfig().ExecName(contractName)))
	ret := &jvmTypes.CheckJVMAddrResp{ExistAlready: exists}
	return ret, nil
}

// GetMStateDB get memorystate db
func (jvm *JVMExecutor) GetMStateDB() *state.MemoryStateDB {
	return jvm.mStateDB
}

func (jvm *JVMExecutor) GetContractAddr() string {
	if jvm.tx != nil {
		return address.GetExecAddress(string(jvm.tx.Execer)).String()
	}
	return address.GetExecAddress(jvm.contract).String()
}

// 从交易信息中获取交易目标地址，在创建合约交易中，此地址为空
func getReceiver(tx *types.Transaction) *address.Address {
	if tx.To == "" {
		return nil
	}

	addr, err := address.NewAddrFromString(tx.To)
	if err != nil {
		log.Error("create address form string error", "string:", tx.To)
		return nil
	}

	return addr
}
