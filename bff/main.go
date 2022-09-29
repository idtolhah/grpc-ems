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
	unit_client            client.UnitClient
	department_client      client.DepartmentClient
	area_client            client.AreaClient
	line_client            client.LineClient
	machine_client         client.MachineClient
	contact_client         client.ContactClient
	asset_equipment_client client.AssetEquipmentClient
	packing_query_client   client.PackingQueryClient
	packing_cmd_client     client.PackingCmdClient
	comment_cmd_client     client.CommentCmdClient
)

// Main
func main() {
	log.Println("Bff Service...")
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

	r.GET("/")

	api := r.Group("/api")
	// Login
	api.POST("/users/login", user_client.Login)
	// Master
	api.GET("/units", unit_client.GetUnits)
	api.GET("/units/:id", unit_client.GetUnit)
	api.GET("/departments", department_client.GetDepartments)
	api.GET("/departments/:id", department_client.GetDepartment)
	api.GET("/areas", area_client.GetAreas)
	api.GET("/areas/:id", area_client.GetArea)
	api.GET("/lines", line_client.GetLines)
	api.GET("/lines/:id", line_client.GetLine)
	api.GET("/machines", machine_client.GetMachines)
	api.GET("/machines/:id", machine_client.GetMachine)
	api.GET("/asset-equipments", asset_equipment_client.GetAssetEquipments)
	api.GET("/asset-equipments/:id", asset_equipment_client.GetAssetEquipment)
	api.GET("/contacts", contact_client.GetContacts)
	// Packing
	api.GET("/packings", packing_query_client.GetPackings)
	api.GET("/packings/:id", packing_query_client.GetPacking)
	api.GET("/packings/summary", packing_query_client.GetSummary)
	api.POST("/packings", packing_cmd_client.CreatePacking)
	api.PUT("/packings/equipment-checkings/:id/comment", comment_cmd_client.CreatePackingComment)
	api.POST("/packings/:id/equipment-checkings", packing_cmd_client.CreateEquipmentChecking)
	api.PUT("/packings/equipment-checkings/:id", packing_cmd_client.UpdateEquipmentChecking)
	// Filtration
	// api.GET("/filtrations", filtration_query_client.GetFiltrations)

	protected := api.Use(auth.IsAuthenticated())
	// Users
	protected.GET("/users", user_client.GetUsers)
	protected.GET("/users/profile", user_client.GetUserDetails)

	r.Run()
}
