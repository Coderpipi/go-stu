package hello

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type Service struct {
	proto.UnimplementedHelloServiceServer
}

var _ proto.HelloServiceServer = (*Service)(nil)

func (h *Service) SayHello(ctx context.Context, request *proto.SayHelloRequest) (*proto.SayHelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["x-request-id"]) == 0 {
		return nil, status.Error(codes.InvalidArgument, "missing request-id")
	}
	requestId := md["x-request-id"]

	header := metadata.New(map[string]string{
		"x-request-start-timestamp": time.Now().String(),
	})

	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, status.Error(codes.Internal, "failed to send header")
	}

	if err := grpc.SetTrailer(ctx, header); err != nil {
		return nil, status.Error(codes.Internal, "failed to set trailer")
	}

	log.Printf("request ID: %v", requestId)

	if request.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "you should input your name!!!")
	}

	return &proto.SayHelloResponse{Message: fmt.Sprintf("Hello %s", request.Name)}, nil
}

func (h *Service) LongRunning(ctx context.Context, request *proto.LongRunningRequest) (*proto.LongRunningResponse, error) {

	select {
	case <-ctx.Done():
		log.Println("context canceled")
		return nil, ctx.Err()
	case <-time.Tick(time.Second * 5):
		log.Println("finished waiting, not end request successfully")
	}
	return &proto.LongRunningResponse{}, nil
}
