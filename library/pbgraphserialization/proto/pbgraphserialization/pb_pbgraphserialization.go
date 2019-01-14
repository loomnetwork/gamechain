// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization/pbgraphserialization.proto

package pbgraphserialization_pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type SerializationId struct {
	SerializationId      uint32   `protobuf:"varint,1,opt,name=serializationId,proto3" json:"serializationId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SerializationId) Reset()         { *m = SerializationId{} }
func (m *SerializationId) String() string { return proto.CompactTextString(m) }
func (*SerializationId) ProtoMessage()    {}
func (*SerializationId) Descriptor() ([]byte, []int) {
	return fileDescriptor_pbgraphserialization_a6ed9de92956f5b6, []int{0}
}
func (m *SerializationId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializationId.Unmarshal(m, b)
}
func (m *SerializationId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializationId.Marshal(b, m, deterministic)
}
func (dst *SerializationId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializationId.Merge(dst, src)
}
func (m *SerializationId) XXX_Size() int {
	return xxx_messageInfo_SerializationId.Size(m)
}
func (m *SerializationId) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializationId.DiscardUnknown(m)
}

var xxx_messageInfo_SerializationId proto.InternalMessageInfo

func (m *SerializationId) GetSerializationId() uint32 {
	if m != nil {
		return m.SerializationId
	}
	return 0
}

type SerializedDebugGraph struct {
	Version              int32      `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Objects              []*any.Any `protobuf:"bytes,2,rep,name=objects" json:"objects,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SerializedDebugGraph) Reset()         { *m = SerializedDebugGraph{} }
func (m *SerializedDebugGraph) String() string { return proto.CompactTextString(m) }
func (*SerializedDebugGraph) ProtoMessage()    {}
func (*SerializedDebugGraph) Descriptor() ([]byte, []int) {
	return fileDescriptor_pbgraphserialization_a6ed9de92956f5b6, []int{1}
}
func (m *SerializedDebugGraph) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializedDebugGraph.Unmarshal(m, b)
}
func (m *SerializedDebugGraph) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializedDebugGraph.Marshal(b, m, deterministic)
}
func (dst *SerializedDebugGraph) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializedDebugGraph.Merge(dst, src)
}
func (m *SerializedDebugGraph) XXX_Size() int {
	return xxx_messageInfo_SerializedDebugGraph.Size(m)
}
func (m *SerializedDebugGraph) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializedDebugGraph.DiscardUnknown(m)
}

var xxx_messageInfo_SerializedDebugGraph proto.InternalMessageInfo

func (m *SerializedDebugGraph) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *SerializedDebugGraph) GetObjects() []*any.Any {
	if m != nil {
		return m.Objects
	}
	return nil
}

type SerializedGraph struct {
	Version              int32    `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Objects              [][]byte `protobuf:"bytes,2,rep,name=objects" json:"objects,omitempty"`
	TypeNames            []string `protobuf:"bytes,3,rep,name=typeNames" json:"typeNames,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SerializedGraph) Reset()         { *m = SerializedGraph{} }
func (m *SerializedGraph) String() string { return proto.CompactTextString(m) }
func (*SerializedGraph) ProtoMessage()    {}
func (*SerializedGraph) Descriptor() ([]byte, []int) {
	return fileDescriptor_pbgraphserialization_a6ed9de92956f5b6, []int{2}
}
func (m *SerializedGraph) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SerializedGraph.Unmarshal(m, b)
}
func (m *SerializedGraph) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SerializedGraph.Marshal(b, m, deterministic)
}
func (dst *SerializedGraph) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SerializedGraph.Merge(dst, src)
}
func (m *SerializedGraph) XXX_Size() int {
	return xxx_messageInfo_SerializedGraph.Size(m)
}
func (m *SerializedGraph) XXX_DiscardUnknown() {
	xxx_messageInfo_SerializedGraph.DiscardUnknown(m)
}

var xxx_messageInfo_SerializedGraph proto.InternalMessageInfo

func (m *SerializedGraph) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *SerializedGraph) GetObjects() [][]byte {
	if m != nil {
		return m.Objects
	}
	return nil
}

func (m *SerializedGraph) GetTypeNames() []string {
	if m != nil {
		return m.TypeNames
	}
	return nil
}

func init() {
	proto.RegisterType((*SerializationId)(nil), "SerializationId")
	proto.RegisterType((*SerializedDebugGraph)(nil), "SerializedDebugGraph")
	proto.RegisterType((*SerializedGraph)(nil), "SerializedGraph")
}

func init() {
	proto.RegisterFile("github.com/loomnetwork/gamechain/library/pbgraphserialization/proto/pbgraphserialization/pbgraphserialization.proto", fileDescriptor_pbgraphserialization_a6ed9de92956f5b6)
}

var fileDescriptor_pbgraphserialization_a6ed9de92956f5b6 = []byte{
	// 254 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xb1, 0x4b, 0xc4, 0x30,
	0x14, 0xc6, 0x39, 0x8b, 0x1e, 0x17, 0x95, 0x83, 0x72, 0x60, 0x4e, 0x1c, 0x4a, 0xa7, 0x4e, 0x09,
	0xe8, 0xe8, 0xa4, 0x08, 0xe2, 0xe2, 0x50, 0x37, 0x17, 0x4d, 0xda, 0x67, 0x1a, 0x6d, 0xf3, 0x42,
	0x92, 0x2a, 0xf5, 0xaf, 0x17, 0x5b, 0xea, 0x79, 0x52, 0xb9, 0x31, 0xdf, 0xfb, 0x7e, 0x1f, 0x3f,
	0x42, 0xbc, 0xd2, 0xa1, 0x6a, 0x25, 0x2b, 0xb0, 0xe1, 0x35, 0x62, 0x63, 0x20, 0x7c, 0xa0, 0x7b,
	0xe3, 0x4a, 0x34, 0x50, 0x54, 0x42, 0x1b, 0x5e, 0x6b, 0xe9, 0x84, 0xeb, 0xb8, 0x95, 0xca, 0x09,
	0x5b, 0x79, 0x70, 0x5a, 0xd4, 0xfa, 0x53, 0x04, 0x8d, 0x86, 0x5b, 0x87, 0x01, 0xff, 0x39, 0x4d,
	0x84, 0xac, 0xef, 0x9f, 0xae, 0x15, 0xa2, 0xaa, 0x61, 0xa0, 0x65, 0xfb, 0xc2, 0x85, 0xe9, 0x86,
	0x53, 0x7a, 0x49, 0x96, 0x0f, 0xbf, 0x89, 0xbb, 0x32, 0xce, 0xc8, 0xd2, 0x6f, 0x47, 0x74, 0x96,
	0xcc, 0xb2, 0xe3, 0xfc, 0x6f, 0x9c, 0x3e, 0x93, 0xd5, 0x08, 0x43, 0x79, 0x03, 0xb2, 0x55, 0xb7,
	0xdf, 0x0a, 0x31, 0x25, 0xf3, 0x77, 0x70, 0x5e, 0xa3, 0xe9, 0xc9, 0xfd, 0x7c, 0x7c, 0xc6, 0x8c,
	0xcc, 0x51, 0xbe, 0x42, 0x11, 0x3c, 0xdd, 0x4b, 0xa2, 0xec, 0xf0, 0x7c, 0xc5, 0x06, 0x37, 0x36,
	0xba, 0xb1, 0x2b, 0xd3, 0xe5, 0x63, 0x29, 0x2d, 0x36, 0x7a, 0x50, 0xee, 0x1a, 0xa7, 0xdb, 0xe3,
	0x47, 0x3f, 0x33, 0xf1, 0x19, 0x59, 0x84, 0xce, 0xc2, 0xbd, 0x68, 0xc0, 0xd3, 0x28, 0x89, 0xb2,
	0x45, 0xbe, 0x09, 0xae, 0xd7, 0x8f, 0x27, 0x53, 0x9f, 0xf7, 0x64, 0xa5, 0x3c, 0xe8, 0xb5, 0x2e,
	0xbe, 0x02, 0x00, 0x00, 0xff, 0xff, 0xa4, 0x1e, 0xfe, 0x9c, 0xb7, 0x01, 0x00, 0x00,
}
