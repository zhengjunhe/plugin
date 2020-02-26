package ethbridge

import (
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"github.com/gogo/protobuf/codec"
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
	oracleClaim, err := types.CreateOracleClaimFromEthClaim(k.cdc, claim)
	if err != nil {
		return oracle.Status{}, err
	}

	status, sdkErr := k.oracleKeeper.ProcessClaim(oracleClaim)
	if sdkErr != nil {
		return oracle.Status{}, sdkErr
	}
	return status, nil
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(claim string) error {
	oracleClaim, err := types.CreateOracleClaimFromOracleString(claim)
	if err != nil {
		return err
	}

	receiverAddress := oracleClaim.CosmosReceiver

	if oracleClaim.ClaimType == types.LockText {
		err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, oracleClaim.Amount)
		if err != nil {
			return err
		}
	}
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiverAddress, oracleClaim.Amount)
	if err != nil {
		panic(err)
	}
	return nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
	if err != nil {
		return err
	}
	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, amount)
	if err != nil {
		panic(err)
	}
	return nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(ctx sdk.Context, cosmosSender sdk.AccAddress, amount sdk.Coins) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, cosmosSender, types.ModuleName, amount)
	if err != nil {
		return err
	}
	return nil
}
