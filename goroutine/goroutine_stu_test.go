package main

import "testing"

func TestChannelStu(t *testing.T) {
	ChannelStu()
}

func TestChannelStu1(t *testing.T) {
	var ch chan int
	ch <- 1
}
