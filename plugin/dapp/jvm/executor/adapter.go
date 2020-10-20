package executor

//#cgo CFLAGS: -I../openjdk/header
//#cgo LDFLAGS: -L../openjdk -ljli
//#cgo LDFLAGS: -ldl -lpthread -lc
//#include <stdlib.h>
//#include <jli.h>
import "C"

import (
	"errors"
	chain33Types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor/state"
	"unsafe"
)

const (
	JLI_SUCCESS = int(0)
	JLI_FAIL    = int(-1)
	TX_EXEC_JOB = C.int(0)
	TX_QUERY_JOB = C.int(1)
	Bool_TRUE  = C.int(1)
	Bool_FALSE = C.int(0)
)

var (
	jvm_init_alreay = false
	consensusType = ""
)

//调用java合约交易
func runJava(contract string, para []string, jvmHandleGo *JVMExecutor, jobType C.int,  chain33Config *chain33Types.Chain33Config) error {
	//第一次调用java合约时，进行jvm的初始化
	initJvm(chain33Config)

	//构建jdk的输入参数
	tx2Exec := append([]string{contract}, para...)
	argc, argv := buildJavaArgument(tx2Exec)
	if TX_EXEC_JOB == jobType {
		//因为query的内在逻辑问题，参数的内存释放由jdk内部进行释放
		defer freeArgument(argc, argv)
	}

	var exception1DPtr *C.char
	exception := &exception1DPtr
	result := C.JLI_Exec_Contract(argc, argv, exception, jobType, (*C.char)(unsafe.Pointer(jvmHandleGo)))
	if int(result) != JLI_SUCCESS {
		exInfo := C.GoString(*exception)
		defer C.free(unsafe.Pointer(*exception))
		log.Debug("adapter::runJava", "java exception", exInfo)
		return errors.New(exInfo)
	}
	return nil
}

func initJvm(chain33Config *chain33Types.Chain33Config) {
	if jvm_init_alreay {
		return
	}

	const_jdkPath := C.CString(jdkPath)
	defer C.free(unsafe.Pointer(const_jdkPath))

	result := C.JLI_Create_JVM(const_jdkPath)
	if int(result) != JLI_SUCCESS {
		panic("Failed to init JLI_Init_JVM")
	}

	state.IsPara = chain33Config.IsPara()
	state.Title = chain33Config.GetTitle()
	consensusType = chain33Config.GetModuleConfig().Consensus.Name

	jvm_init_alreay = true
}

func buildJavaArgument(execPara []string) (C.int, **C.char) {
	argc := C.int(len(execPara))

	nil2dPtr := C.GetNil2dPtr()
	argv := (**C.char)(C.malloc(C.ulong(argc * C.GetPtrSize())))
	if argv == nil2dPtr {
		panic("Failed to malloc for argv")
	}
	//argv [argc]*C.char
	for i, para := range execPara {
		paraCstr := C.CString(para)
		C.SetPtr(argv, paraCstr, C.int(i))
	}
	return argc, argv
}

func freeArgument(argc C.int, argv **C.char) {
	C.FreeArgv(argc, argv)
}

//export SetQueryResult
func SetQueryResult(jvmgo *C.char, exceptionOccurred C.int, info **C.char, count, sizePtr C.int) C.int {
	jvmHandleUintptr := uintptr(unsafe.Pointer(jvmgo))
	exceOcc := false
	if Bool_TRUE == exceptionOccurred {
		exceOcc = true
	}
	var query []string
	for i:=0; i < int(count); i++ {
		ptr := (uintptr)(unsafe.Pointer(info)) + (uintptr)(int(sizePtr) * i)
		infoGO := C.GoString(*(**C.char)(unsafe.Pointer(ptr)))
		query = append(query, infoGO)
	}
	ForwardQueryResult(exceOcc, query, jvmHandleUintptr)

	return 0
}

//用来保存txjvm或者是queryjvm中的env handle

