package common

import (
	"fmt"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"testing"
)

func TestEthAddress_String(t *testing.T) {
	moduleAddress := dapp.ExecAddress(types.ModuleName)
	fmt.Println(moduleAddress)
}
