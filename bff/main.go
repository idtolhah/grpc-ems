package main

import (
	"context"
	"errors"
	"log"
	"time"

	"bff/auth"
	"bff/client"

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

// Call functions
func GetUserDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	userId, ok := c.Get("UserId")
	if !ok {
		client.Response(c, nil, errors.New("invalid user id in token"))
		return
	}
	// fmt.Println("UserId from context: ", userId.(string))

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

func GetAreas(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := area_client.GetAreas(&ctx)
	client.Response(c, data, err)
}

func GetContacts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := contact_client.GetContacts(&ctx)
	client.Response(c, data, err)
}

func GetAssetEquipments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := asset_equipment_client.GetAssetEquipments(&ctx)
	client.Response(c, data, err)
}

// Main
func main() {
	log.Println("Bff Service")

	r := gin.Default()
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"PUT", "GET", "POST"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        24 * time.Hour,
	}))
	// r.GET("/", GetContacts)

	api := r.Group("/api")
	api.POST("/users/login", Login)
	api.GET("/areas", GetAreas)
	api.GET("/contacts", GetContacts)
	api.GET("/asset-equipments", GetAssetEquipments)

	protected := api.Use(auth.IsAuthenticated())
	protected.GET("/users", GetUsers)
	protected.GET("/users/profile", GetUserDetails)

	// consul.RegisterServiceWithConsul("localhost", 50051, "bff-service")
	// http.HandleFunc("/healthcheck", consul.Healthcheck)
	// consul.DeregisterService("bff-service")

	r.Run()
}
