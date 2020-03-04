package utils

// --------------------------------------------------------
//      Utils
//
//      Utils contains utility functionality for the ebrelayer.
// --------------------------------------------------------

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"unicode"

	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

// IsZeroAddress : checks an Ethereum address and returns a bool which indicates if it is the null address
func IsZeroAddress(address common.Address) bool {
	return address == common.HexToAddress(nullAddress)
}

//密码合法性校验,密码长度在8-30位之间。必须是数字+字母的组合
func IsValidPassWord(password string) bool {
	pwLen := len(password)
	if pwLen < 8 || pwLen > 30 {
		return false
	}

	var char bool
	var digit bool
	for _, s := range password {
		if unicode.IsLetter(s) {
			char = true
		} else if unicode.IsDigit(s) {
			digit = true
		} else {
			return false
		}
	}
	return char && digit
}

func decodeInt64(int64bytes []byte) (int64, error) {
	var value types.Int64
	err := types.Decode(int64bytes, &value)
	if err != nil {
		//may be old database format json...
		err = json.Unmarshal(int64bytes, &value.Data)
		if err != nil {
			return -1, types.ErrUnmarshal
		}
	}
	return value.Data, nil
}

func LoadInt64FromDB(key []byte, db dbm.DB) (int64, error) {
	bytes, err := db.Get(key)
	if bytes == nil || err != nil {
		//if err != dbm.ErrNotFoundInDb {
		//	log.Error("LoadInt64FromDB", "error", err)
		//}
		return -1, types.ErrHeightNotExist
	}
	return decodeInt64(bytes)
}

func QueryTxhashes(prefix []byte, db dbm.DB) []string {
	kvdb := dbm.NewKVDB(db)
	hashes, err := kvdb.List([]byte(prefix), nil, 10, 1)
	if nil != err {
		return nil
	}

	var hashStrs []string
	for _, hash := range hashes {
		hashStrs = append(hashStrs, string(hash))
	}
	return hashStrs
}

