package common

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

var Cfg *config

type config struct {
	KafkaServer string `env:"KAFKA_SERVER"`
}

func init() {
	if len(os.Getenv("ENV")) == 0 {
		// if no ENV specified
		envFile := ".env"
		for i := 0; i < 10; i++ {
			if _, err := os.Stat(envFile); err != nil {
				envFile = "../" + envFile
				continue
			}
			break
		}
		// local debugging
		err := godotenv.Load(envFile)
		if err != nil {
			panic(err)
		}

		logrus.Trace("dot env loaded")
	}

	Cfg = &config{}
	if err := env.Parse(Cfg); err != nil {
		panic(err)
	}
	fmt.Println(Cfg)
}
