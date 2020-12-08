#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"

dplatformSenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
# 新增地址 dplatform 需要导入地址 转入 10 bty当收费费
dplatformValidator1="1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
dplatformValidator2="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
dplatformValidator3="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
dplatformValidator4="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
dplatformValidatorKey1="0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
dplatformValidatorKey2="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
dplatformValidatorKey3="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
dplatformValidatorKey4="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10
tokenAddrBty=""
tokenAddr=""
ethUrl=""
DplatformCli=""

function kill_ebrelayerC() {
    #shellcheck disable=SC2154
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayerc_1"
}
function kill_ebrelayerD() {
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayerd_1"
}

function start_ebrelayerA() {
    docker cp "./relayer.toml" "${dockerNamePrefix}_ebrelayera_1":/root/relayer.toml
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1" "/root/ebrelayer" "./ebrelayera.log"
    sleep 5
}

function start_ebrelayerC() {
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayerc_1" "/root/ebrelayer" "./ebrelayerc.log"
    sleep 5
    ${CLIC} relayer unlock -p 123456hzj
    sleep 5
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"
    sleep 10
}
function start_ebrelayerD() {
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayerd_1" "/root/ebrelayer" "./ebrelayerd.log"
    sleep 5
    ${CLID} relayer unlock -p 123456hzj
    sleep 5
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"
    sleep 10
}

function InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} relayer set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartRelayerAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # change EthProvider url
    dockerAddr=$(get_docker_addr "${dockerNamePrefix}_ganachetest_1")
    ethUrl="http://${dockerAddr}:8545"

    # 修改 relayer.toml 配置文件
    updata_relayer_a_toml "${dockerAddr}" "${dockerNamePrefix}_ebrelayera_1" "./relayer.toml"
    # start ebrelayer A
    start_ebrelayerA
    # 部署合约
    InitAndDeploy

    # 获取 BridgeRegistry 地址
    result=$(${CLIA} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "bridgeRegistry" ".addr")

    # kill ebrelayer A
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1"
    sleep 1

    # 修改 relayer.toml 配置文件
    updata_relayer_toml "${BridgeRegistry}" ${maturityDegree} "./relayer.toml"
    # 重启
    start_ebrelayerA

    # start ebrelayer B C D
    for name in b c d; do
        local file="./relayer$name.toml"
        cp './relayer.toml' "${file}"

        # 删除配置文件中不需要的字段
        for deleteName in "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers" "deployerPrivateKey" "deploy"; do
            delete_line "${file}" "${deleteName}"
        done

        sed -i 's/x2ethereum/x2ethereum'${name}'/g' "${file}"

        pushHost=$(get_docker_addr "${dockerNamePrefix}_ebrelayer${name}_1")
        line=$(delete_line_show "${file}" "pushHost")
        sed -i ''"${line}"' a pushHost="http://'"${pushHost}"':20000"' "${file}"

        line=$(delete_line_show "${file}" "pushBind")
        sed -i ''"${line}"' a pushBind="'"${pushHost}"':20000"' "${file}"

        docker cp "${file}" "${dockerNamePrefix}_ebrelayer${name}_1":/root/relayer.toml
        start_docker_ebrelayer "${dockerNamePrefix}_ebrelayer${name}_1" "/root/ebrelayer" "./ebrelayer${name}.log"
    done
    sleep 5

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function EthImportKey() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 重启 ebrelayer 并解锁
    for name in a b c d; do
        # 导入测试地址私钥
        # shellcheck disable=SC2154
        CLI="docker exec ${dockerNamePrefix}_ebrelayer${name}_1 /root/ebcli_A"

        result=$(${CLI} relayer set_pwd -p 123456hzj)

        result=$(${CLI} relayer unlock -p 123456hzj)
        cli_ret "${result}" "unlock"
    done

    result=$(${CLIA} relayer ethereum import_dplatformprivatekey -k "${dplatformValidatorKey1}")
    cli_ret "${result}" "import_dplatformprivatekey"
    result=$(${CLIB} relayer ethereum import_dplatformprivatekey -k "${dplatformValidatorKey2}")
    cli_ret "${result}" "import_dplatformprivatekey"
    result=$(${CLIC} relayer ethereum import_dplatformprivatekey -k "${dplatformValidatorKey3}")
    cli_ret "${result}" "import_dplatformprivatekey"
    result=$(${CLID} relayer ethereum import_dplatformprivatekey -k "${dplatformValidatorKey4}")
    cli_ret "${result}" "import_dplatformprivatekey"

    result=$(${CLIA} relayer dplatform import_privatekey -k "${ethValidatorAddrKeyA}")
    cli_ret "${result}" "A relayer dplatform import_privatekey"
    result=$(${CLIB} relayer dplatform import_privatekey -k "${ethValidatorAddrKeyB}")
    cli_ret "${result}" "B relayer dplatform import_privatekey"
    result=$(${CLIC} relayer dplatform import_privatekey -k "${ethValidatorAddrKeyC}")
    cli_ret "${result}" "C relayer dplatform import_privatekey"
    result=$(${CLID} relayer dplatform import_privatekey -k "${ethValidatorAddrKeyD}")
    cli_ret "${result}" "D relayer dplatform import_privatekey"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chian33 添加验证着及权重
