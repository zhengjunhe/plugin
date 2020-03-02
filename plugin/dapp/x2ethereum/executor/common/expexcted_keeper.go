package common

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(fromAddr string, toAddr string, amt int64) error

	SubtractCoins(addr string, amt int64) (amount, error)
	AddCoins(addr string, amt int64) (amount, error)
}
