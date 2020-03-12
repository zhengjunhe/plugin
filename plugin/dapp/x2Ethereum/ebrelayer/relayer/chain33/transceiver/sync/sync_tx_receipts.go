// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package sync

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	dbm "github.com/33cn/chain33/common/db"
	l "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/rpc/jsonclient"
	"github.com/33cn/chain33/types"
	relayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/rs/cors"
)

var (
	log            = l.New("module", "sync.tx_receipts")
	syncTxReceipts *SyncTxReceipts
)

func StartSyncTxReceipt(cfg *relayerTypes.SyncTxReceiptConfig, syncChan chan<- int64, db dbm.DB) *SyncTxReceipts {
	log.Debug("StartSyncTxReceipt, load config", "para:", cfg)
	log.Debug("SyncTxReceipts started ")

	bind(cfg.Chain33Host, cfg.PushName, "http://"+cfg.PushHost, "proto", cfg.StartSyncHeight)
	syncTxReceipts = NewSyncTxReceipts(db, syncChan)
	go syncTxReceipts.SaveAndSyncTxs2Relayer()
	go startHTTPService(cfg.PushBind, "*")
	return syncTxReceipts
}

func StopSyncTxReceipt() {
	syncTxReceipts.Stop()
}

func startHTTPService(url string, clientHost string) {
	listen, err := net.Listen("tcp", url)
	if err != nil {
		panic(err)
	}
	var handler http.Handler = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//fmt.Println(r.URL, r.Header, r.Body)
			beg := types.Now()
			defer func() {
				log.Info("handler", "cost", types.Since(beg))
			}()

			client := strings.Split(r.RemoteAddr, ":")[0]
			if !checkClient(client, clientHost) {
				log.Error("HandlerFunc", "client", r.RemoteAddr, "expect", clientHost)
				w.Write([]byte(`{"errcode":"-1","result":null,"msg":"reject"}`))
				// unbind 逻辑有问题， 需要的外部处理
				//  切换外部服务时， 可能换 name
				// 收到一个不是client 的请求，很有可能是以前注册过的， 取消掉
				//unbind(client)
				return
			}

			if r.URL.Path == "/" {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(200)
				if len(r.Header["Content-Encoding"]) >= 1 && r.Header["Content-Encoding"][0] == "gzip" {
					gr, err := gzip.NewReader(r.Body)
					//body := make([]byte, r.ContentLength)
					body, err := ioutil.ReadAll(gr)
					//n, err := r.Body.Read(body)
					if err != nil {
						log.Debug("Error while serving JSON request: %v", err)
						return
					}

					err = handleRequest(body)
					if err == nil {
						w.Write([]byte("OK"))
					} else {
						w.Write([]byte(err.Error()))
					}
				}
			}
		})

	co := cors.New(cors.Options{})
	handler = co.Handler(handler)

	http.Serve(listen, handler)
}

func handleRequest(body []byte) error {
	beg := types.Now()
	defer func() {
		log.Info("handleRequest", "cost", types.Since(beg))
	}()

	var req types.TxReceipts4Subscribe
	err := types.Decode(body, &req)
	if err != nil {
		log.Error("handleRequest", "DecodeBlockSeqErr", err)
		return err
	}

	err = pushTxReceipts(&req)
	log.Info("response", "err", err)
	return err
}

func checkClient(addr string, expectClient string) bool {
	if expectClient == "0.0.0.0" || expectClient == "*" {
		return true
	}
	return addr == expectClient
}

func bind(rpcAddr, name, url, encode string, startHeight int64) {
	params := types.SubscribeTxReceipt{
		Name:   name,
		URL:    url,
		Encode: encode,
		LastHeight:startHeight,
		Contract:"coins",
	}
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcAddr, "Chain33.AddSubscribeTxReceipt", params, &res)
	_, err := ctx.RunResult()
	if err != nil {
		fmt.Println("Failed to AddSubscribeTxReceipt to  rpc addr:", rpcAddr, "ReplySubTxReceipt", res)
		panic("bind client failed due to:" + err.Error())
	}
	log.Info("bind", "Succeed to AddSubscribeTxReceipt for rpc address:", rpcAddr)
	fmt.Println("Succeed to AddSubscribeTxReceipt")
}
