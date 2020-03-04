package oracle

import (
	"encoding/json"
	dbm "github.com/33cn/chain33/common/db"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"strings"
)

type Keeper struct {
	db              dbm.KV
	consensusNeeded float64 // The minimum % of stake needed to sign claims in order for consensus to occur
}

func NewKeeper(db dbm.KV, consensusNeeded float64) Keeper {
	if consensusNeeded <= 0 || consensusNeeded > 1 {
		panic(types.ErrMinimumConsensusNeededInvalid)
	}
	return Keeper{
		db:              db,
		consensusNeeded: consensusNeeded,
	}
}

func (k Keeper) GetProphecy(id string) (Prophecy, error) {
	if id == "" {
		return NewEmptyProphecy(), types.ErrInvalidIdentifier
	}

	bz, err := k.db.Get([]byte(id))
	if err != nil {
		return NewEmptyProphecy(), types.ErrProphecyGet
	} else if bz == nil {
		return NewEmptyProphecy(), types.ErrProphecyNotFound
	}
	var dbProphecy DBProphecy
	err = json.Unmarshal(bz, &dbProphecy)
	if err != nil {
		return NewEmptyProphecy(), types2.ErrUnmarshal
	}

	deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
	if err != nil {
		return NewEmptyProphecy(), types.ErrinternalDB
	}
	return deSerializedProphecy, nil
}

// setProphecy saves a prophecy with an initial claim
func (k Keeper) setProphecy(prophecy Prophecy) error {
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
		return types2.ErrMarshal
	}

	err = k.db.Set([]byte(prophecy.ID), serializedProphecyBytes)
	if err != nil {
		return types.ErrSetKV
	}
	return nil
}

func (k Keeper) ProcessClaim(claim types.OracleClaim) (Status, error) {
	activeValidator := k.checkActiveValidator(claim.ValidatorAddress)
	if !activeValidator {
		return Status{}, types.ErrInvalidValidator
	}
	if strings.TrimSpace(claim.Content) == "" {
		return Status{}, types.ErrInvalidClaim
	}
	prophecy, err := k.GetProphecy(claim.ID)
	if err != nil {
		if err != types.ErrProphecyNotFound {
			return Status{}, err
		}
		prophecy = NewProphecy(claim.ID)
	} else {
		if prophecy.Status.Text == StatusText(types.EthBridgeStatus_SuccessStatusText) || prophecy.Status.Text == StatusText(types.EthBridgeStatus_FailedStatusText) {
			return Status{}, types.ErrProphecyFinalized
		}
		if prophecy.ValidatorClaims[claim.ValidatorAddress] != "" {
			return Status{}, types.ErrDuplicateMessage
		}
	}
	prophecy.AddClaim(claim.ValidatorAddress, claim.Content)
	prophecy, err = k.processCompletion(prophecy)
	err = k.setProphecy(prophecy)
	if err != nil {
		return Status{}, err
	}
	return prophecy.Status, nil
}

func (k Keeper) checkActiveValidator(validatorAddress string) bool {
	validatorMap, err := k.GetValidatorArray()
	if err != nil {
		return false
	}

	for _, v := range validatorMap {
		if v.Address == validatorAddress {
			return true
		}
	}
	return false
}

// 计算该prophecy是否达标
func (k Keeper) processCompletion(prophecy Prophecy) (Prophecy, error) {
	address2power := make(map[string]float64)
	validatorArrays, err := k.GetValidatorArray()
	if err != nil {
		return prophecy, err
	}
	for _, validator := range validatorArrays {
		address2power[validator.Address] = validator.Power
	}
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(address2power)
	totalPower, err := k.GetLastTotalPower()
	if err != nil {
		return prophecy, err
	}
	highestConsensusRatio := highestClaimPower / totalPower
	remainingPossibleClaimPower := totalPower - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := highestPossibleClaimPower / totalPower
	if highestConsensusRatio >= k.consensusNeeded {
		prophecy.Status.Text = StatusText(types.EthBridgeStatus_SuccessStatusText)
		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.consensusNeeded {
		prophecy.Status.Text = StatusText(types.EthBridgeStatus_FailedStatusText)
	}
	return prophecy, nil
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower() (power float64, err error) {
	b, err := k.db.Get(types.LastTotalPowerKey)
	if err != nil && err != types2.ErrNotFound {
		return 0, err
	} else if err == types2.ErrNotFound {
		return 0, nil
	}
	err = json.Unmarshal(b, &power)
	if err != nil {
		return 0, types2.ErrUnmarshal
	}
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower() error {
	var totalPower float64
	validatorArrays, err := k.GetValidatorArray()
	if err != nil {
		return err
	}
	for _, validator := range validatorArrays {
		totalPower += validator.Power
	}
	err = k.db.Set(types.LastTotalPowerKey, common.Float64ToBytes(totalPower))
	if err != nil {
		return types.ErrSetKV
	}
	return nil
}

func (k Keeper) GetValidatorArray() ([]ValidatorMap, error) {
	validatorsBytes, err := k.db.Get(types.ValidatorMapsKey)
	if err != nil {
		return nil, err
	}
	var validatorArrays []ValidatorMap
	err = json.Unmarshal(validatorsBytes, &validatorArrays)
	if err != nil {
		return nil, types2.ErrUnmarshal
	}
	return validatorArrays, nil
}

type ValidatorMap struct {
	Address string
	Power   float64
}

func RemoveAddrFromValidatorMap(validatorMap []ValidatorMap, index int) []ValidatorMap {
	return append(validatorMap[:index], validatorMap[index+1:]...)
}

func (k Keeper) SetConsensusNeeded(consensusNeeded float64) {
	k.consensusNeeded = consensusNeeded
	return
}
