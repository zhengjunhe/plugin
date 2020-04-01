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
		TySymbolAssetsLog:          {Ty: reflect.TypeOf(ReceiptQuerySymbolAssets{}), Name: "LogSymbolAssets"},
		TyProphecyLog:              {Ty: reflect.TypeOf(ReceiptEthProphecy{}), Name: "LogEthProphecy"},
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