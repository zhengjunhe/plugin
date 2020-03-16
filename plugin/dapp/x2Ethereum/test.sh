#!/usr/bin/env bash
set -x
CLI="/opt/src/github.com/33cn/plugin/build/chain33-cli"

# 测试流程
#
# bankAddr = 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi
#
#
# 1.Exec_EthBridgeClaim, 测试chain33这端
#
#
#

block_wait() {
    if [[ $# -lt 1 ]]; then
        echo "wrong block_wait parameter"
        exit 1
    fi
    cur_height=$(${CLI} block last_header | jq ".height")
    expect=$((cur_height + ${1}))
    local count=0
    while true; do
        new_height=$(${CLI} block last_header | jq ".height")
        if [[ ${new_height} -ge ${expect} ]]; then
            break
        fi
        count=$((count + 1))
        sleep 1
    done
    sleep 1
    count=$((count + 1))
    echo "wait new block $count s, cur height=$expect,old=$cur_height"
}

# check_tx(ty, hash)
# ty：交易执行结果
#   1：交易打包
#   2：交易执行成功
# hash：交易hash
check_tx() {
    if [[ $# -lt 2 ]]; then
        echo "wrong check_tx parameters"
        exit 1
    fi
    ty=$(${CLI} tx query -s ${2} | jq .receipt.ty)
    if [[ ${ty} != ${1} ]]; then
        echo "check tx error, hash is ${2}"
        exit 1
    fi
}

check_balance() {
    if [[ $# -lt 2 ]]; then
        echo "wrong check_balance parameters"
        exit 1
    fi
    if [[ ${1} != ${2} ]]; then
        echo "error balance, expect ${1}, get ${2}"
        exit 1
    fi
}

# 转币到lns合约
#${CLI} send coins transfer -a 500 -t 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv
#block_wait 2
#${CLI} send coins send_exec -e x2ethereum -a 200 -k 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi
#block_wait 2
#balance=`${CLI} account balance -e x2ethereum -a 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi | jq .balance | sed 's/\"//g'`
#check_balance "200.0000" ${balance}
#echo "check balance on chain ok"

#${CLI} x2ethereum login -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -p 7

# SetConsensusNeeded
hash=`${CLI} send x2ethereum setconsensus -p 0.7 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query consensusNeeded
preConsensus=`${CLI} send x2ethereum query consensus -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .consensusNeed | sed 's/\"//g'`
check_balance 0.7 ${preConsensus}

# SetConsensusNeeded
hash=`${CLI} send x2ethereum setconsensus -p 0.8 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query consensusNeeded
nowConsensus=`${CLI} send x2ethereum query consensus -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .consensusNeed | sed 's/\"//g'`
check_balance 0.8 ${nowConsensus}

# login a address
hash=`${CLI} send x2ethereum login -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -p 7 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# login a address again
hash=`${CLI} send x2ethereum login -a 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -p 6 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}