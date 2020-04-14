package ethtxs

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethcontract/generated"
	"math/big"
)

func GetOperator(client *ethclient.Client, sender, bridgeBank common.Address) (common.Address, error) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		txslog.Error("GetOperator", "Failed to get HeaderByNumber due to:", err.Error())
		return common.Address{}, err
	}

	// Set up CallOpts auth
	auth := bind.CallOpts{
		Pending:     true,
		From:        sender,
		BlockNumber: header.Number,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	bridgeBankInstance, err := generated.NewBridgeBank(bridgeBank, client)
	if err != nil {
		txslog.Error("GetOperator", "Failed to NewBridgeBank to:", err.Error())
		return common.Address{}, err
	}

	return bridgeBankInstance.Operator(&auth)
}

func IsActiveValidator(validator common.Address, valset *generated.Valset) (bool, error) {
	opts := &bind.CallOpts{
		Pending:     true,
		From:        validator,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	isActiveValidator, err := valset.IsActiveValidator(opts, validator)
	if err != nil {
		txslog.Error("IsActiveValidator", "Failed to query IsActiveValidator due to:", err.Error())
		return false, err
	}

	return isActiveValidator, nil
}

func IsProphecyPending(id int64, validator common.Address, chain33Bridge *generated.Chain33Bridge) (bool, error) {
	opts := &bind.CallOpts{
		Pending:     true,
		From:        validator,
		Context:     context.Background(),
	}

	// Initialize BridgeRegistry instance
	active, err := chain33Bridge.IsProphecyClaimActive(opts, big.NewInt(id))
	if err != nil {
		txslog.Error("IsActiveValidatorFromChain33Bridge", "Failed to query IsActiveValidator due to:", err.Error())
		return false, err
	}
	return active, nil
}

func GetBalance(client *ethclient.Client, tokenAddr, owner string) (int64, error) {
	//查询ERC20余额
	if tokenAddr != "" {
		bridgeToken, err := generated.NewBridgeToken(common.HexToAddress(tokenAddr), client)
		if nil != err {
			return 0, err
		}
		ownerAddr := common.HexToAddress(owner)
		opts := &bind.CallOpts{
			Pending: true,
			From:    ownerAddr,
			Context: context.Background(),
		}
		balance, err := bridgeToken.BalanceOf(opts, ownerAddr)
		if nil != err {
			return 0, err
		}
		return balance.Int64(), nil
	}

	//查询ETH余额
	balance, err := client.BalanceAt(context.Background(), common.HexToAddress(owner), nil)
	if nil != err {
		return 0, err
	}
	return balance.Int64(), nil
}



