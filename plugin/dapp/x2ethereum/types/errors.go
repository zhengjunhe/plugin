package types

import "errors"

var (
	ErrInvalidClaimType  = errors.New("invalid claim type provided")
	ErrInvalidEthSymbol  = errors.New("invalid symbol provided, symbol \"eth\" must have null address set as token contract address")
	ErrJSONMarshalling   = errors.New("error marshalling JSON for this claim")
	ErrInvalidChainID    = errors.New("invalid ethereum chain id")
	ErrInvalidEthAddress = errors.New("invalid ethereum address provided, must be a valid hex-encoded Ethereum address")
	ErrInvalidEthNonce   = errors.New("invalid ethereum nonce provided, must be >= 0")
	ErrInvalidAddress    = errors.New("invalid Chain33 address")
)
