package interceptor

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

const (
	userIDCtxKey = "user_id"
)

type Service struct {
	proto.UnimplementedInterceptorServiceServer
}

func (s Service) Unprotected(context.Context, *proto.UnprotectedRequest) (*proto.UnprotectedResponse, error) {
	return &proto.UnprotectedResponse{}, nil
}

func (s Service) Protected(ctx context.Context, request *proto.ProtectedRequest) (*proto.ProtectedResponse, error) {
	userID, ok := ctx.Value(userIDCtxKey).(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.FailedPrecondition, "user id missing from context")
	}

	return &proto.ProtectedResponse{UserId: userID}, nil
}
