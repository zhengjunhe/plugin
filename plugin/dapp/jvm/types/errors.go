package types

import "errors"

var (
	ErrDepthJvm                    = errors.New("max call depth exceeded for Jvm")
	ErrTraceLimitReachedJvm        = errors.New("the number of logs reached the specified limit for Jvm")
	ErrInsufficientBalanceJvm      = errors.New("insufficient balance for transfer for Jvm")
	ErrContractAddressCollisionJvm = errors.New("The Name of contract was used by other contract already")
	ErrAddrNotExistsJvm            = errors.New("address not exists for Jvm")
	ErrTransferBetweenContractsJvm = errors.New("transferring between contracts not supports for Jvm")
	ErrTransferBetweenEOAJvm       = errors.New("transferring between external accounts not supports for Jvm")
	ErrNoCreatorJvm                = errors.New("contract has no creator information for Jvm")
	ErrDestructJvm                 = errors.New("contract has been destructed for Jvm")

	ErrWriteProtectionJvm       = errors.New("Jvm: write protection")
	ErrReturnDataOutOfBoundsJvm = errors.New("Jvm: return data out of bounds")
	ErrExecutionRevertedJvm     = errors.New("Jvm: execution reverted")
	ErrMaxCodeSizeExceededJvm   = errors.New("Jvm: max code size exceeded")
	ErrWrongContractAddr        = errors.New("Jvm: wrong contract addr")
	ErrJvmValidationFail        = errors.New("Jvm: fail to validate byte code")
	ErrJvmNotSupported          = errors.New("Jvm: vm is not supported now")
	ErrJvmContractExecFailed    = errors.New("Jvm: contract executing failed")
	ErrAddrNotExists            = errors.New("Jvm: contract addr not exists")
	ErrContractNotExists        = errors.New("Jvm: contract not exists")
	ErrNoCreator                = errors.New("Jvm: no creator for contract")
	ErrNoCoinsAccount           = errors.New("Jvm: no coins Account")
	ErrTransferBetweenContracts = errors.New("Jvm: not allow to thansfer between contracts account")
	ErrTransferBetweenEOA       = errors.New("Jvm: not allow to thansfer between external account")
	ErrUnserialize              = errors.New("Jvm: unserialize")
	ErrCreateJvmPara            = errors.New("Jvm: wrong parameter for creating new Jvm contract")
	ErrCallJvmPara              = errors.New("Jvm: wrong parameter for creating call Jvm contract")
	ErrWrongContracName         = errors.New("Jvm: Contract name should be a-z and 0-9")
	ErrWrongContracNameLen      = errors.New("Jvm: Contract name length should within [4-16]")
	ErrNoPermission             = errors.New("Jvm: action without permission")
	ErrActionNotSupport         = errors.New("Jvm: action not support")
	ErrAbiNotMatch              = errors.New("Jvm: action data and abi not match due to contract has been updated or abiJson format wrong")
	ErrAbiNotFound              = errors.New("Jvm: abi not found due to wrong contract name")
	ErrWriteJavaClass           = errors.New("Jvm: Failed to write java class to file")
	ErrJavaExecFailed           = errors.New("Jvm: Failed to execute java contract")
	ErrGetJvmFailed             = errors.New("Jvm: Failed to get go-jvm exector")
)
