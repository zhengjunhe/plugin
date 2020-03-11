package common

// Keeper of the supply store
type Keeper struct {
	moduleAddresses map[string]string
}

// NewKeeper creates a new Keeper instance
func NewKeeper(moduleAddressesStr map[string]string) Keeper {
	return Keeper{
		moduleAddresses: moduleAddressesStr,
	}
}
