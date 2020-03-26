// Code generated by protoc-gen-go. DO NOT EDIT.
// source: relayer.proto

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

//以太坊账户信息
// 	 privkey : 账户地址对应的私钥
//	 addr :账户地址
type Account4Relayer struct {
	Privkey              []byte   `protobuf:"bytes,1,opt,name=privkey,proto3" json:"privkey,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account4Relayer) Reset()         { *m = Account4Relayer{} }
func (m *Account4Relayer) String() string { return proto.CompactTextString(m) }
func (*Account4Relayer) ProtoMessage()    {}
func (*Account4Relayer) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{0}
}

func (m *Account4Relayer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account4Relayer.Unmarshal(m, b)
}
func (m *Account4Relayer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account4Relayer.Marshal(b, m, deterministic)
}
func (m *Account4Relayer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account4Relayer.Merge(m, src)
}
func (m *Account4Relayer) XXX_Size() int {
	return xxx_messageInfo_Account4Relayer.Size(m)
}
func (m *Account4Relayer) XXX_DiscardUnknown() {
	xxx_messageInfo_Account4Relayer.DiscardUnknown(m)
}

var xxx_messageInfo_Account4Relayer proto.InternalMessageInfo

func (m *Account4Relayer) GetPrivkey() []byte {
	if m != nil {
		return m.Privkey
	}
	return nil
}

func (m *Account4Relayer) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type ValidatorAddr4EthRelayer struct {
	EthValidator         string   `protobuf:"bytes,1,opt,name=ethValidator,proto3" json:"ethValidator,omitempty"`
	Chain33Validator     string   `protobuf:"bytes,2,opt,name=chain33Validator,proto3" json:"chain33Validator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidatorAddr4EthRelayer) Reset()         { *m = ValidatorAddr4EthRelayer{} }
func (m *ValidatorAddr4EthRelayer) String() string { return proto.CompactTextString(m) }
func (*ValidatorAddr4EthRelayer) ProtoMessage()    {}
func (*ValidatorAddr4EthRelayer) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{1}
}

func (m *ValidatorAddr4EthRelayer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidatorAddr4EthRelayer.Unmarshal(m, b)
}
func (m *ValidatorAddr4EthRelayer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidatorAddr4EthRelayer.Marshal(b, m, deterministic)
}
func (m *ValidatorAddr4EthRelayer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidatorAddr4EthRelayer.Merge(m, src)
}
func (m *ValidatorAddr4EthRelayer) XXX_Size() int {
	return xxx_messageInfo_ValidatorAddr4EthRelayer.Size(m)
}
func (m *ValidatorAddr4EthRelayer) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidatorAddr4EthRelayer.DiscardUnknown(m)
}

var xxx_messageInfo_ValidatorAddr4EthRelayer proto.InternalMessageInfo

func (m *ValidatorAddr4EthRelayer) GetEthValidator() string {
	if m != nil {
		return m.EthValidator
	}
	return ""
}

func (m *ValidatorAddr4EthRelayer) GetChain33Validator() string {
	if m != nil {
		return m.Chain33Validator
	}
	return ""
}

type Txhashes struct {
	Txhash               []string `protobuf:"bytes,1,rep,name=txhash,proto3" json:"txhash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Txhashes) Reset()         { *m = Txhashes{} }
func (m *Txhashes) String() string { return proto.CompactTextString(m) }
func (*Txhashes) ProtoMessage()    {}
func (*Txhashes) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{2}
}

func (m *Txhashes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Txhashes.Unmarshal(m, b)
}
func (m *Txhashes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Txhashes.Marshal(b, m, deterministic)
}
func (m *Txhashes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Txhashes.Merge(m, src)
}
func (m *Txhashes) XXX_Size() int {
	return xxx_messageInfo_Txhashes.Size(m)
}
func (m *Txhashes) XXX_DiscardUnknown() {
	xxx_messageInfo_Txhashes.DiscardUnknown(m)
}

var xxx_messageInfo_Txhashes proto.InternalMessageInfo

func (m *Txhashes) GetTxhash() []string {
	if m != nil {
		return m.Txhash
	}
	return nil
}

type ReqSetPasswd struct {
	OldPassphase         string   `protobuf:"bytes,1,opt,name=oldPassphase,proto3" json:"oldPassphase,omitempty"`
	NewPassphase         string   `protobuf:"bytes,2,opt,name=newPassphase,proto3" json:"newPassphase,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReqSetPasswd) Reset()         { *m = ReqSetPasswd{} }
func (m *ReqSetPasswd) String() string { return proto.CompactTextString(m) }
func (*ReqSetPasswd) ProtoMessage()    {}
func (*ReqSetPasswd) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{3}
}

