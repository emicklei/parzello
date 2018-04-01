// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/queue/parcel.proto

/*
Package queue is a generated protocol buffer package.

It is generated from these files:
	v1/queue/parcel.proto

It has these top-level messages:
	Parcel
*/
package queue

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import v1 "github.com/emicklei/parcello/v1"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Parcel struct {
	Envelope *v1.Envelope `protobuf:"bytes,1,opt,name=envelope" json:"envelope,omitempty"`
	// in flight data
	EnqueueTime  *v1.Timestamp `protobuf:"bytes,2,opt,name=enqueueTime" json:"enqueueTime,omitempty"`
	RequeueCount uint32        `protobuf:"varint,3,opt,name=requeueCount" json:"requeueCount,omitempty"`
	RetryCount   uint32        `protobuf:"varint,4,opt,name=retryCount" json:"retryCount,omitempty"`
}

func (m *Parcel) Reset()                    { *m = Parcel{} }
func (m *Parcel) String() string            { return proto.CompactTextString(m) }
func (*Parcel) ProtoMessage()               {}
func (*Parcel) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Parcel) GetEnvelope() *v1.Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

func (m *Parcel) GetEnqueueTime() *v1.Timestamp {
	if m != nil {
		return m.EnqueueTime
	}
	return nil
}

func (m *Parcel) GetRequeueCount() uint32 {
	if m != nil {
		return m.RequeueCount
	}
	return 0
}

func (m *Parcel) GetRetryCount() uint32 {
	if m != nil {
		return m.RetryCount
	}
	return 0
}

func init() {
	proto.RegisterType((*Parcel)(nil), "queue.Parcel")
}

func init() { proto.RegisterFile("v1/queue/parcel.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2d, 0x33, 0xd4, 0x2f,
	0x2c, 0x4d, 0x2d, 0x4d, 0xd5, 0x2f, 0x48, 0x2c, 0x4a, 0x4e, 0xcd, 0xd1, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0x05, 0x8b, 0x49, 0x09, 0x96, 0x19, 0x42, 0xc5, 0x73, 0xf2, 0x21, 0x32, 0x4a,
	0xcb, 0x19, 0xb9, 0xd8, 0x02, 0xc0, 0x42, 0x42, 0x1a, 0x5c, 0x1c, 0xa9, 0x79, 0x65, 0xa9, 0x39,
	0xf9, 0x05, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x3c, 0x7a, 0x65, 0x86, 0x7a, 0xae,
	0x50, 0xb1, 0x20, 0xb8, 0xac, 0x90, 0x3e, 0x17, 0x77, 0x6a, 0x1e, 0xd8, 0xc8, 0x90, 0xcc, 0xdc,
	0x54, 0x09, 0x26, 0xb0, 0x62, 0x5e, 0x90, 0x62, 0x10, 0xbf, 0xb8, 0x24, 0x31, 0xb7, 0x20, 0x08,
	0x59, 0x85, 0x90, 0x12, 0x17, 0x4f, 0x51, 0x2a, 0x98, 0xeb, 0x9c, 0x5f, 0x9a, 0x57, 0x22, 0xc1,
	0xac, 0xc0, 0xa8, 0xc1, 0x1b, 0x84, 0x22, 0x26, 0x24, 0xc7, 0xc5, 0x55, 0x94, 0x5a, 0x52, 0x54,
	0x09, 0x51, 0xc1, 0x02, 0x56, 0x81, 0x24, 0xe2, 0xa4, 0x1e, 0xa5, 0x9a, 0x9e, 0x59, 0x92, 0x51,
	0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x9f, 0x9a, 0x9b, 0x99, 0x9c, 0x9d, 0x93, 0x9a, 0x09, 0xf7,
	0x8f, 0x3e, 0xcc, 0xe7, 0x49, 0x6c, 0x60, 0x9f, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x33,
	0x6e, 0x2a, 0xed, 0x0c, 0x01, 0x00, 0x00,
}