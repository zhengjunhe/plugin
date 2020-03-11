package types

import "errors"

var (
	ErrInvalidClaimType              = errors.New("invalid claim type provided")
	ErrInvalidEthSymbol              = errors.New("invalid symbol provided, symbol \"eth\" must have null address set as token contract address")
	ErrInvalidChainID                = errors.New("invalid ethereum chain id")
	ErrInvalidEthAddress             = errors.New("invalid ethereum address provided, must be a valid hex-encoded Ethereum address")
	ErrInvalidEthNonce               = errors.New("invalid ethereum nonce provided, must be >= 0")
	ErrInvalidAddress                = errors.New("invalid Chain33 address")
	ErrInvalidIdentifier             = errors.New("invalid identifier provided, must be a nonempty string")
	ErrProphecyNotFound              = errors.New("prophecy with given id not found")
	ErrProphecyGet                   = errors.New("prophecy with given id find error")
	ErrinternalDB                    = errors.New("internal error serializing/deserializing prophecy")
	ErrNoClaims                      = errors.New("cannot create prophecy without initial claim")
	ErrInvalidClaim                  = errors.New("claim cannot be empty string")
	ErrProphecyFinalized             = errors.New("prophecy already finalized")
	ErrDuplicateMessage              = errors.New("already processed message from validator for this id")
	ErrMinimumConsensusNeededInvalid = errors.New("minimum consensus proportion of validator staking power must be > 0 and <= 1")
	ErrInvalidValidator              = errors.New("validator is invalid")
	ErrUnknownAddress                = errors.New("module account does not exist")
	ErrLogOutPowerIsTooBig           = errors.New("log out power is more than which this address saves")
)

//common
var (
	ErrSetKV = errors.New("Set KV error")
)
