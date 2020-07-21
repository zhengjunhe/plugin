package so

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"golang.org/x/net/context"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/jvm/types"
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
//export: StartGrpcClient
func StartGrpcClient(cfg GrpcClientCfg) bool {
	if cfg.ServerAddr == serverAddrRunning && clientActive {
		return true
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

func StopGrpcClient() {
	_ = clientConn.Close()
}

func ExecFrozen(from, execAddr string, amount int64) bool {
	in := &types.TokenOperPara{
		From:                 from,
		ExecAddr:             execAddr,
		Amount:               amount,
	}
	_ , err := jvmClient.ExecFrozen(context.Background(), in)
	if nil != err {
		return false
	}
	return true
}

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
