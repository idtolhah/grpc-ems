package clients

import (
	"context"
	"errors"
	"packing/pb/masterpb"
	"packing/utils"

	"google.golang.org/grpc"
)

var (
	_                       = utils.LoadLocalEnv()
	masterGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	masterGrpcServiceClient masterpb.MasterServiceClient
)

func prepareUnitGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, masterGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...,
	)

	if err != nil {
		masterGrpcServiceClient = nil
		return errors.New("connection to unit gRPC service failed")
	}

	if masterGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	masterGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}
