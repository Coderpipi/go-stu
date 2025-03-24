package config

import (
	"context"

	"github.com/redis/go-redis/v9"
	"url-shorter/internal/cache"
)

type Redis struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func InitRedis() {
	cli := redis.NewClient(&redis.Options{
		DB:   Cfg.Redis.DB,
		Addr: Cfg.Redis.Address,
	})

	if err := cli.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	cache.RedisCli = cli
}
