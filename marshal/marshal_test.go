package marshal

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	animal := Animal{Name: "dog",Age: 3}
	Marshal(animal)

}
