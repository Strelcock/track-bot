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
	var errs AddError
	for _, n := range in.Number {
		carrier, err := s.Carrier(n)
		if err != nil {
			errs.Errs = append(errs.Errs, err.Error())
			continue
		}
		err = s.AddTracker(carrier, n)
		if err != nil {
			errs.Errs = append(errs.Errs, err.Error())
		}
	}

	if errs.Errs != nil {
		return nil, &errs
	}

	return nil, nil
}
