package userdb

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name,omitempty"`
	Email        string             `bson:"email,omitempty"`
	Password     string             `bson:"password,omitempty"`
	IsAdmin      uint               `bson:"is_admin,omitempty"`
	GroupId      string             `bson:"group_id,omitempty"`
	RoleId       uint               `bson:"role_id,omitempty"`
	RefineryId   uint               `bson:"refinery_id,omitempty"`
	AreaId       uint               `bson:"area_id,omitempty"`
	DepartmentId uint               `bson:"department_id,omitempty"`
	CreatedAt    string             `bson:"createdAt,omitempty"`
	UpdatedAt    string             `bson:"updatedAt,omitempty"`
}

var _ = loadLocalEnv()
var (
	db = GetEnv("MONGO_DATABASE")
	// user = GetEnv("MONGO_USER")
	// pwd  = GetEnv("MONGO_PWD")
	coll = GetEnv("MONGO_COLLECTION")
	addr = GetEnv("MONGO_CONN")
)

var Mongo_Client *mongo.Client

func NewClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(addr))
	// .SetAuth(options.Credential{
	// 	AuthSource: db,
	// 	Username:   user,
	// 	Password:   pwd,
	// }))
	if err != nil {
		return nil, errors.New("invalid mongodb options")
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("cannot connect to mongodb instance")
	}
	return client, nil
}

func FindOne(ctx context.Context, condition bson.M) (*User, error) {

	collection := Mongo_Client.Database(db).Collection(coll)

	var data User
	err := collection.FindOne(ctx, condition).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func Find(ctx context.Context) (*[]User, error) {

	collection := Mongo_Client.Database(db).Collection(coll)

	var data []User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var user User
		cursor.Decode(&user)
		data = append(data, user)
	}
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func loadLocalEnv() interface{} {
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
