package hello

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type Service struct {
	proto.UnimplementedHelloServiceServer
}

var _ proto.HelloServiceServer = (*Service)(nil)

func (h Service) SayHello(ctx context.Context, request *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "you should input your name!!!")
	}

	return &proto.SayHelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}
