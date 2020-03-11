package ethbridge

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	supplyKeeper SupplyKeeper
	oracleKeeper OracleKeeper
	db           dbm.KV
}

// NewKeeper creates new instances of the oracle Keeper
func NewKeeper(supplyKeeper SupplyKeeper, oracleKeeper OracleKeeper, db dbm.KV) Keeper {
	return Keeper{
		supplyKeeper: supplyKeeper,
		oracleKeeper: oracleKeeper,
		db:           db,
	}
}

// ProcessClaim processes a new claim coming in from a validator
func (k Keeper) ProcessClaim(claim types.EthBridgeClaim) (oracle.Status, error) {
	oracleClaim, err := CreateOracleClaimFromEthClaim(claim)
	if err != nil {
		elog.Error("CreateEthClaimFromOracleString", "CreateOracleClaimFromOracleString error", err)
		return oracle.Status{}, err
	}

	status, err := k.oracleKeeper.ProcessClaim(oracleClaim)
	if err != nil {
		return oracle.Status{}, err
	}
	return status, nil
}

// ProcessSuccessfulClaim processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessSuccessfulClaim(claim, execAddr string, accDB *account.DB) (*types2.Receipt, error) {
	var receipt *types2.Receipt
	oracleClaim, err := CreateOracleClaimFromOracleString(claim)
	if err != nil {
		elog.Error("CreateEthClaimFromOracleString", "CreateOracleClaimFromOracleString error", err)
		return nil, err
	}

	receiverAddress := oracleClaim.Chain33Receiver

	if oracleClaim.ClaimType == LockText {
		receipt, err = k.supplyKeeper.MintCoins(int64(oracleClaim.Amount), types.ModuleName, execAddr, accDB)
		if err != nil {
			return nil, err
		}
	}
	r, err := k.supplyKeeper.SendCoinsFromModuleToAccount(types.ModuleName, receiverAddress, execAddr, int64(oracleClaim.Amount), accDB)
	if err != nil {
		panic(err)
	}
	receipt.KV = append(receipt.KV, r.KV...)
	receipt.Logs = append(receipt.Logs, r.Logs...)
	return receipt, nil
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(address, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	receipt, err := k.supplyKeeper.SendCoinsFromAccountToModule(address, types.ModuleName, execAddr, amount, accDB)
	if err != nil {
		return nil, err
	}
	r, err := k.supplyKeeper.BurnCoins(amount, types.ModuleName, execAddr, accDB)
	if err != nil {
		panic(err)
	}
	receipt.KV = append(receipt.KV, r.KV...)
	receipt.Logs = append(receipt.Logs, r.Logs...)
	return receipt, nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
func (k Keeper) ProcessLock(address, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	receipt, err := k.supplyKeeper.SendCoinsFromAccountToModule(address, types.ModuleName, execAddr, amount, accDB)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

//todo
// 对于相同的地址该如何处理?
// 现有方案是相同地址power覆盖处理
func (k Keeper) ProcessLogInValidator(address string, power float64) (*types2.Receipt, error) {
	//flg 为true时，则说明有相同地址
	flg := false
	receipt := new(types2.Receipt)

	validatorMaps, err := k.oracleKeeper.GetValidatorArray()
	if err != nil && err != types2.ErrNotFound {
		return nil, err
	}

	elog.Info("ProcessLogInValidator", "pre validatorMaps", validatorMaps, "Add Address", address, "Add power", power)
	var totalPower float64
	for _, p := range validatorMaps {
		if p.Address != address {
			v, _ := json.Marshal(p)
			receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
			totalPower += p.Power
		} else {
			v, _ := json.Marshal(oracle.ValidatorMap{
				Address: address,
				Power:   power,
			})
			receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
			totalPower += power
			flg = true
		}

	}

	if !flg {
		v, _ := json.Marshal(oracle.ValidatorMap{
			Address: address,
			Power:   power,
		})
		receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
		totalPower += power
	}

	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalLastTotalPowerPrefix(), Value: common.Float64ToBytes(totalPower)})

	//validatorMaps = append(validatorMaps, oracle.ValidatorMap{
	//	Address: address,
	//	Power:   power,
	//})
	//validatorMapsBytes, err := json.Marshal(validatorMaps)
	//if err != nil {
	//	return nil, types2.ErrMarshal
	//}
	//err = k.db.Set(types.ValidatorMapsKey, validatorMapsBytes)
	//if err != nil {
	//	return nil, types.ErrSetKV
	//}
	//
	//err = k.db.Set(types.LastTotalPowerKey, common.Float64ToBytes(totalPower))
	//if err != nil {
	//	return nil, types.ErrSetKV
	//}

	return receipt, nil
}

//LogOut的power代表的是减少的量
func (k Keeper) ProcessLogOutValidator(address string, power float64) (*types2.Receipt, error) {
	receipt := new(types2.Receipt)
	//flg 为true时，则说明有相同地址
	flg := false

	validatorMaps, err := k.oracleKeeper.GetValidatorArray()
	if err != nil {
		return nil, err
	}

	elog.Info("ProcessLogOutValidator", "pre validatorMaps", validatorMaps, "Delete Address", address, "Delete power", power)
	var totalPower float64
	for index, p := range validatorMaps {
		if address != p.Address {
			bz, _ := json.Marshal(types.MsgValidator{
				Address: address,
				Power:   power,
			})
			receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: bz})
		} else {
			if p.Power < power {
				return nil, types.ErrLogOutPowerIsTooBig
			} else if p.Power == power {
				oracle.RemoveAddrFromValidatorMap(validatorMaps, index)
				continue
			} else {
				p.Power -= power
				bz, _ := json.Marshal(types.MsgValidator{
					Address: address,
					Power:   p.Power,
				})
				receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: bz})
			}
			flg = true
		}
		totalPower += p.Power
	}

	if flg {
		return nil, types.ErrLogOutAddressNotExist
	}

	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalLastTotalPowerPrefix(), Value: common.Float64ToBytes(totalPower)})

	//validatorMapsBytes, err := json.Marshal(validatorMaps)
	//if err != nil {
	//	return nil, types2.ErrMarshal
	//}
	//_ = k.db.Set(types.ValidatorMapsKey, validatorMapsBytes)

	return receipt, nil
}

func (k Keeper) ProcessSetConsensusNeeded(consensusNeeded float64) (*types2.Receipt, float64, float64, error) {
	receipt := new(types2.Receipt)

	preCon := k.oracleKeeper.GetConsensusNeeded()

	k.oracleKeeper.SetConsensusNeeded(consensusNeeded)

	nowCon := k.oracleKeeper.GetConsensusNeeded()

	elog.Info("ProcessSetConsensusNeeded", "pre ConsensusNeeded", preCon, "now ConsensusNeeded", nowCon)

	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalConsensusNeededPrefix(), Value: common.Float64ToBytes(consensusNeeded)})

	_ = k.db.Set(types.ConsensusNeededKey, common.Float64ToBytes(consensusNeeded))

	return receipt, preCon, nowCon, nil
}
