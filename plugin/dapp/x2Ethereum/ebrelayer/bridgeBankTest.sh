#!/usr/bin/env bash
set -x
CLI="../build/ebcli_A"
tokenAddr=""
bridgeBankAddr=""
chain33SenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33SenderAddrKey="CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
ethOperatorAddrKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
ethReceiverAddr1="0xdb15E7327aDc83F2878624bBD6307f5Af1B477b4"
ethReceiverAddrKey1="1385016736f7379884763f4a39811d1391fa156a7ca017be6afffa52bb327695"
ethReceiverAddr2="0x9cBA1fF8D0b0c9Bc95d5762533F8CddBE795f687"
ethReceiverAddrKey2="4ae589fe3837dcfc90d1c85b8423dc30841525cbebc41dfb537868b0f8376bbf"

InitAndDeploy() {
    result=$(${CLI} relayer set_pwd -n 123456hzj -o kk)
    cli_ret "${result}" "set_pwd"

    result=$(${CLI} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLI} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    result=$(${CLI} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
    cli_ret "${result}" "import_chain33privatekey"

    result=$(${CLI} relayer ethereum import_ethprivatekey -k "${ethOperatorAddrKey}")
    cli_ret "${result}" "import_ethprivatekey"

    echo "Succeed to InitAndDeploy"
}

# eth to chain33
# 在以太坊上锁定资产,然后在 chain33 上铸币,针对 erc20 资产
# 以太坊 brun 资产,balance 对比是否正确
TestETH2Chain33WithErc20Assets() {
    # token4erc20 在 chain33 上先有 token,同时 mint
    tokenSymbol="testc"
    result=$(${CLI} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")

    # 先铸币 1000
    result=$(${CLI} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLI} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 2 chain33 lock 前先审批一下
    result=$(${CLI} relayer ethereum approve -m 100 -k "${ethReceiverAddrKey1}" -t "${tokenAddr}")
    cli_ret "${result}" "approve"

    # ETH 2 chain33 lock
    # -r chain33 receiver addr
    result=$(${CLI} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # ETH 2 chain33 withdraw 40
    # -c 1 burn
    result=$(${CLI} relayer ethereum prophecy -m 40 -a "${chain33SenderAddr}" -c 1 -r "${ethReceiverAddr2}" -s "${tokenSymbol}" -t "${tokenAddr}")
    cli_ret "${result}" "prophecy -m 40"

    walitProphecyFinish "${ethReceiverAddr2}" "${tokenAddr}" 40

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "60"

    echo "Succeed to TestETH2Chain33WithErc20Assets"
}

# 在 chain33 上锁定资产,然后在 ETH 上铸币
TestETH2Chain33WithEthAssets() {
    result=$(${CLI} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # chain33 lock eth
    result=$(${CLI} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance3=$(cli_ret "${result}" "balance" ".balance")

    # ETH 2 chain33 withdraw 50
    result=$(${CLI} relayer ethereum prophecy -m 50 -a "${chain33SenderAddr}" -c 1 -r "${ethReceiverAddr2}" -s eth)
    cli_ret "${result}" "prophecy -m 50"

    walitProphecyFinish "${ethReceiverAddr2}" $((${balance3}+50))

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance4=$(cli_ret "${result}" "balance" ".balance")

    echo "${balance3} ${balance4}"
    if [[ "${balance4}" != $((${balance3}+50)) ]]; then
        echo "wrong balance"
        exit 1
    fi

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "50"

    echo "Succeed to TestETH2Chain33WithEthAssets"
}

# chain33 to eth
# 在 chain33 上锁定资产,然后在以太坊上铸币
# chain33 brun 资产,balance 对比是否正确
TestChain33ToEthAssets() {
    # token4chain33 在 以太坊 上先有 bty
    result=$(${CLI} relayer ethereum token4chain33 -s bty)
    tokenAddr=$(cli_ret "${result}" "token4chain33" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # -c 2 chain33 lock 在以太坊上铸币
    result=$(${CLI} relayer ethereum prophecy -m 100 -a "${chain33SenderAddr}" -c 2 -r "${ethReceiverAddr1}" -s bty -t "${tokenAddr}")
    cli_ret "${result}" "prophecy -m 100"

    walitProphecyFinish "${ethReceiverAddr1}" "${tokenAddr}" "100"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # brun
    result=$(${CLI} relayer ethereum burn -m 10 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "burn"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "90"

    echo "Succeed to TestChain33ToEthAssets"
}

checkProphecyIDActive() {
    while true; do
        pending=$(${CLI} relayer ethereum ispending -i ${1} | jq .isOK)
        if [[ ${pending} == "true" ]]; then
            break
        fi
        sleep 1
    done
}

walitProphecyFinish() {
    local count=0
    while true; do
        if [[ $# -eq 3 ]]; then
            ${CLI} relayer ethereum balance -o "${1}" -t "${2}"
            balance=$(${CLI} relayer ethereum balance -o "${1}" -t "${2}" | jq -r .balance)
            if [[ "${balance}" == "${3}" ]]; then
                break
            fi
        fi

        if [[ $# -eq 2 ]]; then
            ${CLI} relayer ethereum balance -o "${1}"
            balance=$(${CLI} relayer ethereum balance -o "${1}" | jq -r .balance)
            if [[ "${balance}" == "${2}" ]]; then
                break
            fi
        fi

        count=$((${count}+1))
        if [[ "${count}" == 30 ]]; then
            echo "failed to get balance"
            exit 1
        fi

        sleep 1
    done
}

cli_ret() {
    set +x
    if [[ $# -lt 2 ]]; then
        echo "wrong parameter"
        exit 1
    fi

    ok=$(echo "${1}" | jq -r .isOK)
    if [[ ${ok} != "true" ]]; then
        echo "failed to ${2}"
        exit 1
    fi

    local jqMsg=".msg"
    if [[ $# -ge 3 ]]; then
        jqMsg="${3}"
    fi

    msg=$(echo "${1}" | jq -r "${jqMsg}")
    if [[ $# -eq 4 ]]; then
         if [[ "${msg}" != "${4}" ]]; then
          echo "The balance is not correct"
          exit 1
        fi
    fi

    set -x
    echo "${msg}"
}

main () {
    InitAndDeploy
    TestETH2Chain33WithErc20Assets
    TestETH2Chain33WithEthAssets
    TestChain33ToEthAssets
}
main