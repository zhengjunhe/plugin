#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"
source "./relayerPublic.sh"

# ETH 部署合约者的私钥 用于部署合约时签名使用
ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"

maturityDegree=10

Chain33Cli="../../chain33-cli"
chain33BridgeBank=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""
chain33ZBCErc20Addr=""
ethBridgeToeknZBCAddr=""

CLIA="./ebcli_A"
Boss4xCLI="./boss4x"
chain33ID=33

# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== chain33 lock BTY, eth burn BTY ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # 原来的地址金额
    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
#    balance=$(cli_ret "${result}" "balance" ".balance")

    # chain33 lock bty
    hash=$(${Chain33Cli} evm call -f 1 -a 5 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    # 原来的地址金额 减少了 5
    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
#    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}-5" | bc)"
    #balance_ret "${result}" "195.0000"

    # chain33BridgeBank 是否增加了 5
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "5.0000"

    eth_block_wait 2

    # eth 这端 金额是否增加了 5
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    # eth 这端 金额是否减少了 3
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    sleep ${maturityDegree}

     # 接收的地址金额 变成了 3
    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "3.0000"

    # chain33BridgeBank 是否减少了 3
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "2.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 lock ZBC, eth burn ZBC
function TestChain33ToEthZBCAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== chain33 lock ZBC, eth burn ZBC ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # 原来的地址金额
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "0"

    # chain33 lock ZBC
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33ZBCErc20Addr}, 900000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    # chain33BridgeBank 是否增加了 9
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "900000000"

    eth_block_wait 2

    # eth 这端 金额是否增加了 9
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "9"

    # eth burn
    result=$(${CLIA} ethereum burn -m 8 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethBridgeToeknZBCAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    # eth 这端 金额是否减少了 1
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "1"

    sleep ${maturityDegree}

     # 接收的地址金额 变成了 8
     result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33ReceiverAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "800000000"

    # chain33BridgeBank 是否减少了 1
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "100000000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33Assets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 11个
    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    # 查询 ETH 这端 bridgeBank 地址 11
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "11"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11 * le8
    is_equal "${result}" "1100000000"

    # 原来的数额
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33EthTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11-5 * le8
    is_equal "${result}" "600000000"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "6"

    # 比之前多 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "105"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ycc 资产,然后在 chain33 上 burn ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 7个 YCC
    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    # 查询 ETH 这端 bridgeBank 地址 7 YCC
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7 * le8
    is_equal "${result}" "700000000"

    # 原来的数额 0
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33YccTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7-5 * le8
    is_equal "${result}" "200000000"

    # 查询 ETH 这端 bridgeBank 地址 2
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    # 更新后的金额 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2120
function InitAndOfflineDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    # shellcheck disable=SC2154
    result=$(${CLIA} chain33 import_privatekey -k "${chain33DeployKey}")
    cli_ret "${result}" "chain33 import_privatekey"

    result=$(${CLIA} ethereum import_privatekey -k "${ethDeployKey}")
    cli_ret "${result}" "ethereum import_privatekey"

    # 在 chain33 上部署合约
#    result=$(${CLIA} chain33 deploy)
#    cli_ret "${result}" "chain33 deploy"
#    BridgeRegistryOnChain33=$(echo "${result}" | jq -r ".msg")

#    ./boss4x chain33 offline create -f 1 -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae -n "deploy crossx to chain33" -r "1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ, [1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ, 155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6, 13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv, 113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG], [25, 25, 25, 25]" --chainID 33
    # shellcheck disable=SC2154
    ${Boss4xCLI} chain33 offline create -f 1 -k "${chain33DeployKey}" -n "deploy crossx to chain33" -r "${chain33ValidatorA}, [${chain33ValidatorA}, ${chain33ValidatorB}, ${chain33ValidatorC}, ${chain33ValidatorD}], [25, 25, 25, 25]" --chainID 33
    result=$(${Boss4xCLI} chain33 offline send)
    hash=$(echo "${result}" | jq -r ".[6].TxHash")
    check_tx "${Chain33Cli}" "${hash}"
    hash=$(echo "${result}" | jq -r ".[7].TxHash")
    check_tx "${Chain33Cli}" "${hash}"
    BridgeRegistryOnChain33=$(echo "${result}" | jq -r ".[6].ContractAddr")
    # shellcheck disable=SC2034
    multisignChain33Addr=$(echo "${result}" | jq -r ".[7].ContractAddr")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnChain33}.abi"
    chain33BridgeBank=$(${Chain33Cli} evm abi call -c "${chain33DeployAddr}" -b "bridgeBank()" -a "${BridgeRegistryOnChain33}")
    cp Chain33BridgeBank.abi "${chain33BridgeBank}.abi"

