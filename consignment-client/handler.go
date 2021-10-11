package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Alan796/shippy/consignment-service/proto/consignment"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	file = "consignment.json"
)

var client = pb.NewConsignmentService("consignment", service.Client())

func create(w http.ResponseWriter, _ *http.Request) {
	consignment, err := parseFile(file)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not parse file: %v", err)))
		return
	}
	resp, err := client.Create(context.Background(), consignment)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not greet: %v", err)))
		return
	}
	log.Printf("Created: %t", resp.Created)
	getAll, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Could not list consignments: %v", err)))
		return
	}
	consignments, err := json.Marshal(getAll.Consignments)
	if err != nil {
		log.Fatalln(err)
	}
	w.Write(consignments)
}

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &consignment); err != nil {
		return nil, err
	}
	return consignment, nil
}
