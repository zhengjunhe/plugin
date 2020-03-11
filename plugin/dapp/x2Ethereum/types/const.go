package types

var (
	ProphecyKey        = []byte("prefix_for_ProphecyKey")
	EthBridgeClaimKey  = []byte("prefix_for_EthBridgeClaim")
	LockKey            = []byte("prefix_for_LockKey")
	BurnKey            = []byte("prefix_for_BurnKey")
	LastTotalPowerKey  = []byte("prefix_for_LastTotalPowerKey")
	ValidatorMapsKey   = []byte("prefix_for_ValidatorMapsKey")
	ConsensusNeededKey = []byte("prefix_for_ConsensusNeededKey")
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

const ModuleName = "x2ethereumBank"

const DefaultConsensusNeeded = 0.7
