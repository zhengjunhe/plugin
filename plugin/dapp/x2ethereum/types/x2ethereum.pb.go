// Code generated by protoc-gen-go. DO NOT EDIT.
// source: x2ethereum.proto

package types

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EthBridgeStatus int32

const (
	EthBridgeStatus_PendingStatusText EthBridgeStatus = 0
	EthBridgeStatus_SuccessStatusText EthBridgeStatus = 1
	EthBridgeStatus_FailedStatusText  EthBridgeStatus = 2
)

var EthBridgeStatus_name = map[int32]string{
	0: "PendingStatusText",
	1: "SuccessStatusText",
	2: "FailedStatusText",
}
var EthBridgeStatus_value = map[string]int32{
	"PendingStatusText": 0,
	"SuccessStatusText": 1,
	"FailedStatusText":  2,
}

func (x EthBridgeStatus) String() string {
	return proto.EnumName(EthBridgeStatus_name, int32(x))
}
func (EthBridgeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{0}
}

type X2EthereumAction struct {
	// Types that are valid to be assigned to Value:
	//	*X2EthereumAction_EthBridgeClaim
	//	*X2EthereumAction_MsgBurn
	//	*X2EthereumAction_MsgLock
	//	*X2EthereumAction_MsgLogInValidator
	//	*X2EthereumAction_MsgLogOutValidator
	//	*X2EthereumAction_MsgSetConsensusNeeded
	Value                isX2EthereumAction_Value `protobuf_oneof:"value"`
	Ty                   int32                    `protobuf:"varint,10,opt,name=ty" json:"ty,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *X2EthereumAction) Reset()         { *m = X2EthereumAction{} }
func (m *X2EthereumAction) String() string { return proto.CompactTextString(m) }
func (*X2EthereumAction) ProtoMessage()    {}
func (*X2EthereumAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{0}
}
func (m *X2EthereumAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_X2EthereumAction.Unmarshal(m, b)
}
func (m *X2EthereumAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_X2EthereumAction.Marshal(b, m, deterministic)
}
func (dst *X2EthereumAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_X2EthereumAction.Merge(dst, src)
}
func (m *X2EthereumAction) XXX_Size() int {
	return xxx_messageInfo_X2EthereumAction.Size(m)
}
func (m *X2EthereumAction) XXX_DiscardUnknown() {
	xxx_messageInfo_X2EthereumAction.DiscardUnknown(m)
}

var xxx_messageInfo_X2EthereumAction proto.InternalMessageInfo

type isX2EthereumAction_Value interface {
	isX2EthereumAction_Value()
}

type X2EthereumAction_EthBridgeClaim struct {
	EthBridgeClaim *EthBridgeClaim `protobuf:"bytes,1,opt,name=ethBridgeClaim,oneof"`
}
type X2EthereumAction_MsgBurn struct {
	MsgBurn *MsgBurn `protobuf:"bytes,2,opt,name=msgBurn,oneof"`
}
type X2EthereumAction_MsgLock struct {
	MsgLock *MsgLock `protobuf:"bytes,3,opt,name=msgLock,oneof"`
}
type X2EthereumAction_MsgLogInValidator struct {
	MsgLogInValidator *MsgLogInValidator `protobuf:"bytes,4,opt,name=msgLogInValidator,oneof"`
}
type X2EthereumAction_MsgLogOutValidator struct {
	MsgLogOutValidator *MsgLogOutValidator `protobuf:"bytes,5,opt,name=msgLogOutValidator,oneof"`
}
type X2EthereumAction_MsgSetConsensusNeeded struct {
	MsgSetConsensusNeeded *MsgSetConsensusNeeded `protobuf:"bytes,6,opt,name=msgSetConsensusNeeded,oneof"`
}

func (*X2EthereumAction_EthBridgeClaim) isX2EthereumAction_Value()        {}
func (*X2EthereumAction_MsgBurn) isX2EthereumAction_Value()               {}
func (*X2EthereumAction_MsgLock) isX2EthereumAction_Value()               {}
func (*X2EthereumAction_MsgLogInValidator) isX2EthereumAction_Value()     {}
func (*X2EthereumAction_MsgLogOutValidator) isX2EthereumAction_Value()    {}
func (*X2EthereumAction_MsgSetConsensusNeeded) isX2EthereumAction_Value() {}

func (m *X2EthereumAction) GetValue() isX2EthereumAction_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *X2EthereumAction) GetEthBridgeClaim() *EthBridgeClaim {
	if x, ok := m.GetValue().(*X2EthereumAction_EthBridgeClaim); ok {
		return x.EthBridgeClaim
	}
	return nil
}

func (m *X2EthereumAction) GetMsgBurn() *MsgBurn {
	if x, ok := m.GetValue().(*X2EthereumAction_MsgBurn); ok {
		return x.MsgBurn
	}
	return nil
}

func (m *X2EthereumAction) GetMsgLock() *MsgLock {
	if x, ok := m.GetValue().(*X2EthereumAction_MsgLock); ok {
		return x.MsgLock
	}
	return nil
}

func (m *X2EthereumAction) GetMsgLogInValidator() *MsgLogInValidator {
	if x, ok := m.GetValue().(*X2EthereumAction_MsgLogInValidator); ok {
		return x.MsgLogInValidator
	}
	return nil
}

func (m *X2EthereumAction) GetMsgLogOutValidator() *MsgLogOutValidator {
	if x, ok := m.GetValue().(*X2EthereumAction_MsgLogOutValidator); ok {
		return x.MsgLogOutValidator
	}
	return nil
}

func (m *X2EthereumAction) GetMsgSetConsensusNeeded() *MsgSetConsensusNeeded {
	if x, ok := m.GetValue().(*X2EthereumAction_MsgSetConsensusNeeded); ok {
		return x.MsgSetConsensusNeeded
	}
	return nil
}

func (m *X2EthereumAction) GetTy() int32 {
	if m != nil {
		return m.Ty
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*X2EthereumAction) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _X2EthereumAction_OneofMarshaler, _X2EthereumAction_OneofUnmarshaler, _X2EthereumAction_OneofSizer, []interface{}{
		(*X2EthereumAction_EthBridgeClaim)(nil),
		(*X2EthereumAction_MsgBurn)(nil),
		(*X2EthereumAction_MsgLock)(nil),
		(*X2EthereumAction_MsgLogInValidator)(nil),
		(*X2EthereumAction_MsgLogOutValidator)(nil),
		(*X2EthereumAction_MsgSetConsensusNeeded)(nil),
	}
}

func _X2EthereumAction_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*X2EthereumAction)
	// value
	switch x := m.Value.(type) {
	case *X2EthereumAction_EthBridgeClaim:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.EthBridgeClaim); err != nil {
			return err
		}
	case *X2EthereumAction_MsgBurn:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MsgBurn); err != nil {
			return err
		}
	case *X2EthereumAction_MsgLock:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MsgLock); err != nil {
			return err
		}
	case *X2EthereumAction_MsgLogInValidator:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MsgLogInValidator); err != nil {
			return err
		}
	case *X2EthereumAction_MsgLogOutValidator:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MsgLogOutValidator); err != nil {
			return err
		}
	case *X2EthereumAction_MsgSetConsensusNeeded:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MsgSetConsensusNeeded); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("X2EthereumAction.Value has unexpected type %T", x)
	}
	return nil
}

func _X2EthereumAction_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*X2EthereumAction)
	switch tag {
	case 1: // value.ethBridgeClaim
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(EthBridgeClaim)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_EthBridgeClaim{msg}
		return true, err
	case 2: // value.msgBurn
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MsgBurn)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_MsgBurn{msg}
		return true, err
	case 3: // value.msgLock
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MsgLock)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_MsgLock{msg}
		return true, err
	case 4: // value.msgLogInValidator
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MsgLogInValidator)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_MsgLogInValidator{msg}
		return true, err
	case 5: // value.msgLogOutValidator
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MsgLogOutValidator)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_MsgLogOutValidator{msg}
		return true, err
	case 6: // value.msgSetConsensusNeeded
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MsgSetConsensusNeeded)
		err := b.DecodeMessage(msg)
		m.Value = &X2EthereumAction_MsgSetConsensusNeeded{msg}
		return true, err
	default:
		return false, nil
	}
}

func _X2EthereumAction_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*X2EthereumAction)
	// value
	switch x := m.Value.(type) {
	case *X2EthereumAction_EthBridgeClaim:
		s := proto.Size(x.EthBridgeClaim)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *X2EthereumAction_MsgBurn:
		s := proto.Size(x.MsgBurn)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *X2EthereumAction_MsgLock:
		s := proto.Size(x.MsgLock)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *X2EthereumAction_MsgLogInValidator:
		s := proto.Size(x.MsgLogInValidator)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *X2EthereumAction_MsgLogOutValidator:
		s := proto.Size(x.MsgLogOutValidator)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *X2EthereumAction_MsgSetConsensusNeeded:
		s := proto.Size(x.MsgSetConsensusNeeded)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type MsgSetConsensusNeeded struct {
	Power                float64  `protobuf:"fixed64,1,opt,name=power" json:"power,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgSetConsensusNeeded) Reset()         { *m = MsgSetConsensusNeeded{} }
func (m *MsgSetConsensusNeeded) String() string { return proto.CompactTextString(m) }
func (*MsgSetConsensusNeeded) ProtoMessage()    {}
func (*MsgSetConsensusNeeded) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{1}
}
func (m *MsgSetConsensusNeeded) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgSetConsensusNeeded.Unmarshal(m, b)
}
func (m *MsgSetConsensusNeeded) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgSetConsensusNeeded.Marshal(b, m, deterministic)
}
func (dst *MsgSetConsensusNeeded) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetConsensusNeeded.Merge(dst, src)
}
func (m *MsgSetConsensusNeeded) XXX_Size() int {
	return xxx_messageInfo_MsgSetConsensusNeeded.Size(m)
}
func (m *MsgSetConsensusNeeded) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetConsensusNeeded.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetConsensusNeeded proto.InternalMessageInfo

