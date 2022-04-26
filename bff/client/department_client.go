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

type DepartmentClient struct {
}

var (
	_                           = utils.LoadLocalEnv()
	departmentGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	departmentGrpcServiceClient masterpb.MasterServiceClient
)

func prepareDepartmentGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, departmentGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		departmentGrpcServiceClient = nil
		return errors.New("connection to department gRPC service failed")
	}

	if departmentGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// utils.CreateStartPromHttpServer(reg, 9092)

	departmentGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *DepartmentClient) GetDepartments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if utils.GetEnv("USE_CACHE") == "yes" {
		jsonData := cache.GetCacheByKeyDirect("departments")
		if jsonData != nil {
			go log.Println("Use Cache")
			utils.Response(c, jsonData, nil)
			return
		}
	}

	if err := prepareDepartmentGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := departmentGrpcServiceClient.GetDepartments(c, &masterpb.GetDepartmentsRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var departments []masterpb.Department
	for _, u := range res.GetDepartments() {
		departments = append(departments, masterpb.Department{Id: u.Id, Name: u.Name})
	}
	utils.Response(c, &departments, nil)
}

func (uc *DepartmentClient) GetDepartment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	id, _ := utils.GetParam(c, "id")

	if err := prepareDepartmentGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := departmentGrpcServiceClient.GetDepartment(c, &masterpb.GetDepartmentRequest{Id: id})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, res.Department, nil)
}
