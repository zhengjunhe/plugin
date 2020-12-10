#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -e
set -o pipefail

MAIN_HTTP=""

addr_A=19vpbRuz2XtKopQS2ruiVuVZeRdLd5n4t3
addr_B=1FcofeCgU1KYbB8dSa7cV2wjAF2RpMuUQD
source ../dapp-test-common.sh

hashlock_lock() {
    local secret=$1
    tx=$(curl -ksd '{"method":"Dplatform.CreateTransaction","params":[{"execer":"hashlock","actionName":"HashlockLock", "payload":{"secret":"'"${secret}"'","amount":1000000000, "time":75,"toAddr":"'"${addr_B}"'", "returnAddr":"'"${addr_A}"'","fee":100000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatform_SignAndSendTxWait "$tx" "0x1089b7f980fc467f029b7ae301249b36e3b582c911b1af1a24616c83b3563dcb" ${MAIN_HTTP} "$FUNCNAME"
}

hashlock_send() {
    local secret=$1
    tx=$(curl -ksd '{"method":"Dplatform.CreateTransaction","params":[{"execer":"hashlock","actionName":"HashlockSend", "payload":{"secret":"'"${secret}"'","fee":100000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatform_SignAndSendTxWait "$tx" "0xb76a398c3901dfe5c7335525da88fda4df24c11ad11af4332f00c0953cc2910f" ${MAIN_HTTP} "$FUNCNAME"
}

hashlock_unlock() {
    local secret=$1
    tx=$(curl -ksd '{"method":"Dplatform.CreateTransaction","params":[{"execer":"hashlock","actionName":"HashlockUnlock", "payload":{"secret":"'"${secret}"'","fee":100000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatform_SignAndSendTxWait "$tx" "0x1089b7f980fc467f029b7ae301249b36e3b582c911b1af1a24616c83b3563dcb" ${MAIN_HTTP} "$FUNCNAME"
}

init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"
    if [ "$ispara" == true ]; then
        hashlock_addr=$(curl -ksd '{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"user.p.para.hashlock"}]}' ${MAIN_HTTP} | jq -r ".result")
    else
        hashlock_addr=$(curl -ksd '{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"hashlock"}]}' ${MAIN_HTTP} | jq -r ".result")
    fi

    local main_ip=${MAIN_HTTP//8901/28803}
    dplatform_ImportPrivkey "0x1089b7f980fc467f029b7ae301249b36e3b582c911b1af1a24616c83b3563dcb" "19vpbRuz2XtKopQS2ruiVuVZeRdLd5n4t3" "hashlock1" "${main_ip}"
    dplatform_ImportPrivkey "0xb76a398c3901dfe5c7335525da88fda4df24c11ad11af4332f00c0953cc2910f" "1FcofeCgU1KYbB8dSa7cV2wjAF2RpMuUQD" "hashlock2" "$main_ip"

    local hashlock1="19vpbRuz2XtKopQS2ruiVuVZeRdLd5n4t3"
    local hashlock2="1FcofeCgU1KYbB8dSa7cV2wjAF2RpMuUQD"

    if [ "$ispara" == false ]; then
        dplatform_applyCoins "$hashlock1" 12000000000 "${main_ip}"
        dplatform_QueryBalance "${hashlock1}" "$main_ip"

        dplatform_applyCoins "$hashlock2" 12000000000 "${main_ip}"
        dplatform_QueryBalance "${hashlock2}" "$main_ip"
    else
        # tx fee
        dplatform_applyCoins "$hashlock1" 1000000000 "${main_ip}"
        dplatform_QueryBalance "${hashlock1}" "$main_ip"

        dplatform_applyCoins "$hashlock2" 1000000000 "${main_ip}"
        dplatform_QueryBalance "${hashlock2}" "$main_ip"
        local para_ip="${MAIN_HTTP}"
        #para chain import pri key
        dplatform_ImportPrivkey "0x1089b7f980fc467f029b7ae301249b36e3b582c911b1af1a24616c83b3563dcb" "19vpbRuz2XtKopQS2ruiVuVZeRdLd5n4t3" "hashlock1" "$para_ip"
        dplatform_ImportPrivkey "0xb76a398c3901dfe5c7335525da88fda4df24c11ad11af4332f00c0953cc2910f" "1FcofeCgU1KYbB8dSa7cV2wjAF2RpMuUQD" "hashlock2" "$para_ip"

        dplatform_applyCoins "$hashlock1" 12000000000 "${para_ip}"
        dplatform_QueryBalance "${hashlock1}" "$para_ip"
        dplatform_applyCoins "$hashlock2" 12000000000 "${para_ip}"
        dplatform_QueryBalance "${hashlock2}" "$para_ip"
    fi

    dplatform_SendToAddress "$hashlock1" "$hashlock_addr" 10000000000 ${MAIN_HTTP}
    dplatform_QueryExecBalance "${hashlock1}" "hashlock" "$MAIN_HTTP"
    dplatform_SendToAddress "$hashlock2" "$hashlock_addr" 10000000000 ${MAIN_HTTP}
    dplatform_QueryExecBalance "${hashlock2}" "hashlock" "$MAIN_HTTP"

    dplatform_BlockWait 1 "${MAIN_HTTP}"
}

function run_test() {
    dplatform_QueryBalance "$addr_A" "${MAIN_HTTP}"
    dplatform_QueryBalance "$addr_B" "${MAIN_HTTP}"
    hashlock_lock "abc"
    dplatform_QueryBalance "$addr_A" "${MAIN_HTTP}"
    hashlock_send "abc"
    dplatform_QueryBalance "$addr_B" "${MAIN_HTTP}"
    hashlock_unlock "abc"
    hashlock_lock "aef"
    dplatform_QueryBalance "$addr_A" "${MAIN_HTTP}"
    sleep 5
    hashlock_unlock "aef"
    dplatform_BlockWait 1 ${MAIN_HTTP}
    dplatform_QueryBalance "$addr_A" "${MAIN_HTTP}"
}

function main() {
    dplatform_RpcTestBegin hashlock
    MAIN_HTTP="$1"
    echo "ip=$MAIN_HTTP"

    init
    run_test
    dplatform_RpcTestRst hashlock "$CASE_ERR"
}

dplatform_debug_function main "$1"
