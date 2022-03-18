package lru

import (
	"container/list"
)

type Cache struct {
	maxBytes         int64
	usedBytes        int64
	doubleLinkedList *list.List
	cache            map[string]*list.Element

	// 回调函数，当淘汰发生的时候
	EvictedCallback func(key string, value Value)
}

// 存放在list Element中的value
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, callback func(string, Value)) *Cache {
	return &Cache{
		maxBytes:         maxBytes,
		doubleLinkedList: list.New(),
		cache:            make(map[string]*list.Element),
		EvictedCallback:  callback,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.doubleLinkedList.MoveToFront(element)
		// 强制类型转换
		kv := element.Value.(*entry)
		return kv.value, true
	}

	return
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		c.doubleLinkedList.MoveToFront(element)
		kv := element.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
	} else {
		element := c.doubleLinkedList.PushFront(&entry{key, value})
		c.cache[key] = element
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	element := c.doubleLinkedList.Back()
	if element != nil {
		c.doubleLinkedList.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.EvictedCallback != nil {
			c.EvictedCallback(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int {
	return c.doubleLinkedList.Len()
}
