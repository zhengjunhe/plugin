##
##编译solidity，并产生bin文件，abi文件，和相应的go文件


SRC_CONTRACT := sol_contracts
#SRC_CONTRACT := github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/sol_contracts
GO_OUT := generated
PACKAGE := generated

proj := "build"
.PHONY: default build clean registry bridgeBank

default: build

build:
	@abigen --sol $(SRC_CONTRACT)/BridgeRegistry.sol --pkg $(PACKAGE) --out $(GO_OUT)/BridgeRegistry.go
	@abigen --sol $(SRC_CONTRACT)/BridgeBank/BridgeBank.sol --pkg $(PACKAGE) --out $(GO_OUT)/BridgeBank.go
	#@abigen --sol $(SRC_CONTRACT)/CosmosBridge.sol --pkg $(PACKAGE) --out $(GO_OUT)/CosmosBridge.go
	#@abigen --sol $(SRC_CONTRACT)/Oracle.sol --pkg $(PACKAGE) --out $(GO_OUT)/Oracle.go
	#@abigen --sol $(SRC_CONTRACT)/Valset.sol --pkg $(PACKAGE) --out $(GO_OUT)/Valset.go
clean:
	@rm -fr $(GO_OUT)/*

registry:
	@abigen --sol $(SRC_CONTRACT)/BridgeRegistry.sol --pkg $(PACKAGE) --out $(GO_OUT)/BridgeRegistry.go

bridgeBank:
	@abigen --sol $(SRC_CONTRACT)/BridgeBank/BridgeBank.sol --pkg $(PACKAGE) --out $(GO_OUT)/BridgeBank.go
