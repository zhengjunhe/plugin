package common

import (
	"github.com/33cn/chain33/account"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

type amount struct {
	spendable_amount int64
	locked_amount    int64
}

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) SendCoinsFromModuleToAccount(senderModule, recipientAddr, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	senderAddr, ok := k.moduleAddresses[senderModule]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	receipt, err := accDB.ExecTransfer(senderAddr, recipientAddr, execAddr, amt)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k Keeper) SendCoinsFromModuleToModule(senderModule, recipientModule, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	senderAddr, ok := k.moduleAddresses[senderModule]
	if !ok {
		return nil, types.ErrUnknownAddress
	}
	recipientAddr, ok := k.moduleAddresses[recipientModule]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	receipt, err := accDB.ExecTransfer(senderAddr, recipientAddr, execAddr, amt)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k Keeper) SendCoinsFromAccountToModule(senderAddr, recipientModule, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	recipientAddr, ok := k.moduleAddresses[recipientModule]
	if !ok {
		return nil, types.ErrUnknownAddress
	}
	receipt, err := accDB.ExecTransfer(senderAddr, recipientAddr, execAddr, amt)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintCoins(amt int64, moduleName, execAddr string, accDB *account.DB) (*types2.Receipt, error) {

	bankAddr, ok := k.moduleAddresses[moduleName]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	receipt, err := accDB.ExecDeposit(bankAddr, execAddr, amt)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(amt int64, moduleName, execAddr string, accDB *account.DB) (*types2.Receipt, error) {

	bankAddr, ok := k.moduleAddresses[moduleName]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	receipt, err := accDB.ExecWithdraw(bankAddr, execAddr, amt)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}
