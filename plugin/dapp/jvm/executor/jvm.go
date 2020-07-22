package executor

import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor/state"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

type subConfig struct {
	ParaRemoteGrpcClient string `json:"paraRemoteGrpcClient"`
}

// JVMExecutor 执行器结构
type JVMExecutor struct {
	drivers.DriverBase
	mStateDB *state.MemoryStateDB
	tx       *types.Transaction
	txIndex  int
	out2User []*jvmTypes.JVMOutItem
	cp       string
}

var (
	//pMemoryStateDB *state.MemoryStateDB
	//pJvm          *JVMExecutor
	pJvmIndex = uint64(1)
	pJvmMap   map[uint64]*JVMExecutor
	log        = log15.New("module", "execs.jvm")
	cfg        subConfig
)

func initExecType() {
	ety := types.LoadExecutorType(jvmTypes.JvmX)
	ety.InitFuncList(types.ListMethod(&JVMExecutor{}))
	pJvmMap = make(map[uint64]*JVMExecutor)
}

// Init register function
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	if sub != nil {
		types.MustDecode(sub, &cfg)
	}
	drivers.Register(cfg, GetName(), newJVM, cfg.GetDappFork(jvmTypes.JvmX, "Enable"))
	initExecType()
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
	return exec
}

func setJvm4CallbackWithIndex(jvm *JVMExecutor, index uint64) {
	pJvmMap[index] = jvm
}

/////////////////////////LocalDB interface//////////////////////////////////////////
// GetValueSizeFromLocal get value size from local
//export GetValueSizeFromLocal
func GetValueSizeFromLocal(key []byte) int32 {
	log.Debug("Entering GetValueSizeFromLocal")
	jvmIndex := 0
	contractAddrgo := address.GetExecAddress(string(pJvmMap[uint64(jvmIndex)].tx.Execer)).String()
	value := pJvmMap[uint64(jvmIndex)].mStateDB.GetValueFromLocal(contractAddrgo, string(key))
	return int32(len(value))
}

// GetValueFromLocal get value form local for C
//export GetValueFromLocal
func GetValueFromLocal(key []byte) []byte {
	log.Debug("Entering GetValueFromLocal")
	jvmIndex := 0
	contractAddrgo := address.GetExecAddress(string(pJvmMap[uint64(jvmIndex)].tx.Execer)).String()
	value := pJvmMap[uint64(jvmIndex)].mStateDB.GetValueFromLocal(contractAddrgo, string(key))
	if 0 == len(value) {
		log.Debug("Entering Get StateDBGetStateCallback", "get null value for key", string(key))
		return nil
	}
	return value
}

func SetValue2Local(key, value []byte) bool {
	log.Debug("StateDBSetStateCallback", "key", string(key))

	jvmIndex := 0
	contractAddrgo := address.GetExecAddress(string(pJvmMap[uint64(jvmIndex)].tx.Execer)).String()
	return pJvmMap[uint64(jvmIndex)].mStateDB.SetValue2Local(contractAddrgo, string(key), value)
}

func StateDBGetValueSizeCallback(contractAddr string, key []byte) int32 {
	jvmIndex := 0
	log.Debug("Entering StateDBGetValueSize")
	value := pJvmMap[uint64(jvmIndex)].mStateDB.GetState(contractAddr, string(key))
	return int32(len(value))
}

// StateDBGetStateCallback get state db callback C
//export StateDBGetStateCallback
func StateDBGetStateCallback(key []byte) []byte {
	log.Debug("Entering Get StateDBGetStateCallback")
	jvmIndex := 0
	contractAddrgo := address.GetExecAddress(string(pJvmMap[uint64(jvmIndex)].tx.Execer)).String()
	value := pJvmMap[uint64(jvmIndex)].mStateDB.GetState(contractAddrgo, string(key))
	if 0 == len(value) {
		log.Debug("Entering Get StateDBGetStateCallback", "get null value for key", string(key))
		return nil
	}

	return value
}

func StateDBSetStateCallback(key, value []byte) bool {
	log.Debug("StateDBSetStateCallback", "key", string(key), "value in string:",
		"value in slice:", value)
	jvmIndex := 0

	contractAddr := address.GetExecAddress(string(pJvmMap[uint64(jvmIndex)].tx.Execer)).String()
	return pJvmMap[uint64(jvmIndex)].mStateDB.SetState(contractAddr, string(key), value)
}

