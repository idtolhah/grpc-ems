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

type AssetEquipment struct {
	Id             uint   `json:"id"`
	Item           string `json:"item"`
	ItemCheck      string `json:"item_check"`
	CheckingMethod string `json:"checking_method"`
	Tools          string `json:"tools"`
	StandardArea   string `json:"standard_area"`
	Photo          string `json:"photo"`
	LineID         uint   `json:"line_id"`
	MachineID      uint   `json:"machine_id"`
}

type AssetEquipmentClient struct {
}

var (
	_                               = utils.LoadLocalEnv()
	assetEquipmentGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	assetEquipmentGrpcServiceClient masterpb.MasterServiceClient
)

func prepareMasterGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, assetEquipmentGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		assetEquipmentGrpcServiceClient = nil
		return errors.New("connection to Master gRPC service failed")
	}

	if assetEquipmentGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	utils.CreateStartPromHttpServer(reg, 9093)

	assetEquipmentGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *AssetEquipmentClient) GetAssetEquipments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	jsonData := cache.GetCacheByKeyDirect("asset-equipments")
	if jsonData != nil && utils.GetEnv("USE_CACHE") == "yes" {
		utils.Response(c, jsonData, nil)
		return
	}

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := assetEquipmentGrpcServiceClient.GetAssetEquipments(ctx, &masterpb.GetAssetEquipmentsRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var assetEquipments []AssetEquipment
	for _, u := range res.GetAssetequipments() {
		assetEquipments = append(assetEquipments, AssetEquipment{
			Id:             uint(u.Id),
			Item:           u.Item,
			ItemCheck:      u.ItemCheck,
			CheckingMethod: u.CheckingMethod,
			Tools:          u.Tools,
			StandardArea:   u.StandardArea,
			Photo:          u.Photo,
			LineID:         uint(u.LineId),
			MachineID:      uint(u.MachineId),
		})
	}
	utils.Response(c, &assetEquipments, nil)
}
