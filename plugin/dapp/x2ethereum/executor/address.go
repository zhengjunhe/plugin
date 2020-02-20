package executor

import (
	"fmt"
	"github.com/33cn/chain33/common/address"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"reflect"
)

// EthereumAddress defines a standard ethereum address
type EthAddress gethCommon.Address

func NewEthAddressByProto(ethereumAddress string) EthAddress {
	return EthAddress(gethCommon.HexToAddress(ethereumAddress))
}

// NewEthereumAddress is a constructor function for EthereumAddress
func NewEthereumAddress(address string) EthAddress {
	return EthAddress(gethCommon.HexToAddress(address))
}

// Route should return the name of the module
func (ethAddr EthAddress) String() string {
	return gethCommon.Address(ethAddr).String()
}

// MarshalJSON marshals the etherum address to JSON
func (ethAddr EthAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", ethAddr.String())), nil
}

// UnmarshalJSON unmarshals an ethereum address
func (ethAddr *EthAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(gethCommon.Address{}), input, ethAddr[:])
}

type Chain33Address address.Address

func NewChain33AddressByProto(chain33Address string) Chain33Address {
	addr, _ := address.NewAddrFromString(chain33Address)
	return Chain33Address(*addr)
}
