package main

import (
	"bot-gateway/api/trackbot"
	"log"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50051"

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
	bot, err := trackbot.New("8286937197:AAFrfcaG_g_s1Sw5YZKUVgbtxyWbC9M8LWc", uc, tc)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()

}
