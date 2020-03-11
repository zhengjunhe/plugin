// Code generated by protoc-gen-go. DO NOT EDIT.
// source: oracle.proto

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

// EthBridgeClaim is a structure that contains all the data for a particular bridge claim
type OracleClaim struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	ValidatorAddress     string   `protobuf:"bytes,2,opt,name=ValidatorAddress" json:"ValidatorAddress,omitempty"`
	Content              string   `protobuf:"bytes,3,opt,name=Content" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OracleClaim) Reset()         { *m = OracleClaim{} }
func (m *OracleClaim) String() string { return proto.CompactTextString(m) }
func (*OracleClaim) ProtoMessage()    {}
func (*OracleClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_oracle_2eaade929310bf8a, []int{0}
}
func (m *OracleClaim) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OracleClaim.Unmarshal(m, b)
}
func (m *OracleClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OracleClaim.Marshal(b, m, deterministic)
}
func (dst *OracleClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OracleClaim.Merge(dst, src)
}
func (m *OracleClaim) XXX_Size() int {
	return xxx_messageInfo_OracleClaim.Size(m)
}
func (m *OracleClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_OracleClaim.DiscardUnknown(m)
}

var xxx_messageInfo_OracleClaim proto.InternalMessageInfo

func (m *OracleClaim) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *OracleClaim) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *OracleClaim) GetContent() string {
	if m != nil {
		return m.Content
	}
	return ""
}

type AddressArray struct {
	ClaimValidator       []string `protobuf:"bytes,1,rep,name=ClaimValidator" json:"ClaimValidator,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddressArray) Reset()         { *m = AddressArray{} }
func (m *AddressArray) String() string { return proto.CompactTextString(m) }
func (*AddressArray) ProtoMessage()    {}
func (*AddressArray) Descriptor() ([]byte, []int) {
	return fileDescriptor_oracle_2eaade929310bf8a, []int{1}
}
func (m *AddressArray) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressArray.Unmarshal(m, b)
}
func (m *AddressArray) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressArray.Marshal(b, m, deterministic)
}
func (dst *AddressArray) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressArray.Merge(dst, src)
}
func (m *AddressArray) XXX_Size() int {
	return xxx_messageInfo_AddressArray.Size(m)
}
func (m *AddressArray) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressArray.DiscardUnknown(m)
}

var xxx_messageInfo_AddressArray proto.InternalMessageInfo

func (m *AddressArray) GetClaimValidator() []string {
	if m != nil {
		return m.ClaimValidator
	}
	return nil
}

type Prophecy struct {
	ID                   string                   `protobuf:"bytes,1,opt,name=ID" json:"ID,omitempty"`
	Status               *ProphecyStatus          `protobuf:"bytes,2,opt,name=Status" json:"Status,omitempty"`
	ClaimValidators      map[string]*AddressArray `protobuf:"bytes,3,rep,name=ClaimValidators" json:"ClaimValidators,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ValidatorClaims      map[string]string        `protobuf:"bytes,4,rep,name=ValidatorClaims" json:"ValidatorClaims,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Prophecy) Reset()         { *m = Prophecy{} }
func (m *Prophecy) String() string { return proto.CompactTextString(m) }
func (*Prophecy) ProtoMessage()    {}
func (*Prophecy) Descriptor() ([]byte, []int) {
	return fileDescriptor_oracle_2eaade929310bf8a, []int{2}
}
func (m *Prophecy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Prophecy.Unmarshal(m, b)
}
func (m *Prophecy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Prophecy.Marshal(b, m, deterministic)
}
func (dst *Prophecy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Prophecy.Merge(dst, src)
}
func (m *Prophecy) XXX_Size() int {
	return xxx_messageInfo_Prophecy.Size(m)
}
func (m *Prophecy) XXX_DiscardUnknown() {
	xxx_messageInfo_Prophecy.DiscardUnknown(m)
}

var xxx_messageInfo_Prophecy proto.InternalMessageInfo

func (m *Prophecy) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Prophecy) GetStatus() *ProphecyStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *Prophecy) GetClaimValidators() map[string]*AddressArray {
	if m != nil {
		return m.ClaimValidators
	}
	return nil
}

func (m *Prophecy) GetValidatorClaims() map[string]string {
	if m != nil {
		return m.ValidatorClaims
	}
	return nil
}

