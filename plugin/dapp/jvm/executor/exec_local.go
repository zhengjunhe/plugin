package executor

import (
	"github.com/33cn/chain33/types"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

// ExecLocal_CreateJvmContract 本地执行创建Jvm合约
func (jvm *JVMExecutor) ExecLocal_CreateJvmContract(payload *jvmTypes.CreateJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execLocal(tx, receipt, index)
}

// ExecLocal_CallJvmContract 本地执行调用Jvm合约
func (jvm *JVMExecutor) ExecLocal_CallJvmContract(payload *jvmTypes.CallJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execLocal(tx, receipt, index)
}

// ExecLocal_UpdateJvmContract 本地执行更新Jvm合约
func (jvm *JVMExecutor) ExecLocal_UpdateJvmContract(payload *jvmTypes.UpdateJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execLocal(tx, receipt, index)
}

func (Jvm *JVMExecutor) execLocal(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	for _, logItem := range receipt.Logs {
		if jvmTypes.TyLogLocalDataJvm == logItem.Ty {
			data := logItem.Log
			var localData jvmTypes.ReceiptLocalData
			err := types.Decode(data, &localData)
			if err != nil {
				return set, err
			}
			set.KV = append(set.KV, &types.KeyValue{Key: localData.Key, Value: localData.CurValue})
			log.Debug("execLocal_setkv", "key=", string(localData.Key))
		}
	}

	return set, nil
}
