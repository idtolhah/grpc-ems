package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"net"
	"packing/packingdb"
	"packing/packingpb"
	"packing/redis"
	"packing/utils"
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
	packingpb.UnimplementedPackingServiceServer
}

func (*server) GetPackings(ctx context.Context, req *packingpb.GetPackingsRequest) (*packingpb.GetPackingsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []packingpb.Packing
	results, err := db_client.Query("SELECT * FROM packings")
	if err != nil {
		return nil, err
	}
	var packing packingpb.Packing
	for results.Next() {
		err = results.Scan(
			&packing.Id, &packing.FoId, &packing.LineId, &packing.MachineId, &packing.UnitId, &packing.DepartmentId,
			&packing.AreaId, &packing.CompletedAt, &packing.Status, &packing.CreatedAt, &packing.UpdatedAt,
		)
		if err != nil {
			// panic(err.Error())
			log.Println(err)
		}
		data = append(data, packing)
	}
	// Query: End

	var res packingpb.GetPackingsResponse
	for _, d := range data {
		res.Packings = append(res.Packings, &packingpb.Packing{
			Id: d.Id, FoId: d.FoId, LineId: d.LineId, MachineId: d.MachineId, UnitId: d.UnitId, DepartmentId: d.DepartmentId,
			AreaId: d.AreaId, CompletedAt: d.CompletedAt, Status: d.Status, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
		})
	}

	go func() {
		stringData, _ := json.Marshal(res.Packings)
		redis.SendToRedisCacheDirect("packings", string(stringData))
	}()

	return &res, nil
}

func (*server) GetPacking(ctx context.Context, req *packingpb.GetPackingRequest) (*packingpb.GetPackingResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var packing packingpb.Packing
	err := db_client.QueryRow("select * from packings where id=?", req.Id).Scan(&packing.Id)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &packingpb.GetPackingResponse{Packing: &packing, EquipmentCheckings: []*packingpb.EquipmentChecking{}}, nil
}

func (*server) CreatePacking(ctx context.Context, req *packingpb.CreatePackingRequest) (*packingpb.CreatePackingResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := db_client.Exec(
		"insert into packings(fo_id, line_id, machine_id, status, unit_id, department_id, area_id, createdAt, updatedAt) values(?,?,?,?,?,?,?,?,?)",
		req.UserId, req.LineId, req.MachineId, req.StatusSync, req.UnitId, req.DepartmentId, req.AreaId, req.ObservationDatetime, req.ObservationDatetime,
	)

	if err != nil {
		return nil, err
	}

	lastId, errId := res.LastInsertId()
	if errId != nil {
		return nil, errId
	}

	return &packingpb.CreatePackingResponse{Id: lastId}, nil
}

func (*server) CreateEquipmentChecking(ctx context.Context, req *packingpb.CreateEquipmentCheckingRequest) (*packingpb.CreateEquipmentCheckingResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var packing packingpb.Packing
	err := db_client.QueryRow("select fo_id from packings where id=?", req.IdPackagingCheck).Scan(&packing.FoId)
	if err != nil {
		return nil, err
	}
	// Query: End

	// Query: Start
	res, err := db_client.Exec(
		"insert into equipment_checkings(id_equipment_checking_list, packing_id, asset_equipment_id, fo_photo, fo_condition, fo_note, fo_id, createdAt, updatedAt) values(?,?,?,?,?,?,?,?,?)",
		req.IdEquipmentCheckingList, req.IdPackagingCheck, req.IdAssetEquipment, req.Photo, req.Condition, req.Note, packing.FoId,
		req.ObservationDatetime, req.ObservationDatetime,
	)
	if err != nil {
		return nil, err
	}
	// Query: End

	lastId, errId := res.LastInsertId()
	if errId != nil {
		return nil, errId
	}

	return &packingpb.CreateEquipmentCheckingResponse{Id: lastId}, nil
}

func (*server) UpdateEquipmentChecking(ctx context.Context, req *packingpb.UpdateEquipmentCheckingRequest) (*packingpb.UpdateEquipmentCheckingResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := db_client.Exec(
		"update equipment_checkings set ao_conclusion=?, ao_note=?, ao_id=?, ao_created_at=?, updatedAt=? where id=?",
		req.AoConclusion, req.AoNote, req.AoId, req.AoObservationDatetime, req.AoObservationDatetime, req.Id,
	)

	if err != nil {
		return nil, err
	}

	lastId, errId := res.RowsAffected()
	if errId != nil {
		return nil, errId
	}

	return &packingpb.UpdateEquipmentCheckingResponse{Id: lastId}, nil
}

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

// Main
func main() {
	log.Println("Packing Service")
	redis.NewClient()

	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVICE_HOST")+":"+utils.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	db_client, err = packingdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db_client.Close()

	// Prom: Create a HTTP server for prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9096)}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)
	packingpb.RegisterPackingServiceServer(s, &server{})

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
