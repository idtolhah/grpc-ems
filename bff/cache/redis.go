package cache

import (
	"bff/client"
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
)

// Direct
func GetCacheByKeyDirect(key string) interface{} {
	db, err := strconv.Atoi(client.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Invalid db")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     client.GetEnv("REDIS_HOST") + ":" + client.GetEnv("REDIS_PORT"),
		Password: client.GetEnv("REDIS_PWD"),
		DB:       db,
	})

	val, err := redisClient.Get(key).Result()
	if err != nil {
		log.Println(err)
	}

	var jsonData interface{}
	json.Unmarshal([]byte(val), &jsonData)

	return jsonData
}
