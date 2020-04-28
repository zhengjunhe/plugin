package ethtxs

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
	chain33Types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"regexp"
	"strings"
)

// LogLockToEthBridgeClaim : parses and packages a LockEvent struct with a validator address in an EthBridgeClaim msg
func LogLockToEthBridgeClaim(event *events.LockEvent, ethereumChainID int64, bridgeBrankAddr string, decimal int64) (*ebrelayerTypes.EthBridgeClaim, error) {
	recipient := event.To
	if 0 == len(recipient) {
		return nil, ebrelayerTypes.ErrEmptyAddress
	}
	// Symbol formatted to lowercase
	symbol := strings.ToLower(event.Symbol)
	if symbol == "eth" && event.Token != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return nil, ebrelayerTypes.ErrAddress4Eth
	}

	witnessClaim := &ebrelayerTypes.EthBridgeClaim{}
	witnessClaim.EthereumChainID = ethereumChainID
	witnessClaim.BridgeBrankAddr = bridgeBrankAddr
	witnessClaim.Nonce = event.Nonce.Int64()
	witnessClaim.TokenAddr = event.Token.String()
	witnessClaim.Symbol = event.Symbol
	witnessClaim.EthereumSender = event.From.String()
	witnessClaim.Chain33Receiver = string(recipient)
	witnessClaim.Amount = event.Value.String()

	witnessClaim.ClaimType = types.LOCK_CLAIM_TYPE
	witnessClaim.ChainName = types.LOCK_CLAIM
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

func LogBurnToEthBridgeClaim(event *events.BurnEvent, ethereumChainID int64, bridgeBrankAddr string, decimal int64) (*ebrelayerTypes.EthBridgeClaim, error) {
	recipient := event.Chain33Receiver
	if 0 == len(recipient) {
		return nil, ebrelayerTypes.ErrEmptyAddress
	}

	witnessClaim := &ebrelayerTypes.EthBridgeClaim{}
	witnessClaim.EthereumChainID = ethereumChainID
	witnessClaim.BridgeBrankAddr = bridgeBrankAddr
	witnessClaim.Nonce = event.Nonce.Int64()
	witnessClaim.TokenAddr = event.Token.String()
	witnessClaim.Symbol = event.Symbol
	witnessClaim.EthereumSender = event.OwnerFrom.String()
	witnessClaim.Chain33Receiver = string(recipient)
	witnessClaim.Amount = event.Amount.String()
	witnessClaim.ClaimType = types.BURN_CLAIM_TYPE
	witnessClaim.ChainName = types.BURN_CLAIM
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

// BurnLockTxReceiptToChain33Msg : parses data from a Burn/Lock event witnessed on chain33 into a Chain33Msg struct
func BurnLockTxReceiptToChain33Msg(claimType events.Event, receipt *chain33Types.ReceiptData) *events.Chain33Msg {
	// Set up variables
	var chain33Sender []byte
	var ethereumReceiver, tokenContractAddress common.Address
	var symbol string
	var amount *big.Int

	// Iterate over attributes
	for _, log := range receipt.Logs {
		if log.Ty == types.TyChain33ToEthLog || log.Ty == types.TyWithdrawChain33Log {
			txslog.Debug("BurnLockTxReceiptToChain33Msg", "value", string(log.Log))
			var chain33ToEth types.ReceiptChain33ToEth
			err := chain33Types.Decode(log.Log, &chain33ToEth)
			if err != nil {
				return nil
			}
			chain33Sender = []byte(chain33ToEth.Chain33Sender)
			ethereumReceiver = common.HexToAddress(chain33ToEth.EthereumReceiver)
			tokenContractAddress = common.HexToAddress(chain33ToEth.TokenContract)
			symbol = chain33ToEth.EthSymbol
			chain33ToEth.Amount = types.TrimZeroAndDot(chain33ToEth.Amount)
			amount = big.NewInt(1)
			amount, _ = amount.SetString(chain33ToEth.Amount, 10)

			txslog.Info("BurnLockTxReceiptToChain33Msg", "chain33Sender", chain33Sender, "ethereumReceiver", ethereumReceiver.String(), "tokenContractAddress", tokenContractAddress.String(), "symbol", symbol, "amount", amount.String())
			// Package the event data into a Chain33Msg
			chain33Msg := events.NewChain33Msg(claimType, chain33Sender, ethereumReceiver, symbol, amount, tokenContractAddress)
			return &chain33Msg
		}
	}
	return nil
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

// Chain33MsgToProphecyClaim : parses event data from a Chain33Msg, packaging it as a ProphecyClaim
func Chain33MsgToProphecyClaim(event events.Chain33Msg) ProphecyClaim {
	claimType := event.ClaimType
	chain33Sender := event.Chain33Sender
	ethereumReceiver := event.EthereumReceiver
	tokenContractAddress := event.TokenContractAddress
	symbol := strings.ToLower(event.Symbol)
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		Chain33Sender:        chain33Sender,
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
