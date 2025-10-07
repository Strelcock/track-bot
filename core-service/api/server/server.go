package server

import (
	"context"

	"log"
	"net"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {

	return &pb.UserResponse{Resp: "ВЫ УССПЕШНО зарегались"}, nil
}

func New(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})
	log.Printf("Server is listening on %s", lis.Addr())
	if err = s.Serve(lis); err != nil {
		return err
	}
	return nil
}
