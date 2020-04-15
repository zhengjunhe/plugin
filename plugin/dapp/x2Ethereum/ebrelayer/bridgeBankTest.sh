#!/usr/bin/env bash
set -x
CLI="../build/ebcli_A"
prophecyID=1
tokenAddr=""
bridgeBankAddr=""

InitAndDeploy() {
    result=$(${CLI} relayer set_pwd -n 123456hzj -o kk)
    cli_ret "${result}" "set_pwd"

    result=$(${CLI} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLI} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    # 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
    result=$(${CLI} relayer ethereum import_chain33privatekey -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944)
    cli_ret "${result}" "import_chain33privatekey"

    result=$(${CLI} relayer ethereum import_ethprivatekey -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230)
    cli_ret "${result}" "import_ethprivatekey"

    echo "Succeed to InitAndDeploy"
}

# 在以太坊上锁定资产,然后在 chain33 上铸币,针对 erc20 资产
# 以太坊 brun 资产,balance 对比是否正确
TestETH2Chain33WithErc20Assets() {
    # token4erc20 在 chain33 上先有 token,同时 mint
    tokenSymbol="testc"
    result=$(${CLI} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")

    # 先铸币 1000
    result=$(${CLI} relayer ethereum mint -m 1000 -o 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLI} relayer ethereum balance -o 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLI} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 2 chain33 lock 前先审批一下
    result=$(${CLI} relayer ethereum approve -m 100 -k 9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08 -t "${tokenAddr}")
    cli_ret "${result}" "approve"

    # ETH 2 chain33 lock
    # -r chain33 receiver addr
    result=$(${CLI} relayer ethereum lock -m 100 -k 9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08 -r 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # ETH 2 chain33 withdraw 40
    # -c 1 burn, and 2 lock
    result=$(${CLI} relayer ethereum prophecy -m 40 -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -c 1 -r 0xa4ea64a583f6e51c3799335b28a8f0529570a635 -s "${tokenSymbol}" -t "${tokenAddr}")
    cli_ret "${result}" "prophecy -m 40"

    result=$(${CLI} relayer ethereum process -i ${prophecyID})
    cli_ret "${result}" "process -i ${prophecyID}"
    prophecyID=${prophecyID}+1

    result=$(${CLI} relayer ethereum balance -o 0xa4ea64a583f6e51c3799335b28a8f0529570a635 -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "40"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "60"

    # ETH 2 chain33 withdraw 60
    result=$(${CLI} relayer ethereum prophecy -m 60 -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -c 1 -r 0xa4ea64a583f6e51c3799335b28a8f0529570a635 -s "${tokenSymbol}" -t "${tokenAddr}")
    cli_ret "${result}" "prophecy -m 60"

    result=$(${CLI} relayer ethereum process -i ${prophecyID})
    cli_ret "${result}" "process -i ${prophecyID}"
    prophecyID=${prophecyID}+1

    result=$(${CLI} relayer ethereum balance -o 0xa4ea64a583f6e51c3799335b28a8f0529570a635 -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -a "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo "Succeed to TestETH2Chain33WithErc20Assets"
}

TestETH2Chain33WithEthAssets() {
    result=$(${CLI} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLI} relayer ethereum balance -o 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf)
    balance1=$(cli_ret "${result}" "balance" ".balance")

    # ETH 2 chain33 lock eth
    result=$(${CLI} relayer ethereum lock -m 100 -k 9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08 -r 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt)
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf)
    balance2=$(cli_ret "${result}" "balance" ".balance")

    echo "${balance1} ${balance2}"
    if [[ "${balance1}" -ne $((${balance2}+100)) ]]; then
        echo "wrong balance"
        exit 1
    fi

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLI} relayer ethereum balance -o 0xa4ea64a583f6e51c3799335b28a8f0529570a635)
    balance3=$(cli_ret "${result}" "balance" ".balance")

    # ETH 2 chain33 withdraw 50
    result=$(${CLI} relayer ethereum prophecy -m 50 -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -c 1 -r 0xa4ea64a583f6e51c3799335b28a8f0529570a635 -s eth)
    cli_ret "${result}" "prophecy -m 50"

    result=$(${CLI} relayer ethereum process -i ${prophecyID})
    cli_ret "${result}" "process -i ${prophecyID}"
    prophecyID=${prophecyID}+1

    result=$(${CLI} relayer ethereum balance -o 0xa4ea64a583f6e51c3799335b28a8f0529570a635)
    balance4=$(cli_ret "${result}" "balance" ".balance")

    echo "${balance3} ${balance4}"
    if [[ "${balance4}" -ne $((${balance3}+50)) ]]; then
        echo "wrong balance"
        exit 1
    fi

    echo "Succeed to TestETH2Chain33WithEthAssets"
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
}

main