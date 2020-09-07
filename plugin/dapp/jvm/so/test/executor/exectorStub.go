package executor

import (
	"fmt"
	dbm "github.com/33cn/chain33/common/db"
)

var (
	height = int64(10000)
	blockStoreDB dbm.DB
	localDB dbm.DB
)

func init() {
	blockStoreDB = dbm.NewDB("stateDB", "memdb", "chain.cfg.DbPath", 0)
	localDB = dbm.NewDB("localDB", "memdb", "chain.cfg.DbPath", 0)
}

////////////以下接口用于user.jvm.xxx合约内部转账/////////////////////////////
func ExecFrozen(from string, amount int64) bool {
	fmt.Println("ExecFrozen is called", "from=", from, "amout=", amount)
	return true
}

// ExecActive 激活user.jvm.xxx合约addr上的部分余额
func ExecActive(from string, amount int64) bool {
	fmt.Println("ExecActive is called", "from=", from, "amout=", amount)
	return true
}

// ExecTransfer transfer exec
func ExecTransfer(from, to string, amount int64) bool {
	fmt.Println("ExecTransfer is called", "from=", from, "to=", to, "amout=", amount)
	return true
}

// ExecTransferFrozen 冻结的转账
func ExecTransferFrozen(from, to string, amount int64) bool {
	fmt.Println("ExecTransferFrozen is called", "from=", from, "to=", to, "amout=", amount)
	return true
}

func GetRandom() (string, error) {
	return "0x0123456789abc", nil
}

func GetFrom() string {
	return "0x14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
}

func GetHeight() int64 {
	defer func() {
		height++
	}()
	return height
}

func StopTransWithErrInfo(errInfo string) bool {
	fmt.Println("StopTransWithErrInfo:", errInfo)
	return true
}

// db operation
func GetValueFromLocal(key []byte) []byte {
	val, err := localDB.Get(key)
	fmt.Println("GetValueFromLocal", "key=", string(key), ", ", "value=", string(val))
	if nil != err {
		return nil
	}
	return val
}

func SetValue2Local(key, value []byte) bool {
	fmt.Println("SetValue2Local", "key=", string(key), ", ", "value=", string(value))
	if err := localDB.Set(key, value); nil != err {
		return false
	}
	return true
}

func StateDBGetState(key []byte) []byte {
	val, err := blockStoreDB.Get(key)
	fmt.Println("StateDBGetState", "key=", string(key), ", ", "value=", string(val))
	if nil != err {
		return nil
	}
	return val
}

func StateDBSetState(key, value []byte) bool {
	fmt.Println("StateDBSetState", "key=", string(key), ", ", "value=", string(value))
	if err := blockStoreDB.Set(key, value); nil != err {
		return false
	}
	return true
}


