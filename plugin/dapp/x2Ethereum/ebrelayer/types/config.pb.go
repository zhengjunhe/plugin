// Code generated by protoc-gen-go. DO NOT EDIT.
// source: config.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SyncTxConfig struct {
	Chain33Host          string   `protobuf:"bytes,1,opt,name=chain33host,proto3" json:"chain33host,omitempty"`
	PushHost             string   `protobuf:"bytes,2,opt,name=pushHost,proto3" json:"pushHost,omitempty"`
	PushName             string   `protobuf:"bytes,3,opt,name=pushName,proto3" json:"pushName,omitempty"`
	PushBind             string   `protobuf:"bytes,4,opt,name=pushBind,proto3" json:"pushBind,omitempty"`
	MaturityDegree       int32    `protobuf:"varint,5,opt,name=maturityDegree,proto3" json:"maturityDegree,omitempty"`
	Dbdriver             string   `protobuf:"bytes,6,opt,name=dbdriver,proto3" json:"dbdriver,omitempty"`
	DbPath               string   `protobuf:"bytes,7,opt,name=dbPath,proto3" json:"dbPath,omitempty"`
	DbCache              int32    `protobuf:"varint,8,opt,name=dbCache,proto3" json:"dbCache,omitempty"`
	FetchHeightPeriodMs  int64    `protobuf:"varint,9,opt,name=fetchHeightPeriodMs,proto3" json:"fetchHeightPeriodMs,omitempty"`
	StartSyncHeight      int64    `protobuf:"varint,10,opt,name=startSyncHeight,proto3" json:"startSyncHeight,omitempty"`
	StartSyncSequence    int64    `protobuf:"varint,11,opt,name=startSyncSequence,proto3" json:"startSyncSequence,omitempty"`
	StartSyncHash        string   `protobuf:"bytes,12,opt,name=startSyncHash,proto3" json:"startSyncHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncTxConfig) Reset()         { *m = SyncTxConfig{} }
func (m *SyncTxConfig) String() string { return proto.CompactTextString(m) }
func (*SyncTxConfig) ProtoMessage()    {}
func (*SyncTxConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{0}
}

func (m *SyncTxConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTxConfig.Unmarshal(m, b)
}
func (m *SyncTxConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTxConfig.Marshal(b, m, deterministic)
}
func (m *SyncTxConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTxConfig.Merge(m, src)
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
	Loglevel             string   `protobuf:"bytes,1,opt,name=loglevel,proto3" json:"loglevel,omitempty"`
	LogConsoleLevel      string   `protobuf:"bytes,2,opt,name=logConsoleLevel,proto3" json:"logConsoleLevel,omitempty"`
	LogFile              string   `protobuf:"bytes,3,opt,name=logFile,proto3" json:"logFile,omitempty"`
	MaxFileSize          uint32   `protobuf:"varint,4,opt,name=maxFileSize,proto3" json:"maxFileSize,omitempty"`
	MaxBackups           uint32   `protobuf:"varint,5,opt,name=maxBackups,proto3" json:"maxBackups,omitempty"`
	MaxAge               uint32   `protobuf:"varint,6,opt,name=maxAge,proto3" json:"maxAge,omitempty"`
	LocalTime            bool     `protobuf:"varint,7,opt,name=localTime,proto3" json:"localTime,omitempty"`
	Compress             bool     `protobuf:"varint,8,opt,name=compress,proto3" json:"compress,omitempty"`
	CallerFile           bool     `protobuf:"varint,9,opt,name=callerFile,proto3" json:"callerFile,omitempty"`
	CallerFunction       bool     `protobuf:"varint,10,opt,name=callerFunction,proto3" json:"callerFunction,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Log) Reset()         { *m = Log{} }
func (m *Log) String() string { return proto.CompactTextString(m) }
func (*Log) ProtoMessage()    {}
func (*Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{1}
}

