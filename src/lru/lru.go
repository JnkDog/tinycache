package lru

import "container/list"

type Cache struct {
	doubleLinkedList *list
}
