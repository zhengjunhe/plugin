#!/usr/bin/env bash
set -x

source "./ebrelayer/publicTest.sh"

CLI="./build/ebcli_A"
Chain33Cli="/home/lyh/go/src/github.com/33cn/plugin/build/chain33-cli"

tokenAddr=""
BridgeRegistry=""
chain33SenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33SenderAddrKey="CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"

ethValidatorAddrKey="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"
ethReceiverAddr3="0x1919203bA8b325278d28Fb8fFeac49F2CD881A4e"
ethReceiverAddrKey3="62ca4122aac0e6f35bed02fc15c7ddbdaa07f2f2a1821c8b8210b891051e3ee9"

chain33Validator1="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33Validator2="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33Validator3="1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"
BtyReceiever="1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"
ETHContractAddr="0x0000000000000000000000000000000000000000"

InitAndDeploy() {
    echo "=========== $FUNCNAME begin ==========="
    result=$(${CLI} relayer set_pwd -n 123456hzj -o kk)
    cli_ret "${result}" "set_pwd"

    result=$(${CLI} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLI} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    echo "=========== $FUNCNAME end ==========="
}

function ImportKey() {
    echo "=========== $FUNCNAME begin ==========="
    result=$(${CLI} relayer set_pwd -n 123456hzj -o kk)

    result=$(${CLI} relayer unlock -p 123456hzj)
    #cli_ret "${result}" "unlock"

    result=$(${CLI} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
    cli_ret "${result}" "import_chain33privatekey"

    result=$(${CLI} relayer ethereum import_ethprivatekey -k "${ethValidatorAddrKey}")
    cli_ret "${result}" "import_ethprivatekey"

    result=$(${CLI} relayer chain33 import_privatekey -k "${ethValidatorAddrKey}")
    cli_ret "${result}" "import_ethprivatekey"

    echo "=========== $FUNCNAME end ==========="
}

function StartRelayerAndDeploy() {
    kill_ebrelayer "./build/ebrelayer"
    rm -rf "./build/datadir"
    rm -rf "./build/ebrelayer.log"
    start_ebrelayer "./build/ebrelayer" "./build/ebrelayer.log"

    InitAndDeploy
    ImportKey

    # 获取 BridgeRegistry 地址
    result=$(${CLI} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "bridgeRegistry" ".addr")
#    BridgeRegistry="0x5331F912027057fBE8139D91B225246e8159232f"

    kill_ebrelayer "./build/ebrelayer"
    # 修改 relayer.toml 配置文件
    updata_relayer_toml ${BridgeRegistry} "./build/relayer.toml"
    sed -i 's/initPowers=\[25, 25, 25, 25\]/initPowers=\[925, 25, 25, 25\]/g' "./build/relayer.toml"

    # 重启 ebrelayer 并解锁
    start_ebrelayer "./build/ebrelayer" "./build/ebrelayer.log"
    ${CLI} relayer unlock -p 123456hzj
}

# chian33 添加验证着及权重
function InitChain33Vilators() {
    echo "=========== $FUNCNAME begin ==========="
    # SetConsensusThreshold
    hash=$(${Chain33Cli} send x2ethereum setconsensus -p 80 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 2
    check_tx ${Chain33Cli} "${hash}"

    # add a validator
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator1} -p 87 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 2
    check_tx ${Chain33Cli} "${hash}"

    # add a validator again
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator2} -p 6 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 2
    check_tx ${Chain33Cli} "${hash}"

    # add a validator
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator3} -p 7 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 2
    check_tx ${Chain33Cli} "${hash}"

    # query Validators
    totalPower=$(${Chain33Cli} send x2ethereum query validators -v ${chain33Validator1} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g')
    check_number 87 ${totalPower}

    totalPower=$(${Chain33Cli} send x2ethereum query totalpower -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g')
    check_number 100 ${totalPower}

    # cions 转帐到 x2ethereum 合约地址
    hash=$(${Chain33Cli} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 2
    check_tx ${Chain33Cli} "${hash}"

    result=$(${Chain33Cli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)

    echo "=========== $FUNCNAME end ==========="
}

TestChain33ToEthAssets() {
    echo "=========== $FUNCNAME begin ==========="
    # token4chain33 在 以太坊 上先有 bty
    result=$(${CLI} relayer ethereum token4chain33 -s bty)
    tokenAddr=$(cli_ret "${result}" "token4chain33" ".addr")
#    tokenAddr="0x9C3D40A44a2F61Ef8D46fa8C7A731C08FB16cCEF"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # chain33 lock bty
    hash=$(${Chain33Cli} send x2ethereum lock -a 5 -t bty  -r ${ethReceiverAddr1} -q ${tokenAddr} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 4
    check_tx ${Chain33Cli} "${hash}"

    block_wait ${Chain33Cli} 10
    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    return

    # eth transfer
    {
        result=$(${CLI} relayer ethereum transfer -m 1.5 -k "${ethReceiverAddrKey1}" -r "${ethReceiverAddr2}" -t "${tokenAddr}")
        cli_ret "${result}" "transfer"

        result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
        cli_ret "${result}" "balance" ".balance" "3.5000"

        result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
        cli_ret "${result}" "balance" ".balance" "1.5000"
    }

    # eth burn
    result=$(${CLI} relayer ethereum burn -m 0.4 -k "${ethReceiverAddrKey2}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "burn"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1.1000"

    ${Chain33Cli} account balance -a "${chain33SenderAddr}" -e x2ethereum
    cli_ret "${result}" "balance" ".balance" "0.4000"

    result=$(${CLI} relayer ethereum burn -m 0.4 -k "${ethReceiverAddrKey2}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "burn"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0.7000"

    ${Chain33Cli} account balance -a "${chain33SenderAddr}" -e x2ethereum
    cli_ret "${result}" "balance" ".balance" "0.8000"

    echo "=========== $FUNCNAME end ==========="
}

# eth to chain33
# 在以太坊上锁定资产,然后在 chain33 上铸币,针对 eth 资产
TestETH2Chain33Assets() {
    echo "=========== $FUNCNAME begin ==========="
    ${CLI} relayer unlock -p 123456hzj

    result=$(${CLI} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth lock 100
    result=$(${CLI} relayer ethereum lock -m 5 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}")

    # chain33 burn 0.1eth
    hash=$(${Chain33CLI} send x2ethereum burn -a 0.1 -t eth  -r ${ethReceiverAddr1} -q ${ETHContractAddr} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait ${Chain33Cli} 4
    check_tx ${Chain33Cli} "${hash}"

    echo "=========== $FUNCNAME end ==========="
}

TestETH2Chain33Erc20() {
    echo "=========== $FUNCNAME begin ==========="

    ${CLI} relayer unlock -p 123456hzj
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

    # ETH 2 chain33 lock 100
    # -r chain33 receiver addr
    result=$(${CLI} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # chain33 burn 40
    hash=$(${Chain33CLI} send x2ethereum burn -a 40 -t "${tokenSymbol}"  -r ${ethReceiverAddr2} -q ${tokenAddr} -k "${chain33SenderAddr}")
    block_wait ${Chain33Cli} 4
    check_tx ${Chain33Cli} "${hash}"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "40"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "60"

    # burn 60
    hash=$(${Chain33CLI} send x2ethereum burn -a 60 -t "${tokenSymbol}"  -r ${ethReceiverAddr2} -q ${tokenAddr} -k "${chain33SenderAddr}")
    block_wait ${Chain33Cli} 4
    check_tx ${Chain33Cli} "${hash}"

    result=$(${CLI} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLI} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo "=========== $FUNCNAME end ==========="
}

function main() {
   # proxy_off
    StartRelayerAndDeploy
    InitChain33Vilators

    TestChain33ToEthAssets

    #TestETH2Chain33Assets
    #TestETH2Chain33Erc20
}

main





