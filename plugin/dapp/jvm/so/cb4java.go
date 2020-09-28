package main

//#include <stdio.h>
//#include <stdlib.h>
import "C"


import (
	"github.com/33cn/plugin/plugin/dapp/jvm/executor"
	//"github.com/33cn/plugin/plugin/dapp/jvm/so/test/executor"
	"unsafe"
)

const (
	Bool_TRUE  = C.int(1)
	Bool_FALSE = C.int(0)
)

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
	executor.ForwardQueryResult(exceOcc, query, jvmHandleUintptr)

	return 0
}

//用来保存txjvm或者是queryjvm中的env handle
//export BindTxQueryJVMEnvHandle
func BindTxQueryJVMEnvHandle(jvmGoHandle, envHandle *C.char) C.int {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	jvmExecutor := (*executor.JVMExecutor)(unsafe.Pointer(jvmGoHandle))
	if !executor.RecordTxJVMEnv(jvmExecutor, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

/*
 * Account
 */
//export ExecFrozen
func ExecFrozen(from *C.char, amount C.long, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(from))
	fromGoStr := C.GoString(from)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !executor.ExecFrozen(fromGoStr, int64(amount), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecActive
func ExecActive(from *C.char, amount C.long, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(from))
	fromGoStr := C.GoString(from)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !executor.ExecActive(fromGoStr, int64(amount), envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecTransfer
func ExecTransfer(from, to *C.char, amount C.long, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(from))
	defer C.free(unsafe.Pointer(to))
	fromGoStr := C.GoString(from)
	toGoStr := C.GoString(to)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	if !executor.ExecTransfer(fromGoStr, toGoStr, int64(amount), envHandleUintptr) {
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
	random, err := executor.GetRandom(envHandleUintptr)
	if nil != err {
		return nil
	}
	return C.CString(random)
}

//调用者负责释放返回指针内存
//export GetFrom
func GetFrom(envHandle *C.char) *C.char {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	from := executor.GetFrom(envHandleUintptr)
	return C.CString(from)
}

//export GetCurrentHeight
func GetCurrentHeight(envHandle *C.char) C.long {
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	return C.long(executor.GetHeight(envHandleUintptr))
}

//export StopTransWithErrInfo
func StopTransWithErrInfo(errInfo *C.char, envHandle *C.char) C.int {
	defer C.free(unsafe.Pointer(errInfo))
	errInfoStr := C.GoString(errInfo)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	executor.StopTransWithErrInfo(errInfoStr, envHandleUintptr)

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
	if !executor.StateDBSetState(keySlice, valSlice, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//需要调用者释放内存
//export GetFromState
func GetFromState(key *C.char, keySize C.int, valueSize *C.int, envHandle *C.char) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	value := executor.StateDBGetState(keySlice, envHandleUintptr)
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
	if !executor.StateDBSetState([]byte(keyStr), []byte(valueStr), envHandleUintptr) {
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
	valueSlice := executor.StateDBGetState([]byte(keyStr), envHandleUintptr)
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
	if !executor.SetValue2Local(keySlice, valSlice, envHandleUintptr) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export GetFromLocal
func GetFromLocal(key *C.char, keySize C.int, valueSize *C.int, envHandle *C.char) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	envHandleUintptr := uintptr(unsafe.Pointer(envHandle))
	value := executor.GetValueFromLocal(keySlice, envHandleUintptr)
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
	if !executor.SetValue2Local([]byte(keyStr), []byte(valueStr), envHandleUintptr) {
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
	valueSlice := executor.GetValueFromLocal([]byte(keyStr), envHandleUintptr)
	valSize := len(valueSlice)
	if 0 == valSize {
		*size = C.int(0)
		return nil
	}
	*size = C.int(valSize)
	return C.CString(string(valueSlice))
}

func main() {}




