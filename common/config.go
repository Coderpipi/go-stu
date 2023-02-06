package common

import (
	"fmt"
	"github.com/caarlos0/env"
)

var Cfg *config

type config struct {
	KafKaServer string `env:"KAFKA_SERVER"`
}

func init() {
	Cfg = &config{}
	if err := env.Parse(Cfg); err != nil {
		panic(err)
	}
	fmt.Println(Cfg)
}
