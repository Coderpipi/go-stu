package repo

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type (
	IRepo interface {
		DB() *gorm.DB
	}

	IUrlRepo interface {
		IRepo
		CreateUrl(context.Context, CreateUrlParam) (*Url, error)
		GetUrlByShortCode(ctx context.Context, shortCode string) (*Url, error)
		CleanupExpiredUrl(ctx context.Context) error
		IsShortCodeAvailable(ctx context.Context, shortCode string) (bool, error)
	}
)

type (
	UrlRepo struct {
		db *gorm.DB
	}

	CreateUrlParam struct {
		OriginUrl string    `json:"origin_url"`
		ShortCode string    `json:"short_code"`
		IsCustom  bool      `json:"is_custom"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)

func (u *UrlRepo) CleanupExpiredUrl(ctx context.Context) error {
	return u.db.WithContext(ctx).Where("expired_at <= CURRENT_TIMESTAMP").Delete(&Url{}).Error
}

func NewUrlRepo(db *gorm.DB) *UrlRepo {
	return &UrlRepo{
		db: db,
	}
}

func (u *UrlRepo) GetUrlByShortCode(ctx context.Context, shortCode string) (*Url, error) {
	url := &Url{}

	if err := u.db.WithContext(ctx).Debug().Model(url).Where("short_code = ? and expired_at > CURRENT_TIMESTAMP", shortCode).Find(url).Error; err != nil {
		return nil, err
	}

	return url, nil
}

func (u *UrlRepo) DB() *gorm.DB {
	return u.db
}

func (u *UrlRepo) CreateUrl(ctx context.Context, param CreateUrlParam) (*Url, error) {
	url := &Url{
		ShortCode: param.ShortCode,
		ExpiredAt: param.ExpiredAt,
		OriginUrl: param.OriginUrl,
		IsCustom:  param.IsCustom,
	}
	if err := u.db.WithContext(ctx).Model(url).Create(url).Error; err != nil {
		return nil, err
	}

	return url, nil
}

func (u *UrlRepo) IsShortCodeAvailable(ctx context.Context, shortCode string) (bool, error) {
	var cnt int64 = 0
	err := u.db.WithContext(ctx).Model(&Url{}).Where("short_code = ?", shortCode).Count(&cnt).Error
	if err != nil {
		return false, err
	}

	return cnt == 0, nil
}

var _ IUrlRepo = (*UrlRepo)(nil)
