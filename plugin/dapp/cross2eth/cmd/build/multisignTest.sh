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


# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"


chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
#chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

#ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
#ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
#ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10

Chain33Cli="../../chain33-cli"
chain33BridgeBank=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
#chain33EthTokenAddr=""
#ethereumBtyTokenAddr=""
#chain33YccTokenAddr=""
ethereumYccTokenAddr=""
multisignChain33Addr=""
multisignEthAddr=""
ethBridgeToeknYccAddr=""
chain33YccErc20Addr=""

CLIA="./ebcli_A"

# shellcheck disable=SC2034
{
chain33MultisignA="168Sn1DXnLrZHTcAM9stD6t2P49fNuJfJ9"
chain33MultisignB="13KTf57aCkVVJYNJBXBBveiA5V811SrLcT"
chain33MultisignC="1JQwQWsShTHC4zxHzbUfYQK4kRBriUQdEe"
chain33MultisignD="1NHuKqoKe3hyv52PF8XBAyaTmJWAqA2Jbb"
chain33MultisignKeyA="0xcd284cd17456b73619fa609bb9e3105e8eff5d059c5e0b6eb1effbebd4d64144"
chain33MultisignKeyB="0xe892212221b3b58211b90194365f4662764b6d5474ef2961ef77c909e31eeed3"
chain33MultisignKeyC="0x9d19a2e9a440187010634f4f08ce36e2bc7b521581436a99f05568be94dc66ea"
chain33MultisignKeyD="0x45d4ce009e25e6d5e00d8d3a50565944b2e3604aa473680a656b242d9acbff35"

ethMultisignA=0x4c85848a7E2985B76f06a7Ed338FCB3aF94a7DCf
ethMultisignB=0x6F163E6daf0090D897AD7016484f10e0cE844994
ethMultisignC=0xbc333839E37bc7fAAD0137aBaE2275030555101f
ethMultisignD=0x495953A743ef169EC5D4aC7b5F786BF2Bd56aFd5
ethMultisignKeyA=0x5e8aadb91eaa0fce4df0bcc8bd1af9e703a1d6db78e7a4ebffd6cf045e053574
ethMultisignKeyB=0x0504bcb22b21874b85b15f1bfae19ad62fc2ad89caefc5344dc669c57efa60db
ethMultisignKeyC=0x0c61f5a879d70807686e43eccc1f52987a15230ae0472902834af4d1933674f2
ethMultisignKeyD=0x2809477ede1261da21270096776ba7dc68b89c9df5f029965eaa5fe7f0b80697
}

function deployMultisign() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    for name in A B C D; do
        eval chain33MultisignKey=\$chain33MultisignKey${name}
        eval chain33Multisign=\$chain33Multisign${name}
        # shellcheck disable=SC2154
        result=$(${Chain33Cli} account import_key -k "${chain33MultisignKey}" -l multisignAddr$name)
        # shellcheck disable=SC2154
        check_addr "${result}" "${chain33Multisign}"

        # chain33Multisign 要有手续费
        hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Multisign}" -k "${chain33DeployAddr}")
        check_tx "${Chain33Cli}" "${hash}"
        result=$(${Chain33Cli} account balance -a "${chain33Multisign}" -e coins)
        balance_ret "${result}" "10.0000"
    done

    echo -e "${GRE}=========== 部署 chain33 离线钱包合约 ===========${NOC}"
    result=$(${CLIA} chain33 multisign deploy)
    cli_ret "${result}" "chain33 multisign deploy"
    multisignChain33Addr=$(echo "${result}" | jq -r ".msg")

    result=$(${CLIA} chain33 multisign setup -k "${chain33DeployKey}" -o "${chain33MultisignA},${chain33MultisignB},${chain33MultisignC},${chain33MultisignD}")
    cli_ret "${result}" "chain33 multisign setup"

    # multisignChain33Addr 要有手续费
    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${multisignChain33Addr}" -k "${chain33DeployAddr}")
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configOfflineSaveAccount(${multisignChain33Addr})")
    check_tx "${Chain33Cli}" "${hash}"

    echo -e "${GRE}=========== 部署 ETH 离线钱包合约 ===========${NOC}"
    result=$(${CLIA} ethereum multisign deploy)
    cli_ret "${result}" "ethereum multisign deploy"
    multisignEthAddr=$(echo "${result}" | jq -r ".msg")

    result=$(${CLIA} ethereum multisign setup -k "${ethDeployKey}" -o "${ethMultisignA},${ethMultisignB},${ethMultisignC},${ethMultisignD}")
    cli_ret "${result}" "ethereum multisign setup"

    result=$(${CLIA} ethereum multisign set_offline_addr -s "${multisignEthAddr}")
    cli_ret "${result}" "set_offline_addr"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lock_bty_balance () {
    local lockAmount=$1
    local lockAmount2="${1}00000000"
    local bridgeBankBalance=$2
    local multisignBalance=$3

    hash=$(${Chain33Cli} evm call -f 1 -a "${lockAmount}" -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, ${lockAmount2})")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "${bridgeBankBalance}"

    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
    balance_ret "${result}" "${multisignBalance}"
}

