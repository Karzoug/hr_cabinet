package ttlmap

import (
	"sync"
	"time"
)

type ttlmap[K comparable, V any] struct {
	sync.RWMutex
	items map[K]item[V]
	close chan struct{}
}

type item[T any] struct {
	data      T
	expiresAt int64
}

// New creates key-value storage. Simple, but thread-safe.
func New[K comparable, V any](cleaningInterval time.Duration) *ttlmap[K, V] {
	ttlmap := &ttlmap[K, V]{
		items:   make(map[K]item[V]),
		RWMutex: sync.RWMutex{},
		close:   make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(cleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ttlmap.clean()
			case <-ttlmap.close:
				return
			}
		}
	}()

	return ttlmap
}

// Set sets value by key.
func (tm *ttlmap[K, V]) Set(key K, value V, duration time.Duration) error {
	tm.Lock()
	defer tm.Unlock()

	var expires int64
	if duration > 0 {
		expires = time.Now().Add(duration).UnixNano()
	}
	tm.items[key] = item[V]{
		data:      value,
		expiresAt: expires,
	}

	return nil
}

// Get returns value by key.
func (tm *ttlmap[K, V]) Get(key K) (V, bool) {
	tm.Lock()
	defer tm.Unlock()

	item, exists := tm.items[key]
	if !exists {
		return item.data, false
	}
	if item.expiresAt > 0 && time.Now().UnixNano() > item.expiresAt {
		delete(tm.items, key)
		return item.data, false
	}

	return item.data, true
}

// Delete deletes value by key.
func (ttlmap *ttlmap[K, V]) Delete(key K) error {
	delete(ttlmap.items, key)
	return nil
}

// Close closes storage: cleans internal map and stop cleaning work.
func (tm *ttlmap[K, V]) Close() error {
	tm.close <- struct{}{}

	tm.Lock()
	defer tm.Unlock()

	tm.items = make(map[K]item[V])
	return nil
}

func (tm *ttlmap[K, V]) clean() {
	now := time.Now().UnixNano()

	tm.Lock()
	defer tm.Unlock()

	for k, v := range tm.items {
		if now > v.expiresAt && v.expiresAt > 0 {
			delete(tm.items, k)
		}
	}
}
