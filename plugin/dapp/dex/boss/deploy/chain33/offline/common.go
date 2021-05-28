package offline

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

type TxCreateInfo struct {
	privateKey string
	expire     string
	note       string
	fee        int64
	paraName   string
	chainID    int32
}

func Chain33OfflineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain33",
		Short: "create and sign offline tx to deploy and set dex contracts to chain33",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		createERC20ContractCmd(),
		createFactoryCmd(),
		createWeth9Cmd(),
		createRouterCmd(),
		farmofflineCmd(),
	)
	return cmd
}

func writeContractFile(fileName string, content string) {
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
	fmt.Println("tx is written to file: ", fileName)
}

func createContractAndSign(txCreateInfo *TxCreateInfo, code, abi, parameter, contractName string) (string, error) {
	var action evmtypes.EVMContractAction
	bCode, err := common.FromHex(code)
	if err != nil {
		return "", errors.New(contractName + " parse evm code error " + err.Error())
	}
	action = evmtypes.EVMContractAction{Amount: 0, Code: bCode, GasLimit: 0, GasPrice: 0, Note: txCreateInfo.note, Alias: contractName}
	if parameter != "" {
		constructorPara := "constructor(" + parameter + ")"
		packData, err := evmAbi.PackContructorPara(constructorPara, abi)
		if err != nil {
			return "", errors.New(contractName + " " + constructorPara + " Pack Contructor Para error:" + err.Error())
		}
		action.Code = append(action.Code, packData...)
	}
	data, err := createAndSignEvmTx(txCreateInfo.chainID, &action, txCreateInfo.paraName+"evm", txCreateInfo.privateKey, address.ExecAddress(txCreateInfo.paraName+"evm"), txCreateInfo.expire, txCreateInfo.fee)
	if err != nil {
		return "", errors.New(contractName + " create contract error:" + err.Error())
	}
	fmt.Println("The created tx is as below:")
	fmt.Println(data)

	return data, nil
}

func callContractAndSign(txCreateInfo *TxCreateInfo, action *evmtypes.EVMContractAction, contractAddr string) (string, error) {
	data, err := createAndSignEvmTx(txCreateInfo.chainID, action, txCreateInfo.paraName+"evm", txCreateInfo.privateKey, contractAddr, txCreateInfo.expire, txCreateInfo.fee)
	if err != nil {
		return "", errors.New(contractAddr + " call contract error:" + err.Error())
	}
	fmt.Println("The call tx is as created below:")
	fmt.Println(data)

	return data, nil
}