function lockBty() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

#    echo '2:#配置自动转离线钱包(bty, 1000, 50%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(${chain33BtyTokenAddr},BTY,100000000000,50)")
    check_tx "${Chain33Cli}" "${hash}"

    lock_bty_balance 330 "330.0000" "0.0000"
    lock_bty_balance 800 "565.0000" "565.0000"
    lock_bty_balance 500 "532.5000" "1097.5000"

    # transfer test
    hash=$(./ebcli_A chain33 multisign transfer -a 100 -r "${chain33BridgeBank}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}")
#    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)

    hash=$(./ebcli_A chain33 multisign transfer -a 100 -r "${chain33MultisignA}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}")
#    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
    result=$(${Chain33Cli} account balance -a "${chain33MultisignA}" -e evm)

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lock_bty_ycc_balance () {
    local lockAmount="${1}00000000"
    local bridgeBankBalance="${2}00000000"
    local multisignBalance="${3}00000000"
    if [[ "${3}" == "0" ]]; then
        multisignBalance="0"
    fi

    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(${ethDeployAddr}, ${chain33YccErc20Addr}, ${lockAmount})")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "${bridgeBankBalance}"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${multisignChain33Addr}" -b "balanceOf(${multisignChain33Addr})")
    is_equal "${result}" "${multisignBalance}"
}

function lockChain33Ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
#    echo '2:#配置自动转离线钱包(YCC, 100, 60%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(${chain33YccErc20Addr},YCC,10000000000,60)")
    check_tx "${Chain33Cli}" "${hash}"

#    echo 'YCC.0:增加allowance的设置,或者使用relayer工具进行'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33YccErc20Addr}" -p "approve(${chain33BridgeBank}, 330000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    # echo 'YCC.2:#执行add lock操作:addToken2LockList'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "addToken2LockList(${chain33YccErc20Addr}, YCC)")
    check_tx "${Chain33Cli}" "${hash}"

    lock_bty_ycc_balance 30 30 0
    lock_bty_ycc_balance 70 40 60
    lock_bty_ycc_balance 260 120 240
    lock_bty_ycc_balance 10 52 318

    sleep ${maturityDegree}

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknYccAddr}" )
    cli_ret "${result}" "balance" ".balance" "370"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lock_eth_balance() {
    local lockAmount=$1
    local bridgeBankBalance=$2
    local multisignBalance=$3

    result=$(${CLIA} ethereum lock -m "${lockAmount}" -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "${bridgeBankBalance}"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" )
    cli_ret "${result}" "balance" ".balance" "${multisignBalance}"
}

function lockEth() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # echo '2:#配置自动转离线钱包(eth, 20, 50%)'
    result=$(${CLIA} ethereum multisign set_offline_token -s ETH -m 20)
    cli_ret "${result}" "set_offline_token -s ETH -m 20"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "0"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" )
    cli_ret "${result}" "balance" ".balance" "0"

    lock_eth_balance 19 19 0
    lock_eth_balance 1 10 10
    lock_eth_balance 16 13 23

    # transfer
    hash=$(./ebcli_A ethereum multisign transfer -a 10 -r "${ethBridgeBank}" -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}")

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" )

#0x5e8aadb91eaa0fce4df0bcc8bd1af9e703a1d6db78e7a4ebffd6cf045e053574,0x0504bcb22b21874b85b15f1bfae19ad62fc2ad89caefc5344dc669c57efa60db,0x0c61f5a879d70807686e43eccc1f52987a15230ae0472902834af4d1933674f2,0x2809477ede1261da21270096776ba7dc68b89c9df5f029965eaa5fe7f0b80697

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lock_eth_ycc_balance() {
    local lockAmount=$1
    local bridgeBankBalance=$2
    local multisignBalance=$3

    result=$(${CLIA} ethereum lock -m "${lockAmount}" -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 10 个区块
    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "${bridgeBankBalance}"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "${multisignBalance}"
}

function lockEthYcc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # echo '2:#配置自动转离线钱包(ycc, 100, 40%)'
    result=$(${CLIA} ethereum multisign set_offline_token -s YCC -m 100 -p 40 -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "set_offline_token -s YCC -m 100"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    lock_eth_ycc_balance 70 70 0
    lock_eth_ycc_balance 30 60 40
    lock_eth_ycc_balance 60 72 88

    # transfer
    # multisignEthAddr 要有手续费
    ./ebcli_A ethereum transfer -k "${ethDeployKey}" -m 10 -r "${multisignEthAddr}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function mainTest() {
    StartChain33
    start_trufflesuite
    AllRelayerStart

    deployMultisign

#    lockBty
#    lockChain33Ycc
    lockEth
    lockEthYcc
}

mainTest
