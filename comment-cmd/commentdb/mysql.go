package commentdb

import (
	"context"
	"errors"

	"comment/utils"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

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