func (m *ReqSetPasswd) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReqSetPasswd.Unmarshal(m, b)
}
func (m *ReqSetPasswd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReqSetPasswd.Marshal(b, m, deterministic)
}
func (m *ReqSetPasswd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReqSetPasswd.Merge(m, src)
}
func (m *ReqSetPasswd) XXX_Size() int {
	return xxx_messageInfo_ReqSetPasswd.Size(m)
}
func (m *ReqSetPasswd) XXX_DiscardUnknown() {
	xxx_messageInfo_ReqSetPasswd.DiscardUnknown(m)
}

var xxx_messageInfo_ReqSetPasswd proto.InternalMessageInfo

func (m *ReqSetPasswd) GetOldPassphase() string {
	if m != nil {
		return m.OldPassphase
	}
	return ""
}

func (m *ReqSetPasswd) GetNewPassphase() string {
	if m != nil {
		return m.NewPassphase
	}
	return ""
}

type Account4Show struct {
	Privkey              string   `protobuf:"bytes,1,opt,name=privkey,proto3" json:"privkey,omitempty"`
	Addr                 string   `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account4Show) Reset()         { *m = Account4Show{} }
func (m *Account4Show) String() string { return proto.CompactTextString(m) }
func (*Account4Show) ProtoMessage()    {}
func (*Account4Show) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{4}
}

func (m *Account4Show) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account4Show.Unmarshal(m, b)
}
func (m *Account4Show) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account4Show.Marshal(b, m, deterministic)
}
func (m *Account4Show) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account4Show.Merge(m, src)
}
func (m *Account4Show) XXX_Size() int {
	return xxx_messageInfo_Account4Show.Size(m)
}
func (m *Account4Show) XXX_DiscardUnknown() {
	xxx_messageInfo_Account4Show.DiscardUnknown(m)
}

var xxx_messageInfo_Account4Show proto.InternalMessageInfo

func (m *Account4Show) GetPrivkey() string {
	if m != nil {
		return m.Privkey
	}
	return ""
}

func (m *Account4Show) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

type AssetType struct {
	Chain                string   `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	IssueContract        string   `protobuf:"bytes,2,opt,name=issueContract,proto3" json:"issueContract,omitempty"`
	Symbol               string   `protobuf:"bytes,3,opt,name=symbol,proto3" json:"symbol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AssetType) Reset()         { *m = AssetType{} }
func (m *AssetType) String() string { return proto.CompactTextString(m) }
func (*AssetType) ProtoMessage()    {}
func (*AssetType) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{5}
}

func (m *AssetType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AssetType.Unmarshal(m, b)
}
func (m *AssetType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AssetType.Marshal(b, m, deterministic)
}
func (m *AssetType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssetType.Merge(m, src)
}
func (m *AssetType) XXX_Size() int {
	return xxx_messageInfo_AssetType.Size(m)
}
func (m *AssetType) XXX_DiscardUnknown() {
	xxx_messageInfo_AssetType.DiscardUnknown(m)
}

var xxx_messageInfo_AssetType proto.InternalMessageInfo

func (m *AssetType) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

func (m *AssetType) GetIssueContract() string {
	if m != nil {
		return m.IssueContract
	}
	return ""
}

func (m *AssetType) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

type EthBridgeClaim struct {
	EthereumChainID       int32    `protobuf:"varint,1,opt,name=ethereumChainID,proto3" json:"ethereumChainID,omitempty"`
	BridgeContractAddress []byte   `protobuf:"bytes,2,opt,name=bridgeContractAddress,proto3" json:"bridgeContractAddress,omitempty"`
	Nonce                 int64    `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	TokenContractAddress  []byte   `protobuf:"bytes,4,opt,name=tokenContractAddress,proto3" json:"tokenContractAddress,omitempty"`
	Symbol                string   `protobuf:"bytes,5,opt,name=symbol,proto3" json:"symbol,omitempty"`
	EthereumSender        []byte   `protobuf:"bytes,6,opt,name=ethereumSender,proto3" json:"ethereumSender,omitempty"`
	Chain33Receiver       []byte   `protobuf:"bytes,7,opt,name=chain33Receiver,proto3" json:"chain33Receiver,omitempty"`
	ValidatorAddress      []byte   `protobuf:"bytes,8,opt,name=ValidatorAddress,proto3" json:"ValidatorAddress,omitempty"`
	Amount                int64    `protobuf:"varint,9,opt,name=amount,proto3" json:"amount,omitempty"`
	ClaimType             int32    `protobuf:"varint,10,opt,name=claimType,proto3" json:"claimType,omitempty"`
	ChainName             string   `protobuf:"bytes,11,opt,name=chainName,proto3" json:"chainName,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *EthBridgeClaim) Reset()         { *m = EthBridgeClaim{} }
func (m *EthBridgeClaim) String() string { return proto.CompactTextString(m) }
func (*EthBridgeClaim) ProtoMessage()    {}
func (*EthBridgeClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{6}
}

func (m *EthBridgeClaim) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EthBridgeClaim.Unmarshal(m, b)
}
func (m *EthBridgeClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EthBridgeClaim.Marshal(b, m, deterministic)
}
func (m *EthBridgeClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EthBridgeClaim.Merge(m, src)
}
func (m *EthBridgeClaim) XXX_Size() int {
	return xxx_messageInfo_EthBridgeClaim.Size(m)
}
func (m *EthBridgeClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_EthBridgeClaim.DiscardUnknown(m)
}

var xxx_messageInfo_EthBridgeClaim proto.InternalMessageInfo

func (m *EthBridgeClaim) GetEthereumChainID() int32 {
	if m != nil {
		return m.EthereumChainID
	}
	return 0
}

func (m *EthBridgeClaim) GetBridgeContractAddress() []byte {
	if m != nil {
		return m.BridgeContractAddress
	}
	return nil
}

func (m *EthBridgeClaim) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *EthBridgeClaim) GetTokenContractAddress() []byte {
	if m != nil {
		return m.TokenContractAddress
	}
	return nil
}

func (m *EthBridgeClaim) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *EthBridgeClaim) GetEthereumSender() []byte {
	if m != nil {
		return m.EthereumSender
	}
	return nil
}

func (m *EthBridgeClaim) GetChain33Receiver() []byte {
	if m != nil {
		return m.Chain33Receiver
	}
	return nil
}

func (m *EthBridgeClaim) GetValidatorAddress() []byte {
	if m != nil {
		return m.ValidatorAddress
	}
	return nil
}

func (m *EthBridgeClaim) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *EthBridgeClaim) GetClaimType() int32 {
	if m != nil {
		return m.ClaimType
	}
	return 0
}

func (m *EthBridgeClaim) GetChainName() string {
	if m != nil {
		return m.ChainName
	}
	return ""
}

type ImportKeyReq struct {
	PrivateKey           string   `protobuf:"bytes,1,opt,name=privateKey,proto3" json:"privateKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImportKeyReq) Reset()         { *m = ImportKeyReq{} }
func (m *ImportKeyReq) String() string { return proto.CompactTextString(m) }
func (*ImportKeyReq) ProtoMessage()    {}
func (*ImportKeyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{7}
}

func (m *ImportKeyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImportKeyReq.Unmarshal(m, b)
}
func (m *ImportKeyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImportKeyReq.Marshal(b, m, deterministic)
}
func (m *ImportKeyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImportKeyReq.Merge(m, src)
}
func (m *ImportKeyReq) XXX_Size() int {
	return xxx_messageInfo_ImportKeyReq.Size(m)
}
func (m *ImportKeyReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ImportKeyReq.DiscardUnknown(m)
}

var xxx_messageInfo_ImportKeyReq proto.InternalMessageInfo

func (m *ImportKeyReq) GetPrivateKey() string {
	if m != nil {
		return m.PrivateKey
	}
	return ""
}

type RelayerRunStatus struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Details              string   `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelayerRunStatus) Reset()         { *m = RelayerRunStatus{} }
func (m *RelayerRunStatus) String() string { return proto.CompactTextString(m) }
func (*RelayerRunStatus) ProtoMessage()    {}
func (*RelayerRunStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_202a89775a80bd4c, []int{8}
}

func (m *RelayerRunStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayerRunStatus.Unmarshal(m, b)
}
func (m *RelayerRunStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayerRunStatus.Marshal(b, m, deterministic)
}
func (m *RelayerRunStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayerRunStatus.Merge(m, src)
}
func (m *RelayerRunStatus) XXX_Size() int {
	return xxx_messageInfo_RelayerRunStatus.Size(m)
}
func (m *RelayerRunStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayerRunStatus.DiscardUnknown(m)
}

var xxx_messageInfo_RelayerRunStatus proto.InternalMessageInfo

func (m *RelayerRunStatus) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *RelayerRunStatus) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func init() {
	proto.RegisterType((*Account4Relayer)(nil), "types.Account4Relayer")
	proto.RegisterType((*ValidatorAddr4EthRelayer)(nil), "types.ValidatorAddr4EthRelayer")
	proto.RegisterType((*Txhashes)(nil), "types.Txhashes")
	proto.RegisterType((*ReqSetPasswd)(nil), "types.ReqSetPasswd")
	proto.RegisterType((*Account4Show)(nil), "types.Account4Show")
	proto.RegisterType((*AssetType)(nil), "types.assetType")
	proto.RegisterType((*EthBridgeClaim)(nil), "types.EthBridgeClaim")
	proto.RegisterType((*ImportKeyReq)(nil), "types.ImportKeyReq")
	proto.RegisterType((*RelayerRunStatus)(nil), "types.RelayerRunStatus")
}

func init() {
	proto.RegisterFile("relayer.proto", fileDescriptor_202a89775a80bd4c)
}

var fileDescriptor_202a89775a80bd4c = []byte{
	// 500 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x53, 0xdf, 0x6b, 0xdb, 0x30,
	0x10, 0x26, 0x4b, 0x93, 0xd6, 0x37, 0xf7, 0x07, 0xa2, 0x1b, 0x7a, 0x18, 0x23, 0x88, 0x31, 0xc2,
	0x1e, 0xfa, 0xb0, 0xe4, 0x71, 0x30, 0xba, 0xb4, 0x0f, 0xa5, 0x30, 0x86, 0x52, 0xfa, 0x3a, 0x14,
	0xeb, 0x98, 0xbd, 0xda, 0x96, 0x2b, 0xc9, 0xcd, 0xfc, 0xff, 0xec, 0x0f, 0x1d, 0x92, 0xe5, 0x26,
	0x71, 0xc3, 0xde, 0xf4, 0x7d, 0x77, 0xba, 0xbb, 0xef, 0xf4, 0x09, 0x8e, 0x35, 0xe6, 0xa2, 0x41,
	0x7d, 0x51, 0x69, 0x65, 0x15, 0x19, 0xd9, 0xa6, 0x42, 0xc3, 0xbe, 0xc2, 0xe9, 0x65, 0x92, 0xa8,
	0xba, 0xb4, 0x73, 0xde, 0xc6, 0x09, 0x85, 0xc3, 0x4a, 0x67, 0x4f, 0x0f, 0xd8, 0xd0, 0xc1, 0x64,
	0x30, 0x8d, 0x79, 0x07, 0x09, 0x81, 0x03, 0x21, 0xa5, 0xa6, 0xaf, 0x26, 0x83, 0x69, 0xc4, 0xfd,
	0x99, 0xfd, 0x06, 0x7a, 0x2f, 0xf2, 0x4c, 0x0a, 0xab, 0xf4, 0xa5, 0x94, 0x7a, 0x7e, 0x6d, 0xd3,
	0xae, 0x12, 0x83, 0x18, 0x6d, 0xfa, 0x1c, 0xf6, 0xe5, 0x22, 0xbe, 0xc3, 0x91, 0x4f, 0x70, 0x96,
	0xa4, 0x22, 0x2b, 0x67, 0xb3, 0x4d, 0x5e, 0x5b, 0xff, 0x05, 0xcf, 0x18, 0x1c, 0xdd, 0xfd, 0x49,
	0x85, 0x49, 0xd1, 0x90, 0xb7, 0x30, 0xb6, 0xfe, 0x4c, 0x07, 0x93, 0xe1, 0x34, 0xe2, 0x01, 0xb1,
	0x7b, 0x88, 0x39, 0x3e, 0x2e, 0xd1, 0xfe, 0x10, 0xc6, 0xac, 0xa5, 0x9b, 0x41, 0xe5, 0xd2, 0x81,
	0x2a, 0x15, 0x06, 0xbb, 0x19, 0xb6, 0x39, 0x97, 0x53, 0xe2, 0x7a, 0x93, 0xd3, 0xf6, 0xdf, 0xe1,
	0xd8, 0x17, 0x88, 0xbb, 0x45, 0x2d, 0x53, 0xb5, 0xee, 0x6f, 0x29, 0xfa, 0xff, 0x96, 0x7e, 0x42,
	0x24, 0x8c, 0x41, 0x7b, 0xd7, 0x54, 0x48, 0xce, 0x61, 0xe4, 0xa5, 0x85, 0x8b, 0x2d, 0x20, 0x1f,
	0xe0, 0x38, 0x33, 0xa6, 0xc6, 0x85, 0x2a, 0xad, 0x16, 0x89, 0x0d, 0xf7, 0x77, 0x49, 0x27, 0xdb,
	0x34, 0xc5, 0x4a, 0xe5, 0x74, 0xe8, 0xc3, 0x01, 0xb1, 0xbf, 0x43, 0x38, 0xb9, 0xb6, 0xe9, 0x37,
	0x9d, 0xc9, 0x5f, 0xb8, 0xc8, 0x45, 0x56, 0x90, 0x29, 0x9c, 0xa2, 0x4d, 0x51, 0x63, 0x5d, 0x2c,
	0x5c, 0x87, 0x9b, 0x2b, 0xdf, 0x70, 0xc4, 0xfb, 0x34, 0x99, 0xc3, 0x9b, 0x55, 0x7b, 0x31, 0xb4,
	0x71, 0x0f, 0x89, 0xc6, 0xf8, 0x11, 0x62, 0xbe, 0x3f, 0xe8, 0x64, 0x94, 0xaa, 0x4c, 0xd0, 0x4f,
	0x32, 0xe4, 0x2d, 0x20, 0x9f, 0xe1, 0xdc, 0xaa, 0x07, 0x2c, 0xfb, 0xa5, 0x0e, 0x7c, 0xa9, 0xbd,
	0xb1, 0x2d, 0x51, 0xa3, 0x6d, 0x51, 0xe4, 0x23, 0x9c, 0x74, 0xa3, 0x2e, 0xb1, 0x94, 0xa8, 0xe9,
	0xd8, 0x57, 0xe9, 0xb1, 0x4e, 0x69, 0xf0, 0x0a, 0xc7, 0x04, 0xb3, 0x27, 0xd4, 0xf4, 0xd0, 0x27,
	0xf6, 0x69, 0xe7, 0xb6, 0x1d, 0xb7, 0xba, 0xc9, 0x8e, 0x7c, 0xea, 0x0b, 0xde, 0x4d, 0x25, 0x0a,
	0xf7, 0xe0, 0x34, 0xf2, 0x02, 0x03, 0x22, 0xef, 0x20, 0x4a, 0xdc, 0x82, 0xdd, 0x5b, 0x52, 0xf0,
	0x1b, 0xdd, 0x10, 0x3e, 0xea, 0x9a, 0x7e, 0x17, 0x05, 0xd2, 0xd7, 0x5e, 0xce, 0x86, 0x60, 0x17,
	0x10, 0xdf, 0x14, 0x95, 0xd2, 0xf6, 0x16, 0x1b, 0x8e, 0x8f, 0xe4, 0x3d, 0x80, 0xb3, 0x8d, 0xb0,
	0x78, 0xfb, 0x6c, 0xa4, 0x2d, 0x86, 0x5d, 0xc1, 0x59, 0xf8, 0x4c, 0xbc, 0x2e, 0x97, 0x56, 0xd8,
	0xba, 0xdd, 0x96, 0x3f, 0x85, 0xe7, 0x0c, 0xc8, 0x39, 0x52, 0xa2, 0x15, 0x59, 0x6e, 0x82, 0x75,
	0x3a, 0xb8, 0x1a, 0xfb, 0x2f, 0x3f, 0xfb, 0x17, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x0e, 0xf5, 0x03,
	0x03, 0x04, 0x00, 0x00,
}
