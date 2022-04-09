package client

import (
	"context"
	"errors"

	"master/masterdb"
	"master/pb/redispb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type RedisClient struct {
}

var (
	_                      = masterdb.LoadLocalEnv()
	redisGrpcService       = masterdb.GetEnv("REDIS_GRPC_SERVICE")
	redisGrpcServiceClient redispb.RedisServiceClient
)

func prepareRedisGrpcClient(c *context.Context) error {

	conn, err := grpc.DialContext(*c, redisGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)

	if err != nil {
		redisGrpcServiceClient = nil
		return errors.New("connection to redis gRPC service failed")
	}

	if redisGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	redisGrpcServiceClient = redispb.NewRedisServiceClient(conn)
	return nil
}

func (uc *RedisClient) SetCache(c *context.Context, req *redispb.SetRequest) (*string, error) {

	if err := prepareRedisGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := redisGrpcServiceClient.SetCache(context.Background(), &redispb.SetRequest{Key: req.Key, Value: req.Value})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	data := res.Data
	return &data, nil
}
