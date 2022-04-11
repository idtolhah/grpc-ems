package client

import (
	"context"
	"errors"

	"bff/pb/userpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type User struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsAdmin      int32  `json:"is_admin"`
	GroupId      string `json:"group_id"`
	RoleId       uint   `json:"role_id"`
	RefineryId   uint   `json:"refinery_id"`
	AreaId       uint   `json:"area_id"`
	DepartmentId uint   `json:"department_id"`
	CreatedAt    string `json:"createdAt"`
	// UpdatedAt    string `json:"updatedAt"`
}

type UserWithToken struct {
	User
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserClient struct {
}

var (
	_                     = utils.LoadLocalEnv()
	userGrpcService       = utils.GetEnv("USER_GRPC_SERVICE")
	userGrpcServiceClient userpb.UserServiceClient
)

func prepareUserGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	reg, grpcMetrics := utils.GetRegistryMetrics()
	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, userGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		userGrpcServiceClient = nil
		return errors.New("connection to user gRPC service failed")
	}

	if userGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	utils.CreateStartPromHttpServer(reg, 9095)

	userGrpcServiceClient = userpb.NewUserServiceClient(conn)
	return nil
}

func (uc *UserClient) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := prepareUserGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := userGrpcServiceClient.Login(ctx, &userpb.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
	}

	utils.Response(c, &UserWithToken{
		User{Id: res.User.Id,
			Name:         res.User.Name,
			Email:        res.User.Email,
			GroupId:      res.User.GroupId,
			RoleId:       uint(res.User.RoleId),
			RefineryId:   uint(res.User.RefineryId),
			AreaId:       uint(res.User.AreaId),
			DepartmentId: uint(res.User.DepartmentId),
			CreatedAt:    res.User.CreatedAt,
		}, res.Token,
	}, nil)
}

func (uc *UserClient) GetUserDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	// Get UserId from gin context after jwt auth
	userId, ok := c.Get("UserId")
	if !ok {
		utils.Response(c, nil, errors.New("invalid user id in token"))
		return
	}

	if err := prepareUserGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := userGrpcServiceClient.GetUserDetails(ctx, &userpb.GetUserDetailsRequest{Id: userId.(string)})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	utils.Response(c, &User{
		Id:           res.User.Id,
		Name:         res.User.Name,
		Email:        res.User.Email,
		IsAdmin:      res.User.IsAdmin,
		GroupId:      res.User.GroupId,
		RoleId:       uint(res.User.RoleId),
		RefineryId:   uint(res.User.RefineryId),
		AreaId:       uint(res.User.AreaId),
		DepartmentId: uint(res.User.DepartmentId),
		CreatedAt:    res.User.CreatedAt,
	}, err)
}

func (uc *UserClient) GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUserGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	res, err := userGrpcServiceClient.GetUsers(ctx, &userpb.GetUsersRequest{})
	if err != nil {
		utils.Response(c, nil, errors.New(status.Convert(err).Message()))
		return
	}

	var users []User
	for _, u := range res.GetUsers() {
		users = append(users, User{
			Id:           u.Id,
			Name:         u.Name,
			Email:        u.Email,
			GroupId:      u.GroupId,
			RoleId:       uint(u.RoleId),
			RefineryId:   uint(u.RefineryId),
			AreaId:       uint(u.AreaId),
			DepartmentId: uint(u.DepartmentId),
			CreatedAt:    u.CreatedAt,
		})
	}
	utils.Response(c, &users, nil)
}
