package client

import (
	"context"
	"errors"

	"bff/pb/masterpb"

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
	_                               = loadLocalEnv()
	assetEquipmentGrpcService       = GetEnv("MASTER_GRPC_SERVICE")
	assetEquipmentGrpcServiceClient masterpb.MasterServiceClient
)

func prepareMasterGrpcClient(c *context.Context) error {

	// Prom: Get Registry & Metrics
	reg, grpcMetrics := GetRegistryMetrics()
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
	CreateStartPromHttpServer(reg, 9093)

	assetEquipmentGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *AssetEquipmentClient) GetAssetEquipments(c *context.Context) (*[]AssetEquipment, error) {

	if err := prepareMasterGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := assetEquipmentGrpcServiceClient.GetAssetEquipments(*c, &masterpb.GetAssetEquipmentsRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
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
	return &assetEquipments, nil
}
