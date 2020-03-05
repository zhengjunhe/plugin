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
		NameEthBridgeClaimAction:        TyEthBridgeClaimAction,
		NameMsgBurnAction:               TyMsgBurnAction,
		NameMsgLockAction:               TyMsgLockAction,
		NameMsgLogInValidatorAction:     TyMsgLogInValidatorAction,
		NameMsgLogOutValidatorAction:    TyMsgLogOutValidatorAction,
		NameMsgSetConsensusNeededAction: TyMsgSetConsensusNeededAction,
	}
	//定义log的id和具体log类型及名称，填入具体自定义log类型
	logMap = map[int64]*types.LogInfo{
		TyEthBridgeClaimLog:        {Ty: reflect.TypeOf(ReceiptEthBridgeClaim{}), Name: "LogEthBridgeClaimAction"},
		TyMsgBurnLog:               {Ty: reflect.TypeOf(ReceiptBurn{}), Name: "MsgBurnLog"},
		TyMsgLockLog:               {Ty: reflect.TypeOf(ReceiptLock{}), Name: "MsgLockLog"},
		TyMsgLogInValidatorLog:     {Ty: reflect.TypeOf(ReceiptLogIn{}), Name: "MsgLogInValidatorLog"},
		TyMsgLogOutValidatorLog:    {Ty: reflect.TypeOf(ReceiptLogOut{}), Name: "MsgLogOutValidatorLog"},
		TyMsgSetConsensusNeededLog: {Ty: reflect.TypeOf(ReceiptSetConsensusNeeded{}), Name: "MsgSetConsensusNeededLog"},
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