// Output2UserCallback 该接口用于返回查询结果的返回
//export Output2UserCallback
func Output2UserCallback(typeName string, value []byte) {
	log.Debug("Entering Output2UserCallback")
	jvmIndex := 0

	jvmOutItem := &jvmTypes.JVMOutItem{
		ItemType: typeName,
		Data:     value,
	}
	pJvmMap[uint64(jvmIndex)].out2User = append(pJvmMap[uint64(jvmIndex)].out2User, jvmOutItem)

	return
}

////////////以下接口用于user.jvm.xxx合约内部转账/////////////////////////////
func ExecFrozen(from string, amount int64) bool {
	if nil == pJvmMap[uint64(0)] || nil == pJvmMap[uint64(0)].mStateDB {
		log.Error("ExecFrozen failed due to nil handle", "pJvm", pJvmMap[uint64(0)], "pJvmMap[uint64(jvmIndex)].mStateDB", pJvmMap[uint64(0)].mStateDB)
		return jvmTypes.AccountOpFail
	}
	return pJvmMap[0].mStateDB.ExecFrozen(pJvmMap[0].tx, from, amount * jvmTypes.Coin_Precision)
}

// ExecActive 激活user.jvm.xxx合约addr上的部分余额
func ExecActive(from string, amount int64) bool {
	if nil == pJvmMap[0] || nil == pJvmMap[0].mStateDB {
		log.Error("ExecActive failed due to nil handle", "pJvm", pJvmMap[0], "pJvmMap[uint64(jvmIndex)].mStateDB", pJvmMap[0].mStateDB)
		return jvmTypes.AccountOpFail
	}
	return pJvmMap[0].mStateDB.ExecActive(pJvmMap[0].tx, from, amount*jvmTypes.Coin_Precision)
}

// ExecTransfer transfer exec
func ExecTransfer(from, to string, amount int64) bool {
	if nil == pJvmMap[0] || nil == pJvmMap[0].mStateDB {
		log.Error("ExecTransfer failed due to nil handle", "pJvm", pJvmMap[0], "pJvmMap[uint64(jvmIndex)].mStateDB", pJvmMap[0].mStateDB)
		return jvmTypes.AccountOpFail
	}
	return pJvmMap[0].mStateDB.ExecTransfer(pJvmMap[0].tx, from, to, amount * jvmTypes.Coin_Precision)
}

// ExecTransferFrozen 冻结的转账
func ExecTransferFrozen(from, to string, amount int64) bool {
	jvmIndex := 0
	if nil == pJvmMap[uint64(jvmIndex)] || nil == pJvmMap[uint64(jvmIndex)].mStateDB {
		log.Error("ExecTransferFrozen failed due to nil handle", "pJvm", pJvmMap[uint64(jvmIndex)], "pJvmMap[uint64(jvmIndex)].mStateDB", pJvmMap[uint64(jvmIndex)].mStateDB)
		return jvmTypes.AccountOpFail
	}
	return pJvmMap[uint64(jvmIndex)].mStateDB.ExecTransferFrozen(pJvmMap[uint64(jvmIndex)].tx, from, to, int64(amount)*jvmTypes.Coin_Precision)
}

// GetRandom 为jvm用户自定义合约提供随机数，该随机数是64位hash值,返回值为实际返回的长度
func GetRandom() ([]byte, error) {
	jvmIndex := 0
	blockNum := int64(5)
	if nil == pJvmMap[uint64(jvmIndex)] {
		log.Error("GetRandom failed due to nil handle", "pJvm", pJvmMap[uint64(jvmIndex)])
		return nil, errors.New("invalid index")
	}

	req := &types.ReqRandHash{
		ExecName: "ticket",
		BlockNum: blockNum,
		Hash:     pJvmMap[uint64(jvmIndex)].GetLastHash(),
	}
	return pJvmMap[uint64(jvmIndex)].GetExecutorAPI().GetRandNum(req)
}

func GetFrom() string {
	jvmIndex := 0

	if nil == pJvmMap[uint64(jvmIndex)] {
		log.Error("GetFrom failed due to nil handle", "pJvm", pJvmMap[uint64(jvmIndex)])
		return ""
	}
	return pJvmMap[uint64(jvmIndex)].tx.From()
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
