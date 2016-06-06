// Code generated by protoc-gen-go.
// source: local.proto
// DO NOT EDIT!

/*
Package protocol is a generated protocol buffer package.

It is generated from these files:
	local.proto

It has these top-level messages:
	EchoRequest
	EchoReply
*/
package protocol

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

type EchoRequest struct {
	Ping string `protobuf:"bytes,1,opt,name=ping" json:"ping,omitempty"`
}

func (m *EchoRequest) Reset()                    { *m = EchoRequest{} }
func (m *EchoRequest) String() string            { return proto.CompactTextString(m) }
func (*EchoRequest) ProtoMessage()               {}
func (*EchoRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type EchoReply struct {
	Pong string `protobuf:"bytes,1,opt,name=pong" json:"pong,omitempty"`
}

func (m *EchoReply) Reset()                    { *m = EchoReply{} }
func (m *EchoReply) String() string            { return proto.CompactTextString(m) }
func (*EchoReply) ProtoMessage()               {}
func (*EchoReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*EchoRequest)(nil), "protocol.EchoRequest")
	proto.RegisterType((*EchoReply)(nil), "protocol.EchoReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion2

// Client API for LocalScoot service

type LocalScootClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
}

type localScootClient struct {
	cc *grpc.ClientConn
}

func NewLocalScootClient(cc *grpc.ClientConn) LocalScootClient {
	return &localScootClient{cc}
}

func (c *localScootClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error) {
	out := new(EchoReply)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LocalScoot service

type LocalScootServer interface {
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
}

func RegisterLocalScootServer(s *grpc.Server, srv LocalScootServer) {
	s.RegisterService(&_LocalScoot_serviceDesc, srv)
}

func _LocalScoot_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LocalScoot_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protocol.LocalScoot",
	HandlerType: (*LocalScootServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _LocalScoot_Echo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xc9, 0x4f, 0x4e,
	0xcc, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0xc9, 0xf9, 0x39, 0x4a, 0x8a,
	0x5c, 0xdc, 0xae, 0xc9, 0x19, 0xf9, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0x42, 0x5c,
	0x2c, 0x05, 0x99, 0x79, 0xe9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x60, 0xb6, 0x92, 0x3c,
	0x17, 0x27, 0x44, 0x49, 0x41, 0x4e, 0x25, 0x58, 0x41, 0x3e, 0x92, 0x02, 0x20, 0xdb, 0xc8, 0x89,
	0x8b, 0xcb, 0x07, 0x64, 0x78, 0x70, 0x72, 0x7e, 0x7e, 0x89, 0x90, 0x09, 0x17, 0x0b, 0x48, 0xb9,
	0x90, 0xa8, 0x1e, 0xcc, 0x12, 0x3d, 0x24, 0x1b, 0xa4, 0x84, 0xd1, 0x85, 0x81, 0xa6, 0x2a, 0x31,
	0x24, 0xb1, 0x81, 0x45, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x65, 0x1b, 0x8c, 0x52, 0xa7,
	0x00, 0x00, 0x00,
}
