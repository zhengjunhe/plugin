#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"
source "./relayerPublic.sh"

# ETH 部署合约者的私钥 用于部署合约时签名使用
ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# chain33 部署合约者的私钥 用于部署合约时签名使用
chain33DeployAddr="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"
chain33DeployKey="0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"

maturityDegree=10

chain33BridgeBank=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""
BridgeRegistryOnChain33=""
chain33YccErc20Addr=""
BridgeRegistryOnEth=""
ethBridgeToeknYccAddr=""
chain33ZBCErc20Addr=""
ethBridgeToeknZBCAddr=""
chain33ID=33

# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
#ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
# shellcheck disable=SC2034
{
ethValidatorAddrKeyb="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrKeyc="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrKeyd="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"

# 新增地址 chain33 需要导入地址 转入 10 bty当收费费
chain33Validatora="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"
chain33Validatorb="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
chain33Validatorc="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
chain33Validatord="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
#chain33ValidatorKeya="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"
chain33ValidatorKeyb="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
chain33ValidatorKeyc="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
chain33ValidatorKeyd="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
}

#ethUrl=""


function start_docker_ebrelayerA() {
    # shellcheck disable=SC2154
    docker cp "./relayer.toml" "${dockerNamePrefix}_ebrelayera_1":/root/relayer.toml
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1" "/root/ebrelayer" "./ebrelayera.log"
    sleep 5
}

# start ebrelayer B C D
function updata_toml_start_bcd() {
    for name in b c d; do
        local file="./relayer$name.toml"
        cp './relayer.toml' "${file}"

        # 删除配置文件中不需要的字段
        for deleteName in "deploy4chain33" "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers" "deploy" "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers"; do
            delete_line "${file}" "${deleteName}"
        done

        pushNameChange "${file}"

        pushHost=$(get_docker_addr "${dockerNamePrefix}_ebrelayer${name}_1")
        line=$(delete_line_show "${file}" "pushHost")
        sed -i ''"${line}"' a pushHost="http://'"${pushHost}"':20000"' "${file}"

        line=$(delete_line_show "${file}" "pushBind")
        sed -i ''"${line}"' a pushBind="'"${pushHost}"':20000"' "${file}"

        docker cp "${file}" "${dockerNamePrefix}_ebrelayer${name}_1":/root/relayer.toml
        start_docker_ebrelayer "${dockerNamePrefix}_ebrelayer${name}_1" "/root/ebrelayer" "./ebrelayer${name}.log"

        CLI="docker exec ${dockerNamePrefix}_ebrelayer${name}_1 /root/ebcli_A"
        result=$(${CLI} set_pwd -p 123456hzj)
        cli_ret "${result}" "set_pwd"

        result=$(${CLI} unlock -p 123456hzj)
        cli_ret "${result}" "unlock"

        eval chain33ValidatorKey=\$chain33ValidatorKey${name}
        # shellcheck disable=SC2154
        result=$(${CLI} chain33 import_privatekey -k "${chain33ValidatorKey}")
        cli_ret "${result}" "chain33 import_privatekey"

        eval ethValidatorAddrKey=\$ethValidatorAddrKey${name}
        # shellcheck disable=SC2154
        result=$(${CLI} ethereum import_privatekey -k "${ethValidatorAddrKey}")
        cli_ret "${result}" "ethereum import_privatekey"
    done
}

function StartDockerRelayerDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 修改 relayer.toml 配置文件 pushName 字段
    pushNameChange "./relayer.toml"
    # 修改 relayer.toml 配置文件 initPowers
    validators_config

    # change EthProvider url
    dockerAddr=$(get_docker_addr "${dockerNamePrefix}_ganachetest_1")