#    exit 0

    # 在 Eth 上部署合约
#    result=$(${CLIA} ethereum deploy)
#    cli_ret "${result}" "ethereum deploy"
#    BridgeRegistryOnEth=$(echo "${result}" | jq -r ".msg")

#    ./boss4x ethereum offline create_sign -l 500000 -n 0 -g 20000000000 -p "25,25,25,25" -v "0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a,0x0df9a824699bc5878232c9e612fe1a5346a5a368,0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1,0xd9dab021e74ecf475788ed7b61356056b2095830" -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230
    cp ../../../plugin/dapp/cross2eth/boss4x/ethereum/deploy.toml ./deploy.toml
    ./boss4x ethereum offline sign -c ./deploy.toml -g 20000000000 -l 500000
    result=$(./boss4x ethereum offline send -f signed_cross2eth.txt)
    BridgeRegistryOnEth=$(echo "${result}" | jq -r ".[6].ContractAddr")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnEth}.abi"
#    result=$(${CLIA} ethereum bridgeBankAddr)
#    ethBridgeBank=$(echo "${result}" | jq -r ".addr")
#    ethBridgeBank=$(echo "${result}" | jq -r ".[3].ContractAddr")
#    cp EthBridgeBank.abi "${ethBridgeBank}.abi"

    # 修改 relayer.toml 字段
    updata_relayer "BridgeRegistryOnChain33" "${BridgeRegistryOnChain33}" "./relayer.toml"

    line=$(delete_line_show "./relayer.toml" "BridgeRegistry=")
    if [ "${line}" ]; then
        sed -i ''"${line}"' a BridgeRegistry="'"${BridgeRegistryOnEth}"'"' "./relayer.toml"
    fi

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartRelayerAndOfflineDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 修改 relayer.toml 配置文件 pushName 字段
    pushNameChange "./relayer.toml"
    validators_config

    # 删除私钥
#    delete_line "./relayer.toml" "deployerPrivateKey="
#    delete_line "./relayer.toml" "deployerPrivateKey="

    # 启动 ebrelayer
    start_ebrelayerA

    # 导入私钥 部署合约 设置 bridgeRegistry 地址
    InitAndOfflineDeploy

    # 重启
    kill_ebrelayer ebrelayer
    start_ebrelayerA

    result=$(${CLIA} unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    ethBridgeBank=$(./ebcli_A ethereum bridgeBankAddr | jq -r ".addr")
    cp EthBridgeBank.abi "${ethBridgeBank}.abi"

    # start ebrelayer B C D
    updata_toml_start_BCD

    # 设置 token 地址
    InitTokenAddr

    # 设置离线多签数据
    initMultisignAddr
    setupChain33AndEthMultisign
    transferChain33MultisignFee

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2120
function mainTest() {
    if [[ $# -ge 1 && "${1}" != "" ]]; then
        chain33ID="${1}"
    fi

    StartChain33
    start_trufflesuite
#    AllRelayerStart
    kill_all_ebrelayer
    StartRelayerAndOfflineDeploy
#    StartRelayerAndDeploy

    TestChain33ToEthAssets
    TestChain33ToEthZBCAssets
    TestETH2Chain33Assets
    TestETH2Chain33Ycc
}

mainTest "${1}"

