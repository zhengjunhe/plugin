package main

//#cgo CFLAGS: -DTARGET_ARCH_x86 -DLINUX -DTARGET_COMPILER_gcc -DFULL_VERSION=1.8.0_262-internal-hezhengjun_2020_08_04_23_54-b00 -DJDK_MAJOR_VERSION=1 -DJDK_MINOR_VERSION=8
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/share/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/share/vm/prims
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/hotspot/src/cpu/x86/vm
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/share/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/solaris/javavm/export
//#cgo CFLAGS: -I/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/share/native/common
//#cgo CFLAGS: -II/home/hezhengjun/work/c_code/openjdk/dragonwell8/jdk/src/solaris/native/common
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
//char ** GetNil2dPtr() {
//	return (char **)0;
//}
import "C"

import (
	"flag"
	"fmt"
	"strings"
	"unsafe"
)

var (
	cp  = flag.String("c", "", "classpath")
	myTestClass  = flag.String("i", "", "my test class")
	name = flag.String("n", "jvmtest", "classname")
	execPara = flag.String("e", "java -version", "exec para")
)


func main() {
	fmt.Println("Begin to run jvm test with the name")
	fmt.Println(name)

	if *cp == "" {
		*cp = ".:/home/hezhengjun/work/c_code/openjdk/dragonwell8/build/linux-x86_64-normal-server-release/jdk"
	}

	//解析并构建合约argc,argv
	paraSlice := parseArgv(*execPara)
	argc := C.int(len(paraSlice))
	testSize := unsafe.Pointer(C.CString("testPtrSize"))
	defer C.free(testSize)

	nil2dPtr := C.GetNil2dPtr()
	argv0 := (**C.char)(C.malloc(C.ulong(argc * C.GetPtrSize())))
	if argv0 == nil2dPtr {
		panic("Failed to malloc for argv")
	}
	//argv [argc]*C.char
	for i, para := range paraSlice {
		paraCstr := C.CString(para)
		C.SetPtr(argv0, paraCstr, C.int(i))
		defer C.free(unsafe.Pointer(paraCstr))
	}

	FULL_VERSION := C.CString("1.8.0_262-internal-hezhengjun_2020_08_04_23_54-b00")
	DOT_VERSION := C.CString("1.8")
	const_progname := C.CString("java")
	const_launcher := C.CString("openjdk")
	defer C.free(unsafe.Pointer(FULL_VERSION))
	defer C.free(unsafe.Pointer(DOT_VERSION))
	defer C.free(unsafe.Pointer(const_progname))
	defer C.free(unsafe.Pointer(const_launcher))

	C.JLI_Launch(argc, argv0, 0, nil2dPtr, 0, nil2dPtr, FULL_VERSION, DOT_VERSION, const_progname, const_launcher, 0, 1, 0, 0);
}

func parseArgv(para string) []string {
	var paraParsed []string
	para = strings.TrimSpace(para)
	paraParsed = strings.Split(para, " ")
	fmt.Println("the Parsed para is", paraParsed)
	return paraParsed
}
