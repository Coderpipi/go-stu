package kafka

import (
	"log"
	"testing"

	"github.com/IBM/sarama"
)

func TestConsumer(t *testing.T) {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		return
	}
	partition, err := consumer.ConsumePartition("log_collect", 0, sarama.OffsetNewest)
	for {
		msg := <-partition.Messages()
		log.Printf("message offset %d, msg is : %s\n", msg.Offset, string(msg.Value))
	}

}
