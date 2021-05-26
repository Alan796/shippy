package main

import (
	"context"
	pb "github.com/Alan796/shippy/consignment-service/proto/consignment"
	vesselPb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
	"os"
)

const (
	defaultHost    = "localhost:27017"
	dbName         = "shippy"
	collectionName = "consignments"
)

func main() {
	// 创建一个微服务并初始化
	service := micro.NewService(micro.Name("consignment"))
	service.Init()

	// mongodb uri
	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	// 创建mongo客户端连接
	mongoClient, err := CreateMongoClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer mongoClient.Disconnect(context.Background())

	// 创建handler，包括存储和vessel服务的客户端
	collection := mongoClient.Database(dbName).Collection(collectionName)
	repo := &Repository{collection: collection}
	vesselClient := vesselPb.NewVesselService("vessel", service.Client())
	h := &handler{repo: repo, vesselService: vesselClient}

	// 注册handler
	if err := pb.RegisterConsignmentServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	// 运行微服务
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
