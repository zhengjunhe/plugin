package para

import (
	"github.com/33cn/dplatform/queue"
	drivers "github.com/33cn/dplatform/system/mempool"
	"github.com/33cn/dplatform/types"
)

//--------------------------------------------------------------------------------
// Module Mempool

func init() {
	drivers.Reg("para", New)
}

//New 创建price cache 结构的 mempool
func New(cfg *types.Mempool, sub []byte) queue.Module {
	return NewMempool(cfg)
}
