package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"filtration-query/clients"
	"filtration-query/db"
	"filtration-query/pb/filtrationquerypb"
	"filtration-query/redis"
	"filtration-query/utils"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 1. Declare Vars
var (
	timeout   = 10 * time.Second
	db_client *sql.DB
)

// 2. Create gRPC server struct
type server struct {
	filtrationquerypb.UnimplementedFiltrationQueryServiceServer
}

// 3.0. Declare main function
func main() {
	go log.Println("Filtration Service")

	// 3.1. Create redis client
	redis.NewClient()

	// 3.2. Create DB connection
	db_client, err := db.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	go log.Println("Database connected")
	defer db_client.Close()

	// 3.3. Create network listener
	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVICE_HOST")+":"+utils.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	// 3.4. Create gRPC server
	s := grpc.NewServer()
	filtrationquerypb.RegisterFiltrationQueryServiceServer(s, &server{})

	// 3.5. Register server to Reflection
	reflection.Register(s)

	// 3.6. Start gRPC server
	go log.Printf("Server started at %v", lis.Addr().String())
	err = s.Serve(lis)
	if err != nil {
		go log.Println("ERROR:", err.Error())
	}
}

// 4.0. Declare GetFiltrations function
func (*server) GetFiltrations(ctx context.Context, req *filtrationquerypb.GetFiltrationsRequest) (*filtrationquerypb.GetFiltrationsResponse, error) {
	// 4.1. Create context with timeout
	go log.Println("Call GetFiltrations")
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 4.2. Create Query
	// Query: Start
	condition := "WHERE 1=1"
	if req.FilterId != "" {
		condition += " AND filter_id = " + req.FilterId
	}

	// Pagination: Start
	var total int
	var page = 1
	var perpage = 5
	var last_page = 1
	err := db_client.QueryRow(`SELECT COUNT(id) FROM filtrations ` + condition).Scan(&total)
	if err != nil {
		return nil, err
	}
	if req.Page != "" && req.Perpage != "" {
		page, _ = strconv.Atoi(req.Page)
		perpage, _ = strconv.Atoi(req.Perpage)
		offset := (page - 1) * perpage
		last_page = int(math.Ceil(float64(total) / float64(perpage)))
		condition += " LIMIT " + req.Perpage + " OFFSET " + strconv.Itoa(int(offset))
	}
	// Pagination: End

	var data []filtrationquerypb.Filtration
	results, err := db_client.Query(`
		SELECT id, COALESCE(fo_filter_condition, ''), COALESCE(ao_followed_up_date, ''), COALESCE(ao_conclusion, ''), COALESCE(completed_date, ''),
		status, createdAt, 
		unit_id, filter_id, dirt_id
		FROM filtrations ` + condition)
	if err != nil {
		return nil, err
	}

	var filtration filtrationquerypb.Filtration
	for results.Next() {
		err = results.Scan(
			&filtration.Id, &filtration.FoFilterCondition, &filtration.AoFollowedUpDate, &filtration.AoConclusion,
			&filtration.CompletedDate, &filtration.Status, &filtration.CreatedAt,
			&filtration.Unit.Id, &filtration.Filter.Id, &filtration.Dirt.Id,
		)
		if err != nil {
			go log.Println(err)
		}
		data = append(data, filtration)
	}
	// Query: End

	var res filtrationquerypb.GetFiltrationsResponse
	for _, d := range data {
		res.Filtrations = append(res.Filtrations, &filtrationquerypb.Filtration{
			Id: d.Id, FoFilterCondition: d.FoFilterCondition, AoFollowedUpDate: d.AoFollowedUpDate, AoConclusion: d.AoConclusion,
			CompletedDate: d.CompletedDate, Status: d.Status, CreatedAt: d.CreatedAt,
			Unit:   clients.GetUnit(d.Unit.Id),
			Filter: clients.GetFilter(d.Filter.Id),
			Dirt:   clients.GetDirt(d.Dirt.Id),
		})
	}
	res.Total = int64(total)
	res.Page = int64(page)
	res.LastPage = int64(last_page)

	if utils.GetEnv("USE_CACHE") == "yes" {
		go func() {
			byteData, _ := json.Marshal(res.Filtrations)
			redis.SendToRedisCacheDirect("filtrations?page="+req.Page+"&perpage="+req.Perpage, string(byteData))
		}()
		go func() {
			byteTotal, _ := json.Marshal(res.Total)
			redis.SendToRedisCacheDirect("filtrations-total?page="+req.Page+"&perpage="+req.Perpage, string(byteTotal))
		}()
		go func() {
			bytePage, _ := json.Marshal(res.Page)
			redis.SendToRedisCacheDirect("filtrations-page?page="+req.Page+"&perpage="+req.Perpage, string(bytePage))
		}()
		go func() {
			byteLastPage, _ := json.Marshal(res.LastPage)
			redis.SendToRedisCacheDirect("filtrations-last-page?page="+req.Page+"&perpage="+req.Perpage, string(byteLastPage))
		}()
	}

	return &res, nil
}
