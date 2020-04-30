#!/usr/bin/env bash
#set -x
CLI="./ebcli_A"

ethSender0PrivateKey="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"

privateKeys[0]="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
privateKeys[1]="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
privateKeys[2]="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
privateKeys[3]="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
privateKeys[4]="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
privateKeys[5]="1385016736f7379884763f4a39811d1391fa156a7ca017be6afffa52bb327695"
privateKeys[6]="4ae589fe3837dcfc90d1c85b8423dc30841525cbebc41dfb537868b0f8376bbf"
privateKeys[7]="62ca4122aac0e6f35bed02fc15c7ddbdaa07f2f2a1821c8b8210b891051e3ee9"
privateKeys[8]="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
privateKeys[9]="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

perf_lock_eth() {
    count=10
    while true; do
        ethTxHash=$(${CLI} relayer ethereum lock-async -m 0.1 -k "${ethSender0PrivateKey}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
        echo $((10 - count)) "tx hash:" ${ethTxHash}
        count=$((count - 1))
        if [[ ${count} == 0  ]]; then
            break
        fi
    done
}



loop_send() {
    #打印数组长度
    echo ${#privateKeys[@]}

    #while 遍历数组
    i=0
    while  [[ i -lt ${#privateKeys[@]} ]];do
      #echo ${privateKeys[i]}
      ethTxHash=$(${CLI} relayer ethereum lock-async -m 0.1 -k "${privateKeys[i]}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
      echo ${i} "tx hash:" ${ethTxHash}
      let i++
    done
}

main () {
    perf_lock_eth
    loop_send
}

main