package cache

import (
	"bff/client"
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	// redis_client client.RedisClient
	// timeout      = 10 * time.Second
	redisClient *redis.Client
)

// Via gRPC
// func GetCache(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(c, timeout)
// 	defer cancel()

// 	param, err := client.GetParam(c, "key")
// 	if err != nil {
// 		client.Response(c, nil, err)
// 		return
// 	}

// 	data, err := redis_client.GetCache(&ctx, &redispb.GetRequest{Key: param})
// 	if err != nil {
// 		client.Response(c, nil, err)
// 		return
// 	}

// 	var jsonData interface{}
// 	json.Unmarshal([]byte(*data), &jsonData)
// 	// log.Printf("=========>%v\n", jsonData)

// 	client.Response(c, jsonData, err)
// }

// func GetCacheByKey(c *gin.Context, key string) interface{} {
// 	ctx, cancel := context.WithTimeout(c, timeout)
// 	defer cancel()

// 	data, err := redis_client.GetCache(&ctx, &redispb.GetRequest{Key: key})
// 	if err != nil {
// 		client.Response(c, nil, err)
// 		return nil
// 	}

// 	var jsonData interface{}
// 	json.Unmarshal([]byte(*data), &jsonData)
// 	// log.Printf("=========>%v\n", jsonData)

// 	return jsonData
// }

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
