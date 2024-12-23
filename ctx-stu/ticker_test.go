package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	for {
		select {
		case <-time.Tick(time.Second):
			fmt.Println("2")
		}
	}

}