func (m *MsgSetConsensusNeeded) GetPower() float64 {
	if m != nil {
		return m.Power
	}
	return 0
}

type MsgLogInValidator struct {
	Address              string   `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	Power                float64  `protobuf:"fixed64,2,opt,name=power" json:"power,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgLogInValidator) Reset()         { *m = MsgLogInValidator{} }
func (m *MsgLogInValidator) String() string { return proto.CompactTextString(m) }
func (*MsgLogInValidator) ProtoMessage()    {}
func (*MsgLogInValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{2}
}
func (m *MsgLogInValidator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgLogInValidator.Unmarshal(m, b)
}
func (m *MsgLogInValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgLogInValidator.Marshal(b, m, deterministic)
}
func (dst *MsgLogInValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgLogInValidator.Merge(dst, src)
}
func (m *MsgLogInValidator) XXX_Size() int {
	return xxx_messageInfo_MsgLogInValidator.Size(m)
}
func (m *MsgLogInValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgLogInValidator.DiscardUnknown(m)
}

var xxx_messageInfo_MsgLogInValidator proto.InternalMessageInfo

func (m *MsgLogInValidator) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MsgLogInValidator) GetPower() float64 {
	if m != nil {
		return m.Power
	}
	return 0
}

type MsgLogOutValidator struct {
	Address              string   `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
	Power                float64  `protobuf:"fixed64,2,opt,name=power" json:"power,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgLogOutValidator) Reset()         { *m = MsgLogOutValidator{} }
func (m *MsgLogOutValidator) String() string { return proto.CompactTextString(m) }
func (*MsgLogOutValidator) ProtoMessage()    {}
func (*MsgLogOutValidator) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{3}
}
func (m *MsgLogOutValidator) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgLogOutValidator.Unmarshal(m, b)
}
func (m *MsgLogOutValidator) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgLogOutValidator.Marshal(b, m, deterministic)
}
func (dst *MsgLogOutValidator) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgLogOutValidator.Merge(dst, src)
}
func (m *MsgLogOutValidator) XXX_Size() int {
	return xxx_messageInfo_MsgLogOutValidator.Size(m)
}
func (m *MsgLogOutValidator) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgLogOutValidator.DiscardUnknown(m)
}

