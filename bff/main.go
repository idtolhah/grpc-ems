package main

import (
	"context"
	"errors"
	"log"
	"time"

	allowlist "bff/allowlistt"
	"bff/auth"
	"bff/cache"
	"bff/client"
	"bff/rate"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	timeout                = 10 * time.Second
	user_client            client.UserClient
	area_client            client.AreaClient
	contact_client         client.ContactClient
	asset_equipment_client client.AssetEquipmentClient
)

// Call controller functions
// Users
func GetUserDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	userId, ok := c.Get("UserId") // Get UserId from gin context after jwt auth
	if !ok {
		client.Response(c, nil, errors.New("invalid user id in token"))
		return
	}

	data, err := user_client.GetUserDetails(userId.(string), &ctx)
	if err != nil {
		client.Response(c, nil, err)
		return
	}

	client.Response(c, data, err)
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req client.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		client.Response(c, nil, err)
		return
	}

	data, err := user_client.Login(req.Email, req.Password, &ctx)
	if err != nil {
		client.Response(c, nil, err)
		return
	}

	client.Response(c, data, err)
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := user_client.GetUsers(&ctx)
	client.Response(c, data, err)
}

// Master
func GetAreas(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if jsonData := cache.GetCacheByKey(c, "areas"); jsonData != nil {
		client.Response(c, jsonData, nil)
		return
	}

	data, err := area_client.GetAreas(&ctx)
	client.Response(c, data, err)
}

func GetContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if jsonData := cache.GetCacheByKey(c, "contacts"); jsonData != nil {
		client.Response(c, jsonData, nil)
		return
	}

	data, err := contact_client.GetContacts(&ctx)
	client.Response(c, data, err)
}

func GetAssetEquipments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if jsonData := cache.GetCacheByKey(c, "asset-equipments"); jsonData != nil {
		client.Response(c, jsonData, nil)
		return
	}

	data, err := asset_equipment_client.GetAssetEquipments(&ctx)
	client.Response(c, data, err)
}

// Main
func main() {
	log.Println("Bff Service")
	r := gin.Default()

	// enable rate limiter per ip address
	if client.GetEnv("USE_RATE_LIMITER") == "yes" {
		r.ForwardedByClientIP = true
		r.Use(rate.RateLimiter())
	}
	// if need to specify serveral range of allowed sources, use comma to concatenate them
	if client.GetEnv("USE_IP_ALLOWLISTING") == "yes" {
		r.Use(allowlist.CIDR("172.18.0.0/16, 127.0.0.1/32, 192.168.43.1/32"))
	}

	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"PUT", "GET", "POST"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        24 * time.Hour,
	}))

	api := r.Group("/api")
	api.POST("/users/login", Login)
	api.GET("/areas", GetAreas)
	api.GET("/contacts", GetContacts)
	api.GET("/asset-equipments", GetAssetEquipments)
	// api.GET("/cache/:key", GetCache)

	protected := api.Use(auth.IsAuthenticated())
	protected.GET("/users", GetUsers)
	protected.GET("/users/profile", GetUserDetails)

	r.Run()
}
