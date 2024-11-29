package skip_list

import (
	"math/rand"
)

type SkipList struct {
	head *node
}

type node struct {
	nexts    []*node
	key, val int
}

// Get 读操作
func (s *SkipList) Get(key int) (int, bool) {
	if _node := s.search(key); _node != nil {
		return _node.val, true
	}

	return -1, false
}

// Put 写操作
func (s *SkipList) Put(key, val int) {
	// 如果key已经存在, 直接更新
	if _node := s.search(key); _node != nil {
		_node.val = val
		return
	}

	// roll出新节点的高度
	level := s.roll()

	// 新节点高度大于当前跳表最大高度, 则进行补齐
	for len(s.head.nexts)-1 < level {
		s.head.nexts = append(s.head.nexts, nil)
	}

	// 创建新节点
	newNode := &node{
		key:   key,
		val:   val,
		nexts: make([]*node, level+1),
	}

	// 从头结点的最高层出发
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		// 向右遍历, 直到右侧节点不存在或者key 值 > key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		// 插入节点到对应位置
		newNode.nexts[level] = move.nexts[level]
		move.nexts[level] = newNode
	}
}

// Del 删除操作
func (s *SkipList) Del(key int) {
	// 如果不存在, 直接返回
	if _node := s.search(key); _node == nil {
		return
	}

	// 从头结点的最高层出发
	move := s.head

	for level := len(s.head.nexts) - 1; level >= 0; level++ {
		// 向右遍历, 直到右侧节点不存在或者key 值 > key
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}

		if move.nexts[level] == nil || move.nexts[level].key > key {
			continue
		}

		// 走到这代表key必然相等, 且不能直接break, 因为需要对每一层进行调整指针的操作
		move.nexts[level] = move.nexts[level].nexts[level]
	}

	// 对跳表的高度进行更新
	var dif int
	// 倘若某一层已经不存在数据节点, 高度递减
	for level := len(s.head.nexts) - 1; level > 0 && s.head.nexts[level] != nil; level-- {
		dif++
	}
	s.head.nexts = s.head.nexts[:len(s.head.nexts)-dif]
}

// Range 范围查询
// 找到skip list当中 >= start, 且<= end的kv对
func (s *SkipList) Range(start, end int) [][2]int {
	// 首先通过ceiling方法, 找到skip list中key值大于等于start 且最接近于start的节点
	ceilNode := s.ceiling(start)
	// 如果不存在,直接返回
	if ceilNode == nil {
		return [][2]int{}
	}

	// 从ceilNode第一层向右遍历, 把所有位于[start, end]区间的节点返回
	var res [][2]int
	for move := ceilNode; move != nil && move.key <= end; move = move.nexts[0] {
		res = append(res, [2]int{move.key, move.val})
	}

	return res
}

// Ceiling 找到skip list中key值大于等于start 且最接近于target的节点
func (s *SkipList) Ceiling(target int) ([2]int, bool) {
	if ceilNode := s.ceiling(target); ceilNode != nil {
		return [2]int{ceilNode.key, ceilNode.val}, true
	}

	return [2]int{}, false
}

// Floor 找到skip list中key值<=target 且最接近于target的节点
func (s *SkipList) Floor(target int) ([2]int, bool) {
	if floorNode := s.floor(target); s.head.key != floorNode.key {
		return [2]int{floorNode.key, floorNode.val}, true
	}

	return [2]int{}, false
}

// ceiling 找到skip list中key值大于等于target 且最接近于target的节点
func (s *SkipList) ceiling(target int) *node {
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}

		// 如果key值 >= target的kv对存在, 返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 走到这里代表走到最底层都没有找到等于target的kv对, 且停在<key且最接近去key的节点, 其右侧的第一个节点肯定是>且最接近于target的节点
	return move.nexts[0]
}

// floor 找到skip list中key值<=target 且最接近于target的节点
func (s *SkipList) floor(target int) *node {
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < target {
			move = move.nexts[level]
		}

		// 如果key值 >= target的kv对存在, 返回
		if move.nexts[level] != nil && move.nexts[level].key == target {
			return move.nexts[level]
		}
	}

	// 走到这里代表走到最底层都没有找到等于target的kv对, 且停在<key且最接近去key的节点
	return move
}

// roll 决定新节点在跳表中的高度
func (s *SkipList) roll() int {
	var level int
	for rand.Intn(2) > 0 {
		level++
	}
	return level
}

// search 从跳表中检索key 对应的node
func (s *SkipList) search(key int) *node {
	// 每次检索从头部出发
	move := s.head
	for level := len(s.head.nexts) - 1; level >= 0; level-- {
		for move.nexts[level] != nil && move.nexts[level].key < key {
			move = move.nexts[level]
		}
		// 如果key值相等, 则找到了目标直接返回
		if move.nexts[level] != nil && move.nexts[level].key == key {
			return move.nexts[level]
		}

		// 否则进入下一层, 继续遍历
	}

	return nil
}
