// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: jobs/felek.proto

package jobs

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// JobsClient is the client API for Jobs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobsClient interface {
	Start(ctx context.Context, in *JobStartRequest, opts ...grpc.CallOption) (*JobStatus, error)
	Stop(ctx context.Context, in *JobID, opts ...grpc.CallOption) (*JobStatus, error)
	Status(ctx context.Context, in *JobID, opts ...grpc.CallOption) (*JobStatus, error)
	Stdout(ctx context.Context, in *LogsRequest, opts ...grpc.CallOption) (Jobs_StdoutClient, error)
	Stderr(ctx context.Context, in *LogsRequest, opts ...grpc.CallOption) (Jobs_StderrClient, error)
}

type jobsClient struct {
	cc grpc.ClientConnInterface
}

func NewJobsClient(cc grpc.ClientConnInterface) JobsClient {
	return &jobsClient{cc}
}

func (c *jobsClient) Start(ctx context.Context, in *JobStartRequest, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/jobs.Jobs/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobsClient) Stop(ctx context.Context, in *JobID, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/jobs.Jobs/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobsClient) Status(ctx context.Context, in *JobID, opts ...grpc.CallOption) (*JobStatus, error) {
	out := new(JobStatus)
	err := c.cc.Invoke(ctx, "/jobs.Jobs/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobsClient) Stdout(ctx context.Context, in *LogsRequest, opts ...grpc.CallOption) (Jobs_StdoutClient, error) {
	stream, err := c.cc.NewStream(ctx, &Jobs_ServiceDesc.Streams[0], "/jobs.Jobs/Stdout", opts...)
	if err != nil {
		return nil, err
	}
	x := &jobsStdoutClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Jobs_StdoutClient interface {
	Recv() (*LogLine, error)
	grpc.ClientStream
}

type jobsStdoutClient struct {
	grpc.ClientStream
}

func (x *jobsStdoutClient) Recv() (*LogLine, error) {
	m := new(LogLine)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *jobsClient) Stderr(ctx context.Context, in *LogsRequest, opts ...grpc.CallOption) (Jobs_StderrClient, error) {
	stream, err := c.cc.NewStream(ctx, &Jobs_ServiceDesc.Streams[1], "/jobs.Jobs/Stderr", opts...)
	if err != nil {
		return nil, err
	}
	x := &jobsStderrClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Jobs_StderrClient interface {
	Recv() (*LogLine, error)
	grpc.ClientStream
}

type jobsStderrClient struct {
	grpc.ClientStream
}

func (x *jobsStderrClient) Recv() (*LogLine, error) {
	m := new(LogLine)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JobsServer is the server API for Jobs service.
// All implementations must embed UnimplementedJobsServer
// for forward compatibility
type JobsServer interface {
	Start(context.Context, *JobStartRequest) (*JobStatus, error)
	Stop(context.Context, *JobID) (*JobStatus, error)
	Status(context.Context, *JobID) (*JobStatus, error)
	Stdout(*LogsRequest, Jobs_StdoutServer) error
	Stderr(*LogsRequest, Jobs_StderrServer) error
	mustEmbedUnimplementedJobsServer()
}

// UnimplementedJobsServer must be embedded to have forward compatible implementations.
type UnimplementedJobsServer struct {
}

func (UnimplementedJobsServer) Start(context.Context, *JobStartRequest) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
}
func (UnimplementedJobsServer) Stop(context.Context, *JobID) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedJobsServer) Status(context.Context, *JobID) (*JobStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedJobsServer) Stdout(*LogsRequest, Jobs_StdoutServer) error {
	return status.Errorf(codes.Unimplemented, "method Stdout not implemented")
}
func (UnimplementedJobsServer) Stderr(*LogsRequest, Jobs_StderrServer) error {
	return status.Errorf(codes.Unimplemented, "method Stderr not implemented")
}
func (UnimplementedJobsServer) mustEmbedUnimplementedJobsServer() {}

// UnsafeJobsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobsServer will
// result in compilation errors.
type UnsafeJobsServer interface {
	mustEmbedUnimplementedJobsServer()
}

func RegisterJobsServer(s grpc.ServiceRegistrar, srv JobsServer) {
	s.RegisterService(&Jobs_ServiceDesc, srv)
}

func _Jobs_Start_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobStartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServer).Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobs.Jobs/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServer).Start(ctx, req.(*JobStartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobs_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobs.Jobs/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServer).Stop(ctx, req.(*JobID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobs_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobsServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/jobs.Jobs/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobsServer).Status(ctx, req.(*JobID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Jobs_Stdout_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LogsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobsServer).Stdout(m, &jobsStdoutServer{stream})
}

type Jobs_StdoutServer interface {
	Send(*LogLine) error
	grpc.ServerStream
}

type jobsStdoutServer struct {
	grpc.ServerStream
}

func (x *jobsStdoutServer) Send(m *LogLine) error {
	return x.ServerStream.SendMsg(m)
}

func _Jobs_Stderr_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LogsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobsServer).Stderr(m, &jobsStderrServer{stream})
}

type Jobs_StderrServer interface {
	Send(*LogLine) error
	grpc.ServerStream
}

type jobsStderrServer struct {
	grpc.ServerStream
}

func (x *jobsStderrServer) Send(m *LogLine) error {
	return x.ServerStream.SendMsg(m)
}

// Jobs_ServiceDesc is the grpc.ServiceDesc for Jobs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Jobs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "jobs.Jobs",
	HandlerType: (*JobsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Start",
			Handler:    _Jobs_Start_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Jobs_Stop_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Jobs_Status_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stdout",
			Handler:       _Jobs_Stdout_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Stderr",
			Handler:       _Jobs_Stderr_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "jobs/felek.proto",
}
