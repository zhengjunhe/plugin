package testnode

import (
	"github.com/33cn/dplatformos/types"
	"github.com/33cn/dplatformos/util/testnode"
)

/*
1. solo 模式，后台启动一个 主节点
2. 启动一个平行链节点：注意，这个要测试的话，会依赖平行链插件
*/

//ParaNode 平行链节点由两个节点组成
type ParaNode struct {
	Main *testnode.DplatformOSMock
	Para *testnode.DplatformOSMock
}

//NewParaNode 创建一个平行链节点
func NewParaNode(main *testnode.DplatformOSMock, para *testnode.DplatformOSMock) *ParaNode {
	if main == nil {
		main = testnode.New("", nil)
		main.Listen()
	}
	if para == nil {
		cfg := types.NewDplatformOSConfig(DefaultConfig)
		testnode.ModifyParaClient(cfg, main.GetCfg().RPC.GrpcBindAddr)
		para = testnode.NewWithConfig(cfg, nil)
		para.Listen()
	}
	return &ParaNode{Main: main, Para: para}
}

//Close 关闭系统
func (node *ParaNode) Close() {
	node.Para.Close()
	node.Main.Close()
}
