// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package packingcmdpb

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

// PackingCmdServiceClient is the client API for PackingCmdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PackingCmdServiceClient interface {
	// FO
	CreatePacking(ctx context.Context, in *CreatePackingRequest, opts ...grpc.CallOption) (*CreatePackingResponse, error)
	CreateEquipmentChecking(ctx context.Context, in *CreateEquipmentCheckingRequest, opts ...grpc.CallOption) (*CreateEquipmentCheckingResponse, error)
	// AO
	UpdateEquipmentChecking(ctx context.Context, in *UpdateEquipmentCheckingRequest, opts ...grpc.CallOption) (*UpdateEquipmentCheckingResponse, error)
}

type packingCmdServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPackingCmdServiceClient(cc grpc.ClientConnInterface) PackingCmdServiceClient {
	return &packingCmdServiceClient{cc}
}

func (c *packingCmdServiceClient) CreatePacking(ctx context.Context, in *CreatePackingRequest, opts ...grpc.CallOption) (*CreatePackingResponse, error) {
	out := new(CreatePackingResponse)
	err := c.cc.Invoke(ctx, "/packingcmdpb.PackingCmdService/CreatePacking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingCmdServiceClient) CreateEquipmentChecking(ctx context.Context, in *CreateEquipmentCheckingRequest, opts ...grpc.CallOption) (*CreateEquipmentCheckingResponse, error) {
	out := new(CreateEquipmentCheckingResponse)
	err := c.cc.Invoke(ctx, "/packingcmdpb.PackingCmdService/CreateEquipmentChecking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingCmdServiceClient) UpdateEquipmentChecking(ctx context.Context, in *UpdateEquipmentCheckingRequest, opts ...grpc.CallOption) (*UpdateEquipmentCheckingResponse, error) {
	out := new(UpdateEquipmentCheckingResponse)
	err := c.cc.Invoke(ctx, "/packingcmdpb.PackingCmdService/UpdateEquipmentChecking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PackingCmdServiceServer is the server API for PackingCmdService service.
// All implementations must embed UnimplementedPackingCmdServiceServer
// for forward compatibility
type PackingCmdServiceServer interface {
	// FO
	CreatePacking(context.Context, *CreatePackingRequest) (*CreatePackingResponse, error)
	CreateEquipmentChecking(context.Context, *CreateEquipmentCheckingRequest) (*CreateEquipmentCheckingResponse, error)
	// AO
	UpdateEquipmentChecking(context.Context, *UpdateEquipmentCheckingRequest) (*UpdateEquipmentCheckingResponse, error)
	mustEmbedUnimplementedPackingCmdServiceServer()
}

// UnimplementedPackingCmdServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPackingCmdServiceServer struct {
}

func (UnimplementedPackingCmdServiceServer) CreatePacking(context.Context, *CreatePackingRequest) (*CreatePackingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePacking not implemented")
}
func (UnimplementedPackingCmdServiceServer) CreateEquipmentChecking(context.Context, *CreateEquipmentCheckingRequest) (*CreateEquipmentCheckingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEquipmentChecking not implemented")
}
func (UnimplementedPackingCmdServiceServer) UpdateEquipmentChecking(context.Context, *UpdateEquipmentCheckingRequest) (*UpdateEquipmentCheckingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEquipmentChecking not implemented")
}
func (UnimplementedPackingCmdServiceServer) mustEmbedUnimplementedPackingCmdServiceServer() {}

// UnsafePackingCmdServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PackingCmdServiceServer will
// result in compilation errors.
type UnsafePackingCmdServiceServer interface {
	mustEmbedUnimplementedPackingCmdServiceServer()
}

func RegisterPackingCmdServiceServer(s grpc.ServiceRegistrar, srv PackingCmdServiceServer) {
	s.RegisterService(&PackingCmdService_ServiceDesc, srv)
}

func _PackingCmdService_CreatePacking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePackingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingCmdServiceServer).CreatePacking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingcmdpb.PackingCmdService/CreatePacking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingCmdServiceServer).CreatePacking(ctx, req.(*CreatePackingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingCmdService_CreateEquipmentChecking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEquipmentCheckingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingCmdServiceServer).CreateEquipmentChecking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingcmdpb.PackingCmdService/CreateEquipmentChecking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingCmdServiceServer).CreateEquipmentChecking(ctx, req.(*CreateEquipmentCheckingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingCmdService_UpdateEquipmentChecking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEquipmentCheckingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingCmdServiceServer).UpdateEquipmentChecking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingcmdpb.PackingCmdService/UpdateEquipmentChecking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingCmdServiceServer).UpdateEquipmentChecking(ctx, req.(*UpdateEquipmentCheckingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PackingCmdService_ServiceDesc is the grpc.ServiceDesc for PackingCmdService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PackingCmdService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "packingcmdpb.PackingCmdService",
	HandlerType: (*PackingCmdServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePacking",
			Handler:    _PackingCmdService_CreatePacking_Handler,
		},
		{
			MethodName: "CreateEquipmentChecking",
			Handler:    _PackingCmdService_CreateEquipmentChecking_Handler,
		},
		{
			MethodName: "UpdateEquipmentChecking",
			Handler:    _PackingCmdService_UpdateEquipmentChecking_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "packingcmdpb/packingcmd.proto",
}