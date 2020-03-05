package txs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      chain33 Bridge.
// --------------------------------------------------------

import (
	"crypto/ecdsa"
	"math/big"
	"regexp"
	"strings"

	chain33Types "github.com/33cn/chain33/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
)

// LogLockToEthBridgeClaim : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func LogLockToEthBridgeClaim(valAddr []byte, event *events.LockEvent) (ebrelayerTypes.EthBridgeClaim, error) {
	witnessClaim := ebrelayerTypes.EthBridgeClaim{}

	// Nonce type casting (*big.Int -> int)
	nonce := event.Nonce.Int64()

	// Sender type casting (address.common -> string)
	//sender := ebrelayerTypes.NewEthereumAddress(event.From.Hex())

	recipient := event.To
	if 0 == len(recipient) {
		return witnessClaim, ebrelayerTypes.ErrEmptyAddress
	}

	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && event.Token != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return witnessClaim, ebrelayerTypes.ErrAddress4Eth
	}

	// Package the information in a unique EthBridgeClaim
	witnessClaim.Nonce = nonce
	witnessClaim.EthereumSender = event.From.Bytes()
	witnessClaim.ValidatorAddress = valAddr
	witnessClaim.Chain33Receiver = recipient
	witnessClaim.Amount = event.Value.Int64()

	return witnessClaim, nil
}

// BurnLockTxReceiptToChain33Msg : parses data from a Burn/Lock event witnessed on chain33 into a CosmosMsg struct
func BurnLockTxReceiptToChain33Msg(claimType events.Event, receipt *chain33Types.ReceiptData) events.Chain33Msg {
	// Set up variables
	var cosmosSender []byte
	var ethereumReceiver, tokenContractAddress common.Address
	var symbol string
	var amount *big.Int

	// Iterate over attributes
	for _, log := range receipt.Logs {
		// Get (key, value) for each attribute

		switch log.Ty {
		case 1:
			txslog.Debug("BurnLockTxReceiptToChain33Msg", "value", string(log.Log))
		default:

		}

		// Set variable based on value of CosmosMsgAttributeKey
		//switch key {
		//case events.CosmosSender.String():
		//	// Parse sender's Cosmos address
		//	cosmosSender = []byte(val)
		//case events.EthereumReceiver.String():
		//	// Confirm recipient is valid Ethereum address
		//	if !common.IsHexAddress(val) {
		//		log.Fatal("Invalid recipient address:", val)
		//	}
		//	// Parse recipient's Ethereum address
		//	ethereumReceiver = common.HexToAddress(val)
		//case events.Coin.String():
		//	// Parse symbol and amount from coin string
		//	symbol, amount = getSymbolAmountFromCoin(val)
		//case events.TokenContractAddress.String():
		//	// Confirm token contract address is valid Ethereum address
		//	if !common.IsHexAddress(val) {
		//		log.Fatal("Invalid token address:", val)
		//	}
		//	// Parse token contract address
		//	tokenContractAddress = common.HexToAddress(val)
		//}
	}

	// Package the event data into a CosmosMsg
	return events.NewChain33Msg(claimType, cosmosSender, ethereumReceiver, symbol, amount, tokenContractAddress)
}

// ProphecyClaimToSignedOracleClaim : packages and signs a prophecy claim's data, returning a new oracle claim
func ProphecyClaimToSignedOracleClaim(event events.NewProphecyClaimEvent, privateKey *ecdsa.PrivateKey) (*OracleClaim, error) {
	// Parse relevant data into type byte[]
	prophecyID := event.ProphecyID.Bytes()
	sender := event.Chain33Sender
	recipient := []byte(event.EthereumReceiver.Hex())
	token := []byte(event.TokenAddress.Hex())
	amount := event.Amount.Bytes()
	validator := []byte(event.ValidatorAddress.Hex())

	// Generate rawHash using ProphecyClaim data
	hash := GenerateClaimHash(prophecyID, sender, recipient, token, amount, validator)

	// Sign the hash using the active validator's private key
	//signature := SignClaim(hash)
	signature, err := SignClaim4Eth(hash, privateKey)
	if nil != err {
		return nil, err
	}
	// Package the ProphecyID, Message, and Signature into an OracleClaim
	oracleClaim := &OracleClaim{
		ProphecyID: event.ProphecyID,
		Message:    hash,
		Signature:  signature,
	}

	return oracleClaim, nil
}

// Chain33MsgToProphecyClaim : parses event data from a CosmosMsg, packaging it as a ProphecyClaim
func Chain33MsgToProphecyClaim(event events.Chain33Msg) ProphecyClaim {
	claimType := event.ClaimType
	cosmosSender := event.Chain33Sender
	ethereumReceiver := event.EthereumReceiver
	tokenContractAddress := event.TokenContractAddress
	symbol := strings.ToLower(event.Symbol)
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		CosmosSender:         cosmosSender,
		EthereumReceiver:     ethereumReceiver,
		TokenContractAddress: tokenContractAddress,
		Symbol:               symbol,
		Amount:               amount,
	}

	return prophecyClaim
}

// getSymbolAmountFromCoin : Parse (symbol, amount) from coin string
func getSymbolAmountFromCoin(coin string) (string, *big.Int) {
	coinRune := []rune(coin)
	amount := new(big.Int)

	var symbol string

	// Set up regex
	isLetter := regexp.MustCompile(`[a-z]`)

	// Iterate over each rune in the coin string
	for i, char := range coinRune {
		// Regex will match first letter [a-z] (lowercase)
		matched := isLetter.MatchString(string(char))

		// On first match, split the coin into (amount, symbol)
		if matched {
			amount, _ = amount.SetString(string(coinRune[0:i]), 10)
			symbol = string(coinRune[i:])

			break
		}
	}

	return symbol, amount
}
