package main

//#cgo CFLAGS: -DTARGET_ARCH_x86 -DLINUX -DTARGET_COMPILER_gcc -DFULL_VERSION=1.8.0_262-internal-hezhengjun_2020_08_04_23_54-b00 -DJDK_MAJOR_VERSION=1 -DJDK_MINOR_VERSION=8
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/share/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/share/vm/prims
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/cpu/x86/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/share/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/solaris/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/share/native/common
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/solaris/native/common
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/share/bin
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/solaris/bin
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/linux/bin
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
	name = flag.String("n", "jvmtest", "classname")
)

func main() {
	flag.Parse()
	fmt.Println("Begin to run jvm test with the name")
	fmt.Println(*name)

	if *cp == "" {
		*cp = ".:/home/hezhengjun/work/c_code/openjdk/dragonwell8/build/linux-x86_64-normal-server-release/jdk"
	}

	//解析并构建合约argc,argv
	contract2Exec := "java HelloWorld"
	argc, argv := buildJavaArgument(contract2Exec)
	defer freeArgument(argc, argv)

	FULL_VERSION := C.CString("1.8.0_262-internal-hezhengjun_2020_08_04_23_54-b00")
	DOT_VERSION := C.CString("1.8")
	const_progname := C.CString("java")
	const_launcher := C.CString("openjdk")
	defer C.free(unsafe.Pointer(FULL_VERSION))
	defer C.free(unsafe.Pointer(DOT_VERSION))
	defer C.free(unsafe.Pointer(const_progname))
	defer C.free(unsafe.Pointer(const_launcher))

	nil2dPtr := C.GetNil2dPtr()
	C.JLI_Init_JVM(argc, argv, 0, nil2dPtr, 0, nil2dPtr, FULL_VERSION, DOT_VERSION, const_progname, const_launcher, 0, 1, 0, 0)
	//if result != C.int(0) {
	//	fmt.Println("Failed to init jvm")
	//	return
	//}
	fmt.Println("Succeed to init jvm")

	//执行合约１
	contract2Exec = "PrintArgs debug chain33-jvm I am an Engineer !"
	argc1, argv1 := buildJavaArgument(contract2Exec)
	result := C.JLI_Exec_Contract(argc1, argv1)
	fmt.Println(contract2Exec, " is executed with result: ", result)
	defer freeArgument(argc1, argv1)

	//执行合约２
	contract2Exec = "HelloWorld"
	argc, argv = buildJavaArgument(contract2Exec)
	result = C.JLI_Exec_Contract(argc, argv)
	fmt.Println(contract2Exec, " is executed with result: ", result)
	defer freeArgument(argc, argv)

	C.JLI_Debug_Contract()

	//执行结束，将jvm进行销毁
	C.JLI_Detroy_JVM()

	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Println("\n- sig goroute is running")
		<-sigs
		fmt.Println("\n- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()

	for {
		fmt.Println("- Sleeping 10 sec")
		time.Sleep(10 * time.Second)
	}

	//TODO:检查内存泄露
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
