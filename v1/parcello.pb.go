// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v1/parcello.proto

/*
Package v1 is a generated protocol buffer package.

It is generated from these files:
	v1/parcello.proto

It has these top-level messages:
	DeliverRequest
	DeliverResponse
	Envelope
	Timestamp
*/
package v1

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

type DeliverRequest struct {
	Envelope *Envelope `protobuf:"bytes,1,opt,name=envelope" json:"envelope,omitempty"`
}

func (m *DeliverRequest) Reset()                    { *m = DeliverRequest{} }
func (m *DeliverRequest) String() string            { return proto.CompactTextString(m) }
func (*DeliverRequest) ProtoMessage()               {}
func (*DeliverRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DeliverRequest) GetEnvelope() *Envelope {
	if m != nil {
		return m.Envelope
	}
	return nil
}

type DeliverResponse struct {
	ErrorMessage string `protobuf:"bytes,1,opt,name=errorMessage" json:"errorMessage,omitempty"`
}

func (m *DeliverResponse) Reset()                    { *m = DeliverResponse{} }
func (m *DeliverResponse) String() string            { return proto.CompactTextString(m) }
func (*DeliverResponse) ProtoMessage()               {}
func (*DeliverResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DeliverResponse) GetErrorMessage() string {
	if m != nil {
		return m.ErrorMessage
	}
	return ""
}

type Envelope struct {
	// fixed data
	Payload          []byte            `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	Attributes       map[string]string `protobuf:"bytes,2,rep,name=attributes" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	DestinationTopic string            `protobuf:"bytes,3,opt,name=destinationTopic" json:"destinationTopic,omitempty"`
	UndeliveredTopic string            `protobuf:"bytes,4,opt,name=undeliveredTopic" json:"undeliveredTopic,omitempty"`
	PublishAfter     *Timestamp        `protobuf:"bytes,5,opt,name=publishAfter" json:"publishAfter,omitempty"`
}

func (m *Envelope) Reset()                    { *m = Envelope{} }
func (m *Envelope) String() string            { return proto.CompactTextString(m) }
func (*Envelope) ProtoMessage()               {}
func (*Envelope) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Envelope) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Envelope) GetAttributes() map[string]string {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *Envelope) GetDestinationTopic() string {
	if m != nil {
		return m.DestinationTopic
	}
	return ""
}

func (m *Envelope) GetUndeliveredTopic() string {
	if m != nil {
		return m.UndeliveredTopic
	}
	return ""
}

func (m *Envelope) GetPublishAfter() *Timestamp {
	if m != nil {
		return m.PublishAfter
	}
	return nil
}

type Timestamp struct {
	Seconds uint64 `protobuf:"varint,1,opt,name=seconds" json:"seconds,omitempty"`
}

