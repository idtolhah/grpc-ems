package client

import (
	"context"
	"errors"
	"strconv"

	"bff/pb/packingpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type PackingClient struct {
}

var (
	_                        = utils.LoadLocalEnv()
	packingGrpcService       = utils.GetEnv("PACKING_GRPC_SERVICE")
	packingGrpcServiceClient packingpb.PackingServiceClient
)

func preparePackingGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()

	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, packingGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
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
	utils.CreateStartPromHttpServer(reg, 9097)

	packingGrpcServiceClient = packingpb.NewPackingServiceClient(conn)
	return nil
}

func (ac *PackingClient) CreatePacking(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req packingpb.CreatePackingRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := preparePackingGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := packingGrpcServiceClient.CreatePacking(
		ctx,
		&packingpb.CreatePackingRequest{
			UserId: req.UserId, LineId: req.LineId, MachineId: req.MachineId, StatusSync: req.StatusSync, UnitId: req.UnitId,
			DepartmentId: req.DepartmentId, AreaId: req.AreaId, ObservationDatetime: req.ObservationDatetime,
		},
	)

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}

func (ac *PackingClient) CreateEquipmentChecking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req packingpb.CreateEquipmentCheckingRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := preparePackingGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := packingGrpcServiceClient.CreateEquipmentChecking(
		ctx,
		&packingpb.CreateEquipmentCheckingRequest{
			IdEquipmentCheckingList: req.IdEquipmentCheckingList, IdPackagingCheck: req.IdPackagingCheck,
			IdAssetEquipment: req.IdAssetEquipment, Photo: req.Photo, Condition: req.Condition, Note: req.Note,
			Status: req.Status, ObservationDatetime: req.ObservationDatetime,
		},
	)

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}

func (ac *PackingClient) UpdateEquipmentChecking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	param, err := utils.GetParam(c, "ecid")
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	var req packingpb.UpdateEquipmentCheckingRequest
	errReq := c.BindJSON(&req)
	if errReq != nil {
		utils.Response(c, nil, errReq)
		return
	}

	if err := preparePackingGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	ecid, _ := strconv.Atoi(param)
	data, err := packingGrpcServiceClient.UpdateEquipmentChecking(
		ctx,
		&packingpb.UpdateEquipmentCheckingRequest{
			Id: int64(ecid), AoConclusion: req.AoConclusion, AoNote: req.AoNote, AoId: req.AoId,
			AoObservationDatetime: req.AoObservationDatetime,
		},
	)

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}
