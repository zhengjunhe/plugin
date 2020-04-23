#!/usr/bin/env bash
set -x
CLI="/opt/src/github.com/33cn/plugin/build/chain33-cli"

ETHContractAddr="0x40BFE5eD039A9a2Eb42ece2E2CA431bFa7Cf4c42"
BTYContractAddr="0x40BFE5eD039A9a2Eb42ece2E2CA431bFa7Cf4c42"
Ethsender="0xcbfddc6ae318970ba3feeb0541624f95822e413a"
BtyReceiever="1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"
BtyTokenContractAddr="0xbAf2646b8DaD8776fc74Bf4C8d59E6fB3720eddf"

Validator1="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
Validator2="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
Validator3="1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"


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
        echo "wrong #block_wait parameter"
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

check_Number() {
    if [[ $# -lt 2 ]]; then
        echo "wrong check_Number parameters"
        exit 1
    fi
    if [[ ${1} != ${2} ]]; then
        echo "error Number, expect ${1}, get ${2}"
        exit 1
    fi
}

# SetConsensusThreshold
hash=`${CLI} send x2ethereum setconsensus -p 80 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query ConsensusThreshold
nowConsensus=`${CLI} send x2ethereum query consensus -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .nowConsensusThreshold | sed 's/\"//g'`
check_Number 80 ${nowConsensus}

# add a validator
hash=`${CLI} send x2ethereum add -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -p 7 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query Validators
validators=`${CLI} send x2ethereum query validators -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .validators | sed 's/\"//g'`
totalPower=`${CLI} send x2ethereum query validators -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g'`
check_Number 7 ${totalPower}

# add a validator again
hash=`${CLI} send x2ethereum add -a 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -p 6 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query Validators
validators=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .validators | sed 's/\"//g'`
totalPower=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g'`
check_Number 13 ${totalPower}

# remove a validator
hash=`${CLI} send x2ethereum remove -a 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query Validators
validators=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .validators | sed 's/\"//g'`
totalPower=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g'`
check_Number 7 ${totalPower}

# add a validator again
hash=`${CLI} send x2ethereum add -a 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -p 6 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query Validators
validators=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .validators | sed 's/\"//g'`
totalPower=`${CLI} send x2ethereum query validators -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g'`
check_Number 13 ${totalPower}

totalPower=`${CLI} send x2ethereum query totalpower -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g'`
check_Number 13 ${totalPower}

##############################################################
######################## 测试交易 #############################
##############################################################

# ethereum -> chain33

# send a eth2chain33 tx
hash=`${CLI} send x2ethereum create --amount 10 -b ${ETHContractAddr} -e x2ethereum --claimtype 1 -t eth --ethid 0 --nonce 0 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a eth2chain33 tx again
hash=`${CLI} send x2ethereum create --amount 10 -b ${ETHContractAddr} -e x2ethereum --claimtype 1 -t eth --ethid 0 --nonce 0 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`

# TODO
# 后续可以做一笔转账

# send a burn tx
hash=`${CLI} send x2ethereum burn --amount 4 -e x2ethereum -t eth  -r ${Ethsender} -s ${BtyReceiever} -q 0x0000000000000000000000000000000000000000 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a eth2chain33 tx
hash=`${CLI} send x2ethereum create --amount 5 -b ${ETHContractAddr} -e x2ethereum --claimtype 1 -t eth --ethid 0 --nonce 1 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a eth2chain33 tx again
hash=`${CLI} send x2ethereum create --amount 5 -b ${ETHContractAddr} -e x2ethereum --claimtype 1 -t eth --ethid 0 --nonce 1 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`

# chain33 -> ethereum

# send bty to mavl-x2ethereum-bty-12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv
${CLI} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv

# send a chain33eth tx
hash=`${CLI} send x2ethereum lock --amount 5 -e x2ethereum -t bty  -r ${Ethsender} -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -q ${BtyTokenContractAddr} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a withdrawEth tx
hash=`${CLI} send x2ethereum withdraweth --amount 2 -b ${ETHContractAddr} -e x2ethereum --claimtype 2 -t bty  --ethid 0 --nonce 2 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a withdrawEth tx again
hash=`${CLI} send x2ethereum withdraweth --amount 2 -b ${ETHContractAddr} -e x2ethereum --claimtype 2 -t bty --ethid 0 --nonce 2 -r ${BtyReceiever} -s ${Ethsender} -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`


# send a withdrawEth tx
hash=`${CLI} send x2ethereum withdraweth --amount 2 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 2 -t bty --ethid 0 --nonce 3 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a withdrawEth tx again
hash=`${CLI} send x2ethereum withdraweth --amount 2 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 2 -t bty --ethid 0 --nonce 3 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
