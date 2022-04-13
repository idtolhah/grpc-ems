package masterdb

import (
	"context"
	"errors"

	"database/sql"
	"master/utils"

	_ "github.com/go-sql-driver/mysql"
)

// Models
type Area struct {
	Id   uint
	Name string
}

type AssetEquipment struct {
	ID             uint
	Item           string
	ItemCheck      string
	CheckingMethod string
	Tools          string
	StandardArea   string
	Photo          string
	LineID         uint
	MachineID      uint
}

type Contact struct {
	Id     uint
	Title  string
	Number string
	OpTime string
	OpDay  string
	Email  string
}

var _ = utils.LoadLocalEnv()
var (
	db       = utils.GetEnv("MYSQL_DB")
	username = utils.GetEnv("MYSQL_USER")
	password = utils.GetEnv("MYSQL_PASSWORD")
	host     = utils.GetEnv("MYSQL_HOST")
	port     = utils.GetEnv("MYSQL_PORT")
)

func NewClient(ctx context.Context) (*sql.DB, error) {
	url := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	client, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.New("cannot connect to mysql instance")
	}
	return client, nil
}
