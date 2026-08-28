package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mutex   sync.RWMutex
	entries map[string]Entry
	stop    chan struct{}
}

type Entry struct {
	Value     []byte
	Duration  time.Duration
	CreatedAt time.Time
}

func New(interval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]Entry),
		stop:    make(chan struct{}),
	}

	go cache.start(interval)

	return cache
}

func (cache *Cache) start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cache.clean()
		case <-cache.stop:
			return
		}
	}
}

func (cache *Cache) clean() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, entry := range cache.entries {
		if time.Since(entry.CreatedAt) >= entry.Duration {
			delete(cache.entries, key)
		}
	}
}

func (cache *Cache) Stop() {
	close(cache.stop)
}

func (cache *Cache) Remember(key string, resolver func() []byte, duration time.Duration) []byte {
	cache.mutex.RLock()
	entry, ok := cache.entries[key]
	cache.mutex.RUnlock()

	if ok && time.Since(entry.CreatedAt) < entry.Duration {
		return entry.Value
	}

	value := resolver()

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[key] = Entry{
		Value:     value,
		Duration:  duration,
		CreatedAt: time.Now(),
	}

	return value
}

func (cache *Cache) Forget(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.entries, key)
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.RLock()
	entry, ok := cache.entries[key]
	cache.mutex.RUnlock()

	if !ok {
		return nil, false
	}

	if time.Since(entry.CreatedAt) >= entry.Duration {
		cache.Forget(key)
		return nil, false
	}

	return entry.Value, true
}
