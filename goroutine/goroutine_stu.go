package goroutine

import "fmt"

func GetMsg(c <-chan int) {
	msg := <-c
	fmt.Println("go here GetMsg")
	fmt.Println(msg)
}
func ChannelStu() {
	chan1 := make(chan int)
	msg := 1
	go GetMsg(chan1)
	fmt.Println("go here")
	chan1 <- msg
}
