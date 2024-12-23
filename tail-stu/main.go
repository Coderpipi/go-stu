package main

import (
	"log"

	"github.com/hpcloud/tail"
)

func main() {
	fileName := "/Users/pipi/GolandProjects/go-stu/tail-stu/test.log"

	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	file, err := tail.TailFile(fileName, config)
	if err != nil {
		log.Fatal(err)
	}

	// 处理数据
	for msg := range file.Lines {
		log.Println("msg: ", msg.Text)
	}
}
