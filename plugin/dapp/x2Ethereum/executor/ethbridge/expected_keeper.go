package ethbridge

import (
	"github.com/33cn/chain33/account"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(senderModule, recipientAddr, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error)
	SendCoinsFromModuleToModule(senderModule, recipientModule, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error)
	SendCoinsFromAccountToModule(senderAddr, recipientModule, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error)
	MintCoins(amt int64, moduleName, execAddr string, accDB *account.DB) (*types2.Receipt, error)
	BurnCoins(amt int64, moduleName, execAddr string, accDB *account.DB) (*types2.Receipt, error)
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(claim types.OracleClaim) (oracle.Status, error)
	GetProphecy(id string) (oracle.Prophecy, error)
	GetValidatorArray() ([]oracle.ValidatorMap, error)
	SetConsensusNeeded(consensusNeeded float64)
	GetConsensusNeeded() float64
}
