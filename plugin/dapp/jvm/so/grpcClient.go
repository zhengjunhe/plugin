package main

import (
	"C"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"golang.org/x/net/context"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/jvm/types"
	chain33Types "github.com/33cn/chain33/types"
)

var (
	clientActive = false
	serverAddrRunning = ""
	jlog = log.New("jvm", "so")
    jvmClient types.JvmClient
	clientConn *grpc.ClientConn
	//tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	//caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	//serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

type GrpcClientCfg struct {
	ServerAddr string
	Tls bool
	CaFile string
	ServerHostOverride string
}

//export StartDefaultGrpcClient
func StartDefaultGrpcClient(serverAddr string) bool {
	//"127.0.0.1:8802"
	return StartGrpcClient(serverAddr, "", "", false)
}

//export StartGrpcClient
func StartGrpcClient(ServerAddr, CaFile, ServerHostOverride string, tls bool) bool {
	cfg := GrpcClientCfg{
		ServerAddr:ServerAddr,
		Tls:tls,
		CaFile:CaFile,
		ServerHostOverride:ServerHostOverride,
	}

	flag.Parse()
	var opts []grpc.DialOption
	if cfg.Tls {
		if cfg.CaFile == "" {
			cfg.CaFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(cfg.CaFile, cfg.ServerHostOverride)
		if err != nil {
			jlog.Crit("startGrpcClient", "Failed to create TLS credentials", err.Error())
			return false
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	conn, err := grpc.Dial(cfg.ServerAddr, opts...)
	if err != nil {
		jlog.Crit("fail to dial: %v", err)
		return false
	}

	jvmClient = types.NewJvmClient(conn)
	clientActive = true
	clientConn = conn
	return true
}

//export StopGrpcClient
func StopGrpcClient() {
	_ = clientConn.Close()
}

//export ExecFrozen
func ExecFrozen(from string, amount int64) bool {
	in := &types.TokenOperPara{
		From:                 from,
		Amount:               amount,
	}
	_ , err := jvmClient.ExecFrozen(context.Background(), in)
	if nil != err {
		return false
	}
	return true
}

//export ExecActive
func ExecActive(from, execAddr string, amount int64) bool {
	in := &types.TokenOperPara{
		From:                 from,
		ExecAddr:             execAddr,
		Amount:               amount,
	}
	_ , err := jvmClient.ExecActive(context.Background(), in)
	if nil != err {
		return false
	}
	return true
}

//export ExecTransfer
func ExecTransfer(from, to, execAddr string, amount int64) bool {
	in := &types.TokenOperPara{
		From:                 from,
		To:                   to,
		ExecAddr:             execAddr,
		Amount:               amount,
	}
	_ , err := jvmClient.ExecTransfer(context.Background(), in)
	if nil != err {
		return false
	}
	return true
}

//export GetRandom
func GetRandom() []byte {
	in := &chain33Types.ReqNil{}
	random, err := jvmClient.GetRandom(context.Background(), in)
	if nil != err {
		return nil
	}
	return random.Random
}

//export GetFrom
func GetFrom() string {
	in := &chain33Types.ReqNil{}
	from, err := jvmClient.GetFrom(context.Background(), in)
	if nil != err {
		return ""
	}
	return from.From
}

//export SetState
func SetState(key, value []byte) bool {
	in := &types.SetDBPara{
		Key:key,
		Value:value,
	}
	result , err := jvmClient.SetStateDB(context.Background(), in)
	if nil != err {
		return false
	}
	return result.Result
}

//export GetFromState
func GetFromState(key, value []byte) bool {
	in := &types.GetRequest{
		Key:key,
	}
	result , err := jvmClient.GetFromStateDB(context.Background(), in)
	if nil != err {
		return false
	}
	copy(value, result.Value)
	return true
}

//export GetValueSize
func GetValueSize(key []byte) int32 {
	in := &types.GetRequest{
		Key:key,
	}
	result , err := jvmClient.GetValueSize(context.Background(), in)
	if nil != err {
		return 0
	}
	return result.Size
}

//export SetLocalDB
func SetLocalDB(key, value []byte) bool {
	in := &types.SetDBPara{
		Key:key,
		Value:value,
	}
	result , err := jvmClient.SetLocalDB(context.Background(), in)
	if nil != err {
		return false
	}
	return result.Result
}

//export GetFromLocalDB
func GetFromLocalDB(key, value []byte) bool {
	in := &types.GetRequest{
		Key:key,
	}
	result , err := jvmClient.GetFromLocalDB(context.Background(), in)
	if nil != err {
		return false
	}
	copy(value, result.Value)
	return true
}

//export GetLocalValueSize
func GetLocalValueSize(key []byte) int32 {
	in := &types.GetRequest{
		Key:key,
	}
	result , err := jvmClient.GetLocalValueSize(context.Background(), in)
	if nil != err {
		return 0
	}
	return result.Size
}

func main() {}




