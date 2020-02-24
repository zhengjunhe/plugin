package executor

import (
	"github.com/33cn/chain33/common/crypto"
	"time"
)

type Validator struct {
	OperatorAddress         Chain33Address `json:"operator_address" yaml:"operator_address"`       // address of the validator's operator; bech encoded in JSON
	ConsPubKey              crypto.PubKey  `json:"consensus_pubkey" yaml:"consensus_pubkey"`       // the consensus public key of the validator; bech encoded in JSON
	Jailed                  bool           `json:"jailed" yaml:"jailed"`                           // has the validator been jailed from bonded status?
	Status                  sdk.BondStatus `json:"status" yaml:"status"`                           // validator status (bonded/unbonding/unbonded)
	Tokens                  sdk.Int        `json:"tokens" yaml:"tokens"`                           // delegated tokens (incl. self-delegation)
	DelegatorShares         sdk.Dec        `json:"delegator_shares" yaml:"delegator_shares"`       // total shares issued to a validator's delegators
	Description             Description    `json:"description" yaml:"description"`                 // description terms for the validator
	UnbondingHeight         int64          `json:"unbonding_height" yaml:"unbonding_height"`       // if unbonding, height at which this validator has begun unbonding
	UnbondingCompletionTime time.Time      `json:"unbonding_time" yaml:"unbonding_time"`           // if unbonding, min time for the validator to complete unbonding
	Commission              Commission     `json:"commission" yaml:"commission"`                   // commission parameters
	MinSelfDelegation       sdk.Int        `json:"min_self_delegation" yaml:"min_self_delegation"` // validator's self declared minimum self delegation
}
