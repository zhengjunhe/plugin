// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gossip

import (
	"time"
)

// time limit for timeout
var (
	UpdateState                 = 2 * time.Second
	PingTimeout                 = 14 * time.Second
	DefaultSendTimeout          = 10 * time.Second
	DialTimeout                 = 5 * time.Second
	mapUpdateInterval           = 45 * time.Hour
	StreamPingTimeout           = 20 * time.Second
	MonitorPeerInfoInterval     = 10 * time.Second
	MonitorPeerNumInterval      = 30 * time.Second
	MonitorReBalanceInterval    = 15 * time.Minute
	GetAddrFromAddrBookInterval = 5 * time.Second
	GetAddrFromOnlineInterval   = 5 * time.Second
	GetAddrFromGitHubInterval   = 5 * time.Minute
	CheckActivePeersInterVal    = 5 * time.Second
	CheckBlackListInterVal      = 30 * time.Second
	CheckCfgSeedsInterVal       = 1 * time.Minute
)

const (
	msgTx           = 1
	msgBlock        = 2
	tryMapPortTimes = 20
	maxSamIPNum     = 20
)

const (
	//defalutNatPort  = 23802
	maxOutBoundNum  = 25
	stableBoundNum  = 15
	maxAttemps      = 5
	protocol        = "tcp"
	externalPortTag = "externalport"
)

const (
	nodeNetwork = 1
	nodeGetUTXO = 2
	nodeBloom   = 4
)

const (
	// Service service number
	Service int32 = nodeBloom + nodeNetwork + nodeGetUTXO
)

// leveldb 中p2p privkey,addrkey
const (
	addrkeyTag = "addrs"
	privKeyTag = "privkey"
)

//TTL
const (
	DefaultLtTxBroadCastTTL  = 3
	DefaultMaxTxBroadCastTTL = 25
	// 100KB
	defaultMinLtBlockSize = 100
)

// P2pCacheTxSize p2pcache size of transaction
const (
	PeerAddrCacheNum = 1000
	//接收的交易哈希过滤缓存设为mempool最大接收交易量
	TxRecvFilterCacheNum = 10240
	BlockFilterCacheNum  = 50
	//发送过滤主要用于发送时冗余检测, 发送完即可以被删除, 维护较小缓存数
	TxSendFilterCacheNum  = 500
	BlockCacheNum         = 10
	MaxBlockCacheByteSize = 100 * 1024 * 1024
)

// TestNetSeeds test seeds of net
var TestNetSeeds = []string{
	"47.97.223.101:28805",
}

// MainNetSeeds built-in list of seed
var MainNetSeeds = []string{
	"116.62.14.25:28805",
	"114.55.95.234:28805",
	"115.28.184.14:28805",
	"39.106.166.159:28805",
	"39.106.193.172:28805",
	"47.106.114.93:28805",
	"120.76.100.165:28805",
	"120.24.85.66:28805",
	"120.24.92.123:28805",
	"161.117.7.127:28805",
	"161.117.9.54:28805",
	"161.117.5.95:28805",
	"161.117.7.28:28805",
	"161.117.8.242:28805",
	"161.117.6.193:28805",
	"161.117.8.158:28805",
	"47.88.157.209:28805",
	"47.74.215.41:28805",
	"47.74.128.69:28805",
	"47.74.178.226:28805",
	"47.88.154.76:28805",
	"47.74.151.226:28805",
	"47.245.31.41:28805",
	"47.245.57.239:28805",
	"47.245.54.118:28805",
	"47.245.54.121:28805",
	"47.245.56.140:28805",
	"47.245.52.211:28805",
	"47.91.88.195:28805",
	"47.91.72.71:28805",
	"47.91.91.38:28805",
	"47.91.94.224:28805",
	"47.91.75.191:28805",
	"47.254.152.172:28805",
	"47.252.0.181:28805",
	"47.90.246.246:28805",
	"47.90.208.100:28805",
	"47.89.182.70:28805",
	"47.90.207.173:28805",
	"47.89.188.54:28805",
}
