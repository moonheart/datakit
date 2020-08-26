//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: register/InstancePing.proto

package register

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	common "gitlab.jiagouyun.com/cloudcare-tools/datakit/plugins/inputs/traceSkywalking/v2/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ServiceInstancePingPkg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServiceInstanceId   int32  `protobuf:"varint,1,opt,name=serviceInstanceId,proto3" json:"serviceInstanceId,omitempty"`
	Time                int64  `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	ServiceInstanceUUID string `protobuf:"bytes,3,opt,name=serviceInstanceUUID,proto3" json:"serviceInstanceUUID,omitempty"`
}

func (x *ServiceInstancePingPkg) Reset() {
	*x = ServiceInstancePingPkg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_register_InstancePing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceInstancePingPkg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceInstancePingPkg) ProtoMessage() {}

func (x *ServiceInstancePingPkg) ProtoReflect() protoreflect.Message {
	mi := &file_register_InstancePing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceInstancePingPkg.ProtoReflect.Descriptor instead.
func (*ServiceInstancePingPkg) Descriptor() ([]byte, []int) {
	return file_register_InstancePing_proto_rawDescGZIP(), []int{0}
}

func (x *ServiceInstancePingPkg) GetServiceInstanceId() int32 {
	if x != nil {
		return x.ServiceInstanceId
	}
	return 0
}

func (x *ServiceInstancePingPkg) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *ServiceInstancePingPkg) GetServiceInstanceUUID() string {
	if x != nil {
		return x.ServiceInstanceUUID
	}
	return ""
}

var File_register_InstancePing_proto protoreflect.FileDescriptor

var file_register_InstancePing_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x8c, 0x01, 0x0a, 0x16, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x6b, 0x67, 0x12, 0x2c, 0x0a,
	0x11, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12,
	0x30, 0x0a, 0x13, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x55, 0x55, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x55, 0x55, 0x49,
	0x44, 0x32, 0x45, 0x0a, 0x13, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x2e, 0x0a, 0x06, 0x64, 0x6f, 0x50, 0x69,
	0x6e, 0x67, 0x12, 0x17, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x50, 0x69, 0x6e, 0x67, 0x50, 0x6b, 0x67, 0x1a, 0x09, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x22, 0x00, 0x42, 0xa7, 0x01, 0x0a, 0x2d, 0x6f, 0x72, 0x67,
	0x2e, 0x61, 0x70, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x73, 0x6b, 0x79, 0x77, 0x61, 0x6c, 0x6b, 0x69,
	0x6e, 0x67, 0x2e, 0x61, 0x70, 0x6d, 0x2e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x2e, 0x76, 0x32, 0x50, 0x01, 0x5a, 0x57, 0x67, 0x69,
	0x74, 0x6c, 0x61, 0x62, 0x2e, 0x6a, 0x69, 0x61, 0x67, 0x6f, 0x75, 0x79, 0x75, 0x6e, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x63, 0x61, 0x72, 0x65, 0x2d, 0x74, 0x6f, 0x6f,
	0x6c, 0x73, 0x2f, 0x64, 0x61, 0x74, 0x61, 0x6b, 0x69, 0x74, 0x2f, 0x70, 0x6c, 0x75, 0x67, 0x69,
	0x6e, 0x73, 0x2f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x2f, 0x74, 0x72, 0x61, 0x63, 0x65, 0x53,
	0x6b, 0x79, 0x77, 0x61, 0x6c, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x32, 0x2f, 0x72, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0xaa, 0x02, 0x1a, 0x53, 0x6b, 0x79, 0x57, 0x61, 0x6c, 0x6b, 0x69,
	0x6e, 0x67, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_register_InstancePing_proto_rawDescOnce sync.Once
	file_register_InstancePing_proto_rawDescData = file_register_InstancePing_proto_rawDesc
)

func file_register_InstancePing_proto_rawDescGZIP() []byte {
	file_register_InstancePing_proto_rawDescOnce.Do(func() {
		file_register_InstancePing_proto_rawDescData = protoimpl.X.CompressGZIP(file_register_InstancePing_proto_rawDescData)
	})
	return file_register_InstancePing_proto_rawDescData
}

var file_register_InstancePing_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_register_InstancePing_proto_goTypes = []interface{}{
	(*ServiceInstancePingPkg)(nil), // 0: ServiceInstancePingPkg
	(*common.Commands)(nil),        // 1: Commands
}
var file_register_InstancePing_proto_depIdxs = []int32{
	0, // 0: ServiceInstancePing.doPing:input_type -> ServiceInstancePingPkg
	1, // 1: ServiceInstancePing.doPing:output_type -> Commands
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_register_InstancePing_proto_init() }
func file_register_InstancePing_proto_init() {
	if File_register_InstancePing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_register_InstancePing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceInstancePingPkg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_register_InstancePing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_register_InstancePing_proto_goTypes,
		DependencyIndexes: file_register_InstancePing_proto_depIdxs,
		MessageInfos:      file_register_InstancePing_proto_msgTypes,
	}.Build()
	File_register_InstancePing_proto = out.File
	file_register_InstancePing_proto_rawDesc = nil
	file_register_InstancePing_proto_goTypes = nil
	file_register_InstancePing_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ServiceInstancePingClient is the client API for ServiceInstancePing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ServiceInstancePingClient interface {
	DoPing(ctx context.Context, in *ServiceInstancePingPkg, opts ...grpc.CallOption) (*common.Commands, error)
}

type serviceInstancePingClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceInstancePingClient(cc grpc.ClientConnInterface) ServiceInstancePingClient {
	return &serviceInstancePingClient{cc}
}

func (c *serviceInstancePingClient) DoPing(ctx context.Context, in *ServiceInstancePingPkg, opts ...grpc.CallOption) (*common.Commands, error) {
	out := new(common.Commands)
	err := c.cc.Invoke(ctx, "/ServiceInstancePing/doPing", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceInstancePingServer is the server API for ServiceInstancePing service.
type ServiceInstancePingServer interface {
	DoPing(context.Context, *ServiceInstancePingPkg) (*common.Commands, error)
}

// UnimplementedServiceInstancePingServer can be embedded to have forward compatible implementations.
type UnimplementedServiceInstancePingServer struct {
}

func (*UnimplementedServiceInstancePingServer) DoPing(context.Context, *ServiceInstancePingPkg) (*common.Commands, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoPing not implemented")
}

func RegisterServiceInstancePingServer(s *grpc.Server, srv ServiceInstancePingServer) {
	s.RegisterService(&_ServiceInstancePing_serviceDesc, srv)
}

func _ServiceInstancePing_DoPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceInstancePingPkg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceInstancePingServer).DoPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ServiceInstancePing/DoPing",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceInstancePingServer).DoPing(ctx, req.(*ServiceInstancePingPkg))
	}
	return interceptor(ctx, in, info, handler)
}

var _ServiceInstancePing_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ServiceInstancePing",
	HandlerType: (*ServiceInstancePingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "doPing",
			Handler:    _ServiceInstancePing_DoPing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "register/InstancePing.proto",
}
