// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: blockchain.proto

package proto

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

// BlockchainServiceClient is the client API for BlockchainService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BlockchainServiceClient interface {
	Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error)
	CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*CreateWalletResponse, error)
	CreateBlockchain(ctx context.Context, in *CreateBlockchainRequest, opts ...grpc.CallOption) (*CreateBlockchainResponse, error)
	SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*ResponseTransaction, error)
	Mining(ctx context.Context, in *MiningRequest, opts ...grpc.CallOption) (*MiningResponse, error)
	SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*SendBlockResponse, error)
	FindMempool(ctx context.Context, in *FindMempoolRequest, opts ...grpc.CallOption) (*FindMempoolResponse, error)
}

type blockchainServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockchainServiceClient(cc grpc.ClientConnInterface) BlockchainServiceClient {
	return &blockchainServiceClient{cc}
}

func (c *blockchainServiceClient) Send(ctx context.Context, in *SendRequest, opts ...grpc.CallOption) (*SendResponse, error) {
	out := new(SendResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) CreateWallet(ctx context.Context, in *CreateWalletRequest, opts ...grpc.CallOption) (*CreateWalletResponse, error) {
	out := new(CreateWalletResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/CreateWallet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) CreateBlockchain(ctx context.Context, in *CreateBlockchainRequest, opts ...grpc.CallOption) (*CreateBlockchainResponse, error) {
	out := new(CreateBlockchainResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/CreateBlockchain", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) SendTransaction(ctx context.Context, in *SendTransactionRequest, opts ...grpc.CallOption) (*ResponseTransaction, error) {
	out := new(ResponseTransaction)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/SendTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) Mining(ctx context.Context, in *MiningRequest, opts ...grpc.CallOption) (*MiningResponse, error) {
	out := new(MiningResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/Mining", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) SendBlock(ctx context.Context, in *SendBlockRequest, opts ...grpc.CallOption) (*SendBlockResponse, error) {
	out := new(SendBlockResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/SendBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainServiceClient) FindMempool(ctx context.Context, in *FindMempoolRequest, opts ...grpc.CallOption) (*FindMempoolResponse, error) {
	out := new(FindMempoolResponse)
	err := c.cc.Invoke(ctx, "/blockchain.BlockchainService/FindMempool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockchainServiceServer is the server API for BlockchainService service.
// All implementations must embed UnimplementedBlockchainServiceServer
// for forward compatibility
type BlockchainServiceServer interface {
	Send(context.Context, *SendRequest) (*SendResponse, error)
	CreateWallet(context.Context, *CreateWalletRequest) (*CreateWalletResponse, error)
	CreateBlockchain(context.Context, *CreateBlockchainRequest) (*CreateBlockchainResponse, error)
	SendTransaction(context.Context, *SendTransactionRequest) (*ResponseTransaction, error)
	Mining(context.Context, *MiningRequest) (*MiningResponse, error)
	SendBlock(context.Context, *SendBlockRequest) (*SendBlockResponse, error)
	FindMempool(context.Context, *FindMempoolRequest) (*FindMempoolResponse, error)
	//mustEmbedUnimplementedBlockchainServiceServer()
}

// UnimplementedBlockchainServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBlockchainServiceServer struct {
}

func (UnimplementedBlockchainServiceServer) Send(context.Context, *SendRequest) (*SendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedBlockchainServiceServer) CreateWallet(context.Context, *CreateWalletRequest) (*CreateWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWallet not implemented")
}
func (UnimplementedBlockchainServiceServer) CreateBlockchain(context.Context, *CreateBlockchainRequest) (*CreateBlockchainResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBlockchain not implemented")
}
func (UnimplementedBlockchainServiceServer) SendTransaction(context.Context, *SendTransactionRequest) (*ResponseTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendTransaction not implemented")
}
func (UnimplementedBlockchainServiceServer) Mining(context.Context, *MiningRequest) (*MiningResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Mining not implemented")
}
func (UnimplementedBlockchainServiceServer) SendBlock(context.Context, *SendBlockRequest) (*SendBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendBlock not implemented")
}
func (UnimplementedBlockchainServiceServer) FindMempool(context.Context, *FindMempoolRequest) (*FindMempoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindMempool not implemented")
}
func (UnimplementedBlockchainServiceServer) mustEmbedUnimplementedBlockchainServiceServer() {}

// UnsafeBlockchainServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockchainServiceServer will
// result in compilation errors.
type UnsafeBlockchainServiceServer interface {
	mustEmbedUnimplementedBlockchainServiceServer()
}

func RegisterBlockchainServiceServer(s grpc.ServiceRegistrar, srv BlockchainServiceServer) {
	s.RegisterService(&BlockchainService_ServiceDesc, srv)
}

func _BlockchainService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).Send(ctx, req.(*SendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_CreateWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).CreateWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/CreateWallet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).CreateWallet(ctx, req.(*CreateWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_CreateBlockchain_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBlockchainRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).CreateBlockchain(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/CreateBlockchain",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).CreateBlockchain(ctx, req.(*CreateBlockchainRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_SendTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).SendTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/SendTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).SendTransaction(ctx, req.(*SendTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_Mining_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MiningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).Mining(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/Mining",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).Mining(ctx, req.(*MiningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_SendBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).SendBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/SendBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).SendBlock(ctx, req.(*SendBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BlockchainService_FindMempool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindMempoolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServiceServer).FindMempool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blockchain.BlockchainService/FindMempool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServiceServer).FindMempool(ctx, req.(*FindMempoolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BlockchainService_ServiceDesc is the grpc.ServiceDesc for BlockchainService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BlockchainService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blockchain.BlockchainService",
	HandlerType: (*BlockchainServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _BlockchainService_Send_Handler,
		},
		{
			MethodName: "CreateWallet",
			Handler:    _BlockchainService_CreateWallet_Handler,
		},
		{
			MethodName: "CreateBlockchain",
			Handler:    _BlockchainService_CreateBlockchain_Handler,
		},
		{
			MethodName: "SendTransaction",
			Handler:    _BlockchainService_SendTransaction_Handler,
		},
		{
			MethodName: "Mining",
			Handler:    _BlockchainService_Mining_Handler,
		},
		{
			MethodName: "SendBlock",
			Handler:    _BlockchainService_SendBlock_Handler,
		},
		{
			MethodName: "FindMempool",
			Handler:    _BlockchainService_FindMempool_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blockchain.proto",
}
