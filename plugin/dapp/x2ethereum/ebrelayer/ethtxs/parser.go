package ethtxs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      dplatform Bridge.
// --------------------------------------------------------

import (
	"math/big"
	"strings"

	dplatformTypes "github.com/33cn/dplatform/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/events"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"github.com/ethereum/go-ethereum/common"
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
	witnessClaim.DplatformReceiver = string(recipient)

	if decimal > 8 {
		event.Value = event.Value.Quo(event.Value, big.NewInt(int64(types.MultiplySpecifyTimes(1, decimal-8))))
	} else {
		event.Value = event.Value.Mul(event.Value, big.NewInt(int64(types.MultiplySpecifyTimes(1, 8-decimal))))
	}
	witnessClaim.Amount = event.Value.String()

	witnessClaim.ClaimType = types.LockClaimType
	witnessClaim.ChainName = types.LockClaim
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

//LogBurnToEthBridgeClaim ...
func LogBurnToEthBridgeClaim(event *events.BurnEvent, ethereumChainID int64, bridgeBrankAddr string, decimal int64) (*ebrelayerTypes.EthBridgeClaim, error) {
	recipient := event.DplatformReceiver
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
	witnessClaim.DplatformReceiver = string(recipient)
	witnessClaim.Amount = event.Amount.String()
	witnessClaim.ClaimType = types.BurnClaimType
	witnessClaim.ChainName = types.BurnClaim
	witnessClaim.Decimal = decimal

	return witnessClaim, nil
}

// ParseBurnLockTxReceipt : parses data from a Burn/Lock event witnessed on dplatform into a DplatformMsg struct
func ParseBurnLockTxReceipt(claimType events.Event, receipt *dplatformTypes.ReceiptData) *events.DplatformMsg {
	// Set up variables
	var dplatformSender []byte
	var ethereumReceiver, tokenContractAddress common.Address
	var symbol string
	var amount *big.Int

	// Iterate over attributes
	for _, log := range receipt.Logs {
		if log.Ty == types.TyDplatformToEthLog || log.Ty == types.TyWithdrawDplatformLog {
			txslog.Debug("ParseBurnLockTxReceipt", "value", string(log.Log))
			var dplatformToEth types.ReceiptDplatformToEth
			err := dplatformTypes.Decode(log.Log, &dplatformToEth)
			if err != nil {
				return nil
			}
			dplatformSender = []byte(dplatformToEth.DplatformSender)
			ethereumReceiver = common.HexToAddress(dplatformToEth.EthereumReceiver)
			tokenContractAddress = common.HexToAddress(dplatformToEth.TokenContract)
			symbol = dplatformToEth.IssuerDotSymbol
			dplatformToEth.Amount = types.TrimZeroAndDot(dplatformToEth.Amount)
			amount = big.NewInt(1)
			amount, _ = amount.SetString(dplatformToEth.Amount, 10)
			if dplatformToEth.Decimals > 8 {
				amount = amount.Mul(amount, big.NewInt(int64(types.MultiplySpecifyTimes(1, dplatformToEth.Decimals-8))))
			} else {
				amount = amount.Quo(amount, big.NewInt(int64(types.MultiplySpecifyTimes(1, 8-dplatformToEth.Decimals))))
			}

			txslog.Info("ParseBurnLockTxReceipt", "dplatformSender", dplatformSender, "ethereumReceiver", ethereumReceiver.String(), "tokenContractAddress", tokenContractAddress.String(), "symbol", symbol, "amount", amount.String())
			// Package the event data into a DplatformMsg
			dplatformMsg := events.NewDplatformMsg(claimType, dplatformSender, ethereumReceiver, symbol, amount, tokenContractAddress)
			return &dplatformMsg
		}
	}
	return nil
}

// DplatformMsgToProphecyClaim : parses event data from a DplatformMsg, packaging it as a ProphecyClaim
func DplatformMsgToProphecyClaim(event events.DplatformMsg) ProphecyClaim {
	claimType := event.ClaimType
	dplatformSender := event.DplatformSender
	ethereumReceiver := event.EthereumReceiver
	tokenContractAddress := event.TokenContractAddress
	symbol := strings.ToLower(event.Symbol)
	amount := event.Amount

	prophecyClaim := ProphecyClaim{
		ClaimType:            claimType,
		DplatformSender:        dplatformSender,
		EthereumReceiver:     ethereumReceiver,
		TokenContractAddress: tokenContractAddress,
		Symbol:               symbol,
		Amount:               amount,
	}

	return prophecyClaim
}
