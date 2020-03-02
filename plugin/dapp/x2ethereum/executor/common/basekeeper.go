package common

import (
	"encoding/json"
	"github.com/33cn/chain33/common/db"
	types2 "github.com/33cn/chain33/types"
)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	db db.KV
	// list of addresses that are restricted from receiving transactions
	blacklistedAddrs map[string]bool
}

// NewBaseSendKeeper returns a new BaseSendKeeper.
func NewBaseSendKeeper(db db.KV, blacklistedAddrs map[string]bool) BaseSendKeeper {

	return BaseSendKeeper{
		db:               db,
		blacklistedAddrs: blacklistedAddrs,
	}
}

// SendCoins moves coins from one account to another
func (keeper BaseSendKeeper) SendCoins(fromAddr string, toAddr string, amt int64) error {
	_, err := keeper.SubtractCoins(fromAddr, amt)
	if err != nil {
		return err
	}

	_, err = keeper.AddCoins(toAddr, amt)
	if err != nil {
		return err
	}

	return nil
}

func (keeper BaseSendKeeper) SubtractCoins(addr string, amt int64) (amount, error) {
	if amt < 0 {
		return amount{}, types2.ErrAmount
	}

	addrBytes, err := json.Marshal(addr)
	if err != nil {
		return amount{}, types2.ErrMarshal
	}
	old, err := keeper.db.Get(addrBytes)
	if err != nil && err != types2.ErrNotFound {
		return amount{}, err
	} else if err == types2.ErrNotFound {
		//stateDB中没有该地址的余额
		a := amount{
			spendable_amount: 0,
			locked_amount:    0,
		}
		amountBytes, err := json.Marshal(a)
		if err != nil {
			return amount{}, types2.ErrMarshal
		}
		_ = keeper.db.Set(addrBytes, amountBytes)
		return amount{}, err
	}

	var amounts amount
	err = json.Unmarshal(old, &amounts)
	if err != nil {
		return amount{}, types2.ErrUnmarshal
	}

	if amounts.spendable_amount < amt {
		return amount{}, types2.ErrInsufficientBalance
	}

	amounts.spendable_amount -= amt

	err = keeper.SetCoins(addr, amounts)
	if err != nil {
		return amount{}, err
	}

	return amounts, err
}

// AddCoins adds amt to the coins at the addr.
func (keeper BaseSendKeeper) AddCoins(addr string, amt int64) (amount, error) {
	if amt < 0 {
		return amount{}, types2.ErrAmount
	}

	addrBytes, err := json.Marshal(addr)
	if err != nil {
		return amount{}, types2.ErrMarshal
	}
	old, err := keeper.db.Get(addrBytes)
	if err != nil && err != types2.ErrNotFound {
		return amount{}, err
	} else if err == types2.ErrNotFound {
		//stateDB中没有该地址的余额
		a := amount{
			spendable_amount: 0,
			locked_amount:    0,
		}
		amountBytes, err := json.Marshal(a)
		if err != nil {
			return amount{}, types2.ErrMarshal
		}
		_ = keeper.db.Set(addrBytes, amountBytes)
		return amount{}, err
	}

	var amounts amount
	err = json.Unmarshal(old, &amounts)
	if err != nil {
		return amount{}, types2.ErrUnmarshal
	}

	amounts.spendable_amount += amt

	err = keeper.SetCoins(addr, amounts)
	if err != nil {
		return amount{}, err
	}
	return amounts, err
}

// SetCoins sets the coins at the addr.
func (keeper BaseSendKeeper) SetCoins(addr string, amt amount) error {
	if amt.locked_amount < 0 || amt.spendable_amount < 0 {
		return types2.ErrAmount
	}

	addrBytes, err := json.Marshal(addr)
	if err != nil {
		return types2.ErrMarshal
	}
	old, err := keeper.db.Get(addrBytes)
	if err != nil && err != types2.ErrNotFound {
		return err
	} else if err == types2.ErrNotFound {
		//stateDB中没有该地址的余额
		a := amount{
			spendable_amount: 0,
			locked_amount:    0,
		}
		amountBytes, err := json.Marshal(a)
		if err != nil {
			return types2.ErrMarshal
		}
		_ = keeper.db.Set(addrBytes, amountBytes)
		return nil
	}

	var amounts amount
	err = json.Unmarshal(old, &amounts)
	if err != nil {
		return types2.ErrUnmarshal
	}

	amounts.spendable_amount = amt.spendable_amount
	amounts.locked_amount = amt.locked_amount
	amountBytes, err := json.Marshal(amounts)
	if err != nil {
		return types2.ErrMarshal
	}
	_ = keeper.db.Set(addrBytes, amountBytes)

	return nil
}
