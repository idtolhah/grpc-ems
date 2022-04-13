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

func Error_credentials() error {
	return status.Error(codes.Internal, "Invalid Credentials!")
}

func LoadLocalEnv() interface{} {
	if _, runningInContainer := os.LookupEnv("CONTAINER"); !runningInContainer {
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
