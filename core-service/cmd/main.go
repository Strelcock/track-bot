package main

import (
	"core-service/api/server"
	"core-service/config"
	"core-service/internal/infrastructure/postgres"
	"core-service/internal/usecase/service"
	"log"
)

func main() {
	cfg := config.Load()
	db := postgres.New(cfg.DSN)

	service := service.New(db)

	s := server.New(service)

	err := s.Listen(":50051")
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
