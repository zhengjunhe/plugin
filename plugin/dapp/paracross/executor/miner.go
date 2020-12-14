// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/paracross/types"
	"github.com/pkg/errors"
)

//当前miner tx不需要校验上一个区块的衔接性，因为tx就是本节点发出，高度，preHash等都在本区块里面的blockchain做了校验
//note: 平行链的Miner从Height=1开始， 创世区块不挖矿
func (a *action) Miner(miner *pt.ParacrossMinerAction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()
	if miner.Status.Title != cfg.GetTitle() || miner.Status.MainBlockHash == nil {
		return nil, pt.ErrParaMinerExecErr
	}

	var logs []*types.ReceiptLog
	var receipt = &pt.ReceiptParacrossMiner{}

	log := &types.ReceiptLog{}
	log.Ty = pt.TyLogParacrossMiner
	receipt.Status = miner.Status

	log.Log = types.Encode(receipt)
	logs = append(logs, log)

	minerReceipt := &types.Receipt{Ty: types.ExecOk, KV: nil, Logs: logs}

	on, err := a.isSelfConsensOn(miner)
	if err != nil {
		return nil, err
	}
	//自共识后才挖矿
	if on {
		r, err := a.issueCoins(miner)
		if err != nil {
			return nil, err
		}

		minerReceipt = mergeReceipt(minerReceipt, r)
	}

	return minerReceipt, nil
}

func (a *action) isSelfConsensOn(miner *pt.ParacrossMinerAction) (bool, error) {
	cfg := a.api.GetConfig()
	//ForkParaInitMinerHeight高度后，默认全部挖矿，产生在paracross执行器地址，如果自共识分阶段，也只是分阶段奖励，挖矿一直产生
	if cfg.IsDappFork(a.height, pt.ParaX, pt.ForkParaFullMinerHeight) {
		return true, nil
	}

	isSelfConsensOn := miner.IsSelfConsensus

	//自共识分阶段使能，综合考虑挖矿奖励和共识分配奖励，判断是否自共识使能需要采用共识的高度，而不能采用当前区块高度a.height
	//考虑自共识使能区块高度100，如果采用区块高度判断，则在100高度可能收到80~99的20条共识交易，这20条交易在100高度参与共识，则无奖励可分配，而且共识高度将是80而不是100
	//采用共识高度miner.Status.Height判断，则严格执行了产生奖励和分配奖励，且共识高度从100开始

	if cfg.IsDappFork(a.height, pt.ParaX, pt.ForkParaSelfConsStages) {
		var err error
		isSelfConsensOn, err = isSelfConsOn(a.db, miner.Status.Height)
		if err != nil && errors.Cause(err) != pt.ErrKeyNotExist {
			clog.Error("paracross miner getConsensus ", "height", miner.Status.Height, "err", err)
			return false, err
		}
	}
	return isSelfConsensOn, nil
}

const (
	normalMiner = iota
	halveMiner
	customMiner
)

type rewardValFn func(cfg *types.Chain33Config, height int64) (int64, int64, int64)

var getConfigRewards = make(map[int]rewardValFn)

func init() {
	getConfigRewards[normalMiner] = getNormalReward

}

func (a *action) issueCoins(miner *pt.ParacrossMinerAction) (*types.Receipt, error) {
	cfg := a.api.GetConfig()

	mode := int(cfg.MGInt("mver.consensus.paracross.minerMode", a.height))

	if getConfigRewards[mode] == nil {
		panic("getTotalReward not be set depend on consensus.paracross.minerMode")
	}

	coinReward, coinFundReward, _ := getConfigRewards[mode](cfg, a.height)
	totalReward := coinReward + coinFundReward
	if totalReward > 0 {
		issueReceipt, err := a.coinsAccount.ExecIssueCoins(a.execaddr, totalReward)
		if err != nil {
			clog.Error("paracross miner issue err", "height", miner.Status.Height,
				"execAddr", a.execaddr, "amount", totalReward)
			return nil, err
		}
		return issueReceipt, nil
	}
	return nil, nil
}