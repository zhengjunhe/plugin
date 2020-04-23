package ethtxs

import (
	"crypto/ecdsa"
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	txslog = log15.New("ethereum relayer", "ethtxs")
)

const (
	// GasLimit : the gas limit in Gwei used for transactions sent with TransactOpts
	GasLimit        = uint64(600000)
	GasLimit4Deploy = uint64(0) //此处需要设置为0,让交易自行估计,否则将会导致部署失败,TODO:其他解决途径后续调研解决
)

// RelayProphecyClaimToEthereum : relays the provided ProphecyClaim to Chain33Bridge contract on the Ethereum network
func RelayProphecyClaimToEthereum(oracleInstance *generated.Oracle, client *ethclient.Client, sender common.Address, event events.Event, claim ProphecyClaim, privateKey *ecdsa.PrivateKey, chain33TxHash []byte) (txhash string, err error) {
	txslog.Info("RelayProphecyClaimToEthereum","sender", sender.String(), "event", event, "claim", claim)

	auth, err := PrepareAuth(client, privateKey, sender)
	if nil != err {
		txslog.Error("RelayProphecyClaimToEthereum", "PrepareAuth err", err.Error())
		return "", err
	}

	claimID := crypto.Keccak256Hash(chain33TxHash, claim.Chain33Sender, claim.EthereumReceiver.Bytes(), claim.TokenContractAddress.Bytes(), claim.Amount.Bytes())

	// Sign the hash using the active validator's private key
	signature, err := SignClaim4Eth(claimID, privateKey)
	if nil != err {
		return "", err
	}

	tx, err := oracleInstance.NewOracleClaim(auth, uint8(claim.ClaimType), claim.Chain33Sender, claim.EthereumReceiver, claim.TokenContractAddress, claim.Symbol, claim.Amount, claimID, signature)
	if nil != err {
		txslog.Error("RelayProphecyClaimToEthereum", "NewOracleClaim failed due to:", err.Error())
		return "", err
	}

	txhash = tx.Hash().Hex()
	txslog.Info("RelayProphecyClaimToEthereum", "NewProphecyClaim tx hash:", txhash)
	err = waitEthTxFinished(client, tx.Hash(), "ProphecyClaimToEthereum")
	if nil != err {
		return txhash, err
	}
	return "", nil
}

// RelayOracleClaimToEthereum : relays the provided OracleClaim to Oracle contract on the Ethereum network
func RelayOracleClaimToEthereum(provider string, sender, contractAddress common.Address, event events.Event, claim *OracleClaim, privateKey *ecdsa.PrivateKey) (txhash string, err error) {
	// Initialize client service, validator's tx auth, and target contract address
	//client, auth, target, err := initRelayConfig(provider, sender, contractAddress, event, privateKey)
	//if nil != err {
	//	return "", err
	//}
	//
	//// Initialize Oracle instance
	//oracleInstance, err := oracle.NewOracle(*target, client)
	//if err != nil {
	//	txslog.Error("RelayOracleClaimToEthereum", "NewOracle failed due to:", err.Error())
	//	return "", err
	//}
	//
	//// Send transaction
	//tx, err := oracleInstance.NewOracleClaim(auth, claim.ProphecyID, claim.Message, claim.Signature)
	//if err != nil {
	//	txslog.Error("RelayOracleClaimToEthereum", "NewOracleClaim failed due to:", err.Error())
	//	return "", err
	//}
	//txhash = tx.Hash().Hex()
	//txslog.Info("RelayOracleClaimToEthereum", "NewOracleClaim tx hash:", txhash)
	//err = waitEthTxFinished(client, tx.Hash(), "RelayOracleClaimToEthereum")
	//if nil != err {
	//	return txhash, err
	//}

	//return txhash, nil
	return "", nil
}