package token

import (
	"context"
	"os"

	"github.com/fatih/structs"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-stu/proto"
)

const (
	claimID = "user_id"
)

type Issuer interface {
	IssueToken(ctx context.Context, userID string) (string, error)
}

type jwtCredentials struct {
	tokenIssuer Issuer
}

func NewJwtCredentials(issuer Issuer) *jwtCredentials {
	if issuer == nil {
		panic("issuer must be set")
	}

	return &jwtCredentials{tokenIssuer: issuer}
}

func (j jwtCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	info, ok := credentials.RequestInfoFromContext(ctx)

	if !ok || info.Method != proto.InterceptorService_Protected_FullMethodName {
		return nil, nil

	}

	token, err := j.tokenIssuer.IssueToken(ctx, "user-id-123")
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"authorization": token,
	}, nil
}

func (j jwtCredentials) RequireTransportSecurity() bool {
	return false
}

type Service struct {
	proto.UnimplementedTokenServiceServer
}

func (s Service) Validate(ctx context.Context, request *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	claims, ok := ctx.Value(claimsKey).(map[string]string)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "claims missing from context")
	}

	return &proto.ValidateResponse{Claims: claims}, nil
}

func (s Service) IssueToken(ctx context.Context, request *proto.IssueTokenRequest) (*proto.IssueTokenResponse, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok || jwtSecret == "" {
		return nil, status.Error(codes.Internal, "env missing var JWT_SECRET")
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(structs.Map(&Token{Sub: "lbwnb.com", Iss: request.UserId, Aud: "public"})))
	token, err := claims.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to issue token")
	}

	return &proto.IssueTokenResponse{Token: token}, nil
}
