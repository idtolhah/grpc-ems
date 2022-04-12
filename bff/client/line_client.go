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

type LineClient struct {
}

var (
	_                     = utils.LoadLocalEnv()
	lineGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	lineGrpcServiceClient masterpb.MasterServiceClient
)

func prepareLineGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, lineGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		lineGrpcServiceClient = nil
		return errors.New("connection to line gRPC service failed")
	}

	if lineGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// utils.CreateStartPromHttpServer(reg, 9092)

	lineGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *LineClient) GetLines(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	jsonData := cache.GetCacheByKeyDirect("lines")
	if jsonData != nil && utils.GetEnv("USE_CACHE") == "yes" {
		utils.Response(c, jsonData, nil)
		return
	}

	if err := prepareLineGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := lineGrpcServiceClient.GetLines(c, &masterpb.GetLinesRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var lines []masterpb.Line
	for _, u := range res.GetLines() {
		lines = append(lines, masterpb.Line{Id: u.Id, Name: u.Name})
	}
	utils.Response(c, &lines, nil)
}

func (uc *LineClient) GetLine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	id, _ := utils.GetParam(c, "id")

	if err := prepareLineGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := lineGrpcServiceClient.GetLine(c, &masterpb.GetLineRequest{Id: id})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, res.Line, nil)
}
