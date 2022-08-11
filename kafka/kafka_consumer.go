package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
)

type ConsumerGroup struct {
}

func (c ConsumerGroup) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerGroup) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 处理消息
func (c ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		log.Printf("Message Topic: %v,  Value: %s, TimeStamp: %v", msg.Topic, msg.Value, msg.Timestamp)
	}

	return nil
}

// SaramaConsumerGroup 消费者组
func SaramaConsumerGroup() {
	config := sarama.NewConfig()
	// 是否设置消费者错误消息回传管道, 开启则消费者发生错误是会将错误消息通过.Error管道回传
	config.Consumer.Return.Errors = true
	// 指定kafka的版本
	config.Version = sarama.V1_1_0_0
	// 设置是否自动提交消费标记, 开启之后需要设置自动提交的间隔时间
	config.Consumer.Offsets.AutoCommit.Enable = false
	// 设置新加入的消费者从那个位置开始消费,默认从最新的数据开始消费
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 创建ConsumerGroup
	consumerGroup, err := sarama.NewConsumerGroup([]string{"106.52.172.56:9092"}, "test_group", config)
	if err != nil {
		log.Fatalln("cannot create ConsumerGroup")
	}
	defer consumerGroup.Close()
	// 跟踪回传的错误消息
	go func() {
		for err := range consumerGroup.Errors() {
			log.Fatalln(err.Error())
		}
	}()

	// 监听消息队列, 消费消息
	log.Println("ConsumerGroup Start")

	//
	for {
		// 定义主题
		topics := []string{"topic_test1"}
		handler := ConsumerGroup{}

		// 消费消息
		err := consumerGroup.Consume(context.TODO(), topics, handler)
		if err != nil {
			panic(err)
		}
	}

}

// SaramaConsumer 消费者
func SaramaConsumer() {
	consumer, err := sarama.NewConsumer([]string{"106.52.172.56:9092"}, nil)
	if err != nil {
		// 关闭消费者
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
	}
	// Consumer消费
	partition, err := consumer.ConsumePartition("topic_test1", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect success...")
	defer func() {
		if err := partition.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partition.Messages():
			log.Printf("message offset %d, msg is : %s\n", msg.Offset, string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
	log.Printf("Consumed: %d\n", consumed)
}
