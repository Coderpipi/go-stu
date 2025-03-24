package repo

import "time"

type Url struct {
	ID        int64     `json:"id"`
	OriginUrl string    `json:"origin_url"`
	ShortCode string    `json:"short_code"`
	IsCustom  bool      `json:"is_custom"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (Url) TableName() string {
	return "urls"
}
