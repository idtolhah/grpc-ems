// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package packingpb

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

// PackingServiceClient is the client API for PackingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PackingServiceClient interface {
	// All
	GetPackings(ctx context.Context, in *GetPackingsRequest, opts ...grpc.CallOption) (*GetPackingsResponse, error)
	GetPacking(ctx context.Context, in *GetPackingRequest, opts ...grpc.CallOption) (*GetPackingResponse, error)
	// FO
	CreatePacking(ctx context.Context, in *CreatePackingRequest, opts ...grpc.CallOption) (*CreatePackingResponse, error)
	CreateEquipmentChecking(ctx context.Context, in *CreateEquipmentCheckingRequest, opts ...grpc.CallOption) (*CreateEquipmentCheckingResponse, error)
	// AO
	UpdateEquipmentChecking(ctx context.Context, in *UpdateEquipmentCheckingRequest, opts ...grpc.CallOption) (*UpdateEquipmentCheckingResponse, error)
	// QFS
	CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
}

type packingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPackingServiceClient(cc grpc.ClientConnInterface) PackingServiceClient {
	return &packingServiceClient{cc}
}

func (c *packingServiceClient) GetPackings(ctx context.Context, in *GetPackingsRequest, opts ...grpc.CallOption) (*GetPackingsResponse, error) {
	out := new(GetPackingsResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/GetPackings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingServiceClient) GetPacking(ctx context.Context, in *GetPackingRequest, opts ...grpc.CallOption) (*GetPackingResponse, error) {
	out := new(GetPackingResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/GetPacking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingServiceClient) CreatePacking(ctx context.Context, in *CreatePackingRequest, opts ...grpc.CallOption) (*CreatePackingResponse, error) {
	out := new(CreatePackingResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/CreatePacking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingServiceClient) CreateEquipmentChecking(ctx context.Context, in *CreateEquipmentCheckingRequest, opts ...grpc.CallOption) (*CreateEquipmentCheckingResponse, error) {
	out := new(CreateEquipmentCheckingResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/CreateEquipmentChecking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingServiceClient) UpdateEquipmentChecking(ctx context.Context, in *UpdateEquipmentCheckingRequest, opts ...grpc.CallOption) (*UpdateEquipmentCheckingResponse, error) {
	out := new(UpdateEquipmentCheckingResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/UpdateEquipmentChecking", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *packingServiceClient) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, "/packingpb.PackingService/CreateComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PackingServiceServer is the server API for PackingService service.
// All implementations must embed UnimplementedPackingServiceServer
// for forward compatibility
type PackingServiceServer interface {
	// All
	GetPackings(context.Context, *GetPackingsRequest) (*GetPackingsResponse, error)
	GetPacking(context.Context, *GetPackingRequest) (*GetPackingResponse, error)
	// FO
	CreatePacking(context.Context, *CreatePackingRequest) (*CreatePackingResponse, error)
	CreateEquipmentChecking(context.Context, *CreateEquipmentCheckingRequest) (*CreateEquipmentCheckingResponse, error)
	// AO
	UpdateEquipmentChecking(context.Context, *UpdateEquipmentCheckingRequest) (*UpdateEquipmentCheckingResponse, error)
	// QFS
	CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
	mustEmbedUnimplementedPackingServiceServer()
}

// UnimplementedPackingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPackingServiceServer struct {
}

func (UnimplementedPackingServiceServer) GetPackings(context.Context, *GetPackingsRequest) (*GetPackingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackings not implemented")
}
func (UnimplementedPackingServiceServer) GetPacking(context.Context, *GetPackingRequest) (*GetPackingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPacking not implemented")
}
func (UnimplementedPackingServiceServer) CreatePacking(context.Context, *CreatePackingRequest) (*CreatePackingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePacking not implemented")
}
func (UnimplementedPackingServiceServer) CreateEquipmentChecking(context.Context, *CreateEquipmentCheckingRequest) (*CreateEquipmentCheckingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateEquipmentChecking not implemented")
}
func (UnimplementedPackingServiceServer) UpdateEquipmentChecking(context.Context, *UpdateEquipmentCheckingRequest) (*UpdateEquipmentCheckingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEquipmentChecking not implemented")
}
func (UnimplementedPackingServiceServer) CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedPackingServiceServer) mustEmbedUnimplementedPackingServiceServer() {}

// UnsafePackingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PackingServiceServer will
// result in compilation errors.
type UnsafePackingServiceServer interface {
	mustEmbedUnimplementedPackingServiceServer()
}

func RegisterPackingServiceServer(s grpc.ServiceRegistrar, srv PackingServiceServer) {
	s.RegisterService(&PackingService_ServiceDesc, srv)
}

func _PackingService_GetPackings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPackingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).GetPackings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/GetPackings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).GetPackings(ctx, req.(*GetPackingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingService_GetPacking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPackingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).GetPacking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/GetPacking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).GetPacking(ctx, req.(*GetPackingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingService_CreatePacking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePackingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).CreatePacking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/CreatePacking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).CreatePacking(ctx, req.(*CreatePackingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingService_CreateEquipmentChecking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateEquipmentCheckingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).CreateEquipmentChecking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/CreateEquipmentChecking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).CreateEquipmentChecking(ctx, req.(*CreateEquipmentCheckingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingService_UpdateEquipmentChecking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEquipmentCheckingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).UpdateEquipmentChecking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/UpdateEquipmentChecking",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).UpdateEquipmentChecking(ctx, req.(*UpdateEquipmentCheckingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PackingService_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PackingServiceServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/packingpb.PackingService/CreateComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PackingServiceServer).CreateComment(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PackingService_ServiceDesc is the grpc.ServiceDesc for PackingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PackingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "packingpb.PackingService",
	HandlerType: (*PackingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPackings",
			Handler:    _PackingService_GetPackings_Handler,
		},
		{
			MethodName: "GetPacking",
			Handler:    _PackingService_GetPacking_Handler,
		},
		{
			MethodName: "CreatePacking",
			Handler:    _PackingService_CreatePacking_Handler,
		},
		{
			MethodName: "CreateEquipmentChecking",
			Handler:    _PackingService_CreateEquipmentChecking_Handler,
		},
		{
			MethodName: "UpdateEquipmentChecking",
			Handler:    _PackingService_UpdateEquipmentChecking_Handler,
		},
		{
			MethodName: "CreateComment",
			Handler:    _PackingService_CreateComment_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "packingpb/packing.proto",
}