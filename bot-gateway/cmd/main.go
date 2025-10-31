package main

import (
	"bot-gateway/api/bot"
	"bot-gateway/api/queue"
	"bot-gateway/internal/infrastructure/kafka"
	"bot-gateway/internal/usecase/kafkaservice"
	"log"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50051"

const (
	userMessageTopic = "user.status.push"
	broker           = "localhost:9092"
)

func main() {

	//new grpc client
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	uc := pb.NewUserServiceClient(conn)
	tc := pb.NewTrackServiceClient(conn)

	//new bot
	bot, err := bot.New("8286937197:AAFrfcaG_g_s1Sw5YZKUVgbtxyWbC9M8LWc", uc, tc)
	if err != nil {
		log.Fatal(err)
	}
	cons := kafka.New([]string{broker}, userMessageTopic, "B1")
	kafkaService := kafkaservice.New(cons)
	q := queue.New(kafkaService)

	go q.Read(bot)

	bot.Start()

}
