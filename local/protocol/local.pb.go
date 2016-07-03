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
	Command
	ProcessStatus
	StatusQuery
	SnapshotCreateReq
	SnapshotCreateResp
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

type ProcessState int32

const (
	ProcessState_UNKNOWN   ProcessState = 0
	ProcessState_PENDING   ProcessState = 1
	ProcessState_RUNNING   ProcessState = 2
	ProcessState_COMPLETED ProcessState = 3
	ProcessState_FAILED    ProcessState = 4
)

var ProcessState_name = map[int32]string{
	0: "UNKNOWN",
	1: "PENDING",
	2: "RUNNING",
	3: "COMPLETED",
	4: "FAILED",
}
var ProcessState_value = map[string]int32{
	"UNKNOWN":   0,
	"PENDING":   1,
	"RUNNING":   2,
	"COMPLETED": 3,
	"FAILED":    4,
}

func (x ProcessState) String() string {
	return proto.EnumName(ProcessState_name, int32(x))
}
func (ProcessState) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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

type Command struct {
	Argv []string          `protobuf:"bytes,1,rep,name=argv" json:"argv,omitempty"`
	Env  map[string]string `protobuf:"bytes,2,rep,name=env" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// Timeout in nanoseconds
	Timeout int64 `protobuf:"varint,3,opt,name=timeout" json:"timeout,omitempty"`
}

func (m *Command) Reset()                    { *m = Command{} }
func (m *Command) String() string            { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()               {}
func (*Command) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Command) GetEnv() map[string]string {
	if m != nil {
		return m.Env
	}
	return nil
}

type ProcessStatus struct {
	RunId     string       `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
	State     ProcessState `protobuf:"varint,2,opt,name=state,enum=protocol.ProcessState" json:"state,omitempty"`
	StdoutRef string       `protobuf:"bytes,3,opt,name=stdout_ref,json=stdoutRef" json:"stdout_ref,omitempty"`
	StderrRef string       `protobuf:"bytes,4,opt,name=stderr_ref,json=stderrRef" json:"stderr_ref,omitempty"`
	ExitCode  int32        `protobuf:"varint,5,opt,name=exit_code,json=exitCode" json:"exit_code,omitempty"`
	Error     string       `protobuf:"bytes,6,opt,name=error" json:"error,omitempty"`
}

