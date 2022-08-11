package marshal

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {

	m1 := map[int]interface{}{1: "John", 2: 10, 3: [0]int{}}

	b, err := json.Marshal(m1)
	if err != nil {
		fmt.Println("json.Marshal failed:", err)
		return
	}

	fmt.Println("b:", string(b))
}