#    ethUrl="http://${dockerAddr}:8545"

    # 修改 relayer.toml 配置文件
    updata_relayer_a_toml "${dockerAddr}" "${dockerNamePrefix}_ebrelayera_1" "./relayer.toml"

    # para
    # shellcheck disable=SC2155
    local line=$(delete_line_show "./relayer.toml" "chain33Host")
    # 在第 line 行后面 新增合约地址
    docker_chain33_ip=$(get_docker_addr "${dockerNamePrefix}_chain33_1")
    sed -i ''"${line}"' a chain33Host="http://'"${docker_chain33_ip}"':8901"' "./relayer.toml"

    # shellcheck disable=SC2155
    local line=$(delete_line_show "./relayer.toml" "ChainName")
    # 在第 line 行后面 新增合约地址
    sed -i ''"${line}"' a ChainName="user.p.para."' "./relayer.toml"

    # 启动 ebrelayer
    start_docker_ebrelayerA

    # 导入私钥 部署合约 设置 bridgeRegistry 地址
    InitAndDeploy

    docker cp "${BridgeRegistryOnChain33}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${BridgeRegistryOnChain33}.abi
    docker cp "${chain33BridgeBank}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chain33BridgeBank}.abi
    docker cp "${BridgeRegistryOnEth}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${BridgeRegistryOnEth}.abi
    docker cp "${ethBridgeBank}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${ethBridgeBank}.abi

    # 重启
    # kill ebrelayer A
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1"
    sleep 1
