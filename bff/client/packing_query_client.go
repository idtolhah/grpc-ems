package client

import (
	"context"
	"errors"
	"strconv"

	"bff/pb/packingquerypb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Packing struct {
	Id           int64  `json:"id,omitempty"`
	FoId         string `json:"fo_id,omitempty"`
	LineId       int32  `json:"line_id,omitempty"`
	MachineId    int32  `json:"machine_id,omitempty"`
	UnitId       int32  `json:"unit_id,omitempty"`
	DepartmentId int32  `json:"department_id,omitempty"`
	AreaId       int32  `json:"area_id,omitempty"`
	CompletedAt  string `json:"completed_at,omitempty"`
	Status       int32  `json:"status,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
}

type PackingQueryClient struct {
}

var (
	_                             = utils.LoadLocalEnv()
	packingQueryGrpcService       = utils.GetEnv("PACKING_QUERY_GRPC_SERVICE")
	packingQueryGrpcServiceClient packingquerypb.PackingQueryServiceClient
)

func preparePackingGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()

	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, packingQueryGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		packingQueryGrpcServiceClient = nil
		return errors.New("connection to packing gRPC service failed")
	}

	if packingQueryGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	utils.CreateStartPromHttpServer(reg, 9097)

	packingQueryGrpcServiceClient = packingquerypb.NewPackingQueryServiceClient(conn)
	return nil
}

func (ac *PackingQueryClient) GetPackings(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := preparePackingGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := packingQueryGrpcServiceClient.GetPackings(ctx, &packingquerypb.GetPackingsRequest{
		Nophoto: c.Query("nophoto"), Page: c.Query("page"), Perpage: c.Query("perpage"), Status: c.Query("status"),
		LineId: c.Query("line_id"), MachineId: c.Query("machine_id"), Date: c.Query("date"),
		FoCondition: c.Query("fo_condition"), AoConclusion: c.Query("ao_conclusion"),
	})

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	var packings []Packing
	for _, d := range data.Packings {
		packings = append(packings, Packing{
			Id: d.Id, FoId: d.FoId, LineId: d.LineId, MachineId: d.MachineId, UnitId: d.UnitId, DepartmentId: d.DepartmentId,
			AreaId: d.AreaId, CompletedAt: d.CompletedAt, Status: d.Status, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		})
	}

	utils.ResponsePaged(c, packings, int(data.Total), int(data.Page), int(data.LastPage), err)
}

func (ac *PackingQueryClient) GetPacking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	param, err := utils.GetParam(c, "id")
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := preparePackingGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	id, _ := strconv.Atoi(param)
	data, err := packingQueryGrpcServiceClient.GetPacking(ctx, &packingquerypb.GetPackingRequest{Id: int64(id)})

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}
