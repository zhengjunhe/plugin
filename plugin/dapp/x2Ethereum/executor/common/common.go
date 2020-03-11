package common

import (
	"encoding/binary"
	log "github.com/33cn/chain33/common/log/log15"
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
