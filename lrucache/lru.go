//go:build !solution

package lrucache

import "container/list"

func New(cap int) Cache {
	return NewCacheImpl(cap)
}

type entry struct {
	key   int
	value int
}

type CacheImpl struct {
	capacity int
	data     map[int]*list.Element
	history  *list.List
}

func NewCacheImpl(capacity int) *CacheImpl {
	return &CacheImpl{
		capacity: capacity,
		data:     make(map[int]*list.Element, capacity),
		history:  list.New(),
	}
}

func (c *CacheImpl) Get(key int) (int, bool) {
	if elem, found := c.data[key]; found {
		c.history.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return 0, false
}

func (c *CacheImpl) Set(key, value int) {
	if c.capacity <= 0 {
		return
	}

	if elem, found := c.data[key]; found {
		elem.Value.(*entry).value = value
		c.history.MoveToFront(elem)
	} else {
		if c.history.Len() == c.capacity {
			c.evict()
		}
		elem := c.history.PushFront(&entry{key: key, value: value})
		c.data[key] = elem
	}
}

func (c *CacheImpl) evict() {
	lru := c.history.Back()
	if lru != nil {
		c.history.Remove(lru)
		delete(c.data, lru.Value.(*entry).key)
	}
}

func (c *CacheImpl) Range(f func(key, value int) bool) {
	for elem := c.history.Back(); elem != nil; elem = elem.Prev() {
		kv := elem.Value.(*entry)
		if !f(kv.key, kv.value) {
			break
		}
	}
}

func (c *CacheImpl) Clear() {
	c.history.Init()
	c.data = make(map[int]*list.Element)
}
