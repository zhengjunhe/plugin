package main

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
import "C"

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

var (
	cp   = flag.String("c", "", "classpath")
	name = flag.String("n", "jvm_run_tx_query", "classname")
	tx_job = C.int(0)
	query_job = C.int(1)
)

func main() {
	flag.Parse()
	fmt.Println("Begin to run jvm test with the name")
	fmt.Println(*name)

	if *cp == "" {
		*cp = ".:/home/hezhengjun/work/c_code/openjdk/jvm4chain33/build/linux-x86_64-normal-server-release/images/j2sdk-image/lib:/home/hezhengjun/work/c_code/openjdk/jvm4chain33/build/linux-x86_64-normal-server-release/images/j2sdk-image/jre/lib:com/fuzamei/chain33/*"
	}
	const_cp := C.CString(*cp)
	defer C.free(unsafe.Pointer(const_cp))

	if result := C.JLI_Create_JVM(const_cp); int(result) != 0 {
		panic("Failed to call JLI_Create_JVM")
	}
	fmt.Println("Succeed to create jvm via JLI_Create_JVM")

	exceptionInfo := get2DPtr()
	jvmGo := const_cp

	//contract2Exec := "-jar HelloWorld.jar"
	//execContract(contract2Exec, jvmGo, exceptionInfo)

	contract2Exec := "TxAndQuery0 random"
	execContract(contract2Exec, jvmGo, exceptionInfo)

	//contract2Query := "TxAndQuery0 get query0 value0"
	//queryContract(contract2Query, jvmGo, exceptionInfo)
	//contract2Query = "TxAndQuery1 get query1 value1"
	//queryContract(contract2Query, jvmGo, exceptionInfo)
	//contract2Query = "TxAndQuery1 get query2 value2"
	//queryContract(contract2Query, jvmGo, exceptionInfo)

	i := 0
	for ; i < 4; i++ {
		j := i
		go func(int) {
			index := j & 0x01
			contract2Query := fmt.Sprintf("TxAndQuery%d get query%d value%d", index, j, j)
			fmt.Println("*** ^^^^ To query for:", contract2Query)
			queryContract(contract2Query, jvmGo, exceptionInfo)
			for {
				fmt.Println("Query Go Routine- Sleeping 1 sec in goroutine:", j)
				time.Sleep(5 * time.Second)
			}
		}(j)
	}

	time.Sleep(2 * time.Second)

	for ; i < 8; i++ {
		j := i
		go func(int) {
			index := j & 0x01
			contract2Query := fmt.Sprintf("TxAndQuery%d get query%d value%d", index, j, j)
			fmt.Println("*** ^^^^ To query for:", contract2Query)
			queryContract(contract2Query, jvmGo, exceptionInfo)
			for {
				fmt.Println("Query Go Routine- Sleeping 1 sec in goroutine:", j)
				time.Sleep(5 * time.Second)
			}
		}(j)
	}

	//contract2Exec = "-jar PrintArgs.jar debug chain33-jvm I am an Engineer !"
	//execContract(contract2Exec, jvmGo, exceptionInfo)
	//
	//contract2Exec = "-jar TestQuery.jar random"
	//execContract(contract2Exec, jvmGo, exceptionInfo)

	//contract2Query := "-jar TestClasspath.jar cpget"
	//queryContract(contract2Query, jvmGo, exceptionInfo)
	//
	//contract2Query = "-jar TestQuery.jar get"
	//queryContract(contract2Query, jvmGo, exceptionInfo)

	//

	//
	////contract2Exec = "-jar TestQuery.jar get"
	////queryContract(contract2Exec, jvmGo, exceptionInfo)
	//
	//contract2Exec = "-jar TestQuery.jar getl"
	//queryContract(contract2Exec, jvmGo, exceptionInfo)
	//
	//contract2Exec = "-jar TestDB.jar get"
	//queryContract(contract2Exec, jvmGo, exceptionInfo)


	//contract2Exec := "-jar HelloWorld.jar"
	//argc0, argv0 := buildJavaArgument(contract2Exec)
	//result := C.JLI_Exec_Contract(argc0, argv0, exceptionInfo, tx_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc0, argv0)

	//执行合约１
	//contract2Exec = "PrintArgs debug chain33-jvm I am an Engineer !"
	//argc1, argv1 := buildJavaArgument(contract2Exec)
	//result = C.JLI_Exec_Contract(argc1, argv1, exceptionInfo, tx_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc1, argv1)
	//
	////执行合约２
	//contract2Exec := "HelloWorld"
	//argc2, argv2 := buildJavaArgument(contract2Exec)
	//result := C.JLI_Exec_Contract(argc2, argv2, exceptionInfo, tx_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc2, argv2)
	//
	//contract2Exec = "-jar TestQuery.jar currentheight"
	//argc_exe, argv_exe := buildJavaArgument(contract2Exec)
	//result = C.JLI_Exec_Contract(argc_exe, argv_exe, exceptionInfo, tx_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc_exe, argv_exe)
	//
	////查询合约 1
	//contract2Exec = "-jar TestQuery.jar get"
	//argc3, argv3 := buildJavaArgument(contract2Exec)
	//result = C.JLI_Exec_Contract(argc3, argv3, exceptionInfo, query_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc3, argv3)
	//
	//contract2Exec = "-jar TestQuery.jar getl"
	//argc4, argv4 := buildJavaArgument(contract2Exec)
	//result = C.JLI_Exec_Contract(argc4, argv4, exceptionInfo, query_job, jvmGo)
	//fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc4, argv4)

	//执行结束，将jvm进行销毁
	//C.JLI_Detroy_JVM()

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Println("\n- sig goroute is running")
		<-sigs
		fmt.Println("\n- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()

	for {
		fmt.Println("Main Goroutine- Sleeping 2 sec")
		time.Sleep(2 * time.Second)
	}

	//TODO:检查内存泄露
}

func execContract(contract2Exec string, jvmGo *C.char, exceptionInfo **C.char) {
	//contract2Exec = "-jar TestQuery.jar getl"
	argc4, argv4 := buildJavaArgument(contract2Exec)
	result := C.JLI_Exec_Contract(argc4, argv4, exceptionInfo, tx_job, jvmGo)
	fmt.Println(contract2Exec, " is executed with result: ", result)
	defer freeArgument(argc4, argv4)
}

func queryContract(contract2Exec string, jvmGo *C.char, exceptionInfo **C.char) {
	//contract2Exec = "-jar TestQuery.jar getl"
	argc4, argv4 := buildJavaArgument(contract2Exec)
	result := C.JLI_Exec_Contract(argc4, argv4, exceptionInfo, query_job, jvmGo)
	fmt.Println(contract2Exec, " is executed with result: ", result)
	//defer freeArgument(argc4, argv4)
}

func buildJavaArgument(execPara string) (C.int, **C.char) {
	//解析并构建合约argc,argv
	paraSlice := parseArgv(execPara)
	argc := C.int(len(paraSlice))
	nil2dPtr := C.GetNil2dPtr()
	argv := (**C.char)(C.malloc(C.ulong(argc * C.GetPtrSize())))
	if argv == nil2dPtr {
		panic("Failed to malloc for argv")
	}
	//argv [argc]*C.char
	for i, para := range paraSlice {
		paraCstr := C.CString(para)
		C.SetPtr(argv, paraCstr, C.int(i))
	}
	return argc, argv
}

func get2DPtr() **C.char {
	nil2dPtr := C.GetNil2dPtr()
	ptr2D := (**C.char)(C.malloc(C.ulong(C.GetPtrSize())))
	if ptr2D == nil2dPtr {
		panic("Failed to malloc for argv")
	}

	return ptr2D
}

func freeArgument(argc C.int, argv **C.char) {
	C.FreeArgv(argc, argv)
}

func parseArgv(para string) []string {
	var paraParsed []string
	para = strings.TrimSpace(para)
	paraParsed = strings.Split(para, " ")
	fmt.Println("the Parsed para is", paraParsed)
	return paraParsed
}
