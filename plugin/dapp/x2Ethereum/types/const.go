package types

var (
	LastTotalPowerKey  = []byte{0x12}
	ValidatorMapsKey   = []byte{0x13}
	ConsensusNeededKey = []byte{0x14}
)

// log for x2ethereum
// log类型id值
const (
	TyUnknownLog = iota + 100
	TyEthBridgeClaimLog
	TyMsgBurnLog
	TyMsgLockLog
	TyMsgLogInValidatorLog
	TyMsgLogOutValidatorLog
	TyMsgSetConsensusNeededLog
)

// action类型id和name，这些常量可以自定义修改
const (
	TyUnknowAction = iota + 100
	TyEthBridgeClaimAction
	TyMsgBurnAction
	TyMsgLockAction
	TyMsgLogInValidatorAction
	TyMsgLogOutValidatorAction
	TyMsgSetConsensusNeededAction

	NameEthBridgeClaimAction        = "EthBridgeClaim"
	NameMsgBurnAction               = "MsgBurn"
	NameMsgLockAction               = "MsgLock"
	NameMsgLogInValidatorAction     = "MsgLogInValidator"
	NameMsgLogOutValidatorAction    = "MsgLogOutValidator"
	NameMsgSetConsensusNeededAction = "MsgSetConsensusNeeded"
)

const ModuleName = "x2ethereum"

const DefaultConsensusNeeded = 0.7
