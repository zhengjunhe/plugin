#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
source ../dapp-test-common.sh

MAIN_HTTP=""
txhash=""

function query_unfreezeID() {
    dplatform_BlockWait 1 "$MAIN_HTTP"

    local times=10
    while true; do
        req='{"method":"Dplatform.QueryTransaction","params":[{"hash":"'"$txhash"'"}]}'
        ret=$(curl -ksd "$req" ${MAIN_HTTP})
        tx=$(jq -r ".result.tx.hash" <<<"$ret")
        #   echo "====query tx= ${txhash}, return=$ret "
        if [ "${tx}" != "${txhash}" ]; then
            dplatform_BlockWait 1 "${MAIN_HTTP}"
            times=$((times - 1))
            if [ $times -le 0 ]; then
                echo -e "${RED}====query tx=$txhash failed${NOC}"
                echo "req=$req"
                curl -ksd "$req" ${MAIN_HTTP}
                exit 1
            fi
        else
            uid=${tx#0x}
            echo -e "${GRE}====query tx=$txhash  success${NOC}"
            break
        fi
    done
}

function init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"
    exec_name="unfreeze"
    symbol="dpom"

    beneficiary_key=0xa2ec1c6274723c021daa8792f4d0d52ffa0eff0fd47c9c6c1d1dd618762dc178
    beneficiary=1qpAv7H4C5JBgVQffDRbQKti7ibdM2TfU

    owner=1CK51xZ1wNkrzAhGyDuFayxeQXHg3gqcVS
    owner_key=0x3b0d7f65b35da1c394891ba7a8ce0f070ccef6818e3f7ca9c203776013b3a4b0

    dplatform_ImportPrivkey "${beneficiary_key}" "${beneficiary}" "unfreeze_beneficiary" "${MAIN_HTTP}"
    dplatform_ImportPrivkey "${owner_key}" "${owner}" "unfreeze_owner" "${MAIN_HTTP}"

    dplatform_applyCoins "${beneficiary}" 10000000000 "${MAIN_HTTP}"
    dplatform_applyCoins "${owner}" 10000000000 "${MAIN_HTTP}"

    if [ "$ispara" == true ]; then
        exec_name="user.p.para."${exec_name}
        symbol="para"

        local main_ip=${MAIN_HTTP//8901/28803}
        dplatform_applyCoins "${beneficiary}" 10000000000 "${main_ip}"
        dplatform_applyCoins "${owner}" 10000000000 "${main_ip}"
    fi

    exec_addr=$(curl -ksd '{"method":"Dplatform.ConvertExectoAddr","params":[{"execname":"'${exec_name}'"}]}' ${MAIN_HTTP} | jq -r ".result")
    echo "exec_addr=${exec_addr}"
    dplatform_SendToAddress "$owner" "$exec_addr" 5000000000 "${MAIN_HTTP}"
    dplatform_SendToAddress "$beneficiary" "$exec_addr" 5000000000 "${MAIN_HTTP}"
    dplatform_BlockWait 1 "${MAIN_HTTP}"
}

function CreateRawUnfreezeCreate() {
    req='{"jsonrpc": "2.0", "method" :  "unfreeze.CreateRawUnfreezeCreate" , "params":[{"startTime":10000,"assetExec":"coins","assetSymbol":"'$symbol'","totalCount":400000000,"beneficiary":"'$beneficiary'","means":"FixAmount","fixAmount": {"period":10,"amount":1000000}}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatform_SignAndSendTx "$RETURN_RESP" "$owner_key" "${MAIN_HTTP}"
    query_unfreezeID
}

function CreateRawUnfreezeWithdraw() {
    sleep 10
    req='{"method":"unfreeze.CreateRawUnfreezeWithdraw","params":[{"unfreezeID":"'${uid}'"}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatform_SignAndSendTx "$RETURN_RESP" "${beneficiary_key}" "${MAIN_HTTP}"
}

function CreateRawUnfreezeTerminate() {
    req='{"method":"unfreeze.CreateRawUnfreezeTerminate","params":[{"unfreezeID":"'${uid}'"}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatform_SignAndSendTx "$RETURN_RESP" "$owner_key" "${MAIN_HTTP}"
    dplatform_BlockWait 2 "${MAIN_HTTP}"
}

function GetUnfreeze() {
    req='{"method":"unfreeze.GetUnfreeze","params":[{"data":"'${uid}'"}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME"
}

function GetUnfreezeWithdraw() {
    req='{"method":"unfreeze.GetUnfreezeWithdraw","params":[{"data":"'${uid}'"}]}'
    dplatform_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME"
}

function run_testcases() {
    CreateRawUnfreezeCreate
    CreateRawUnfreezeWithdraw
    GetUnfreeze
    GetUnfreezeWithdraw
    CreateRawUnfreezeTerminate
}

function rpc_test() {
    dplatform_RpcTestBegin unfreeze
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    init
    run_testcases
    dplatform_RpcTestRst unfreeze "$CASE_ERR"
}

dplatform_debug_function rpc_test "$1"
