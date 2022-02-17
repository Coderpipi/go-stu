package marshal

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Animal struct {
	Name string `json:"name" gorm:"gcc"`
	Age  uint8  `json:age`
}

func Marshal(model Animal) {
	str, err := json.Marshal(model)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	fmt.Println(string(str))
	//val := reflect.ValueOf(model)
	//if val.CanSet() {
	//	fmt.Println("可以设置")
	//} else {
	//	name := val.Field(0)
	//	fmt.Println(name.String())
	//}
	typ := reflect.TypeOf(model)
	lookup, ok := typ.Field(0).Tag.Lookup("json")
	fmt.Println(lookup,ok)

}
