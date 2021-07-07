#在chain33部署操作手册
##步骤一: 离线创建7笔部署跨链合约的交易
```
交易1: 部署合约: Valset
交易2: 部署合约: EthereumBridge
交易3: 部署合约: Oracle
交易4: 部署合约: BridgeBank
交易5: 在合约EthereumBridge中设置BridgeBank合约地址
交易6: 在合约EthereumBridge中设置Oracle合约地址
交易7: 部署合约: BridgeRegistry

./boss4x chain33 offline create -f 1 -k 0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944 -n "deploy crossx to chain33" -r "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt, [12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv, 13qgxAxLeLYefgCUrbNRE36QUXB8SMkAVS, 1LfjSMbWAgjGp2FVgSeJVg4B1Y81r23wBu, 155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6], [25, 25, 25, 25]" --chainID 33

-f, --fee float: 交易费设置，因为只是少量几笔交易，且部署交易消耗gas较多，直接设置1个代币即可
-k, --key string: 部署人的私钥，用于对交易签名
-n, --note string: 备注信息 
-r, --valset string: 构造函数参数,严格按照该格式输入'addr, [addr, addr, addr, addr], [25, 25, 25, 25]',其中第一个地址为部署人私钥对应地址，后面4个地址为不同验证人的地址，4个数字为不同验证人的权重
--chainID 平行链的chainID
执行之后会将7笔交易写入到文件：deployCrossX2Chain33.txt

```

##步骤二: 串行发送7笔部署跨链合约的交易
```
./boss4x chain33 offline send
```