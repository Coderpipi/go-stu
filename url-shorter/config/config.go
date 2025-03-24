package config

import (
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var Cfg *Config

type (
	ServerConfig struct {
		Addr         string        `mapstructure:"addr"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	}

	App struct {
		BaseUrl         string        `mapstructure:"base_url"`
		DefaultDuration time.Duration `mapstructure:"default_duration"`
		CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
	}

	ShortCodeConfig struct {
		Length int `mapstructure:"length"`
	}

	Config struct {
		App             App             `mapstructure:"app"`
		ShortCodeConfig ShortCodeConfig `mapstructure:"shortcode_config"`
		Redis           Redis           `mapstructure:"redis"`
		DBConfig        DBConfig        `mapstructure:"db"`
		ServerConfig    ServerConfig    `mapstructure:"server_config"`
	}
)

func InitConfig(filePath string) {
	viper.SetConfigFile(filePath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		panic(err)
	}
	slog.Info("init config success")

	InitRedis()
	slog.Info("init redis success")
}
