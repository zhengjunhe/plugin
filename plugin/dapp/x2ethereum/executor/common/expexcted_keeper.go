package common

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(process func(exported.Account) (stop bool))
	GetAccount(sdk.AccAddress) exported.Account
	SetAccount(exported.Account)
	NewAccount(exported.Account) exported.Account
}

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoins(fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	DelegateCoins(fromAdd, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	UndelegateCoins(fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error

	SubtractCoins(addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
	AddCoins(addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
}
