package executor

import (
	"encoding/json"
	excutor "github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"strings"
)

// GetProphecy gets the entire prophecy data struct for a given id
func (a *action) GetProphecy(id string) (excutor.Prophecy, error) {
	if id == "" {
		return excutor.NewEmptyProphecy(), types.ErrInvalidIdentifier
	}
	//store to localDB
	bz, err := a.db.Get([]byte(id))
	if err != nil {
		return excutor.NewEmptyProphecy(), types.ErrProphecyGet
	} else if bz == nil {
		return excutor.NewEmptyProphecy(), types.ErrProphecyNotFound
	}
	var dbProphecy excutor.DBProphecy
	err = json.Unmarshal(bz, &dbProphecy)
	if err != nil {
		return excutor.NewEmptyProphecy(), types.ErrUnmarshal
	}

	deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return excutor.NewEmptyProphecy(), types.ErrinternalDB
	}
	return deSerializedProphecy, nil
}

// setProphecy saves a prophecy with an initial claim
func (a *action) setProphecy(ctx sdk.Context, prophecy types.Prophecy) sdk.Error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier(k.Codespace())
	}
	if len(prophecy.ClaimValidators) == 0 {
		return types.ErrNoClaims(k.Codespace())
	}
	store := ctx.KVStore(k.storeKey)
	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		return types.ErrInternalDB(k.Codespace(), err)
	}
	store.Set([]byte(prophecy.ID), k.cdc.MustMarshalBinaryBare(serializedProphecy))
	return nil
}

// ProcessClaim TODO: write description
func (a *executor.action) ProcessClaim(ctx sdk.Context, claim types.Claim) (types.Status, sdk.Error) {
	activeValidator := k.checkActiveValidator(ctx, claim.ValidatorAddress)
	if !activeValidator {
		return types.Status{}, types.ErrInvalidValidator(k.Codespace())
	}
	if strings.TrimSpace(claim.Content) == "" {
		return types.Status{}, types.ErrInvalidClaim(k.Codespace())
	}
	prophecy, err := k.GetProphecy(ctx, claim.ID)
	if err != nil {
		if err.Code() != types.CodeProphecyNotFound {
			return types.Status{}, err
		}
		prophecy = types.NewProphecy(claim.ID)
	} else {
		if prophecy.Status.Text == types.SuccessStatusText || prophecy.Status.Text == types.FailedStatusText {
			return types.Status{}, types.ErrProphecyFinalized(k.Codespace())
		}
		if prophecy.ValidatorClaims[claim.ValidatorAddress.String()] != "" {
			return types.Status{}, types.ErrDuplicateMessage(k.Codespace())
		}
	}
	prophecy.AddClaim(claim.ValidatorAddress, claim.Content)
	prophecy = k.processCompletion(ctx, prophecy)
	err = k.setProphecy(ctx, prophecy)
	if err != nil {
		return types.Status{}, err
	}
	return prophecy.Status, nil
}

func (a *executor.action) checkActiveValidator(ctx sdk.Context, validatorAddress sdk.ValAddress) bool {
	validator, found := k.stakeKeeper.GetValidator(ctx, validatorAddress)
	if !found {
		return false
	}
	bondStatus := validator.GetStatus()
	return bondStatus == sdk.Bonded
}

// processCompletion looks at a given prophecy an assesses whether the claim with the highest power on that prophecy has enough
// power to be considered successful, or alternatively, will never be able to become successful due to not enough validation power being
// left to push it over the threshold required for consensus.
func (a *executor.action) processCompletion(ctx sdk.Context, prophecy types.Prophecy) types.Prophecy {
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(ctx, k.stakeKeeper)
	totalPower := k.stakeKeeper.GetLastTotalPower(ctx)
	highestConsensusRatio := float64(highestClaimPower) / float64(totalPower.Int64())
	remainingPossibleClaimPower := totalPower.Int64() - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := float64(highestPossibleClaimPower) / float64(totalPower.Int64())
	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.Text = types.SuccessStatusText
		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.consensusNeeded {
		prophecy.Status.Text = types.FailedStatusText
	}
	return prophecy
}
