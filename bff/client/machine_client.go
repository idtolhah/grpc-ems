package client

import (
	"context"
	"errors"
	"log"

	"bff/cache"
	"bff/pb/masterpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type MachineClient struct {
}

var (
	_                        = utils.LoadLocalEnv()
	machineGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	machineGrpcServiceClient masterpb.MasterServiceClient
)

func prepareMachineGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, machineGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		machineGrpcServiceClient = nil
		return errors.New("connection to machine gRPC service failed")
	}

	if machineGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// utils.CreateStartPromHttpServer(reg, 9092)

	machineGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *MachineClient) GetMachines(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if utils.GetEnv("USE_CACHE") == "yes" {
		jsonData := cache.GetCacheByKeyDirect("machines")
		if jsonData != nil {
			go log.Println("Use Cache")
			utils.Response(c, jsonData, nil)
			return
		}
	}

	if err := prepareMachineGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := machineGrpcServiceClient.GetMachines(c, &masterpb.GetMachinesRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var machines []masterpb.Machine
	for _, u := range res.GetMachines() {
		machines = append(machines, masterpb.Machine{Id: u.Id, Name: u.Name})
	}
	utils.Response(c, &machines, nil)
}

func (uc *MachineClient) GetMachine(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	id, _ := utils.GetParam(c, "id")

	if err := prepareMachineGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := machineGrpcServiceClient.GetMachine(c, &masterpb.GetMachineRequest{Id: id})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, res.Machine, nil)
}
