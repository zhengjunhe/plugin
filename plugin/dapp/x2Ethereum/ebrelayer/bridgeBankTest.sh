#!/usr/bin/env bash
set -x
CLI="../build/ebcli_A"
prophecyID=1

InitAndDeploy() {
  result=$(${CLI} relayer set_pwd -n 123456hzj -o kk | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to set_pwd"
        exit 1
  fi
  result=$(${CLI} relayer unlock -p 123456hzj | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to unlock"
        exit 1
  fi
  result=$(${CLI} relayer ethereum deploy | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to deploy"
        exit 1
  fi
  result=$(${CLI} relayer ethereum import_chain33privatekey -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944 | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to import_chain33privatekey"
        exit 1
  fi
  result=$(${CLI} relayer ethereum import_ethprivatekey -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230 | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to import_ethprivatekey"
        exit 1
  fi
  echo "Succeed to InitAndDeploy"
}
#chain33 asset ---> Ethereum
#测试在chain33上锁定资产,然后在以太坊上铸币
#发行token="BTY"
#NewProphecyClaim lock
#铸币NewOracleClaim,
#ProcessBridgeProphecy
#Bridge token minting (for locked chain33 assets)
TestBrigeTokenMint4Chain33Assets() {
  #创建token，获取token地址, 该处创建的token为bridge token，即原始发行方为chain33
  result=$(${CLI} relayer ethereum token -s abc | jq .msg)
  #Token address:xxxx
  tokenAddr=${result#*:}
  tokenAddr=${tokenAddr%\"*}

  result=$(${CLI} relayer ethereum prophecy -a ${tokenAddr} -t 0x0c05ba5c230fdaa503b53702af1962e08d0c60bf -c 2 -s bty | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to ${CLI} relayer ethereum prophecy"
        exit 1
  fi

  echo "wait relayer 5 seconds to make sure validator vote for the prophecy"
  sleep 5

  #确认
  checkProphecyIDActive ${prophecyID}

  #后期改为自动处理，不需要手动递交
  result=$(${CLI} relayer ethereum process -i ${prophecyID} | jq .isOK)
  if [[ ${result} != "true" ]]; then
        echo "failed to ${CLI} relayer ethereum process -i 1"
        exit 1
  fi

  #balance:297
  result=$(${CLI} relayer ethereum balance -o 0x0df9a824699bc5878232c9e612fe1a5346a5a368 -a ${tokenAddr} | jq .msg)
  balance=${result#*:}
  balance=${balance%\"*}
  if [[ ${balance} != "99" ]]; then
        echo "The balance is not correct"
        exit 1
  fi
  echo "Succeed to TestBrigeTokenMint4Chain33Assets"
}

checkProphecyIDActive() {
    while true; do
        pending=$(${CLI} relayer ethereum ispending -i ${1} | jq .isOK)
        if [[ ${pending} == "true" ]]; then
            break
        fi
        sleep 1
    done
}


main () {
    #InitAndDeploy
    TestBrigeTokenMint4Chain33Assets
}

main