func (m *Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Log.Unmarshal(m, b)
}
func (m *Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Log.Marshal(b, m, deterministic)
}
func (m *Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Log.Merge(m, src)
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
	Title                string        `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	SyncTxConfig         *SyncTxConfig `protobuf:"bytes,2,opt,name=syncTxConfig,proto3" json:"syncTxConfig,omitempty"`
	Log                  *Log          `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	JrpcBindAddr         string        `protobuf:"bytes,4,opt,name=jrpcBindAddr,proto3" json:"jrpcBindAddr,omitempty"`
	EthProvider          string        `protobuf:"bytes,5,opt,name=ethProvider,proto3" json:"ethProvider,omitempty"`
	BridgeRegistry       string        `protobuf:"bytes,6,opt,name=bridgeRegistry,proto3" json:"bridgeRegistry,omitempty"`
	Deploy               *Deploy       `protobuf:"bytes,7,opt,name=Deploy,proto3" json:"Deploy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RelayerConfig) Reset()         { *m = RelayerConfig{} }
func (m *RelayerConfig) String() string { return proto.CompactTextString(m) }
func (*RelayerConfig) ProtoMessage()    {}
func (*RelayerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{2}
}

func (m *RelayerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayerConfig.Unmarshal(m, b)
}
func (m *RelayerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayerConfig.Marshal(b, m, deterministic)
}
func (m *RelayerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayerConfig.Merge(m, src)
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

type SyncTxReceiptConfig struct {
	Chain33Host          string   `protobuf:"bytes,1,opt,name=chain33host,proto3" json:"chain33host,omitempty"`
	PushHost             string   `protobuf:"bytes,2,opt,name=pushHost,proto3" json:"pushHost,omitempty"`
	PushName             string   `protobuf:"bytes,3,opt,name=pushName,proto3" json:"pushName,omitempty"`
	PushBind             string   `protobuf:"bytes,4,opt,name=pushBind,proto3" json:"pushBind,omitempty"`
	StartSyncHeight      int64    `protobuf:"varint,5,opt,name=startSyncHeight,proto3" json:"startSyncHeight,omitempty"`
	StartSyncSequence    int64    `protobuf:"varint,6,opt,name=startSyncSequence,proto3" json:"startSyncSequence,omitempty"`
	StartSyncHash        string   `protobuf:"bytes,7,opt,name=startSyncHash,proto3" json:"startSyncHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncTxReceiptConfig) Reset()         { *m = SyncTxReceiptConfig{} }
func (m *SyncTxReceiptConfig) String() string { return proto.CompactTextString(m) }
func (*SyncTxReceiptConfig) ProtoMessage()    {}
func (*SyncTxReceiptConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{3}
}

func (m *SyncTxReceiptConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncTxReceiptConfig.Unmarshal(m, b)
}
func (m *SyncTxReceiptConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncTxReceiptConfig.Marshal(b, m, deterministic)
}
func (m *SyncTxReceiptConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncTxReceiptConfig.Merge(m, src)
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
	//操作管理员地址
	OperatorAddr string `protobuf:"bytes,1,opt,name=operatorAddr,proto3" json:"operatorAddr,omitempty"`
	//合约部署人员私钥，用于部署合约时签名使用
	DeployerPrivateKey string `protobuf:"bytes,2,opt,name=deployerPrivateKey,proto3" json:"deployerPrivateKey,omitempty"`
	//验证人地址
	ValidatorsAddr []string `protobuf:"bytes,3,rep,name=validatorsAddr,proto3" json:"validatorsAddr,omitempty"`
	//验证人权重
	InitPowers           []int64  `protobuf:"varint,4,rep,packed,name=initPowers,proto3" json:"initPowers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Deploy) Reset()         { *m = Deploy{} }
func (m *Deploy) String() string { return proto.CompactTextString(m) }
func (*Deploy) ProtoMessage()    {}
func (*Deploy) Descriptor() ([]byte, []int) {
	return fileDescriptor_3eaf2c85e69e9ea4, []int{4}
}

func (m *Deploy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Deploy.Unmarshal(m, b)
}
func (m *Deploy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Deploy.Marshal(b, m, deterministic)
}
func (m *Deploy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deploy.Merge(m, src)
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

func init() {
	proto.RegisterFile("config.proto", fileDescriptor_3eaf2c85e69e9ea4)
}

var fileDescriptor_3eaf2c85e69e9ea4 = []byte{
	// 645 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x54, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0x56, 0x9b, 0x6d, 0xb7, 0x75, 0x5b, 0x10, 0x5e, 0x84, 0x2c, 0xb4, 0x42, 0x51, 0x05, 0xa8,
	0x07, 0x54, 0xa1, 0xdd, 0x03, 0xe7, 0xfd, 0x11, 0x5a, 0x89, 0x05, 0x55, 0xee, 0xbe, 0x80, 0xeb,
	0xcc, 0x26, 0x06, 0x27, 0x0e, 0xb6, 0x5b, 0x1a, 0xde, 0x80, 0x23, 0xef, 0xc0, 0x89, 0xb7, 0xe0,
	0xcd, 0x90, 0x9d, 0xa4, 0x4d, 0x4b, 0x0f, 0x70, 0xe2, 0xd6, 0xef, 0xfb, 0x46, 0x53, 0xcf, 0xf7,
	0xcd, 0x04, 0x0d, 0xb9, 0xca, 0xee, 0x45, 0x3c, 0xcd, 0xb5, 0xb2, 0x0a, 0x77, 0x6c, 0x91, 0x83,
	0x19, 0xff, 0x0c, 0xd0, 0x70, 0x5e, 0x64, 0xfc, 0x6e, 0x7d, 0xe5, 0x55, 0x1c, 0xa2, 0x01, 0x4f,
	0x98, 0xc8, 0xce, 0xcf, 0x13, 0x65, 0x2c, 0x69, 0x85, 0xad, 0x49, 0x9f, 0x36, 0x29, 0xfc, 0x14,
	0xf5, 0xf2, 0xa5, 0x49, 0x6e, 0x9c, 0xdc, 0xf6, 0xf2, 0x06, 0xd7, 0xda, 0x07, 0x96, 0x02, 0x09,
	0xb6, 0x9a, 0xc3, 0xb5, 0x76, 0x29, 0xb2, 0x88, 0x1c, 0x6d, 0x35, 0x87, 0xf1, 0x4b, 0xf4, 0x20,
	0x65, 0x76, 0xa9, 0x85, 0x2d, 0xae, 0x21, 0xd6, 0x00, 0xa4, 0x13, 0xb6, 0x26, 0x1d, 0xba, 0xc7,
	0xba, 0x1e, 0xd1, 0x22, 0xd2, 0x62, 0x05, 0x9a, 0x74, 0xcb, 0x1e, 0x35, 0xc6, 0x4f, 0x50, 0x37,
	0x5a, 0xcc, 0x98, 0x4d, 0xc8, 0xb1, 0x57, 0x2a, 0x84, 0x09, 0x3a, 0x8e, 0x16, 0x57, 0x8c, 0x27,
	0x40, 0x7a, 0xbe, 0x69, 0x0d, 0xf1, 0x6b, 0x74, 0x72, 0x0f, 0x96, 0x27, 0x37, 0x20, 0xe2, 0xc4,
	0xce, 0x40, 0x0b, 0x15, 0xbd, 0x37, 0xa4, 0x1f, 0xb6, 0x26, 0x01, 0x3d, 0x24, 0xe1, 0x09, 0x7a,
	0x68, 0x2c, 0xd3, 0xd6, 0x59, 0x56, 0x4a, 0x04, 0xf9, 0xea, 0x7d, 0x1a, 0xbf, 0x42, 0x8f, 0x36,
	0xd4, 0x1c, 0x3e, 0x2f, 0x21, 0xe3, 0x40, 0x06, 0xbe, 0xf6, 0x4f, 0x01, 0x3f, 0x47, 0xa3, 0x6d,
	0x03, 0x66, 0x12, 0x32, 0xf4, 0x23, 0xec, 0x92, 0xe3, 0x5f, 0x6d, 0x14, 0xdc, 0xaa, 0xd8, 0xb9,
	0x20, 0x55, 0x2c, 0x61, 0x05, 0xb2, 0x0a, 0x68, 0x83, 0xdd, 0x0b, 0xa5, 0x8a, 0xaf, 0x54, 0x66,
	0x94, 0x84, 0x5b, 0x5f, 0x52, 0x86, 0xb4, 0x4f, 0x3b, 0x5f, 0xa4, 0x8a, 0xdf, 0x0a, 0x59, 0x47,
	0x55, 0x43, 0xb7, 0x03, 0x29, 0x5b, 0xbb, 0x9f, 0x73, 0xf1, 0x15, 0x7c, 0x58, 0x23, 0xda, 0xa4,
	0xf0, 0x33, 0x84, 0x52, 0xb6, 0xbe, 0x64, 0xfc, 0xd3, 0x32, 0x37, 0x3e, 0xab, 0x11, 0x6d, 0x30,
	0x2e, 0x8b, 0x94, 0xad, 0x2f, 0x62, 0xf0, 0x29, 0x8d, 0x68, 0x85, 0xf0, 0x29, 0xea, 0x4b, 0xc5,
	0x99, 0xbc, 0x13, 0x29, 0xf8, 0x98, 0x7a, 0x74, 0x4b, 0xb8, 0xb9, 0xb8, 0x4a, 0x73, 0x0d, 0xc6,
	0xf8, 0xa8, 0x7a, 0x74, 0x83, 0xdd, 0x3f, 0x72, 0x26, 0x25, 0x68, 0xff, 0xe0, 0xbe, 0x57, 0x1b,
	0x8c, 0xdb, 0xa0, 0x0a, 0x2d, 0x33, 0x6e, 0x85, 0xca, 0x7c, 0x30, 0x3d, 0xba, 0xc7, 0x8e, 0xbf,
	0xb7, 0xd1, 0x88, 0x82, 0x64, 0x05, 0xe8, 0x6a, 0xe3, 0x1f, 0xa3, 0x8e, 0x15, 0x56, 0x42, 0x65,
	0x65, 0x09, 0xf0, 0x1b, 0x34, 0x34, 0x8d, 0xbb, 0xf0, 0x26, 0x0e, 0xce, 0x4e, 0xa6, 0xfe, 0x6c,
	0xa6, 0xcd, 0x93, 0xa1, 0x3b, 0x85, 0xf8, 0x14, 0x05, 0x52, 0xc5, 0xde, 0xd2, 0xc1, 0x19, 0xaa,
	0xea, 0x6f, 0x55, 0x4c, 0x1d, 0x8d, 0xc7, 0x68, 0xf8, 0x51, 0xe7, 0xdc, 0x2d, 0xfd, 0x45, 0x14,
	0xe9, 0xea, 0x10, 0x76, 0x38, 0x67, 0x3f, 0xd8, 0x64, 0xa6, 0xd5, 0x4a, 0x44, 0xa0, 0xbd, 0xbb,
	0x7d, 0xda, 0xa4, 0xdc, 0xb0, 0x0b, 0x2d, 0xa2, 0x18, 0x28, 0xc4, 0xc2, 0x58, 0x5d, 0x54, 0xc7,
	0xb0, 0xc7, 0xe2, 0x17, 0xa8, 0x7b, 0x0d, 0xb9, 0x54, 0x85, 0xf7, 0x7a, 0x70, 0x36, 0xaa, 0x9e,
	0x53, 0x92, 0xb4, 0x12, 0xc7, 0xdf, 0xda, 0xe8, 0xa4, 0x9c, 0x88, 0x02, 0x07, 0x91, 0xdb, 0xff,
	0xfa, 0x2d, 0x38, 0x70, 0x63, 0x9d, 0x7f, 0xb8, 0xb1, 0xee, 0x5f, 0xdf, 0xd8, 0xf1, 0xa1, 0x1b,
	0xfb, 0xd1, 0xaa, 0x3d, 0x73, 0x59, 0xa9, 0x1c, 0x34, 0xb3, 0x4a, 0xfb, 0xac, 0xca, 0xf9, 0x77,
	0x38, 0x3c, 0x45, 0x38, 0xf2, 0xd5, 0xa0, 0x67, 0x5a, 0xac, 0x98, 0x85, 0x77, 0x50, 0x54, 0x56,
	0x1c, 0x50, 0x5c, 0x72, 0x2b, 0x26, 0x45, 0xe4, 0x1a, 0x18, 0xdf, 0x35, 0x08, 0x03, 0x97, 0xdc,
	0x2e, 0xeb, 0xd6, 0x5d, 0x64, 0xc2, 0xce, 0xd4, 0x17, 0xd0, 0x86, 0x1c, 0x85, 0xc1, 0x24, 0xa0,
	0x0d, 0x66, 0xd1, 0xf5, 0x5f, 0xf1, 0xf3, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x07, 0x89, 0x22,
	0xf1, 0xd5, 0x05, 0x00, 0x00,
}
