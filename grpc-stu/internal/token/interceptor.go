package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

const (
	claimsKey = "claims"
)

type middleware struct {
	secret []byte
}

func NewMiddleware(secret string) *middleware {
	if secret == "" {
		panic("secret cannot be empty")
	}

	return &middleware{secret: []byte(secret)}
}

func (m *middleware) UnaryAuthMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	if info.FullMethod == proto.TokenService_IssueToken_FullMethodName {
		return handler(ctx, req)
	}

	token, err := getToken(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "token must be provide")
	}

	verifiedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
		}

		return m.secret, nil
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "invalid token")
	}

	if !verifiedToken.Valid {
		return nil, status.Error(codes.PermissionDenied, "invalid token")
	}

	claimsMap := make(map[string]string)
	claims, ok := verifiedToken.Claims.(jwt.MapClaims)
	if ok {
		for key, val := range claims {
			claimsMap[key] = cast.ToString(val)
		}
	}

	ctx = context.WithValue(ctx, claimsKey, claimsMap)

	return handler(ctx, req)
}

func getToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) != 1 {
		return "", errors.New("token not found in metadata")
	}

	return md["authorization"][0], nil
}
