package skip_list

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkList(t *testing.T) {
	head := &nd[int, string]{
		key:   0,
		val:   "",
		nexts: nil,
	}

	sklist := SkList[int, string]{
		head: head,
	}
	sklist.Insert(1, "1")
	sklist.Insert(2, "2")
	sklist.Insert(3, "3")
	sklist.Insert(4, "4")
	sklist.Insert(5, "5")
	e, exist := sklist.Find(1)
	assert.True(t, exist)
	t.Log(e)

	sklist.Remove(4)
	find, b := sklist.Find(4)
	assert.False(t, b)
	assert.Equal(t, "", find)

	n := sklist.ceiling(4)
	assert.Equal(t, 5, n.key)

	floor := sklist.floor(3)
	assert.Equal(t, 3, floor.key)
	n2 := sklist.floor(6)
	assert.Equal(t, 5, n2.key)

	for k, v := range sklist.All() {
		fmt.Printf(" key: %d, val: %s\n", k, v)
	}
}

func TestSkList_floor(t *testing.T) {
	head := &nd[int, string]{
		key:   0,
		val:   "",
		nexts: nil,
	}

	sklist := SkList[int, string]{
		head: head,
	}

	testcases := func() (s []struct {
		key      int
		val      string
		expected int
	}) {
		for i := 1; i <= 10; i++ {
			v := struct {
				key      int
				val      string
				expected int
			}{
				key:      i,
				val:      fmt.Sprintf("%d", i),
				expected: i,
			}
			sklist.Insert(v.key*2, v.val)
			s = append(s, v)
		}
		return
	}()
	for _, c := range testcases {
		n := sklist.floor(c.key)
		t.Logf("key: %d, expected: %d, node: %v", c.key, c.expected, n)
	}
}

func TestSkList_ceiling(t *testing.T) {
	head := &nd[int, string]{
		key:   0,
		val:   "",
		nexts: nil,
	}

	sklist := SkList[int, string]{
		head: head,
	}

	testcases := func() (s []struct {
		key      int
		val      string
		expected int
	}) {
		for i := 1; i <= 10; i++ {
			v := struct {
				key      int
				val      string
				expected int
			}{
				key:      i,
				val:      fmt.Sprintf("%d", i),
				expected: i,
			}
			sklist.Insert(v.key*2, v.val)
			s = append(s, v)
		}
		return
	}()
	for _, c := range testcases {
		n := sklist.ceiling(c.key)
		t.Logf("key: %d, expected: %d, node: %v", c.key, c.expected, n)
	}
}

func TestSkListAll(t *testing.T) {
	head := &nd[int, string]{
		key:   0,
		val:   "",
		nexts: nil,
	}

	sklist := SkList[int, string]{
		head: head,
	}

	for i := 1; i <= 10; i++ {
		v := struct {
			key      int
			val      string
			expected int
		}{
			key:      i,
			val:      fmt.Sprintf("%d", i),
			expected: i,
		}
		sklist.Insert(v.key*2, v.val)
	}

	for k, v := range sklist.All() {
		t.Logf("key: %d,val: %s", k, v)
	}

	for k, v := range sklist.Range(1, 4) {
		t.Logf("key: %d,val: %s", k, v)
	}

}
