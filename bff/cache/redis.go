package cache

import (
	"bff/utils"
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
	db, err := strconv.Atoi(utils.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Invalid db")
	}
	redisClient = redis.NewClient(&redis.Options{
		// Addr:     utils.GetEnv("REDIS_HOST") + ":" + utils.GetEnv("REDIS_PORT"),
		Addr:     utils.GetEnv("REDIS_HOST") + ":6379",
		Password: utils.GetEnv("REDIS_PWD"),
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
