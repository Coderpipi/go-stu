package config

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Cfg *Config

type (
	Config struct {
		BasePath string `env:"BASE_PATH"`
		Kafka    `yaml:"kafka" mapstructure:"kafka"`
		Collect  `yaml:"collect" mapstructure:"collect"`
	}

	Kafka struct {
		Addr        []string `yaml:"addr" mapstructure:"addr"`
		MsgChanSize int      `mapstructure:"msg-chan-size"`
	}

	Collect struct {
		LogFilePath string `yaml:"log-file-path" mapstructure:"log-file-path"`
	}
)

func InitConfig() error {
	cfg := &Config{}
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config error: %w", err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("parse config error: %w", err)
	}

	Cfg = cfg

	logrus.Info(cfg)

	return nil
}