function InitDplatformVilators() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 导入 dplatformValidators 私钥生成地址
    result=$(${DplatformCli} account import_key -k ${dplatformValidatorKey1} -l validator1)
    check_addr "${result}" ${dplatformValidator1}
    result=$(${DplatformCli} account import_key -k ${dplatformValidatorKey2} -l validator2)
    check_addr "${result}" ${dplatformValidator2}
    result=$(${DplatformCli} account import_key -k ${dplatformValidatorKey3} -l validator3)
    check_addr "${result}" ${dplatformValidator3}
    result=$(${DplatformCli} account import_key -k ${dplatformValidatorKey4} -l validator4)
    check_addr "${result}" ${dplatformValidator4}

    # SetConsensusThreshold
    hash=$(${DplatformCli} send x2ethereum setconsensus -p 80 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"

    # add a validator
    hash=$(${DplatformCli} send x2ethereum add -a ${dplatformValidator1} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    hash=$(${DplatformCli} send x2ethereum add -a ${dplatformValidator2} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    hash=$(${DplatformCli} send x2ethereum add -a ${dplatformValidator3} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    hash=$(${DplatformCli} send x2ethereum add -a ${dplatformValidator4} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"

    # query Validators
    totalPower=$(${DplatformCli} send x2ethereum query totalpower -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g')
    check_number 100 "${totalPower}"

    # cions 转帐到 x2ethereum 合约地址
    hash=$(${DplatformCli} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)

    check_tx "${DplatformCli}" "${hash}"
    result=$(${DplatformCli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "200.0000"

    # dplatformValidator 要有手续费
    hash=$(${DplatformCli} send coins transfer -a 10 -t "${dplatformValidator1}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    result=$(${DplatformCli} account balance -a "${dplatformValidator1}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${DplatformCli} send coins transfer -a 10 -t "${dplatformValidator2}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    result=$(${DplatformCli} account balance -a "${dplatformValidator2}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${DplatformCli} send coins transfer -a 10 -t "${dplatformValidator3}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    result=$(${DplatformCli} account balance -a "${dplatformValidator3}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${DplatformCli} send coins transfer -a 10 -t "${dplatformValidator4}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${DplatformCli}" "${hash}"
    result=$(${DplatformCli} account balance -a "${dplatformValidator4}" -e coins)
    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestDplatformToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # token4dplatform 在 以太坊 上先有 bty
    result=$(${CLIA} relayer ethereum token4dplatform -s coins.bty)
    tokenAddrBty=$(cli_ret "${result}" "token4dplatform" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # dplatform lock bty
    hash=$(${DplatformCli} send x2ethereum lock -a 5 -t coins.bty -r ${ethReceiverAddr1} -q "${tokenAddrBty}" --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    result=$(${DplatformCli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "195.0000"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 5 -k "${ethReceiverAddrKey1}" -r "${dplatformSenderAddr}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} account balance -a "${dplatformSenderAddr}" -e x2ethereum)
    balance_ret "${result}" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to dplatform
# 在以太坊上锁定资产,然后在 dplatform 上铸币,针对 eth 资产
function TestETH2DplatformAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.1 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0.1"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    hash=$(${DplatformCli} send x2ethereum burn -a 0.1 -t eth -r ${ethReceiverAddr2} --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    result=$(${DplatformCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.1" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2DplatformErc20() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20 在 dplatform 上先有 token,同时 mint
    tokenSymbol="testc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")

    # 先铸币 1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${dplatformValidator1}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} x2ethereum balance -s "${dplatformValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    # dplatform burn 100
    hash=$(${DplatformCli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr}" --node_addr "${ethUrl}" -k "${dplatformValidator1}")
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    result=$(${DplatformCli} x2ethereum balance -s "${dplatformValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestDplatformToEthAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [ "${tokenAddrBty}" == "" ]; then
        # token4dplatform 在 以太坊 上先有 bty
        result=$(${CLIA} relayer ethereum token4dplatform -s coins.bty)
        tokenAddrBty=$(cli_ret "${result}" "token4dplatform" ".addr")
    fi

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # dplatform lock bty
    hash=$(${DplatformCli} send x2ethereum lock -a 5 -t coins.bty -r ${ethReceiverAddr2} -q "${tokenAddrBty}" --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 5 -k "${ethReceiverAddrKey2}" -r "${dplatformValidator1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} account balance -a "${dplatformValidator1}" -e x2ethereum)
    balance_ret "${result}" "0"

    start_ebrelayerD

    result=$(${DplatformCli} account balance -a "${dplatformValidator1}" -e x2ethereum)
    balance_ret "${result}" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to dplatform
# 在以太坊上锁定资产,然后在 dplatform 上铸币,针对 eth 资产
function TestETH2DplatformAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.1 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${DplatformCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0.1"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    kill_ebrelayerC
    kill_ebrelayerD

    hash=$(${DplatformCli} send x2ethereum burn -a 0.1 -t eth -r ${ethReceiverAddr2} --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    result=$(${DplatformCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.1" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2DplatformErc20Kill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20 在 dplatform 上先有 token,同时 mint
    tokenSymbol="testcc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr2=$(cli_ret "${result}" "token4erc20" ".addr")

    # 先铸币 1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${dplatformValidator1}" -t "${tokenAddr2}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${DplatformCli} x2ethereum balance -s "${dplatformValidator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${DplatformCli} x2ethereum balance -s "${dplatformValidator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    kill_ebrelayerC
    kill_ebrelayerD

    # dplatform burn 100
    hash=$(${DplatformCli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr2}" --node_addr "${ethUrl}" -k "${dplatformValidator1}")
    block_wait "${DplatformCli}" $((maturityDegree + 2))
    check_tx "${DplatformCli}" "${hash}"

    result=$(${DplatformCli} x2ethereum balance -s "${dplatformValidator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function AllRelayerMainTest() {
    set +e
    docker_dplatform_ip=$(get_docker_addr "${dockerNamePrefix}_dplatform_1")
    DplatformCli="./dplatform-cli --rpc_laddr http://${docker_dplatform_ip}:8801"

    CLIA="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A"
    CLIB="docker exec ${dockerNamePrefix}_ebrelayerb_1 /root/ebcli_A"
    CLIC="docker exec ${dockerNamePrefix}_ebrelayerc_1 /root/ebcli_A"
    CLID="docker exec ${dockerNamePrefix}_ebrelayerd_1 /root/ebcli_A"
    echo "${CLIA}"

    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    # init
    StartRelayerAndDeploy
    InitDplatformVilators
    EthImportKey

    # test
    TestDplatformToEthAssets
    TestETH2DplatformAssets
    TestETH2DplatformErc20

    # kill relayer and start relayer
    TestDplatformToEthAssetsKill
    TestETH2DplatformAssetsKill
    TestETH2DplatformErc20Kill

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
