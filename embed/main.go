package main

import (
	_ "embed"
	"fmt"
)

//go:embed hello.txt
var S string

const (
	a = 1 << iota
	b
	c
)

func main() {
	fmt.Println(S)
	fmt.Println("a = ", a)
	fmt.Println("b = ", b)
	fmt.Println("c = ", c)
}
