package clients

import (
	"context"
	"errors"
	"packing/pb/masterpb"
	"packing/pb/userpb"
	"packing/utils"

	"google.golang.org/grpc"
)

var (
	_                       = utils.LoadLocalEnv()
	masterGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	masterGrpcServiceClient masterpb.MasterServiceClient
	userGrpcService         = utils.GetEnv("USER_GRPC_SERVICE")
	userGrpcServiceClient   userpb.UserServiceClient
)

func prepareMasterGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, masterGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...,
	)

	if err != nil {
		masterGrpcServiceClient = nil
		return errors.New("connection to master gRPC service failed")
	}

	if masterGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	masterGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func prepareUserGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, userGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...,
	)

	if err != nil {
		userGrpcServiceClient = nil
		return errors.New("connection to user gRPC service failed")
	}

	if userGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	userGrpcServiceClient = userpb.NewUserServiceClient(conn)
	return nil
}
