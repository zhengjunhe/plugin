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

# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
#ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
# shellcheck disable=SC2034
#{
#ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
#ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
#ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
#}

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
#chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

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

CLIA="./ebcli_A"

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
    balance_ret "${result}" "5.0000"

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    # eth 这端 金额是否增加了 5
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    kill_ebrelayerC

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    # eth 这端 金额是否减少了 3
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "0.0000"

    start_ebrelayerC
    start_ebrelayerD

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

    kill_ebrelayerC
    kill_ebrelayerD

    # ETH 这端 lock 11个
    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2

    # 查询 ETH 这端 bridgeBank 地址 11
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "11"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额 原来是0
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11 * le8
    is_equal "${result}" "0"

    start_ebrelayerC

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
    start_ebrelayerD

    # 比之前多 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "105"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11-5 * le8
    is_equal "${result}" "600000000"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "6"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33YccKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # ETH 这端 lock 7个 YCC
    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    eth_block_wait 2
    sleep ${maturityDegree}

    # 查询 ETH 这端 bridgeBank 地址 7 YCC
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额 地址原来是 0
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7 * le8
    is_equal "${result}" "0"

    start_ebrelayerC

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
    start_ebrelayerD

    # 更新后的金额 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7-5 * le8
    is_equal "${result}" "200000000"

    # 查询 ETH 这端 bridgeBank 地址 2
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function mainTest() {
    StartChain33
    start_trufflesuite
    AllRelayerStart

    TestChain33ToEthAssetsKill
    TestETH2Chain33AssetsKill
    TestETH2Chain33YccKill
}

mainTest
