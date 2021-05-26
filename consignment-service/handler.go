package main

import (
	"context"
	pb "github.com/Alan796/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"log"
)

type handler struct {
	repo          repository
	vesselService vesselPb.VesselService
}

// Create 创建consignment
func (h *handler) Create(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := h.vesselService.FindAvailable(ctx, &vesselPb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	if err := h.repo.Create(ctx, req); err != nil {
		return err
	}

	res.Created = true
	res.Consignment = req
	return nil
}

// GetAll 获取所有consignment
func (h *handler) GetAll(ctx context.Context, req *pb.GetAllRequest, res *pb.Response) error {
	consignments, err := h.repo.GetAll(ctx)
	if err != nil {
		return err
	}
	res.Consignments = consignments
	return nil
}
