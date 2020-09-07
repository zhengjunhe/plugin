package main

//#include <stdio.h>
//#include <stdlib.h>
import "C"

import (
	//"github.com/33cn/plugin/plugin/dapp/jvm/executor"
	"github.com/33cn/plugin/plugin/dapp/jvm/so/test/executor"
	"unsafe"
)

//var Nil_info_ptr *C.char

const (
	Bool_TRUE  = C.int(1)
	Bool_FALSE = C.int(0)
	//NIL_INFO   = "magic-data-nil-0000"
)

//func init() {
//	Nil_info_ptr = C.CString(NIL_INFO)
//}
/*
 * Account
 */
//export ExecFrozen
func ExecFrozen(from *C.char, amount C.long) C.int {
	defer C.free(unsafe.Pointer(from))
	fromGoStr := C.GoString(from)
	if !executor.ExecFrozen(fromGoStr, int64(amount)) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecActive
func ExecActive(from *C.char, amount C.long) C.int {
	defer C.free(unsafe.Pointer(from))
	fromGoStr := C.GoString(from)
	if !executor.ExecActive(fromGoStr, int64(amount)) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export ExecTransfer
func ExecTransfer(from, to *C.char, amount C.long) C.int {
	defer C.free(unsafe.Pointer(from))
	defer C.free(unsafe.Pointer(to))
	fromGoStr := C.GoString(from)
	toGoStr := C.GoString(to)
	if !executor.ExecTransfer(fromGoStr, toGoStr, int64(amount)) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

/*
 * blockchain misc
 */
//调用者负责释放返回指针内存
//export GetRandom
func GetRandom() *C.char {
	random, err := executor.GetRandom()
	if nil != err {
		return nil
	}
	return C.CString(random)
}

//调用者负责释放返回指针内存
//export GetFrom
func GetFrom() *C.char {
	from := executor.GetFrom()
	return C.CString(from)
}

//export GetCurrentHeight
func GetCurrentHeight() C.long {
	return C.long(executor.GetHeight())
}

//export StopTransWithErrInfo
func StopTransWithErrInfo(errInfo *C.char) C.int {
	defer C.free(unsafe.Pointer(errInfo))
	errInfoStr := C.GoString(errInfo)
	executor.StopTransWithErrInfo(errInfoStr)

	return Bool_TRUE
}

/*
 * State DB
 */
//export SetState
func SetState(key *C.char, keySize C.int, value *C.char, valueSize C.int) C.int {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	valSlice := C.GoBytes(unsafe.Pointer(value), valueSize)
	if !executor.StateDBSetState(keySlice, valSlice) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//需要调用者释放内存
//export GetFromState
func GetFromState(key *C.char, keySize C.int, valueSize *C.int) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	value := executor.StateDBGetState(keySlice)
	*valueSize = C.int(len(value))
	return (*C.char)(C.CBytes(value))
}

//export SetStateInStr
func SetStateInStr(key *C.char, value *C.char) C.int {
	defer C.free(unsafe.Pointer(key))
	defer C.free(unsafe.Pointer(value))
	keyStr := C.GoString(key)
	valueStr := C.GoString(value)

	if !executor.StateDBSetState([]byte(keyStr), []byte(valueStr)) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//调用者负责释放返回指针内存
//export GetFromStateInStr
func GetFromStateInStr(key *C.char, size *C.int) *C.char {
	defer C.free(unsafe.Pointer(key))
	keyStr := C.GoString(key)
	if "" == keyStr {
		*size = C.int(0)
		return nil
	}
	valueSlice := executor.StateDBGetState([]byte(keyStr))
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
func SetLocal(key *C.char, keySize C.int, value *C.char, valueSize C.int) C.int {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	valSlice := C.GoBytes(unsafe.Pointer(value), valueSize)
	if !executor.SetValue2Local(keySlice, valSlice) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//export GetFromLocal
func GetFromLocal(key *C.char, keySize C.int, valueSize *C.int) *C.char {
	keySlice := C.GoBytes(unsafe.Pointer(key), keySize)
	value := executor.GetValueFromLocal(keySlice)
	*valueSize = C.int(len(value))
	return (*C.char)(C.CBytes(value))
}

//export SetLocalInStr
func SetLocalInStr(key *C.char, value *C.char) C.int {
	defer C.free(unsafe.Pointer(key))
	defer C.free(unsafe.Pointer(value))
	keyStr := C.GoString(key)
	valueStr := C.GoString(value)

	if !executor.SetValue2Local([]byte(keyStr), []byte(valueStr)) {
		return Bool_FALSE
	}
	return Bool_TRUE
}

//调用者负责释放返回指针内存
//export GetFromLocalInStr
func GetFromLocalInStr(key *C.char, size *C.int) *C.char {
	defer C.free(unsafe.Pointer(key))
	keyStr := C.GoString(key)
	if "" == keyStr {
		*size = C.int(0)
		return nil
	}
	valueSlice := executor.GetValueFromLocal([]byte(keyStr))
	valSize := len(valueSlice)
	if 0 == valSize {
		*size = C.int(0)
		return nil
	}
	*size = C.int(valSize)
	return C.CString(string(valueSlice))
}

func main() {}




