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
	// 检查cache是否存在数据
	if ele, ok := c.cache[key]; ok {
		// 移动到链表首部
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest 淘汰最少访问的节点
func (c *Catch) RemoveOldest() {
	// 获取队尾节点
	ele := c.ll.Back()
	if ele != nil {
		// 移除链表节点
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 删除cache数据
		delete(c.cache, kv.key)
		// 去除缓存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 被清除时执行
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Catch) Add(key string, value Value) {
	// 检查缓存是否存在
	if ele, ok := c.cache[key]; ok {
		// 移动到链表首部
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		// 更新数据
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 不存在，放在链表首部
		ele := c.ll.PushFront(&entry{key, value})
		// 存cache
		c.cache[key] = ele
		// 增加缓存值
		c.nbytes += int64(value.Len()) + int64(len(key))
	}
	if c.maxByte != 0 && c.maxByte < c.nbytes {
		// 移除链表尾部节点
		c.RemoveOldest()
	}
}

func (c *Catch) Len() int {
	return c.ll.Len()
}
