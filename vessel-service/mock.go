package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
)

func mock(file string, repo *Repository) error {
	// document是否有记录
	count, err := repo.collection.CountDocuments(context.Background(), bson.D{})
	if err != nil {
		return fmt.Errorf("counting document failed: %s", err)
	}
	if count != 0 {
		return nil
	}
	// 从文件读取，写入mongodb
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read vessel mock file failed: %s", err)
	}
	var vessel *Vessel
	if err := json.Unmarshal(data, &vessel); err != nil {
		return fmt.Errorf("unmarshal vessel mock file failed: %s", err)
	}
	_, err = repo.collection.InsertOne(context.Background(), vessel)
	if err != nil {
		return fmt.Errorf("insert document failed: %s", err)
	}
	return nil
}
