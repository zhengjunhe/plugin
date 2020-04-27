#!/usr/bin/env bash
set -x

# 解锁
function unlock_relayer() {
    for name in A B C D
    do
        local CLI="../build/ebcli_$name"
        ${CLI} relayer unlock -p 123456hzj
    done
}

# 判断结果是否正确
function cli_ret() {
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
#          exit 1
        fi
    fi

    set -x
    echo "${msg}"
}

# 判断结果是否错误
function cli_ret_err() {
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
function delete_line() {
    line=$(cat -n "${1}"|grep "${2}"|awk '{print $1}')
    if [[ "${line}" != "" ]]; then
        sed -i "${line}"'d' "${1}" # 删除行
    fi
}

# 查询关键字所在行然后删除 ${1}文件名称 ${2}关键字
function delete_line_show() {
    local line=$(cat -n "${1}"|grep "${2}"|awk '{print $1}')
    if [[ "${line}" != "" ]]; then
        sed -i "${line}"'d' "${1}" # 删除行
        line=$((line - 1))
    fi
    echo "${line}"
}

# 后台启动 ebrelayer 进程 $1进程名称 $2进程信息输出重定向文件
function start_ebrelayer() {
    # 参数如果小于 2 直接报错
    if [[ $# -lt 2 ]]; then
        echo "wrong parameter"
        exit 1
    fi

    # 判断可执行文件是否存在
    if [ ! -x "${1}" ];then
        echo "${1} not exist"
        exit 1
    fi

    # 后台启动程序
    nohup "${1}" >"${2}" 2>&1 &
    sleep 1

    pid=$(ps -ef | grep "${1}" | grep -v 'grep' | awk '{print $2}')
    if [ "${pid}" == "" ];then
        echo "start ${1} failed"
        exit 1
    fi
}

# 杀死进程ebrelayer 进程 $1进程名称
function kill_ebrelayer() {
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
    sleep 1
}

function kill_all_ebrelayer() {
    for name in A B C D
    do
        local ebrelayer="./../build/$name/ebrelayer"
        kill_ebrelayer "${ebrelayer}"
    done
}

# 区块等待 $1:cli 路径  $2:等待高度
function block_wait() {
    local CLI=${1}

    if [[ $# -lt 1 ]]; then
        echo "wrong block_wait parameter"
        exit 1
    fi

    local cur_height=$(${CLI} block last_header | jq ".height")
    local expect=$((cur_height + ${2}))
    local count=0
    while true; do
        new_height=$(${CLI} block last_header | jq ".height")
        if [[ ${new_height} -ge ${expect} ]]; then
            break
        fi

        count=$((count + 1))
        sleep 1

        if [[ ${count} -ge 30 ]]; then
           break
        fi
    done

    count=$((count + 1))
    echo "wait new block $count s, cur height=$expect,old=$cur_height"
}

# 检查交易是否执行成功 $1:cli 路径  $2:交易hash
function check_tx() {
    local CLI=${1}

    if [[ $# -lt 2 ]]; then
        echo "wrong check_tx parameters"
        exit 1
    fi
    ty=$(${CLI} tx query -s ${2} | jq .receipt.ty)
    if [[ ${ty} != 2 ]]; then
        echo "check tx error, hash is ${2}"
        exit 1
    fi
}

function check_number() {
    if [[ $# -lt 2 ]]; then
        echo "wrong check number parameters"
        exit 1
    fi
    if [[ ${1} != ${2} ]]; then
        echo "error number, expect ${1}, get ${2}"
        exit 1
    fi
}

# 更新配置文件 $1 为 BridgeRegistry 合约地址 $2 relayer.toml 地址
function updata_relayer_toml() {
    local BridgeRegistry=${1}
    local file=${2}

    local chain33Host=$(docker inspect build_chain33_1 | jq ".[].NetworkSettings.Networks.build_default.IPAddress" | sed 's/\"//g')
    local pushHost=$(ifconfig eth0 | grep "inet " | awk '{ print $2}' | awk -F: '{print $2}')

    local line=$(delete_line_show ${file} "chain33Host")
    # 在第 line 行后面 新增合约地址
    sed -i ''${line}' a chain33Host="http://'${chain33Host}':8801"' "${file}"

    line=$(delete_line_show ${file} "pushHost")
    sed -i ''${line}' a pushHost="http://'${pushHost}':20000"' "${file}"

    line=$(delete_line_show ${file} "BridgeRegistry")
    sed -i ''${line}' a BridgeRegistry="'${BridgeRegistry}'"' "${file}"

    #sed -i 's/#BridgeRegistry=\"0x40BFE5eD039A9a2Eb42ece2E2CA431bFa7Cf4c42\"/BridgeRegistry=\"'${BridgeRegistry}'\"/g' "./build/relayer.toml"
    #sed -i 's/192.168.64.2/'${chain33Host}'/g' "./build/relayer.toml"
    #sed -i 's/192.168.3.156/'${pushHost}'/g' "./build/relayer.toml"
}

# $1 CLI
function wait_prophecy_finish() {
    local CLI=${1}
    set +x
    local count=0
    while true; do
        if [[ $# -eq 4 ]]; then
            ${CLI} relayer ethereum balance -o "${2}" -t "${3}"
            balance=$(${CLI} relayer ethereum balance -o "${2}" -t "${3}" | jq -r .balance)
            if [[ "${balance}" == "${4}" ]]; then
                break
            fi
        fi
        if [[ $# -eq 3 ]]; then
            ${CLI} relayer ethereum balance -o "${2}"
            balance=$(${CLI} relayer ethereum balance -o "${2}" | jq -r .balance)
            if [[ "${balance}" == "${3}" ]]; then
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
    set -x
}