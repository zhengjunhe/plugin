package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/33cn/chain33/system/crypto/secp256k1"
	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	evmAbi "github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

type TxCreateInfo struct {
	PrivateKey string
	Expire     string
	Note       string
	Fee        int64
	ParaName   string
	ChainID    int32
}

func CreateContractAndSign(txCreateInfo *TxCreateInfo, code, abi, parameter, contractName string) (string, []byte, error) {
	var action evmtypes.EVMContractAction
	bCode, err := common.FromHex(code)
	if err != nil {
		return "", nil, errors.New(contractName + " parse evm code error " + err.Error())
	}
	action = evmtypes.EVMContractAction{Amount: 0, Code: bCode, GasLimit: 0, GasPrice: 0, Note: txCreateInfo.Note, Alias: contractName}
	if parameter != "" {
		constructorPara := "constructor(" + parameter + ")"
		packData, err := evmAbi.PackContructorPara(constructorPara, abi)
		if err != nil {
			return "", nil, errors.New(contractName + " " + constructorPara + " Pack Contructor Para error:" + err.Error())
		}
		action.Code = append(action.Code, packData...)
	}
	data, txHash, err := CreateAndSignEvmTx(txCreateInfo.ChainID, &action, txCreateInfo.ParaName+"evm", txCreateInfo.PrivateKey, address.ExecAddress(txCreateInfo.ParaName+"evm"), txCreateInfo.Expire, txCreateInfo.Fee)
	if err != nil {
		return "", nil, errors.New(contractName + " create contract error:" + err.Error())
	}

	return data, txHash, nil
}

func CreateAndSignEvmTx(chainID int32, action proto.Message, execer, privateKeyStr, contract2call, expire string, fee int64) (string, []byte, error) {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: 0, To: contract2call}

	expireInt64, err := types.ParseExpire(expire)
	if nil != err {
		return "", nil, err
	}

	if expireInt64 > types.ExpireBound {
		if expireInt64 < int64(time.Second*120) {
			expireInt64 = int64(time.Second * 120)
		}
		//用秒数来表示的时间
		tx.Expire = types.Now().Unix() + expireInt64/int64(time.Second)
	} else {
		tx.Expire = expireInt64
	}

	tx.Fee = int64(1e7)
	if tx.Fee < fee {
		tx.Fee += fee
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	tx.ChainID = chainID

	var driver secp256k1.Driver
	privateKeySli, err := common.FromHex(privateKeyStr)
	if nil != err {
		return "", nil, err
	}
	privateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return "", nil, err
	}

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)

	return dataStr, tx.Hash(), nil
}

func WriteContractFile(fileName string, content string) {
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
	fmt.Println("tx is written to file: ", fileName)
}

func CallContractAndSign(txCreateInfo *TxCreateInfo, action *evmtypes.EVMContractAction, contractAddr string) (string, []byte, error) {
	data, txHash, err := CreateAndSignEvmTx(txCreateInfo.ChainID, action, txCreateInfo.ParaName+"evm", txCreateInfo.PrivateKey, contractAddr, txCreateInfo.Expire, txCreateInfo.Fee)
	if err != nil {
		return "", nil, errors.New(contractAddr + " call contract error:" + err.Error())
	}
	//fmt.Println("The call tx is as created below:")
	//fmt.Println(data)

	return data, txHash, nil
}

func ParseFileInJson(file string, result interface{}) error {
	_, err := os.Stat(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(b, result)

}

func WriteToFileInJson(fileName string, content interface{}) {
	jbytes, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, jbytes, 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
}
