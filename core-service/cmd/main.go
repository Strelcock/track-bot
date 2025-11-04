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
	"fmt"
	"log"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	trackStatusTopic = "track.status.get"
	userMessageTopic = "user.status.push"
)

func main() {
	cfg := config.Load()
	db := postgres.New(cfg.DSN)

	fmt.Println(cfg)

	consumer := kafka.NewCons([]string{cfg.Broker}, trackStatusTopic, "A1")
	producer := kafka.NewProd([]string{cfg.Broker}, userMessageTopic)

	uService := uservice.New(db.UserDb)
	tService := tservice.New(db.TrackDb)

	consService := kservice.NewCons(consumer, db.TrackDb)
	prodService := kservice.NewProd(producer)

	q := queue.New(consService, prodService)
	go q.ServeMessages()

	conn, err := grpc.NewClient(cfg.Tracker, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println(err)
	}

	tc := pb.NewTrackerClient(conn)
	s := server.New(uService, tService, tc)

	err = s.Listen(":50051")
	if err != nil {
		log.Fatalf("failed to serve: %s", err.Error())
	}
}
