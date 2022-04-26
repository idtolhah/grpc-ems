package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"

	"net"
	"packing/db"
	"packing/packingcmdpb"
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

type NullString string

func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return errors.New("Column is not a string")
	}
	*s = NullString(strVal)
	return nil
}
func (s NullString) Value() (driver.Value, error) {
	if len(s) == 0 { // if nil or empty string
		return nil, nil
	}
	return string(s), nil
}

type server struct {
	packingcmdpb.UnimplementedPackingCmdServiceServer
}

func (*server) CreatePacking(ctx context.Context, req *packingcmdpb.CreatePackingRequest) (*packingcmdpb.CreatePackingResponse, error) {
	go log.Println("Call CreatePacking")
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

	if utils.GetEnv("USE_CACHE") == "yes" {
		redis.DeleteRedisCacheDirect("packings")
	}

	return &packingcmdpb.CreatePackingResponse{Id: lastId}, nil
}

func (*server) CreateEquipmentChecking(ctx context.Context, req *packingcmdpb.CreateEquipmentCheckingRequest) (*packingcmdpb.CreateEquipmentCheckingResponse, error) {
	go log.Println("Call CreateEquipmentChecking")
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var packing db.Packing
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

	return &packingcmdpb.CreateEquipmentCheckingResponse{Id: lastId}, nil
}

func (*server) UpdateEquipmentChecking(ctx context.Context, req *packingcmdpb.UpdateEquipmentCheckingRequest) (*packingcmdpb.UpdateEquipmentCheckingResponse, error) {
	go log.Println("Call UpdateEquipmentChecking")
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

	return &packingcmdpb.UpdateEquipmentCheckingResponse{Id: lastId}, nil
}

func init() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues("Test")
}

// Main
func main() {
	log.Println("Packing Command Service")
	redis.NewClient()

	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVICE_HOST")+":"+utils.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	db_client, err = db.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db_client.Close()

	// Prom: Create a HTTP server for prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", 9098)}

	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		grpc.UnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
	)
	packingcmdpb.RegisterPackingCmdServiceServer(s, &server{})

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
