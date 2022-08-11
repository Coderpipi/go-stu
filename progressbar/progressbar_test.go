package main

import (
	"github.com/schollz/progressbar/v3"
	"testing"
	"time"
)

func TestSimpleProgressBar(t *testing.T) {
	r := [100]struct{}{}
	bar := progressbar.Default(int64(len(r)))
	for range r {
		bar.Add(1)
		time.Sleep(20 * time.Millisecond)
	}

}
