// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package language_agent_v2

import (
	context "context"
	common "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs/skywalking/v2/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// JVMMetricReportServiceClient is the client API for JVMMetricReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JVMMetricReportServiceClient interface {
	Collect(ctx context.Context, in *JVMMetricCollection, opts ...grpc.CallOption) (*common.Commands, error)
}

type jVMMetricReportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewJVMMetricReportServiceClient(cc grpc.ClientConnInterface) JVMMetricReportServiceClient {
	return &jVMMetricReportServiceClient{cc}
}

func (c *jVMMetricReportServiceClient) Collect(ctx context.Context, in *JVMMetricCollection, opts ...grpc.CallOption) (*common.Commands, error) {
	out := new(common.Commands)
	err := c.cc.Invoke(ctx, "/skywalking.v2.JVMMetricReportService/collect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// JVMMetricReportServiceServer is the server API for JVMMetricReportService service.
// All implementations must embed UnimplementedJVMMetricReportServiceServer
// for forward compatibility
type JVMMetricReportServiceServer interface {
	Collect(context.Context, *JVMMetricCollection) (*common.Commands, error)
	mustEmbedUnimplementedJVMMetricReportServiceServer()
}

// UnimplementedJVMMetricReportServiceServer must be embedded to have forward compatible implementations.
type UnimplementedJVMMetricReportServiceServer struct {
}

func (UnimplementedJVMMetricReportServiceServer) Collect(context.Context, *JVMMetricCollection) (*common.Commands, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Collect not implemented")
}
func (UnimplementedJVMMetricReportServiceServer) mustEmbedUnimplementedJVMMetricReportServiceServer() {
}

// UnsafeJVMMetricReportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JVMMetricReportServiceServer will
// result in compilation errors.
type UnsafeJVMMetricReportServiceServer interface {
	mustEmbedUnimplementedJVMMetricReportServiceServer()
}

func RegisterJVMMetricReportServiceServer(s grpc.ServiceRegistrar, srv JVMMetricReportServiceServer) {
	s.RegisterService(&_JVMMetricReportService_serviceDesc, srv)
}

func _JVMMetricReportService_Collect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JVMMetricCollection)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JVMMetricReportServiceServer).Collect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/skywalking.v2.JVMMetricReportService/collect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JVMMetricReportServiceServer).Collect(ctx, req.(*JVMMetricCollection))
	}
	return interceptor(ctx, in, info, handler)
}

var _JVMMetricReportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "skywalking.v2.JVMMetricReportService",
	HandlerType: (*JVMMetricReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "collect",
			Handler:    _JVMMetricReportService_Collect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugins/inputs/skywalking/v2/language-agent-v2/JVMMetric.proto",
}
