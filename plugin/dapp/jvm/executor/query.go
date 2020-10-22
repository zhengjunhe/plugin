package executor

import (
	"bytes"
	"os"
	"strings"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
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

//查询java合约状态
func (jvm *JVMExecutor) Query_JavaContract(in *jvmTypes.JVMQueryReq) (types.Message, error) {
	if in == nil {
		return nil, types.ErrInvalidParam
	}
	//兼容contract_name和user.jvm.contract_name
	if !bytes.Contains([]byte(in.Contract), []byte(jvmTypes.UserJvmX)) {
		in.Contract = jvmTypes.UserJvmX + in.Contract
	}

	jvm.prepareQueryContext([]byte(in.Contract))
	jvm.queryChan = make(chan QueryResult, 1)

	log.Debug("jvm call", "Para Query_JavaContract", in)

	contractName := in.Contract
	jvm.contract = in.Contract
	//将执行器名字设置为jvm，是为了能够获取java合约字节码
	jvm.mStateDB.SetCurrentExecutorName(jvmTypes.JvmX)
	userJvmAddr := address.ExecAddress(contractName)
	contractAccount := jvm.mStateDB.GetAccount(userJvmAddr)
	temp := strings.Split(contractName, ".")
	contractName = temp[len(temp) - 1]
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
	jvm.mStateDB.SetCurrentExecutorName(string(jvm.GetAPI().GetConfig().GetParaExec([]byte(in.Contract))))

	log.Debug("Query_JavaContract", "ContractName", contractName, "Para", in.Para)
	//2nd step: just call contract
	//在此处将gojvm指针传递到c实现的jvm中，进行回调的时候用来区分是获取数据时，使用执行db还是查询db
	_ = runJava(contractName, in.Para, jvm, TX_QUERY_JOB, jvm.GetAPI().GetConfig())

	//阻塞并等待查询结果的返回
	queryResult := <-jvm.queryChan
	log.Debug("Query_JavaContract::Finish query", "Success", !queryResult.exceptionOccurred, "info", queryResult.info)

	response := &jvmTypes.JVMQueryResponse{
		Success:!queryResult.exceptionOccurred,
		Result:queryResult.info,
	}
	return response, nil
}
