package server

import (
	"context"
	"core-service/internal/domain/track"
	"core-service/internal/domain/user"
	"core-service/internal/usecase/tservice"
	"core-service/internal/usecase/uservice"
	"fmt"
	"strings"

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
	tr := track.New(in.Number, in.User)

	//subscribe to changes
	errCh := make(chan error, 16)

	go func(errCh chan error) {

		_, err := s.TrackerClient.ServeTrack(ctx, &pb.ToTracker{
			Number: in.Number,
		})
		if err != nil {
			errCh <- err
			log.Println("hi from secong goroutine with error")
		}

		close(errCh)
	}(errCh)

	//write to db
	err := s.TrackService.Create(tr)
	if err != nil {
		return &pb.TrackResponse{
			Status: "Что-то пошло не так в ядре",
		}, err
	}

	//response
	err = <-errCh

	if err != nil {

		if strings.Contains(err.Error(), "TrackerAlreadyExists") { //sub
			return &pb.TrackResponse{
				Status: "Вы Успешно подписались на уже отслеживаемый заказ",
			}, nil
		}

		return &pb.TrackResponse{
			Status: fmt.Sprintf("Заказ %s добавить не удалось: %s", in.Number, err.Error()),
		}, err

	}

	return &pb.TrackResponse{
		Status: fmt.Sprintf("Заказ %s добавлен успешно\n", in.Number),
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
