package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	// dummy node
	// [2:2 4:4 6:6 12:2 14:4 16:6 22:2 24:4 26:6]
	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}
	// 赛一个节点
	hash.Add("8")
	testCases["27"] = "8"
	t.Log(testCases)
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yield %s", k, v)
		}
	}
}
