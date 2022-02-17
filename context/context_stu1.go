package context

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var channel chan bool

func f1(ctx context.Context) {
	defer wg.Done()
LOOP:
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		//case <-channel:
		//	break LOOP
		case <-ctx.Done():
			fmt.Println("f1 准备退出。。。。。。")
			break LOOP
		default:
			fmt.Println("Hello World")

		}
	}
	fmt.Println("f1已经退出。。。。。。")

}
func inSeq() int {
	i := 0
	return func() int {
		i++
		return i
	}()
}
//func main() {
//	ctx, cancel := context.WithCancel(context.Background())
//	//channel = make(chan bool)
//	wg.Add(1)
//	go f1(ctx)
//	time.Sleep(5 * time.Second)
//	cancel()
//	wg.Wait()
//
//}
