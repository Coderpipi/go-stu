package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

// ASyncProducer 异步生产者Goroutines
func ASyncProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer([]string{"106.52.172.56:9092"}, config)
	if err != nil {
		panic(err)
	}

	// 捕捉系统信号SIGINT,产生中断
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		// 用于线程计数，控制执行的goroutine顺序
		wg                                  sync.WaitGroup
		enqueued, successes, producerErrors int
	)
	// 同时只允许一个goroutine并发
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			producerErrors++
		}
	}()
	i := 0
ProducerLoop:
	for {
		time.Sleep(2 * time.Second)
		i++
		info := fmt.Sprintf("Message: %d", i)
		message := &sarama.ProducerMessage{Topic: "topic_test1", Value: sarama.StringEncoder(info)}
		// 开始发送消息
		select {
		// 发送一条消息到主题消息队列
		case producer.Input() <- message:
			enqueued++
		case <-signals:
			producer.AsyncClose()
			break ProducerLoop
		}
	}
	wg.Wait()
	log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)
}

// SyncProducerSelect 异步生产消息
func SyncProducerSelect() {
	producer, err := sarama.NewAsyncProducer([]string{"106.52.107.56:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var enqueued, producerErrors int
ProducerLoop:
	for {
		select {
		case producer.Input() <- &sarama.ProducerMessage{Topic: "topic_test1", Key: nil, Value: sarama.StringEncoder("testing 123")}:
			enqueued++
		case err := <-producer.Errors():
			log.Println("Failed to produce message", err)
			producerErrors++
		case <-signals:
			break ProducerLoop
		}
	}
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, producerErrors)
}

// SaramaProducer 同步生产模式
func SaramaProducer() {
	producer, err := sarama.NewSyncProducer([]string{"106.52.107.56:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	msg := &sarama.ProducerMessage{Topic: "my_topic", Value: sarama.StringEncoder("testing 123")}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}

}
