package main

import (
	"context"
	"log"
	"net"
	"time"
	"user/userdb"
	"user/userpb"
	"user/utils"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2/bson"
)

var (
	timeout = time.Second
)

type server struct {
	userpb.UnimplementedUserServiceServer
}

// Implementations
func (*server) GetUserDetails(ctx context.Context, req *userpb.GetUserDetailsRequest) (*userpb.GetUserDetailsResponse, error) {
	go log.Println("Call GetUserDetails")
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data userdb.User
	uid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		// Non-hex id
		collection := userdb.Mongo_Client.Database(userdb.DB).Collection(userdb.Coll)
		errQuery := collection.FindOne(ctx, bson.M{"id": req.GetId()}).Decode(&data)
		if errQuery != nil {
			return nil, errQuery
		}
	} else {
		// Hex id
		collection := userdb.Mongo_Client.Database(userdb.DB).Collection(userdb.Coll)
		errQuery := collection.FindOne(ctx, bson.M{"_id": uid}).Decode(&data)
		if errQuery != nil {
			return nil, errQuery
		}
	}
	// Query: End

	return &userpb.GetUserDetailsResponse{
		User: &userpb.User{
			Id:           data.Id,
			Name:         data.Name,
			Email:        data.Email,
			IsAdmin:      int32(data.IsAdmin),
			GroupId:      data.GroupId,
			RoleId:       int32(data.RoleId),
			RefineryId:   int32(data.RefineryId),
			AreaId:       int32(data.AreaId),
			DepartmentId: int32(data.DepartmentId),
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
		},
	}, nil
}

func (*server) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	go log.Println("Call GetUsers")
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data []userdb.User
	collection := userdb.Mongo_Client.Database(userdb.DB).Collection(userdb.Coll)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var user userdb.User
		cursor.Decode(&user)
		data = append(data, user)
	}

	if err != nil {
		return nil, err
	}
	// Query: End

	if err != nil {
		return nil, utils.Error_response(err)
	}

	var users []*userpb.User
	for _, d := range data {
		users = append(
			users,
			&userpb.User{
				Id:           d.Id,
				Name:         d.Name,
				Email:        d.Email,
				IsAdmin:      int32(d.IsAdmin),
				GroupId:      d.GroupId,
				RoleId:       int32(d.RoleId),
				RefineryId:   int32(d.RefineryId),
				AreaId:       int32(d.AreaId),
				DepartmentId: int32(d.DepartmentId),
				CreatedAt:    d.CreatedAt,
				UpdatedAt:    d.UpdatedAt,
			},
		)
	}

	return &userpb.GetUsersResponse{Users: users}, nil
}

func (*server) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	go log.Println("Call login")
	_, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Query: Start
	var data userdb.User
	collection := userdb.Mongo_Client.Database(userdb.DB).Collection(userdb.Coll)

	errQuery := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&data)
	if errQuery != nil {
		return nil, utils.Error_credentials()
	}
	// Query: End

	errMatch := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if errMatch != nil {
		return nil, utils.Error_credentials()
	}

	// uid, err := primitive.ObjectIDFromHex(data.Id)
	// if err != nil {
	// 	return nil, utils.Error_response(err)
	// }

	token, errToken := generateToken(data.Id)
	if errToken != nil {
		return nil, utils.Error_response(errToken)
	}

	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:           data.Id,
			Name:         data.Name,
			Email:        data.Email,
			IsAdmin:      int32(data.IsAdmin),
			GroupId:      data.GroupId,
			RoleId:       int32(data.RoleId),
			RefineryId:   int32(data.RefineryId),
			AreaId:       int32(data.AreaId),
			DepartmentId: int32(data.DepartmentId),
			CreatedAt:    data.CreatedAt,
			UpdatedAt:    data.UpdatedAt,
		},
		Token: *token,
	}, nil
}

// Utils
var jwtKey = []byte(utils.GetEnv("APP_JWT_SECRET"))

type Claims struct {
	UserId string `json:"id"`
	jwt.StandardClaims
}

func generateToken(userId string) (*string, error) {
	claims := &Claims{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), Id: userId},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, utils.Error_response(err)
	}
	return &tokenString, nil
}

// Main
func main() {
	log.Println("User Service")

	lis, err := net.Listen("tcp", utils.GetEnv("GRPC_SERVICE_HOST")+":"+utils.GetEnv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
	defer lis.Close()

	userdb.Mongo_Client, err = userdb.NewClient(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer userdb.Mongo_Client.Disconnect(context.Background())

	s := grpc.NewServer()
	userpb.RegisterUserServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("Server started at %v", lis.Addr().String())

	err = s.Serve(lis)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}
}
