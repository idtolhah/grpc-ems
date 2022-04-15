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

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &packingquerypb.Unit{}
	}

	res, err := masterGrpcServiceClient.GetUnit(c, &masterpb.GetUnitRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Unit{}
	}

	return &packingquerypb.Unit{Id: res.Unit.Id, Name: res.Unit.Name}
}

func GetDepartment(id int32) *packingquerypb.Department {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &packingquerypb.Department{}
	}

	res, err := masterGrpcServiceClient.GetDepartment(c, &masterpb.GetDepartmentRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Department{}
	}

	return &packingquerypb.Department{Id: res.Department.Id, Name: res.Department.Name}
}

func GetArea(id int32) *packingquerypb.Area {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &packingquerypb.Area{}
	}

	res, err := masterGrpcServiceClient.GetArea(c, &masterpb.GetAreaRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Area{}
	}

	return &packingquerypb.Area{Id: res.Area.Id, Name: res.Area.Name}
}

func GetLine(id int32) *packingquerypb.Line {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
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

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &packingquerypb.Machine{}
	}

	res, err := masterGrpcServiceClient.GetMachine(c, &masterpb.GetMachineRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.Machine{}
	}

	return &packingquerypb.Machine{Id: res.Machine.Id, Name: res.Machine.Name}
}

func GetAssetEquipment(id int32) *packingquerypb.AssetEquipment {
	c := context.Background()
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	if err := prepareMasterGrpcClient(&ctx); err != nil {
		return &packingquerypb.AssetEquipment{}
	}

	res, err := masterGrpcServiceClient.GetAssetEquipment(c, &masterpb.GetAssetEquipmentRequest{Id: strconv.Itoa(int(id))})
	if err != nil {
		return &packingquerypb.AssetEquipment{}
	}

	return &packingquerypb.AssetEquipment{
		Id: res.Assetequipment.Id, Item: res.Assetequipment.Item, ItemCheck: res.Assetequipment.ItemCheck,
		CheckingMethod: res.Assetequipment.CheckingMethod, Tools: res.Assetequipment.Tools, StandardArea: res.Assetequipment.StandardArea,
		Photo: res.Assetequipment.Photo, LineId: res.Assetequipment.LineId, MachineId: res.Assetequipment.MachineId,
	}
}