//export BindTxQueryJVMEnvHandle
func BindTxQueryJVMEnvHandle(jvmGoHandle, envHandle *C.char) C.int {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	jvmExecutor := (*JVMExecutor)(unsafe.Pointer(jvmGoHandle))
	if !recordTxJVMEnv(jvmExecutor, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

/*
 * Account
 */
//export ExecFrozen
func ExecFrozen(from *C.char, amount C.long, envHandle *C.char) C.int {
	fromGoStr := C.GoString(from)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !execFrozen(fromGoStr, int64(amount), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecActive
func ExecActive(from *C.char, amount C.long, envHandle *C.char) C.int {
	fromGoStr := C.GoString(from)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !execActive(fromGoStr, int64(amount), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecTransfer
func ExecTransfer(from, to *C.char, amount C.long, envHandle *C.char) C.int {
	fromGoStr := C.GoString(from)
	toGoStr := C.GoString(to)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !execTransfer(fromGoStr, toGoStr, int64(amount), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

/*
 * blockchain misc
 */
//调用者负责释放返回指针内存
//export GetRandom
func GetRandom(envHandle *C.char) *C.char {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	random, err := getRandom(envHandleUintptr)
	if nil != err {
		return nil
	}
	return C.CString(random)
}

//调用者负责释放返回指针内存
//export GetFrom
func GetFrom(envHandle *C.char) *C.char {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	from := getFrom(envHandleUintptr)
	return C.CString(from)
}

//export GetCurrentHeight
func GetCurrentHeight(envHandle *C.char) C.long {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	return C.long(getHeight(envHandleUintptr))
}

//export StopTransWithErrInfo
func StopTransWithErrInfo(errInfo *C.char, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(errInfo))
	errInfoStr := C.GoString(errInfo)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	stopTransWithErrInfo(errInfoStr, envHandleUintptr)

	return Bool_TRUE
}

/*
 * State DB
 */
//export SetState
func SetState(key *C.char, keySize C.int, value *C.char, valueSize C.int, envHandle *C.char) C.int {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	valSlice := C.GoBytes(unsafe.Pointer(value), valueSize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !stateDBSetState(keySlice, valSlice, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//需要调用者释放内存
//export GetFromState
func GetFromState(key *C.char, keySize C.int, valueSize *C.int, envHandle *C.char) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	value := stateDBGetState(keySlice, envHandleUintptr)
	*valueSize = C.int(len(value))
	return (*C.char)(C.CBytes(value))
}

//export SetStateInStr
func SetStateInStr(key *C.char, value *C.char, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(key))
	defer C.free(unsafe.Pointer(value))
	keyStr := C.GoString(key)
	valueStr := C.GoString(value)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !stateDBSetState([]byte(keyStr), []byte(valueStr), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//调用者负责释放返回指针内存
//export GetFromStateInStr
func GetFromStateInStr(key *C.char, size *C.int, envHandle *C.char) *C.char {
	defer C.free(unsafe.Pointer(key))
	keyStr := C.GoString(key)
	if "" == keyStr {
		*size = C.int(0)
		return nil
	}
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	valueSlice := stateDBGetState([]byte(keyStr), envHandleUintptr)
	valSize := len(valueSlice)
	if 0 == valSize {
		*size = C.int(0)
		return nil
	}
	*size = C.int(valSize)
	return C.CString(string(valueSlice))
}

/*
 * Local DB
 */
//export SetLocal
func SetLocal(key *C.char, keySize C.int, value *C.char, valueSize C.int, envHandle *C.char) C.int {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	valSlice := C.GoBytes(unsafe.Pointer(value), valueSize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !setValue2Local(keySlice, valSlice, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export GetFromLocal
func GetFromLocal(key *C.char, keySize C.int, valueSize *C.int, envHandle *C.char) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	value := getValueFromLocal(keySlice, envHandleUintptr)
	*valueSize = C.int(len(value))
	return (*C.char)(C.CBytes(value))
}

//export SetLocalInStr
func SetLocalInStr(key *C.char, value *C.char, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(key))
	defer C.free(unsafe.Pointer(value))
	keyStr := C.GoString(key)
	valueStr := C.GoString(value)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !setValue2Local([]byte(keyStr), []byte(valueStr), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//调用者负责释放返回指针内存
//export GetFromLocalInStr
func GetFromLocalInStr(key *C.char, size *C.int, envHandle *C.char) *C.char {
	defer C.free(unsafe.Pointer(key))
	keyStr := C.GoString(key)
	if "" == keyStr {
		*size = C.int(0)
		return nil
	}
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	valueSlice := getValueFromLocal([]byte(keyStr), envHandleUintptr)
	valSize := len(valueSlice)
	if 0 == valSize {
		*size = C.int(0)
		return nil
	}
	*size = C.int(valSize)
	return C.CString(string(valueSlice))
}
