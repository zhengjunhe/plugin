// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package valnode

import (
	"github.com/33cn/dplatform/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/valnode/commands"
	"github.com/33cn/plugin/plugin/dapp/valnode/executor"
	"github.com/33cn/plugin/plugin/dapp/valnode/rpc"
	"github.com/33cn/plugin/plugin/dapp/valnode/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.ValNodeX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.ValCmd,
		RPC:      rpc.Init,
	})
}