var xxx_messageInfo_MsgLogOutValidator proto.InternalMessageInfo

func (m *MsgLogOutValidator) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MsgLogOutValidator) GetPower() float64 {
	if m != nil {
		return m.Power
	}
	return 0
}

// EthBridgeClaim is a structure that contains all the data for a particular bridge claim
type EthBridgeClaim struct {
	EthereumChainID       int64    `protobuf:"varint,1,opt,name=EthereumChainID" json:"EthereumChainID,omitempty"`
	BridgeContractAddress string   `protobuf:"bytes,2,opt,name=BridgeContractAddress" json:"BridgeContractAddress,omitempty"`
	Nonce                 int64    `protobuf:"varint,3,opt,name=Nonce" json:"Nonce,omitempty"`
	LocalCoinSymbol       string   `protobuf:"bytes,4,opt,name=localCoinSymbol" json:"localCoinSymbol,omitempty"`
	LocalCoinExec         string   `protobuf:"bytes,5,opt,name=localCoinExec" json:"localCoinExec,omitempty"`
	TokenContractAddress  string   `protobuf:"bytes,6,opt,name=TokenContractAddress" json:"TokenContractAddress,omitempty"`
	EthereumSender        string   `protobuf:"bytes,7,opt,name=EthereumSender" json:"EthereumSender,omitempty"`
	Chain33Receiver       string   `protobuf:"bytes,8,opt,name=Chain33Receiver" json:"Chain33Receiver,omitempty"`
	ValidatorAddress      string   `protobuf:"bytes,9,opt,name=ValidatorAddress" json:"ValidatorAddress,omitempty"`
	Amount                uint64   `protobuf:"varint,10,opt,name=Amount" json:"Amount,omitempty"`
	ClaimType             int64    `protobuf:"varint,11,opt,name=ClaimType" json:"ClaimType,omitempty"`
	EthSymbol             string   `protobuf:"bytes,12,opt,name=EthSymbol" json:"EthSymbol,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *EthBridgeClaim) Reset()         { *m = EthBridgeClaim{} }
func (m *EthBridgeClaim) String() string { return proto.CompactTextString(m) }
func (*EthBridgeClaim) ProtoMessage()    {}
func (*EthBridgeClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{4}
}
func (m *EthBridgeClaim) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EthBridgeClaim.Unmarshal(m, b)
}
func (m *EthBridgeClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EthBridgeClaim.Marshal(b, m, deterministic)
}
func (dst *EthBridgeClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EthBridgeClaim.Merge(dst, src)
}
func (m *EthBridgeClaim) XXX_Size() int {
	return xxx_messageInfo_EthBridgeClaim.Size(m)
}
func (m *EthBridgeClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_EthBridgeClaim.DiscardUnknown(m)
}

var xxx_messageInfo_EthBridgeClaim proto.InternalMessageInfo

func (m *EthBridgeClaim) GetEthereumChainID() int64 {
	if m != nil {
		return m.EthereumChainID
	}
	return 0
}

func (m *EthBridgeClaim) GetBridgeContractAddress() string {
	if m != nil {
		return m.BridgeContractAddress
	}
	return ""
}

func (m *EthBridgeClaim) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *EthBridgeClaim) GetLocalCoinSymbol() string {
	if m != nil {
		return m.LocalCoinSymbol
	}
	return ""
}

func (m *EthBridgeClaim) GetLocalCoinExec() string {
	if m != nil {
		return m.LocalCoinExec
	}
	return ""
}

func (m *EthBridgeClaim) GetTokenContractAddress() string {
	if m != nil {
		return m.TokenContractAddress
	}
	return ""
}

func (m *EthBridgeClaim) GetEthereumSender() string {
	if m != nil {
		return m.EthereumSender
	}
	return ""
}

func (m *EthBridgeClaim) GetChain33Receiver() string {
	if m != nil {
		return m.Chain33Receiver
	}
	return ""
}

func (m *EthBridgeClaim) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *EthBridgeClaim) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *EthBridgeClaim) GetClaimType() int64 {
	if m != nil {
		return m.ClaimType
	}
	return 0
}

func (m *EthBridgeClaim) GetEthSymbol() string {
	if m != nil {
		return m.EthSymbol
	}
	return ""
}

// OracleClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
type OracleClaimContent struct {
	Chain33Receiver      string   `protobuf:"bytes,1,opt,name=Chain33Receiver" json:"Chain33Receiver,omitempty"`
	Amount               uint64   `protobuf:"varint,2,opt,name=Amount" json:"Amount,omitempty"`
	ClaimType            int64    `protobuf:"varint,3,opt,name=ClaimType" json:"ClaimType,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OracleClaimContent) Reset()         { *m = OracleClaimContent{} }
func (m *OracleClaimContent) String() string { return proto.CompactTextString(m) }
func (*OracleClaimContent) ProtoMessage()    {}
func (*OracleClaimContent) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{5}
}
func (m *OracleClaimContent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OracleClaimContent.Unmarshal(m, b)
}
func (m *OracleClaimContent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OracleClaimContent.Marshal(b, m, deterministic)
}
func (dst *OracleClaimContent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OracleClaimContent.Merge(dst, src)
}
func (m *OracleClaimContent) XXX_Size() int {
	return xxx_messageInfo_OracleClaimContent.Size(m)
}
func (m *OracleClaimContent) XXX_DiscardUnknown() {
	xxx_messageInfo_OracleClaimContent.DiscardUnknown(m)
}

