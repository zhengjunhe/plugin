package common

import (
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"math"
	"math/big"
	"testing"
)

func TestAddToStringMap(t *testing.T) {
	bn := big.NewInt(1)
	ss := types.TrimZeroAndDot(fmt.Sprintf("%.0f", types.MultiplySpecifyTimes(math.Trunc(5*1e4), 14)))
	bn, ok := bn.SetString(ss, 10)
	fmt.Println(bn, ok)
}