func init() {
	proto.RegisterType((*OracleClaim)(nil), "types.OracleClaim")
	proto.RegisterType((*AddressArray)(nil), "types.AddressArray")
	proto.RegisterType((*Prophecy)(nil), "types.Prophecy")
	proto.RegisterMapType((map[string]*AddressArray)(nil), "types.Prophecy.ClaimValidatorsEntry")
	proto.RegisterMapType((map[string]string)(nil), "types.Prophecy.ValidatorClaimsEntry")
}

func init() { proto.RegisterFile("oracle.proto", fileDescriptor_oracle_2eaade929310bf8a) }

var fileDescriptor_oracle_2eaade929310bf8a = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x51, 0x4b, 0xc3, 0x30,
	0x10, 0x80, 0x69, 0xe3, 0xa6, 0xbb, 0x8d, 0x59, 0x62, 0x85, 0xd0, 0xa7, 0x52, 0x44, 0xaa, 0x60,
	0x1f, 0x2a, 0x88, 0xf8, 0x36, 0x37, 0x1f, 0xf6, 0xa2, 0x52, 0x41, 0x9f, 0x63, 0x1b, 0xd8, 0xb0,
	0x6b, 0x4a, 0x9a, 0x8a, 0xfd, 0x03, 0xfe, 0x6e, 0x59, 0x9a, 0x8d, 0x9a, 0xf5, 0xad, 0x77, 0xf7,
	0xdd, 0x77, 0x77, 0x25, 0x30, 0xe1, 0x82, 0xa6, 0x39, 0x8b, 0x4a, 0xc1, 0x25, 0xc7, 0x03, 0xd9,
	0x94, 0xac, 0xf2, 0x9c, 0x9f, 0x98, 0xc9, 0x15, 0x13, 0xac, 0xde, 0xb4, 0x85, 0x20, 0x85, 0xf1,
	0x8b, 0x02, 0xe7, 0x39, 0x5d, 0x6f, 0xf0, 0x14, 0xec, 0xe5, 0x82, 0x58, 0xbe, 0x15, 0x8e, 0x12,
	0x7b, 0xb9, 0xc0, 0xd7, 0xe0, 0xbc, 0xd3, 0x7c, 0x9d, 0x51, 0xc9, 0xc5, 0x2c, 0xcb, 0x04, 0xab,
	0x2a, 0x62, 0xab, 0xea, 0x41, 0x1e, 0x13, 0x38, 0x9e, 0xf3, 0x42, 0xb2, 0x42, 0x12, 0xa4, 0x90,
	0x5d, 0x18, 0xdc, 0xc1, 0x44, 0x43, 0x33, 0x21, 0x68, 0x83, 0x2f, 0x61, 0xaa, 0xc6, 0xed, 0x15,
	0xc4, 0xf2, 0x51, 0x38, 0x4a, 0x8c, 0x6c, 0xf0, 0x8b, 0xe0, 0xe4, 0x55, 0xf0, 0x72, 0xc5, 0xd2,
	0xe6, 0x60, 0xb5, 0x1b, 0x18, 0xbe, 0x49, 0x2a, 0xeb, 0x76, 0xa1, 0x71, 0x7c, 0x1e, 0xa9, 0x1b,
	0xa3, 0x5d, 0x43, 0x5b, 0x4c, 0x34, 0x84, 0x9f, 0xe1, 0xf4, 0xbf, 0xbd, 0x22, 0xc8, 0x47, 0xe1,
	0x38, 0xbe, 0x30, 0xfa, 0x22, 0x03, 0x7b, 0x2a, 0xa4, 0x68, 0x12, 0xb3, 0x79, 0xeb, 0xdb, 0x47,
	0xaa, 0x56, 0x91, 0xa3, 0x7e, 0x9f, 0x81, 0x69, 0x9f, 0x91, 0xf5, 0x3e, 0xc0, 0xed, 0x1b, 0x8c,
	0x1d, 0x40, 0x5f, 0xac, 0xd1, 0x77, 0x6f, 0x3f, 0xf1, 0x15, 0x0c, 0xbe, 0x69, 0x5e, 0x33, 0x7d,
	0xf7, 0x99, 0x9e, 0xd7, 0xfd, 0xc3, 0x49, 0x4b, 0x3c, 0xd8, 0xf7, 0x96, 0xf7, 0x08, 0x6e, 0xdf,
	0x06, 0x3d, 0x62, 0xb7, 0x2b, 0x1e, 0x75, 0x1c, 0x9f, 0x43, 0xf5, 0x58, 0x6e, 0xff, 0x02, 0x00,
	0x00, 0xff, 0xff, 0xeb, 0x69, 0x11, 0x7c, 0x55, 0x02, 0x00, 0x00,
}
