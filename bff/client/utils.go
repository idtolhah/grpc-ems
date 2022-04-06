package client

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("CONTAINER"); !runningInContainer {
		err := godotenv.Load(".env")
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
