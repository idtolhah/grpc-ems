package clients

import (
	"context"
	"filtration-query/pb/filtrationquerypb"
	"filtration-query/pb/masterpb"
	"strconv"
	"time"
)

var timeout = time.Second

func GetUnit(id int32) *filtrationquerypb.Unit {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &filtrationquerypb.Unit{}
	}

	res, err := masterGrpcServiceClient.GetUnit(c, &masterpb.GetUnitRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &filtrationquerypb.Unit{}
	}

	return &filtrationquerypb.Unit{Name: res.Unit.Name}
}

func GetFilter(id int32) *filtrationquerypb.Filter {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &filtrationquerypb.Filter{}
	}

	res, err := masterGrpcServiceClient.GetFilter(c, &masterpb.GetFilterRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &filtrationquerypb.Filter{}
	}

	return &filtrationquerypb.Filter{Code: res.Filter.Code}
}

func GetDirt(id int32) *filtrationquerypb.Dirt {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &filtrationquerypb.Dirt{}
	}

	res, err := masterGrpcServiceClient.GetDirt(c, &masterpb.GetDirtRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &filtrationquerypb.Dirt{}
	}

	return &filtrationquerypb.Dirt{Name: res.Dirt.Name}
}
