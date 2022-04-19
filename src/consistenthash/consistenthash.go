package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func(data []byte) uint32

type Map struct {
	hash     HashFunc
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn HashFunc) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add 插入节点到一致性哈希结构中
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// strconv.Itoa(num) ---> string + string
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 因为是一个环形的代码所以需要取余数 因为可能回出现 idx == len(m.keys)
	// 因为如果hash>= keys[len - 1] idx = len(m.keys) ---> 变为0节点那个
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
