package common

// Keeper of the supply store
type Keeper struct {
	bk              BankKeeper
	moduleAddresses map[string]string
}

// NewKeeper creates a new Keeper instance
func NewKeeper(bk BankKeeper, moduleAddressesStr map[string]string) Keeper {
	return Keeper{
		bk:              bk,
		moduleAddresses: moduleAddressesStr,
	}
}
