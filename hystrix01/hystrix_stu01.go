package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"net/http"
	"time"
)

type Handle struct{}

func (h *Handle) ServeHTTP(r http.ResponseWriter, request *http.Request) {
	h.Common(r, request)
}

func (h *Handle) Common(r http.ResponseWriter, request *http.Request) {
	hystrix.ConfigureCommand("my_command", hystrix.CommandConfig{
		Timeout:                3000,
		MaxConcurrentRequests:  10,
		SleepWindow:            5000,
		RequestVolumeThreshold: 20,
		ErrorPercentThreshold:  30,
	})
	msg := "success"
	_ = hystrix.Do("my_command", func() error {
		_, err := http.Get("https://www.baidu.com")
		if err != nil {
			fmt.Printf("请求失败:%v\n", err)
			return err
		}
		return nil
	}, func(err error) error {
		fmt.Printf("handle  error:%v\n", err)
		msg = err.Error()
		return nil
	})
	r.Write([]byte(msg))
}

func main() {
	fmt.Println(2 ^ 2)
	time.Sleep(50 * time.Millisecond)
	http.ListenAndServe(":8090", &Handle{})
}
