package ethtxs

import (
	"fmt"
	"testing"

	"github.com/33cn/dplatform/client/mocks"
	dplatformCommon "github.com/33cn/dplatform/common"
	_ "github.com/33cn/dplatform/system"
	"github.com/33cn/dplatform/system/crypto/secp256k1"
	dplatformTypes "github.com/33cn/dplatform/types"
	"github.com/33cn/dplatform/util/testnode"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	chainTestCfg = dplatformTypes.NewDplatformConfig(dplatformTypes.GetDefaultCfgstring())
)

func Test_RelayToDplatform(t *testing.T) {
	var tx dplatformTypes.Transaction
	var ret dplatformTypes.Reply
	ret.IsOk = true

	mockapi := &mocks.QueueProtocolAPI{}
	// 这里对需要mock的方法打桩,Close是必须的，其它方法根据需要
	mockapi.On("Close").Return()
	mockapi.On("AddPushSubscribe", mock.Anything).Return(&ret, nil)
	mockapi.On("CreateTransaction", mock.Anything).Return(&tx, nil)
	mockapi.On("SendTx", mock.Anything).Return(&ret, nil)
	mockapi.On("SendTransaction", mock.Anything).Return(&ret, nil)
	mockapi.On("GetConfig", mock.Anything).Return(chainTestCfg, nil)

	mock33 := testnode.New("", mockapi)
	defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	// 这里必须设置监听端口，默认的是无效值
	rpcCfg.JrpcBindAddr = "127.0.0.1:8801"
	mock33.GetRPC().Listen()

	dplatformPrivateKeyStr := "0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
	var driver secp256k1.Driver
	privateKeySli, err := dplatformCommon.FromHex(dplatformPrivateKeyStr)
	require.Nil(t, err)

	priKey, err := driver.PrivKeyFromBytes(privateKeySli)
	require.Nil(t, err)

	claim := &ebrelayerTypes.EthBridgeClaim{}

	fmt.Println("======================= testRelayLockToDplatform =======================")
	_, err = RelayLockToDplatform(priKey, claim, "http://127.0.0.1:8801")
	require.Nil(t, err)

	fmt.Println("======================= testRelayBurnToDplatform =======================")
	_, err = RelayBurnToDplatform(priKey, claim, "http://127.0.0.1:8801")
	require.Nil(t, err)
}
