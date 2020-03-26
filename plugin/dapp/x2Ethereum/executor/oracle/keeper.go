package oracle

import (
	"encoding/json"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"strings"
)

var (
	//日志
	olog = log.New("module", "oracle")
)

type Keeper struct {
	db                 dbm.KV
	ConsensusThreshold float64 // The minimum % of stake needed to sign claims in order for consensus to occur
}

func NewKeeper(db dbm.KV, ConsensusThreshold float64) Keeper {
	if ConsensusThreshold <= 0 || ConsensusThreshold > 1 {
		panic(types.ErrMinimumConsensusNeededInvalid)
	}
	return Keeper{
		db:                 db,
		ConsensusThreshold: ConsensusThreshold,
	}
}

func (k *Keeper) GetProphecy(id string) (Prophecy, error) {
	if id == "" {
		return NewEmptyProphecy(), types.ErrInvalidIdentifier
	}

	bz, err := k.db.Get(types.CalProphecyPrefix())
	if err != nil && err != types2.ErrNotFound {
		return NewEmptyProphecy(), types.ErrProphecyGet
	} else if err == types2.ErrNotFound {
		return NewEmptyProphecy(), types.ErrProphecyNotFound
	}
	var dbProphecys []DBProphecy
	var dbProphecy DBProphecy
	err = json.Unmarshal(bz, &dbProphecys)
	if err != nil {
		return NewEmptyProphecy(), types2.ErrUnmarshal
	}

	var exist bool
	for _, p := range dbProphecys {
		if p.ID == id {
			dbProphecy = p
			exist = true
			break
		}
	}

	if exist {
		deSerializedProphecy, err := dbProphecy.DeserializeFromDB()
		if err != nil {
			return NewEmptyProphecy(), types.ErrinternalDB
		}
		return deSerializedProphecy, nil
	} else {
		return NewEmptyProphecy(), types.ErrProphecyNotFound
	}
}

// setProphecy saves a prophecy with an initial claim
func (k *Keeper) setProphecy(prophecy Prophecy) error {
	err := k.checkProphecy(prophecy)
	if err != nil {
		return err
	}

	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		return types.ErrinternalDB
	}

	bz, err := k.db.Get(types.CalProphecyPrefix())
	if err != nil && err != types2.ErrNotFound {
		return types.ErrProphecyGet
	}

	var dbProphecys []DBProphecy
	if err != types2.ErrNotFound {
		err = json.Unmarshal(bz, &dbProphecys)
		if err != nil {
			return types2.ErrUnmarshal
		}
	}

	var exist bool
	for index, dbP := range dbProphecys {
		if dbP.ID == serializedProphecy.ID {
			exist = true
			dbProphecys[index] = serializedProphecy
			break
		}
	}
	if !exist {
		dbProphecys = append(dbProphecys, serializedProphecy)
	}

	serializedProphecyBytes, err := json.Marshal(dbProphecys)
	if err != nil {
		return types2.ErrMarshal
	}

	err = k.db.Set(types.CalProphecyPrefix(), serializedProphecyBytes)
	if err != nil {
		return types.ErrSetKV
	}
	return nil
}

// modifyProphecy saves a modified prophecy
func (k *Keeper) modifyProphecy(prophecy Prophecy) error {
	err := k.checkProphecy(prophecy)
	if err != nil {
		return err
	}

	serializedProphecy, err := prophecy.SerializeForDB()
	if err != nil {
		return types.ErrinternalDB
	}

	bz, err := k.db.Get(types.CalProphecyPrefix())
	if err != nil && err != types2.ErrNotFound {
		return types.ErrProphecyGet
	}

	var dbProphecys []DBProphecy
	if err != types2.ErrNotFound {
		err = json.Unmarshal(bz, &dbProphecys)
		if err != nil {
			return types2.ErrUnmarshal
		}
	}

	for index, dbP := range dbProphecys {
		if dbP.ID == serializedProphecy.ID {
			dbProphecys[index] = serializedProphecy
			break
		}
	}

	serializedProphecyBytes, err := json.Marshal(dbProphecys)
	if err != nil {
		return types2.ErrMarshal
	}

	err = k.db.Set(types.CalProphecyPrefix(), serializedProphecyBytes)
	if err != nil {
		return types.ErrSetKV
	}
	return nil
}

func (k *Keeper) checkProphecy(prophecy Prophecy) error {
	if prophecy.ID == "" {
		return types.ErrInvalidIdentifier
	}
	if len(prophecy.ClaimValidators) == 0 {
		return types.ErrNoClaims
	}
	return nil
}

