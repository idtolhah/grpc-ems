package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"redis-cache/client"
	"redis-cache/redispb"
	"strconv"

	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Redis Connection Client
var redisClient *redis.Client
var _ = client.LoadLocalEnv()

type server struct {
	redispb.UnimplementedRedisServiceServer
}

func NewClient() {
	db, err := strconv.Atoi(client.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Invalid db")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     client.GetEnv("REDIS_HOST") + ":" + client.GetEnv("REDIS_PORT"),
		Password: client.GetEnv("REDIS_PWD"),
		DB:       db,
	})
	log.Printf("New Client %v\n", redisClient)
	redisClient.FlushAllAsync()
	log.Println("All Cache Flushed")
}

func (*server) SetCache(c context.Context, req *redispb.SetRequest) (*redispb.SetResponse, error) {
	log.Println("Called SetCache")

	err := redisClient.Set(req.Key, req.Value, 0).Err()
	if err != nil {
		fmt.Println(err)
	}

	val, err := redisClient.Get(req.Key).Result()
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("=======>%v\n", val)

	return &redispb.SetResponse{Data: val}, nil
}

func (*server) GetCache(c context.Context, req *redispb.GetRequest) (*redispb.GetResponse, error) {
	log.Println("Called GetCache")

	val, err := redisClient.Get(req.Key).Result()
	if err != nil {
		fmt.Println(err)
	}
	// log.Printf("=======>%v\n", val)

	var res redispb.GetResponse
	res.Data = val

	return &res, nil
}

// Main
func main() {
	// Grpc
	log.Println("Redis Service")
	NewClient()

	lis, err := net.Listen("tcp", client.GetEnv("GRPC_SERVICE_HOST")+":"+client.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	s := grpc.NewServer()
	redispb.RegisterRedisServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}
