#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -o pipefail

MAIN_HTTP=""
GAME_ID=""
PASSWD="ABCD"
HASH_VALUE=$(echo -n "ABCD1" | sha256sum | awk '{print $1}')

PRIVA_A="0xfa21dc33a6144c546537580d28d894355d1e9af7292be175808b0f5737c30849"
PRIVA_B="0x213286d352b01fd740b6eaeb78a4fd316d743dd51d2f12c6789977430a41e0c7"

EXECTOR=""
source ../dapp-test-common.sh

function dplatform_GetExecAddr() {
    #获取GAME合约地址
    req='{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"'"$1"'"}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME"
}

function CreateGameTx() {
    local amount=$1
    local hash_value=$2
    local req='{"method":"Dplatform.CreateTransaction","params":[{"execer":"'"${EXECTOR}"'", "actionName":"createGame", "payload":{"amount": '"${amount}"',"hashType":"sha256","hashValue":"'"${hash_value}"'"}}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatform_SignAndSendTx "${RETURN_RESP}" "${PRIVA_A}" "${MAIN_HTTP}"
    GAME_ID=$RAW_TX_HASH

    echo_rst "CreateGame query_tx" "$?"
}

function MatchGameTx() {
    local gameId=$1
    local req='{"method":"Dplatform.CreateTransaction","params":[{"execer":"'"${EXECTOR}"'", "actionName":"matchGame", "payload":{"gameId": "'"${gameId}"'","guess":2}}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "MatchGame createRawTx" ".result"
    dplatform_SignAndSendTx "${RETURN_RESP}" "${PRIVA_B}" "${MAIN_HTTP}"
    echo_rst "MatchGame query_tx" "$?"
}

function CloseGameTx() {
    local gameId=$1
    local secret=$2
    local req='{"method":"Dplatform.CreateTransaction","params":[{"execer":"'"${EXECTOR}"'", "actionName":"closeGame", "payload":{"gameId": "'"${gameId}"'","secret":"'"${secret}"'","result":1}}]}'

    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "CloseGame createRawTx" ".result"
    dplatform_SignAndSendTx "${RETURN_RESP}" "${PRIVA_A}" "${MAIN_HTTP}"
    echo_rst "CloseGame query_tx" "$?"
}

function CancleGameTx() {
    local gameId=$1
    local req='{"method":"Dplatform.CreateTransaction","params":[{"execer":"'"${EXECTOR}"'", "actionName":"cancelGame", "payload":{"gameId": "'"${gameId}"'"}}]}'

    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "CancleGame createRawTx" ".result"
    dplatform_SignAndSendTx "${RETURN_RESP}" "${PRIVA_A}" "${MAIN_HTTP}"
    echo_rst "CancleGame query_tx" "$?"
}

function QueryGameByStatus() {
    local status=$1
    local req='{"method":"Dplatform.Query","params":[{"execer":"'"${EXECTOR}"'","funcName":"QueryGameListByStatusAndAddr","payload":{"status":'"${status}"',"address":""}}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME" ".result.games"
}

function QueryGameByGameId() {
    local gameId=$1
    local status=$2
    local req='{"method":"Dplatform.Query","params":[{"execer":"'"${EXECTOR}"'","funcName":"QueryGameById","payload":{"gameId":"'"${gameId}"'"}}]}'
    resok='(.error|not) and (.result.game.status = "'"${status}"'")'
    dplatform_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
}

function init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"

    local game_addr=""
    if [ "$ispara" == "true" ]; then
        EXECTOR="user.p.para.game"
        game_addr=$(curl -ksd '{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"user.p.para.game"}]}' ${MAIN_HTTP} | jq -r ".result")
    else
        EXECTOR="game"
        game_addr=$(curl -ksd '{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"game"}]}' ${MAIN_HTTP} | jq -r ".result")
    fi
    echo "gameAddr=${game_addr}"

    local main_ip=${MAIN_HTTP//8901/28803}
    #main chain import pri key
    dplatform_ImportPrivkey "$PRIVA_A" "16Z3haNPQd9wrnFDw19rtpbgnN2xynNT9f" "game1" "${main_ip}"
    dplatform_ImportPrivkey "$PRIVA_B" "16GXRfd9xj3XYMDti4y4ht7uzwoh55gZEc" "game2" "$main_ip"

    local ACCOUNT_A="16Z3haNPQd9wrnFDw19rtpbgnN2xynNT9f"
    local ACCOUNT_B="16GXRfd9xj3XYMDti4y4ht7uzwoh55gZEc"

    if [ "$ispara" == false ]; then
        dplatform_applyCoins "$ACCOUNT_A" 12000000000 "${main_ip}"
        dplatform_QueryBalance "${ACCOUNT_A}" "$main_ip"

        dplatform_applyCoins "$ACCOUNT_B" 12000000000 "${main_ip}"
        dplatform_QueryBalance "${ACCOUNT_B}" "$main_ip"
    else
        dplatform_applyCoins "$ACCOUNT_A" 1000000000 "${main_ip}"
        dplatform_QueryBalance "${ACCOUNT_A}" "$main_ip"

        dplatform_applyCoins "$ACCOUNT_B" 1000000000 "${main_ip}"
        dplatform_QueryBalance "${ACCOUNT_B}" "$main_ip"
        local para_ip="${MAIN_HTTP}"
        dplatform_ImportPrivkey "$PRIVA_A" "16Z3haNPQd9wrnFDw19rtpbgnN2xynNT9f" "game1" "$para_ip"
        dplatform_ImportPrivkey "$PRIVA_B" "16GXRfd9xj3XYMDti4y4ht7uzwoh55gZEc" "game2" "$para_ip"

        dplatform_applyCoins "$ACCOUNT_A" 12000000000 "${para_ip}"
        dplatform_QueryBalance "${ACCOUNT_A}" "$para_ip"
        dplatform_applyCoins "$ACCOUNT_B" 12000000000 "${para_ip}"
        dplatform_QueryBalance "${ACCOUNT_B}" "$para_ip"
    fi

    dplatform_SendToAddress "${ACCOUNT_B}" "$game_addr" 5000000000 "${MAIN_HTTP}"
    dplatform_QueryExecBalance "${ACCOUNT_B}" "game" "$MAIN_HTTP"
    dplatform_SendToAddress "${ACCOUNT_A}" "$game_addr" 5000000000 "${MAIN_HTTP}"
    dplatform_QueryExecBalance "${ACCOUNT_A}" "game" "$MAIN_HTTP"

    dplatform_BlockWait 1 "$MAIN_HTTP"
}

function run_test() {
    local ip=$1
    CreateGameTx 1000000000 "${HASH_VALUE}"
    QueryGameByGameId "${GAME_ID}" 1
    QueryGameByStatus 1
    MatchGameTx "${GAME_ID}"
    QueryGameByGameId "${GAME_ID}" 2
    QueryGameByStatus 2
    CloseGameTx "${GAME_ID}" "${PASSWD}"
    QueryGameByGameId "${GAME_ID}" 4
    QueryGameByStatus 4
    CreateGameTx 500000000 "${HASH_VALUE}"
    QueryGameByGameId "${GAME_ID}" 1
    CancleGameTx "${GAME_ID}"
    QueryGameByGameId "${GAME_ID}" 3
    QueryGameByStatus 3
}

function main() {
    dplatform_RpcTestBegin game
    local ip=$1
    MAIN_HTTP=$ip

    init
    run_test "$MAIN_HTTP"
    dplatform_RpcTestRst game "$CASE_ERR"
}

dplatform_debug_function main "$1"
