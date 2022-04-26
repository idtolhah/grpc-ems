package client

import (
	"context"
	"errors"
	"strconv"

	"bff/pb/packingcmdpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type PackingCmdClient struct {
}

var (
	_                           = utils.LoadLocalEnv()
	packingcmdGrpcService       = utils.GetEnv("PACKING_CMD_GRPC_SERVICE")
	packingcmdGrpcServiceClient packingcmdpb.PackingCmdServiceClient
)

func preparePackingCmdGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()

	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, packingcmdGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		packingcmdGrpcServiceClient = nil
		return errors.New("connection to packingcmd gRPC service failed")
	}

	if packingcmdGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	utils.CreateStartPromHttpServer(reg, 9099)

	packingcmdGrpcServiceClient = packingcmdpb.NewPackingCmdServiceClient(conn)
	return nil
}

func (ac *PackingCmdClient) CreatePacking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req packingcmdpb.CreatePackingRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := preparePackingCmdGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := packingcmdGrpcServiceClient.CreatePacking(
		ctx,
		&packingcmdpb.CreatePackingRequest{
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

func (ac *PackingCmdClient) CreateEquipmentChecking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req packingcmdpb.CreateEquipmentCheckingRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := preparePackingCmdGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := packingcmdGrpcServiceClient.CreateEquipmentChecking(
		ctx,
		&packingcmdpb.CreateEquipmentCheckingRequest{
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

func (ac *PackingCmdClient) UpdateEquipmentChecking(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	param, err := utils.GetParam(c, "id")
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	var req packingcmdpb.UpdateEquipmentCheckingRequest
	errReq := c.BindJSON(&req)
	if errReq != nil {
		utils.Response(c, nil, errReq)
		return
	}

	if err := preparePackingCmdGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	id, _ := strconv.Atoi(param)
	data, err := packingcmdGrpcServiceClient.UpdateEquipmentChecking(
		ctx,
		&packingcmdpb.UpdateEquipmentCheckingRequest{
			Id: int64(id), AoConclusion: req.AoConclusion, AoNote: req.AoNote, AoId: req.AoId,
			AoObservationDatetime: req.AoObservationDatetime,
		},
	)

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}
