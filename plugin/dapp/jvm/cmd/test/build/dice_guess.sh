#!/usr/bin/env bash
#set -x

#cli="docker exec build_chain33_1 ./chain33-cli"
cli="./chain33-cli"

# chain33 区块等待 $1:cli 路径  $2:等待高度
function block_wait() {
    set +x
    local CLI=${1}

    if [[ $# -lt 1 ]]; then
        echo -e "${RED}wrong block_wait parameter${NOC}"
        exit_cp_file
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
    done

    count=$((count + 1))
    set -x
    echo -e "${GRE}chain33 wait new block $count s, cur height=$expect,old=$cur_height${NOC}"
}

#创建账户，并充值
function setupAccount() {
    $cli account import_key -l player -k 0x7b2800cdecd978ab0e877f7e3734b9d0b11d864fa51d9b623d7bdbd76c16a40d

    echo "transfer to 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
    $cli send  coins transfer -a 1000 -n "t1000" -t 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv
    echo 'transfer to 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU'
    $cli send  coins transfer -a 100 -n "t100" -t 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
}

#部署合约
function deployJavaContract() {
    for contract in Guess Dice
    do
        echo "Deploy contract for $contract"
        $cli send jvm create -x $contract -n "deploy $contract" -d . -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
    done
}

function depositAndStartGame() {
    for contract in Guess Dice
    do
        echo "transfer to user.jvm.$contract"
        $cli send coins send_exec -a 20 -e user.jvm.$contract -n send2exec -k 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU

        #开始游戏
        echo 'send tx to startGame'
        $cli send jvm call -e $contract -x startGame -n "call $contract startGame" -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt

        #投注
        echo 'send tx to playGame'
        $cli send jvm call -e $contract -x playGame -r "6 2" -n "call $contract playGame" -k 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU
    done
    block_wait $cli 12
}



function closeGame() {
    for contract in Guess Dice
    do
        echo "close $contract"
        $cli send jvm call -e $contract -x closeGame -n "call $contract closeGame" -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
    done
    block_wait $cli 1
}

expectQueryRes[0]="guessNum=6,ticketNum=2"
expectQueryRes[1]="diceNum=6,ticketNum=2"

function queryGame() {
    i=0
    for contract in Guess Dice
    do
        echo "query get${contract}RecordByRound"
        result=$($cli jvm query -e user.jvm.$contract -r "get${contract}RecordByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1")
        if [ "${result}" == "${privateKeys[i]}" ]; then
            echo "Succeed to do query from user.jvm.$contract"
        else
            echo -e "${RED}error query via get${contract}RecordByRound, expect "${privateKeys[i]}", get ${result}${NOC}"
        fi
        let i++

        echo "query getBonusByRound"
        $cli jvm query -e user.jvm.$contract -r "getBonusByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1"

        echo "query getLuckNumByRound"
        $cli jvm query -e user.jvm.$contract -r "getLuckNumByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1"
    done
}

function dice_game_test() {
    setupAccount
    deployJavaContract
    depositAndStartGame
    closeGame
    queryGame
}