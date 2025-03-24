package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"url-shorter/internal/api"
	"url-shorter/internal/cache"
	"url-shorter/internal/model"
	"url-shorter/internal/repo"
)

type (
	ShortCodeGenerator interface {
		GenerateShortCode() string
	}
)

type UrlService struct {
	repo               repo.IUrlRepo
	shortCodeGenerator ShortCodeGenerator
	defaultDuration    time.Duration
	baseUrl            string
}

func (s *UrlService) DeleteUrl(ctx context.Context) error {
	slog.Info("start delete expired url...")
	return s.repo.CleanupExpiredUrl(ctx)
}

func NewUrlService(db *gorm.DB, generator ShortCodeGenerator, duration time.Duration, baseUrl string) *UrlService {
	return &UrlService{
		repo:               repo.NewUrlRepo(db),
		shortCodeGenerator: generator,
		defaultDuration:    duration,
		baseUrl:            baseUrl,
	}
}

func (s *UrlService) GetUrl(ctx context.Context, shortCode string) (string, error) {
	data, err := cache.RedisCli.Get(ctx, shortCode).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	url := &repo.Url{}
	if len(data) > 0 {
		if err := json.Unmarshal(data, url); err != nil {
			return "", err
		}

		if url.OriginUrl != "" {
			return url.OriginUrl, nil
		}
	}

	url, err = s.repo.GetUrlByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	data, err = json.Marshal(url)
	if err != nil {
		return "", err
	}

	if err := cache.RedisCli.Set(ctx, url.ShortCode, data, time.Until(url.ExpiredAt)).Err(); err != nil {
		return "", err
	}

	return url.OriginUrl, nil

}

func (s *UrlService) CreateUrl(ctx context.Context, req *model.CreateUrlRequest) (*model.CreateUrlResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	param := repo.CreateUrlParam{
		OriginUrl: req.OriginUrl,
		ExpiredAt: time.Now().Add(s.defaultDuration),
	}

	if req.CustomCode != "" {
		available, err := s.repo.IsShortCodeAvailable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}

		if !available {
			return nil, fmt.Errorf("custom code is not available")
		}
		param.IsCustom = true
		param.ShortCode = req.CustomCode
	} else {
		code, err := s.getShortCode(ctx, 5)
		if err != nil {
			return nil, err
		}

		param.ShortCode = code
	}

	if req.Duration > 0 {
		param.ExpiredAt = time.Now().Add(time.Duration(req.Duration))
	}

	url, err := s.repo.CreateUrl(ctx, param)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	if err := cache.RedisCli.Set(ctx, url.ShortCode, data, time.Until(url.ExpiredAt)).Err(); err != nil {
		return nil, err
	}

	return &model.CreateUrlResponse{
		ShortUrl:  s.baseUrl + "/" + url.ShortCode,
		ExpiredAt: url.ExpiredAt,
	}, nil

}

func (s *UrlService) getShortCode(ctx context.Context, retryCount int) (string, error) {
	if retryCount <= 0 {
		return "", errors.New("重试次数过多")
	}

	shortCode := s.shortCodeGenerator.GenerateShortCode()
	available, err := s.repo.IsShortCodeAvailable(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if available {
		return shortCode, nil
	}

	return s.getShortCode(ctx, retryCount-1)
}

var _ api.IUrlService = (*UrlService)(nil)
