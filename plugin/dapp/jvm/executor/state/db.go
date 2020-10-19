package state

import (
	"github.com/33cn/chain33/types"
	chain33db "github.com/33cn/chain33/common/db"
	"strings"
)

var (
	currentExecTxHash string
	localDB chain33db.DB
)

func newMemDB() chain33db.DB {
	memdb, err := chain33db.NewGoMemDB("", "", 0)
	if err != nil {
		panic(err)
	}
	return memdb
}

func setCurrentTx(txhashNew string) {
	if 0 == strings.Compare(txhashNew, currentExecTxHash) {
		return
	}
	currentExecTxHash = txhashNew
	localDB = newMemDB()
}

func getLocalValue(key []byte, txHash string) ([]byte, error) {
	setCurrentTx(txHash)
	return localDB.Get(key)
}

func setLocalValue(key, value []byte, txHash string) error {
	setCurrentTx(txHash)
	return localDB.Set(key, value)
}

func GetAllLocalKeyValues() []*types.KeyValue {
	goMemDB, ok := localDB.(*chain33db.GoMemDB)
	if !ok {
		return nil
	}

	var kvs []*types.KeyValue
	it := goMemDB.DB().NewIterator(nil)
	for it.Next() {
		kvs = append(kvs, &types.KeyValue{Key:it.Key(), Value:it.Value()})
	}
	return kvs
}


