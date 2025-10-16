package server

import (
	"context"
	"core-service/internal/domain/user"
	"core-service/internal/usecase/uservice"
	"strings"

	"log"
	"net"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	*userServer
	*trackServer
}

type trackServer struct {
	pb.UnimplementedTrackServiceServer
}

type userServer struct {
	pb.UnimplementedUserServiceServer
	UserService *uservice.UserService
}

func New(uService *uservice.UserService) *server {
	return &server{
		&userServer{
			UserService: uService,
		},
		&trackServer{},
	}
}

func (s *server) CreateUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	user := user.New(in.Id, in.Name, false)
	err := s.UserService.Create(user)
	if err != nil {
		return nil, err
	}
	return &pb.UserResponse{Resp: "ВЫ УССПЕШНО зарегались"}, nil
}

func (s *server) IsAdmin(ctx context.Context, in *pb.AdminRequest) (*pb.AdminResponse, error) {
	admin, err := s.UserService.IsAdmin(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AdminResponse{IsAdmin: admin}, nil
}

func (s *server) AddTrack(ctx context.Context, in *pb.TrackRequest) (*pb.TrackResponse, error) {
	resp := strings.Join(in.Number, "; ")
	return &pb.TrackResponse{
		Number: resp,
		Status: "Ok",
	}, nil
}

func (s *server) Listen(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)
	pb.RegisterTrackServiceServer(grpcServer, s)
	log.Printf("Server is listening on %s", lis.Addr())
	if err = grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
