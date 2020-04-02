package ethtxs

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"math/big"

	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// GenerateClaimHash : Generates an OracleClaim hash from a ProphecyClaim's event data
func GenerateClaimHash(prophecyID []byte, sender []byte, recipient []byte, token []byte, amount []byte, validator []byte) common.Hash {
	// Generate a hash containing the information
	rawHash := crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)

	// Cast hash to hex encoded string
	return rawHash
}

func SignClaim4Eth(hash common.Hash, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	rawSignature, _ := prefixMessage(hash, privateKey)
	signature := hexutil.Bytes(rawSignature)
	return signature, nil
}

func prefixMessage(message common.Hash, key *ecdsa.PrivateKey) ([]byte, []byte) {
	prefixed := solsha3.SoliditySHA3WithPrefix(message[:])
	sig, err := secp256k1.Sign(prefixed, math.PaddedBigBytes(key.D, 32))
	if err != nil {
		panic(err)
	}

	return sig, prefixed
}

func loadPrivateKey(privateKey []byte) (key *ecdsa.PrivateKey, err error) {
	key, err = crypto.ToECDSA(privateKey)
	if nil != err {
		return nil, err
	}
	return
}

// LoadSender : uses the validator's private key to load the validator's address
func LoadSender(privateKey *ecdsa.PrivateKey) (address common.Address, err error) {
	// Parse public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, ebrelayerTypes.ErrPublicKeyType
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return fromAddress, nil
}

func PrepareAuth(backend bind.ContractBackend, privateKey *ecdsa.PrivateKey, transactor common.Address) (*bind.TransactOpts, error) {
	if nil == privateKey || nil == backend {
		txslog.Error("PrepareAuth", "nil input parameter", "backend", backend, "privateKey", privateKey)
		return nil, errors.New("nil input parameter")
	}

	ctx := context.Background()
	nonce, err := backend.PendingNonceAt(context.Background(), transactor)
	if err != nil {
		txslog.Error("PrepareAuth", "Failed to PendingNonceAt due to:", err.Error())
		return nil, errors.New("Failed to get nonce")
	}

	gasPrice, err := backend.SuggestGasPrice(ctx)
	if err != nil {
		txslog.Error("PrepareAuth", "Failed to SuggestGasPrice due to:", err.Error())
		return nil, errors.New("Failed to get suggest gas price")
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = GasLimit4Deploy
	auth.GasPrice = gasPrice

	return auth, nil
}
