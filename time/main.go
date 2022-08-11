package main

import (
	"fmt"
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

func Ticker() {
	ticker := time.NewTicker(time.Second)
	i := 1
Loop:
	for {
		select {
		case <-ticker.C:
			fmt.Printf("1s执行任务: %d\n", i)
			if i == 3 {
				break Loop
			}
			i++
		}
	}
}
func Timer() {
	// 延时任务
	timer := time.NewTimer(3 * time.Second)
	i := 1
Loop:
	for {
		select {
		case <-timer.C:
			{
				fmt.Println("3秒后执行任务:", i)
				break Loop
			}
		}
	}
}