var xxx_messageInfo_OracleClaimContent proto.InternalMessageInfo

func (m *OracleClaimContent) GetChain33Receiver() string {
	if m != nil {
		return m.Chain33Receiver
	}
	return ""
}

func (m *OracleClaimContent) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *OracleClaimContent) GetClaimType() int64 {
	if m != nil {
		return m.ClaimType
	}
	return 0
}

// MsgBurn defines a message for burning coins and triggering a related event
type MsgBurn struct {
	EthereumChainID      int64    `protobuf:"varint,1,opt,name=EthereumChainID" json:"EthereumChainID,omitempty"`
	TokenContract        string   `protobuf:"bytes,2,opt,name=TokenContract" json:"TokenContract,omitempty"`
	Chain33Sender        string   `protobuf:"bytes,3,opt,name=Chain33Sender" json:"Chain33Sender,omitempty"`
	EthereumReceiver     string   `protobuf:"bytes,4,opt,name=EthereumReceiver" json:"EthereumReceiver,omitempty"`
	Amount               uint64   `protobuf:"varint,5,opt,name=Amount" json:"Amount,omitempty"`
	LocalCoinSymbol      string   `protobuf:"bytes,6,opt,name=localCoinSymbol" json:"localCoinSymbol,omitempty"`
	LocalCoinExec        string   `protobuf:"bytes,7,opt,name=localCoinExec" json:"localCoinExec,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgBurn) Reset()         { *m = MsgBurn{} }
func (m *MsgBurn) String() string { return proto.CompactTextString(m) }
func (*MsgBurn) ProtoMessage()    {}
func (*MsgBurn) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{6}
}
func (m *MsgBurn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgBurn.Unmarshal(m, b)
}
func (m *MsgBurn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgBurn.Marshal(b, m, deterministic)
}
func (dst *MsgBurn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgBurn.Merge(dst, src)
}
func (m *MsgBurn) XXX_Size() int {
	return xxx_messageInfo_MsgBurn.Size(m)
}
func (m *MsgBurn) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgBurn.DiscardUnknown(m)
}

var xxx_messageInfo_MsgBurn proto.InternalMessageInfo

func (m *MsgBurn) GetEthereumChainID() int64 {
	if m != nil {
		return m.EthereumChainID
	}
	return 0
}

func (m *MsgBurn) GetTokenContract() string {
	if m != nil {
		return m.TokenContract
	}
	return ""
}

func (m *MsgBurn) GetChain33Sender() string {
	if m != nil {
		return m.Chain33Sender
	}
	return ""
}

func (m *MsgBurn) GetEthereumReceiver() string {
	if m != nil {
		return m.EthereumReceiver
	}
	return ""
}

func (m *MsgBurn) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *MsgBurn) GetLocalCoinSymbol() string {
	if m != nil {
		return m.LocalCoinSymbol
	}
	return ""
}

func (m *MsgBurn) GetLocalCoinExec() string {
	if m != nil {
		return m.LocalCoinExec
	}
	return ""
}

// MsgLock defines a message for locking coins and triggering a related event
type MsgLock struct {
	EthereumChainID      int64    `protobuf:"varint,1,opt,name=EthereumChainID" json:"EthereumChainID,omitempty"`
	TokenContract        string   `protobuf:"bytes,2,opt,name=TokenContract" json:"TokenContract,omitempty"`
	Chain33Sender        string   `protobuf:"bytes,3,opt,name=Chain33Sender" json:"Chain33Sender,omitempty"`
	EthereumReceiver     string   `protobuf:"bytes,4,opt,name=EthereumReceiver" json:"EthereumReceiver,omitempty"`
	Amount               uint64   `protobuf:"varint,5,opt,name=Amount" json:"Amount,omitempty"`
	LocalCoinSymbol      string   `protobuf:"bytes,6,opt,name=localCoinSymbol" json:"localCoinSymbol,omitempty"`
	LocalCoinExec        string   `protobuf:"bytes,7,opt,name=localCoinExec" json:"localCoinExec,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MsgLock) Reset()         { *m = MsgLock{} }
func (m *MsgLock) String() string { return proto.CompactTextString(m) }
func (*MsgLock) ProtoMessage()    {}
func (*MsgLock) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{7}
}
func (m *MsgLock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MsgLock.Unmarshal(m, b)
}
func (m *MsgLock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MsgLock.Marshal(b, m, deterministic)
}
func (dst *MsgLock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgLock.Merge(dst, src)
}
func (m *MsgLock) XXX_Size() int {
	return xxx_messageInfo_MsgLock.Size(m)
}
func (m *MsgLock) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgLock.DiscardUnknown(m)
}

