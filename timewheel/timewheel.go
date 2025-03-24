package main

import (
	_ "embed"
	"fmt"
	"runtime"
	"time"
)

type (
	Task struct {
		CallFunc func()

		// 定时任务的唯一标识键
		Key string

		ExecuteAtTime time.Time
	}
)

func NewTask(key string, fn func(), t time.Time) *Task {
	// 返回一个指向Task结构体的指针，该结构体包含任务的详细信息
	return &Task{
		// 设置任务的唯一标识符
		Key: key,
		// 设置任务的执行函数
		CallFunc: fn,
		// 设置任务执行的指定时间
		ExecuteAtTime: t,
	}
}

func main() {
	// wheel := NewSingletonTimeWheel(0, time.Second)
	// go wheel.Run()
	// defer wheel.Stop()
	// wheel.AddTask(NewTask("1", func() {
	// 	fmt.Println("你好")
	// }, time.Now().Add(time.Second)))
	// wheel.AddTask(NewTask("2", func() {
	// 	fmt.Println("你好, 10")
	// }, time.Now().Add(time.Second*10)))

	wheel := NewRTimeWheel("localhost:6379")
	go wheel.Run()
	err := wheel.AddTask("1", NewRTaskElement("http://localhost:8080/ping", nil, nil), time.Now().Add(time.Minute))
	if err != nil {
		fmt.Println(err)
	}

	runtime.Gosched()
	for {

	}
}
