package executor

import (
	"encoding/json"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/ethbridge"
	types2 "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

type action struct {
	api          client.QueueProtocolAPI
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	index        int32
	execaddr     string
}

func newAction(a *x2ethereum, tx *types.Transaction, index int32) *action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &action{a.GetAPI(), a.GetCoinsAccount(), a.GetStateDB(), hash, fromaddr,
		a.GetBlockTime(), a.GetHeight(), index, dapp.ExecAddress(string(tx.Execer))}
}

func (a *action) procMsgEthBridgeClaim(ethBridgeClaim *types2.EthBridgeClaim) (*types.Receipt, error) {
	msgEthBridgeClaim := executor.NewMsgCreateEthBridgeClaim(*ethBridgeClaim)
	if err := msgEthBridgeClaim.ValidateBasic(); err != nil {
		return nil, err
	}
	oracleClaim, err := executor.CreateOracleClaimFromEthClaim(*ethBridgeClaim)
	if err != nil {
		return nil, err
	}

	status, err := a.ProcessClaim(oracleClaim)
	if err != nil {
		return nil, err
	}

	msgEthBridgeClaimBytes, err := json.Marshal(msgEthBridgeClaim)
	if err != nil {
		return nil, types.ErrMarshal
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return nil, types.ErrMarshal
	}
	receipt := &types.Receipt{KV: make([]*types.KeyValue, 0)}
	receipt.KV = append(receipt.KV, &types.KeyValue{Key: msgEthBridgeClaimBytes, Value: statusBytes})

	receipt.Ty = types.ExecOk
}

func (a *action) procMsgLock(msgLock *types2.MsgLock) (*types.Receipt, error) {

}

func (a *action) procMsgBurn(msgBurn *types2.MsgBurn) (*types.Receipt, error) {

}
