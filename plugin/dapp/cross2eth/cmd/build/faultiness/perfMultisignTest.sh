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
chain33DeployAddr="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"
#chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

#maturityDegree=10

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
chain33ID=33

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

chain33BridgeBank=16A3uxgPqCv5pVkKqtdVnv2As6DbfRVZRH
multisignChain33Addr=1b193HbfvVUunUL2DVXrqt9jnbAWwLjcT

function lockBty() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
#    echo '2:#配置自动转离线钱包(bty, 100, 50%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(${chain33BtyTokenAddr},BTY,10000000000,50)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
#    balance_ret "${result}" "0"
    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
#    balance_ret "${result}" "0"

    for (( i = 0; i < 1000; i++ )); do
        echo "${i}"
        lock_bty_multisign 1
        sleep 1
    done

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
#    balance_ret "${result}" "50"
    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
#    balance_ret "${result}" "950"

#    # transfer test
#    hash=$(${CLIA} chain33 multisign transfer -a 100 -r "${chain33BridgeBank}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}" | jq -r ".msg")
#    check_tx "${Chain33Cli}" "${hash}"
#    sleep 2
#    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
#    balance_ret "${result}" "997.5000"
#    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
#    balance_ret "${result}" "632.5000"
#
#    hash=$(${CLIA} chain33 multisign transfer -a 100 -r "${chain33MultisignA}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}" | jq -r ".msg")
#    check_tx "${Chain33Cli}" "${hash}"
#    sleep 2
#    result=$(${Chain33Cli} account balance -a "${multisignChain33Addr}" -e evm)
#    balance_ret "${result}" "897.5000"
#    result=$(${Chain33Cli} account balance -a "${chain33MultisignA}" -e evm)
#    balance_ret "${result}" "100.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lockChain33Ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
#    echo '2:#配置自动转离线钱包(YCC, 100, 60%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(${chain33YccErc20Addr},YCC,10000000000,60)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    lock_chain33_ycc_multisign 30 30 0
    lock_chain33_ycc_multisign 70 40 60
    lock_chain33_ycc_multisign 260 120 240
    lock_chain33_ycc_multisign 10 52 318

     # transfer test
    hash=$(${CLIA} chain33 multisign transfer -a 10 -r "${chain33BridgeBank}" -t "${chain33YccErc20Addr}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}" | jq -r ".msg")
    check_tx "${Chain33Cli}" "${hash}"
    sleep 2
    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "6200000000"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${multisignChain33Addr}" -b "balanceOf(${multisignChain33Addr})")
    is_equal "${result}" "30800000000"

    hash=$(${CLIA} chain33 multisign transfer -a 5 -r "${chain33MultisignA}" -t "${chain33YccErc20Addr}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}" | jq -r ".msg")
    check_tx "${Chain33Cli}" "${hash}"
    sleep 2
    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33MultisignA}" -b "balanceOf(${chain33MultisignA})")
    is_equal "${result}" "500000000"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${multisignChain33Addr}" -b "balanceOf(${multisignChain33Addr})")
    is_equal "${result}" "30300000000"

    # 判断 ETH 这端是否金额一致
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknYccAddr}" )
    cli_ret "${result}" "balance" ".balance" "370"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
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

    lock_eth_multisign 19 19 0
    lock_eth_multisign 1 10 10
    lock_eth_multisign 16 13 23

    # transfer
    ${CLIA} ethereum multisign transfer -a 3 -r "${ethBridgeBank}" -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}"
    sleep 2
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "16"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}")
    cli_ret "${result}" "balance" ".balance" "20"

    # transfer
    ${CLIA} ethereum multisign transfer -a 5 -r "${ethMultisignA}" -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}"
    sleep 2
    result=$(${CLIA} ethereum balance -o "${ethMultisignA}")
    cli_ret "${result}" "balance" ".balance" "105"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}")
    cli_ret "${result}" "balance" ".balance" "15"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
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

    lock_ethereum_ycc_multisign 70 70 0
    lock_ethereum_ycc_multisign 30 60 40
    lock_ethereum_ycc_multisign 60 72 88

    # transfer
    # multisignEthAddr 要有手续费
    ./ebcli_A ethereum transfer -k "${ethDeployKey}" -m 10 -r "${multisignEthAddr}"

     # transfer
    ${CLIA} ethereum multisign transfer -a 8 -r "${ethBridgeBank}" -t "${ethereumYccTokenAddr}" -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}"
    sleep 2
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "80"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "80"

    # transfer
    ${CLIA} ethereum multisign transfer -a 10 -r "${ethMultisignA}" -t "${ethereumYccTokenAddr}" -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}"
    sleep 2
    result=$(${CLIA} ethereum balance -o "${ethMultisignA}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "10"
    result=$(${CLIA} ethereum balance -o "${multisignEthAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "70"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2120
function mainTest() {
    if [[ $# -ge 1 && "${1}" != "" ]]; then
        chain33ID="${1}"
    fi
    StartChain33
    start_trufflesuite
    AllRelayerStart

    deployMultisign

    lockBty
#    lockChain33Ycc
#    lockEth
#    lockEthYcc
}

mainTest "${1}"
#lockBty
