package main

import (
	"fmt"
	"strconv"
)

type User struct {
	Age  int
	Name string
}

func (u *User) String() string {
	return u.Name + "\t" + strconv.Itoa(u.Age)
}

func (u *User) UpdateAge(age int) {
	u.Age = age
}
func main() {
	//u := new(User)
	u := User{
		Name: "Lbw",
		Age:  20,
	}
	u.UpdateAge(100)
	fmt.Println(u.String())
	x, y := 1, 2
	fmt.Println(x &^ y)
}
