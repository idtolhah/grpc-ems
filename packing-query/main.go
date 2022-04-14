package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"net"
	"packing/db"
	"packing/pb/masterpb"
	"packing/pb/packingquerypb"
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

	_                       = utils.LoadLocalEnv()
	masterGrpcService       = utils.GetEnv("MASTER_GRPC_SERVICE")
	masterGrpcServiceClient masterpb.MasterServiceClient
)

type server struct {
	packingquerypb.UnimplementedPackingQueryServiceServer
}

func prepareUnitGrpcClient(c *context.Context) error {
	conn, err := grpc.DialContext(*c, masterGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...,
	)

	if err != nil {
		masterGrpcServiceClient = nil
		return errors.New("connection to unit gRPC service failed")
	}

	if masterGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	masterGrpcServiceClient = masterpb.NewMasterServiceClient(conn)
	return nil
}

func GetUnit(id int32) *packingquerypb.Unit {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Unit{}
	}

	res, err := masterGrpcServiceClient.GetUnit(c, &masterpb.GetUnitRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Unit{}
	}

	return &packingquerypb.Unit{Id: res.Unit.Id, Name: res.Unit.Name}
}

func GetLine(id int32) *packingquerypb.Line {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Line{}
	}

	res, err := masterGrpcServiceClient.GetLine(c, &masterpb.GetLineRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Line{}
	}

	return &packingquerypb.Line{Id: res.Line.Id, Name: res.Line.Name}
}

func GetMachine(id int32) *packingquerypb.Machine {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Machine{}
	}

	res, err := masterGrpcServiceClient.GetMachine(c, &masterpb.GetMachineRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Machine{}
	}

	return &packingquerypb.Machine{Id: res.Machine.Id, Name: res.Machine.Name}
}

func GetEquipmentCheckings(id int64) []*packingquerypb.EquipmentChecking {
	var data []*packingquerypb.EquipmentChecking
	results, err := db_client.Query(`SELECT ao_created_at FROM equipment_checkings WHERE id=?`, id)
	if err != nil {
		return []*packingquerypb.EquipmentChecking{}
	}
	var equipment_checking packingquerypb.EquipmentChecking
	for results.Next() {
		err = results.Scan(
			&equipment_checking.AoCreatedAt,
		)
		if err != nil {
			log.Println(err)
		}
		data = append(data, &equipment_checking)
	}
	return data
}

func (*server) GetPackings(ctx context.Context, req *packingquerypb.GetPackingsRequest) (*packingquerypb.GetPackingsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	condition := "WHERE 1=1"
	if req.LineId != "" {
		condition += " AND line_id = " + req.LineId
	}
	if req.MachineId != "" {
		condition += " AND machine_id = " + req.MachineId
	}

	// Pagination: Start
	var total int
	var page = 1
	var perpage = 5
	var last_page = 1
	err := db_client.QueryRow(`SELECT COUNT(id) FROM packings ` + condition).Scan(&total)
	if err != nil {
		return nil, err
	}
	if req.Page != "" && req.Perpage != "" {
		page, _ = strconv.Atoi(req.Page)
		perpage, _ = strconv.Atoi(req.Perpage)
		offset := (page - 1) * perpage
		last_page = int(math.Ceil(float64(total / perpage)))
		condition += " LIMIT " + req.Perpage + " OFFSET " + strconv.Itoa(int(offset))
	}
	// Pagination: End

	var data []packingquerypb.Packing
	results, err := db_client.Query(`
		SELECT id, fo_id, line_id, machine_id, unit_id, department_id, area_id, COALESCE(completed_at,''), status, 
		createdAt, updatedAt FROM packings ` + condition)
	if err != nil {
		return nil, err
	}

	var packing packingquerypb.Packing
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

	var res packingquerypb.GetPackingsResponse
	for _, d := range data {
		res.Packings = append(res.Packings, &packingquerypb.Packing{
			Id: d.Id, FoId: d.FoId, LineId: d.LineId, MachineId: d.MachineId, UnitId: d.UnitId, DepartmentId: d.DepartmentId,
			AreaId: d.AreaId, CompletedAt: d.CompletedAt, Status: d.Status, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt,
			EquipmentCheckings: GetEquipmentCheckings(d.Id),
			Unit:               GetUnit(d.UnitId),
			Line:               GetLine(d.LineId),
			Machine:            GetMachine(d.MachineId),
		})
	}
	res.Total = int64(total)
	res.Page = int64(page)
	res.LastPage = int64(last_page)

	go func() {
		stringData, _ := json.Marshal(res.Packings)
		redis.SendToRedisCacheDirect("packings", string(stringData))
	}()

	return &res, nil
}

func (*server) GetPacking(ctx context.Context, req *packingquerypb.GetPackingRequest) (*packingquerypb.GetPackingResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var packing packingquerypb.Packing
	err := db_client.QueryRow(`
		select id, fo_id, line_id, machine_id, unit_id, department_id, area_id, COALESCE(completed_at,''), status, 
		createdAt, updatedAt from packings where id=?`, req.Id).Scan(
		&packing.Id, &packing.FoId, &packing.LineId, &packing.MachineId, &packing.UnitId, &packing.DepartmentId,
		&packing.AreaId, &packing.CompletedAt, &packing.Status, &packing.CreatedAt, &packing.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	var equipmentCheckings []*packingquerypb.EquipmentChecking
	results, err := db_client.Query(`
		SELECT id, id_equipment_checking_list, packing_id, asset_equipment_id, 
		fo_id,  fo_photo, fo_condition, COALESCE(fo_note,''), 
		COALESCE(ao_id,''), COALESCE(ao_conclusion, 0), COALESCE(ao_note,''), COALESCE(ao_created_at,''), 
		COALESCE(mo_id,''), COALESCE(mo_note,''), COALESCE(mo_repair_photo,''), COALESCE(mo_created_at,''), 
		COALESCE(mr_id,''), COALESCE(mr_comment,''), COALESCE(mr_created_at,''), COALESCE(createdAt,''), COALESCE(updatedAt,'') 
		FROM equipment_checkings where packing_id=?`, packing.Id)
	if err != nil {
		return nil, err
	}

	var ec packingquerypb.EquipmentChecking
	for results.Next() {
		err = results.Scan(
			&ec.Id, &ec.IdEquipmentCheckingList, &ec.PackingId, &ec.AssetEquipmentId, &ec.FoId, &ec.FoPhoto, &ec.FoCondition, &ec.FoNote,
			&ec.AoId, &ec.AoConclusion, &ec.AoNote, &ec.AoCreatedAt, &ec.MoId, &ec.MoRepairPhoto, &ec.MoNote, &ec.MoCreatedAt,
			&ec.MrId, &ec.MrComment, &ec.MrCreatedAt, &ec.CreatedAt, &ec.UpdatedAt,
		)
		if err != nil {
			log.Println(err)
		}
		equipmentCheckings = append(equipmentCheckings, &ec)
	}

	if err != nil {
		return nil, err
	}
	// Query: End

	return &packingquerypb.GetPackingResponse{
		Id: packing.Id, FoId: packing.FoId, LineId: packing.LineId, MachineId: packing.MachineId, UnitId: packing.UnitId,
		DepartmentId: packing.DepartmentId, AreaId: packing.AreaId, CompletedAt: packing.CompletedAt, Status: packing.Status,
		CreatedAt: packing.CreatedAt, UpdatedAt: packing.UpdatedAt,
		EquipmentCheckings: equipmentCheckings,
	}, nil
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

	db_client, err = db.NewClient(context.Background())
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
	packingquerypb.RegisterPackingQueryServiceServer(s, &server{})

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
