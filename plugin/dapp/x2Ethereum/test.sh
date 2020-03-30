#!/usr/bin/env bash
set -x
CLI="/opt/src/github.com/33cn/plugin/cli/chain33-cli"

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

# send a eth2chain33 tx
hash=`${CLI} send x2ethereum create --amount 10 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 0 -t eth -g eth --ethid 0 --nonce 0 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a eth2chain33 tx again
hash=`${CLI} send x2ethereum create --amount 10 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 0 -t eth -g eth --ethid 0 --nonce 0 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`

Amount=`${CLI} send x2ethereum query symbolamount -s eth -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq '.totalAmount' | sed 's/\"//g'`
check_Number 10 ${Amount}

# send a withdrawEth tx
hash=`${CLI} send x2ethereum withdraweth --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 1 -t eth -g eth --ethid 0 --nonce 1 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a withdrawEth tx again
hash=`${CLI} send x2ethereum withdraweth --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 1 -t eth -g eth --ethid 0 --nonce 1 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`

Amount=`${CLI} send x2ethereum query symbolamount -s eth -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq '.totalAmount' | sed 's/\"//g'`
check_Number 5 ${Amount}

# send a eth2chain33 tx
hash=`${CLI} send x2ethereum create --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 0 -t eth -g eth --ethid 0 --nonce 2 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# send a eth2chain33 tx again
hash=`${CLI} send x2ethereum create --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum --claimtype 0 -t eth -g eth --ethid 0 --nonce 2 -r 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -s 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -q 0x0000000000000000000000000000000000000000 -v 1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

# query prophecy
ProphecyID=`${CLI} tx query -s ${hash}  | jq '.receipt.logs[].log | select((.ProphecyID !="") and (.ProphecyID !=null)) | .ProphecyID' | sed 's/\"//g'`
Prophecy=`${CLI} send x2ethereum query prophecy -i ${ProphecyID} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`

Amount=`${CLI} send x2ethereum query symbolamount -s eth -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq '.totalAmount' | sed 's/\"//g'`
check_Number 10 ${Amount}

# send bty to mavl-x2ethereum-bty-12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv
${CLI} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv

# send a chain33eth tx
hash=`${CLI} send x2ethereum lock --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum -t bty  -g eth --ethid 0 --nonce 0 -r 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -q 0x0000000000000000000000000000000000000000 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

Amount=`${CLI} send x2ethereum query symbolamount -s eth -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq '.totalAmount' | sed 's/\"//g'`
check_Number 10 ${Amount}

# send a chain33eth tx
hash=`${CLI} send x2ethereum burn --amount 5 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -e x2ethereum -t bty -g eth --ethid 0 --nonce 0 -r 0x7B95B6EC7EbD73572298cEf32Bb54FA408207359 -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -q 0x0000000000000000000000000000000000000000 -b 0xC4cE93a5699c68241fc2fB503Fb0f21724A624BB -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv`
#block_wait 2
check_tx 2 ${hash}

Amount=`${CLI} send x2ethereum query symbolamount -s eth -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq '.totalAmount' | sed 's/\"//g'`
check_Number 10 ${Amount}
