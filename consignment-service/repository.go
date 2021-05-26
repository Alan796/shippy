package main

import (
	"context"
	pb "github.com/Alan796/shippy/consignment-service/proto/consignment"
	"go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
	Create(context.Context, *pb.Consignment) error
	GetAll(context.Context) ([]*pb.Consignment, error)
}

// Repository 存储库
type Repository struct {
	collection *mongo.Collection
}

// Create 创建consignment
func (repo *Repository) Create(ctx context.Context, consignment *pb.Consignment) error {
	_, err := repo.collection.InsertOne(ctx, MarshalConsignment(consignment))
	return err
}

// GetAll 获取所有consignment
func (repo *Repository) GetAll(ctx context.Context) ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment
	cursor, err := repo.collection.Find(ctx, nil, nil)
	for cursor.Next(ctx) {
		var consignment *Consignment
		if err := cursor.Decode(&consignment); err != nil { // todo &
			return nil, err
		}
		consignments = append(consignments, UnmarshalConsignment(consignment))
	}
	return consignments, err
}

// Consignment mongodb model
type Consignment struct {
	Id          string       `json:"id"`
	Weight      int32        `json:"weight"`
	Description string       `json:"description"`
	Containers  []*Container `json:"containers"`
	VesselId    string       `json:"vessel_id"`
}

// Container mongodb model
type Container struct {
	Id         string `json:"id"`
	CustomerId string `json:"customer_id"`
	UserId     string `json:"user_id"`
}

// MarshalConsignment Consignment由pb model转为mongodb model
func MarshalConsignment(consignment *pb.Consignment) *Consignment {
	containers := MarshalContainers(consignment.Containers)
	return &Consignment{
		Id:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  containers,
		VesselId:    consignment.VesselId,
	}
}

// MarshalContainers Containers由pb model转为mongodb model
func MarshalContainers(containers []*pb.Container) []*Container {
	var collection []*Container
	for _, container := range containers {
		collection = append(collection, MarshalContainer(container))
	}
	return collection
}

// MarshalContainer Container由pb model转为mongodb model
func MarshalContainer(container *pb.Container) *Container {
	return &Container{
		Id:         container.Id,
		CustomerId: container.CustomerId,
		UserId:     container.UserId,
	}
}

// UnmarshalConsignment Consignment由mongodb model转为pb model
func UnmarshalConsignment(consignment *Consignment) *pb.Consignment {
	containers := UnmarshalContainers(consignment.Containers)
	return &pb.Consignment{
		Id:          consignment.Id,
		Weight:      consignment.Weight,
		Description: consignment.Description,
		Containers:  containers,
		VesselId:    consignment.VesselId,
	}
}

// UnmarshalContainers Containers由mongodb model转为pb model
func UnmarshalContainers(containers []*Container) []*pb.Container {
	var collection []*pb.Container
	for _, container := range containers {
		collection = append(collection, UnmarshalContainer(container))
	}
	return collection
}

// UnmarshalContainer Container由mongodb model转为pb model
func UnmarshalContainer(container *Container) *pb.Container {
	return &pb.Container{
		Id:         container.Id,
		CustomerId: container.CustomerId,
		UserId:     container.UserId,
	}
}
