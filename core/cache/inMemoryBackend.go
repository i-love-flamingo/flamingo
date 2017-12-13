package cache

import (
	"log"
	"time"

	"github.com/hashicorp/golang-lru"
)

const lurkerPeriod = 1 * time.Minute

type (
	inMemoryCache struct {
		pool *lru.TwoQueueCache
	}

	inMemoryCacheEntry struct {
		valid time.Time
		data  interface{}
	}
)

func NewInMemoryCache() *inMemoryCache {
	cache, _ := lru.New2Q(100)

	m := &inMemoryCache{
		pool: cache,
	}
	go m.lurker()
	return m
}

func (m *inMemoryCache) Get(key string) (*CacheEntry, bool) {
	entry, ok := m.pool.Get(key)
	if !ok {
		return nil, ok
	}
	return entry.(inMemoryCacheEntry).data.(*CacheEntry), ok
}

func (m *inMemoryCache) Set(key string, entry *CacheEntry) {
	m.pool.Add(key, inMemoryCacheEntry{
		data:  entry,
		valid: entry.Gracetime,
	})
}

func (m *inMemoryCache) Purge(key string) {
	m.pool.Remove(key)
}

func (m *inMemoryCache) PurgeTags(tags []string) {
	panic("implement me")
}

func (m *inMemoryCache) Flush() {
	m.pool.Purge()
}

func (m *inMemoryCache) lurker() {
	for tick := range time.Tick(lurkerPeriod) {
		for _, key := range m.pool.Keys() {
			item, ok := m.pool.Peek(key)
			if ok && item.(inMemoryCacheEntry).valid.Before(time.Now()) {
				m.pool.Remove(key)
				log.Println("cache lurker", tick, "cleared", key)
				break
			}
		}
	}
}