#    kill_ebrelayer ebrelayer
    start_docker_ebrelayerA

    result=$(${CLIA} unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    # start ebrelayer B C D
    updata_toml_start_bcd

    # 设置 token 地址
    InitTokenAddr
    docker cp "${chain33EthTokenAddr}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chain33EthTokenAddr}.abi
    docker cp "${chain33YccTokenAddr}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chain33YccTokenAddr}.abi
    docker cp "${chain33YccErc20Addr}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chain33YccErc20Addr}.abi
    docker cp "${ethBridgeToeknYccAddr}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${ethBridgeToeknYccAddr}.abi

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== chain33 lock BTY, eth burn BTY ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # 原来的地址金额
    result=$(${Chain33Cli} asset balance -a "${chain33DeployAddr}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "100.0000"

    # chain33 lock bty
    hash=$(${Chain33Cli} evm call -f 1 -a 5 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    # 原来的地址金额 减少了 5
    result=$(${Chain33Cli} asset balance -a "${chain33DeployAddr}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "95.0000"

    # chain33BridgeBank 是否增加了 5
    result=$(${Chain33Cli} asset balance -a "${chain33BridgeBank}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "5.0000"

    sleep 2
#    eth_block_wait 2 "${ethUrl}"

    # eth 这端 金额是否增加了 5
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    sleep 2
#    eth_block_wait 2 "${ethUrl}"

    # eth 这端 金额是否减少了 3
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    sleep ${maturityDegree}

    # 接收的地址金额 变成了 3
    result=$(${Chain33Cli} asset balance -a "${chain33ReceiverAddr}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "3.0000"

    # chain33BridgeBank 是否减少了 3
    result=$(${Chain33Cli} asset balance -a "${chain33BridgeBank}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "2.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 lock ZBC, eth burn ZBC
function TestChain33ToEthZBCAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== chain33 lock ZBC, eth burn ZBC ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # 原来的地址金额
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "0"

    # chain33 lock ZBC
    hash=$(${Chain33Cli} evm call -f 1 -c "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33ZBCErc20Addr}, 900000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    # chain33BridgeBank 是否增加了 9
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "900000000"

    sleep 2
#    eth_block_wait 2

    # eth 这端 金额是否增加了 9
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "9"

    # eth burn
    result=$(${CLIA} ethereum burn -m 8 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethBridgeToeknZBCAddr}" ) #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    sleep 2
#    eth_block_wait 2

    # eth 这端 金额是否减少了 1
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "1"

    sleep ${maturityDegree}

    # 接收的地址金额 变成了 8
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33ReceiverAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "800000000"

    # chain33BridgeBank 是否减少了 1
    result=$(${Chain33Cli} evm abi call -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "100000000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33Assets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 11个
    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    sleep 2
#    eth_block_wait 2 "${ethUrl}"

    # 查询 ETH 这端 bridgeBank 地址 11
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "11"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11 * le8
    is_equal "${result}" "1100000000"

    # 原来的数额
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33EthTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 11-5 * le8
    is_equal "${result}" "600000000"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" )
    cli_ret "${result}" "balance" ".balance" "6"

    # 比之前多 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "105"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ycc 资产,然后在 chain33 上 burn ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 7个 YCC
    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

     # eth 等待 10 个区块
    sleep 2
#    eth_block_wait 2 "${ethUrl}"

    # 查询 ETH 这端 bridgeBank 地址 7 YCC
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    sleep ${maturityDegree}

    # chain33 chain33EthTokenAddr（ETH合约中）查询 lock 金额
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7 * le8
    is_equal "${result}" "700000000"

    # 原来的数额 0
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "0"

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33YccTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm abi call -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 7-5 * le8
    is_equal "${result}" "200000000"

    # 查询 ETH 这端 bridgeBank 地址 2
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    # 更新后的金额 5
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function initPara() {
    # para add
    hash=$(${Para8901Cli}  send coins transfer -a 10000 -n test -t "${chain33ReceiverAddr}" -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944)
    check_tx "${Para8901Cli}" "${hash}"

    Chain33Cli=${Para8901Cli}
    InitChain33Validator

    # 先把 bty 转入到 paracross 合约中
    hash=$(${MainCli} send coins send_exec -e paracross -a 1000 -k "${chain33DeployKey}")
    check_tx "${MainCli}" "${hash}"

    # 主链中的 bty 夸链到 平行链中
    hash=$(${Para8801Cli} send para cross_transfer -a 200 -e coins -s bty -t "${chain33DeployAddr}" -k "${chain33DeployKey}")
    check_tx "${Para8801Cli}" "${hash}"
    check_tx "${Para8901Cli}" "${hash}"
    result=$(${Para8901Cli} asset balance -a "${chain33DeployAddr}" --asset_exec paracross --asset_symbol coins.bty | jq -r .balance)
    is_equal "${result}" "200.0000"

    # 把平行链中的 bty 转入 平行链中的 evm 合约
    hash=$(${Para8901Cli} send para transfer_exec -a 100 -e user.p.para.evm -s coins.bty -k "${chain33DeployKey}")
    check_tx "${Para8901Cli}" "${hash}"
    result=$(${Para8901Cli} asset balance -a "${chain33DeployAddr}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "100.0000"
}

function AllRelayerMainTest() {
    set +e
    docker_chain33_ip=$(get_docker_addr "${dockerNamePrefix}_chain33_1")
    MainCli="./chain33-cli --rpc_laddr http://${docker_chain33_ip}:8801"
    Para8801Cli="./chain33-cli --rpc_laddr http://${docker_chain33_ip}:8901 --paraName user.p.para."
    Para8901Cli="./chain33-cli --rpc_laddr http://${docker_chain33_ip}:8901 --paraName user.p.para."

    # shellcheck disable=SC2034
    {
        CLIA="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A"
        CLIB="docker exec ${dockerNamePrefix}_ebrelayerb_1 /root/ebcli_A"
        CLIC="docker exec ${dockerNamePrefix}_ebrelayerc_1 /root/ebcli_A"
        CLID="docker exec ${dockerNamePrefix}_ebrelayerd_1 /root/ebcli_A"
    }

    echo "${CLIA}"

    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    # shellcheck disable=SC2120
    if [[ $# -ge 2 ]]; then
        chain33ID="${2}"
    fi

    # init
    Chain33Cli=${MainCli}
    InitChain33Validator
    # para add
    initPara

    Chain33Cli=${Para8901Cli}
    StartDockerRelayerDeploy

    # test
    TestChain33ToEthAssets
    TestChain33ToEthZBCAssets
    TestETH2Chain33Assets
    TestETH2Chain33Ycc

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}


