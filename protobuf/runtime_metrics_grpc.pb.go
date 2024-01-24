// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: protobuf/runtime_metrics.proto

package pb

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

const (
	MetricsService_UpdateMetrics_FullMethodName = "/pb.MetricsService/UpdateMetrics"
	MetricsService_UpdateMetric_FullMethodName  = "/pb.MetricsService/UpdateMetric"
	MetricsService_GetMetric_FullMethodName     = "/pb.MetricsService/GetMetric"
	MetricsService_PingDB_FullMethodName        = "/pb.MetricsService/PingDB"
)

// MetricsServiceClient is the client API for MetricsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MetricsServiceClient interface {
	UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error)
	UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error)
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	PingDB(ctx context.Context, in *PingDBRequest, opts ...grpc.CallOption) (*PingDBResponse, error)
}

type metricsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMetricsServiceClient(cc grpc.ClientConnInterface) MetricsServiceClient {
	return &metricsServiceClient{cc}
}

func (c *metricsServiceClient) UpdateMetrics(ctx context.Context, in *UpdateMetricsRequest, opts ...grpc.CallOption) (*UpdateMetricsResponse, error) {
	out := new(UpdateMetricsResponse)
	err := c.cc.Invoke(ctx, MetricsService_UpdateMetrics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*UpdateMetricResponse, error) {
	out := new(UpdateMetricResponse)
	err := c.cc.Invoke(ctx, MetricsService_UpdateMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, MetricsService_GetMetric_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricsServiceClient) PingDB(ctx context.Context, in *PingDBRequest, opts ...grpc.CallOption) (*PingDBResponse, error) {
	out := new(PingDBResponse)
	err := c.cc.Invoke(ctx, MetricsService_PingDB_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricsServiceServer is the server API for MetricsService service.
// All implementations must embed UnimplementedMetricsServiceServer
// for forward compatibility
type MetricsServiceServer interface {
	UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error)
	UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error)
	GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
	PingDB(context.Context, *PingDBRequest) (*PingDBResponse, error)
	mustEmbedUnimplementedMetricsServiceServer()
}

// UnimplementedMetricsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMetricsServiceServer struct {
}

func (UnimplementedMetricsServiceServer) UpdateMetrics(context.Context, *UpdateMetricsRequest) (*UpdateMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetrics not implemented")
}
func (UnimplementedMetricsServiceServer) UpdateMetric(context.Context, *UpdateMetricRequest) (*UpdateMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedMetricsServiceServer) GetMetric(context.Context, *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedMetricsServiceServer) PingDB(context.Context, *PingDBRequest) (*PingDBResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PingDB not implemented")
}
func (UnimplementedMetricsServiceServer) mustEmbedUnimplementedMetricsServiceServer() {}

// UnsafeMetricsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MetricsServiceServer will
// result in compilation errors.
type UnsafeMetricsServiceServer interface {
	mustEmbedUnimplementedMetricsServiceServer()
}

func RegisterMetricsServiceServer(s grpc.ServiceRegistrar, srv MetricsServiceServer) {
	s.RegisterService(&MetricsService_ServiceDesc, srv)
}

func _MetricsService_UpdateMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).UpdateMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_UpdateMetrics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).UpdateMetrics(ctx, req.(*UpdateMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_UpdateMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).UpdateMetric(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_GetMetric_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricsService_PingDB_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingDBRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricsServiceServer).PingDB(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MetricsService_PingDB_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricsServiceServer).PingDB(ctx, req.(*PingDBRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MetricsService_ServiceDesc is the grpc.ServiceDesc for MetricsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MetricsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.MetricsService",
	HandlerType: (*MetricsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateMetrics",
			Handler:    _MetricsService_UpdateMetrics_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _MetricsService_UpdateMetric_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _MetricsService_GetMetric_Handler,
		},
		{
			MethodName: "PingDB",
			Handler:    _MetricsService_PingDB_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/runtime_metrics.proto",
}