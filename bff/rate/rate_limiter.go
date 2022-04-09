package rate

import (
	"bff/client"
	"log"

	"github.com/gin-gonic/gin"
	libredis "github.com/go-redis/redis/v8"

	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimiter() gin.HandlerFunc {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted("10-S")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Create a redis client.
	option, err := libredis.ParseURL("redis://" + client.GetEnv("REDIS_HOST") + ":" + client.GetEnv("REDIS_PORT") + "/" + client.GetEnv("REDIS_DB"))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	client := libredis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_gin_example",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Create a new middleware with the limiter instance.
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}
