package main

import (
	"github.com/micro/go-micro/v2"
	"log"
	"net/http"
)

var service = micro.NewService()

func main() {
	service.Init()
	http.HandleFunc("/", create)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}
