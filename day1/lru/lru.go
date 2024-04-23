package lru

import (
	"container/list"
)

type Catch struct {
	// 允许的最大内存
	maxByte int64
	// 当前已经使用的内存
	nbytes int64
	// 定义双向链表
	ll *list.List
	// map[string]*list.Element，键是字符串，值是双向链表中对应节点的指针
	cache map[string]*list.Element
	// 被清除时执行
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxByte int64, onEvicted func(string, Value)) *Catch {
	return &Catch{
		maxByte:   maxByte,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Catch) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 淘汰最少访问的节点
func (c *Catch) RemoveOldest() {
	// 获取队首节点表示最少访问
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 被清除时执行
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Catch) Add(key string, value Value) {
	// 如果存在需要得到原来的
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(value.Len()) + int64(len(key))
	}
	if c.maxByte != 0 && c.maxByte < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Catch) Len() int {
	return c.ll.Len()
}