func (k *Keeper) ProcessClaim(claim types.OracleClaim) (Status, error) {
	activeValidator := k.checkActiveValidator(claim.ValidatorAddress)
	if !activeValidator {
		return Status{}, types.ErrInvalidValidator
	}
	if strings.TrimSpace(claim.Content) == "" {
		return Status{}, types.ErrInvalidClaim
	}
	var claimContent types.OracleClaimContent
	err := json.Unmarshal([]byte(claim.Content), &claimContent)
	if err != nil {
		return Status{}, types2.ErrUnmarshal
	}
	prophecy, err := k.GetProphecy(claim.ID)
	if err != nil {
		if err != types.ErrProphecyNotFound {
			return Status{}, err
		}
		prophecy = NewProphecy(claim.ID)
	} else {
		if claimContent.ClaimType == common.LockText {
			if prophecy.Status.Text == StatusText(types.EthBridgeStatus_SuccessStatusText) {
				return Status{}, types.ErrProphecyFinalized
			}
			for _, vc := range prophecy.ValidatorClaims {
				if vc.Validator == claim.ValidatorAddress && vc.Claim != "" {
					return Status{}, types.ErrDuplicateMessage
				}
			}
		} else if claimContent.ClaimType == common.BurnText {
			if prophecy.Status.Text == StatusText(types.EthBridgeStatus_WithdrawedStatusText) {
				return Status{}, types.ErrProphecyFinalized
			}
		}
	}
	prophecy.AddClaim(claim.ValidatorAddress, claim.Content)
	prophecy, err = k.processCompletion(&prophecy, claimContent.ClaimType)
	err = k.setProphecy(prophecy)
	if err != nil {
		return Status{}, err
	}
	return prophecy.Status, nil
}

func (k *Keeper) checkActiveValidator(validatorAddress string) bool {
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
func (k *Keeper) processCompletion(prophecy *Prophecy, claimType int64) (Prophecy, error) {
	address2power := make(map[string]int64)
	validatorArrays, err := k.GetValidatorArray()
	if err != nil {
		return *prophecy, err
	}
	for _, validator := range validatorArrays {
		address2power[validator.Address] = validator.Power
	}
	highestClaim, highestClaimPower, totalClaimsPower := prophecy.FindHighestClaim(address2power)
	totalPower, err := k.GetLastTotalPower()
	if err != nil {
		return *prophecy, err
	}
	highestConsensusRatio := highestClaimPower / totalPower
	remainingPossibleClaimPower := totalPower - totalClaimsPower
	highestPossibleClaimPower := highestClaimPower + remainingPossibleClaimPower
	highestPossibleConsensusRatio := highestPossibleClaimPower / totalPower
	olog.Info("processCompletion", "highestConsensusRatio", highestConsensusRatio, "ConsensusThreshold", k.ConsensusThreshold, "highestPossibleConsensusRatio", highestPossibleConsensusRatio)
	if highestConsensusRatio >= k.ConsensusThreshold {
		if claimType == common.LockText {
			prophecy.Status.Text = StatusText(types.EthBridgeStatus_SuccessStatusText)
		} else {
			prophecy.Status.Text = StatusText(types.EthBridgeStatus_WithdrawedStatusText)
		}

		prophecy.Status.FinalClaim = highestClaim
	} else if highestPossibleConsensusRatio < k.ConsensusThreshold {
		prophecy.Status.Text = StatusText(types.EthBridgeStatus_FailedStatusText)
	}
	return *prophecy, nil
}

// Load the last total validator power.
func (k *Keeper) GetLastTotalPower() (float64, error) {
	b, err := k.db.Get(types.CalLastTotalPowerPrefix())
	if err != nil && err != types2.ErrNotFound {
		return 0, err
	} else if err == types2.ErrNotFound {
		return 0, nil
	}
	var powers types.ReceiptQueryTotalPower
	err = json.Unmarshal(b, &powers)
	if err != nil {
		return 0, types2.ErrUnmarshal
	}
	return float64(powers.TotalPower), nil
}

// Set the last total validator power.
func (k *Keeper) SetLastTotalPower() error {
	var totalPower int64
	validatorArrays, err := k.GetValidatorArray()
	if err != nil {
		return err
	}
	for _, validator := range validatorArrays {
		totalPower += validator.Power
	}
	totalP := types.ReceiptQueryTotalPower{
		TotalPower: totalPower,
	}
	totalPBytes, _ := json.Marshal(totalP)
	err = k.db.Set(types.CalLastTotalPowerPrefix(), totalPBytes)
	if err != nil {
		return types.ErrSetKV
	}
	return nil
}

func (k *Keeper) GetValidatorArray() ([]types.MsgValidator, error) {
	validatorsBytes, err := k.db.Get(types.CalValidatorMapsPrefix())
	if err != nil {
		return nil, err
	}
	var validatorArrays []types.MsgValidator
	err = json.Unmarshal(validatorsBytes, &validatorArrays)
	if err != nil {
		return nil, types2.ErrUnmarshal
	}
	return validatorArrays, nil
}

func RemoveAddrFromValidatorMap(validatorMap []types.MsgValidator, index int) []types.MsgValidator {
	return append(validatorMap[:index], validatorMap[index+1:]...)
}

func (k *Keeper) SetConsensusThreshold(ConsensusThreshold float64) {
	k.ConsensusThreshold = ConsensusThreshold
	olog.Info("SetConsensusNeeded", "nowConsensusNeeded", k.ConsensusThreshold)
	return
}

func (k *Keeper) GetConsensusThreshold() float64 {
	return k.ConsensusThreshold
}
