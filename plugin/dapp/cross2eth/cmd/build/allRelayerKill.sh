#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"

# ETH 部署合约者的私钥 用于部署合约时签名使用
ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
#ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
# shellcheck disable=SC2034
{
ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
}

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
#chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

# 新增地址 chain33 需要导入地址 转入 10 bty当收费费
# shellcheck disable=SC2034
chain33ValidatorA="1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
chain33ValidatorB="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
chain33ValidatorC="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
chain33ValidatorD="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
chain33ValidatorKeyA="0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
# shellcheck disable=SC2034
{
chain33ValidatorKeyB="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
chain33ValidatorKeyC="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
chain33ValidatorKeyD="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
}
ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
#ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
#ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10

Chain33Cli="../../chain33-cli"
BridgeRegistryOnChain33=""
chain33BridgeBank=""
BridgeRegistryOnEth=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""

CLIA="./ebcli_A"
# shellcheck disable=SC2034
CLIB="./ebcli_B"
CLIC="./ebcli_C"
CLID="./ebcli_D"

function kill_ebrelayerC() {
    kill_ebrelayer ebrelayer_C
    sleep 1
}
function kill_ebrelayerD() {
    kill_ebrelayer ebrelayer_D
    sleep 1
}

function start_ebrelayerC() {
    nohup ./relayer_B/ebrelayer ./relayer_B/relayer.toml &
    sleep 2
    ${CLIC} relayer unlock -p 123456hzj
    sleep ${maturityDegree}
}
function start_ebrelayerD() {
    nohup ./relayer_C/ebrelayer ./relayer_C/relayer.toml &
    sleep 2
    ${CLID} relayer unlock -p 123456hzj
    sleep ${maturityDegree}
}

function InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} chain33 import_privatekey -k "${chain33ValidatorKeyA}")
    cli_ret "${result}" "chain33 import_privatekey"

    result=$(${CLIA} ethereum import_privatekey -k "${ethDeployKey}")
    cli_ret "${result}" "ethereum import_privatekey"

    # 在 chain33 上部署合约
    result=$(${CLIA} chain33 deploy)
    cli_ret "${result}" "chain33 deploy"
    BridgeRegistryOnChain33=$(echo "${result}" | jq -r ".msg")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnChain33}.abi"
    chain33BridgeBank=$(${Chain33Cli} evm abi call -c "${chain33DeployAddr}" -b "bridgeBank()" -a "${BridgeRegistryOnChain33}")
    cp BridgeBank.abi "${chain33BridgeBank}.abi"

    # 在 Eth 上部署合约
    result=$(${CLIA} ethereum deploy)
    cli_ret "${result}" "ethereum deploy"
    BridgeRegistryOnEth=$(echo "${result}" | jq -r ".msg")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${BridgeRegistryOnEth}.abi"
    result=$(${CLIA} ethereum bridgeBankAddr)
    ethBridgeBank=$(echo "${result}" | jq -r ".addr")
    cp BridgeBank.abi "${ethBridgeBank}.abi"

    # 修改 relayer.toml 字段
    updata_relayer "BridgeRegistryOnChain33" "${BridgeRegistryOnChain33}" "./relayer.toml"

    line=$(delete_line_show "./relayer.toml" "BridgeRegistry=")
    sed -i ''"${line}"' a BridgeRegistry="'"${BridgeRegistryOnEth}"'"' "./relayer.toml"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function InitTokenAddr() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 在 Ethereum 上创建 bridgeToken BTY
    result=$(${CLIA} ethereum token create-bridge-token -s BTY)
    cli_ret "${result}" "ethereum token create-bridge-token -s BTY"
    ethereumBtyTokenAddr=$(echo "${result}" | jq -r .addr)

    # 在 chain33 上创建 bridgeToken ETH
    ${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "createNewBridgeToken(ETH)"
    sleep 1
    chain33EthTokenAddr=$(${Chain33Cli} evm abi call -a "${chain33BridgeBank}" -c "${chain33DeployAddr}" -b "getToken2address(ETH)")
    echo "ETH Token Addr= ${chain33EthTokenAddr}"
    cp BridgeToken.abi "${chain33EthTokenAddr}.abi"

    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33EthTokenAddr}" -b "symbol()")
    is_equal "${result}" "ETH"

    # 在chain33上创建bridgeToken YCC
    ${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "createNewBridgeToken(YCC)"
    sleep 1
    chain33YccTokenAddr=$(${Chain33Cli} evm abi call -a "${chain33BridgeBank}" -c "${chain33DeployAddr}" -b "getToken2address(YCC)")
    echo "YCC Token Addr = ${chain33YccTokenAddr}"
    cp BridgeToken.abi "${chain33YccTokenAddr}.abi"

    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33YccTokenAddr}" -b "symbol()")
    is_equal "${result}" "YCC"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function start_ebrelayerA() {
    nohup ./ebrelayer ./relayer.toml &
    sleep 2
}

function StartRelayerAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 修改 relayer.toml 配置文件 pushName 字段
    pushNameChange "./relayer.toml"

    # 修改 relayer.toml 配置文件 initPowers
    {
        # shellcheck disable=SC2155
        line=$(delete_line_show "./relayer.toml" "initPowers=\[96, 1, 1, 1\]")
        sed -i ''"${line}"' a initPowers=[25, 25, 25, 25]' "./relayer.toml"
        line=$(delete_line_show "./relayer.toml" "initPowers=\[96, 1, 1, 1\]")
        sed -i ''"${line}"' a initPowers=[25, 25, 25, 25]' "./relayer.toml"

        line=$(delete_line_show "./relayer.toml" "validatorsAddr=\[\"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
        sed -i ''"${line}"' a validatorsAddr=['\""${chain33ValidatorA}"\"', '\""${chain33ValidatorB}"\"', '\""${chain33ValidatorC}"\"', '\""${chain33ValidatorD}"\"']' "./relayer.toml"
    }

    # 启动 ebrelayer
    start_ebrelayerA

    # 导入私钥 部署合约 设置 bridgeRegistry 地址
    InitAndDeploy

    # 重启
    kill_ebrelayer ebrelayer
    start_ebrelayerA

    result=$(${CLIA} unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    # start ebrelayer B C D
    bind_port=9901
    push_port=20000
    for name in B C D; do
        local file="./relayer_$name/relayer.toml"
        cp './relayer.toml' "${file}"

        # 删除配置文件中不需要的字段
        for deleteName in "deploy4chain33" "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers" "deploy" "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers"; do
            delete_line "${file}" "${deleteName}"
        done

        bind_port=$((bind_port + 1))
        line=$(delete_line_show "./relayer_$name/relayer.toml" "JrpcBindAddr")
        sed -i ''"${line}"' a JrpcBindAddr="localhost:'${bind_port}'"' "./relayer_$name/relayer.toml"

        push_port=$((push_port + 1))
        line=$(delete_line_show "./relayer_$name/relayer.toml" "pushHost")
        sed -i ''"${line}"' a pushHost="http://localhost:'${push_port}'"' "./relayer_$name/relayer.toml"
        line=$(delete_line_show "./relayer_$name/relayer.toml" "pushBind")
        sed -i ''"${line}"' a pushBind="0.0.0.0:'${push_port}'"' "./relayer_$name/relayer.toml"

        sleep 1
        pushNameChange "./relayer_$name/relayer.toml"

        nohup ./relayer_$name/ebrelayer ./relayer_$name/relayer.toml &
        sleep 2

        CLI="./ebcli_$name"
        result=$(${CLI} set_pwd -p 123456hzj)
        cli_ret "${result}" "set_pwd"

        result=$(${CLI} unlock -p 123456hzj)
        cli_ret "${result}" "unlock"

        eval chain33ValidatorKey=\$chain33ValidatorKey${name}
        # shellcheck disable=SC2154
        result=$(${CLI} chain33 import_privatekey -k "${chain33ValidatorKey}")
        cli_ret "${result}" "chain33 import_privatekey"

        eval ethValidatorAddrKey=\$ethValidatorAddrKey${name}
        # shellcheck disable=SC2154
        result=$(${CLI} ethereum import_privatekey -k "${ethValidatorAddrKey}")
        cli_ret "${result}" "ethereum import_privatekey"
    done

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

    # 导入 chain33Validators 私钥生成地址
    for name in A B C D; do
        eval chain33ValidatorKey=\$chain33ValidatorKey${name}
        eval chain33Validator=\$chain33Validator${name}
        result=$(${Chain33Cli} account import_key -k "${chain33ValidatorKey}" -l validator$name)
        # shellcheck disable=SC2154
        check_addr "${result}" "${chain33Validator}"

        # chain33Validator 要有手续费
        hash=$(${Chain33Cli} send coins transfer -a 1000 -t "${chain33Validator}" -k "${chain33DeployAddr}")
        check_tx "${Chain33Cli}" "${hash}"
        result=$(${Chain33Cli} account balance -a "${chain33Validator}" -e coins)
        balance_ret "${result}" "1000.0000"
    done

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain33 lock bty
    hash=$(${Chain33Cli} evm call -f 1 -a 5 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)")
    check_tx "${Chain33Cli}" "${hash}"

    # chain33BridgeBank 是否增加了 5
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "0.0000"

    start_ebrelayerC

    # chain33BridgeBank 是否增加了 5
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "5.0000"

    eth_block_wait 2

    # eth 这端 金额是否增加了 5
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    kill_ebrelayerC

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "0.0000"

    start_ebrelayerD

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

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33AssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerD

    # ETH 这端 lock 11个
    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "0"

    sleep ${maturityDegree}

    start_ebrelayerC

    # 查询 ETH 这端 bridgeBank 地址 11
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "11"

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11 * le8
    is_equal "${result}" "1100000000"

    # 原来的数额
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    kill_ebrelayerC

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33EthTokenAddr}"

    sleep ${maturityDegree}

    # 原来的数额
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    start_ebrelayerC

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

function TestETH2Chain33YccKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # eth 铸币
    result=$(${CLIA} ethereum deploy_erc20 -c "${ethDeployAddr}" -n YCC -s YCC -m 33000000000000000000)
    cli_ret "${result}" "ethereum deploy_erc20 -s YCC"
    ethereumYccTokenAddr=$(echo "${result}" | jq -r .msg)

    result=$(${CLIA} ethereum token add_lock_list -s YCC -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "add_lock_list"

    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC

    # ETH 这端 lock 7个 YCC
    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2
    sleep ${maturityDegree}

    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    # 查询 ETH 这端 bridgeBank 地址 7 YCC
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7 * le8
    is_equal "${result}" "700000000"

    # 原来的数额 0
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33YccTokenAddr}"

    sleep ${maturityDegree}

    # 原来的数额 0
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

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

function StartChain33() {
    kill_ebrelayer chain33
    sleep 2

    # delete chain33 datadir
    rm ../../datadir ../../logs -rf

    nohup ../../chain33 -f ./ci/cross2eth/test.toml &

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
    InitChain33

    start_trufflesuite

    kill_all_ebrelayer
    StartRelayerAndDeploy

    # 设置 token 地址
    InitTokenAddr

    TestChain33ToEthAssetsKill
    TestETH2Chain33AssetsKill
    TestETH2Chain33YccKill
}

mainTest
