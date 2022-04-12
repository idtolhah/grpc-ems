package client

import (
	"context"
	"errors"

	"bff/cache"
	"bff/pb/masterpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type UnitClient struct {
}

var (
	_                     = utils.LoadLocalEnv()
	unitGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	unitGrpcServiceClient masterpb.MasterServiceClient
)

func prepareUnitGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, unitGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		unitGrpcServiceClient = nil
		return errors.New("connection to unit gRPC service failed")
	}

	if unitGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// utils.CreateStartPromHttpServer(reg, 9092)

	unitGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *UnitClient) GetUnits(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	jsonData := cache.GetCacheByKeyDirect("units")
	if jsonData != nil && utils.GetEnv("USE_CACHE") == "yes" {
		utils.Response(c, jsonData, nil)
		return
	}

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := unitGrpcServiceClient.GetUnits(c, &masterpb.GetUnitsRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var units []masterpb.Unit
	for _, u := range res.GetUnits() {
		units = append(units, masterpb.Unit{Id: u.Id, Name: u.Name})
	}
	utils.Response(c, &units, nil)
}

func (uc *UnitClient) GetUnit(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	id, _ := utils.GetParam(c, "id")

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := unitGrpcServiceClient.GetUnit(c, &masterpb.GetUnitRequest{Id: id})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, res.Unit, nil)
}
