package oracle

import (
	"encoding/json"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

type Prophecy struct {
	ID              string                   `json:"id"`
	Status          Status                   `json:"status"`
	ClaimValidators []*types.ClaimValidators `json:"claim_validators"` //This is a mapping from a claim to the list of validators that made that claim.
	ValidatorClaims []*types.ValidatorClaims `json:"validator_claims"` //This is a mapping from a validator bech32 address to their claim
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(id string) Prophecy {
	return Prophecy{
		ID:              id,
		Status:          NewStatus(StatusText(types.EthBridgeStatus_PendingStatusText), ""),
		ClaimValidators: *new([]*types.ClaimValidators),
		ValidatorClaims: *new([]*types.ValidatorClaims),
	}
}

// NewEmptyProphecy returns a blank prophecy, used with errors
func NewEmptyProphecy() Prophecy {
	return NewProphecy("")
}

// DBProphecy is what the prophecy becomes when being saved to the database. Tendermint/Amino does not support maps so we must serialize those variables into bytes.
type DBProphecy struct {
	ID              string `json:"id"`
	Status          Status `json:"status"`
	ClaimValidators []byte `json:"claim_validators"`
	ValidatorClaims []byte `json:"validator_claims"`
}

// SerializeForDB serializes a prophecy into a DBProphecy
func (prophecy Prophecy) SerializeForDB() (DBProphecy, error) {
	claimValidators, err := json.Marshal(prophecy.ClaimValidators)
	if err != nil {
		return DBProphecy{}, err
	}

	validatorClaims, err := json.Marshal(prophecy.ValidatorClaims)
	if err != nil {
		return DBProphecy{}, err
	}

	return DBProphecy{
		ID:              prophecy.ID,
		Status:          prophecy.Status,
		ClaimValidators: claimValidators,
		ValidatorClaims: validatorClaims,
	}, nil
}

// DeserializeFromDB deserializes a DBProphecy into a prophecy
func (dbProphecy DBProphecy) DeserializeFromDB() (Prophecy, error) {
	claimValidators := new([]*types.ClaimValidators)
	err := json.Unmarshal(dbProphecy.ClaimValidators, &claimValidators)
	if err != nil {
		return Prophecy{}, err
	}

	validatorClaims := new([]*types.ValidatorClaims)
	err = json.Unmarshal(dbProphecy.ValidatorClaims, &validatorClaims)
	if err != nil {
		return Prophecy{}, err
	}

	return Prophecy{
		ID:              dbProphecy.ID,
		Status:          dbProphecy.Status,
		ClaimValidators: *claimValidators,
		ValidatorClaims: *validatorClaims,
	}, nil
}

// AddClaim adds a given claim to this prophecy
func (prophecy *Prophecy) AddClaim(validator string, claim string) {
	claimValidators := new(types.StringMap)
	if len(prophecy.ClaimValidators) == 0 {
		prophecy.ClaimValidators = append(prophecy.ClaimValidators, &types.ClaimValidators{
			Claim: claim,
			Validators: &types.StringMap{
				Validators: []string{validator},
			},
		})
	} else {
		for index, cv := range prophecy.ClaimValidators {
			if cv.Claim == claim {
				claimValidators = cv.Validators
				prophecy.ClaimValidators[index].Validators = AddToStringMap(claimValidators, validator)
				break
			}
		}
	}

	if len(prophecy.ValidatorClaims) == 0 {
		prophecy.ValidatorClaims = append(prophecy.ValidatorClaims, &types.ValidatorClaims{
			Validator: validator,
			Claim:     claim,
		})
	} else {
		for index, vc := range prophecy.ValidatorClaims {
			if vc.Validator == validator {
				prophecy.ValidatorClaims[index].Claim = claim
				break
			} else {
				prophecy.ValidatorClaims = append(prophecy.ValidatorClaims, &types.ValidatorClaims{
					Validator: validator,
					Claim:     claim,
				})
			}
		}
	}

}

// FindHighestClaim looks through all the existing claims on a given prophecy. It adds up the total power across
// all claims and returns the highest claim, power for that claim, and total power claimed on the prophecy overall.
func (prophecy *Prophecy) FindHighestClaim(validators map[string]int64) (string, float64, float64) {
	totalClaimsPower := int64(0)
	highestClaimPower := int64(-1)
	highestClaim := ""
	for _, claimValidators := range prophecy.ClaimValidators {
		claimPower := int64(0)
		for _, validatorAddr := range claimValidators.Validators.Validators {
			validatorPower := validators[validatorAddr]
			claimPower += validatorPower
		}
		totalClaimsPower += claimPower
		if claimPower > highestClaimPower {
			highestClaimPower = claimPower
			highestClaim = claimValidators.Claim
		}
	}
	return highestClaim, float64(highestClaimPower), float64(totalClaimsPower)
}

// Status is a struct that contains the status of a given prophecy
type Status struct {
	Text       StatusText `json:"text"`
	FinalClaim string     `json:"final_claim"`
}

// NewStatus returns a new Status with the given data contained
func NewStatus(text StatusText, finalClaim string) Status {
	return Status{
		Text:       text,
		FinalClaim: finalClaim,
	}
}

func AddToStringMap(in *types.StringMap, validator string) *types.StringMap {
	inStringMap := append(in.GetValidators(), validator)
	stringMapRes := new(types.StringMap)
	stringMapRes.Validators = inStringMap
	return stringMapRes
}
