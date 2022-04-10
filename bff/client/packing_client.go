package client

import (
	"context"
	"errors"
	"time"

	"bff/pb/packingpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Packing struct {
	Id           int64  `json:"id"`
	FoId         string `json:"fo_id"`
	LineId       int32  `json:"line_id"`
	MachineId    int32  `json:"machine_id"`
	UnitId       int32  `json:"unit_id"`
	DepartmentId int32  `json:"department_id"`
	AreaId       int32  `json:"area_id"`
	CompletedAt  string `json:"completed_at"`
	Status       int32  `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type CreatePackingRequest struct {
	UserId              string `json:"user_id"`
	LineId              int32  `json:"line_id"`
	MachineId           int32  `json:"machine_id"`
	StatusSync          int32  `json:"status_sync"`
	ObservationDatetime string `json:"observation_datetime"`
	UnitId              int32  `json:"unit_id"`
	DepartmentId        int32  `json:"department_id"`
	AreaId              int32  `json:"area_id"`
}

type PackingClient struct {
}

var (
	timeout                  = 10 * time.Second
	_                        = loadLocalEnv()
	packingGrpcService       = GetEnv("PACKING_GRPC_SERVICE")
	packingGrpcServiceClient packingpb.PackingServiceClient
)

func preparePackingGrpcClient(c *context.Context) error {

	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := GetRegistryMetrics()

	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, packingGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		packingGrpcServiceClient = nil
		return errors.New("connection to packing gRPC service failed")
	}

	if packingGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// CreateStartPromHttpServer(reg, 9094)

	packingGrpcServiceClient = packingpb.NewPackingServiceClient(conn)
	return nil
}

func (ac *PackingClient) CreatePacking(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req CreatePackingRequest
	err := c.BindJSON(&req)
	if err != nil {
		Response(c, nil, err)
		return
	}

	if err := preparePackingGrpcClient(&ctx); err != nil {
		return
	}

	data, err := packingGrpcServiceClient.CreatePacking(ctx, &packingpb.CreatePackingRequest{
		UserId:              req.UserId,
		LineId:              req.LineId,
		MachineId:           req.MachineId,
		StatusSync:          req.StatusSync,
		UnitId:              req.UnitId,
		DepartmentId:        req.DepartmentId,
		AreaId:              req.AreaId,
		ObservationDatetime: req.ObservationDatetime,
	})

	if err != nil {
		return
	}

	Response(c, data, err)
}
