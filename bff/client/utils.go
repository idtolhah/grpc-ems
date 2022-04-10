package client

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func loadLocalEnv() interface{} {
	_, runningInContainer := os.LookupEnv("CONTAINER")
	log.Printf("runningInContainer:%v\n", runningInContainer)
	if !runningInContainer {
		err := godotenv.Load(".env.local")
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("Environment variable not found: ", key)
	}
	return value
}

func Response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		log.Println("Server Error Occured:", err)
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"data": data, "error": errorMessage})
}

func GetParam(c *gin.Context, param string) (string, error) {
	p := c.Param(param)
	if len(p) == 0 {
		return "", errors.New("invalid parameter: " + p)
	}
	return p, nil
}

// Prom
func GetRegistryMetrics() (*prometheus.Registry, *grpc_prometheus.ClientMetrics) {
	// Prom: Create a metrics registry.
	reg := prometheus.NewRegistry()
	// Prom: Create some standard client metrics.
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	// Prom: Register client metrics to registry.
	reg.MustRegister(grpcMetrics)

	return reg, grpcMetrics
}

func CreateStartPromHttpServer(reg *prometheus.Registry, port uint) {
	// Prom: Create a HTTP server for prometheus.
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", port)}
	// Prom: Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()
}
