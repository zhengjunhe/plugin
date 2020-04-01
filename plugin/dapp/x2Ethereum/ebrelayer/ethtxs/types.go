package ethtxs

import (
	"math/big"

	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/events"
	"github.com/ethereum/go-ethereum/common"
)

const (
	X2Eth      = "x2ethereum"
	BurnAction = "burn"
	LockAction = "lock"
)

// OracleClaim : contains data required to make an OracleClaim
type OracleClaim struct {
	ProphecyID *big.Int
	Message    [32]byte
	Signature  []byte
}

// ProphecyClaim : contains data required to make an ProphecyClaim
type ProphecyClaim struct {
	ClaimType            events.Event
	CosmosSender         []byte
	EthereumReceiver     common.Address
	TokenContractAddress common.Address
	Symbol               string
	Amount               *big.Int
}
