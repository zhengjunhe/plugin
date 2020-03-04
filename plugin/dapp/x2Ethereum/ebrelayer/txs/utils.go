package txs

import (
	"crypto/ecdsa"
	solsha3 "github.com/miguelmota/go-solidity-sha3"

    ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// GenerateClaimHash : Generates an OracleClaim hash from a ProphecyClaim's event data
func GenerateClaimHash(prophecyID []byte, sender []byte, recipient []byte, token []byte, amount []byte, validator []byte) string {
	// Generate a hash containing the information
	rawHash := crypto.Keccak256Hash(prophecyID, sender, recipient, token, amount, validator)

	// Cast hash to hex encoded string
	return rawHash.Hex()
}

func SignClaim4Eth(hash string, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	rawSignature, _ := prefixMessage(hash, privateKey)
	signature := hexutil.Bytes(rawSignature)
	return signature, nil
}

func prefixMessage(message string, key *ecdsa.PrivateKey) ([]byte, []byte) {
	// Turn the message into a 32-byte hash
	hash := solsha3.SoliditySHA3(solsha3.String(message))
	// Prefix and then hash to mimic behavior of eth_sign
	prefixed := solsha3.SoliditySHA3(solsha3.String("\x19Ethereum Signed Message:\n32"), solsha3.Bytes32(hash))
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
func LoadSender(privateKey *ecdsa.PrivateKey) (address *common.Address, err error) {
	// Parse public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, ebrelayerTypes.ErrPublicKeyType
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return &fromAddress, nil
}
