package lru

import (
	"reflect"
	"testing"
)

type String string

// go test 有没有before这种方法
//var c *Cache
func (d String) Len() int {
	return len(d)
}

func TestNew(t *testing.T) {
	//c := new(Cache)
	//types.AssertableTo(&Cache{}, c)
}

func TestCache_Add_if_not_full_and_get(t *testing.T) {
	//c = &Cache{}
	c := New(5, nil)
	c.Add("r", String("111"))
	if key, ok := c.Get("r"); ok {
		t.Logf("%s found", key)
	} else {
		t.Fatalf("%s found", key)
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	size := len(k1 + k2 + v1 + v2)
	lru := New(int64(size), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest faild")
	}
}

func TestEvictedCallback(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))
	expect := []string{"key1", "k2"}
	// 比较
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Callback test fail")
	}
}
