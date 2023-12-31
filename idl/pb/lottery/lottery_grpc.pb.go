// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.6
// source: lottery.proto

package lottery

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
	LotteryService_InitAward_FullMethodName     = "/LotteryService/InitAward"
	LotteryService_Draw_FullMethodName          = "/LotteryService/Draw"
	LotteryService_ListAwardInfo_FullMethodName = "/LotteryService/ListAwardInfo"
	LotteryService_ToMysql_FullMethodName       = "/LotteryService/ToMysql"
)

// LotteryServiceClient is the client API for LotteryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LotteryServiceClient interface {
	InitAward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Draw(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*DrawResponse, error)
	ListAwardInfo(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*ListAwardInfoResponse, error)
	ToMysql(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ToMysqlResponse, error)
}

type lotteryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLotteryServiceClient(cc grpc.ClientConnInterface) LotteryServiceClient {
	return &lotteryServiceClient{cc}
}

func (c *lotteryServiceClient) InitAward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, LotteryService_InitAward_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotteryServiceClient) Draw(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*DrawResponse, error) {
	out := new(DrawResponse)
	err := c.cc.Invoke(ctx, LotteryService_Draw_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotteryServiceClient) ListAwardInfo(ctx context.Context, in *DrawRequest, opts ...grpc.CallOption) (*ListAwardInfoResponse, error) {
	out := new(ListAwardInfoResponse)
	err := c.cc.Invoke(ctx, LotteryService_ListAwardInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotteryServiceClient) ToMysql(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ToMysqlResponse, error) {
	out := new(ToMysqlResponse)
	err := c.cc.Invoke(ctx, LotteryService_ToMysql_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LotteryServiceServer is the server API for LotteryService service.
// All implementations must embed UnimplementedLotteryServiceServer
// for forward compatibility
type LotteryServiceServer interface {
	InitAward(context.Context, *Empty) (*Empty, error)
	Draw(context.Context, *DrawRequest) (*DrawResponse, error)
	ListAwardInfo(context.Context, *DrawRequest) (*ListAwardInfoResponse, error)
	ToMysql(context.Context, *Empty) (*ToMysqlResponse, error)
	mustEmbedUnimplementedLotteryServiceServer()
}

// UnimplementedLotteryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLotteryServiceServer struct {
}

func (UnimplementedLotteryServiceServer) InitAward(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitAward not implemented")
}
func (UnimplementedLotteryServiceServer) Draw(context.Context, *DrawRequest) (*DrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Draw not implemented")
}
func (UnimplementedLotteryServiceServer) ListAwardInfo(context.Context, *DrawRequest) (*ListAwardInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAwardInfo not implemented")
}
func (UnimplementedLotteryServiceServer) ToMysql(context.Context, *Empty) (*ToMysqlResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ToMysql not implemented")
}
func (UnimplementedLotteryServiceServer) mustEmbedUnimplementedLotteryServiceServer() {}

// UnsafeLotteryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LotteryServiceServer will
// result in compilation errors.
type UnsafeLotteryServiceServer interface {
	mustEmbedUnimplementedLotteryServiceServer()
}

func RegisterLotteryServiceServer(s grpc.ServiceRegistrar, srv LotteryServiceServer) {
	s.RegisterService(&LotteryService_ServiceDesc, srv)
}

func _LotteryService_InitAward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).InitAward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_InitAward_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).InitAward(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotteryService_Draw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).Draw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_Draw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).Draw(ctx, req.(*DrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotteryService_ListAwardInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).ListAwardInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_ListAwardInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).ListAwardInfo(ctx, req.(*DrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotteryService_ToMysql_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).ToMysql(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_ToMysql_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).ToMysql(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LotteryService_ServiceDesc is the grpc.ServiceDesc for LotteryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LotteryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LotteryService",
	HandlerType: (*LotteryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InitAward",
			Handler:    _LotteryService_InitAward_Handler,
		},
		{
			MethodName: "Draw",
			Handler:    _LotteryService_Draw_Handler,
		},
		{
			MethodName: "ListAwardInfo",
			Handler:    _LotteryService_ListAwardInfo_Handler,
		},
		{
			MethodName: "ToMysql",
			Handler:    _LotteryService_ToMysql_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lottery.proto",
}
