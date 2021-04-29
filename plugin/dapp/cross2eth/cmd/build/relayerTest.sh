#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

# ETH 部署合约者的私钥 用于部署合约时签名使用
ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrA="0x92c8b16afd6d423652559c6e266cbe1c29bfd84f"
ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"
#ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
#ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
#ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
# 新增地址 chain33 需要导入地址 转入 10 bty当收费费
#chain33Validator1="1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
#chain33Validator2="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
#chain33Validator3="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
#chain33Validator4="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
#chain33ValidatorKey1="0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
#chain33ValidatorKey2="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
#chain33ValidatorKey3="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
#chain33ValidatorKey4="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
#ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
#ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
#ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

#ethUrl=""

maturityDegree=5
Chain33Cli="../../chain33-cli"

CLIA="./ebcli_A"
BridgeRegistryOnChain33=""
chain33BridgeBank=""
BridgeRegistryOnEth=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""

function InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} relayer set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} relayer chain33 import_privatekey -k "${chain33DeployKey}")
    cli_ret "${result}" "chain33 import_privatekey"

    result=$(${CLIA} relayer ethereum import_privatekey -k "${ethDeployKey}")
    cli_ret "${result}" "ethereum import_privatekey"

    # 在 chain33 上部署合约
    result=$(${CLIA} relayer chain33 deploy)
    cli_ret "${result}" "chain33 deploy"
    BridgeRegistryOnChain33=$(echo "${result}" | jq -r ".msg")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnChain33}.abi"
    chain33BridgeBank=$(${Chain33Cli} evm abi call -c "${chain33DeployAddr}" -b "bridgeBank()" -a "${BridgeRegistryOnChain33}")
    cp BridgeBank.abi "${chain33BridgeBank}.abi"

    # 在 Eth 上部署合约
    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "ethereum deploy"
    BridgeRegistryOnEth=$(echo "${result}" | jq -r ".msg")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnEth}.abi"
    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    ethBridgeBank=$(echo "${result}" | jq -r ".addr")
    cp BridgeBank.abi "${ethBridgeBank}.abi"

    # 修改 relayer.toml 字段
    updata_relayer "BridgeRegistryOnChain33" "${BridgeRegistryOnChain33}" "./relayer.toml"

    line=$(delete_line_show "./relayer.toml" "BridgeRegistry=")
    sed -i ''"${line}"' a BridgeRegistry="'"${BridgeRegistryOnEth}"'"' "./relayer.toml"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function InitTokenAddr() {
    # set chain33 BTY Token
    result=$(${CLIA} relayer chain33 token set -s BTY -t "${chain33BtyTokenAddr}")
    cli_ret "${result}" "chain33 token set -s BTY"
    result=$(${CLIA} relayer chain33 token show | jq -r .tokenAddress[0].address)
    is_equal "${result}" "${chain33BtyTokenAddr}"

    # 在 Ethereum 上创建 bridgeToken BTY
    result=$(${CLIA} relayer ethereum token create-bridge-token -s BTY)
    cli_ret "${result}" "ethereum token create-bridge-token -s BTY"

    ethereumBtyTokenAddr=$(echo "${result}" | jq -r .addr)
    result=$(${CLIA} relayer ethereum token set -s BTY -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "ethereum token set -s BTY"
    result=$(${CLIA} relayer ethereum token show | jq -r .tokenAddress[0].address)
    is_equal "${result}" "${ethereumBtyTokenAddr}"

    # 在 chain33 上创建 bridgeToken ETH
    ${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "createNewBridgeToken(ETH)"
    chain33EthTokenAddr=$(${Chain33Cli} evm abi call -a "${chain33BridgeBank}" -c "${chain33DeployAddr}" -b "getToken2address(ETH)")
    echo "ETH Token Addr= ${chain33EthTokenAddr}"
    cp BridgeToken.abi "${chain33EthTokenAddr}.abi"

    result=$(${Chain33Cli} evm abi call -a ${chain33EthTokenAddr} -c ${chain33EthTokenAddr} -b "symbol()")
    is_equal "${result}" "ETH"

    # 设置 token 地址
    result=$(${CLIA} relayer chain33 token set -s ETH -t ${chain33EthTokenAddr})
    cli_ret "${result}" "chain33 token set -s ETH"

    result=$(${CLIA} relayer chain33 token show | jq -r .tokenAddress[1].address)
    is_equal "${result}" "${chain33EthTokenAddr}"

# YCC
}

function start_ebrelayerA() {
    nohup ./ebrelayer ./relayer.toml &
    sleep 2
}

function StartRelayerAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 修改 relayer.toml 配置文件 pushName 字段
    pushNameChange "./relayer.toml"

    # 启动 ebrelayer
    start_ebrelayerA

    # 导入私钥 部署合约 设置 bridgeRegistry 地址
    InitAndDeploy

    # 重启
    kill_ebrelayer ebrelayer
    start_ebrelayerA

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chian33 初始化准备
function InitChain33() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 转账到 EVM  合约中
    hash=$(${Chain33Cli} send coins send_exec -e evm -a 1000000 -k "${chain33DeployAddr}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
    balance_ret "${result}" "1000000.0000"

    # chain33Validator 要有手续费
#    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Validator1}" -k "${chain33DeployAddr}")
#    check_tx "${Chain33Cli}" "${hash}"
#    result=$(${Chain33Cli} account balance -a "${chain33Validator1}" -e coins)
#    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} relayer ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # 原来的地址金额
    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
    balance=$(cli_ret "${result}" "balance" ".balance")

    # chain33 lock bty
    hash=$(${Chain33Cli} evm call -f 1 -a 5 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)")
    check_tx "${Chain33Cli}" "${hash}"

    # 原来的地址金额 减少了 5
    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}-5" | bc)"
    #balance_ret "${result}" "195.0000"

    # chain33BridgeBank 是否增加了 5
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "5.0000"

    eth_block_wait 2

    # eth 这端 金额是否增加了 5
    result=$(${CLIA} relayer ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    # eth 这端 金额是否减少了 3
    result=$(${CLIA} relayer ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
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

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33Assets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    ./ebcli_A relayer ethereum balance -o "${ethDeployAddr}"
    ./ebcli_A relayer ethereum balance -o "${ethValidatorAddrA}"

    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} relayer ethereum balance -o "${ethBridgeBank}" )
    #cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 0.0001个
    result=$(${CLIA} relayer ethereum lock -m 0.0001 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    # 查询 ETH 这端 bridgeBank 地址 0.0001
    result=$(${CLIA} relayer ethereum balance -o "${ethBridgeBank}" )
    #cli_ret "${result}" "balance" ".balance" "0.0001"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 0.0001 * le8
    #is_equal "${result}" "10000"

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    #0x"${chain33ReceiverAddrKey}" 是 地址 "${chain33ReceiverAddr}" 的私钥
#                                                                                                           0x"${ethDeployAddr}"
    #${CLIA} relayer chain33 burn -m 2 -k 0x"${chain33ReceiverAddrKey}" -r "${ethDeployAddr}" -t "${chain33EthTokenAddr}"
    ${CLIA} relayer chain33 burn -m 0.0001 -k "${chain33ReceiverAddrKey}" -r "${ethDeployAddr}" -t "${chain33EthTokenAddr}"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
        #balance_ret "${result}" "10000"
    echo "check the balance on Ethereum"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} relayer ethereum balance -o "${ethBridgeBank}" )
#    cli_ret "${result}" "balance" ".balance" "0.0001"

# 比之前多 0.0001
./ebcli_A relayer ethereum balance -o "${ethDeployAddr}"


    ./ebcli_A relayer ethereum balance -o "${ethValidatorAddrA}"



    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Erc20() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartChain33() {
    kill_ebrelayer chain33
    sleep 2

    # delete chain33 datadir
    rm ../../datadir ../../logs -rf

    nohup ../../chain33 -f ./test.toml &

    ps -ef | grep chain33

    sleep 1

    # init
    ${Chain33Cli}  seed save -p 1314fuzamei -s "tortoise main civil member grace happy century convince father cage beach hip maid merry rib"
    ${Chain33Cli}  wallet unlock -p 1314fuzamei -t 0
    ${Chain33Cli}  account import_key -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944 -l returnAddr
    ${Chain33Cli}  account import_key -k "${chain33ReceiverAddrKey}" -l minerAddr
    ${Chain33Cli}  send coins transfer -a 10000 -n test -t "${chain33ReceiverAddr}" -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944
}

function mainTest() {
    StartChain33

    start_trufflesuite

    kill_ebrelayer ebrelayer
    sleep 10

    rm datadir/ logs/ -rf
    StartRelayerAndDeploy

    # 设置 token 地址
    InitTokenAddr
    InitChain33

    TestChain33ToEthAssets
    TestETH2Chain33Assets
}

mainTest