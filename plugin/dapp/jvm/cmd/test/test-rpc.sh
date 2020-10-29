#!/usr/bin/env bash
#shellcheck disable=SC2128
#shellcheck source=/dev/null
set -x
source ../dapp-test-common.sh


function rpc_test() {
    set +e
    set -x
    echo "rpc_test for dapp jvm"
}

chain33_debug_function rpc_test "$1" "$2"