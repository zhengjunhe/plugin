package executor

import (
	"bytes"
	"encoding/hex"
	"os"
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
func (jvm *JVMExecutor) Query_CallJvmContract(in *jvmTypes.CallJvmContract) (types.Message, error) {
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
				ActionData: in.ActionData,
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

//查询java合约状态
func (jvm *JVMExecutor) Query_JavaContract(in *jvmTypes.JVMQueryReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	jvm.prepareQueryContext([]byte(jvmTypes.JvmX))
	jvm.queryChan = make(chan QueryResult, 1)

	execer := types.GetRealExecName([]byte(in.Contract))
	if bytes.HasPrefix(execer, []byte(jvmTypes.UserJvmX)) {
		execer = execer[len(jvmTypes.UserJvmX):]
	}

	jvm.mStateDB.SetCurrentExecutorName(jvmTypes.JvmX)

	log.Debug("jvm call", "Para Query_JavaContract", in)

	contractName := in.Contract
	userJvmAddr := address.ExecAddress(contractName)
	contractAccount := jvm.mStateDB.GetAccount(userJvmAddr)
	jarPath := "./" + contractName + ".jar"
	jarFileExist := true
	//判断jar文件是否存在
	_, err := os.Stat(jarPath)
	if err != nil && !os.IsExist(err) {
		jarFileExist = false
	}

	if !jarFileExist {
		javaClassfile, err := os.OpenFile(jarPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		code := contractAccount.Data.GetCode()
		if len(code) == 0 {
			log.Error("call jvm contract ", "failed to get code from contract", contractName)
			return nil, jvmTypes.ErrWrongContractAddr
		}

		writeLen, err := javaClassfile.Write(code)
		if writeLen != len(code) {
			return nil, jvmTypes.ErrWriteJavaClass
		}
		if closeErr := javaClassfile.Close(); nil != closeErr {
			return nil, closeErr
		}
	}

	//将当前合约执行名字修改为user.jvm.xxx
	jvm.mStateDB.SetCurrentExecutorName(string(jvm.GetAPI().GetConfig().GetParaExec([]byte(contractName))))

	log.Debug("Query_JavaContract", "ContractName", contractName, "Para", in.Para)
	//2nd step: just call contract
	//在此处将gojvm指针传递到c实现的jvm中，进行回调的时候用来区分是获取数据时，使用执行db还是查询db
	_ = runJava(contractName, in.Para, jvm, TX_EXEC_JOB, jvm.GetAPI().GetConfig())

	//阻塞并等待查询结果的返回
	queryResult := <-jvm.queryChan
	log.Debug("Query_JavaContract::Finish query", "Success", !queryResult.exceptionOccurred, "info", queryResult.info)

	response := &jvmTypes.JVMQueryResponse{
		Success:!queryResult.exceptionOccurred,
		Result:queryResult.info,
	}
	return response, nil
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
