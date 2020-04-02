package executor

import (
	"fmt"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util"
	types2 "github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	apimock "github.com/33cn/chain33/client/mocks"
)

var chainTestCfg = types.NewChain33Config(types.GetDefaultCfgstring())

func init() {
	Init(types2.X2ethereumX, chainTestCfg, nil)
}

//--------------------------合约管理员账户操作-------------------------//
// 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv "data": "0x4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"
// 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt "data": "0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

var bankAddr = "1BqP2vHkYNjSgdnTqm7pGbnphLhtEhuJFi"

var add = "1Lu8XmStYvWwfNqiQ3nNK34R9FfH4kRpBV" // 合约地址
//fromaddr 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
var privFrom = getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944")

//to 1Mcx9PczwPQ79tDnYzw62SEQifPwXH84yN
var privTo = getprivkey("BC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944")

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func TestX2ethereum_Exec_AddValidator(t *testing.T) {
	add := types2.MsgValidator{
		Address: "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
		Power:   7,
	}

	x := x2ethereum{
		drivers.DriverBase{},
	}
	_, sdb, kvdb := util.CreateTestDB()
	x.SetLocalDB(kvdb)
	api := new(apimock.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(chainTestCfg, nil)
	//x.SetAPI(api)
	x.SetStateDB(sdb)

	tx := &types.Transaction{}
	tx.Execer = []byte(types2.X2ethereumX)
	tx.To = address.ExecAddress(types2.X2ethereumX)
	tx.Nonce = 1
	tx.Sign(types.SECP256K1, privFrom)

	action1 := newAction(&x, tx, 0)
	fmt.Println("***", &action1.keeper)
	//_, _ = action1.procMsgSetConsensusThreshold(&types2.MsgConsensusThreshold{ConsensusThreshold: 80})
	//msg, err := x.Query_GetConsensusThreshold(&types2.QueryConsensusThresholdParams{})
	//fmt.Println("=", msg, err)
	//assert.NotEqual(t, action1, nil)

	receipt, err := action1.procAddValidator(&add)
	fmt.Println(err, receipt)
	assert.NoError(t, err)

	//tx.Nonce = 2
	//action1 = newAction(&x, tx, 1)

	add2 := types2.MsgValidator{
		Address: "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
		Power:   8,
	}
	receipt, err = action1.procAddValidator(&add2)
	fmt.Println(err, receipt)
	//	fmt.Println("0000", err)
	assert.NoError(t, err)

	//	add3 := types2.QueryValidatorsParams{
	//		Validator: "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
	//	}

	//	message, err := x.Query_GetValidators(&add3)
	//	fmt.Println(message, err)

	//	fmt.Println(receipt, err)

}
