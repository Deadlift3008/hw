package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	currentItem, ok := c.items[key]

	if ok {
		currentItem.Value.(*cacheItem).value = value

		c.queue.MoveToFront(currentItem)

		return true
	}

	newItem := c.queue.PushFront(&cacheItem{key, value})

	if c.queue.Len() > c.capacity {
		last := c.queue.Back()

		c.queue.Remove(last)
		delete(c.items, last.Value.(*cacheItem).key)
	}

	c.items[key] = newItem

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	currentItem, ok := c.items[key]

	if ok {
		c.queue.MoveToFront(currentItem)

		return currentItem.Value.(*cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
