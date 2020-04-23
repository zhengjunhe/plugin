package types

import (
	"encoding/json"
	"fmt"
	"github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
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
func (x *x2ethereumType) ActionName(tx *types.Transaction) string {
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

func GetDecimalsFromDB(addr string, db db.KV) (int64, error) {
	res, err := db.Get(CalAddr2DecimalsPrefix())
	if err != nil {
		return 0, err
	}
	var addr2Decimals map[string]int64
	err = json.Unmarshal(res, &addr2Decimals)
	if err != nil {
		return 0, err
	}
	if d, ok := addr2Decimals[addr]; ok {
		return d, nil
	}
	return 0, types.ErrNotFound
}

func GetDecimals(addr string) (int64, error) {
	if addr == "0x0000000000000000000000000000000000000000" || addr == "" {
		return 18, nil
	}
	Hashprefix := "0x313ce567"
	postData := fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","method":"eth_call","params":[{"to":"%s", "data":"%s"},"latest"]}`, addr, Hashprefix)

	retryTimes := 0
RETRY:
	res, err := sendToServer(EthNodeUrl, strings.NewReader(postData))
	if err != nil {
		tlog.Error("GetDecimals", "error:", err.Error())
		if retryTimes > 3 {
			return 0, err
		}
		retryTimes++
		goto RETRY
	}
	js, err := simplejson.NewJson(res)
	if err != nil {
		tlog.Error("GetDecimals", "NewJson error:", err.Error())
		if retryTimes > 3 {
			return 0, err
		}
		retryTimes++
		goto RETRY
	}
	result := js.Get("result").MustString()

	decimals, err := strconv.ParseInt(result, 0, 64)
	if err != nil {
		if retryTimes > 3 {
			return 0, err
		}
		retryTimes++
		goto RETRY
	}
	return decimals, nil
}

func sendToServer(url string, req io.Reader) ([]byte, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	var request *http.Request
	var err error

	request, err = http.NewRequest("POST", url, req)

	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
