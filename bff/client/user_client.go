package client

import (
	"context"
	"errors"

	"bff/pb/userpb"

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
	_                     = loadLocalEnv()
	userGrpcService       = GetEnv("USER_GRPC_SERVICE")
	userGrpcServiceClient userpb.UserServiceClient
)

func prepareUserGrpcClient(c *context.Context) error {

	conn, err := grpc.DialContext(*c, userGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock()}...)

	if err != nil {
		userGrpcServiceClient = nil
		return errors.New("connection to user gRPC service failed")
	}

	if userGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	userGrpcServiceClient = userpb.NewUserServiceClient(conn)
	return nil
}

func (uc *UserClient) GetUserDetails(id string, c *context.Context) (*User, error) {

	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.GetUserDetails(*c, &userpb.GetUserDetailsRequest{Id: id})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}
	return &User{
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
	}, nil
}

func (uc *UserClient) GetUsers(c *context.Context) (*[]User, error) {

	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.GetUsers(*c, &userpb.GetUsersRequest{})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
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
	return &users, nil
}

func (uc *UserClient) Login(username string, password string, c *context.Context) (*UserWithToken, error) {

	if err := prepareUserGrpcClient(c); err != nil {
		return nil, err
	}

	res, err := userGrpcServiceClient.Login(*c, &userpb.LoginRequest{Email: username, Password: password})
	if err != nil {
		return nil, errors.New(status.Convert(err).Message())
	}

	return &UserWithToken{
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
	}, nil
}
