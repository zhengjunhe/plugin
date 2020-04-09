package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/types"
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
				ClaimValidators: dbPD.ClaimValidators,
				ValidatorClaims: dbPD.ValidatorClaims,
			}
			return prophecy, nil
		}
	}
	return nil, types2.ErrInvalidProphecyID
}

func (x *x2ethereum) Query_GetValidators(in *types2.QueryValidatorsParams) (types.Message, error) {
	validatorsKey := types2.CalValidatorMapsPrefix()

	var v []*types2.MsgValidator
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
				val[0] = vv
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
		var totalPower int64
		for _, vv := range v {
			totalPower += vv.Power
		}
		validatorsRes.Validators = v
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
	err = json.Unmarshal(totalPowerBytes, &totalPower)
	if err != nil {
		return nil, types.ErrUnmarshal
	}
	return totalPower, nil
}

func (x *x2ethereum) Query_GetConsensusThreshold(in *types2.QueryConsensusThresholdParams) (types.Message, error) {
	consensus := &types2.ReceiptSetConsensusThreshold{}
	consensusKey := types2.CalConsensusThresholdPrefix()

	consensusBytes, err := x.GetStateDB().Get(consensusKey)
	if err != nil {
		elog.Error("Query_GetConsensusNeeded", "GetConsensusNeeded Err", err)
		return nil, err
	}
	err = json.Unmarshal(consensusBytes, &consensus)
	if err != nil {
		return nil, types.ErrUnmarshal
	}
	return consensus, nil
}

func (x *x2ethereum) Query_GetSymbolTotalAmount(in *types2.QuerySymbolAssetsParams) (types.Message, error) {
	symbolAmount := &types2.ReceiptQuerySymbolAssets{}
	symbolAmountKey := types2.CalTokenSymbolTotalAmountPrefix(in.TokenSymbol, types2.DirectionType[in.Direction])

	totalAmountBytes, err := x.GetStateDB().Get(symbolAmountKey)
	if err != nil {
		elog.Error("Query_GetSymbolTotalAmount", "GetSymbolTotalAmount Err", err)
		return nil, err
	}
	err = json.Unmarshal(totalAmountBytes, &symbolAmount)
	if err != nil {
		return nil, types.ErrUnmarshal
	}
	return symbolAmount, nil
}

func (x *x2ethereum) Query_GetSymbolTotalAmountByTxType(in *types2.QuerySymbolAssetsByTxTypeParams) (types.Message, error) {
	symbolAmount := &types2.ReceiptQuerySymbolAssetsByTxType{}
	symbolAmountKey := types2.CalTokenSymbolTotalLockOrBurnAmount(in.TokenSymbol, types2.DirectionType[in.Direction], in.TxType)

	totalAmountBytes, err := x.GetLocalDB().Get(symbolAmountKey)
	if err != nil {
		elog.Error("Query_GetSymbolTotalAmountByTxType", "GetSymbolTotalAmountByTxType Err", err)
		return nil, err
	}
	err = types.Decode(totalAmountBytes, symbolAmount)
	if err != nil {
		return nil, types.ErrUnmarshal
	}
	return symbolAmount, nil
}
