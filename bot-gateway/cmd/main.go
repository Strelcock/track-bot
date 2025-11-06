package main

import (
	"bot-gateway/api/bot"
	"bot-gateway/api/queue"
	"bot-gateway/config"
	"bot-gateway/internal/infrastructure/kafka"
	"bot-gateway/internal/usecase/kafkaservice"
	"log"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userMessageTopic = "user.status.push"
)

func main() {
	cfg := config.Load()

	//fmt.Println(cfg)
	//new grpc client
	conn, err := grpc.NewClient(cfg.Core, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	uc := pb.NewUserServiceClient(conn)
	tc := pb.NewTrackServiceClient(conn)

	//new bot
	bot, err := bot.New(cfg.Token, uc, tc)
	if err != nil {
		log.Fatal(err)
	}
	cons := kafka.New([]string{cfg.Broker}, userMessageTopic, "B1")
	kafkaService := kafkaservice.New(cons)
	q := queue.New(kafkaService)

	go q.Read(bot)

	bot.Start()

}
