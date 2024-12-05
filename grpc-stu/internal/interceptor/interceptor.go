package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

type Validator interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}

type middleware struct {
	validator Validator
}

func NewMiddleware(validator Validator) (*middleware, error) {
	if validator == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &middleware{validator: validator}, nil
}

func (m *middleware) UnaryAuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// 检查rpc方法, 当他是protected方法才做校验
	if info.FullMethod != proto.InterceptorService_Protected_FullMethodName {
		return handler(ctx, req)
	}

	token, err := getTokenFormMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "token must be provided")
	}

	userID, err := m.validator.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "invalid token")
	}

	ctx = context.WithValue(ctx, "user_id", userID)
	return handler(ctx, req)

}

func getTokenFormMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || len(md["authorization"]) != 1 {
		return "", errors.New("token not found in metadata")
	}

	return md["authorization"][0], nil
}
