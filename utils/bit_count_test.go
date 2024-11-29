package main

import (
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBitCount(t *testing.T) {
	Convey("test bit count", t, func() {
		n := BitCount(12)
		So(n, should.Equal, 2)
	})
}
