package packingdb

import (
	"context"
	"errors"
	"log"
	"os"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Models
type Packing struct {
	Id           int64
	FoId         string
	LineId       int32
	MachineId    int32
	UnitId       int32
	DepartmentId int32
	AreaId       int32
	CompletedAt  string
	Status       int32
	CreatedAt    string
	UpdatedAt    string
}

var _ = LoadLocalEnv()
var (
	db       = GetEnv("MYSQL_DB")
	username = GetEnv("MYSQL_USER")
	password = GetEnv("MYSQL_PASSWORD")
	host     = GetEnv("MYSQL_HOST")
	port     = GetEnv("MYSQL_PORT")
)

func NewClient(ctx context.Context) (*sql.DB, error) {
	url := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	client, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.New("cannot connect to mysql instance")
	}
	return client, nil
}

func FindPackings(client *sql.DB, ctx context.Context) (*[]Packing, error) {
	var data []Packing
	results, err := client.Query("SELECT * FROM packings")
	if err != nil {
		return nil, err
	}
	var packing Packing
	for results.Next() {
		err = results.Scan(
			&packing.Id, &packing.FoId, &packing.LineId, &packing.MachineId, &packing.UnitId, &packing.DepartmentId,
			&packing.AreaId, &packing.CompletedAt, &packing.Status, &packing.CreatedAt, &packing.UpdatedAt,
		)
		if err != nil {
			// panic(err.Error())
			log.Println(err)
		}
		data = append(data, packing)
	}
	return &data, nil
}

// Utils
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
