package main

import (
	"core-service/api/server"
	"log"
)

func main() {
	err := server.New(":50051")
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
