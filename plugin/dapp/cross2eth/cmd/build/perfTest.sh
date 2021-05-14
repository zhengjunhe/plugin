#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
# shellcheck disable=SC2178
set -x

source "./publicTest.sh"
source "./allRelayerPublic.sh"

ethAddress[0]=0x92C8b16aFD6d423652559C6E266cBE1c29Bfd84f
ethAddress[1]=0x0df9a824699Bc5878232C9e612fE1A5346a5A368
ethAddress[2]=0xcB074CB21cdDDF3ce9c3C0a7AC4497d633C9D9f1
ethAddress[3]=0xd9dAb021e74EcF475788ed7b61356056B2095830
ethAddress[4]=0xdb15E7327aDc83F2878624bBD6307f5Af1B477b4
ethAddress[5]=0x9cBA1fF8D0b0c9Bc95d5762533F8CddBE795f687
ethAddress[6]=0x1919203bA8b325278d28Fb8fFeac49F2CD881A4e
ethAddress[7]=0xA4Ea64a583F6e51C3799335b28a8F0529570A635
ethAddress[8]=0x0C05bA5c230fDaA503b53702aF1962e08D0C60BF

privateKeys[0]=3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e
privateKeys[1]=a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400
privateKeys[2]=bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e
privateKeys[3]=c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b
privateKeys[4]=1385016736f7379884763f4a39811d1391fa156a7ca017be6afffa52bb327695
privateKeys[5]=4ae589fe3837dcfc90d1c85b8423dc30841525cbebc41dfb537868b0f8376bbf
privateKeys[6]=62ca4122aac0e6f35bed02fc15c7ddbdaa07f2f2a1821c8b8210b891051e3ee9
privateKeys[7]=355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71
privateKeys[8]=9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08

# ETH 部署合约者的私钥 用于部署合约时签名使用
#ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
#chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

maturityDegree=10

Chain33Cli="../../chain33-cli"
chain33BridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""

CLIA="./ebcli_A"

