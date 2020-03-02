package ethbridge

import (
	"github.com/33cn/chain33/account"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	supplyKeeper SupplyKeeper
	oracleKeeper OracleKeeper
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(supplyKeeper SupplyKeeper, oracleKeeper OracleKeeper) Keeper {
	return Keeper{
		supplyKeeper: supplyKeeper,
		oracleKeeper: oracleKeeper,
	}
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(claim types.EthBridgeClaim) (oracle.Status, error) {
	oracleClaim, err := CreateOracleClaimFromEthClaim(claim)
	if err != nil {
		return oracle.Status{}, err
	}

	status, err := k.oracleKeeper.ProcessClaim(oracleClaim)
	if err != nil {
		return oracle.Status{}, err
	}
	return status, nil
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(claim, execAddr string, accDB *account.DB) (*types2.Receipt, error) {
	var receipt *types2.Receipt
	oracleClaim, err := CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return nil, err
	}

	receiverAddress := oracleClaim.Chain33Receiver

	if oracleClaim.ClaimType == LockText {
		receipt, err = k.supplyKeeper.MintCoins(int64(oracleClaim.Amount), types.ModuleName, execAddr, accDB)
		if err != nil {
			return nil, err
		}
	}
	r, err := k.supplyKeeper.SendCoinsFromModuleToAccount(types.ModuleName, receiverAddress, execAddr, int64(oracleClaim.Amount), accDB)
	if err != nil {
		panic(err)
	}
	receipt.KV = append(receipt.KV, r.KV...)
	receipt.Logs = append(receipt.Logs, r.Logs...)
	return receipt, nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(address, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	receipt, err := k.supplyKeeper.SendCoinsFromAccountToModule(address, types.ModuleName, execAddr, amount, accDB)
	if err != nil {
		return nil, err
	}
	r, err := k.supplyKeeper.BurnCoins(amount, types.ModuleName, execAddr, accDB)
	if err != nil {
		panic(err)
	}
	receipt.KV = append(receipt.KV, r.KV...)
	receipt.Logs = append(receipt.Logs, r.Logs...)
	return receipt, nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(address, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	receipt, err := k.supplyKeeper.SendCoinsFromAccountToModule(address, types.ModuleName, execAddr, amount, accDB)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
