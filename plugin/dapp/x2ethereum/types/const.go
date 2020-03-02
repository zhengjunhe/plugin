package types

var (
	LastTotalPowerKey = []byte{0x12}
	ValidatorMaps     = []byte{0x13}
)

//log for x2ethereum
const (
	TyLogMsgEthBridgeClaim = 350
	TyLogRelayRevokeCreate = 351
	TyLogRelayAccept       = 352
	TyLogRelayRevokeAccept = 353
	TyLogRelayConfirmTx    = 354
	TyLogRelayFinishTx     = 355
	TyLogRelayRcvBTCHead   = 356
)

const ModuleName = "x2ethereum"
