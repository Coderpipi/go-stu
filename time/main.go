package main

import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
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

// Ticker 定时任务
func Ticker() {
	runtime.GOMAXPROCS(3)
	ticker := time.NewTicker(time.Second)
	i := 1
	pointer := unsafe.Pointer(&i)
	fmt.Println(pointer)
Loop:
	for {
		select {
		case <-ticker.C:
			fmt.Printf("每隔1s执行任务: %d\n", i)
			if i == 3 {
				break Loop
			}
			i++
		}
	}
}

// Timer 延时任务
func Timer() {
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
