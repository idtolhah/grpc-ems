package client

import (
	"context"
	"errors"

	"bff/pb/masterpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Area struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type AreaClient struct {
}

var (
	_                     = loadLocalEnv()
	areaGrpcService       = GetEnv("MASTER_GRPC_SERVICE")
	areaGrpcServiceClient masterpb.MasterServiceClient
)

func prepareAreaGrpcClient(c *context.Context) error {

	// Prom: Get Registry & Metrics
	reg, grpcMetrics := GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, areaGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		areaGrpcServiceClient = nil
		return errors.New("connection to area gRPC service failed")
	}

	if areaGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	CreateStartPromHttpServer(reg, 9092)

	areaGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *AreaClient) GetAreas(c *context.Context) (*[]Area, error) {

	if err := prepareAreaGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := areaGrpcServiceClient.GetAreas(*c, &masterpb.GetAreasRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	var areas []Area
	for _, u := range res.GetAreas() {
		areas = append(areas, Area{Id: u.Id, Name: u.Name})
	}
	return &areas, nil
}
