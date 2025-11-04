package server

import (
	"context"
	"core-service/internal/domain/track"
	"core-service/internal/domain/user"
	"core-service/internal/usecase/tservice"
	"core-service/internal/usecase/uservice"
	"fmt"

	"log"
	"net"

	"github.com/Strelcock/pb/bot/pb"
	"google.golang.org/grpc"
)

type server struct {
	*userServer
	*trackServer
	pb.TrackerClient
}

type toTracker struct {
	pb.TrackServiceClient
}

type trackServer struct {
	pb.UnimplementedTrackServiceServer
	TrackService *tservice.TrackService
}

type userServer struct {
	pb.UnimplementedUserServiceServer
	UserService *uservice.UserService
}

func New(uService *uservice.UserService, tService *tservice.TrackService, client pb.TrackerClient) *server {
	return &server{
		&userServer{
			UserService: uService,
		},
		&trackServer{
			TrackService: tService,
		},
		client,
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
	//create new tracks
	fmt.Println(in.Number)
	var tracks = []track.Track{}
	for _, n := range in.Number {
		tr := track.New(n, in.User)
		tracks = append(tracks, *tr)
	}

	//subscribe to changes
	errCh := make(chan error, 16)

	go func(errCh chan error) {
		_, err := s.TrackerClient.ServeTrack(ctx, &pb.ToTracker{
			Number: in.Number,
		})
		if err != nil {
			errCh <- err
		}
		close(errCh)
	}(errCh)

	//write to db
	err := s.TrackService.Create(tracks)
	if err != nil {
		return &pb.TrackResponse{
			Status: "Что-то пошло не так в ядре",
		}, err
	}

	err = <-errCh

	//response
	if err != nil {
		return &pb.TrackResponse{
			Status: "Что-то пошло не так на сервисе отслеживания",
		}, err
	}
	return &pb.TrackResponse{
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
