#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null

source "./dockerRelayerTest.sh"
#source "./perf_test.sh"

function cross2eth() {
    if [ "${2}" == "init" ]; then
        return
    elif [ "${2}" == "config" ]; then
        return
    elif [ "${2}" == "test" ]; then
        echo "========================== cross2eth test =========================="
        set +e
        set -x
        AllRelayerMainTest 20
#        perf_test_main 1
        echo "========================== cross2eth test end =========================="
    fi
}
