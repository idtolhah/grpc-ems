// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package filtrationquerypb

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

// FiltrationQueryServiceClient is the client API for FiltrationQueryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FiltrationQueryServiceClient interface {
	GetFiltrations(ctx context.Context, in *GetFiltrationsRequest, opts ...grpc.CallOption) (*GetFiltrationsResponse, error)
}

type filtrationQueryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFiltrationQueryServiceClient(cc grpc.ClientConnInterface) FiltrationQueryServiceClient {
	return &filtrationQueryServiceClient{cc}
}

func (c *filtrationQueryServiceClient) GetFiltrations(ctx context.Context, in *GetFiltrationsRequest, opts ...grpc.CallOption) (*GetFiltrationsResponse, error) {
	out := new(GetFiltrationsResponse)
	err := c.cc.Invoke(ctx, "/filtrationquerypb.FiltrationQueryService/GetFiltrations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FiltrationQueryServiceServer is the server API for FiltrationQueryService service.
// All implementations must embed UnimplementedFiltrationQueryServiceServer
// for forward compatibility
type FiltrationQueryServiceServer interface {
	GetFiltrations(context.Context, *GetFiltrationsRequest) (*GetFiltrationsResponse, error)
	mustEmbedUnimplementedFiltrationQueryServiceServer()
}

// UnimplementedFiltrationQueryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFiltrationQueryServiceServer struct {
}

func (UnimplementedFiltrationQueryServiceServer) GetFiltrations(context.Context, *GetFiltrationsRequest) (*GetFiltrationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFiltrations not implemented")
}
func (UnimplementedFiltrationQueryServiceServer) mustEmbedUnimplementedFiltrationQueryServiceServer() {
}

// UnsafeFiltrationQueryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FiltrationQueryServiceServer will
// result in compilation errors.
type UnsafeFiltrationQueryServiceServer interface {
	mustEmbedUnimplementedFiltrationQueryServiceServer()
}

func RegisterFiltrationQueryServiceServer(s grpc.ServiceRegistrar, srv FiltrationQueryServiceServer) {
	s.RegisterService(&FiltrationQueryService_ServiceDesc, srv)
}

func _FiltrationQueryService_GetFiltrations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFiltrationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FiltrationQueryServiceServer).GetFiltrations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filtrationquerypb.FiltrationQueryService/GetFiltrations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FiltrationQueryServiceServer).GetFiltrations(ctx, req.(*GetFiltrationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FiltrationQueryService_ServiceDesc is the grpc.ServiceDesc for FiltrationQueryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FiltrationQueryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "filtrationquerypb.FiltrationQueryService",
	HandlerType: (*FiltrationQueryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetFiltrations",
			Handler:    _FiltrationQueryService_GetFiltrations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pb/filtrationquerypb/filtrationquery.proto",
}
