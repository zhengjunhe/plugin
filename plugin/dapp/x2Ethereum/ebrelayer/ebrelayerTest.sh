#!/usr/bin/env bash
set -x
CLIA="../build/ebcli_A"
CLIB="../build/ebcli_B"
CLIC="../build/ebcli_C"
CLID="../build/ebcli_D"


tokenAddr=""
bridgeBankAddr=""
chain33SenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
chain33SenderAddrKey="CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
ethOperatorAddrKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"

ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

ethReceiverAddr3="0x1919203bA8b325278d28Fb8fFeac49F2CD881A4e"
ethReceiverAddrKey3="62ca4122aac0e6f35bed02fc15c7ddbdaa07f2f2a1821c8b8210b891051e3ee9"

ethValidatorAddrA="0x92c8b16afd6d423652559c6e266cbe1c29bfd84f"
ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
ethValidatorAddrB="0x0df9a824699bc5878232c9e612fe1a5346a5a368"
ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrC="0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1"
ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrD="0xd9dab021e74ecf475788ed7b61356056b2095830"
ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"

prophecyTx0="0x772260c98aec81b3e235af47c355db720f60e751cce100fed6f334e1b1530bde"

# 初始化部署合约
InitAndDeploy() {
    echo "=========== $FUNCNAME begin ==========="
    # 创建文件夹及拷贝
    mkdir '../build/A' '../build/B' '../build/C' '../build/D'
    cp '../build/relayer.toml' '../build/A/relayer.toml'
    cp '../build/ebrelayer' '../build/A/ebrelayer'
    startEbrelayer "./../build/A/ebrelayer" "./../build/A/ebrelayer.log"

    result=$(${CLIA} relayer set_pwd -n 123456hzj -o kk)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    result=$(${CLIA} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
    cli_ret "${result}" "import_chain33privatekey"

    result=$(${CLIA} relayer ethereum import_ethprivatekey -k "${ethValidatorAddrKeyA}")
    cli_ret "${result}" "import_ethprivatekey"

    echo "=========== $FUNCNAME end ==========="
}

# 初始化 B C D 文件夹下的文容
function InitConfigFile() {
    echo "=========== $FUNCNAME begin ==========="
    # 获取 BridgeRegistry 地址
    result=$(${CLIA} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "token4chain33" ".addr")
    port=9901
    for name in B C D
    do
        file="../build/"$name"/relayer.toml"
        cp '../build/relayer.toml' "${file}"
        cp '../build/ebrelayer' "../build/"$name"/ebrelayer"

        # 删除配置文件中不需要的字段
        for deleteName in "BridgeRegistry" "operatorAddr" "validatorsAddr" "initPowers"
        do
            deleteLine "${file}" "${deleteName}"
        done

        # 在第 5 行后面 新增合约地址
        sed -i '5 a BridgeRegistry="'${BridgeRegistry}'"' "${file}"

        # 替换端口
        port=$((${port} + 1))
        sed -i 's/localhost:9901/localhost:'${port}'/g' "${file}"
    done
    echo "=========== $FUNCNAME end ==========="
}

# 启动 B C D 的 ebrelayer 服务,导入私钥
function ImportCBDKey() {
    echo "=========== $FUNCNAME begin ==========="
    for name in B C D
    do
        startEbrelayer "./../build/"$name"/ebrelayer" "./../build/"$name"/ebrelayer.log"

        # 导入测试地址私钥
        CLI="../build/ebcli_$name"

        result=$(${CLI} relayer set_pwd -n 123456hzj -o kk)
        cli_ret "${result}" "set_pwd"

        result=$(${CLI} relayer unlock -p 123456hzj)
        cli_ret "${result}" "unlock"

        result=$(${CLI} relayer ethereum import_chain33privatekey -k "${chain33SenderAddrKey}")
        cli_ret "${result}" "import_chain33privatekey"
    done

    result=$(${CLIB} relayer ethereum import_ethprivatekey -k "${ethValidatorAddrKeyB}")
    cli_ret "${result}" "import_ethprivatekeyB"

    result=$(${CLIC} relayer ethereum import_ethprivatekey -k "${ethValidatorAddrKeyC}")
    cli_ret "${result}" "import_ethprivatekeyC"

    result=$(${CLID} relayer ethereum import_ethprivatekey -k "${ethValidatorAddrKeyD}")
    cli_ret "${result}" "import_ethprivatekeyD"
    echo "=========== $FUNCNAME end ==========="
}

# chain33 到 eth,chian33 lock 100,必须 A B C D 中有三个都lock,才能成功
TestChain33ToEth() {
    echo "=========== $FUNCNAME begin ==========="
    # token4chain33 在 以太坊 上先有 bty
    result=$(${CLIA} relayer ethereum token4chain33 -s btyi)
    tokenAddr=$(cli_ret "${result}" "token4chain33" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    for name in A B C D
    do
        CLI="../build/ebcli_$name"
        # -c 2 chain33 lock 100
        result=$(${CLI} relayer ethereum prophecy -i "${prophecyTx0}" -m 1 -a "${chain33SenderAddr}" -c 2 -r "${ethReceiverAddr1}" -s bty -t "${tokenAddr}")
        cli_ret "${result}" "prophecy -m 1"
        #sleep 15

        if [[ "${name}" == "A" || "${name}" == "B" ]]; then
            result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
            cli_ret "${result}" "balance" ".balance" "0"
        fi

        if [[ "${name}" == "C" || "${name}" == "D" ]]; then
            result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
            cli_ret "${result}" "balance" ".balance" "1"
        fi
    done
    echo "=========== $FUNCNAME end ==========="
}

# 解锁
function unlock() {
    for name in A B C D
    do
        CLI="../build/ebcli_$name"
        result=$(${CLI} relayer unlock -p 123456hzj)
    done
}

# 判断结果是否正确
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

# 判断结果是否错误
cli_ret_err() {
    #set +x
    ok=$(echo "${1}" | jq -r .isOK)
    echo "${ok}"
    if [[ "${ok}" == "true" ]]; then
        echo "isOK is true"
        exit 1
    fi
    #set -x
}

# 查询关键字所在行然后删除 ${1}文件名称 ${2}关键字
function deleteLine() {
    line=$(cat -n "${1}"|grep "${2}"|awk '{print $1}')
        if [[ "${line}" != "" ]]; then
            sed -i "${line}"'d' "${1}" # 删除行
        fi
}

# 后台启动 ebrelayer 进程 $1进程名称 $2进程信息输出重定向文件
function startEbrelayer() {
    # 参数如果小于 2 直接报错
    if [[ $# -lt 2 ]]; then
        echo "wrong parameter"
        return
    fi

    # 判断可执行文件是否存在
    if [ ! -x "${1}" ];then
        echo "${1} not exist"
        return
    fi

    # 后台启动程序
    nohup "${1}" >"${2}" 2>&1 &
    sleep 1

    pid=$(ps -ef | grep "${1}" | grep -v 'grep' | awk '{print $2}')
    if [ "${pid}" == "" ];then
        echo "start ${1} failed"
        return
    fi
}

# 杀死进程ebrelayer 进程 $1进程名称
function killEbrelayer() {
    pid=$(ps -ef | grep "${1}" | grep -v 'grep' | awk '{print $2}')
    if [ "${pid}" == "" ];then
        echo "not find ${1} pid"
        return
    fi

    kill "${pid}"
    pid=$(ps -ef | grep "${1}" | grep -v 'grep' | awk '{print $2}')
    if [ "${pid}" != "" ];then
        echo "kill ${1} failed"
        kill -9 "${pid}"
    fi
}

function KillAllEbrelayer() {
    for name in A B C D
    do
        CLI="./../build/$name/ebrelayer"
        killEbrelayer "${CLI}"
    done
}

main () {
#    InitAndDeploy
#
#    InitConfigFile
#    ImportCBDKey
#    TestChain33ToEth

    KillAllEbrelayer

}
main