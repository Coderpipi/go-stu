package main

import "fmt"

func main() {
	defer A()
	defer B()
	defer C()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	panic("panic main")
	fmt.Println("main complete")
}

func A() {
	if err := recover(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("A")
}

func B() {
	panic("panic B")
	fmt.Println("B")
}

func C() {
	fmt.Println("C")
}
