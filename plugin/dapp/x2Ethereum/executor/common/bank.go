package common

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/system/dapp"
	types2 "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/types"
)

// SendCoinsFromModuleToAccount transfers coins from a ModuleAccount to an AccAddress
func (k *Keeper) SendCoinsFromModuleToAccount(tokenSymbol, recipientAddr, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	senderAddr, ok := k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	clog.Info("SendCoinsFromModuleToAccount", "from", senderAddr, "to", recipientAddr)
	receipt, err := accDB.Transfer(senderAddr, recipientAddr, amt)
	if err != nil {
		clog.Error("SendCoinsFromModuleToAccount", "ExecTransfer error ", err)
		return nil, err
	}
	return receipt, nil
}

// SendCoinsFromModuleToModule transfers coins from a ModuleAccount to another
func (k *Keeper) SendCoinsFromModuleToModule(senderTokenSymbol, recipientTokenSymbol, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	senderAddr, ok := k.moduleAddresses[types.CalTokenSymbol(senderTokenSymbol)]
	if !ok {
		return nil, types.ErrUnknownAddress
	}
	recipientAddr, ok := k.moduleAddresses[types.CalTokenSymbol(recipientTokenSymbol)]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	clog.Info("SendCoinsFromModuleToModule", "from", senderAddr, "to", recipientAddr)
	receipt, err := accDB.Transfer(senderAddr, recipientAddr, amt)
	if err != nil {
		clog.Error("SendCoinsFromModuleToModule", "ExecTransfer error ", err)
		return nil, err
	}
	return receipt, nil
}

// SendCoinsFromAccountToModule transfers coins from an AccAddress to a ModuleAccount
func (k *Keeper) SendCoinsFromAccountToModule(senderAddr, tokenSymbol, execAddr string, amt int64, accDB *account.DB) (*types2.Receipt, error) {
	recipientAddr, ok := k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	if !ok {
		return nil, types.ErrUnknownAddress
	}

	clog.Info("SendCoinsFromAccountToModule", "from", senderAddr, "to", recipientAddr)
	receipt, err := accDB.Transfer(senderAddr, recipientAddr, amt)
	if err != nil {
		clog.Error("SendCoinsFromAccountToModule", "ExecTransfer error ", err)
		return nil, err
	}
	return receipt, nil
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k *Keeper) MintCoins(amt int64, tokenSymbol string, accDB *account.DB) (*types2.Receipt, error) {
	bankAddr, ok := k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	if !ok {
		err := k.AddAddressMap(tokenSymbol)
		if err != nil {
			return nil, err
		}
		bankAddr = k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	}
	clog.Info("MintCoins", "to", bankAddr, "amount", amt)

	receipt, err := accDB.Mint(bankAddr, amt)
	if err != nil {
		clog.Error("MintCoins", "ExecDeposit error ", err)
		return nil, err
	}

	return receipt, nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k *Keeper) BurnCoins(amt int64, tokenSymbol string, accDB *account.DB) (*types2.Receipt, error) {

	bankAddr, ok := k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	if !ok {
		return nil, types.ErrUnknownAddress
	}
	clog.Info("BurnCoins", "burn", bankAddr, "amount", amt)

	receipt, err := accDB.Burn(bankAddr, amt)
	if err != nil {
		clog.Error("BurnCoins", "Burn error ", err)
		return nil, err
	}

	return receipt, nil
}

func (k *Keeper) AddAddressMap(tokenSymbol string) error {

	_, ok := k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)]
	if !ok {
		addr := dapp.ExecAddress(types.CalTokenSymbol(tokenSymbol))
		k.moduleAddresses[types.CalTokenSymbol(tokenSymbol)] = addr
		return nil
	}

	return nil
}
