package config

import (
	"fmt"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
)

func InitKafka() error {

	// 生产者配置
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Partitioner = sarama.NewRandomPartitioner
	conf.Producer.Return.Successes = true

	// 链接 kafka
	client, err := sarama.NewSyncProducer(Cfg.Kafka.Addr, conf)
	if err != nil {
		return fmt.Errorf("create producer failed, err: %w", err)
	}

	kafka.Producer = client

	logrus.Info("kafka init success...")

	kafka.MsgChan = make(chan *sarama.ProducerMessage, Cfg.MsgChanSize)

	return nil

}
