package client

import (
	"context"
	"errors"

	"bff/pb/masterpb"

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
	_                        = loadLocalEnv()
	contactGrpcService       = GetEnv("MASTER_GRPC_SERVICE")
	contactGrpcServiceClient masterpb.MasterServiceClient
)

func prepareContactGrpcClient(c *context.Context) error {

	// Prom: Get Registry & Metrics
	reg, grpcMetrics := GetRegistryMetrics()
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
	CreateStartPromHttpServer(reg, 9094)

	contactGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func (uc *ContactClient) GetContacts(c *context.Context) (*[]Contact, error) {

	if err := prepareContactGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := contactGrpcServiceClient.GetContacts(*c, &masterpb.GetContactsRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
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
	return &contacts, nil
}
