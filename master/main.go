package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"master/masterdb"
	"master/masterpb"
	"master/redis"
	"master/utils"
	"net"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	timeout   = time.Second
	db_client *sql.DB

	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_say_hello_method_handle_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
)

type server struct {
	masterpb.UnimplementedMasterServiceServer
}

func (*server) GetAreas(ctx context.Context, req *masterpb.GetAreasRequest) (*masterpb.GetAreasResponse, error) {
	// log.Println("Called GetAreas")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindAreas(db_client, c)
	if err != nil {
		return nil, utils.Error_response(err)
	}

	var res masterpb.GetAreasResponse
	for _, d := range *data {
		res.Areas = append(res.Areas, &masterpb.Area{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Areas)
		redis.SendToRedisCacheDirect("areas", string(stringData))
	}()

	return &res, nil
}

func (*server) GetAssetEquipments(ctx context.Context, req *masterpb.GetAssetEquipmentsRequest) (*masterpb.GetAssetEquipmentsResponse, error) {
	// log.Println("Called GetAssetEquipments")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindAssetEquipments(db_client, c)
	if err != nil {
		return nil, utils.Error_response(err)
	}

	var res masterpb.GetAssetEquipmentsResponse
	for _, d := range *data {
		res.Assetequipments = append(res.Assetequipments, &masterpb.AssetEquipment{Id: int32(d.ID), Item: d.Item, ItemCheck: d.ItemCheck, CheckingMethod: d.CheckingMethod, Tools: d.Tools, StandardArea: d.StandardArea, Photo: d.Photo, LineId: int32(d.LineID), MachineId: int32(d.MachineID)})
	}

	go func() {
		stringData, _ := json.Marshal(res.Assetequipments)
		redis.SendToRedisCacheDirect("asset-equipments", string(stringData))
	}()

	return &res, nil
}

func (*server) GetContacts(ctx context.Context, req *masterpb.GetContactsRequest) (*masterpb.GetContactsResponse, error) {
	// log.Println("Called GetContacts")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := masterdb.FindContacts(db_client, c)
	if err != nil {
		return nil, utils.Error_response(err)
	}

	var res masterpb.GetContactsResponse
	for _, d := range *data {
		res.Contacts = append(res.Contacts, &masterpb.Contact{Id: int32(d.Id), Title: d.Title, Number: d.Number, Optime: d.OpTime, Opday: d.OpDay, Email: d.Email})
	}

	go func() {
		stringData, _ := json.Marshal(res.Contacts)
		redis.SendToRedisCacheDirect("contacts", string(stringData))
	}()

	return &res, nil
}

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

// Main
func main() {
	log.Println("Master Service")
	redis.NewClient()

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

	// Prom: Create a HTTP server for prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9091)}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)
	masterpb.RegisterMasterServiceServer(s, &server{})

	// Prom: Initialize all metrics.
	grpcMetrics.InitializeMetrics(s)
	// Prom: Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// Register server to Reflection
	reflection.Register(s)

	log.Printf("Server started at %v", lis.Addr().String())
	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}
