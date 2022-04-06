package main

import (
	"context"
	"database/sql"
	"log"
	"master/masterdb"
	"master/masterpb"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	timeout   = time.Second
	db_client *sql.DB
)

type server struct {
	masterpb.UnimplementedMasterServiceServer
}

// Server functions
func (*server) GetAreas(ctx context.Context, req *masterpb.GetAreasRequest) (*masterpb.GetAreasResponse, error) {
	log.Println("Called GetAreas")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindAreas(db_client, c)
	if err != nil {
		return nil, error_response(err)
	}

	var res masterpb.GetAreasResponse
	for _, d := range *data {
		res.Areas = append(res.Areas, &masterpb.Area{Id: int32(d.Id), Name: d.Name})
	}

	return &res, nil
}

func (*server) GetAssetEquipments(ctx context.Context, req *masterpb.GetAssetEquipmentsRequest) (*masterpb.GetAssetEquipmentsResponse, error) {
	log.Println("Called GetAssetEquipments")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindAssetEquipments(db_client, c)
	if err != nil {
		return nil, error_response(err)
	}

	var res masterpb.GetAssetEquipmentsResponse
	for _, d := range *data {
		res.Assetequipments = append(res.Assetequipments, &masterpb.AssetEquipment{Id: int32(d.ID), Item: d.Item, ItemCheck: d.ItemCheck, CheckingMethod: d.CheckingMethod, Tools: d.Tools, StandardArea: d.StandardArea, Photo: d.Photo, LineId: int32(d.LineID), MachineId: int32(d.MachineID)})
	}

	return &res, nil
}

func (*server) GetContacts(ctx context.Context, req *masterpb.GetContactsRequest) (*masterpb.GetContactsResponse, error) {
	log.Println("Called GetContacts")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindContacts(db_client, c)
	if err != nil {
		return nil, error_response(err)
	}

	var res masterpb.GetContactsResponse
	for _, d := range *data {
		res.Contacts = append(res.Contacts, &masterpb.Contact{Id: int32(d.Id), Title: d.Title, Number: d.Number, Optime: d.OpTime, Opday: d.OpDay, Email: d.Email})
	}

	return &res, nil
}

// Utils
func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

// Main
func main() {
	log.Println("Master Service")

	lis, err := net.Listen("tcp", masterdb.GetEnv("GRPC_SERVICE_HOST")+":"+masterdb.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	db_client, err = masterdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db_client.Close()

	s := grpc.NewServer()
	masterpb.RegisterMasterServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}

}
