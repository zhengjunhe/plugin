package types

import (
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/log/log15"
)

const (
	CheckNameExistsFunc       = "CheckContractNameExist"
	EstimateGasJvm           = "EstimateGasJvm"
	JvmDebug                 = "JvmDebug"
	JvmGetAbi                = "JvmGetAbi"
	ConvertJSON2Abi           = "ConvertJson2Abi"
	JvmGetContractTable      = "JvmGetContractTable"
	JvmFuzzyGetContractTable = "JvmFuzzyGetContractTable"
	Success                   = int(0)
	Exception_Fail            = int(1)
	JvmX                     = "jvm"
	UserJvmX                 = "user.jvm."
	CreateJvmContractStr     = "CreateJvmContract"
	CallJvmContractStr       = "CallJvmContract"
	UpdateJvmContractStr     = "UpdateJvmContract"
	//NameRegExp             = "[a-z0-9]"^[a-z]+\[[0-9]+\]$
	NameRegExp             = "^[a-z0-9]+$"
	AccountOpFail          = false
	AccountOpSuccess       = true
	RetryNum               = int(10)
	GRPCRecSize            = 5 * 30 * 1024 * 1024
	Coin_Precision   int64 = (1e4)
	MaxCodeSize            = 2 * 1024 * 1024
)

type JvmContratOpType int
const (
	CreateJvmContractAction = 1 + iota
	CallJvmContractAction
	UpdateJvmContractAction
)

// log for Jvm
const (
	// TyLogContractDataJvm 合约代码日志
	TyLogContractDataJvm = iota + 100
	// TyLogContractStateJvm 合约状态数据日志
	TyLogContractStateJvm
	// TyLogCallContractJvm 合约调用日志
	TyLogCallContractJvm
	// TyLogStateChangeItemJvm 合约状态变化的日志
	TyLogStateChangeItemJvm
	// TyLogCreateUserJvmContract 合约创建用户的日志
	TyLogCreateUserJvmContract
	// TyLogUpdateUserJvmContract 合约更新用户的日志
	TyLogUpdateUserJvmContract
	// TyLogLocalDataJvm 合约本地数据日志
	TyLogLocalDataJvm

	// TyLogOutputItemJvm 用于Jvm合约输出可读信息的日志记录，尤其是query的相关信息
	// 为什么不将该种信息类型的获取不放置在query中呢，因为query的操作
	// 中是不含交易费的，如果碰到恶意的Jvm合约，输出无限长度的信息，
	// 会对我们的Jvm合约系统的安全性造成威胁，基于这样的考虑我们
	TyLogOutputItemJvm
)

// ContractLog 合约在日志，对应EVM中的Log指令，可以生成指定的日志信息
// 目前这些日志只是在合约执行完成时进行打印，没有其它用途
type ContractLog struct {
	// 合约地址
	Address address.Address

	// 对应交易哈希
	TxHash common.Hash

	// 日志序号
	Index int

	// 此合约提供的主题信息
	Topics []common.Hash

	// 日志数据
	Data []byte
}

// PrintLog 合约日志打印格式
func (log *ContractLog) PrintLog() {
	log15.Debug("!Contract Log!", "Contract address", log.Address.String(), "TxHash", log.TxHash.Bytes(), "Log Index", log.Index, "Log Topics", log.Topics, "Log Data", common.ToHex(log.Data))
}


