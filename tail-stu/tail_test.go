package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteLog(t *testing.T) {
	file, err := os.OpenFile("/Users/pipi/GolandProjects/go-stu/tail-stu/test.log", os.O_RDWR|os.O_APPEND, 0777)
	assert.NoError(t, err)
	defer file.Close()
	for i := range 10 {
		file.WriteString(fmt.Sprintf("log%d\n", i))
	}
}
