package executor

import (
	"github.com/33cn/chain33/types"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

// ExecDelLocal_CreateJvmContract 本地撤销执行创建Jvm合约
func (jvm *JVMExecutor) ExecDelLocal_CreateJvmContract(payload *jvmTypes.CreateJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execDelLocal(tx, receipt, index)
}

// ExecDelLocal_CallJvmContract 本地撤销执行调用Jvm合约
func (jvm *JVMExecutor) ExecDelLocal_CallJvmContract(payload *jvmTypes.CallJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execDelLocal(tx, receipt, index)
}

// ExecDelLocal_UpdateJvmContract 本地撤销执行更新Jvm合约
func (jvm *JVMExecutor) ExecDelLocal_UpdateJvmContract(payload *jvmTypes.UpdateJvmContract, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return jvm.execDelLocal(tx, receipt, index)
}

func (Jvm *JVMExecutor) execDelLocal(tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	set := &types.LocalDBSet{}
	// 需要将Exec中生成的合约状态变更信息从localdb中恢复
	for _, logItem := range receipt.Logs {
		if jvmTypes.TyLogLocalDataJvm == logItem.Ty {
			data := logItem.Log
			var localData jvmTypes.ReceiptLocalData
			err := types.Decode(data, &localData)
			if err != nil {
				return set, err
			}
			set.KV = append(set.KV, &types.KeyValue{Key: localData.Key, Value: localData.PreValue})
		}
	}

	return set, nil
}
