package config

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type Service struct {
	name string

	proto.UnimplementedConfigServiceServer
}

func NewService(name string) *Service {
	if name == "" {
		panic("name is required")
	}

	return &Service{name: name}
}

func (s Service) ConfigLongRunning(ctx context.Context, request *proto.ConfigLongRunningRequest) (*proto.ConfigLongRunningResponse, error) {
	time.Sleep(time.Second * 3)
	return &proto.ConfigLongRunningResponse{}, nil
}

func (s Service) Flaky(context.Context, *proto.FlakyRequest) (*proto.FlakyResponse, error) {
	if rand.IntN(3) != 0 {
		log.Println("error response returned")
		return nil, status.Error(codes.Internal, "flaky error occurred")
	}

	log.Println("flaky response returned")

	return &proto.FlakyResponse{}, nil
}

func (s Service) GetServerAddress(context.Context, *proto.GetServerAddressRequest) (*proto.GetServerAddressResponse, error) {
	log.Printf("request received: %s", s.name)
	return &proto.GetServerAddressResponse{Address: s.name}, nil
}
