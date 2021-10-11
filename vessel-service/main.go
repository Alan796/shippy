package main

import (
	"context"
	pb "github.com/Alan796/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro/v2"
	"log"
	"os"
)

const (
	defaultHost    = "localhost:27017"
	dbName         = "shippy"
	collectionName = "vessels"
	vesselMockFile = "vessel.json"
)

func main() {
	// 创建一个微服务并初始化
	service := micro.NewService(micro.Name("vessel"))
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

	collection := mongoClient.Database(dbName).Collection(collectionName)
	repo := &Repository{collection: collection}
	// 测试环境mock vessel数据
	if err := mock(vesselMockFile, repo); err != nil {
		log.Panic(err)
	}
	h := &handler{repo: repo}

	// 注册handler
	if err := pb.RegisterVesselServiceHandler(service.Server(), h); err != nil {
		log.Panic(err)
	}

	// 运行微服务
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
