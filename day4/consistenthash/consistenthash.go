package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 实现一致性哈希算法

type Hash func(data []byte) uint32

// Map Hash 函数 hash；
type Map struct {
	hash Hash
	// 虚拟节点倍数 replicas
	replicas int
	// 哈希环 keys
	keys []int
	// 虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称
	hashMap map[int]string
}

// New 允许自定义虚拟节点倍数和 Hash 函数
func New(replicas int, fu Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fu,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		// 采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，默认为 crc32.ChecksumIEEE 算法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	// 对每一个真实哈希节点创建 m.replicas 个虚拟节点，
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 记录 i + key 所对应的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 添加到哈希环上
			m.keys = append(m.keys, hash)
			// 在 hashMap 存储虚拟节点和真实节点的对应关系
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