function loop_send_lock_bty() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm | jq -r ".balance" | sed 's/\"//g')

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[$i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumBtyTokenAddr}" | jq -r ".balance")

        hash=$(${Chain33Cli} evm call -f 1 -a 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethAddress[i]}, ${chain33BtyTokenAddr}, 100000000)")
        check_tx "${Chain33Cli}" "${hash}"

        i=$((i+1))
    done

    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumBtyTokenAddr}" | jq -r ".balance")
        res=$((nowEthBalance - preEthBalance[i]))
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" ${res}
        check_number "${res}" 1
        i=$((i+1))
    done
    nowChain33Balance=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm | jq -r ".balance" | sed 's/\"//g')
    diff=$(echo "$preChain33Balance - $nowChain33Balance" | bc)
    check_number "${diff}" "${#privateKeys[@]}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function loop_send_burn_bty() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm | jq -r ".balance" | sed 's/\"//g')

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[$i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumBtyTokenAddr}" | jq -r ".balance")
        result=$(${CLIA} ethereum burn -m 1 -k "${privateKeys[i]}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" )
        cli_ret "${result}" "burn"
        i=$((i+1))
    done

    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumBtyTokenAddr}" | jq -r ".balance")
        res=$((preEthBalance[i] - nowEthBalance))
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" ${res}
        check_number "${res}" 1
        i=$((i+1))
    done
    nowChain33Balance=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm | jq -r ".balance" | sed 's/\"//g')
    diff=$(echo "$nowChain33Balance - $preChain33Balance" | bc)
    check_number "${diff}" "${#privateKeys[@]}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function loop_send_lock_eth() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[$i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" | jq -r ".balance")
        result=$(${CLIA} ethereum lock -m 1 -k "${privateKeys[i]}" -r "${chain33ReceiverAddr}")
        cli_ret "${result}" "lock"
        i=$((i+1))
    done

    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" | jq -r ".balance")
        res=$(echo "${preEthBalance[i]} - $nowEthBalance" | bc)
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" ${res}
        diff=$(echo "$res >= 1"| bc) # 浮点数比较 判断是否大于1 大于返回1 小于返回0
        if [ "${diff}" -ne 1 ]; then
            echo -e "${RED}error number, expect greater than 1, get ${res}${NOC}"
            exit 1
        fi
        i=$((i+1))
    done
    nowChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    diff=$(echo "$nowChain33Balance - $preChain33Balance" | bc)
    diff=$(echo "$diff / 100000000" | bc)
    check_number "${diff}" "${#privateKeys[@]}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function loop_send_burn_eth() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[$i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" | jq -r ".balance")
        ethTxHash=$(${CLIA} chain33 burn -m 1 -k "${chain33ReceiverAddrKey}" -r "${ethAddress[i]}" -t "${chain33EthTokenAddr}" | jq -r ".msg")
        echo ${i} "burn chain33 tx hash:" "${ethTxHash}"
        i=$((i+1))
    done

    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" | jq -r ".balance")
        res=$(echo "$nowEthBalance - ${preEthBalance[i]}" | bc)
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" ${res}
        diff=$(echo "$res >= 1"| bc) # 浮点数比较 判断是否大于1 大于返回1 小于返回0
        if [ "${diff}" -ne 1 ]; then
            echo -e "${RED}error number, expect greater than 1, get ${res}${NOC}"
            exit 1
        fi
        i=$((i+1))
    done
    nowChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    diff=$(echo "$preChain33Balance - $nowChain33Balance" | bc)
    diff=$(echo "$diff / 100000000" | bc)
    check_number "${diff}" "${#privateKeys[@]}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function loop_send_lock_ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    # 先往每个ETH地址中导入token
    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        ethTxHash=$(${CLIA} ethereum transfer -m 10 -k "${ethDeployKey}" -r "${ethAddress[i]}" -t "${ethereumYccTokenAddr}" | jq -r ".msg")
        echo ${i} "burn chain33 tx hash:" "${ethTxHash}"
        i=$((i+1))
    done

    sleep 2

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumYccTokenAddr}" | jq -r ".balance")
        ethTxHash=$(${CLIA} ethereum lock -m 1 -k "${privateKeys[i]}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}" | jq -r ".msg")
        echo ${i} "lock ycc tx hash:" "${ethTxHash}"
        i=$((i+1))
    done
    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumYccTokenAddr}" | jq -r ".balance")
        res=$(echo "${preEthBalance[i]} - $nowEthBalance" | bc)
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" "${res}"
        check_number "${res}" 1
        i=$((i+1))
    done

    nowChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    diff=$((nowChain33Balance - preChain33Balance))
    diff=$(echo "$diff / 100000000" | bc)
    check_number "${diff}" "${#privateKeys[@]}"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function loop_send_burn_ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    preChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        preEthBalance[i]=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumYccTokenAddr}" | jq -r ".balance")
        ethTxHash=$(${CLIA} chain33 burn -m 1 -k "${chain33ReceiverAddrKey}" -r "${ethAddress[i]}" -t "${chain33YccTokenAddr}" | jq -r ".msg")
        echo ${i} "burn chain33 tx hash:" "${ethTxHash}"
        i=$((i+1))
    done

    eth_block_wait $((maturityDegree + 2))

    i=0
    while [[ i -lt ${#privateKeys[@]} ]]; do
        nowEthBalance=$(${CLIA} ethereum balance -o "${ethAddress[i]}" -t "${ethereumYccTokenAddr}" | jq -r ".balance")
        res=$((nowEthBalance - preEthBalance[i]))
        echo ${i} "preBalance" "${preEthBalance[i]}" "nowBalance" "${nowEthBalance}" "diff" ${res}
        check_number "${res}" 1
        i=$((i+1))
    done
    nowChain33Balance=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    diff=$(echo "$preChain33Balance - $nowChain33Balance" | bc)
    diff=$(echo "$diff / 100000000" | bc)
    check_number "${diff}" "${#privateKeys[@]}"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function perf_test_main() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    StartChain33

    start_trufflesuite

    kill_all_ebrelayer
    StartRelayerAndDeploy

    loop_send_lock_bty
    loop_send_burn_bty
    loop_send_lock_eth
    loop_send_burn_eth
    loop_send_lock_ycc
    loop_send_burn_ycc

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

perf_test_main 10