var xxx_messageInfo_MsgLock proto.InternalMessageInfo

func (m *MsgLock) GetEthereumChainID() int64 {
	if m != nil {
		return m.EthereumChainID
	}
	return 0
}

func (m *MsgLock) GetTokenContract() string {
	if m != nil {
		return m.TokenContract
	}
	return ""
}

func (m *MsgLock) GetChain33Sender() string {
	if m != nil {
		return m.Chain33Sender
	}
	return ""
}

func (m *MsgLock) GetEthereumReceiver() string {
	if m != nil {
		return m.EthereumReceiver
	}
	return ""
}

func (m *MsgLock) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *MsgLock) GetLocalCoinSymbol() string {
	if m != nil {
		return m.LocalCoinSymbol
	}
	return ""
}

func (m *MsgLock) GetLocalCoinExec() string {
	if m != nil {
		return m.LocalCoinExec
	}
	return ""
}

type QueryEthProphecyParams struct {
	EthereumChainID       int64    `protobuf:"varint,1,opt,name=EthereumChainID" json:"EthereumChainID,omitempty"`
	BridgeContractAddress string   `protobuf:"bytes,2,opt,name=BridgeContractAddress" json:"BridgeContractAddress,omitempty"`
	Nonce                 int64    `protobuf:"varint,3,opt,name=Nonce" json:"Nonce,omitempty"`
	Symbol                string   `protobuf:"bytes,4,opt,name=Symbol" json:"Symbol,omitempty"`
	TokenContractAddress  string   `protobuf:"bytes,5,opt,name=TokenContractAddress" json:"TokenContractAddress,omitempty"`
	EthereumSender        string   `protobuf:"bytes,6,opt,name=EthereumSender" json:"EthereumSender,omitempty"`
	XXX_NoUnkeyedLiteral  struct{} `json:"-"`
	XXX_unrecognized      []byte   `json:"-"`
	XXX_sizecache         int32    `json:"-"`
}

func (m *QueryEthProphecyParams) Reset()         { *m = QueryEthProphecyParams{} }
func (m *QueryEthProphecyParams) String() string { return proto.CompactTextString(m) }
func (*QueryEthProphecyParams) ProtoMessage()    {}
func (*QueryEthProphecyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{8}
}
func (m *QueryEthProphecyParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryEthProphecyParams.Unmarshal(m, b)
}
func (m *QueryEthProphecyParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryEthProphecyParams.Marshal(b, m, deterministic)
}
func (dst *QueryEthProphecyParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEthProphecyParams.Merge(dst, src)
}
func (m *QueryEthProphecyParams) XXX_Size() int {
	return xxx_messageInfo_QueryEthProphecyParams.Size(m)
}
func (m *QueryEthProphecyParams) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEthProphecyParams.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEthProphecyParams proto.InternalMessageInfo

func (m *QueryEthProphecyParams) GetEthereumChainID() int64 {
	if m != nil {
		return m.EthereumChainID
	}
	return 0
}

func (m *QueryEthProphecyParams) GetBridgeContractAddress() string {
	if m != nil {
		return m.BridgeContractAddress
	}
	return ""
}

func (m *QueryEthProphecyParams) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *QueryEthProphecyParams) GetSymbol() string {
	if m != nil {
		return m.Symbol
	}
	return ""
}

func (m *QueryEthProphecyParams) GetTokenContractAddress() string {
	if m != nil {
		return m.TokenContractAddress
	}
	return ""
}

func (m *QueryEthProphecyParams) GetEthereumSender() string {
	if m != nil {
		return m.EthereumSender
	}
	return ""
}

