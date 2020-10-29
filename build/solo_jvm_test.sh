#!/usr/bin/env bash

cli="./chain33-cli"
contract=Guess
#初始化
$cli seed save -p hzjhzj123 -s "hurdle civil burden caught lamp spoon confirm admit plug gate bracket paddle eight merry repair"
$cli wallet unlock -p hzjhzj123

#创建账户
$cli account import_key -l opener -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944
$cli account import_key -l player -k 0x7b2800cdecd978ab0e877f7e3734b9d0b11d864fa51d9b623d7bdbd76c16a40d



echo "Begin to test contract $contract"
#部署合约
$cli send jvm create -x $contract -d . -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt

#转账
echo "send coins transfer -a 1000"
$cli send coins transfer -a 1000 -n "t1000" -t 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
echo "send coins send_exec -a 30"
$cli send coins send_exec -a 30 -e user.jvm.$contract -n send2exec -k 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU

#开始游戏
echo "send jvm call -e $contract -x startGame"
$cli send jvm call -e $contract -x startGame -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt

#投注
echo "send jvm call -e $contract -x playGame"
$cli send jvm call -e $contract -x playGame -r "6 2" -k 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU



#转账，增加高度
for ((a=0;a<=12;a++))
do
echo "$a" tx to transfer
$cli send  coins transfer -a 1 -n "t1000" -t 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU -k 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
done

##开奖
echo "close $contract"
$cli send jvm call -e $contract -x closeGame -k 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU
#
##查看信息
$cli jvm query -e user.jvm.Guess -r "getGuessRecordByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1"
$cli jvm query -e user.jvm.Guess -r "getBonusByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1"
$cli jvm query -e user.jvm.Guess -r "getLuckNumByRound 1PrTWtT1Bzhg2L8jjVKU7ohxHVXLU4NMEU 1"
