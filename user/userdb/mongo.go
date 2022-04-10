package userdb

import (
	"context"
	"errors"

	"user/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

var _ = utils.LoadLocalEnv()
var (
	DB = utils.GetEnv("MONGO_DATABASE")
	// user = GetEnv("MONGO_USER")
	// pwd  = GetEnv("MONGO_PWD")
	Coll = utils.GetEnv("MONGO_COLLECTION")
	Addr = utils.GetEnv("MONGO_CONN")
)

var Mongo_Client *mongo.Client

func NewClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(Addr))
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
