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
func (a *action) setProphecy(prophecy excutor.Prophecy) error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier
	}
	if len(prophecy.ClaimValidators) == 0 {
		return types.ErrNoClaims
	}
	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		return types.ErrinternalDB
	}
	serializedProphecyBytes, err := json.Marshal(serializedProphecy)
	if err != nil {
		return types.ErrMarshal
	}
	_ = a.db.Set([]byte(prophecy.ID), serializedProphecyBytes)
	return nil
}

// ProcessClaim TODO: validator hasn't implement
func (a *action) ProcessClaim(claim types.OracleClaim) (excutor.Status, error) {
	if strings.TrimSpace(claim.Content) == "" {
		return excutor.Status{}, types.ErrInvalidClaim
	}
	prophecy, err := a.GetProphecy(claim.ID)
	if err != nil {
		if err != types.ErrProphecyNotFound {
			return excutor.Status{}, err
		}
		prophecy = excutor.NewProphecy(claim.ID)
	} else {
		if prophecy.Status.Text == excutor.StatusText(types.EthBridgeStatus_SuccessStatusText) || prophecy.Status.Text == excutor.StatusText(types.EthBridgeStatus_FailedStatusText) {
			return excutor.Status{}, types.ErrProphecyFinalized
		}
		if prophecy.ValidatorClaims[claim.ValidatorAddress] != "" {
			return excutor.Status{}, types.ErrDuplicateMessage
		}
	}
	prophecy.AddClaim(NewChain33Address(claim.ValidatorAddress), claim.Content)
	prophecy = a.processCompletion(prophecy)
	err = a.setProphecy(prophecy)
	if err != nil {
		return excutor.Status{}, err
	}
	return prophecy.Status, nil
}

//todo
//func (a *action) checkActiveValidator(validatorAddress Chain33Address) bool {
//	validator, found := a.GetValidator(validatorAddress)
//	if !found {
//		return false
//	}
//	bondStatus := validator.GetStatus()
//	return bondStatus == sdk.Bonded
//}
//
//func (a *action) GetValidator(addr Chain33Address) ()

// processCompletion looks at a given prophecy an assesses whether the claim with the highest power on that prophecy has enough
// power to be considered successful, or alternatively, will never be able to become successful due to not enough validation power being
// left to push it over the threshold required for consensus.
func (a *action) processCompletion(prophecy types.Prophecy) types.Prophecy {
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(ctx, k.stakeKeeper)
	totalPower := k.stakeKeeper.GetLastTotalPower(ctx)
	highestConsensusRatio := float64(highestClaimPower) / float64(totalPower.Int64())
	remainingPossibleClaimPower := totalPower.Int64() - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := float64(highestPossibleClaimPower) / float64(totalPower.Int64())
	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.Text = types.EthBridgeStatus_SuccessStatusText
		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.consensusNeeded {
		prophecy.Status.Text = types.EthBridgeStatus_FailedStatusText
	}
	return prophecy
}
