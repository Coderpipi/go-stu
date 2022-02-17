package stringslice

import (
	"fmt"
	"github.com/caarlos0/env"
	"testing"
)

type Cfg struct {
	Env string `env:"ENV"`
}

func TestEnv(t *testing.T) {
	var cfg Cfg
	env.Parse(&cfg)
	fmt.Println(cfg)
}

//func _TestStringSliceEqual(t *testing.T) {
//
//	// 第一行表示执行的测试函数
//	// 第二行表示测试的结果， 绿色代表通过，黄色x代表测试结果与预期不符
//	Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
//		a := []string{"hello", "go convey"}
//		b := []string{"hello", "go convey"}
//		So(StringSliceEqual(a, b), ShouldBeFalse)
//	})
//
//}
