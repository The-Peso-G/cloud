// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/service.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() { proto.RegisterFile("pb/service.proto", fileDescriptor_6ff5ab49d8a5fcc4) }

var fileDescriptor_6ff5ab49d8a5fcc4 = []byte{
	// 370 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xbb, 0x4e, 0x02, 0x41,
	0x14, 0x86, 0xd9, 0x46, 0x93, 0xb1, 0x21, 0x93, 0x78, 0xc9, 0x36, 0x1a, 0x0a, 0x35, 0x0a, 0xb3,
	0x78, 0xc1, 0x98, 0x60, 0xa3, 0x28, 0xdb, 0x58, 0x61, 0xb4, 0xb0, 0x31, 0x3b, 0xc3, 0x61, 0x9d,
	0x08, 0x3b, 0xc3, 0x5c, 0x30, 0x76, 0x74, 0xfa, 0x0a, 0xbe, 0xad, 0x61, 0x97, 0x55, 0x11, 0x07,
	0xc1, 0x7a, 0xbe, 0xef, 0x3f, 0xff, 0x9c, 0xe4, 0xa0, 0xa2, 0xa4, 0x81, 0x06, 0x35, 0xe0, 0x0c,
	0x88, 0x54, 0xc2, 0x08, 0xbc, 0x21, 0x58, 0x87, 0xb0, 0xae, 0xb0, 0x6d, 0x12, 0x2b, 0xc9, 0xe2,
	0xc8, 0xc0, 0x73, 0xf4, 0x42, 0x24, 0xf5, 0xcf, 0x63, 0x6e, 0x1e, 0x2d, 0x25, 0x4c, 0xf4, 0x82,
	0x58, 0x54, 0x04, 0xeb, 0x04, 0x29, 0x17, 0x8c, 0xb8, 0xca, 0x18, 0x0c, 0x24, 0x0d, 0x58, 0x97,
	0x43, 0x62, 0x1a, 0x22, 0xe9, 0xf0, 0xd8, 0xaa, 0xc8, 0x70, 0x91, 0x64, 0xe1, 0x7e, 0x6d, 0xde,
	0x88, 0x36, 0x8c, 0x2a, 0xe9, 0x4c, 0x3b, 0x1c, 0x2e, 0xa3, 0x95, 0x50, 0x49, 0x16, 0x66, 0x00,
	0x7e, 0x40, 0x28, 0x04, 0x73, 0x99, 0x31, 0x78, 0x9f, 0xb8, 0x2a, 0x93, 0x2f, 0xaa, 0x05, 0x7d,
	0x0b, 0xda, 0xf8, 0x5b, 0x6e, 0x38, 0x23, 0x4b, 0x85, 0xaa, 0x87, 0xfb, 0xa8, 0x18, 0x82, 0x69,
	0x81, 0x16, 0x56, 0x31, 0xb8, 0xe6, 0xc9, 0x93, 0xc6, 0x07, 0x33, 0xc7, 0x4c, 0xb0, 0xf9, 0xb0,
	0x6d, 0xb7, 0xf2, 0x9d, 0x4f, 0x47, 0xbe, 0x7b, 0xc8, 0x6f, 0x81, 0x51, 0x1c, 0x06, 0x90, 0x3f,
	0x36, 0x95, 0xe8, 0x65, 0xad, 0x70, 0x7d, 0x56, 0x94, 0xcb, 0xca, 0x7b, 0x9c, 0xfd, 0x4f, 0xd6,
	0x52, 0x24, 0x1a, 0x4a, 0x05, 0x3c, 0xf4, 0xd0, 0xfa, 0x4f, 0x50, 0xdf, 0x45, 0x5d, 0x0b, 0x1a,
	0x9f, 0xce, 0x9f, 0x3d, 0x56, 0xf2, 0x56, 0x3b, 0x7f, 0x6f, 0x27, 0x15, 0xd2, 0xf5, 0xbc, 0x79,
	0x68, 0xf5, 0x56, 0xb6, 0x23, 0x33, 0x55, 0xa0, 0xe6, 0x8e, 0x99, 0x14, 0x26, 0xa7, 0x9f, 0x2c,
	0xaa, 0x7d, 0x6e, 0x83, 0x23, 0x7c, 0x63, 0xa9, 0x66, 0x8a, 0x53, 0x68, 0x0a, 0x75, 0x35, 0x80,
	0xc4, 0x68, 0x5c, 0x76, 0xe7, 0x4d, 0xd3, 0xfe, 0xa6, 0x9b, 0x4e, 0x89, 0x52, 0x61, 0xd7, 0xab,
	0x7a, 0xf8, 0xd5, 0x43, 0x6b, 0x21, 0x98, 0xc6, 0xf4, 0x41, 0xe1, 0x63, 0x77, 0xc2, 0x2f, 0x78,
	0xfe, 0xeb, 0xda, 0x82, 0x56, 0xfe, 0xe9, 0x8b, 0xf2, 0xfd, 0xde, 0x9c, 0xb7, 0x5b, 0x97, 0x94,
	0x2e, 0xa5, 0x77, 0x7b, 0xf4, 0x11, 0x00, 0x00, 0xff, 0xff, 0x00, 0xf6, 0xe7, 0x5e, 0x5f, 0x04,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GrpcGatewayClient is the client API for GrpcGateway service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GrpcGatewayClient interface {
	// Get all devices
	GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (GrpcGateway_GetDevicesClient, error)
	// Get resource links of devices.
	GetResourceLinks(ctx context.Context, in *GetResourceLinksRequest, opts ...grpc.CallOption) (GrpcGateway_GetResourceLinksClient, error)
	RetrieveResourceFromDevice(ctx context.Context, in *RetrieveResourceFromDeviceRequest, opts ...grpc.CallOption) (*RetrieveResourceFromDeviceResponse, error)
	// Retrieve resources values from resource shadow
	RetrieveResourcesValues(ctx context.Context, in *RetrieveResourcesValuesRequest, opts ...grpc.CallOption) (GrpcGateway_RetrieveResourcesValuesClient, error)
	// Update resource values
	UpdateResourcesValues(ctx context.Context, in *UpdateResourceValuesRequest, opts ...grpc.CallOption) (*UpdateResourceValuesResponse, error)
	// Subscribe to events
	SubscribeForEvents(ctx context.Context, opts ...grpc.CallOption) (GrpcGateway_SubscribeForEventsClient, error)
	// Get client configuration
	GetClientConfiguration(ctx context.Context, in *ClientConfigurationRequest, opts ...grpc.CallOption) (*ClientConfigurationResponse, error)
}

type grpcGatewayClient struct {
	cc *grpc.ClientConn
}

func NewGrpcGatewayClient(cc *grpc.ClientConn) GrpcGatewayClient {
	return &grpcGatewayClient{cc}
}

func (c *grpcGatewayClient) GetDevices(ctx context.Context, in *GetDevicesRequest, opts ...grpc.CallOption) (GrpcGateway_GetDevicesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GrpcGateway_serviceDesc.Streams[0], "/ocf.cloud.grpcgateway.pb.GrpcGateway/GetDevices", opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcGatewayGetDevicesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GrpcGateway_GetDevicesClient interface {
	Recv() (*Device, error)
	grpc.ClientStream
}

type grpcGatewayGetDevicesClient struct {
	grpc.ClientStream
}

func (x *grpcGatewayGetDevicesClient) Recv() (*Device, error) {
	m := new(Device)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcGatewayClient) GetResourceLinks(ctx context.Context, in *GetResourceLinksRequest, opts ...grpc.CallOption) (GrpcGateway_GetResourceLinksClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GrpcGateway_serviceDesc.Streams[1], "/ocf.cloud.grpcgateway.pb.GrpcGateway/GetResourceLinks", opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcGatewayGetResourceLinksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GrpcGateway_GetResourceLinksClient interface {
	Recv() (*ResourceLink, error)
	grpc.ClientStream
}

type grpcGatewayGetResourceLinksClient struct {
	grpc.ClientStream
}

func (x *grpcGatewayGetResourceLinksClient) Recv() (*ResourceLink, error) {
	m := new(ResourceLink)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcGatewayClient) RetrieveResourceFromDevice(ctx context.Context, in *RetrieveResourceFromDeviceRequest, opts ...grpc.CallOption) (*RetrieveResourceFromDeviceResponse, error) {
	out := new(RetrieveResourceFromDeviceResponse)
	err := c.cc.Invoke(ctx, "/ocf.cloud.grpcgateway.pb.GrpcGateway/RetrieveResourceFromDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcGatewayClient) RetrieveResourcesValues(ctx context.Context, in *RetrieveResourcesValuesRequest, opts ...grpc.CallOption) (GrpcGateway_RetrieveResourcesValuesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GrpcGateway_serviceDesc.Streams[2], "/ocf.cloud.grpcgateway.pb.GrpcGateway/RetrieveResourcesValues", opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcGatewayRetrieveResourcesValuesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GrpcGateway_RetrieveResourcesValuesClient interface {
	Recv() (*ResourceValue, error)
	grpc.ClientStream
}

type grpcGatewayRetrieveResourcesValuesClient struct {
	grpc.ClientStream
}

func (x *grpcGatewayRetrieveResourcesValuesClient) Recv() (*ResourceValue, error) {
	m := new(ResourceValue)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcGatewayClient) UpdateResourcesValues(ctx context.Context, in *UpdateResourceValuesRequest, opts ...grpc.CallOption) (*UpdateResourceValuesResponse, error) {
	out := new(UpdateResourceValuesResponse)
	err := c.cc.Invoke(ctx, "/ocf.cloud.grpcgateway.pb.GrpcGateway/UpdateResourcesValues", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcGatewayClient) SubscribeForEvents(ctx context.Context, opts ...grpc.CallOption) (GrpcGateway_SubscribeForEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GrpcGateway_serviceDesc.Streams[3], "/ocf.cloud.grpcgateway.pb.GrpcGateway/SubscribeForEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcGatewaySubscribeForEventsClient{stream}
	return x, nil
}

type GrpcGateway_SubscribeForEventsClient interface {
	Send(*SubscribeForEvents) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type grpcGatewaySubscribeForEventsClient struct {
	grpc.ClientStream
}

func (x *grpcGatewaySubscribeForEventsClient) Send(m *SubscribeForEvents) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcGatewaySubscribeForEventsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcGatewayClient) GetClientConfiguration(ctx context.Context, in *ClientConfigurationRequest, opts ...grpc.CallOption) (*ClientConfigurationResponse, error) {
	out := new(ClientConfigurationResponse)
	err := c.cc.Invoke(ctx, "/ocf.cloud.grpcgateway.pb.GrpcGateway/GetClientConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcGatewayServer is the server API for GrpcGateway service.
type GrpcGatewayServer interface {
	// Get all devices
	GetDevices(*GetDevicesRequest, GrpcGateway_GetDevicesServer) error
	// Get resource links of devices.
	GetResourceLinks(*GetResourceLinksRequest, GrpcGateway_GetResourceLinksServer) error
	RetrieveResourceFromDevice(context.Context, *RetrieveResourceFromDeviceRequest) (*RetrieveResourceFromDeviceResponse, error)
	// Retrieve resources values from resource shadow
	RetrieveResourcesValues(*RetrieveResourcesValuesRequest, GrpcGateway_RetrieveResourcesValuesServer) error
	// Update resource values
	UpdateResourcesValues(context.Context, *UpdateResourceValuesRequest) (*UpdateResourceValuesResponse, error)
	// Subscribe to events
	SubscribeForEvents(GrpcGateway_SubscribeForEventsServer) error
	// Get client configuration
	GetClientConfiguration(context.Context, *ClientConfigurationRequest) (*ClientConfigurationResponse, error)
}

// UnimplementedGrpcGatewayServer can be embedded to have forward compatible implementations.
type UnimplementedGrpcGatewayServer struct {
}

func (*UnimplementedGrpcGatewayServer) GetDevices(req *GetDevicesRequest, srv GrpcGateway_GetDevicesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetDevices not implemented")
}
func (*UnimplementedGrpcGatewayServer) GetResourceLinks(req *GetResourceLinksRequest, srv GrpcGateway_GetResourceLinksServer) error {
	return status.Errorf(codes.Unimplemented, "method GetResourceLinks not implemented")
}
func (*UnimplementedGrpcGatewayServer) RetrieveResourceFromDevice(ctx context.Context, req *RetrieveResourceFromDeviceRequest) (*RetrieveResourceFromDeviceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RetrieveResourceFromDevice not implemented")
}
func (*UnimplementedGrpcGatewayServer) RetrieveResourcesValues(req *RetrieveResourcesValuesRequest, srv GrpcGateway_RetrieveResourcesValuesServer) error {
	return status.Errorf(codes.Unimplemented, "method RetrieveResourcesValues not implemented")
}
func (*UnimplementedGrpcGatewayServer) UpdateResourcesValues(ctx context.Context, req *UpdateResourceValuesRequest) (*UpdateResourceValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateResourcesValues not implemented")
}
func (*UnimplementedGrpcGatewayServer) SubscribeForEvents(srv GrpcGateway_SubscribeForEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeForEvents not implemented")
}
func (*UnimplementedGrpcGatewayServer) GetClientConfiguration(ctx context.Context, req *ClientConfigurationRequest) (*ClientConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClientConfiguration not implemented")
}

func RegisterGrpcGatewayServer(s *grpc.Server, srv GrpcGatewayServer) {
	s.RegisterService(&_GrpcGateway_serviceDesc, srv)
}

func _GrpcGateway_GetDevices_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDevicesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GrpcGatewayServer).GetDevices(m, &grpcGatewayGetDevicesServer{stream})
}

type GrpcGateway_GetDevicesServer interface {
	Send(*Device) error
	grpc.ServerStream
}

type grpcGatewayGetDevicesServer struct {
	grpc.ServerStream
}

func (x *grpcGatewayGetDevicesServer) Send(m *Device) error {
	return x.ServerStream.SendMsg(m)
}

func _GrpcGateway_GetResourceLinks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetResourceLinksRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GrpcGatewayServer).GetResourceLinks(m, &grpcGatewayGetResourceLinksServer{stream})
}

type GrpcGateway_GetResourceLinksServer interface {
	Send(*ResourceLink) error
	grpc.ServerStream
}

type grpcGatewayGetResourceLinksServer struct {
	grpc.ServerStream
}

func (x *grpcGatewayGetResourceLinksServer) Send(m *ResourceLink) error {
	return x.ServerStream.SendMsg(m)
}

func _GrpcGateway_RetrieveResourceFromDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetrieveResourceFromDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcGatewayServer).RetrieveResourceFromDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocf.cloud.grpcgateway.pb.GrpcGateway/RetrieveResourceFromDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcGatewayServer).RetrieveResourceFromDevice(ctx, req.(*RetrieveResourceFromDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcGateway_RetrieveResourcesValues_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RetrieveResourcesValuesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GrpcGatewayServer).RetrieveResourcesValues(m, &grpcGatewayRetrieveResourcesValuesServer{stream})
}

type GrpcGateway_RetrieveResourcesValuesServer interface {
	Send(*ResourceValue) error
	grpc.ServerStream
}

type grpcGatewayRetrieveResourcesValuesServer struct {
	grpc.ServerStream
}

func (x *grpcGatewayRetrieveResourcesValuesServer) Send(m *ResourceValue) error {
	return x.ServerStream.SendMsg(m)
}

func _GrpcGateway_UpdateResourcesValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateResourceValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcGatewayServer).UpdateResourcesValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocf.cloud.grpcgateway.pb.GrpcGateway/UpdateResourcesValues",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcGatewayServer).UpdateResourcesValues(ctx, req.(*UpdateResourceValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcGateway_SubscribeForEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcGatewayServer).SubscribeForEvents(&grpcGatewaySubscribeForEventsServer{stream})
}

