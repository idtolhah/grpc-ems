package client

import (
	"context"
	"errors"
	"log"
	"time"

	"bff/cache"
	"bff/pb/masterpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
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
	timeout               = 10 * time.Second
	_                     = utils.LoadLocalEnv()
	areaGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	areaGrpcServiceClient masterpb.MasterServiceClient
)

func prepareAreaGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()
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
	utils.CreateStartPromHttpServer(reg, 9092)

	areaGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *AreaClient) GetAreas(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if utils.GetEnv("USE_CACHE") == "yes" {
		jsonData := cache.GetCacheByKeyDirect("areas")
		if jsonData != nil {
			go log.Println("Use Cache")
			utils.Response(c, jsonData, nil)
			return
		}
	}

	if err := prepareAreaGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := areaGrpcServiceClient.GetAreas(c, &masterpb.GetAreasRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var areas []Area
	for _, u := range res.GetAreas() {
		areas = append(areas, Area{Id: u.Id, Name: u.Name})
	}
	utils.Response(c, &areas, nil)
}

func (uc *AreaClient) GetArea(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	id, _ := utils.GetParam(c, "id")

	if err := prepareAreaGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := areaGrpcServiceClient.GetArea(c, &masterpb.GetAreaRequest{Id: id})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, res.Area, nil)
}