func (m *Timestamp) Reset()                    { *m = Timestamp{} }
func (m *Timestamp) String() string            { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()               {}
func (*Timestamp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Timestamp) GetSeconds() uint64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func init() {
	proto.RegisterType((*DeliverRequest)(nil), "v1.DeliverRequest")
	proto.RegisterType((*DeliverResponse)(nil), "v1.DeliverResponse")
	proto.RegisterType((*Envelope)(nil), "v1.Envelope")
	proto.RegisterType((*Timestamp)(nil), "v1.Timestamp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DeliveryService service

type DeliveryServiceClient interface {
	Deliver(ctx context.Context, in *DeliverRequest, opts ...grpc.CallOption) (*DeliverResponse, error)
}

type deliveryServiceClient struct {
	cc *grpc.ClientConn
}

func NewDeliveryServiceClient(cc *grpc.ClientConn) DeliveryServiceClient {
	return &deliveryServiceClient{cc}
}

func (c *deliveryServiceClient) Deliver(ctx context.Context, in *DeliverRequest, opts ...grpc.CallOption) (*DeliverResponse, error) {
	out := new(DeliverResponse)
	err := grpc.Invoke(ctx, "/v1.DeliveryService/Deliver", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DeliveryService service

type DeliveryServiceServer interface {
	Deliver(context.Context, *DeliverRequest) (*DeliverResponse, error)
}

func RegisterDeliveryServiceServer(s *grpc.Server, srv DeliveryServiceServer) {
	s.RegisterService(&_DeliveryService_serviceDesc, srv)
}

func _DeliveryService_Deliver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeliverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeliveryServiceServer).Deliver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.DeliveryService/Deliver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeliveryServiceServer).Deliver(ctx, req.(*DeliverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DeliveryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.DeliveryService",
	HandlerType: (*DeliveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deliver",
			Handler:    _DeliveryService_Deliver_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/parcello.proto",
}

func init() { proto.RegisterFile("v1/parcello.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x92, 0x41, 0xcf, 0xd2, 0x40,
	0x10, 0x86, 0x43, 0x01, 0x81, 0xa1, 0x0a, 0xae, 0x1e, 0x1a, 0x62, 0x22, 0x36, 0x31, 0x21, 0x1e,
	0x4a, 0x5a, 0x63, 0x62, 0x88, 0x1e, 0x30, 0x72, 0xf4, 0x52, 0x39, 0x79, 0xdb, 0xb6, 0x23, 0x6c,
	0xd8, 0x76, 0xd7, 0xdd, 0xed, 0x26, 0xfd, 0x4f, 0xfe, 0xc8, 0x2f, 0x6d, 0xa1, 0xa1, 0x7c, 0xb7,
	0x9d, 0x67, 0xde, 0x99, 0xbc, 0x3b, 0x33, 0xf0, 0xda, 0x86, 0x5b, 0x49, 0x55, 0x8a, 0x9c, 0x8b,
	0x40, 0x2a, 0x61, 0x04, 0x71, 0x6c, 0xe8, 0xef, 0xe0, 0xd5, 0x4f, 0xe4, 0xcc, 0xa2, 0x8a, 0xf1,
	0x5f, 0x89, 0xda, 0x90, 0x0d, 0x4c, 0xb1, 0xb0, 0xc8, 0x85, 0x44, 0x6f, 0xb0, 0x1e, 0x6c, 0xe6,
	0x91, 0x1b, 0xd8, 0x30, 0x38, 0x5c, 0x59, 0xdc, 0x65, 0xfd, 0x2f, 0xb0, 0xe8, 0x6a, 0xb5, 0x14,
	0x85, 0x46, 0xe2, 0x83, 0x8b, 0x4a, 0x09, 0xf5, 0x0b, 0xb5, 0xa6, 0xa7, 0xb6, 0xc1, 0x2c, 0xee,
	0x31, 0xff, 0xbf, 0x03, 0xd3, 0x5b, 0x37, 0xe2, 0xc1, 0x44, 0xd2, 0x8a, 0x0b, 0x9a, 0x35, 0x5a,
	0x37, 0xbe, 0x85, 0xe4, 0x1b, 0x00, 0x35, 0x46, 0xb1, 0xa4, 0x34, 0xa8, 0x3d, 0x67, 0x3d, 0xdc,
	0xcc, 0xa3, 0x77, 0xf7, 0x4e, 0x82, 0x7d, 0x97, 0x3e, 0x14, 0x46, 0x55, 0xf1, 0x9d, 0x9e, 0x7c,
	0x82, 0x65, 0x86, 0xda, 0xb0, 0x82, 0x1a, 0x26, 0x8a, 0xa3, 0x90, 0x2c, 0xf5, 0x86, 0x8d, 0x99,
	0x67, 0xbc, 0xd6, 0x96, 0x45, 0xd6, 0xfe, 0x04, 0xb3, 0x56, 0x3b, 0x6a, 0xb5, 0x8f, 0x9c, 0x84,
	0xe0, 0xca, 0x32, 0xe1, 0x4c, 0x9f, 0xf7, 0x7f, 0x0d, 0x2a, 0x6f, 0xdc, 0x4c, 0xe8, 0x65, 0xed,
	0xeb, 0xc8, 0x72, 0xd4, 0x86, 0xe6, 0x32, 0xee, 0x49, 0x56, 0xdf, 0x61, 0xf1, 0xe0, 0x94, 0x2c,
	0x61, 0x78, 0xc1, 0xea, 0x3a, 0x9d, 0xfa, 0x49, 0xde, 0xc2, 0xd8, 0x52, 0x5e, 0xa2, 0xe7, 0x34,
	0xac, 0x0d, 0x76, 0xce, 0xd7, 0x81, 0xff, 0x11, 0x66, 0x5d, 0xe7, 0x7a, 0x5c, 0x1a, 0x53, 0x51,
	0x64, 0xba, 0x29, 0x1e, 0xc5, 0xb7, 0x30, 0x3a, 0x74, 0xcb, 0xa8, 0x7e, 0xa3, 0xb2, 0x2c, 0x45,
	0x12, 0xc1, 0xe4, 0x8a, 0x08, 0xa9, 0x0d, 0xf6, 0x17, 0xbd, 0x7a, 0xd3, 0x63, 0xed, 0x02, 0x7f,
	0x7c, 0xf8, 0xf3, 0xfe, 0xc4, 0xcc, 0xb9, 0x4c, 0x82, 0x54, 0xe4, 0x5b, 0xcc, 0x59, 0x7a, 0xe1,
	0xc8, 0xba, 0xcb, 0xd9, 0xda, 0x30, 0x79, 0xd1, 0x5c, 0xcf, 0xe7, 0xa7, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x8a, 0x60, 0xe8, 0xec, 0x52, 0x02, 0x00, 0x00,
}
