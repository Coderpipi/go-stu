package main

import (
	"fmt"
)

func GetMsg(c chan<- int) {
	fmt.Println("go here GetMsg")
	fmt.Println("get msg goroutine exit!!!")
	c <- 1
}

func ChannelStu() {
	chan1 := make(chan int)
	//msg := 1
	go GetMsg(chan1)
	<-chan1
	fmt.Println("main goroutine exit!!!")

}
