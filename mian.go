package main

import (
	"fmt"
	"log"
	"net/http"
	tc "tinycache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	tc.NewGroup("scores", 2<<10, tc.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[slowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}

			return nil, fmt.Errorf("%s not exists", key)
		}))

	addr := "localhost:4399"
	peers := tc.NewHTTPPool(addr)
	log.Println("tiny cache is running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
