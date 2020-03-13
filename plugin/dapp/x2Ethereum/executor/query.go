package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

func (x *x2ethereum) Query_GetEthProphecy(in *types2.QueryEthProphecyParams) (types.Message, error) {
	prophecy := &types2.ReceiptEthProphecy{}
	prophecyKey := types2.CalProphecyPrefix()

	var dbProphecy []oracle.DBProphecy
	val, err := x.GetStateDB().Get(prophecyKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(val, &dbProphecy)
	if err != nil {
		return nil, types.ErrUnmarshal
	}

	for _, dbP := range dbProphecy {
		if dbP.ID == in.ID {
			dbPD, err := dbP.DeserializeFromDB()
			if err != nil {
				return nil, err
			}
			prophecy = &types2.ReceiptEthProphecy{
				ID: in.ID,
				Status: &types2.ProphecyStatus{
					Text:       types2.EthBridgeStatus(dbP.Status.Text),
					FinalClaim: dbP.Status.FinalClaim,
				},
				ClaimValidators: stringArray2StringMap(dbPD.ClaimValidators),
				ValidatorClaims: dbPD.ValidatorClaims,
			}
		}
		return prophecy, nil
	}
	return nil, types2.ErrInvalidProphecyID
}

func (x *x2ethereum) Query_GetValidators(in *types2.QueryValidatorsParams) (types.Message, error) {
	validatorsKey := types2.CalValidatorMapsPrefix()

	var v []oracle.ValidatorMap
	vBytes, err := x.GetStateDB().Get(validatorsKey)
	if err != nil {
		elog.Error("Query_GetValidators", "GetValidators Err", err)
		return nil, err
	}

	err = json.Unmarshal(vBytes, &v)
	if err != nil {
		return nil, types.ErrUnmarshal
	}

	if in.Validator != "" {
		validatorsRes := new(types2.ReceiptQueryValidator)
		for _, vv := range v {
			if vv.Address == in.Validator {
				val := make([]*types2.MsgValidator, 1)
				val[0] = &types2.MsgValidator{
					Address: vv.Address,
					Power:   vv.Power,
				}
				validatorsRes = &types2.ReceiptQueryValidator{
					Validators: val,
					TotalPower: vv.Power,
				}
				return validatorsRes, nil
			}
		}
		// 未知的地址
		return nil, types2.ErrInvalidValidator
	} else {
		validatorsRes := new(types2.ReceiptQueryValidator)
		val := make([]*types2.MsgValidator, len(v))
		var totalPower float64
		for index, vv := range v {
			val[index] = &types2.MsgValidator{
				Address: vv.Address,
				Power:   vv.Power,
			}
			totalPower += vv.Power
		}
		validatorsRes.Validators = val
		validatorsRes.TotalPower = totalPower
		return validatorsRes, nil
	}
}

func (x *x2ethereum) Query_GetTotalPower(in *types2.QueryTotalPowerParams) (types.Message, error) {
	totalPower := &types2.ReceiptQueryTotalPower{}
	totalPowerKey := types2.CalLastTotalPowerPrefix()

	totalPowerBytes, err := x.GetStateDB().Get(totalPowerKey)
	if err != nil {
		elog.Error("Query_GetTotalPower", "GetTotalPower Err", err)
		return nil, err
	}
	totalPower.TotalPower = common.BytesToFloat64(totalPowerBytes)
	return totalPower, nil
}

func (x *x2ethereum) Query_GetConsensusNeeded(in *types2.QueryConsensusNeededParams) (types.Message, error) {
	consensus := &types2.ReceiptQueryConsensusNeeded{}
	consensusKey := types2.CalConsensusNeededPrefix()

	var consensusTemp types2.MsgSetConsensusNeeded
	consensusTempBytes, err := x.GetStateDB().Get(consensusKey)
	if err != nil {
		elog.Error("Query_GetConsensusNeeded", "GetConsensusNeeded Err", err)
		return nil, err
	}
	err = json.Unmarshal(consensusTempBytes, &consensusTemp)
	if err != nil {
		return nil, types.ErrUnmarshal
	}
	consensus.ConsensusNeed = consensusTemp.ConsensusNeed
	return consensus, nil
}

func stringArray2StringMap(in map[string][]string) map[string]*types2.StringMap {
	res := make(map[string]*types2.StringMap, len(in))
	for key, value := range in {
		sm := new(types2.StringMap)
		sm.Validators = value
		res[key] = sm
	}
	return res
}
