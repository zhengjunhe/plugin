package ethbridge

import (
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(address common.Chain33Address) authexported.Account
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(name string, amt sdk.Coins) error
	BurnCoins(name string, amt sdk.Coins) error
	SetModuleAccount(supplyexported.ModuleAccountI)
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(claim types.OracleClaim) (oracle.Status, error)
	GetProphecy(id string) (oracle.Prophecy, error)
}
