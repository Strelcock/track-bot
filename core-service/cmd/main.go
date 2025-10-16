package main

import (
	"core-service/api/server"
	"core-service/config"
	"core-service/internal/infrastructure/postgres"
	"core-service/internal/usecase/uservice"
	"log"
)

func main() {
	cfg := config.Load()
	db := postgres.New(cfg.DSN)

	uService := uservice.New(db.UserDb)

	s := server.New(uService)

	err := s.Listen(":50051")
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
