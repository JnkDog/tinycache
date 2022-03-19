package src

import (
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"tom":  "630",
	"jack": "589",
	"sam":  "567",
}

func TestGet(t *testing.T) {
	loadCount := make(map[string]int, len(db))
	tiny := NewGroup("test", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCount[key]; !ok {
					loadCount[key] = 0
				}
				loadCount[key] += 1
				return []byte(v), nil
			}

			return nil, fmt.Errorf("%s not exists", key)
		}))

	for k, v := range db {
		if view, err := tiny.Get(k); err != nil || view.String() != v {
			t.Fatalf("test fail")
		}
		if _, err := tiny.Get(k); err != nil || loadCount[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := tiny.Get("unknown"); err == nil {
		t.Fatalf("the value should nil. but %s got", view)
	}
}
