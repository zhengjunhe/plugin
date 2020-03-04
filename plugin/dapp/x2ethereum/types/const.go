package types

var (
	LastTotalPowerKey  = []byte{0x12}
	ValidatorMapsKey   = []byte{0x13}
	ConsensusNeededKey = []byte{0x14}
)

//log for x2ethereum
const (
	TyLogMsgEthBridgeClaim     = 350
	TyLogMsgLock               = 351
	TyLogMsgBurn               = 352
	TyLogMsgLogInValidator     = 353
	TyLogMsgLogOutValidator    = 354
	TyLogMsgSetConsensusNeeded = 355
)

const ModuleName = "x2ethereum"

const DefaultConsensusNeeded = 0.7

const AddressPowerPrefix = "power-"
