package executor

import (
	"os"
	"os/exec"
	"fmt"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor/contract"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
)

// Exec_CreateJvmContract 创建合约
func (jvm *JVMExecutor) Exec_CreateJvmContract(createJvmContract *jvmTypes.CreateJvmContract, tx *types.Transaction, index int) (*types.Receipt, error) {
	jvm.prepareExecContext(tx, index)
	// 使用随机生成的地址作为合约地址（这个可以保证每次创建的合约地址不会重复，不存在冲突的情况）

	contractAddr := address.GetExecAddress(createJvmContract.Name)
	contractAddrInStr := contractAddr.String()
	if !jvm.mStateDB.Empty(contractAddrInStr) {
		return nil, jvmTypes.ErrContractAddressCollisionJvm
	}

	log.Debug("jvm create", "new created jvm contract addr =", contractAddrInStr)

	codeSize := len(createJvmContract.GetCode())
	if codeSize > jvmTypes.MaxCodeSize {
		return nil, jvmTypes.ErrMaxCodeSizeExceededJvm
	}
	// 此处暂时不考虑消息发送签名的处理，chain33在mempool中对签名做了检查
	from := address.PubKeyToAddress(jvm.tx.GetSignature().GetPubkey())
	to := getReceiver(jvm.tx)
	if to == nil {
		return nil, types.ErrInvalidAddress
	}

	snapshot := jvm.mStateDB.Snapshot()

	// 创建新的合约对象，包含双方地址以及合约代码，可用Gas信息
	contract := contract.NewContract(contract.AccountRef(*from), contract.AccountRef(*contractAddr), 0)
	contract.SetCallCode(*contractAddr, common.BytesToHash(common.Sha256(createJvmContract.Code)), createJvmContract.Code)

	// 创建一个新的账户对象（合约账户）
	jvm.mStateDB.CreateAccount(contractAddrInStr, contract.CallerAddress.String(), createJvmContract.Name)
	jvm.mStateDB.SetCodeAndAbi(contractAddrInStr, createJvmContract.Code, []byte(createJvmContract.Abi))

	receipt, err := jvm.GenerateExecReceipt(
		snapshot,
		createJvmContract.Name,
		contract.CallerAddress.String(),
		contractAddrInStr,
		jvmTypes.CreateJvmContractAction)
	log.Debug("jvm create", "receipt", receipt, "err info", err)

	return receipt, err
}

// Exec_CallJvmContract 调用合约
func (jvm *JVMExecutor) Exec_CallJvmContract(callJvmContract *jvmTypes.CallJvmContract, tx *types.Transaction, index int) (*types.Receipt, error) {
	jvm.prepareExecContext(tx, index)
	//因为在真正地执行user.jvm.xxx合约前，还需要通过Jvm合约平台获取其合约字节码，
	//所以需要先将其合约名字设置为Jvm
	jvm.mStateDB.SetCurrentExecutorName(jvmTypes.JvmX)

	log.Debug("jvm call", "Para CallJvmContract", callJvmContract)

	userJvmAddr := address.ExecAddress(string(tx.Execer))
	contractAccount := jvm.mStateDB.GetAccount(userJvmAddr)
	code := contractAccount.Data.GetCode()
	if len(code) == 0 {
		log.Error("call jvm contract ", "failed to get code&abi from contract", string(tx.Execer))
		return nil, jvmTypes.ErrWrongContractAddr
	}

	//将当前合约执行名字修改为user.jvm.xxx
	jvm.mStateDB.SetCurrentExecutorName(string(jvm.GetAPI().GetConfig().GetParaExec(tx.Execer)))
	snapshot := jvm.mStateDB.Snapshot()
	setJvm4CallbackWithIndex(jvm, 0)

	//1st step: create apply context
	log.Debug("jvm call para", "ActionData", callJvmContract.ActionData,
		"ContractName", string(tx.Execer),
		"ActionName", callJvmContract.ActionName)
	caller := tx.From()
	actineInfo := string(callJvmContract.ActionData)
	actineInfo = fmt.Sprintf("%s:%s:from=%s", callJvmContract.ActionName, actineInfo, caller)
	//2nd step: just call contract
	errinfo := runJava(code, jvm.cp, string(tx.Execer), actineInfo)
	//合约执行失败
	if errinfo != nil {
		jvm.mStateDB.RevertToSnapshot(snapshot)
		log.Error("call jvm contract", "failed to call contract due to", errinfo)
		return nil, errinfo
	}

	receipt, err := jvm.GenerateExecReceipt(
		snapshot,
		contractAccount.GetExecName(),
		caller,
		userJvmAddr,
		jvmTypes.CallJvmContractAction)
	log.Debug("jvm call", "receipt", receipt, "err info", err)

	return receipt, err
}

