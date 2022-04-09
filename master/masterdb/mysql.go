package masterdb

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

var _ = LoadLocalEnv()
var (
	db       = GetEnv("MYSQL_DB")
	username = GetEnv("MYSQL_USER")
	password = GetEnv("MYSQL_PASSWORD")
	host     = GetEnv("MYSQL_HOST")
	port     = GetEnv("MYSQL_PORT")
)

func NewClient(ctx context.Context) (*sql.DB, error) {
	// url := username + ":" + password + "@" + host + "/" + db
	// url := username + "@/" + db
	url := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + db
	client, err := sql.Open("mysql", url)
	if err != nil {
		return nil, errors.New("cannot connect to mysql instance")
	}
	return client, nil
}

func FindAreas(client *sql.DB, ctx context.Context) (*[]Area, error) {
	var data []Area
	results, err := client.Query("SELECT id, name FROM areas")
	if err != nil {
		return nil, err
	}
	var area Area
	for results.Next() {
		err = results.Scan(&area.Id, &area.Name)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, area)
	}
	return &data, nil
}

func FindAssetEquipments(client *sql.DB, ctx context.Context) (*[]AssetEquipment, error) {
	var data []AssetEquipment
	results, err := client.Query("SELECT id, item, item_check, checking_method, tools, standard_area, photo, line_id, machine_id FROM asset_equipments")
	if err != nil {
		return nil, err
	}
	var a AssetEquipment
	for results.Next() {
		err = results.Scan(&a.ID, &a.Item, &a.ItemCheck, &a.CheckingMethod, &a.Tools, &a.StandardArea, &a.Photo, &a.LineID, &a.MachineID)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, a)
	}
	return &data, nil
}

func FindContacts(client *sql.DB, ctx context.Context) (*[]Contact, error) {
	var data []Contact
	results, err := client.Query("SELECT * FROM contacts")
	if err != nil {
		return nil, err
	}
	var contact Contact
	for results.Next() {
		err = results.Scan(&contact.Id, &contact.Title, &contact.Number, &contact.OpTime, &contact.OpDay, &contact.Email)
		if err != nil {
			panic(err.Error())
		}
		data = append(data, contact)
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