type QueryEthProphecyResponse struct {
	ID                   string            `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Status               *ProphecyStatus   `protobuf:"bytes,2,opt,name=Status" json:"Status,omitempty"`
	Claims               []*EthBridgeClaim `protobuf:"bytes,3,rep,name=Claims" json:"Claims,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *QueryEthProphecyResponse) Reset()         { *m = QueryEthProphecyResponse{} }
func (m *QueryEthProphecyResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEthProphecyResponse) ProtoMessage()    {}
func (*QueryEthProphecyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{9}
}
func (m *QueryEthProphecyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryEthProphecyResponse.Unmarshal(m, b)
}
func (m *QueryEthProphecyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryEthProphecyResponse.Marshal(b, m, deterministic)
}
func (dst *QueryEthProphecyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEthProphecyResponse.Merge(dst, src)
}
func (m *QueryEthProphecyResponse) XXX_Size() int {
	return xxx_messageInfo_QueryEthProphecyResponse.Size(m)
}
func (m *QueryEthProphecyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEthProphecyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEthProphecyResponse proto.InternalMessageInfo

func (m *QueryEthProphecyResponse) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *QueryEthProphecyResponse) GetStatus() *ProphecyStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *QueryEthProphecyResponse) GetClaims() []*EthBridgeClaim {
	if m != nil {
		return m.Claims
	}
	return nil
}

