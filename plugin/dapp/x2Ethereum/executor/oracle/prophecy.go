package oracle

import (
	"encoding/json"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

type Prophecy struct {
	ID              string              `json:"id"`
	Status          Status              `json:"status"`
	ClaimValidators map[string][]string `json:"claim_validators"` //This is a mapping from a claim to the list of validators that made that claim.
	ValidatorClaims map[string]string   `json:"validator_claims"` //This is a mapping from a validator bech32 address to their claim
}

// NewProphecy returns a new Prophecy, initialized in pending status with an initial claim
func NewProphecy(id string) Prophecy {
	return Prophecy{
		ID:              id,
		Status:          NewStatus(StatusText(types.EthBridgeStatus_PendingStatusText), ""),
		ClaimValidators: make(map[string][]string),
		ValidatorClaims: make(map[string]string),
	}
}

// NewEmptyProphecy returns a blank prophecy, used with errors
func NewEmptyProphecy() Prophecy {
	return NewProphecy("")
}

func NewProphecyByproto(prophecy types.Prophecy) Prophecy {
	claimValidators := make(map[string][]string)
	for claim, addresses := range prophecy.ClaimValidators {
		addressArrays := make([]string, 0)
		for _, addr := range addresses.ClaimValidator {
			addressArrays = append(addressArrays, addr)
		}
		claimValidators[claim] = addressArrays
	}
	return Prophecy{
		ID:              prophecy.ID,
		Status:          NewStatus(StatusText(prophecy.Status.Text), prophecy.Status.FinalClaim),
		ClaimValidators: claimValidators,
		ValidatorClaims: prophecy.ValidatorClaims,
	}
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
	var claimValidators map[string][]string
	err := json.Unmarshal(dbProphecy.ClaimValidators, &claimValidators)
	if err != nil {
		return Prophecy{}, err
	}

	var validatorClaims map[string]string
	err = json.Unmarshal(dbProphecy.ValidatorClaims, &validatorClaims)
	if err != nil {
		return Prophecy{}, err
	}

	return Prophecy{
		ID:              dbProphecy.ID,
		Status:          dbProphecy.Status,
		ClaimValidators: claimValidators,
		ValidatorClaims: validatorClaims,
	}, nil
}

// AddClaim adds a given claim to this prophecy
func (prophecy Prophecy) AddClaim(validator string, claim string) {
	claimValidators := prophecy.ClaimValidators[claim]
	prophecy.ClaimValidators[claim] = append(claimValidators, validator)

	prophecy.ValidatorClaims[validator] = claim
}

// FindHighestClaim looks through all the existing claims on a given prophecy. It adds up the total power across
// all claims and returns the highest claim, power for that claim, and total power claimed on the prophecy overall.
func (prophecy Prophecy) FindHighestClaim(validators map[string]float64) (string, float64, float64) {
	totalClaimsPower := float64(0)
	highestClaimPower := float64(-1)
	highestClaim := ""
	for claim, validatorAddrs := range prophecy.ClaimValidators {
		claimPower := float64(0)
		for _, validatorAddr := range validatorAddrs {
			validatorPower := validators[validatorAddr]
			claimPower += validatorPower
		}
		totalClaimsPower += claimPower
		if claimPower > highestClaimPower {
			highestClaimPower = claimPower
			highestClaim = claim
		}
	}
	return highestClaim, highestClaimPower, totalClaimsPower
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
