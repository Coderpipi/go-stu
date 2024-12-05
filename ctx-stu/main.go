package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx1, cancel1 := context.WithCancel(context.Background())
	cancel1()
	ctx2, _ := context.WithTimeout(ctx1, time.Second*3)

	select {
	case <-ctx2.Done():
		fmt.Println("触发")
	}
}
