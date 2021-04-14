package chain33

import (
	"errors"
	"fmt"
	"time"

	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/contracts/contracts4chain33/generated"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func DeployAndInit2Chain33(rpcLaddr, paraChainName string, para4deploy *DeployPara4Chain33) (*X2EthDeployInfo, error) {
	deployer := para4deploy.Deployer.String()
	deployInfo := &X2EthDeployInfo{}
	var err error
	constructorPara := ""
	paraLen := len(para4deploy.InitValidators)

	var valsetAddr string
	var ethereumBridgeAddr string
	var oracleAddr string
	var bridgeBankAddr string

	//x2EthContracts.Valset, deployInfo.Valset, err = DeployValset(client, para.DeployPrivateKey, para.Deployer, para.Operator, para.InitValidators, para.InitPowers)
	//constructor(
	//	address _operator,
	//	address[] memory _initValidators,
	//	uint256[] memory _initPowers
	//)
	if 1 == paraLen {
		constructorPara = fmt.Sprintf("constructor(%s, %s, %d)", para4deploy.Operator.String(),
			para4deploy.InitValidators[0].String(),
			para4deploy.InitPowers[0].Int64())
	} else if 4 == paraLen {
		constructorPara = fmt.Sprintf("constructor(%s, %s, %s, %s, %s, %d, %d, %d, %d)", para4deploy.Operator.String(),
			para4deploy.InitValidators[0].String(), para4deploy.InitValidators[1].String(), para4deploy.InitValidators[2].String(), para4deploy.InitValidators[3].String(),
			para4deploy.InitPowers[0].Int64(), para4deploy.InitPowers[1].Int64(), para4deploy.InitPowers[2].Int64(), para4deploy.InitPowers[3].Int64())
	} else {
		panic(fmt.Sprintf("Not support valset with parameter count=%d", paraLen))
	}

	deployValsetHash, err := deploySingleContract(ethcommon.FromHex(generated.ValsetBin), generated.ValsetABI, constructorPara, "valset", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployValset due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("\nDeployValset tx hash:", deployInfo.Valset.TxHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("DeployValset timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployValsetHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy valset tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy valset failed due to" + ", ty = " + data)
				}
				valsetAddr = getContractAddr(deployer, deployValsetHash)
				fmt.Println("Succeed to deploy valset with address =", valsetAddr, "\n")
				goto deployEthereumBridge
			}
		}
	}

deployEthereumBridge:
	//x2EthContracts.Chain33Bridge, deployInfo.Chain33Bridge, err = DeployChain33Bridge(client, para.DeployPrivateKey, para.Deployer, para.Operator, deployInfo.Valset.Address)
	//constructor(
	//	address _operator,
	//	address _valset
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s)", para4deploy.Operator.String(), valsetAddr)
	deployEthereumBridgeHash, err := deploySingleContract(ethcommon.FromHex(generated.EthereumBridgeBin), generated.EthereumBridgeABI, constructorPara, "EthereumBridge", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to deployEthereumBridge due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("\nDeploy EthereumBridge Hash tx hash:", deployEthereumBridgeHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployEthereumBridge timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployEthereumBridgeHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy EthereumBridge tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy EthereumBridge failed due to" + ", ty = " + data)
				}
				ethereumBridgeAddr = getContractAddr(deployer, deployEthereumBridgeHash)
				fmt.Println("Succeed to deploy EthereumBridge with address =", ethereumBridgeAddr, "\n")
				goto deployOracle
			}
		}
	}

