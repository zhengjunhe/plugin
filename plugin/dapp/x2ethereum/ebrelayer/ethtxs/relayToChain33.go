package ethtxs

// ------------------------------------------------------------
//	Relay : Builds and encodes EthBridgeClaim Msgs with the
//  	specified variables, before presenting the unsigned
//      transaction to validators for optional signing.
//      Once signed, the data packets are sent as transactions
//      on the dplatform Bridge.
// ------------------------------------------------------------

import (
	"github.com/33cn/dplatform/common"
	dplatformCrypto "github.com/33cn/dplatform/common/crypto"
	"github.com/33cn/dplatform/rpc/jsonclient"
	rpctypes "github.com/33cn/dplatform/rpc/types"
	dplatformTypes "github.com/33cn/dplatform/types"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

// RelayLockToDplatform : RelayLockToDplatform applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockToDplatform(privateKey dplatformCrypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL string) (string, error) {
	var res string

	params := &types.Eth2Dplatform{
		EthereumChainID:       claim.EthereumChainID,
		BridgeContractAddress: claim.BridgeBrankAddr,
		Nonce:                 claim.Nonce,
		IssuerDotSymbol:       claim.Symbol,
		TokenContractAddress:  claim.TokenAddr,
		EthereumSender:        claim.EthereumSender,
		DplatformReceiver:       claim.DplatformReceiver,
		Amount:                claim.Amount,
		ClaimType:             int64(claim.ClaimType),
		Decimals:              claim.Decimal,
	}

	pm := rpctypes.CreateTxIn{
		Execer:     X2Eth,
		ActionName: types.NameEth2DplatformAction,
		Payload:    dplatformTypes.MustPBToJSON(params),
	}
	ctx := jsonclient.NewRPCCtx(rpcURL, "Dplatform.CreateTransaction", pm, &res)
	_, _ = ctx.RunResult()

	data, err := common.FromHex(res)
	if err != nil {
		return "", err
	}
	var tx dplatformTypes.Transaction
	err = dplatformTypes.Decode(data, &tx)
	if err != nil {
		return "", err
	}

	if tx.Fee == 0 {
		tx.Fee, err = tx.GetRealFee(1e5)
		if err != nil {
			return "", err
		}
	}
	//构建交易，验证人validator用来向dplatform合约证明自己验证了该笔从以太坊向dplatform跨链转账的交易
	tx.Sign(dplatformTypes.SECP256K1, privateKey)

	txData := dplatformTypes.Encode(&tx)
	dataStr := common.ToHex(txData)
	pms := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}
	var txhash string

	ctx = jsonclient.NewRPCCtx(rpcURL, "Dplatform.SendTransaction", pms, &txhash)
	_, err = ctx.RunResult()
	return txhash, err
}

//RelayBurnToDplatform ...
func RelayBurnToDplatform(privateKey dplatformCrypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL string) (string, error) {
	var res string

	params := &types.Eth2Dplatform{
		EthereumChainID:       claim.EthereumChainID,
		BridgeContractAddress: claim.BridgeBrankAddr,
		Nonce:                 claim.Nonce,
		IssuerDotSymbol:       claim.Symbol,
		TokenContractAddress:  claim.TokenAddr,
		EthereumSender:        claim.EthereumSender,
		DplatformReceiver:       claim.DplatformReceiver,
		Amount:                claim.Amount,
		ClaimType:             int64(claim.ClaimType),
		Decimals:              claim.Decimal,
	}

	pm := rpctypes.CreateTxIn{
		Execer:     X2Eth,
		ActionName: types.NameWithdrawEthAction,
		Payload:    dplatformTypes.MustPBToJSON(params),
	}
	ctx := jsonclient.NewRPCCtx(rpcURL, "Dplatform.CreateTransaction", pm, &res)
	_, _ = ctx.RunResult()

	data, err := common.FromHex(res)
	if err != nil {
		return "", err
	}
	var tx dplatformTypes.Transaction
	err = dplatformTypes.Decode(data, &tx)
	if err != nil {
		return "", err
	}

	if tx.Fee == 0 {
		tx.Fee, err = tx.GetRealFee(1e5)
		if err != nil {
			return "", err
		}
	}
	//构建交易，验证人validator用来向dplatform合约证明自己验证了该笔从以太坊向dplatform跨链转账的交易
	tx.Sign(dplatformTypes.SECP256K1, privateKey)

	txData := dplatformTypes.Encode(&tx)
	dataStr := common.ToHex(txData)
	pms := rpctypes.RawParm{
		Token: "BTY",
		Data:  dataStr,
	}
	var txhash string

	ctx = jsonclient.NewRPCCtx(rpcURL, "Dplatform.SendTransaction", pms, &txhash)
	_, err = ctx.RunResult()
	return txhash, err
}
