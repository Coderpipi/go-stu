package lru

import "container/list"

type (
	// Cache LRU Cache
	Cache struct {
		maxBytes int64
		nBytes   int64
		ll       *list.List
		cache    map[string]*list.Element

		// 当被淘汰时执行的策略
		OnEvicted func(key string, value Value)
	}

	entry struct {
		key   string
		value Value
	}

	Value interface {
		Len() int
	}
)

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return nil, false
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(len(kv.key) + kv.value.Len())
		kv.value = value
	} else {
		ele = c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key) + value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