type GrpcGateway_SubscribeForEventsServer interface {
	Send(*Event) error
	Recv() (*SubscribeForEvents, error)
	grpc.ServerStream
}

type grpcGatewaySubscribeForEventsServer struct {
	grpc.ServerStream
}

func (x *grpcGatewaySubscribeForEventsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcGatewaySubscribeForEventsServer) Recv() (*SubscribeForEvents, error) {
	m := new(SubscribeForEvents)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GrpcGateway_GetClientConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcGatewayServer).GetClientConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ocf.cloud.grpcgateway.pb.GrpcGateway/GetClientConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcGatewayServer).GetClientConfiguration(ctx, req.(*ClientConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GrpcGateway_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ocf.cloud.grpcgateway.pb.GrpcGateway",
	HandlerType: (*GrpcGatewayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RetrieveResourceFromDevice",
			Handler:    _GrpcGateway_RetrieveResourceFromDevice_Handler,
		},
		{
			MethodName: "UpdateResourcesValues",
			Handler:    _GrpcGateway_UpdateResourcesValues_Handler,
		},
		{
			MethodName: "GetClientConfiguration",
			Handler:    _GrpcGateway_GetClientConfiguration_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDevices",
			Handler:       _GrpcGateway_GetDevices_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetResourceLinks",
			Handler:       _GrpcGateway_GetResourceLinks_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RetrieveResourcesValues",
			Handler:       _GrpcGateway_RetrieveResourcesValues_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeForEvents",
			Handler:       _GrpcGateway_SubscribeForEvents_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/service.proto",
}
