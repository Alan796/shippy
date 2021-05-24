package main

import (
	"context"
	pb "github.com/Alan796/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository 模拟一个数据库
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type consignmentService struct {
	repo          repository
	vesselService vesselPb.VesselService
}

func (s *consignmentService) Create(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	vesselResponse, err := s.vesselService.FindAvailable(context.Background(), &vesselPb.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	req.VesselId = vesselResponse.Vessel.Id

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *consignmentService) GetAll(ctx context.Context, req *pb.GetAllRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	srv := micro.NewService(micro.Name("consignment"))
	srv.Init()

	repo := &Repository{}
	vesselClient := vesselPb.NewVesselService("vessel", srv.Client())

	if err := pb.RegisterConsignmentServiceHandler(srv.Server(), &consignmentService{repo: repo, vesselService: vesselClient}); err != nil {
		log.Panic(err)
	}
	if err := srv.Run(); err != nil {
		log.Panic(err)
	}
}
