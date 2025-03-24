package model

import "time"

type (
	CreateUrlRequest struct {
		OriginUrl  string `json:"origin_url" binding:"required,url"`
		CustomCode string `json:"custom_code,omitempty" binding:"omitempty,min=4,max=10,alphanum"`
		Duration   int    `json:"duration,omitempty" binding:"omitempty,min=1,max=100"`
	}

	CreateUrlResponse struct {
		ShortUrl  string    `json:"short_url"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)
