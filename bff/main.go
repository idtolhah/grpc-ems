package main

import (
	"log"
	"time"

	allowlist "bff/allowlist"
	"bff/auth"
	"bff/client"
	"bff/rate"
	"bff/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	user_client            client.UserClient
	area_client            client.AreaClient
	contact_client         client.ContactClient
	asset_equipment_client client.AssetEquipmentClient
	packing_client         client.PackingClient
)

// Main
func main() {
	log.Println("Bff Service")
	r := gin.Default()

	// enable rate limiter per ip address
	if utils.GetEnv("USE_RATE_LIMITER") == "yes" {
		r.ForwardedByClientIP = true
		r.Use(rate.RateLimiter())
	}
	// if need to specify serveral range of allowed sources, use comma to concatenate them
	if utils.GetEnv("USE_IP_ALLOWLISTING") == "yes" {
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
	// Login
	api.POST("/users/login", user_client.Login)
	// Master
	api.GET("/areas", area_client.GetAreas)
	api.GET("/contacts", contact_client.GetContacts)
	api.GET("/asset-equipments", asset_equipment_client.GetAssetEquipments)

	protected := api.Use(auth.IsAuthenticated())
	// Users
	protected.GET("/users", user_client.GetUsers)
	protected.GET("/users/profile", user_client.GetUserDetails)
	// Packing
	protected.POST("/packings", packing_client.CreatePacking)
	protected.POST("/packings/:id/equipment-checkings", packing_client.CreateEquipmentChecking)
	protected.PUT("/packings/:id/equipment-checkings/:ecid", packing_client.UpdateEquipmentChecking)

	r.Run()
}
