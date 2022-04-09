package cache

import (
	"bff/client"
	"bff/pb/redispb"
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	timeout      = 10 * time.Second
	redis_client client.RedisClient
)

func GetCache(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	param, err := client.GetParam(c, "key")
	if err != nil {
		client.Response(c, nil, err)
		return
	}

	data, err := redis_client.GetCache(&ctx, &redispb.GetRequest{Key: param})
	if err != nil {
		client.Response(c, nil, err)
		return
	}

	var jsonData interface{}
	json.Unmarshal([]byte(*data), &jsonData)
	// log.Printf("=========>%v\n", jsonData)

	client.Response(c, jsonData, err)
}

func GetCacheByKey(c *gin.Context, key string) interface{} {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := redis_client.GetCache(&ctx, &redispb.GetRequest{Key: key})
	if err != nil {
		client.Response(c, nil, err)
		return nil
	}

	var jsonData interface{}
	json.Unmarshal([]byte(*data), &jsonData)
	// log.Printf("=========>%v\n", jsonData)

	return jsonData
}
