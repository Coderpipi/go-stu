package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	dial, err := net.DialTimeout("tcp", "103.235.46.40:12345", 3*time.Second)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer func(dial net.Conn) {
		err := dial.Close()
		if err != nil {

		}
	}(dial)
	fmt.Println("test complete")

}
