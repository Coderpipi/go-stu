package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	Timer()
}

func TestTicker(t *testing.T) {
	Ticker()
}

func BenchmarkTimer(b *testing.B) {
	timer := time.After(3 * time.Second)
	b.ResetTimer()
	b.StartTimer()
	time.Sleep(4 * time.Second)
	select {
	case <-timer:
		fmt.Println("exit!!!")
	}
	b.StopTimer()
}
