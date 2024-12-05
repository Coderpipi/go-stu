package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	claimID = "id"
)

type Service struct {
	secret []byte
}

func NewService(secret string) *Service {
	if secret == "" {
		panic("secret must not be empty")
	}

	return &Service{secret: []byte(secret)}
}

func (s *Service) IssueToken(ctx context.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		claimID: userID,
	})

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signed, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (string, error) {
	decodeToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected siging method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token:%w", err)
	}

	claims, ok := decodeToken.Claims.(jwt.MapClaims)
	if !decodeToken.Valid || !ok {
		return "", errors.New("failed to extract claims")
	}

	id, ok := claims[claimID].(string)
	if !ok {
		return "", errors.New("cannot extract user ID from claims")
	}

	return id, nil
}
