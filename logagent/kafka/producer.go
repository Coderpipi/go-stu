package kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Producer sarama.SyncProducer
	MsgChan  chan *sarama.ProducerMessage
)

func SendMessage() {
	for {
		select {
		case msg, ok := <-MsgChan:
			if !ok {
				return
			}

			if _, _, err := Producer.SendMessage(msg); err != nil {
				logrus.WithField("topic", msg.Topic).WithError(err).Error("send message failed")
			}

		}
	}
}
