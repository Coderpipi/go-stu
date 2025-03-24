package skip_list

import (
	"cmp"
	"iter"
	"math/rand/v2"
	"slices"
)

// SkList 自己实现 skip list
type (
	SkList[K cmp.Ordered, E any] struct {
		// head, 跳表头节点, 不存储任何数据,只指向下一个 node 节点
		head *nd[K, E]
	}

	nd[K cmp.Ordered, E any] struct {
		key   K
		val   E
		nexts []*nd[K, E]
	}
)

// Find 查找跳表对应 key 的 val
func (s *SkList[K, E]) Find(key K) (e E, exist bool) {
	if n := s.search(key); n != nil {
		return n.val, true
	}

	return e, false
}

func (s *SkList[K, E]) Insert(k K, v E) {
	// 如果节点已经存在, 直接更新
	if n := s.search(k); n != nil {
		n.val = v
		return
	}

	// 获取新节点所在高度的 index
	high := s.roll()

	// 如果当前跳表高度不足,补齐
	for high > len(s.head.nexts)-1 {
		s.head.nexts = append(s.head.nexts, nil)
	}

	// 创建新节点
	_node := &nd[K, E]{
		key:   k,
		val:   v,
		nexts: make([]*nd[K, E], high+1),
	}

	h := s.head

	for i := len(h.nexts) - 1; i >= 0; i-- {
		for h.nexts[i] != nil && h.nexts[i].key < k {
			h = h.nexts[i]
		}

		// 当前高度 <= 新节点插入高度时, 将新节点插入对应位置
		if i <= high {
			_node.nexts[i] = h.nexts[i]
			h.nexts[i] = _node
		}
	}

}

func (s *SkList[K, E]) Remove(k K) {
	// 不存在,直接返回
	if n := s.search(k); n == nil {
		return
	}

	h := s.head

	for i := len(h.nexts) - 1; i >= 0; i-- {
		for h.nexts[i] != nil && h.nexts[i].key < k {
			h = h.nexts[i]
		}

		if h.nexts[i] == nil || h.nexts[i].key > k {
			continue
		}

		// 走到这里必然找到了 k, 调整指针
		h.nexts[i] = h.nexts[i].nexts[i]
	}

	// 调整跳表的高度, 因为删除后,可能出现某些高度一个节点都没有了
	dif := 0
	for i := len(s.head.nexts) - 1; i >= 0 && s.head.nexts[i] == nil; i-- {
		dif++
	}

	s.head.nexts = s.head.nexts[:len(s.head.nexts)-dif]
}

// All 整个跳表的迭代器
func (s *SkList[K, E]) All() iter.Seq2[K, E] {

	return func(yield func(K, E) bool) {
		if s.empty() {
			return
		}

		sn := s.head.nexts[0]
		if sn == nil {
			return
		}

		for n := sn; n.nexts[0] != nil; n = n.nexts[0] {
			if !yield(n.key, n.val) {
				return
			}
		}
	}
}

func (s *SkList[K, E]) Range(start, end K) iter.Seq2[K, E] {

	return func(yield func(K, E) bool) {
		if s.empty() {
			return
		}

		sn := s.ceiling(start)

		if sn == nil {
			return
		}

		for n := sn; n.nexts[0] != nil && n.key <= end; n = n.nexts[0] {
			if !yield(n.key, n.val) {
				return
			}
		}
	}
}

// ceiling 找到第一个 >= t的节点
func (s *SkList[K, E]) ceiling(t K) *nd[K, E] {
	if s.empty() {
		return nil
	}

	n := s.head
	for i, v := range slices.Backward(s.head.nexts) {
		for v != nil && v.key < t {
			v = v.nexts[i]
			n = v
		}

		if v != nil && v.key == t {
			return v
		}
	}

	return n.nexts[0]
}

// floor 找到最后一个 <= t的节点
func (s *SkList[K, E]) floor(t K) *nd[K, E] {
	n := s.head
	for i, v := range slices.Backward(s.head.nexts) {
		for v != nil && v.key < t {
			n = v
			v = v.nexts[i]
		}

		if v != nil && v.key == t {
			return v
		}
	}

	if n == s.head {
		return nil
	}

	return n
}

// search 核心查找方法
func (s *SkList[K, E]) search(k K) *nd[K, E] {
	if s.empty() {
		return nil
	}

	h := s.head
	for i := len(h.nexts) - 1; i >= 0; i-- {
		// 查找当前层,遍历到当前层第一个 >=k 的位置停下
		for h.nexts[i] != nil && h.nexts[i].key < k {
			h = h.nexts[i]
		}

		// 如果key 相等,找到对应节点, 返回
		if h.nexts[i] != nil && h.nexts[i].key == k {
			return h.nexts[i]
		}

		// 否则下降到下一层继续查找, 即i--操作
	}
	return nil
}

func (s *SkList[K, E]) roll() int {
	high := 0
	for rand.IntN(2) > 0 {
		high++
	}
	return high
}

func (s *SkList[K, E]) empty() bool {
	return s == nil || s.head == nil || len(s.head.nexts) == 0
}
