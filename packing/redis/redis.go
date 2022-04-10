package redis

import (
	"log"
	"packing/packingdb"
	"packing/utils"
	"strconv"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func NewClient() {
	db, err := strconv.Atoi(packingdb.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Invalid db")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     packingdb.GetEnv("REDIS_HOST") + ":" + packingdb.GetEnv("REDIS_PORT"),
		Password: packingdb.GetEnv("REDIS_PWD"),
		DB:       db,
	})
	log.Printf("New Client %v\n", redisClient)
	redisClient.FlushAllAsync()
	log.Println("All Cache Flushed")
}

func SendToRedisCacheDirect(key string, val string) {
	err := redisClient.Set(key, val, 0).Err()
	if err != nil {
		utils.Error_response(err)
		return
	}
	log.Println("Data Cached")
}
