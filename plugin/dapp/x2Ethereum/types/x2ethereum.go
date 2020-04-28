package types

import (
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"reflect"
)

/*
 * 交易相关类型定义
 * 交易action通常有对应的log结构，用于交易回执日志记录
 * 每一种action和log需要用id数值和name名称加以区分
 */

var (
	//X2ethereumX 执行器名称定义
	X2ethereumX = "x2ethereum"
	//定义actionMap
	actionMap = map[string]int32{
		NameEth2Chain33Action:           TyEth2Chain33Action,
		NameWithdrawEthAction:           TyWithdrawEthAction,
		NameWithdrawChain33Action:       TyWithdrawChain33Action,
		NameChain33ToEthAction:          TyChain33ToEthAction,
		NameAddValidatorAction:          TyAddValidatorAction,
		NameRemoveValidatorAction:       TyRemoveValidatorAction,
		NameModifyPowerAction:           TyModifyPowerAction,
		NameSetConsensusThresholdAction: TySetConsensusThresholdAction,
		NameTransferAction:              TyTransferAction,
	}
	//定义log的id和具体log类型及名称，填入具体自定义log类型
	logMap = map[int64]*types.LogInfo{
		TyEth2Chain33Log:           {Ty: reflect.TypeOf(ReceiptEth2Chain33{}), Name: "LogEth2Chain33"},
		TyWithdrawEthLog:           {Ty: reflect.TypeOf(ReceiptEth2Chain33{}), Name: "LogWithdrawEth"},
		TyWithdrawChain33Log:       {Ty: reflect.TypeOf(ReceiptChain33ToEth{}), Name: "LogWithdrawChain33"},
		TyChain33ToEthLog:          {Ty: reflect.TypeOf(ReceiptChain33ToEth{}), Name: "LogChain33ToEth"},
		TyAddValidatorLog:          {Ty: reflect.TypeOf(ReceiptValidator{}), Name: "LogAddValidator"},
		TyRemoveValidatorLog:       {Ty: reflect.TypeOf(ReceiptValidator{}), Name: "LogRemoveValidator"},
		TyModifyPowerLog:           {Ty: reflect.TypeOf(ReceiptValidator{}), Name: "LogModifyPower"},
		TySetConsensusThresholdLog: {Ty: reflect.TypeOf(ReceiptSetConsensusThreshold{}), Name: "LogSetConsensusThreshold"},
		TyProphecyLog:              {Ty: reflect.TypeOf(ReceiptEthProphecy{}), Name: "LogEthProphecy"},
		TyTransferLog:              {Ty: reflect.TypeOf(types.ReceiptAccountTransfer{}), Name: "LogTransfer"},
	}
	tlog = log.New("module", "x2ethereum.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(X2ethereumX))
	//注册合约启用高度
	types.RegFork(X2ethereumX, InitFork)
	types.RegExec(X2ethereumX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(X2ethereumX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(X2ethereumX, NewType(cfg))
}

type x2ethereumType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *x2ethereumType {
	c := &x2ethereumType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (x *x2ethereumType) GetName() string {
	return X2ethereumX
}

// GetPayload 获取合约action结构
func (x *x2ethereumType) GetPayload() types.Message {
	return &X2EthereumAction{}
}

// GeTypeMap 获取合约action的id和name信息
func (x *x2ethereumType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap 获取合约log相关信息
func (x *x2ethereumType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}

// ActionName get PrivacyType action name
func (x x2ethereumType) ActionName(tx *types.Transaction) string {
	var action X2EthereumAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		return "unknown-x2ethereum-err"
	}
	tlog.Info("ActionName", "ActionName", action.GetActionName())
	return action.GetActionName()
}

// GetActionName get action name
func (action *X2EthereumAction) GetActionName() string {
	if action.Ty == TyEth2Chain33Action && action.GetEth2Chain33() != nil {
		return "Eth2Chain33"
	} else if action.Ty == TyWithdrawEthAction && action.GetWithdrawEth() != nil {
		return "WithdrawEth"
	} else if action.Ty == TyWithdrawChain33Action && action.GetWithdrawChain33() != nil {
		return "WithdrawChain33"
	} else if action.Ty == TyChain33ToEthAction && action.GetChain33ToEth() != nil {
		return "Chain33ToEth"
	} else if action.Ty == TyAddValidatorAction && action.GetAddValidator() != nil {
		return "AddValidator"
	} else if action.Ty == TyRemoveValidatorAction && action.GetRemoveValidator() != nil {
		return "RemoveValidator"
	} else if action.Ty == TyModifyPowerAction && action.GetModifyPower() != nil {
		return "ModifyPower"
	} else if action.Ty == TySetConsensusThresholdAction && action.GetSetConsensusThreshold() != nil {
		return "SetConsensusThreshold"
	}
	return "unknown-x2ethereum"
}
