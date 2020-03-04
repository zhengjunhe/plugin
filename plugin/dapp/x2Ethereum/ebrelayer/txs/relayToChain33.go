package txs

// ------------------------------------------------------------
//	Relay : Builds and encodes EthBridgeClaim Msgs with the
//  	specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the chain33 Bridge.
// ------------------------------------------------------------

import (
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	chain33Types "github.com/33cn/chain33/types"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/chain33"
	"github.com/33cn/chain33/common"
	chain33Crypto "github.com/33cn/chain33/common/crypto"
)

// RelayLockToChain33 : RelayLockToChain33 applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockToChain33(
	privateKey chain33Crypto.PrivKey,
	chainID string,
	claim *ebrelayerTypes.EthBridgeClaim,
	rpcUrl string,
) (string, error) {

	tx := &chain33Types.Transaction{}
	//构建交易，验证人validator用来向chain33合约证明自己验证了该笔从以太坊向chain33跨链转账的交易
	tx.Execer = []byte(chainID + "." +chain33.X2Eth)
	tx.Sign(chain33Types.SECP256K1, privateKey)

	txData := chain33Types.Encode(tx)
	dataStr := common.ToHex(txData)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}
	var txhash string

	ctx := jsonclient.NewRPCCtx(rpcUrl, "Chain33.SendTransaction", params, &txhash)
	_, err := ctx.RunResult()
	return txhash, err
}
