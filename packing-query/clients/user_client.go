package clients

import (
	"context"
	"packing/pb/packingquerypb"
	"packing/pb/userpb"
)

func GetUser(id string) *packingquerypb.User {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUserGrpcClient(&ctx); err != nil {
		return &packingquerypb.User{}
	}

	res, err := userGrpcServiceClient.GetUserDetails(c, &userpb.GetUserDetailsRequest{Id: id})
	if err != nil {
		return &packingquerypb.User{}
	}

	return &packingquerypb.User{Id: res.User.Id, Name: res.User.Name, GroupId: res.User.GroupId}
}
