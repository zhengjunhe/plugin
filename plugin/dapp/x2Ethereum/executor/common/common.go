package common

import (
	"encoding/binary"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"math"
)

var (
	//日志
	clog = log.New("module", "common")
)

func Float64ToBytes(float float64) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, math.Float64bits(float))
	return result
}

func BytesToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(bytes))
}

func AddressIsEmpty(address string) bool {
	if address == "" {
		return true
	}

	var aa2 string
	return address == aa2
}

func AddToStringMap(in *types.StringMap, validator string) *types.StringMap {
	inStringMap := append(in.GetValidators(), validator)
	stringMapRes := new(types.StringMap)
	stringMapRes.Validators = inStringMap
	return stringMapRes
}
