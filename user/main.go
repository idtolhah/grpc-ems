package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"user/userdb"
	"user/userpb"

	"github.com/golang-jwt/jwt"
	consulapi "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	log.Println("Called GetUserDetails, Id", req.Id)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, error_response(err)
	}

	res, err := userdb.FindOne(c, bson.M{"_id": uid})
	if err != nil {
		return nil, error_response(err)
	}

	return &userpb.GetUserDetailsResponse{
		User: &userpb.User{
			Id:           res.Id.Hex(),
			Name:         res.Name,
			Email:        res.Email,
			IsAdmin:      int32(res.IsAdmin),
			GroupId:      res.GroupId,
			RoleId:       int32(res.RoleId),
			RefineryId:   int32(res.RefineryId),
			AreaId:       int32(res.AreaId),
			DepartmentId: int32(res.DepartmentId),
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,
		},
	}, nil
}

func (*server) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	log.Println("Called GetUsers")

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	data, err := userdb.Find(c)

	if err != nil {
		return nil, error_response(err)
	}

	var users []*userpb.User
	for _, d := range *data {
		users = append(
			users,
			&userpb.User{
				Id:           d.Id.Hex(),
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
	log.Printf("Called Login, Email: %v", req.Email)

	c, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	res, err := userdb.FindOne(c, bson.M{"email": req.Email})
	if err != nil {
		return nil, error_credentials()
	}

	errMatch := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(req.Password))
	if errMatch != nil {
		return nil, error_credentials()
	}

	uid, err := primitive.ObjectIDFromHex(res.Id.Hex())
	if err != nil {
		return nil, error_response(err)
	}

	token, errToken := generateToken(uid.Hex())
	if errToken != nil {
		return nil, error_response(errToken)
	}

	return &userpb.LoginResponse{
		User: &userpb.User{
			Id:           res.Id.Hex(),
			Name:         res.Name,
			Email:        res.Email,
			IsAdmin:      int32(res.IsAdmin),
			GroupId:      res.GroupId,
			RoleId:       int32(res.RoleId),
			RefineryId:   int32(res.RefineryId),
			AreaId:       int32(res.AreaId),
			DepartmentId: int32(res.DepartmentId),
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,
		},
		Token: *token,
	}, nil
}

// Utils
var jwtKey = []byte(userdb.GetEnv("APP_JWT_SECRET"))

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
		return nil, error_response(err)
	}
	return &tokenString, nil
}

func error_response(err error) error {
	log.Println("ERROR:", err.Error())
	return status.Error(codes.Internal, err.Error())
}

func error_credentials() error {
	return status.Error(codes.Internal, "Invalid Credentials!")
}

// Consul
func registerServiceWithConsul() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "user-service"   //replace with service id
	registration.Name = "user-service" //replace with service name
	address := hostname()
	registration.Address = address
	port, err := strconv.Atoi(port()[1:len(port())])
	if err != nil {
		log.Fatalln(err)
	}
	registration.Port = port
	registration.Check = new(consulapi.AgentServiceCheck)
	registration.Check.HTTP = fmt.Sprintf("http://%s:%v/healthcheck",
		address, port)
	registration.Check.Interval = "5s"
	registration.Check.Timeout = "3s"
	consul.Agent().ServiceRegister(registration)
}

func hostname() string {
	hn, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	return hn
}

func port() string {
	p := os.Getenv("PRODUCT_SERVICE_PORT")
	if len(strings.TrimSpace(p)) == 0 {
		return ":8100"
	}
	return fmt.Sprintf(":%s", p)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `product service is good`)
}

// Main
func main() {
	log.Println("User Service")
	registerServiceWithConsul()

	lis, err := net.Listen("tcp", userdb.GetEnv("GRPC_SERVICE_HOST")+":"+userdb.GetEnv("GRPC_SERVICE_PORT"))
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
