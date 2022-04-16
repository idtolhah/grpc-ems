package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func LoadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("GRPC_SERVICE_HOST"); !runningInContainer {
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
