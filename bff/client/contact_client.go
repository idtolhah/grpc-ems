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

type Contact struct {
	Id     int32  `json:"id"`
	Title  string `json:"name"`
	Number string `json:"number"`
	OpTime string `json:"optime"`
	OpDay  string `json:"opday"`
	Email  string `json:"email"`
}

type ContactClient struct {
}

var (
	_                        = utils.LoadLocalEnv()
	contactGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	contactGrpcServiceClient masterpb.MasterServiceClient
)

func prepareContactGrpcClient(c *context.Context) error {

	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, contactGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		contactGrpcServiceClient = nil
		return errors.New("connection to contact gRPC service failed")
	}

	if contactGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	utils.CreateStartPromHttpServer(reg, 9094)

	contactGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *ContactClient) GetContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	jsonData := cache.GetCacheByKeyDirect("contacts")
	if jsonData != nil && utils.GetEnv("USE_CACHE") == "yes" {
		utils.Response(c, jsonData, nil)
		return
	}

	if err := prepareContactGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := contactGrpcServiceClient.GetContacts(ctx, &masterpb.GetContactsRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var contacts []Contact
	for _, u := range res.GetContacts() {
		contacts = append(contacts, Contact{
			Id:     u.Id,
			Title:  u.Title,
			Number: u.Number,
			OpTime: u.Optime,
			OpDay:  u.Opday,
			Email:  u.Email,
		})
	}
	utils.Response(c, &contacts, err)
}
