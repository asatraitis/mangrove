package utils

import "sync"

type CacheMap[K comparable, V any] map[K]V
type Cache[K comparable, V any] interface {
	GetValue(K) V
	SetValue(K, V)
}
type cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data CacheMap[K, V]
}

func NewCache[K comparable, V any]() Cache[K, V] {
	return &cache[K, V]{
		data: make(CacheMap[K, V]),
	}
}

func (c *cache[K, V]) GetValue(key K) V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}
func (c *cache[K, V]) SetValue(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
