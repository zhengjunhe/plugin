package rpc

import (
	"context"
	"github.com/33cn/chain33/types"
	jvmTypes "github.com/33cn/plugin/plugin/dapp/jvm/types"
	"github.com/33cn/plugin/plugin/dapp/jvm/executor"
)

func (c *channelClient) ExecFrozen(ctx context.Context, para *jvmTypes.TokenOperPara) (*jvmTypes.OpResult, error) {
	result := executor.ExecFrozen(para.From, para.Amount)
	return &jvmTypes.OpResult{Result:result}, nil
}

func (c *channelClient) ExecActive(ctx context.Context, para *jvmTypes.TokenOperPara) (*jvmTypes.OpResult, error) {
	result := executor.ExecActive(para.From, para.Amount)
	return &jvmTypes.OpResult{Result:result}, nil
}

func (c *channelClient) ExecTransfer(ctx context.Context, para *jvmTypes.TokenOperPara) (*jvmTypes.OpResult, error) {
	result := executor.ExecTransfer(para.From, para.To, para.Amount)
	return &jvmTypes.OpResult{Result:result}, nil
}

func (c *channelClient) GetRandom(ctx context.Context, para *types.ReqNil)(*jvmTypes.RandomData, error) {
	random, err := executor.GetRandom()
	return &jvmTypes.RandomData{Random:random}, err
}

func (c *channelClient) GetFrom(ctx context.Context, para *types.ReqNil) (*jvmTypes.FromAddr, error) {
	from := executor.GetFrom()
	return &jvmTypes.FromAddr{From:from}, nil
}

func (c *channelClient) SetStateDB(ctx context.Context, para *jvmTypes.SetDBPara) (*jvmTypes.OpResult, error) {
	result := executor.StateDBSetStateCallback(para.Key, para.Value)
	return &jvmTypes.OpResult{Result: result}, nil
}

func (c *channelClient) GetFromStateDB(ctx context.Context, para *jvmTypes.GetRequest) (*jvmTypes.GetResponse, error) {
	result := executor.StateDBGetStateCallback(para.Key)
	return &jvmTypes.GetResponse{Value: result}, nil
}

func (c *channelClient) GetValueSize(ctx context.Context, para *jvmTypes.GetRequest) (*jvmTypes.ValueSize, error) {
	size := executor.StateDBGetValueSizeCallback(para.ContractAddr, para.Key)
	return &jvmTypes.ValueSize{Size:size}, nil
}

func (c *channelClient) SetLocalDB(ctx context.Context, para *jvmTypes.SetDBPara) (*jvmTypes.OpResult, error) {
	result := executor.SetValue2Local(para.Key, para.Value)
	return &jvmTypes.OpResult{Result: result}, nil
}

func (c *channelClient) GetFromLocalDB(ctx context.Context, para *jvmTypes.GetRequest) (*jvmTypes.GetResponse, error) {
	result := executor.GetValueFromLocal(para.Key)
	return &jvmTypes.GetResponse{Value: result}, nil
}

func (c *channelClient) GetLocalValueSize(ctx context.Context, para *jvmTypes.GetRequest) (*jvmTypes.ValueSize, error) {
	size := executor.GetValueSizeFromLocal(para.Key)
	return &jvmTypes.ValueSize{Size:size}, nil
}