type ProphecyStatus struct {
	Text                 EthBridgeStatus `protobuf:"varint,1,opt,name=Text,enum=types.EthBridgeStatus" json:"Text,omitempty"`
	FinalClaim           string          `protobuf:"bytes,2,opt,name=FinalClaim" json:"FinalClaim,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ProphecyStatus) Reset()         { *m = ProphecyStatus{} }
func (m *ProphecyStatus) String() string { return proto.CompactTextString(m) }
func (*ProphecyStatus) ProtoMessage()    {}
func (*ProphecyStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_x2ethereum_ce546ae61fce443d, []int{10}
}
func (m *ProphecyStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProphecyStatus.Unmarshal(m, b)
}
func (m *ProphecyStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProphecyStatus.Marshal(b, m, deterministic)
}
func (dst *ProphecyStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProphecyStatus.Merge(dst, src)
}
func (m *ProphecyStatus) XXX_Size() int {
	return xxx_messageInfo_ProphecyStatus.Size(m)
}
func (m *ProphecyStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ProphecyStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ProphecyStatus proto.InternalMessageInfo

func (m *ProphecyStatus) GetText() EthBridgeStatus {
	if m != nil {
		return m.Text
	}
	return EthBridgeStatus_PendingStatusText
}

func (m *ProphecyStatus) GetFinalClaim() string {
	if m != nil {
		return m.FinalClaim
	}
	return ""
}

func init() {
	proto.RegisterType((*X2EthereumAction)(nil), "types.X2ethereumAction")
	proto.RegisterType((*MsgSetConsensusNeeded)(nil), "types.MsgSetConsensusNeeded")
	proto.RegisterType((*MsgLogInValidator)(nil), "types.MsgLogInValidator")
	proto.RegisterType((*MsgLogOutValidator)(nil), "types.MsgLogOutValidator")
	proto.RegisterType((*EthBridgeClaim)(nil), "types.EthBridgeClaim")
	proto.RegisterType((*OracleClaimContent)(nil), "types.OracleClaimContent")
	proto.RegisterType((*MsgBurn)(nil), "types.MsgBurn")
	proto.RegisterType((*MsgLock)(nil), "types.MsgLock")
	proto.RegisterType((*QueryEthProphecyParams)(nil), "types.QueryEthProphecyParams")
	proto.RegisterType((*QueryEthProphecyResponse)(nil), "types.QueryEthProphecyResponse")
	proto.RegisterType((*ProphecyStatus)(nil), "types.ProphecyStatus")
	proto.RegisterEnum("types.EthBridgeStatus", EthBridgeStatus_name, EthBridgeStatus_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for X2Ethereum service

type X2EthereumClient interface {
}

type x2EthereumClient struct {
	cc *grpc.ClientConn
}

func NewX2EthereumClient(cc *grpc.ClientConn) X2EthereumClient {
	return &x2EthereumClient{cc}
}

// Server API for X2Ethereum service

type X2EthereumServer interface {
}

func RegisterX2EthereumServer(s *grpc.Server, srv X2EthereumServer) {
	s.RegisterService(&_X2Ethereum_serviceDesc, srv)
}

var _X2Ethereum_serviceDesc = grpc.ServiceDesc{
	ServiceName: "types.x2ethereum",
	HandlerType: (*X2EthereumServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "x2ethereum.proto",
}

func init() { proto.RegisterFile("x2ethereum.proto", fileDescriptor_x2ethereum_ce546ae61fce443d) }

var fileDescriptor_x2ethereum_ce546ae61fce443d = []byte{
	// 772 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x56, 0x4f, 0x4f, 0xdb, 0x4a,
	0x10, 0x4f, 0x6c, 0xe2, 0xbc, 0x0c, 0x90, 0x17, 0x56, 0x24, 0xf2, 0x93, 0xd0, 0x13, 0x8a, 0xd0,
	0x13, 0x8a, 0x04, 0x87, 0xf0, 0xee, 0x15, 0x84, 0xa0, 0xa0, 0xb6, 0x40, 0x37, 0x69, 0xd5, 0x43,
	0x2f, 0xc6, 0x1e, 0x25, 0x16, 0xce, 0x3a, 0xb2, 0xd7, 0x94, 0x5c, 0x7b, 0xea, 0xbd, 0x9f, 0xa4,
	0xdf, 0xa7, 0x1f, 0xa6, 0xf2, 0x78, 0x9d, 0xd8, 0x8e, 0x91, 0xca, 0xa9, 0x97, 0xde, 0xb2, 0xbf,
	0xf9, 0xcd, 0x6f, 0xe7, 0x5f, 0xc6, 0x0b, 0xad, 0xa7, 0x3e, 0xca, 0x19, 0x06, 0x18, 0xcd, 0x4f,
	0x17, 0x81, 0x2f, 0x7d, 0x56, 0x93, 0xcb, 0x05, 0x86, 0xdd, 0xef, 0x3a, 0xb4, 0x3e, 0xae, 0x6c,
	0xe7, 0xb6, 0x74, 0x7d, 0xc1, 0x5e, 0x41, 0x13, 0xe5, 0xec, 0x22, 0x70, 0x9d, 0x29, 0x0e, 0x3c,
	0xcb, 0x9d, 0x9b, 0xd5, 0xc3, 0xea, 0xf1, 0x76, 0xbf, 0x7d, 0x4a, 0x4e, 0xa7, 0xc3, 0x9c, 0x71,
	0x54, 0xe1, 0x05, 0x3a, 0xeb, 0x41, 0x7d, 0x1e, 0x4e, 0x2f, 0xa2, 0x40, 0x98, 0x1a, 0x79, 0x36,
	0x95, 0xe7, 0xdb, 0x04, 0x1d, 0x55, 0x78, 0x4a, 0x50, 0xdc, 0x37, 0xbe, 0xfd, 0x60, 0xea, 0x45,
	0x6e, 0x8c, 0x2a, 0x6e, 0xfc, 0x93, 0x8d, 0x60, 0x8f, 0x7e, 0x4e, 0xaf, 0xc5, 0x07, 0xcb, 0x73,
	0x1d, 0x4b, 0xfa, 0x81, 0xb9, 0x45, 0x5e, 0x66, 0xd6, 0x2b, 0x6b, 0x1f, 0x55, 0xf8, 0xa6, 0x13,
	0x7b, 0x0d, 0x2c, 0x01, 0x6f, 0x23, 0xb9, 0x96, 0xaa, 0x91, 0xd4, 0x3f, 0x39, 0xa9, 0x2c, 0x61,
	0x54, 0xe1, 0x25, 0x6e, 0x6c, 0x02, 0xed, 0x79, 0x38, 0x1d, 0xa3, 0x1c, 0xf8, 0x22, 0x44, 0x11,
	0x46, 0xe1, 0x0d, 0xa2, 0x83, 0x8e, 0x69, 0x90, 0xde, 0xc1, 0x5a, 0x6f, 0x93, 0x33, 0xaa, 0xf0,
	0x72, 0x67, 0xd6, 0x04, 0x4d, 0x2e, 0x4d, 0x38, 0xac, 0x1e, 0xd7, 0xb8, 0x26, 0x97, 0x17, 0x75,
	0xa8, 0x3d, 0x5a, 0x5e, 0x84, 0xdd, 0x13, 0x68, 0x97, 0x4a, 0xb1, 0x7d, 0xa8, 0x2d, 0xfc, 0xcf,
	0x18, 0x50, 0xbb, 0xaa, 0x3c, 0x39, 0x74, 0x07, 0xb0, 0xb7, 0x51, 0x14, 0x66, 0x42, 0xdd, 0x72,
	0x9c, 0x00, 0xc3, 0x90, 0xc8, 0x0d, 0x9e, 0x1e, 0xd7, 0x22, 0x5a, 0x56, 0xe4, 0x12, 0xd8, 0x66,
	0x39, 0x5e, 0xac, 0xf2, 0x43, 0x87, 0x66, 0x7e, 0x78, 0xd8, 0x31, 0xfc, 0x3d, 0x54, 0xd3, 0x37,
	0x98, 0x59, 0xae, 0xb8, 0xbe, 0x24, 0x29, 0x9d, 0x17, 0x61, 0xf6, 0x3f, 0xb4, 0x95, 0xa3, 0x2f,
	0x64, 0x60, 0xd9, 0xf2, 0x5c, 0x5d, 0xad, 0xd1, 0xd5, 0xe5, 0xc6, 0x38, 0x90, 0x1b, 0x5f, 0xd8,
	0x48, 0xc3, 0xa5, 0xf3, 0xe4, 0x10, 0xdf, 0xea, 0xf9, 0xb6, 0xe5, 0x0d, 0x7c, 0x57, 0x8c, 0x97,
	0xf3, 0x7b, 0xdf, 0xa3, 0x31, 0x6a, 0xf0, 0x22, 0xcc, 0x8e, 0x60, 0x77, 0x05, 0x0d, 0x9f, 0xd0,
	0xa6, 0x19, 0x69, 0xf0, 0x3c, 0xc8, 0xfa, 0xb0, 0x3f, 0xf1, 0x1f, 0x50, 0x14, 0x43, 0x33, 0x88,
	0x5c, 0x6a, 0x63, 0xff, 0x51, 0x2d, 0x28, 0xc5, 0x31, 0x0a, 0x07, 0x03, 0xb3, 0x4e, 0xec, 0x02,
	0x1a, 0xc7, 0x4a, 0x25, 0x38, 0x3b, 0xe3, 0x68, 0xa3, 0xfb, 0x88, 0x81, 0xf9, 0x57, 0x12, 0x6b,
	0x01, 0x66, 0x3d, 0x68, 0xad, 0x7a, 0x93, 0x46, 0xd0, 0x20, 0xea, 0x06, 0xce, 0x3a, 0x60, 0x9c,
	0xcf, 0xfd, 0x48, 0x48, 0x9a, 0xb0, 0x2d, 0xae, 0x4e, 0xec, 0x00, 0x1a, 0xd4, 0x98, 0xc9, 0x72,
	0x81, 0xe6, 0x36, 0xd5, 0x6c, 0x0d, 0xc4, 0xd6, 0xa1, 0x9c, 0xa9, 0x8a, 0xed, 0x90, 0xf4, 0x1a,
	0xe8, 0x4a, 0x60, 0xb7, 0x81, 0x65, 0x7b, 0x49, 0x6b, 0xe3, 0x7c, 0x51, 0xc8, 0xb2, 0xf8, 0xab,
	0xe5, 0xf1, 0xaf, 0x63, 0xd2, 0x9e, 0x8f, 0x49, 0x2f, 0xc4, 0xd4, 0xfd, 0xa6, 0x41, 0x5d, 0xed,
	0x95, 0x17, 0x4c, 0xd3, 0x11, 0xec, 0xe6, 0xba, 0xa2, 0xa6, 0x28, 0x0f, 0xc6, 0x2c, 0x15, 0xa4,
	0x6a, 0x91, 0x9e, 0xb0, 0x72, 0x60, 0x5c, 0xf7, 0x54, 0x7e, 0x95, 0x62, 0x32, 0x4e, 0x1b, 0x78,
	0x26, 0xc7, 0x5a, 0x2e, 0xc7, 0x92, 0x89, 0x34, 0x7e, 0x71, 0x22, 0xeb, 0x25, 0x13, 0x99, 0x56,
	0x85, 0xd6, 0xe6, 0x9f, 0xaa, 0xa4, 0x55, 0xf9, 0xa2, 0x41, 0xe7, 0x5d, 0x84, 0xc1, 0x72, 0x28,
	0x67, 0x77, 0x81, 0xbf, 0x98, 0xa1, 0xbd, 0xbc, 0xb3, 0x02, 0x6b, 0x1e, 0xfe, 0xa6, 0x45, 0xd4,
	0x01, 0x23, 0xb7, 0x7f, 0xd4, 0xe9, 0xd9, 0x85, 0x52, 0x7b, 0xd1, 0x42, 0x31, 0xca, 0x16, 0x4a,
	0xf7, 0x6b, 0x15, 0xcc, 0x62, 0x11, 0x38, 0x86, 0x8b, 0xf8, 0x6b, 0x12, 0x7f, 0x75, 0x54, 0xe6,
	0x0d, 0xae, 0x5d, 0x5f, 0xb2, 0x13, 0x30, 0xc6, 0xd2, 0x92, 0x51, 0xa8, 0xbe, 0xe4, 0xe9, 0x1b,
	0x20, 0x75, 0x4c, 0x8c, 0x5c, 0x91, 0x62, 0x3a, 0xfd, 0x33, 0x43, 0x53, 0x3f, 0xd4, 0x9f, 0x7d,
	0x32, 0x70, 0x45, 0xea, 0x7e, 0x82, 0x66, 0x5e, 0x88, 0xf5, 0x60, 0x6b, 0x82, 0x4f, 0x92, 0x22,
	0x68, 0xf6, 0x3b, 0x45, 0x77, 0x75, 0x1d, 0x71, 0xd8, 0xbf, 0x00, 0x57, 0xae, 0xb0, 0xbc, 0xe4,
	0x8d, 0x92, 0x54, 0x3f, 0x83, 0xf4, 0xde, 0x53, 0x4b, 0xb3, 0x8e, 0xac, 0x0d, 0x7b, 0x77, 0x28,
	0x1c, 0x57, 0x4c, 0x13, 0x20, 0xd6, 0x69, 0x55, 0x62, 0x78, 0x1c, 0xd9, 0x36, 0x86, 0x61, 0x06,
	0xae, 0xb2, 0x7d, 0x68, 0x5d, 0x59, 0xae, 0x87, 0x4e, 0x06, 0xd5, 0xfa, 0x3b, 0x00, 0xeb, 0xe7,
	0xd4, 0xbd, 0x41, 0xef, 0xa9, 0xb3, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa9, 0x07, 0x8c, 0xa0,
	0x63, 0x09, 0x00, 0x00,
}
