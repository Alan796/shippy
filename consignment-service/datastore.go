package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var maxRetry int32 = 3
var retryGap = 2 * time.Second

// CreateMongoClient 创建mongo客户端连接
func CreateMongoClient(ctx context.Context, uri string, retry int32) (*mongo.Client, error) {
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err := conn.Ping(ctx, nil); err != nil {
		if retry >= maxRetry {
			return nil, err
		}
		retry = retry + 1
		time.Sleep(retryGap)
		return CreateMongoClient(ctx, uri, retry)
	}
	return conn, err
}
