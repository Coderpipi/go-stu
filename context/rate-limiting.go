package context
//
//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	// 实现一个简单的限流策略
//	// 模拟瞬间有5个请求进来
//	requests := make(chan int, 5)
//	for i := 0; i < 5; i++ {
//		requests <- i
//	}
//
//	close(requests)
//	// 设置一个limiter，每200毫秒往通道里发送一个值
//	limiter := time.Tick(200 * time.Millisecond)
//	for req := range requests {
//		<-limiter
//		fmt.Println("request", req, time.Now())
//	}
//	fmt.Println("=================================")
//	// 第二种方式,临时对某一段业务进行速率限制，不影响整体逻辑
//	// 通过通道缓冲
//	burstyLimiter := make(chan time.Time, 3)
//	// 首先将通道占满阻塞
//	for i := 0; i < 3; i++ {
//		burstyLimiter <- time.Now()
//	}
//
//	// 另起一个goroutine,每200毫秒往通道里放值
//	go func() {
//		for t := range time.Tick(200 * time.Millisecond) {
//			burstyLimiter <- t
//		}
//	}()
//
//	// 现在模拟有5个请求，同时打进来
//	burstyRequests := make(chan int, 5)
//	for i :=0 ; i < 5; i++ {
//		burstyRequests <- i
//	}
//	close(burstyRequests)
//	for req := range burstyRequests {
//		<- burstyLimiter
//		fmt.Println("request", req, time.Now())
//	}
//}