func (m *ProcessStatus) Reset()                    { *m = ProcessStatus{} }
func (m *ProcessStatus) String() string            { return proto.CompactTextString(m) }
func (*ProcessStatus) ProtoMessage()               {}
func (*ProcessStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type StatusQuery struct {
	RunId string `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
}

func (m *StatusQuery) Reset()                    { *m = StatusQuery{} }
func (m *StatusQuery) String() string            { return proto.CompactTextString(m) }
func (*StatusQuery) ProtoMessage()               {}
func (*StatusQuery) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type SnapshotCreateReq struct {
	FromDir string `protobuf:"bytes,1,opt,name=from_dir,json=fromDir" json:"from_dir,omitempty"`
}

func (m *SnapshotCreateReq) Reset()                    { *m = SnapshotCreateReq{} }
func (m *SnapshotCreateReq) String() string            { return proto.CompactTextString(m) }
func (*SnapshotCreateReq) ProtoMessage()               {}
func (*SnapshotCreateReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SnapshotCreateResp struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *SnapshotCreateResp) Reset()                    { *m = SnapshotCreateResp{} }
func (m *SnapshotCreateResp) String() string            { return proto.CompactTextString(m) }
func (*SnapshotCreateResp) ProtoMessage()               {}
func (*SnapshotCreateResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func init() {
	proto.RegisterType((*EchoRequest)(nil), "protocol.EchoRequest")
	proto.RegisterType((*EchoReply)(nil), "protocol.EchoReply")
	proto.RegisterType((*Command)(nil), "protocol.Command")
	proto.RegisterType((*ProcessStatus)(nil), "protocol.ProcessStatus")
	proto.RegisterType((*StatusQuery)(nil), "protocol.StatusQuery")
	proto.RegisterType((*SnapshotCreateReq)(nil), "protocol.SnapshotCreateReq")
	proto.RegisterType((*SnapshotCreateResp)(nil), "protocol.SnapshotCreateResp")
	proto.RegisterEnum("protocol.ProcessState", ProcessState_name, ProcessState_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for LocalScoot service

type LocalScootClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoReply, error)
	Run(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ProcessStatus, error)
	Status(ctx context.Context, in *StatusQuery, opts ...grpc.CallOption) (*ProcessStatus, error)
	SnapshotCreate(ctx context.Context, in *SnapshotCreateReq, opts ...grpc.CallOption) (*SnapshotCreateResp, error)
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

func (c *localScootClient) Run(ctx context.Context, in *Command, opts ...grpc.CallOption) (*ProcessStatus, error) {
	out := new(ProcessStatus)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Run", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localScootClient) Status(ctx context.Context, in *StatusQuery, opts ...grpc.CallOption) (*ProcessStatus, error) {
	out := new(ProcessStatus)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/Status", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *localScootClient) SnapshotCreate(ctx context.Context, in *SnapshotCreateReq, opts ...grpc.CallOption) (*SnapshotCreateResp, error) {
	out := new(SnapshotCreateResp)
	err := grpc.Invoke(ctx, "/protocol.LocalScoot/SnapshotCreate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LocalScoot service

type LocalScootServer interface {
	Echo(context.Context, *EchoRequest) (*EchoReply, error)
	Run(context.Context, *Command) (*ProcessStatus, error)
	Status(context.Context, *StatusQuery) (*ProcessStatus, error)
	SnapshotCreate(context.Context, *SnapshotCreateReq) (*SnapshotCreateResp, error)
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

func _LocalScoot_Run_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Run(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Run",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Run(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalScoot_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).Status(ctx, req.(*StatusQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocalScoot_SnapshotCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SnapshotCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocalScootServer).SnapshotCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protocol.LocalScoot/SnapshotCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocalScootServer).SnapshotCreate(ctx, req.(*SnapshotCreateReq))
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
		{
			MethodName: "Run",
			Handler:    _LocalScoot_Run_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _LocalScoot_Status_Handler,
		},
		{
			MethodName: "SnapshotCreate",
			Handler:    _LocalScoot_SnapshotCreate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("local.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x52, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x8d, 0xe3, 0x7c, 0x79, 0x42, 0xa3, 0x74, 0xa0, 0x60, 0x52, 0x10, 0x60, 0x71, 0x40, 0x08,
	0xe5, 0x90, 0x22, 0x84, 0x7a, 0x43, 0x89, 0x41, 0x11, 0xa9, 0x1b, 0x36, 0x54, 0x1c, 0x23, 0x63,
	0x6f, 0x5a, 0x8b, 0xc4, 0xeb, 0xee, 0xae, 0x23, 0x72, 0xe5, 0x7f, 0xf0, 0x73, 0xf8, 0x5f, 0xec,
	0xda, 0x4e, 0x63, 0x42, 0xe1, 0xe4, 0x79, 0xf3, 0xde, 0xcc, 0xbc, 0x1d, 0x0f, 0xb4, 0x97, 0x2c,
	0xf0, 0x97, 0xfd, 0x84, 0x33, 0xc9, 0xb0, 0x95, 0x7d, 0x02, 0xb6, 0x74, 0x9e, 0x41, 0xdb, 0x0d,
	0xae, 0x18, 0xa1, 0xd7, 0x29, 0x15, 0x12, 0x11, 0x6a, 0x49, 0x14, 0x5f, 0xda, 0xc6, 0x53, 0xe3,
	0x85, 0x45, 0xb2, 0xd8, 0x79, 0x02, 0x56, 0x2e, 0x49, 0x96, 0x9b, 0x4c, 0xc0, 0x4a, 0x02, 0x15,
	0x3b, 0x3f, 0x0d, 0x68, 0x0e, 0xd9, 0x6a, 0xe5, 0xc7, 0xa1, 0xe6, 0x7d, 0x7e, 0xb9, 0x56, 0xbc,
	0xa9, 0x79, 0x1d, 0xe3, 0x2b, 0x30, 0x69, 0xbc, 0xb6, 0xab, 0x2a, 0xd5, 0x1e, 0xf4, 0xfa, 0xdb,
	0xd9, 0xfd, 0xa2, 0xa6, 0xef, 0xc6, 0x6b, 0x37, 0x96, 0x7c, 0x43, 0xb4, 0x0c, 0x6d, 0x68, 0xca,
	0x68, 0x45, 0x59, 0x2a, 0x6d, 0x53, 0x0d, 0x31, 0xc9, 0x16, 0xf6, 0xde, 0x40, 0x6b, 0x2b, 0xc5,
	0x2e, 0x98, 0xdf, 0xe8, 0xa6, 0xb0, 0xa1, 0x43, 0xbc, 0x07, 0xf5, 0xb5, 0xbf, 0x4c, 0xa9, 0x9a,
	0xa3, 0x73, 0x39, 0x38, 0xad, 0xbe, 0x35, 0x9c, 0x5f, 0x06, 0x1c, 0x4c, 0x39, 0x0b, 0xa8, 0x10,
	0x33, 0xe9, 0xcb, 0x54, 0xe0, 0x11, 0x34, 0x78, 0x1a, 0xcf, 0xa3, 0xb0, 0x68, 0x50, 0x57, 0x68,
	0x1c, 0x2a, 0xa3, 0x75, 0xa1, 0x04, 0x79, 0x8b, 0xce, 0xe0, 0xfe, 0xce, 0x6a, 0xa9, 0x9c, 0x92,
	0x5c, 0x84, 0x8f, 0x01, 0x84, 0x0c, 0x95, 0xb1, 0x39, 0xa7, 0x8b, 0xcc, 0xab, 0x45, 0xac, 0x3c,
	0x43, 0xe8, 0xa2, 0xa0, 0x29, 0xe7, 0x19, 0x5d, 0xbb, 0xa1, 0x55, 0x46, 0xd3, 0xc7, 0x60, 0xd1,
	0xef, 0x91, 0x9c, 0x07, 0x2c, 0xa4, 0x76, 0x5d, 0xb1, 0x75, 0xd2, 0xd2, 0x89, 0xa1, 0xc2, 0xfa,
	0x2d, 0x4a, 0xc6, 0xb8, 0xdd, 0xc8, 0xed, 0x65, 0xc0, 0x79, 0x0e, 0xed, 0xdc, 0xff, 0xa7, 0x94,
	0xaa, 0x15, 0xdc, 0xfe, 0x08, 0xa7, 0x0f, 0x87, 0xb3, 0xd8, 0x4f, 0xc4, 0x15, 0x93, 0x43, 0x4e,
	0xb5, 0x5f, 0x7a, 0x8d, 0x0f, 0xa1, 0xb5, 0xe0, 0x6c, 0x35, 0x0f, 0x23, 0x5e, 0xa8, 0x9b, 0x1a,
	0x8f, 0x22, 0xdd, 0x15, 0xf7, 0xf5, 0x22, 0xc1, 0x0e, 0x54, 0x6f, 0x1a, 0xab, 0xe8, 0xe5, 0x14,
	0xee, 0x94, 0x77, 0x80, 0x6d, 0x68, 0x5e, 0x78, 0x1f, 0xbd, 0xf3, 0x2f, 0x5e, 0xb7, 0xa2, 0xc1,
	0xd4, 0xf5, 0x46, 0x63, 0xef, 0x43, 0xd7, 0xd0, 0x80, 0x5c, 0x78, 0x9e, 0x06, 0x55, 0x3c, 0x00,
	0x6b, 0x78, 0x7e, 0x36, 0x9d, 0xb8, 0x9f, 0xdd, 0x51, 0xd7, 0x44, 0x80, 0xc6, 0xfb, 0x77, 0xe3,
	0x89, 0x8a, 0x6b, 0x83, 0x1f, 0x55, 0x80, 0x89, 0xbe, 0xc9, 0x59, 0xc0, 0x98, 0xc4, 0xd7, 0x50,
	0xd3, 0x57, 0x86, 0x47, 0xbb, 0xa5, 0x97, 0x0e, 0xb3, 0x77, 0x77, 0x3f, 0xad, 0x8e, 0xd1, 0xa9,
	0xe0, 0x09, 0x98, 0x24, 0x8d, 0xf1, 0xf0, 0xaf, 0xa3, 0xea, 0x3d, 0xb8, 0xf5, 0xe7, 0xa5, 0x42,
	0x15, 0x9d, 0x42, 0x63, 0x7b, 0x07, 0x3b, 0x51, 0x69, 0xb3, 0xff, 0xab, 0x3d, 0x83, 0xce, 0x9f,
	0xdb, 0xc2, 0xe3, 0x52, 0x8f, 0xfd, 0xbd, 0xf7, 0x1e, 0xfd, 0x9b, 0x14, 0x89, 0x53, 0xf9, 0xda,
	0xc8, 0xe8, 0x93, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x4a, 0xc4, 0xca, 0x9e, 0x03, 0x00,
	0x00,
}
