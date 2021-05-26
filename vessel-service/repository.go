package main

import (
	"context"
	pb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	FindAvailable(context.Context, *pb.Specification) (*pb.Vessel, error)
	Create(context.Context, *pb.Vessel) error
}

// Repository 存储库
type Repository struct {
	collection *mongo.Collection
}

// FindAvailable 查找符合要求的vessel
func (repo *Repository) FindAvailable(ctx context.Context, spec *pb.Specification) (*pb.Vessel, error) {
	marshaledSpec := MarshalSpecification(spec)
	vessel := &Vessel{}
	filter := bson.D{
		{
			"capacity",
			bson.D{{
				"$gte",
				marshaledSpec.Capacity,
			}},
		},
		{
			"max_weight",
			bson.D{{
				"$gte",
				marshaledSpec.MaxWeight,
			}},
		},
	}
	if err := repo.collection.FindOne(ctx, filter).Decode(vessel); err != nil {
		return nil, err
	}
	return UnmarshalVessel(vessel), nil
}

// Create 创建vessel
func (repo *Repository) Create(ctx context.Context, vessel *pb.Vessel) error {
	_, err := repo.collection.InsertOne(ctx, MarshalVessel(vessel))
	return err
}

// Specification mongodb model
type Specification struct {
	Capacity  int32
	MaxWeight int32
}

// Vessel mongodb model
type Vessel struct {
	Id        string // todo `json`
	Capacity  int32
	Name      string
	Available bool
	OwnerId   string
	MaxWeight int32
}

// MarshalSpecification Specification由pb model转为mongodb model
func MarshalSpecification(spec *pb.Specification) *Specification {
	return &Specification{
		Capacity:  spec.Capacity,
		MaxWeight: spec.MaxWeight,
	}
}

// MarshalVessel Vessel由pb model转为mongodb model
func MarshalVessel(vessel *pb.Vessel) *Vessel {
	return &Vessel{
		Id:        vessel.Id,
		Capacity:  vessel.Capacity,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerId:   vessel.OwnerId,
		MaxWeight: vessel.MaxWeight,
	}
}

// UnmarshalVessel Vessel由mongodb model转为pb model
func UnmarshalVessel(vessel *Vessel) *pb.Vessel {
	return &pb.Vessel{
		Id:        vessel.Id,
		Capacity:  vessel.Capacity,
		Name:      vessel.Name,
		Available: vessel.Available,
		OwnerId:   vessel.OwnerId,
		MaxWeight: vessel.MaxWeight,
	}
}
