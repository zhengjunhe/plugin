package types

import "errors"

var (
	ErrContractAddressCollisionJvm = errors.New("The Name of contract was used by other contract already")
	ErrMaxCodeSizeExceededJvm   = errors.New("Jvm: max code size exceeded")
	ErrWrongContractAddr        = errors.New("Jvm: wrong contract addr")
	ErrAddrNotExists            = errors.New("Jvm: contract addr not exists")
	ErrContractNotExists        = errors.New("Jvm: contract not exists")
	ErrNoCreator                = errors.New("Jvm: no creator for contract")
	ErrWrongContracName         = errors.New("Jvm: Contract name should be a-z and 0-9")
	ErrWrongContracNameLen      = errors.New("Jvm: Contract name length should within [4-16]")
	ErrNoPermission             = errors.New("Jvm: action without permission")
	ErrActionNotSupport         = errors.New("Jvm: action not support")
	ErrWriteJavaClass           = errors.New("Jvm: Failed to write java class to file")
	ErrGetJvmFailed             = errors.New("Jvm: Failed to get go-jvm exector")
	ErrSelLocalNotAllowed       = errors.New("Jvm: Set Local DB only Allowed during tx exec")
	ErrJvmCodeString            = errors.New("Jvm: Wrong jvm code string")
)
