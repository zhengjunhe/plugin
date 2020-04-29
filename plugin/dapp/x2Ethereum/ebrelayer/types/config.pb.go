// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package types

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SyncTxConfig struct {
	Chain33Host          string   `protobuf:"bytes,1,opt,name=chain33host" json:"chain33host,omitempty"`
	PushHost             string   `protobuf:"bytes,2,opt,name=pushHost" json:"pushHost,omitempty"`
	PushName             string   `protobuf:"bytes,3,opt,name=pushName" json:"pushName,omitempty"`
	PushBind             string   `protobuf:"bytes,4,opt,name=pushBind" json:"pushBind,omitempty"`
	MaturityDegree       int32    `protobuf:"varint,5,opt,name=maturityDegree" json:"maturityDegree,omitempty"`
	Dbdriver             string   `protobuf:"bytes,6,opt,name=dbdriver" json:"dbdriver,omitempty"`
	DbPath               string   `protobuf:"bytes,7,opt,name=dbPath" json:"dbPath,omitempty"`
	DbCache              int32    `protobuf:"varint,8,opt,name=dbCache" json:"dbCache,omitempty"`
	FetchHeightPeriodMs  int64    `protobuf:"varint,9,opt,name=fetchHeightPeriodMs" json:"fetchHeightPeriodMs,omitempty"`
	StartSyncHeight      int64    `protobuf:"varint,10,opt,name=startSyncHeight" json:"startSyncHeight,omitempty"`
	StartSyncSequence    int64    `protobuf:"varint,11,opt,name=startSyncSequence" json:"startSyncSequence,omitempty"`
	StartSyncHash        string   `protobuf:"bytes,12,opt,name=startSyncHash" json:"startSyncHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncTxConfig) Reset()         { *m = SyncTxConfig{} }
func (m *SyncTxConfig) String() string { return proto.CompactTextString(m) }
func (*SyncTxConfig) ProtoMessage()    {}
func (*SyncTxConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_de1752c779fc3508, []int{0}
}
func (m *SyncTxConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTxConfig.Unmarshal(m, b)
}
func (m *SyncTxConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTxConfig.Marshal(b, m, deterministic)
}
func (dst *SyncTxConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTxConfig.Merge(dst, src)
}
func (m *SyncTxConfig) XXX_Size() int {
	return xxx_messageInfo_SyncTxConfig.Size(m)
}
func (m *SyncTxConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncTxConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SyncTxConfig proto.InternalMessageInfo

func (m *SyncTxConfig) GetChain33Host() string {
	if m != nil {
		return m.Chain33Host
	}
	return ""
}

func (m *SyncTxConfig) GetPushHost() string {
	if m != nil {
		return m.PushHost
	}
	return ""
}

func (m *SyncTxConfig) GetPushName() string {
	if m != nil {
		return m.PushName
	}
	return ""
}

func (m *SyncTxConfig) GetPushBind() string {
	if m != nil {
		return m.PushBind
	}
	return ""
}

func (m *SyncTxConfig) GetMaturityDegree() int32 {
	if m != nil {
		return m.MaturityDegree
	}
	return 0
}

func (m *SyncTxConfig) GetDbdriver() string {
	if m != nil {
		return m.Dbdriver
	}
	return ""
}

func (m *SyncTxConfig) GetDbPath() string {
	if m != nil {
		return m.DbPath
	}
	return ""
}

func (m *SyncTxConfig) GetDbCache() int32 {
	if m != nil {
		return m.DbCache
	}
	return 0
}

func (m *SyncTxConfig) GetFetchHeightPeriodMs() int64 {
	if m != nil {
		return m.FetchHeightPeriodMs
	}
	return 0
}

func (m *SyncTxConfig) GetStartSyncHeight() int64 {
	if m != nil {
		return m.StartSyncHeight
	}
	return 0
}

func (m *SyncTxConfig) GetStartSyncSequence() int64 {
	if m != nil {
		return m.StartSyncSequence
	}
	return 0
}

func (m *SyncTxConfig) GetStartSyncHash() string {
	if m != nil {
		return m.StartSyncHash
	}
	return ""
}

type Log struct {
	Loglevel             string   `protobuf:"bytes,1,opt,name=loglevel" json:"loglevel,omitempty"`
	LogConsoleLevel      string   `protobuf:"bytes,2,opt,name=logConsoleLevel" json:"logConsoleLevel,omitempty"`
	LogFile              string   `protobuf:"bytes,3,opt,name=logFile" json:"logFile,omitempty"`
	MaxFileSize          uint32   `protobuf:"varint,4,opt,name=maxFileSize" json:"maxFileSize,omitempty"`
	MaxBackups           uint32   `protobuf:"varint,5,opt,name=maxBackups" json:"maxBackups,omitempty"`
	MaxAge               uint32   `protobuf:"varint,6,opt,name=maxAge" json:"maxAge,omitempty"`
	LocalTime            bool     `protobuf:"varint,7,opt,name=localTime" json:"localTime,omitempty"`
	Compress             bool     `protobuf:"varint,8,opt,name=compress" json:"compress,omitempty"`
	CallerFile           bool     `protobuf:"varint,9,opt,name=callerFile" json:"callerFile,omitempty"`
	CallerFunction       bool     `protobuf:"varint,10,opt,name=callerFunction" json:"callerFunction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}
func (*Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_de1752c779fc3508, []int{1}
}
func (m *Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Log.Unmarshal(m, b)
}
func (m *Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Log.Marshal(b, m, deterministic)
}
func (dst *Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Log.Merge(dst, src)
}
func (m *Log) XXX_Size() int {
	return xxx_messageInfo_Log.Size(m)
}
func (m *Log) XXX_DiscardUnknown() {
	xxx_messageInfo_Log.DiscardUnknown(m)
}

var xxx_messageInfo_Log proto.InternalMessageInfo

func (m *Log) GetLoglevel() string {
	if m != nil {
		return m.Loglevel
	}
	return ""
}

func (m *Log) GetLogConsoleLevel() string {
	if m != nil {
		return m.LogConsoleLevel
	}
	return ""
}

func (m *Log) GetLogFile() string {
	if m != nil {
		return m.LogFile
	}
	return ""
}

func (m *Log) GetMaxFileSize() uint32 {
	if m != nil {
		return m.MaxFileSize
	}
	return 0
}

func (m *Log) GetMaxBackups() uint32 {
	if m != nil {
		return m.MaxBackups
	}
	return 0
}

func (m *Log) GetMaxAge() uint32 {
	if m != nil {
		return m.MaxAge
	}
	return 0
}

func (m *Log) GetLocalTime() bool {
	if m != nil {
		return m.LocalTime
	}
	return false
}

func (m *Log) GetCompress() bool {
	if m != nil {
		return m.Compress
	}
	return false
}

func (m *Log) GetCallerFile() bool {
	if m != nil {
		return m.CallerFile
	}
	return false
}

func (m *Log) GetCallerFunction() bool {
	if m != nil {
		return m.CallerFunction
	}
	return false
}

type RelayerConfig struct {
	Title                string        `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	SyncTxConfig         *SyncTxConfig `protobuf:"bytes,2,opt,name=syncTxConfig" json:"syncTxConfig,omitempty"`
	Log                  *Log          `protobuf:"bytes,3,opt,name=log" json:"log,omitempty"`
	JrpcBindAddr         string        `protobuf:"bytes,4,opt,name=jrpcBindAddr" json:"jrpcBindAddr,omitempty"`
	EthProvider          string        `protobuf:"bytes,5,opt,name=ethProvider" json:"ethProvider,omitempty"`
	BridgeRegistry       string        `protobuf:"bytes,6,opt,name=bridgeRegistry" json:"bridgeRegistry,omitempty"`
	Deploy               *Deploy       `protobuf:"bytes,7,opt,name=deploy" json:"deploy,omitempty"`
	EthMaturityDegree    int32         `protobuf:"varint,8,opt,name=ethMaturityDegree" json:"ethMaturityDegree,omitempty"`
	EthBlockFetchPeriod  int32         `protobuf:"varint,9,opt,name=ethBlockFetchPeriod" json:"ethBlockFetchPeriod,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RelayerConfig) Reset()         { *m = RelayerConfig{} }
func (m *RelayerConfig) String() string { return proto.CompactTextString(m) }
func (*RelayerConfig) ProtoMessage()    {}
func (*RelayerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_de1752c779fc3508, []int{2}
}
func (m *RelayerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayerConfig.Unmarshal(m, b)
}
func (m *RelayerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayerConfig.Marshal(b, m, deterministic)
}
func (dst *RelayerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayerConfig.Merge(dst, src)
}
func (m *RelayerConfig) XXX_Size() int {
	return xxx_messageInfo_RelayerConfig.Size(m)
}
func (m *RelayerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_RelayerConfig proto.InternalMessageInfo

func (m *RelayerConfig) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *RelayerConfig) GetSyncTxConfig() *SyncTxConfig {
	if m != nil {
		return m.SyncTxConfig
	}
	return nil
}

func (m *RelayerConfig) GetLog() *Log {
	if m != nil {
		return m.Log
	}
	return nil
}

func (m *RelayerConfig) GetJrpcBindAddr() string {
	if m != nil {
		return m.JrpcBindAddr
	}
	return ""
}

func (m *RelayerConfig) GetEthProvider() string {
	if m != nil {
		return m.EthProvider
	}
	return ""
}

func (m *RelayerConfig) GetBridgeRegistry() string {
	if m != nil {
		return m.BridgeRegistry
	}
	return ""
}

func (m *RelayerConfig) GetDeploy() *Deploy {
	if m != nil {
		return m.Deploy
	}
	return nil
}

func (m *RelayerConfig) GetEthMaturityDegree() int32 {
	if m != nil {
		return m.EthMaturityDegree
	}
	return 0
}

func (m *RelayerConfig) GetEthBlockFetchPeriod() int32 {
	if m != nil {
		return m.EthBlockFetchPeriod
	}
	return 0
}

type SyncTxReceiptConfig struct {
	Chain33Host          string   `protobuf:"bytes,1,opt,name=chain33host" json:"chain33host,omitempty"`
	PushHost             string   `protobuf:"bytes,2,opt,name=pushHost" json:"pushHost,omitempty"`
	PushName             string   `protobuf:"bytes,3,opt,name=pushName" json:"pushName,omitempty"`
	PushBind             string   `protobuf:"bytes,4,opt,name=pushBind" json:"pushBind,omitempty"`
	StartSyncHeight      int64    `protobuf:"varint,5,opt,name=startSyncHeight" json:"startSyncHeight,omitempty"`
	StartSyncSequence    int64    `protobuf:"varint,6,opt,name=startSyncSequence" json:"startSyncSequence,omitempty"`
	StartSyncHash        string   `protobuf:"bytes,7,opt,name=startSyncHash" json:"startSyncHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncTxReceiptConfig) Reset()         { *m = SyncTxReceiptConfig{} }
func (m *SyncTxReceiptConfig) String() string { return proto.CompactTextString(m) }
func (*SyncTxReceiptConfig) ProtoMessage()    {}
func (*SyncTxReceiptConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_de1752c779fc3508, []int{3}
}
func (m *SyncTxReceiptConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTxReceiptConfig.Unmarshal(m, b)
}
func (m *SyncTxReceiptConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTxReceiptConfig.Marshal(b, m, deterministic)
}
func (dst *SyncTxReceiptConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTxReceiptConfig.Merge(dst, src)
}
func (m *SyncTxReceiptConfig) XXX_Size() int {
	return xxx_messageInfo_SyncTxReceiptConfig.Size(m)
}
func (m *SyncTxReceiptConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncTxReceiptConfig.DiscardUnknown(m)
}

var xxx_messageInfo_SyncTxReceiptConfig proto.InternalMessageInfo

func (m *SyncTxReceiptConfig) GetChain33Host() string {
	if m != nil {
		return m.Chain33Host
	}
	return ""
}

func (m *SyncTxReceiptConfig) GetPushHost() string {
	if m != nil {
		return m.PushHost
	}
	return ""
}

func (m *SyncTxReceiptConfig) GetPushName() string {
	if m != nil {
		return m.PushName
	}
	return ""
}

func (m *SyncTxReceiptConfig) GetPushBind() string {
	if m != nil {
		return m.PushBind
	}
	return ""
}

func (m *SyncTxReceiptConfig) GetStartSyncHeight() int64 {
	if m != nil {
		return m.StartSyncHeight
	}
	return 0
}

func (m *SyncTxReceiptConfig) GetStartSyncSequence() int64 {
	if m != nil {
		return m.StartSyncSequence
	}
	return 0
}

func (m *SyncTxReceiptConfig) GetStartSyncHash() string {
	if m != nil {
		return m.StartSyncHash
	}
	return ""
}

type Deploy struct {
	// 操作管理员地址
	OperatorAddr string `protobuf:"bytes,1,opt,name=operatorAddr" json:"operatorAddr,omitempty"`
	// 合约部署人员私钥，用于部署合约时签名使用
	DeployerPrivateKey string `protobuf:"bytes,2,opt,name=deployerPrivateKey" json:"deployerPrivateKey,omitempty"`
	// 验证人地址
	ValidatorsAddr []string `protobuf:"bytes,3,rep,name=validatorsAddr" json:"validatorsAddr,omitempty"`
	// 验证人权重
	InitPowers           []int64  `protobuf:"varint,4,rep,packed,name=initPowers" json:"initPowers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Deploy) Reset()         { *m = Deploy{} }
func (m *Deploy) String() string { return proto.CompactTextString(m) }
func (*Deploy) ProtoMessage()    {}
func (*Deploy) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_de1752c779fc3508, []int{4}
}
func (m *Deploy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Deploy.Unmarshal(m, b)
}
func (m *Deploy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Deploy.Marshal(b, m, deterministic)
}
func (dst *Deploy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deploy.Merge(dst, src)
}
func (m *Deploy) XXX_Size() int {
	return xxx_messageInfo_Deploy.Size(m)
}
func (m *Deploy) XXX_DiscardUnknown() {
	xxx_messageInfo_Deploy.DiscardUnknown(m)
}

var xxx_messageInfo_Deploy proto.InternalMessageInfo

func (m *Deploy) GetOperatorAddr() string {
	if m != nil {
		return m.OperatorAddr
	}
	return ""
}

func (m *Deploy) GetDeployerPrivateKey() string {
	if m != nil {
		return m.DeployerPrivateKey
	}
	return ""
}

func (m *Deploy) GetValidatorsAddr() []string {
	if m != nil {
		return m.ValidatorsAddr
	}
	return nil
}

func (m *Deploy) GetInitPowers() []int64 {
	if m != nil {
		return m.InitPowers
	}
	return nil
}

func init() {
	proto.RegisterType((*SyncTxConfig)(nil), "types.SyncTxConfig")
	proto.RegisterType((*Log)(nil), "types.Log")
	proto.RegisterType((*RelayerConfig)(nil), "types.RelayerConfig")
	proto.RegisterType((*SyncTxReceiptConfig)(nil), "types.SyncTxReceiptConfig")
	proto.RegisterType((*Deploy)(nil), "types.Deploy")
}

func init() { proto.RegisterFile("config.proto", fileDescriptor_config_de1752c779fc3508) }

var fileDescriptor_config_de1752c779fc3508 = []byte{
	// 680 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x55, 0xcd, 0x6e, 0x1a, 0x3b,
	0x18, 0x15, 0x4c, 0x20, 0x60, 0xe0, 0x5e, 0x5d, 0xe7, 0xaa, 0x1a, 0x55, 0x51, 0x85, 0x50, 0x5b,
	0xb1, 0xa8, 0x50, 0x95, 0x2c, 0xba, 0xce, 0x8f, 0xa2, 0x48, 0x4d, 0x2a, 0xe4, 0xe4, 0x05, 0x8c,
	0xe7, 0xcb, 0x8c, 0x1b, 0xcf, 0x78, 0x6a, 0x1b, 0xca, 0xf4, 0x0d, 0xfa, 0x1e, 0x5d, 0xf5, 0x2d,
	0xfa, 0x48, 0x7d, 0x83, 0xca, 0xdf, 0x0c, 0x30, 0x10, 0x16, 0xed, 0xaa, 0x3b, 0xce, 0x39, 0xd6,
	0x87, 0x7d, 0xce, 0x77, 0x80, 0xf4, 0x85, 0xce, 0x1e, 0x64, 0x3c, 0xc9, 0x8d, 0x76, 0x9a, 0xb6,
	0x5c, 0x91, 0x83, 0x1d, 0x7d, 0x0f, 0x48, 0xff, 0xae, 0xc8, 0xc4, 0xfd, 0xf2, 0x02, 0x55, 0x3a,
	0x24, 0x3d, 0x91, 0x70, 0x99, 0x9d, 0x9e, 0x26, 0xda, 0xba, 0xb0, 0x31, 0x6c, 0x8c, 0xbb, 0xac,
	0x4e, 0xd1, 0xe7, 0xa4, 0x93, 0xcf, 0x6d, 0x72, 0xed, 0xe5, 0x26, 0xca, 0x6b, 0xbc, 0xd2, 0x3e,
	0xf0, 0x14, 0xc2, 0x60, 0xa3, 0x79, 0xbc, 0xd2, 0xce, 0x65, 0x16, 0x85, 0x07, 0x1b, 0xcd, 0x63,
	0xfa, 0x9a, 0xfc, 0x93, 0x72, 0x37, 0x37, 0xd2, 0x15, 0x97, 0x10, 0x1b, 0x80, 0xb0, 0x35, 0x6c,
	0x8c, 0x5b, 0x6c, 0x87, 0xf5, 0x33, 0xa2, 0x59, 0x64, 0xe4, 0x02, 0x4c, 0xd8, 0x2e, 0x67, 0xac,
	0x30, 0x7d, 0x46, 0xda, 0xd1, 0x6c, 0xca, 0x5d, 0x12, 0x1e, 0xa2, 0x52, 0x21, 0x1a, 0x92, 0xc3,
	0x68, 0x76, 0xc1, 0x45, 0x02, 0x61, 0x07, 0x87, 0xae, 0x20, 0x7d, 0x4b, 0x8e, 0x1e, 0xc0, 0x89,
	0xe4, 0x1a, 0x64, 0x9c, 0xb8, 0x29, 0x18, 0xa9, 0xa3, 0x5b, 0x1b, 0x76, 0x87, 0x8d, 0x71, 0xc0,
	0xf6, 0x49, 0x74, 0x4c, 0xfe, 0xb5, 0x8e, 0x1b, 0xe7, 0x2d, 0x2b, 0xa5, 0x90, 0xe0, 0xe9, 0x5d,
	0x9a, 0xbe, 0x21, 0xff, 0xad, 0xa9, 0x3b, 0xf8, 0x34, 0x87, 0x4c, 0x40, 0xd8, 0xc3, 0xb3, 0x4f,
	0x05, 0xfa, 0x92, 0x0c, 0x36, 0x03, 0xb8, 0x4d, 0xc2, 0x3e, 0x3e, 0x61, 0x9b, 0x1c, 0xfd, 0x68,
	0x92, 0xe0, 0x46, 0xc7, 0xde, 0x05, 0xa5, 0x63, 0x05, 0x0b, 0x50, 0x55, 0x40, 0x6b, 0xec, 0x6f,
	0xa8, 0x74, 0x7c, 0xa1, 0x33, 0xab, 0x15, 0xdc, 0xe0, 0x91, 0x32, 0xa4, 0x5d, 0xda, 0xfb, 0xa2,
	0x74, 0x7c, 0x25, 0xd5, 0x2a, 0xaa, 0x15, 0xf4, 0x3b, 0x90, 0xf2, 0xa5, 0xff, 0x78, 0x27, 0xbf,
	0x00, 0x86, 0x35, 0x60, 0x75, 0x8a, 0xbe, 0x20, 0x24, 0xe5, 0xcb, 0x73, 0x2e, 0x1e, 0xe7, 0xb9,
	0xc5, 0xac, 0x06, 0xac, 0xc6, 0xf8, 0x2c, 0x52, 0xbe, 0x3c, 0x8b, 0x01, 0x53, 0x1a, 0xb0, 0x0a,
	0xd1, 0x63, 0xd2, 0x55, 0x5a, 0x70, 0x75, 0x2f, 0x53, 0xc0, 0x98, 0x3a, 0x6c, 0x43, 0xf8, 0x77,
	0x09, 0x9d, 0xe6, 0x06, 0xac, 0xc5, 0xa8, 0x3a, 0x6c, 0x8d, 0xfd, 0x37, 0x0a, 0xae, 0x14, 0x18,
	0xbc, 0x70, 0x17, 0xd5, 0x1a, 0xe3, 0x37, 0xa8, 0x42, 0xf3, 0x4c, 0x38, 0xa9, 0x33, 0x0c, 0xa6,
	0xc3, 0x76, 0xd8, 0xd1, 0xcf, 0x26, 0x19, 0x30, 0x50, 0xbc, 0x00, 0x53, 0x6d, 0xfc, 0xff, 0xa4,
	0xe5, 0xa4, 0x53, 0x50, 0x59, 0x59, 0x02, 0xfa, 0x8e, 0xf4, 0x6d, 0xad, 0x17, 0x68, 0x62, 0xef,
	0xe4, 0x68, 0x82, 0xb5, 0x99, 0xd4, 0x2b, 0xc3, 0xb6, 0x0e, 0xd2, 0x63, 0x12, 0x28, 0x1d, 0xa3,
	0xa5, 0xbd, 0x13, 0x52, 0x9d, 0xbf, 0xd1, 0x31, 0xf3, 0x34, 0x1d, 0x91, 0xfe, 0x47, 0x93, 0x0b,
	0xbf, 0xf4, 0x67, 0x51, 0x64, 0xaa, 0x22, 0x6c, 0x71, 0xde, 0x7e, 0x70, 0xc9, 0xd4, 0xe8, 0x85,
	0x8c, 0xc0, 0xa0, 0xbb, 0x5d, 0x56, 0xa7, 0xfc, 0x63, 0x67, 0x46, 0x46, 0x31, 0x30, 0x88, 0xa5,
	0x75, 0xa6, 0xa8, 0xca, 0xb0, 0xc3, 0xd2, 0x57, 0xa4, 0x1d, 0x41, 0xae, 0x74, 0x81, 0x5e, 0xf7,
	0x4e, 0x06, 0xd5, 0x75, 0x2e, 0x91, 0x64, 0x95, 0xe8, 0x77, 0x15, 0x5c, 0x72, 0xbb, 0x5d, 0xc0,
	0xb2, 0x2b, 0x4f, 0x05, 0xdf, 0x1a, 0x70, 0xc9, 0xb9, 0xd2, 0xe2, 0xf1, 0xca, 0x57, 0xa4, 0x2c,
	0x07, 0x46, 0xd2, 0x62, 0xfb, 0xa4, 0xd1, 0xd7, 0x26, 0x39, 0x2a, 0x1d, 0x63, 0x20, 0x40, 0xe6,
	0xee, 0xaf, 0xfe, 0xd6, 0xec, 0xe9, 0x70, 0xeb, 0x0f, 0x3a, 0xdc, 0xfe, 0xed, 0x0e, 0x1f, 0xee,
	0xeb, 0xf0, 0xb7, 0x06, 0x69, 0x97, 0xf6, 0xfb, 0x5d, 0xd0, 0x39, 0x18, 0xee, 0xb4, 0xc1, 0x5d,
	0x28, 0xdf, 0xbf, 0xc5, 0xd1, 0x09, 0xa1, 0x65, 0x48, 0x60, 0xa6, 0x46, 0x2e, 0xb8, 0x83, 0xf7,
	0x50, 0x54, 0x56, 0xec, 0x51, 0xfc, 0x66, 0x2c, 0xb8, 0x92, 0x91, 0x1f, 0x60, 0x71, 0x6a, 0x30,
	0x0c, 0xfc, 0x66, 0x6c, 0xb3, 0xbe, 0x4e, 0x32, 0x93, 0x6e, 0xaa, 0x3f, 0x83, 0xb1, 0xe1, 0xc1,
	0x30, 0x18, 0x07, 0xac, 0xc6, 0xcc, 0xda, 0xf8, 0x2f, 0x71, 0xfa, 0x2b, 0x00, 0x00, 0xff, 0xff,
	0x09, 0xa3, 0xd6, 0xaa, 0x35, 0x06, 0x00, 0x00,
}