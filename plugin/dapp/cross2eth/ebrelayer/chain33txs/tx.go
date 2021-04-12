package chain33txs

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	chain33Crypto "github.com/33cn/chain33/common/crypto"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	types "github.com/33cn/chain33/types"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/utils"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/protobuf/proto"
)

var chain33txLog = log.New("module", "chain33_txs")

// RelayLockToChain33 : RelayLockToChain33 applies validator's signature to an EthBridgeClaim message
//		containing information about an event on the Ethereum blockchain before relaying to the Bridge
func RelayLockBurnToChain33(privateKey chain33Crypto.PrivKey, privateKey_ecdsa *ecdsa.PrivateKey, claim *ebrelayerTypes.EthBridgeClaim, rpcURL string) (string, error) {
	nonceBytes := big.NewInt(claim.Nonce).Bytes()
	amountBytes := big.NewInt(claim.Amount).Bytes()
	claimID := crypto.Keccak256Hash(nonceBytes, []byte(claim.EthereumSender), []byte(claim.Chain33Receiver), []byte(claim.Symbol), amountBytes)

	// Sign the hash using the active validator's private key
	signature, err := utils.SignClaim4Evm(claimID, privateKey_ecdsa)
	if nil != err {
		return "", err
	}
	parameter := fmt.Sprintf("newOracleClaim(%d, %s, %s, %s, %s, %s, %s, %s)",
		claim.ClaimType,
		claim.EthereumSender,
		claim.Chain33Receiver,
		claim.TokenAddr,
		claim.Symbol,
		claim.Amount,
		claimID,
		signature)

	return RelayEvmTx2Chain33(privateKey, claim, parameter, rpcURL)
}

func createEvmTx(privateKey chain33Crypto.PrivKey, action proto.Message, execer, to string, fee int64) string {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: fee, To: to}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	tx.ChainID = 33

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)
	return dataStr
}

func RelayEvmTx2Chain33(privateKey chain33Crypto.PrivKey, claim *ebrelayerTypes.EthBridgeClaim, parameter, rpcURL string) (string, error) {
	note := fmt.Sprintf("RelayLockToChain33 by validator:%s with nonce:%d",
		address.PubKeyToAddr(privateKey.PubKey().Bytes()),
		claim.Nonce)

	action := evmtypes.EVMContractAction{Amount: 0, GasLimit: 0, GasPrice: 0, Note: note, Abi: parameter}

	feeInt64 := int64(1e7)
	toAddr := ""
	wholeEvm := claim.ChainName + ".evm"
	//name表示发给哪个执行器
	data := createEvmTx(privateKey, &action, wholeEvm, toAddr, feeInt64)
	params := rpctypes.RawParm{
		Token: "BTY",
		Data:  data,
	}
	var txhash string

	ctx := jsonclient.NewRPCCtx(rpcURL, "Chain33.SendTransaction", params, &txhash)
	_, err := ctx.RunResult()
	return txhash, err
}
