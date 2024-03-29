package redis

import (
	"log"
	"master/utils"
	"strconv"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func NewClient() {
	db, err := strconv.Atoi(utils.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatalln("Invalid db")
	}

	go log.Printf("REDIS_HOST====>%v", utils.GetEnv("REDIS_HOST"))
	go log.Printf("REDIS_PORT====>%v", utils.GetEnv("REDIS_PORT"))

	redisClient = redis.NewClient(&redis.Options{
		// Addr:     utils.GetEnv("REDIS_HOST") + ":" + utils.GetEnv("REDIS_PORT"),
		Addr:     utils.GetEnv("REDIS_HOST") + ":6379",
		Password: utils.GetEnv("REDIS_PWD"),
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
