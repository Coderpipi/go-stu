package main

import (
	"awesomeProject/request"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type User struct {
	Id   uint64
	Name string
}
type Order struct {
	Id         uint64
	OrderId    uint64
	OrderPrice uint64
}

func Timer() {
	//timer := time.NewTimer(3 * time.Second)
	ticker := time.NewTicker(3 * time.Second)
	i := 1
	for {
		select {
		case <-ticker.C:
			{
				fmt.Println("3秒后执行任务:", i)
			}
		}
		if i == 3 {
			ticker.Stop()
			break
		}
		i++
	}
	//for range ticker.C {
	//	fmt.Println("每隔3秒执行任务")
	//}
	//ticker.Stop()

	//var user = User{Id: 1, Name: "lbw"}
	//var order Order = Order{Id:2, OrderId: 11111, OrderPrice: 1000}
	//var a interface{} = user
	//typ := reflect.TypeOf(a)
	//val := reflect.ValueOf(&user.Id)

	//if typ.String() == "main.Order" {
	//	fmt.Println("this type is Order")
	//} else {
	//	fmt.Println("I don't know")
	//}
	//fmt.Println(order)
	//fmt.Println(typ)
	//fmt.Println(val.Kind())
	//fmt.Println(val.Kind())
	//if val.Kind() == reflect.Uint64 {
	//
	//	fmt.Println("this is a uint64")
	//	val.SetUint(1)
	//} else if val.Kind() == reflect.Ptr {
	//	fmt.Println("this is a Pointer")
	//
	//
	//	val.Elem().SetUint(2)
	//	fmt.Println(user)
	//}
}

type ConsumerGroupHandler struct {
}

func (c ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费者消费消息
func (c ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message Claimed: value = %s, timestamp = %v, topic = %s", msg.Value, msg.Timestamp, msg.Topic)
		// 将消息成功标记为处理，然后会自动提交
		session.MarkMessage(msg, "")
	}
	return nil
}

// SaramaConsumerGroup 消费者组
func SaramaConsumerGroup() {
	// 使用默认配置文件
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Version = sarama.V1_1_0_0
	// 是否开启自动提交，不管是否是自动还是手动，都需要标记之后才能提交
	config.Consumer.Offsets.CommitInterval = 1 * time.Second
	// 设置每次新增的消费者未找到消费者组的消费位移时从哪端开始消费
	// OffsetOldest: 表示总是拉取当前消息队列里最老的消息开始消费
	// OffsetNewest: 表示总是拉取当前消息队列里最新的消息开始消费
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// 创建consumerGroup
	group, err := sarama.NewConsumerGroup([]string{"106.52.172.56:9092"}, "octopus", config)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = group.Close()
	}()

	// 错误跟踪 Track Error
	go func() {
		for err := range group.Errors() {
			fmt.Println("err: ", err)
		}
	}()
	// 表示监听消息队列，准备消费信息
	fmt.Println("Consumed Start")

	// 迭代消费者会话
	ctx := context.Background()

	for {
		// 定义主题
		topics := []string{"image_sync_v2"}
		handler := ConsumerGroupHandler{}

		// 开始消费消息
		// 第三个参数，一个实现ConsumeGroupHandler的结构体
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
		fmt.Println("Kfldfkdsl;f")
	}
}

// SaramaConsumer 消费者
func SaramaConsumer(name string) {
	consumer, err := sarama.NewConsumer([]string{"106.52.172.56:9092"}, nil)
	if err != nil {
		panic(err)
	}
	// 关闭消费者
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	// ConsumerPartition
	// param1: 消费者订阅哪个主题
	partitionConsumer, err := consumer.ConsumePartition("image_sync_v2", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println(name + " connect success...")
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	// 捕捉操作系统SIGINT信号，造成一个中断，以至关闭消费者
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	consumed := 0
	// 开始循环监听主题是否有消息
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumer %s message offset %d, msg is : %s\n", name, msg.Offset, string(msg.Value))
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}
	log.Printf("Consumed: %d\n", consumed)
}

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
		str := new(string)
		*str = "则是一张图片"
		req1 := request.Request{
			ID:      "1",
			Type:    "JPEG",
			Src:     "http://localhost:8888/1.jpg",
			StoreID: "222",
			Alt:     str,
		}
		req2 := request.Request{
			ID:      "1",
			Type:    "JPEG",
			Src:     "http://localhost:8888/1.jpg",
			StoreID: "222",
			Alt:     str,
		}
		images := make([]request.Request, 0)
		images = append(images, req1, req2)
		info, _ := json.Marshal(images)
		message := &sarama.ProducerMessage{Topic: "image_sync_v2", Value: sarama.ByteEncoder(info)}
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

// SyncProducerSelect
// @Description: 异步生产消息
//
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
		case producer.Input() <- &sarama.ProducerMessage{Topic: "my_topic", Key: nil, Value: sarama.StringEncoder("testing 123")}:
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
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}

}
func f() (r int) {
	// 这里面的r是一份拷贝值, 无法对外面的r修改
	defer func() {
		r = r + 5
	}()
	return 1
}
func round(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}

func main() {
	//timeOut := int(3 * time.Second)
	//fmt.Println(timeOut)
	//生产者
	//ASyncProducer()
	//go SaramaProducer()
	//fmt.Println(round(7, 4))
	//go SyncProducerSelect()
	//消费者
	//SaramaConsumerGroup()
	//go SaramaConsumer("小明")
	//SaramaConsumer("小红")
	//time.Sleep(5 * time.Second)
	//res := make(chan bool)
	//go func() {
	//	defer close(res)
	//	time.Sleep(3 * time.Second)
	//}()
	//select {
	//case <-res:
	//	fmt.Println("3秒后接收到消息,退出")
	//}

}
