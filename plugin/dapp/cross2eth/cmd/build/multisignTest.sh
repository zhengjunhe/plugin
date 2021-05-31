#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"
source "./allRelayerPublic.sh"

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
chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
#ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
#ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10

Chain33Cli="../../chain33-cli"
chain33BridgeBank=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""
multisignAddr=""

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
}

function deployMultisign() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    for name in A B C D; do
        eval chain33MultisignKey=\$chain33MultisignKey${name}
        eval chain33Multisign=\$chain33Multisign${name}
        result=$(${Chain33Cli} account import_key -k "${chain33MultisignKey}" -l multisignAddr$name)
        # shellcheck disable=SC2154
        check_addr "${result}" "${chain33Multisign}"

        # chain33Multisign 要有手续费
        hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Multisign}" -k "${chain33DeployAddr}")
        check_tx "${Chain33Cli}" "${hash}"
        result=$(${Chain33Cli} account balance -a "${chain33Multisign}" -e coins)
        balance_ret "${result}" "10.0000"
    done

    echo 'multisign１:部署离线钱包合约'
    result=$(${CLIA} chain33 multisign deploy)
    cli_ret "${result}" "chain33 multisign deploy"
    multisignAddr=$(echo "${result}" | jq -r ".msg")

    result=$(${CLIA} chain33 multisign setup -k "${chain33DeployKey}" -o "${chain33MultisignA},${chain33MultisignB},${chain33MultisignC},${chain33MultisignD}")
    cli_ret "${result}" "chain33 multisign setup"

    # multisignAddr 要有手续费
    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${multisignAddr}" -k "${chain33DeployAddr}")
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${multisignAddr}" -e coins)
    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lockBty() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    echo '1:#配置离线钱包地址'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configOfflineSaveAccount(${multisignAddr})")
    check_tx "${Chain33Cli}" "${hash}"

    echo '2:#配置自动转离线钱包(bty, 1000, 50%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(1111111111111111111114oLvT2,BTY,100000000000,50)")
    check_tx "${Chain33Cli}" "${hash}"

    echo '3:#在chain33侧lock bty, 不需要addlock，执行完成之后，在chain33侧的multisign没有增加 bty 余额'
    hash=$(${Chain33Cli} evm call -f 1 -a 330 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, 1111111111111111111114oLvT2, 33000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignAddr}" -e evm)
    balance_ret "${result}" "0.0000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "330.0000"

    echo '4:#在#执行完成之后，在chain33侧的multisign增加了 bty 余额，具体的数量　＝　(之前执行该笔交易执行的chain33BridgeBank的BTY余额　+ 800 ) * 50%'
    hash=$(${Chain33Cli} evm call -f 1 -a 800 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, 1111111111111111111114oLvT2, 80000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignAddr}" -e evm)
    balance_ret "${result}" "565.0000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "565.0000"

    hash=$(${Chain33Cli} evm call -f 1 -a 500 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, 1111111111111111111114oLvT2, 50000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignAddr}" -e evm)
    balance_ret "${result}" "1097.5000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "532.5000"

    # transfer test
    hash=$(./ebcli_A chain33 multisign transfer -a 100 -r "${chain33BridgeBank}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${multisignAddr}" -e evm)


    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)

    result=$(${Chain33Cli} account balance -a "${chain33MultisignA}" -e evm)

    hash=$(./ebcli_A chain33 multisign transfer -a 100 -r "${chain33MultisignA}" -k "${chain33MultisignKeyA},${chain33MultisignKeyB},${chain33MultisignKeyC},${chain33MultisignKeyD}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33MultisignA}" -e evm)

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lockEth() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"




    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function lockYcc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} chain33 token create -s YCC -o 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt)
    cli_ret "${result}" "chain33 token create -s YCC"
    chain33YccErc20Addr=$(echo "${result}" | jq -r .msg)

    cp ./ERC20.abi "${chain33YccErc20Addr}.abi"

    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33DeployAddr})"

#    echo '1:#配置离线钱包地址'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configOfflineSaveAccount(${multisignAddr})")
    check_tx "${Chain33Cli}" "${hash}"

#    echo '2:#配置自动转离线钱包(YCC, 1000, 60%)'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "configLockedTokenOfflineSave(${chain33YccErc20Addr},YCC,100000000000,60)")
    check_tx "${Chain33Cli}" "${hash}"

#    echo 'YCC.0:增加allowance的设置,或者使用relayer工具进行'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33YccErc20Addr}" -p "approve(${chain33BridgeBank}, 330000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    echo 'YCC.1:#在chain33侧lock YCC, 因为需要提前addlock，所以lock失败,chain33BridgeBank的ＹＣＣ余额没有发生变化'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, ${chain33YccErc20Addr}, 50000000000)")

    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33DeployAddr})"
    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})"

#    echo 'YCC.2:#执行add lock操作:addToken2LockList'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "addToken2LockList(${chain33YccErc20Addr}, YCC)")
    check_tx "${Chain33Cli}" "${hash}"

#    echo '3:#在chain33侧lock YCC, 执行完成之后，在chain33侧的multisign没有增加YCC余额'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, ${chain33YccErc20Addr}, 50000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33DeployAddr})"
    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})"

#    echo '6:#在#执行完成之后，在chain33侧的multisign增加了ＢＴＹ余额，具体的数量　＝　(之前执行该笔交易执行的chain33BridgeBank的BTY余额　+ 800 ) * 50%'
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e ${chain33BridgeBank} -p "lock(8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, ${chain33YccErc20Addr}, 60000000000)")
    check_tx "${Chain33Cli}" "${hash}"

    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33DeployAddr})"
    ${Chain33Cli} evm abi call -a "${chain33YccErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33DeployAddr})"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function mainTest() {
    StartChain33

    start_trufflesuite

    kill_all_ebrelayer
    StartRelayerAndDeploy

    deployMultisign
#    lockBty
    lockYcc


}

mainTest
