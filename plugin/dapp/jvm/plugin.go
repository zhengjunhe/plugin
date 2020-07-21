package jvm

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor"
	jvmtypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/rpc"

	// init auto test
	"github.com/33cn/plugin/plugin/dapp/jvm/commands"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     jvmtypes.JvmX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.JvmCmd,
		RPC:      rpc.Init,
	})
}
