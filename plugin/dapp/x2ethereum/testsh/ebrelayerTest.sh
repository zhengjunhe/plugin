#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x

# 只启动 eth 这一端的测试
# 启动4个 relayer  每个权重一样

source "./publicTest.sh"
CLIA="../build/ebcli_A"

tokenAddr=""
chain33SenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33SenderAddrKey="CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
prophecyTx0="0x772260c98aec81b3e235af47c355db720f60e751cce100fed6f334e1b1530bde"

# 初始化部署合约
InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 创建文件夹及拷贝
    rm -rf '../build/A' '../build/B' '../build/C' '../build/D'
    mkdir '../build/A' '../build/B' '../build/C' '../build/D'
    cp '../build/relayer.toml' '../build/A/relayer.toml'
    cp '../build/ebrelayer' '../build/A/ebrelayer'
    start_ebrelayer "./../build/A/ebrelayer" "./../build/A/ebrelayer.log"

    result=$(${CLIA} relayer set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "deploy"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# 初始化 B C D 文件夹下的文容
function InitConfigFile() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 获取 BridgeRegistry 地址
    result=$(${CLIA} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "bridgeRegistry" ".addr")

    port=9901
    for name in B C D; do
        file="../build/$name/relayer.toml"
        cp '../build/relayer.toml' "${file}"
        cp '../build/ebrelayer' "../build/$name/ebrelayer"

        # 删除配置文件中不需要的字段
        for deleteName in "BridgeRegistry" "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers" "deployerPrivateKey" "deploy"; do
            delete_line "${file}" "${deleteName}"
        done

        # 在第 5 行后面 新增合约地址
        sed -i '5 a BridgeRegistry="'"${BridgeRegistry}"'"' "${file}"

        # 替换端口
        port=$((port + 1))
        sed -i 's/localhost:9901/localhost:'${port}'/g' "${file}"
    done
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# 启动 B C D 的 ebrelayer 服务,导入私钥
function ImportCBDKey() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    for name in B C D; do
        start_ebrelayer "./../build/$name/ebrelayer" "./../build/$name/ebrelayer.log"

        # 导入测试地址私钥
        CLI="../build/ebcli_$name"

        result=$(${CLI} relayer set_pwd -p 123456hzj)
        cli_ret "${result}" "set_pwd"

        result=$(${CLI} relayer unlock -p 123456hzj)
        cli_ret "${result}" "unlock"

        result=$(${CLI} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
        cli_ret "${result}" "import_chain33privatekey"
    done

    result=$(${CLIA} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
    cli_ret "${result}" "import_chain33privatekey"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 到 eth,chian33 lock 100,必须 A B C D 中有三个都lock,才能成功
TestChain33ToEth() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # token4chain33 在 以太坊 上先有 bty
    result=$(${CLIA} relayer ethereum token4chain33 -s coins.bty)
    tokenAddr=$(cli_ret "${result}" "token4chain33" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "coins.bty")
    cli_ret "${result}" "balance" ".balance" "0"

    for name in A B C D; do
        CLI="../build/ebcli_$name"
        # -c 2 chain33 lock 100
        result=$(${CLI} relayer ethereum prophecy -i "${prophecyTx0}" -m 100 -a "${chain33SenderAddr}" -c 2 -r "${ethReceiverAddr1}" -s coins.bty -t "${tokenAddr}")
        cli_ret "${result}" "prophecy -m 1"

        if [[ ${name} == "A" || ${name} == "B" ]]; then
            result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "coins.bty")
            cli_ret "${result}" "balance" ".balance" "0"
        elif [[ ${name} == "C" || ${name} == "D" ]]; then
            result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "coins.bty")
            cli_ret "${result}" "balance" ".balance" "100"
        fi
    done

    # transfer 10
    result=$(${CLIA} relayer ethereum transfer -m 10 -k "${ethReceiverAddrKey1}" -r "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "transfer"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "coins.bty")
    cli_ret "${result}" "balance" ".balance" "90"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "coins.bty")
    cli_ret "${result}" "balance" ".balance" "10"

    # burn 90
    result=$(${CLIA} relayer ethereum burn -m 90 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "burn"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

main() {
    kill_all_ebrelayer

    InitAndDeploy

    InitConfigFile
    ImportCBDKey
    TestChain33ToEth

    kill_all_ebrelayer
}
main