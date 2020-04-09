package common

import (
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	num      = 10
	numBytes = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x40}
)

func TestFloat64ToBytes(t *testing.T) {
	nb := Float64ToBytes(float64(num))
	assert.Equal(t, numBytes, nb)
}

func TestBytesToFloat64(t *testing.T) {
	n := BytesToFloat64(numBytes)
	assert.Equal(t, num, n)
}

func TestAddressIsEmpty(t *testing.T) {
	flg := AddressIsEmpty("")
	assert.Equal(t, true, flg)
}

func TestAddToStringMap(t *testing.T) {
	in := new(types.StringMap)
	res := AddToStringMap(in, "validator")
	assert.Equal(t, &types.StringMap{
		Validators: []string{"validator"},
	}, res)
}
