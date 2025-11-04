package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTrackerServer
}

func New() *server {
	return &server{}
}

func (s *server) Listen(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTrackerServer(grpcServer, s)
	log.Printf("Server is listening on %s", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *server) ServeTrack(ctx context.Context, in *pb.ToTracker) (*pb.Empty, error) {
	fmt.Println(in.Number)
	return nil, nil
}
