#!/usr/bin/env bash
#shellcheck disable=SC2128
#shellcheck source=/dev/null
set -x
source ../dapp-test-common.sh

contract="user.jvm.Guess"
#privkey="CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
privkey="0x4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"
jvm_privkey="0x7b2800cdecd978ab0e877f7e3734b9d0b11d864fa51d9b623d7bdbd76c16a40d"
MAIN_HTTP=""
exector="jvm"
#main_addr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
main_addr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
jvm_addr="1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU"

function init() {
    req='{"method":"Chain33.ImportPrivkey", "params":[{"privkey":"'${jvm_privkey}'", "label":"player"}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result.label=="player") and (.result.acc.addr == "'${jvm_addr}'")' "$FUNCNAME"

    #转账
    echo "send coins transfer -a 1000"
    chain33_SendToAddress "${main_addr}" "${jvm_addr}" 100000000000 ${MAIN_HTTP}
}

function create_contract() {
    echo "Begin to test contract $contract"
    local CLI="docker exec ${dockerNamePrefix}_chain33_1 /root/chain33-cli"
    code=$(${CLI} jvm code -x Guess -d ./ | jq -r ".code")

    local req='{"method":"Chain33.CreateTransaction","params":[{"execer":"'"${exector}"'", "actionName":"CreateJvmContract","payload":{"name": "'"${contract}"'","code":"'"${code}"'"}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "CreateJvmContract" ".result"
    chain33_SignAndSendTx "${RETURN_RESP}" "${privkey}" "${MAIN_HTTP}"
    echo_rst "CreateJvmContract query_tx" "$?"
}

function transfer() {
    echo "send coins send_exec -a 30"
    contract_addr=$(curl -ksd '{"method":"Chain33.ConvertExectoAddr","params":[{"execname":"'${contract}'"}]}' ${MAIN_HTTP} | jq -r ".result")
    #2  存钱到合约地址
    tx_hash=$(curl -ksd '{"method":"Chain33.CreateRawTransaction","params":[{"to":"'"$contract_addr"'","amount":3000000000,"note":"send2exec"}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$tx_hash" "$jvm_privkey" ${MAIN_HTTP}
}

function start_game() {
    #开始游戏
    echo "send jvm call -e $contract -x startGame"
    local req='{"method":"Chain33.CreateTransaction","params":[{"execer":"'"${contract}"'", "actionName":"CallJvmContract", "payload":{"Name": "'"${contract}"'","actionData":["startGame"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "startGame" ".result"
    chain33_SignAndSendTx "${RETURN_RESP}" "$privkey" ${MAIN_HTTP}
}

function play_game() {
    #投注
    echo "send jvm call -e $contract -x playGame"
    local req='{"method":"Chain33.CreateTransaction","params":[{"execer":"'"${contract}"'", "actionName":"CallJvmContract","payload":{"Name": "'"${contract}"'","actionData":["playGame", "6", "2"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "playGame" ".result"
    chain33_SignAndSendTx "${RETURN_RESP}" "$jvm_privkey" ${MAIN_HTTP}
}

function close_game() {
    chain33_BlockWait 12 ${MAIN_HTTP}

    echo "close $contract"
    local req='{"method":"Chain33.CreateTransaction","params":[{"execer":"'"${contract}"'", "actionName":"CallJvmContract","payload":{"Name": "'"${contract}"'","actionData":["closeGame"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "closeGame" ".result"
    chain33_SignAndSendTx "${RETURN_RESP}" "$jvm_privkey" ${MAIN_HTTP}

    chain33_BlockWait 10 ${MAIN_HTTP}
}

function query() {
    #查看信息
    local req='{"method":"Chain33.Query","params":[{"execer":"'"${exector}"'", "funcName":"JavaContract","payload":{"contract": "'"${contract}"'","para":["getGuessRecordByRound","1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU","1"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "Query" ".result"
    check=$(echo "${RETURN_RESP}" | jq -r ".result")
    if [ "${check}" != "["'guessNum=6,ticketNum=2'"]" ]; then
        echo -e "${RED}error query via get${contract}RecordByRound, expect guessNum=6,ticketNum=2 , get $RETURN_RESP${NOC}"
    fi

    local req='{"method":"Chain33.Query","params":[{"execer":"'"${exector}"'", "funcName":"JavaContract","payload":{"contract": "'"${contract}"'","para":["getBonusByRound","1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU","1"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "Query" ".result"

    local req='{"method":"Chain33.Query","params":[{"execer":"'"${exector}"'", "funcName":"JavaContract","payload":{"contract": "'"${contract}"'","para":["getLuckNumByRound","1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU","1"]}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "Query" ".result"
}

function rpc_test() {
    set +e
    set -x
    chain33_RpcTestBegin jvm
    MAIN_HTTP="$1"
    dockerNamePrefix="$2"
    echo "main_ip=$MAIN_HTTP"

    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    if [ "$ispara" == false ]; then
#        local req='{"method":"Chain33.Query","params":[{"execer":"'"${exector}"'", "funcName":"CheckContractNameExist","payload":{"JvmContractName": "'"${contract}"'"}}]}'
#        contract_exist=$(curl -ksd "$req" ${MAIN_HTTP} | jq -r ".result.existAlready")
#        if [ "${contract_exist}" != "true" ]; then
            init
            create_contract
            transfer
            start_game
            play_game
            close_game
#        fi
        query
    fi
    chain33_RpcTestRst jvm "$CASE_ERR"
}

chain33_debug_function rpc_test "$1" "$2"
