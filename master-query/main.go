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

// Unit
func (*server) GetUnits(ctx context.Context, req *masterpb.GetUnitsRequest) (*masterpb.GetUnitsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterpb.Unit
	results, err := db_client.Query("SELECT id, name FROM refineries")
	if err != nil {
		return nil, err
	}
	var unit masterpb.Unit
	for results.Next() {
		err = results.Scan(&unit.Id, &unit.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, unit)
	}
	// Query: End

	var res masterpb.GetUnitsResponse
	for _, d := range data {
		res.Units = append(res.Units, &masterpb.Unit{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Units)
		redis.SendToRedisCacheDirect("units", string(stringData))
	}()

	return &res, nil
}

func (*server) GetUnit(ctx context.Context, req *masterpb.GetUnitRequest) (*masterpb.GetUnitResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var unit masterpb.Unit
	err := db_client.QueryRow(`SELECT id, name FROM refineries WHERE id=?`, req.Id).Scan(&unit.Id, &unit.Name)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetUnitResponse{Unit: &masterpb.Unit{Id: unit.Id, Name: unit.Name}}, nil
}

// Department
func (*server) GetDepartments(ctx context.Context, req *masterpb.GetDepartmentsRequest) (*masterpb.GetDepartmentsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterpb.Department
	results, err := db_client.Query("SELECT id, name FROM departments")
	if err != nil {
		return nil, err
	}
	var department masterpb.Department
	for results.Next() {
		err = results.Scan(&department.Id, &department.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, department)
	}
	// Query: End

	var res masterpb.GetDepartmentsResponse
	for _, d := range data {
		res.Departments = append(res.Departments, &masterpb.Department{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Departments)
		redis.SendToRedisCacheDirect("departments", string(stringData))
	}()

	return &res, nil
}

func (*server) GetDepartment(ctx context.Context, req *masterpb.GetDepartmentRequest) (*masterpb.GetDepartmentResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var department masterpb.Department
	err := db_client.QueryRow(`SELECT id, name FROM departments WHERE id=?`, req.Id).Scan(&department.Id, &department.Name)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetDepartmentResponse{Department: &masterpb.Department{Id: department.Id, Name: department.Name}}, nil
}

// Area
func (*server) GetAreas(ctx context.Context, req *masterpb.GetAreasRequest) (*masterpb.GetAreasResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterdb.Area
	results, err := db_client.Query("SELECT id, name FROM areas")
	if err != nil {
		return nil, err
	}
	var area masterdb.Area
	for results.Next() {
		err = results.Scan(&area.Id, &area.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, area)
	}
	// Query: End

	var res masterpb.GetAreasResponse
	for _, d := range data {
		res.Areas = append(res.Areas, &masterpb.Area{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Areas)
		redis.SendToRedisCacheDirect("areas", string(stringData))
	}()

	return &res, nil
}

func (*server) GetArea(ctx context.Context, req *masterpb.GetAreaRequest) (*masterpb.GetAreaResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var area masterpb.Area
	err := db_client.QueryRow(`SELECT id, name FROM areas WHERE id=?`, req.Id).Scan(&area.Id, &area.Name)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetAreaResponse{Area: &masterpb.Area{Id: area.Id, Name: area.Name}}, nil
}

// Line
func (*server) GetLines(ctx context.Context, req *masterpb.GetLinesRequest) (*masterpb.GetLinesResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterpb.Line
	results, err := db_client.Query("SELECT id, name FROM liness")
	if err != nil {
		return nil, err
	}
	var line masterpb.Line
	for results.Next() {
		err = results.Scan(&line.Id, &line.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, line)
	}
	// Query: End

	var res masterpb.GetLinesResponse
	for _, d := range data {
		res.Lines = append(res.Lines, &masterpb.Line{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Lines)
		redis.SendToRedisCacheDirect("lines", string(stringData))
	}()

	return &res, nil
}

func (*server) GetLine(ctx context.Context, req *masterpb.GetLineRequest) (*masterpb.GetLineResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var line masterpb.Line
	err := db_client.QueryRow(`SELECT id, name FROM liness WHERE id=?`, req.Id).Scan(&line.Id, &line.Name)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetLineResponse{Line: &masterpb.Line{Id: line.Id, Name: line.Name}}, nil
}

// Machine
func (*server) GetMachines(ctx context.Context, req *masterpb.GetMachinesRequest) (*masterpb.GetMachinesResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterpb.Machine
	results, err := db_client.Query("SELECT id, name FROM machines")
	if err != nil {
		return nil, err
	}
	var machine masterpb.Machine
	for results.Next() {
		err = results.Scan(&machine.Id, &machine.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, machine)
	}
	// Query: End

	var res masterpb.GetMachinesResponse
	for _, d := range data {
		res.Machines = append(res.Machines, &masterpb.Machine{Id: int32(d.Id), Name: d.Name})
	}

	go func() {
		stringData, _ := json.Marshal(res.Machines)
		redis.SendToRedisCacheDirect("machines", string(stringData))
	}()

	return &res, nil
}

func (*server) GetMachine(ctx context.Context, req *masterpb.GetMachineRequest) (*masterpb.GetMachineResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var machine masterpb.Machine
	err := db_client.QueryRow(`SELECT id, name FROM machines WHERE id=?`, req.Id).Scan(&machine.Id, &machine.Name)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetMachineResponse{Machine: &masterpb.Machine{Id: machine.Id, Name: machine.Name}}, nil
}

// Asset Equipment
func (*server) GetAssetEquipments(ctx context.Context, req *masterpb.GetAssetEquipmentsRequest) (*masterpb.GetAssetEquipmentsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterdb.AssetEquipment
	results, err := db_client.Query("SELECT id, item, item_check, checking_method, tools, standard_area, photo, line_id, machine_id FROM asset_equipments")
	if err != nil {
		return nil, err
	}
	var a masterdb.AssetEquipment
	for results.Next() {
		err = results.Scan(&a.ID, &a.Item, &a.ItemCheck, &a.CheckingMethod, &a.Tools, &a.StandardArea, &a.Photo, &a.LineID, &a.MachineID)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, a)
	}
	// Query: End

	var res masterpb.GetAssetEquipmentsResponse
	for _, d := range data {
		res.Assetequipments = append(res.Assetequipments, &masterpb.AssetEquipment{Id: int32(d.ID), Item: d.Item, ItemCheck: d.ItemCheck, CheckingMethod: d.CheckingMethod, Tools: d.Tools, StandardArea: d.StandardArea, Photo: d.Photo, LineId: int32(d.LineID), MachineId: int32(d.MachineID)})
	}

	go func() {
		stringData, _ := json.Marshal(res.Assetequipments)
		redis.SendToRedisCacheDirect("asset-equipments", string(stringData))
	}()

	return &res, nil
}

func (*server) GetAssetEquipment(ctx context.Context, req *masterpb.GetAssetEquipmentRequest) (*masterpb.GetAssetEquipmentResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var assetEquipment masterpb.AssetEquipment
	err := db_client.QueryRow(`SELECT id, item, item_check, checking_method, line_id, machine_id, standard_area, tools, photo FROM asset_equipments WHERE id=?`, req.Id).Scan(
		&assetEquipment.Id, &assetEquipment.Item, &assetEquipment.ItemCheck, &assetEquipment.CheckingMethod,
		&assetEquipment.LineId, &assetEquipment.MachineId, &assetEquipment.StandardArea, &assetEquipment.Tools,
		&assetEquipment.Photo,
	)
	if err != nil {
		return nil, err
	}
	// Query: End

	return &masterpb.GetAssetEquipmentResponse{Assetequipment: &masterpb.AssetEquipment{
		Id: assetEquipment.Id, Item: assetEquipment.Item, ItemCheck: assetEquipment.ItemCheck, CheckingMethod: assetEquipment.CheckingMethod,
		LineId: assetEquipment.LineId, MachineId: assetEquipment.MachineId, StandardArea: assetEquipment.StandardArea,
		Tools: assetEquipment.Tools, Photo: assetEquipment.Photo,
	}}, nil
}

// Contact
func (*server) GetContacts(ctx context.Context, req *masterpb.GetContactsRequest) (*masterpb.GetContactsResponse, error) {
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []masterdb.Contact
	results, err := db_client.Query("SELECT * FROM contacts")
	if err != nil {
		return nil, err
	}
	var contact masterdb.Contact
	for results.Next() {
		err = results.Scan(&contact.Id, &contact.Title, &contact.Number, &contact.OpTime, &contact.OpDay, &contact.Email)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, contact)
	}
	// Query: End

	var res masterpb.GetContactsResponse
	for _, d := range data {
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

	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVICE_HOST")+":"+utils.GetEnv("GRPC_SERVICE_PORT"))
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
