package main

import "C"
import (
	"fmt"
)

//export CallbackFromOpenjdk
func CallbackFromOpenjdk(strInfo *C.char) {
	strDeubgInGo := C.GoString(strInfo)
	fmt.Println("Succeed to Call back from openjdk", strDeubgInGo)
}

//func main() {}
