package executor

//#cgo CFLAGS: -DTARGET_ARCH_x86 -DLINUX -DTARGET_COMPILER_gcc -DFULL_VERSION=1.8.0_262-internal-hezhengjun_2020_08_04_23_54-b00 -DJDK_MAJOR_VERSION=1 -DJDK_MINOR_VERSION=8
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/hotspot/src/share/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/hotspot/src/share/vm/prims
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/hotspot/src/cpu/x86/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/share/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/solaris/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/share/native/common
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/solaris/native/common
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/share/bin
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/solaris/bin
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/jvm4chain33/jdk/src/linux/bin
//#cgo LDFLAGS: -L/home/hezhengjun/work/go/src/github.com/33cn/plugin/plugin/dapp/jvm/so -ljli
//#cgo LDFLAGS: -ldl -lpthread -lc
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//#include <stdint.h>
//#include <defines.h>
//int GetPtrSize() {
//	return sizeof(char *);
//}
//void SetPtr(char **ptr, char *value, int index) {
//	ptr[index] = value;
//}
//void FreeArgv(int argc, char **argv) {
//    int i = 0;
//    for (; i < argc; i++) {
//        free(argv[i]);
//    }
//    free(argv);
//}
//char ** GetNil2dPtr() {
//	return (char **)0;
//}
//void * GetVoidPtr(char *voidPtr) {
//	return (void *)voidPtr;
//}
import "C"


import (
	"errors"
	"unsafe"
)

const (
	JLI_SUCCESS = int(0)
	JLI_FAIL    = int(-1)
	TX_EXEC_JOB = C.int(0)
	TX_QUERY_JOB = C.int(1)
)

var (
	jvm_init_alreay = false
)

//调用java合约交易
func runJava(contract, action string, para []string, cp string, jvmHandleGo *JVMExecutor) error {
	//第一次调用java合约时，进行jvm的初始化
	initJvm(cp)

	//执行合约
	// -jar contract action para0 para1 ...
	//tx2Exec := "PrintArgs debug chain33-jvm I am an Engineer !"
	tx2Exec := append([]string{"-jar", contract, action}, para...)
	argc, argv := buildJavaArgument(tx2Exec)
	defer freeArgument(argc, argv)
	var exception1DPtr *C.char
	exception := &exception1DPtr
	result := C.JLI_Exec_Contract(argc, argv, exception, TX_EXEC_JOB, (*C.char)(unsafe.Pointer(jvmHandleGo)))
	if int(result) != JLI_SUCCESS {
		exInfo := C.GoString(*exception)
		defer C.free(unsafe.Pointer(*exception))
		log.Debug("adapter::runJava", "java exception", exInfo)
		return errors.New(exInfo)
	}
	return nil
}

func initJvm(cp string) {
	if jvm_init_alreay {
		return
	}
	const_cp := C.CString(cp)
	defer C.free(unsafe.Pointer(const_cp))
	result := C.JLI_Create_TxJVM(const_cp)
	if int(result) != JLI_SUCCESS {
		panic("Failed to init JLI_Init_JVM")
	}
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
