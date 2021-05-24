package main

import (
	"context"
	"errors"
	pb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
)

type repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type Repository struct {
	vessels []*pb.Vessel
}

func (repo *Repository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("no vessel found by that spec")
}

type service struct {
	repo repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Vessel = vessel
	return nil
}

func main() {
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Titanic", MaxWeight: 200000, Capacity: 500},
		{Id: "vessel002", Name: "Pearl", MaxWeight: 100000, Capacity: 300},
	}
	repo := &Repository{vessels: vessels}

	srv := micro.NewService(micro.Name("vessel"))
	srv.Init()

	if err := pb.RegisterVesselServiceHandler(srv.Server(), &service{repo: repo}); err != nil {
		log.Panic(err)
	}
	if err := srv.Run(); err != nil {
		log.Panic(err)
	}
}
