package server

import (
	"context"
	"core-service/internal/domain/user"
	"core-service/internal/usecase/service"

	"log"
	"net"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	Service *service.UserService
}

func New(service *service.UserService) *server {
	return &server{
		Service: service,
	}
}

func (s *server) CreateUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := user.New(in.Id, in.Name)
	err := s.Service.Create(user)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Resp: "ВЫ УССПЕШНО зарегались"}, nil
}

func (s *server) Listen(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)
	log.Printf("Server is listening on %s", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
