package common

import (
	"encoding/json"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"testing"
)

func TestBytesToFloat64(t *testing.T) {
	var validatorMaps []types.MsgValidator
	validatorMaps = append(validatorMaps, types.MsgValidator{
		Address: "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
		Power:   7,
	})
	validatorMaps = append(validatorMaps, types.MsgValidator{
		Address: "1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi",
		Power:   6,
	})
	validatorMaps = append(validatorMaps, types.MsgValidator{
		Address: "3333333333333333333333333",
		Power:   8,
	})
	RemoveAddrFromValidatorMap(validatorMaps, 1)
	vv, _ := json.Marshal(validatorMaps)
	println(string(vv))
}

func RemoveAddrFromValidatorMap(validatorMap []types.MsgValidator, index int) []types.MsgValidator {
	return append(validatorMap[:index], validatorMap[index+1:]...)
}
