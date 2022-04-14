package clients

import (
	"context"
	"packing/pb/masterpb"
	"packing/pb/packingquerypb"
	"strconv"
	"time"
)

var timeout = time.Second

func GetUnit(id int32) *packingquerypb.Unit {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Unit{}
	}

	res, err := masterGrpcServiceClient.GetUnit(c, &masterpb.GetUnitRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Unit{}
	}

	return &packingquerypb.Unit{Id: res.Unit.Id, Name: res.Unit.Name}
}

func GetLine(id int32) *packingquerypb.Line {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Line{}
	}

	res, err := masterGrpcServiceClient.GetLine(c, &masterpb.GetLineRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Line{}
	}

	return &packingquerypb.Line{Id: res.Line.Id, Name: res.Line.Name}
}

func GetMachine(id int32) *packingquerypb.Machine {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareUnitGrpcClient(&ctx); err != nil {
		return &packingquerypb.Machine{}
	}

	res, err := masterGrpcServiceClient.GetMachine(c, &masterpb.GetMachineRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Machine{}
	}

	return &packingquerypb.Machine{Id: res.Machine.Id, Name: res.Machine.Name}
}
