package main

import (
	"context"
	pb "github.com/Alan796/shippy/vessel-service/proto/vessel"
)

type handler struct {
	repo repository
}

// FindAvailable 查找符合要求的vessel
func (h *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := h.repo.FindAvailable(ctx, req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}

// Create 创建vessel
func (h *handler) Create(ctx context.Context, req *pb.Vessel, res *pb.Response) error {
	if err := h.repo.Create(ctx, req); err != nil {
		return err
	}

	res.Created = true
	res.Vessel = req
	return nil
}
