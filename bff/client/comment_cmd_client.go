package client

import (
	"context"
	"errors"

	"bff/pb/commentcmdpb"
	"bff/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type CreatePackingCommentRequest struct {
	EquipmentCheckingId int32  `json:"equipment_checking_id"`
	MrComment           string `json:"mr_comment"`
	MrId                string `json:"mr_id"`
}

type CommentCmdClient struct {
}

var (
	_                           = utils.LoadLocalEnv()
	commentcmdGrpcService       = utils.GetEnv("COMMENT_CMD_GRPC_SERVICE")
	commentcmdGrpcServiceClient commentcmdpb.CommentCmdServiceClient
)

func prepareCommentCmdGrpcClient(c *context.Context) error {
	// Prom: Get Registry & Metrics
	// reg, grpcMetrics := utils.GetRegistryMetrics()

	// Prom: Create a insecure gRPC channel to communicate with the server.
	conn, err := grpc.DialContext(*c, commentcmdGrpcService, []grpc.DialOption{
		grpc.WithInsecure(),
		// grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
		// grpc.WithStreamInterceptor(grpcMetrics.StreamClientInterceptor()),
		grpc.WithBlock()}...,
	)

	if err != nil {
		commentcmdGrpcServiceClient = nil
		return errors.New("connection to commentcmd gRPC service failed")
	}

	if commentcmdGrpcServiceClient != nil {
		conn.Close()
		return nil
	}

	// Prom
	// utils.CreateStartPromHttpServer(reg, 9099)

	commentcmdGrpcServiceClient = commentcmdpb.NewCommentCmdServiceClient(conn)
	return nil
}

func (ac *CommentCmdClient) CreatePackingComment(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req commentcmdpb.CreatePackingCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	if err := prepareCommentCmdGrpcClient(&ctx); err != nil {
		utils.Response(c, nil, err)
		return
	}

	data, err := commentcmdGrpcServiceClient.CreatePackingComment(
		ctx,
		&commentcmdpb.CreatePackingCommentRequest{
			EquipmentCheckingId: req.EquipmentCheckingId,
			MrComment:           req.MrComment,
			MrId:                req.MrId,
		},
	)

	if err != nil {
		utils.Response(c, nil, err)
		return
	}

	utils.Response(c, data, err)
}
