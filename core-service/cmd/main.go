package main

import (
	"core-service/api/queue"
	"core-service/api/server"
	"core-service/config"
	"core-service/internal/infrastructure/kafka"
	"core-service/internal/infrastructure/postgres"
	"core-service/internal/usecase/kservice"
	"core-service/internal/usecase/tservice"
	"core-service/internal/usecase/uservice"
	"log"
)

const (
	trackStatusTopic = "track.status.get"
	userMessageTopic = "user.status.push"
	broker           = "localhost:9092"
)

func main() {
	cfg := config.Load()
	db := postgres.New(cfg.DSN)

	consumer := kafka.NewCons([]string{broker}, trackStatusTopic, "A1")
	producer := kafka.NewProd([]string{broker}, userMessageTopic)

	uService := uservice.New(db.UserDb)
	tService := tservice.New(db.TrackDb)

	consService := kservice.NewCons(consumer, db.TrackDb)
	prodService := kservice.NewProd(producer)

	q := queue.New(consService, prodService)
	go q.ServeMessages()

	s := server.New(uService, tService)

	err := s.Listen(":50051")
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