// Exec_UpdateJvmContract 创建合约
func (jvm *JVMExecutor) Exec_UpdateJvmContract(in *jvmTypes.UpdateJvmContract, tx *types.Transaction, index int) (*types.Receipt, error) {
	jvm.prepareExecContext(tx, index)

	// 使用随机生成的地址作为合约地址（这个可以保证每次创建的合约地址不会重复，不存在冲突的情况）
	contractAddr := address.GetExecAddress(in.Name)
	contractAddrInStr := contractAddr.String()
	if !jvm.mStateDB.Exist(contractAddrInStr) {
		return nil, jvmTypes.ErrContractNotExists
	}
	//只有创建合约的人可以更新合约
	manager := jvm.mStateDB.GetAccount(contractAddrInStr).GetCreator()
	if tx.From() != manager {
		log.Error("update contract", "tx from:", tx.From(), "manager:", manager)
		return nil, jvmTypes.ErrNoPermission
	}
	log.Debug("jvm update", "updated jvm contract addr =", contractAddrInStr)

	codeSize := len(in.GetCode())
	if codeSize > jvmTypes.MaxCodeSize {
		return nil, jvmTypes.ErrMaxCodeSizeExceededJvm
	}

	// 此处暂时不考虑消息发送签名的处理，chain33在mempool中对签名做了检查
	from := address.PubKeyToAddress(jvm.tx.GetSignature().GetPubkey())
	to := getReceiver(jvm.tx)
	if to == nil {
		return nil, types.ErrInvalidAddress
	}

	snapshot := jvm.mStateDB.Snapshot()
	// 创建新的合约对象，包含双方地址以及合约代码，可用Gas信息
	contract := contract.NewContract(contract.AccountRef(*from), contract.AccountRef(*contractAddr), 0)
	contract.SetCallCode(*contractAddr, common.BytesToHash(common.Sha256(in.Code)), in.Code)
	jvm.mStateDB.SetCodeAndAbi(contractAddrInStr, in.Code, []byte(in.Abi))

	receipt, err := jvm.GenerateExecReceipt(
		snapshot,
		in.Name,
		contract.CallerAddress.String(),
		contractAddrInStr,
		jvmTypes.UpdateJvmContractAction)
	log.Debug("jvm create", "receipt", receipt, "err info", err)

	return receipt, err
}

func runJava(code []byte, cp, contractName string, actionData string) error {
	//java -cp .:../test/testclasses/src/main/java:/usr/lib/jdk/jdk1.8.0_251/jre/lib/jna.jar Chain33DBop
	//cmd := exec.Command("ls", "-a")
	//para := []string{"-cp .:../test/testclasses/src/main/java:/usr/lib/jdk/jdk1.8.0_251/jre/lib/jna.jar Chain33DBop"}
	//cmd := exec.Command("java", para...)
	//cmd := exec.Command("find", "-name", "runjava.go")
	//cmd := exec.Command("java", "-cp", ".:../test/testclasses/src/main/java:/usr/lib/jdk/jdk1.8.0_251/jre/lib/jna.jar", "Chain33DBop")
	javaClassfile, err := os.OpenFile("./"+contractName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer javaClassfile.Close()
	writeLen, err := javaClassfile.Write(code)
	if writeLen != len(code) {
		return jvmTypes.ErrWriteJavaClass
	}
	//准备执行
	//将当前临时保存javaclass的路径添加到cp中
	cp = cp + ":./"
	cmd := exec.Command("java", "-cp", cp, contractName, actionData)
	//获取输出对象，可以从该对象中读取输出结果
	out, err := cmd.Output()
	if err != nil {
		println("error:", err.Error())
		return jvmTypes.ErrJavaExecFailed
	}
	log.Debug("runJava", "java contract result ", string(out))
	return nil
}