deployOracle:
	//constructor(
	//	address _operator,
	//	address _valset,
	//	address _ethereumBridge
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s)", para4deploy.Operator.String(), valsetAddr, ethereumBridgeAddr)
	//x2EthContracts.Oracle, deployInfo.Oracle, err = DeployOracle(client, para.DeployPrivateKey, para.Deployer, para.Operator, deployInfo.Valset.Address, deployInfo.Chain33Bridge.Address)
	deployOracleHash, err := deploySingleContract(ethcommon.FromHex(generated.OracleBin), generated.OracleABI, constructorPara, "Oracle", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployOracle due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("DeployOracle tx hash:", deployOracleHash)

		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployOracle timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployOracleHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy Oracle tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy Oracle failed due to" + ", ty = " + data)
				}
				oracleAddr = getContractAddr(deployer, deployOracleHash)
				fmt.Println("Succeed to deploy EthereumBridge with address =", oracleAddr, "\n")
				goto deployBridgeBank
			}
		}
	}
	/////////////////////////////////////
deployBridgeBank:
	//constructor (
	//	address _operatorAddress,
	//	address _oracleAddress,
	//	address _ethereumBridgeAddress
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s)", para4deploy.Operator.String(), oracleAddr, ethereumBridgeAddr)
	deployBridgeBankHash, err := deploySingleContract(ethcommon.FromHex(generated.BridgeBankBin), generated.BridgeBankABI, constructorPara, "BridgeBank", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to DeployBridgeBank due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("deployBridgeBank tx hash:", deployBridgeBankHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployBridgeBank timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for Deploy BridgeBank tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("Deploy BridgeBank failed due to" + ", ty = " + data)
				}
				bridgeBankAddr = getContractAddr(deployer, deployBridgeBankHash)
				fmt.Println("Succeed to deploy BridgeBank with address =", bridgeBankAddr, "\n")
				goto settingBridgeBank
			}
		}
	}

settingBridgeBank:
	////////////////////////
	//function setBridgeBank(
	//	address payable _bridgeBank
	//)
	callPara := fmt.Sprintf("setBridgeBank(%s)", bridgeBankAddr)
	settingBridgeBankHash, err := sendTx2Evm(callPara, rpcLaddr, ethereumBridgeAddr, paraChainName, deployer)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to settingBridgeBank due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("setBridgeBank tx hash:", settingBridgeBankHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("setBridgeBank timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for setBridgeBank tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("setBridgeBank failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to setBridgeBank ")
				goto setOracle
			}
		}
	}

setOracle:
	//function setOracle(
	//	address _oracle
	//)
	callPara = fmt.Sprintf("setOracle(%s)", oracleAddr)
	setOracleHash, err := sendTx2Evm(callPara, rpcLaddr, ethereumBridgeAddr, paraChainName, deployer)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to setOracle due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("setOracle tx hash:", setOracleHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("setOracle timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for setOracle tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("setOracle failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to setOracle ")
				goto deployBridgeRegistry
			}
		}
	}

deployBridgeRegistry:
	//constructor(
	//	address _ethereumBridge,
	//	address _bridgeBank,
	//	address _oracle,
	//	address _valset
	//)
	constructorPara = fmt.Sprintf("constructor(%s, %s, %s, %s)", ethereumBridgeAddr, bridgeBankAddr, oracleAddr, valsetAddr)
	deployBridgeRegistryHash, err := deploySingleContract(ethcommon.FromHex(generated.BridgeRegistryBin), generated.BridgeRegistryABI, constructorPara, "BridgeRegistry", paraChainName, para4deploy.Deployer.String(), rpcLaddr)
	if nil != err {
		chain33txLog.Error("DeployAndInit", "failed to deployBridgeRegistry due to:", err.Error())
		return nil, err
	}
	{
		fmt.Println("deployBridgeRegistryHash tx hash:", deployBridgeRegistryHash)
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("deployBridgeRegistry timeout")
			case <-oneSecondtimeout.C:
				data, _ := getTxByHashesRpc(deployBridgeBankHash, rpcLaddr)
				if data == "" {
					fmt.Println("No receipt received yet for deployBridgeRegistry tx and continue to wait")
					continue
				} else if data != "2" {
					return nil, errors.New("deployBridgeRegistry failed due to" + ", ty = " + data)
				}
				fmt.Println("Succeed to deployBridgeRegistry")
				goto finished
			}
		}
	}
finished:

	return deployInfo, nil
}
