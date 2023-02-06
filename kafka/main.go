package main

import (
	"fmt"
	"time"
)

func main() {
	timeOut := int(3 * time.Second)
	fmt.Println(timeOut)
	// 生产者
	ASyncProducer()
	go SaramaProducer()
	go SyncProducerSelect()
	// 消费者
	SaramaConsumerGroup()
	go SaramaConsumer()
	SaramaConsumer()
	time.Sleep(5 * time.Second)
	res := make(chan bool)
	go func() {
		defer close(res)
		time.Sleep(3 * time.Second)
	}()
	select {
	case <-res:
		fmt.Println("3秒后接收到消息,退出")
	}

}
