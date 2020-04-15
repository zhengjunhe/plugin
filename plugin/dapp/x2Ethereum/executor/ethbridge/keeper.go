package ethbridge

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/address"
	dbm "github.com/33cn/chain33/common/db"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/executor/oracle"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

type Keeper struct {
	oracleKeeper OracleKeeper
	db           dbm.KV
}

func NewKeeper(oracleKeeper OracleKeeper, db dbm.KV) Keeper {
	return Keeper{
		oracleKeeper: oracleKeeper,
		db:           db,
	}
}

// 处理接收到的ethchain33请求
func (k Keeper) ProcessClaim(claim types.Eth2Chain33) (oracle.Status, error) {
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

// 处理经过审核的关于Lock的claim
func (k Keeper) ProcessSuccessfulClaimForLock(claim, execAddr, tokenSymbol string, accDB *account.DB) (*types2.Receipt, error) {
	var receipt *types2.Receipt
	oracleClaim, err := CreateOracleClaimFromOracleString(claim)
	if err != nil {
		elog.Error("CreateEthClaimFromOracleString", "CreateOracleClaimFromOracleString error", err)
		return nil, err
	}

	receiverAddress := oracleClaim.Chain33Receiver

	if oracleClaim.ClaimType == common.LockText {
		//铸币到相关的tokenSymbolBank账户下
		receipt, err = accDB.Mint(execAddr, int64(oracleClaim.Amount))
		if err != nil {
			return nil, err
		}
		r, err := accDB.ExecDeposit(receiverAddress, execAddr, int64(oracleClaim.Amount))
		if err != nil {
			return nil, err
		}
		receipt.KV = append(receipt.KV, r.KV...)
		receipt.Logs = append(receipt.Logs, r.Logs...)
		return receipt, nil
	}
	return nil, types.ErrInvalidClaimType
}

// 处理经过审核的关于Burn的claim
func (k Keeper) ProcessSuccessfulClaimForBurn(claim, execAddr, tokenSymbol string, accDB *account.DB) (*types2.Receipt, error) {
	receipt := new(types2.Receipt)
	oracleClaim, err := CreateOracleClaimFromOracleString(claim)
	if err != nil {
		elog.Error("CreateEthClaimFromOracleString", "CreateOracleClaimFromOracleString error", err)
		return nil, err
	}

	senderAddr := oracleClaim.Chain33Receiver

	if oracleClaim.ClaimType == common.BurnText {
		receipt, err = accDB.ExecTransfer(address.ExecAddress(tokenSymbol), senderAddr, execAddr, int64(oracleClaim.Amount))
		if err != nil {
			return nil, err
		}

		return receipt, nil
	}
	return nil, types.ErrInvalidClaimType
}

// ProcessBurn processes the burn of bridged coins from the given sender
func (k Keeper) ProcessBurn(address, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	receipt, err := accDB.ExecWithdraw(execAddr, address, amount)
	if err != nil {
		return nil, err
	}

	r, err := accDB.Burn(execAddr, amount)
	if err != nil {
		return nil, err
	}
	receipt.KV = append(receipt.KV, r.KV...)
	receipt.Logs = append(receipt.Logs, r.Logs...)
	return receipt, nil
}

// ProcessLock processes the lockup of cosmos coins from the given sender
// accDB = mavl-coins-bty-addr
func (k Keeper) ProcessLock(address, to, execAddr string, amount int64, accDB *account.DB) (*types2.Receipt, error) {
	// 转到 mavl-coins-bty-execAddr:addr
	receipt, err := accDB.ExecTransfer(address, to, execAddr, amount*1e8)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

//todo
// 对于相同的地址该如何处理?
// 现有方案是相同地址就报错
func (k Keeper) ProcessAddValidator(address string, power int64) (*types2.Receipt, error) {
	receipt := new(types2.Receipt)

	validatorMaps, err := k.oracleKeeper.GetValidatorArray()
	if err != nil && err != types2.ErrNotFound {
		return nil, err
	}

	elog.Info("ProcessLogInValidator", "pre validatorMaps", validatorMaps, "Add Address", address, "Add power", power)
	var totalPower int64
	for _, p := range validatorMaps {
		if p.Address != address {
			totalPower += p.Power
		} else {
			return nil, types.ErrAddressExists
		}
	}

	validatorMaps = append(validatorMaps, types.MsgValidator{
		Address: address,
		Power:   power,
	})
	v, _ := json.Marshal(validatorMaps)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
	totalPower += power

	totalP := types.ReceiptQueryTotalPower{
		TotalPower: totalPower,
	}
	totalPBytes, _ := json.Marshal(totalP)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalLastTotalPowerPrefix(), Value: totalPBytes})
	return receipt, nil
}

func (k Keeper) ProcessRemoveValidator(address string) (*types2.Receipt, error) {
	var exist bool
	receipt := new(types2.Receipt)

	validatorMaps, err := k.oracleKeeper.GetValidatorArray()
	if err != nil {
		return nil, err
	}

	elog.Info("ProcessLogOutValidator", "pre validatorMaps", validatorMaps, "Delete Address", address)
	var totalPower int64
	var validatorRes []types.MsgValidator
	for _, p := range validatorMaps {
		if address != p.Address {
			validatorRes = append(validatorRes, p)
			totalPower += p.Power
		} else {
			//oracle.RemoveAddrFromValidatorMap(validatorMaps, index)
			exist = true
			continue
		}
	}

	if !exist {
		return nil, types.ErrAddressNotExist
	}

	v, _ := json.Marshal(validatorRes)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
	totalP := types.ReceiptQueryTotalPower{
		TotalPower: totalPower,
	}
	totalPBytes, _ := json.Marshal(totalP)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalLastTotalPowerPrefix(), Value: totalPBytes})
	return receipt, nil
}

//这里的power指的是修改后的power
func (k Keeper) ProcessModifyValidator(address string, power int64) (*types2.Receipt, error) {
	var exist bool
	receipt := new(types2.Receipt)

	validatorMaps, err := k.oracleKeeper.GetValidatorArray()
	if err != nil {
		return nil, err
	}

	elog.Info("ProcessModifyValidator", "pre validatorMaps", validatorMaps, "Modify Address", address, "Modify power", power)
	var totalPower int64
	for index, p := range validatorMaps {
		if address != p.Address {
			totalPower += p.Power
		} else {
			validatorMaps[index].Power = power
			exist = true
			totalPower += power
		}
	}

	if !exist {
		return nil, types.ErrAddressNotExist
	}

	v, _ := json.Marshal(validatorMaps)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalValidatorMapsPrefix(), Value: v})
	totalP := types.ReceiptQueryTotalPower{
		TotalPower: totalPower,
	}
	totalPBytes, _ := json.Marshal(totalP)
	receipt.KV = append(receipt.KV, &types2.KeyValue{Key: types.CalLastTotalPowerPrefix(), Value: totalPBytes})

	return receipt, nil
}

func (k Keeper) ProcessSetConsensusNeeded(ConsensusThreshold float64) (float64, float64, error) {
	preCon := k.oracleKeeper.GetConsensusThreshold()
	k.oracleKeeper.SetConsensusThreshold(ConsensusThreshold)
	nowCon := k.oracleKeeper.GetConsensusThreshold()

	elog.Info("ProcessSetConsensusNeeded", "pre ConsensusThreshold", preCon, "now ConsensusThreshold", nowCon)

	return preCon, nowCon, nil
}
