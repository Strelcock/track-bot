package server

import (
	"context"
	"log"
	"net"
	"tracker/api/hook"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTrackerServer
	*hook.Sender
}

func New(sender *hook.Sender) *server {
	return &server{Sender: sender}
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

	carrier, err := s.Carrier(in.Number)
	if err != nil {
		return nil, NewAddError(in.Number, err)
	}

	err = s.AddTracker(carrier, in.Number)
	if err != nil {
		return nil, NewAddError(in.Number, err)
	}

	return nil, nil
}
