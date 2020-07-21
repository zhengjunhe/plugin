package executor

import (
	"bytes"
	"encoding/hex"
	"regexp"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

// Query_CheckContractNameExist 确认是否存在该Jvm合约，
func (jvm *JVMExecutor) Query_CheckContractNameExist(in *jvmTypes.CheckJVMContractNameReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	jvm.prepareQueryContext([]byte(jvmTypes.JvmX))
	return jvm.checkContractNameExists(in)
}

// Query_CreateJvmContract 创建Jvm 合约
func (jvm *JVMExecutor) Query_CreateJvmContract(in *jvmTypes.CreateContractReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}

	execer := types.GetRealExecName([]byte(in.Name))
	if bytes.HasPrefix(execer, []byte(jvmTypes.UserJvmX)) {
		execer = execer[len(jvmTypes.UserJvmX):]
	}

	execerStr := string(execer)
	nameReg, err := regexp.Compile(jvmTypes.NameRegExp)
	if nil != err || !nameReg.MatchString(execerStr) {
		return nil, jvmTypes.ErrWrongContracName
	}

	if len(execerStr) > 16 || len(execerStr) < 4 {
		return nil, jvmTypes.ErrWrongContracNameLen
	}
	action := &jvmTypes.JVMContractAction{
		Value: &jvmTypes.JVMContractAction_CreateJvmContract{
			CreateJvmContract: &jvmTypes.CreateJvmContract{
				Code:     in.Code,
				Abi:      in.Abi,
				Name:     jvm.GetAPI().GetConfig().ExecName(jvmTypes.UserJvmX + execerStr),
				Note:     in.Note,
			},
		},
		Ty: jvmTypes.CreateJvmContractAction,
	}

	createRsp, err := createRawJvmTx(jvm.GetAPI().GetConfig(), action, jvmTypes.JvmX, in.Fee)
	result := hex.EncodeToString(types.Encode(createRsp))
	relpydata := &types.ReplyString{Data: result}
	return relpydata, err
}

// Query_CallJvmContract 调用创建的合约
func (jvm *JVMExecutor) Query_CallJvmContract(in *jvmTypes.CallContractReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	jvm.prepareQueryContext([]byte(jvmTypes.JvmX))

	contractName := in.Name
	if !bytes.Contains([]byte(contractName), []byte(jvmTypes.UserJvmX)) {
		contractName = jvmTypes.UserJvmX + contractName
	}

	action := &jvmTypes.JVMContractAction{
		Value: &jvmTypes.JVMContractAction_CallJvmContract{
			CallJvmContract: &jvmTypes.CallJvmContract{
				Note:       in.Note,
				ActionName: in.ActionName,
				ActionData: []byte(in.DataInJson),
			},
		},
		Ty: jvmTypes.CallJvmContractAction,
	}
	createRsp, err := createRawJvmTx(jvm.GetAPI().GetConfig(), action, contractName, in.Fee)
	result := hex.EncodeToString(types.Encode(createRsp))
	replydata := &types.ReplyString{Data: result}
	return replydata, err
}

// Query_CreateJvmContract 创建Jvm 合约
func (jvm *JVMExecutor) Query_UpdateJvmContract(in *jvmTypes.UpdateContractReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}

	execer := types.GetRealExecName([]byte(in.Name))
	if bytes.HasPrefix(execer, []byte(jvmTypes.UserJvmX)) {
		execer = execer[len(jvmTypes.UserJvmX):]
	}

	execerStr := string(execer)
	nameReg, err := regexp.Compile(jvmTypes.NameRegExp)
	if err != nil || !nameReg.MatchString(execerStr) {
		return nil, jvmTypes.ErrWrongContracName
	}

	if len(execerStr) > 16 || len(execerStr) < 4 {
		return nil, jvmTypes.ErrWrongContracNameLen
	}

	action := &jvmTypes.JVMContractAction{
		Value: &jvmTypes.JVMContractAction_UpdateJvmContract{
			UpdateJvmContract: &jvmTypes.UpdateJvmContract{
				Code:     in.Code,
				Abi:      in.Abi,
				Name:     jvm.GetAPI().GetConfig().ExecName(jvmTypes.UserJvmX + execerStr),
				Note:     in.Note,
			},
		},
		Ty: jvmTypes.UpdateJvmContractAction,
	}
	createRsp, err := createRawJvmTx(jvm.GetAPI().GetConfig(), action, jvmTypes.JvmX, in.Fee)
	result := hex.EncodeToString(types.Encode(createRsp))
	relpydata := &types.ReplyString{Data: result}
	return relpydata, err
}

func createRawJvmTx(chain33Config *types.Chain33Config, action proto.Message, JvmName string, fee int64) (*types.Transaction, error) {
	tx := &types.Transaction{
		Execer:  []byte(chain33Config.ExecName(JvmName)),
		Payload: types.Encode(action),
		To:      address.ExecAddress(chain33Config.ExecName(JvmName)),
	}
	tx, err := types.FormatTx(chain33Config, string(tx.Execer), tx)
	if err != nil {
		return nil, err
	}
	if tx.Fee < fee {
		tx.Fee = fee
	}
	return tx, nil
}
