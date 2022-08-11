package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

type Size struct {
	a [6]int32
	b rune
}

func f() (r int) {
	// 如果返回值不命名, defer func 中的r是一份拷贝值, 无法对外面的r修改
	defer func() {
		r = r + 5
	}()
	return
}

// round 字节填充测试函数
func round(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}
func TestRound(t *testing.T) {
	size := unsafe.Sizeof(Size{})
	fmt.Println(size)
	assert.True(t, size == round(25, 4))
}

func TestR(t *testing.T) {
	r := f()
	t.Log(r)
	assert.Equal(t, r, 5, "not equal 5")
}

func TestASyncProducer(t *testing.T) {
	go ASyncProducer()
	SaramaConsumerGroup()
}
