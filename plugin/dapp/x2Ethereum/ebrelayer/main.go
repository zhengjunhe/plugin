package main

import (
	"io"
	"context"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/limits"
	"github.com/prometheus/common/log"
	relayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	logf "github.com/33cn/chain33/common/log"
	chain33Types "github.com/33cn/chain33/types"
	tml "github.com/BurntSushi/toml"
	chain33Relayer "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/chain33"
	ethRelayer "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/ethereum"
	relayer "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer"
	dbm "github.com/33cn/chain33/common/db"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	configPath = flag.String("f", "", "configfile")
	versionCmd = flag.Bool("v", false, "version")
	IPWhiteListMap = make(map[string]bool)
)

func main() {
	flag.Parse()
	if *versionCmd {
		fmt.Println(relayerTypes.Version4Relayer)
		return
	}
	if *configPath == "" {
		*configPath = "relayer.toml"
	}

	err := os.Chdir(pwd())
	if err != nil {
		panic(err)
	}
	d, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Info("current dir:", "dir", d)
	err = limits.SetLimits()
	if err != nil {
		panic(err)
	}

	//set config: lns 用 lns.toml 这个配置文件
	cfg := initCfg(*configPath)
	log.Info("Starting FUZAMEI Chain33-X-Ethereum relayer software:", "Name:", cfg.Title)

	logf.SetFileLog(convertLogCfg(cfg.Log))

	_, cancel := context.WithCancel(context.Background())
	//创建blockchain服务，用于接收chain33的区块推送，过滤转发，以及转发闪电钱包的相关交易
	var wg sync.WaitGroup
	db := dbm.NewDB("relayer_db_service", cfg.SyncTxConfig.Dbdriver, cfg.SyncTxConfig.DbPath, cfg.SyncTxConfig.DbCache)
	chain33RelayerService := chain33Relayer.StartChain33Relayer(cfg.SyncTxConfig, db)
    ethRelayerService := ethRelayer.StartEthereumRelayer(cfg.SyncTxConfig.Chain33Host, db)
	relayerManager := relayer.NewRelayerManager(chain33RelayerService, ethRelayerService, db)

	startRpcServer(cfg.JrpcBindAddr, relayerManager)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		cancel()
		wg.Wait()
		os.Exit(0)
	}()
}

func convertLogCfg(log *relayerTypes.Log) *chain33Types.Log {
	return &chain33Types.Log{
		Loglevel:        log.Loglevel,
		LogConsoleLevel: log.LogConsoleLevel,
		LogFile:         log.LogFile,
		MaxFileSize:     log.MaxFileSize,
		MaxBackups:      log.MaxBackups,
		MaxAge:          log.MaxAge,
		LocalTime:       log.LocalTime,
		Compress:        log.Compress,
		CallerFile:      log.CallerFile,
		CallerFunction:  log.CallerFunction,
	}
}

func pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func initCfg(path string) *relayerTypes.RelayerConfig {
	var cfg relayerTypes.RelayerConfig
	if _, err := tml.DecodeFile(path, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println(cfg)
	return &cfg
}

func IsIPWhiteListEmpty() bool {
	return len(IPWhiteListMap) == 0
}

//判断ipAddr是否在ip地址白名单中
func IsInIPWhitelist(ipAddrPort string) bool {
	ipAddr, _, err := net.SplitHostPort(ipAddrPort)
	if err != nil {
		return false
	}
	ip := net.ParseIP(ipAddr)
	if ip.IsLoopback() {
		return true
	}
	if _, ok := IPWhiteListMap[ipAddr]; ok {
		return true
	}
	return false
}

type RPCServer struct {
	*rpc.Server
}

func (r *RPCServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Info("ServeHTTP", "request address", req.RemoteAddr)
	if !IsIPWhiteListEmpty() {
		if !IsInIPWhitelist(req.RemoteAddr) {
			log.Info("ServeHTTP", "refuse connect address", req.RemoteAddr)
			w.WriteHeader(401)
			return
		}
	}
	r.Server.ServeHTTP(w, req)
}

func (r *RPCServer) HandleHTTP(rpcPath, debugPath string) {
	http.Handle(rpcPath, r)
}

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

func startRpcServer(address string, api interface{}) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("监听失败，端口可能已经被占用")
		panic(err)
	}
	srv := &RPCServer{rpc.NewServer()}
	srv.Server.Register(api)
	srv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(200)
			err := srv.ServeRequest(serverCodec)
			if err != nil {
				log.Debug("http", "Error while serving JSON request: %v", err)
				return
			}
		}
	})
	http.Serve(listener, handler)
}