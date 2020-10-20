package executor

import (
	"fmt"
	dbm "github.com/33cn/chain33/common/db"
	"unsafe"
)

var (
	height = int64(10000)
	blockStoreDB dbm.DB
	localDB dbm.DB
)

type JVMExecutor struct {
	contract string
	txIndex  int
	cp       string
}

func init() {
	blockStoreDB = dbm.NewDB("stateDB", "memdb", "chain.cfg.DbPath", 0)
	localDB = dbm.NewDB("localDB", "memdb", "chain.cfg.DbPath", 0)
}

////////////以下接口用于user.jvm.xxx合约内部转账/////////////////////////////
func ExecFrozen(from string, amount int64, envHandle uintptr) bool {
	fmt.Println("ExecFrozen is called", "from=", from, "amout=", amount)
	fmt.Println("ExecFrozen", "envHandle=", envHandle)
	return true
}

// ExecActive 激活user.jvm.xxx合约addr上的部分余额
func ExecActive(from string, amount int64, envHandle uintptr) bool {
	fmt.Println("ExecActive is called", "from=", from, "amout=", amount)
	fmt.Println("ExecActive", "envHandle=", envHandle)
	return true
}

// ExecTransfer transfer exec
func ExecTransfer(from, to string, amount int64, envHandle uintptr) bool {
	fmt.Println("ExecTransfer is called", "from=", from, "to=", to, "amout=", amount)
	fmt.Println("ExecTransfer", "envHandle=", envHandle)
	return true
}

// ExecTransferFrozen 冻结的转账
func ExecTransferFrozen(from, to string, amount int64, envHandle uintptr) bool {
	fmt.Println("ExecTransferFrozen is called", "from=", from, "to=", to, "amout=", amount)
	fmt.Println("ExecTransferFrozen", "envHandle=", envHandle)
	return true
}

func GetRandom(envHandle uintptr) (string, error) {
	fmt.Println("GetRandom", "envHandle=", envHandle)
	return "0x0123456789abc", nil
}

func GetFrom(envHandle uintptr) string {
	fmt.Println("GetFrom", "envHandle=", envHandle)
	return "0x14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
}

func GetHeight(envHandle uintptr) int64 {
	fmt.Println("GetHeight", "envHandle=", envHandle)
	defer func() {
		height++
	}()
	return height
}

func StopTransWithErrInfo(errInfo string, envHandle uintptr) bool {
	fmt.Println("StopTransWithErrInfo:", errInfo)
	fmt.Println("StopTransWithErrInfo", "envHandle=", envHandle)
	return true
}

func ForwardQueryResult(exceptionOccurred bool, info []string, jvmHandle uintptr) bool {
	for i, info := range info {
		fmt.Println("ForwardQueryResult", "index=", i, ", info=", info, ", exceptionOccurred=", exceptionOccurred)
	}

	fmt.Println("ForwardQueryResult", "jvmHandle=", jvmHandle)
	return true
}

func RecordTxJVMEnv(jvm *JVMExecutor, envHandle uintptr ) bool {
	//jvmExecutor := (*JVMExecutor)(unsafe.Pointer(goHandle))

	fmt.Println("RecordTxJVMEnv", "jvm=", uintptr(unsafe.Pointer(jvm)), ", ", "envHandle=", envHandle)
	return true
}

// db operation
func GetValueFromLocal(key []byte, envHandle uintptr) []byte {
	val, err := localDB.Get(key)
	fmt.Println("GetValueFromLocal", "key=", string(key), ", ", "value=", string(val))
	fmt.Println("GetValueFromLocal", "envHandle=", envHandle)
	if nil != err {
		return nil
	}
	return val
}

func SetValue2Local(key, value []byte, envHandle uintptr) bool {
	fmt.Println("SetValue2Local", "key=", string(key), ", ", "value=", string(value))
	fmt.Println("SetValue2Local", "envHandle=", envHandle)
	if err := localDB.Set(key, value); nil != err {
		return false
	}
	return true
}

func StateDBGetState(key []byte, envHandle uintptr) []byte {
	val, err := blockStoreDB.Get(key)
	fmt.Println("StateDBGetState", "key=", string(key), ", ", "value=", string(val))
	fmt.Println("StateDBGetState", "envHandle=", envHandle)
	if nil != err {
		return nil
	}
	return val
}

func StateDBSetState(key, value []byte, envHandle uintptr) bool {
	fmt.Println("StateDBSetState", "key=", string(key), ", ", "value=", string(value))
	fmt.Println("StateDBSetState", "envHandle=", envHandle)
	if err := blockStoreDB.Set(key, value); nil != err {
		return false
	}
	return true
}


