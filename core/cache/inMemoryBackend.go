package cache

import (
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

const lurkerPeriod = 1 * time.Minute

type (
	inMemoryCache struct {
		pool *lru.TwoQueueCache[string, inMemoryCacheEntry]
	}

	inMemoryCacheEntry struct {
		valid time.Time
		data  interface{}
	}
)

// NewInMemoryCache creates a new lru TwoQueue backed cache backend
func NewInMemoryCache() Backend {
	cache, _ := lru.New2Q[string, inMemoryCacheEntry](100)

	m := &inMemoryCache{
		pool: cache,
	}
	go m.lurker()
	return m
}

// Get tries to get an object from cache
func (m *inMemoryCache) Get(key string) (*Entry, bool) {
	entry, ok := m.pool.Get(key)
	if !ok {
		return nil, ok
	}
	return entry.data.(*Entry), ok
}

// Set a cache entry with a key
func (m *inMemoryCache) Set(key string, entry *Entry) error {
	m.pool.Add(key, inMemoryCacheEntry{
		data:  entry,
		valid: entry.Meta.gracetime,
	})

	return nil
}

// Purge a cache key
func (m *inMemoryCache) Purge(key string) error {
	m.pool.Remove(key)

	return nil
}

// PurgeTags purges all entries with matching tags from the cache
func (m *inMemoryCache) PurgeTags(_ []string) error {
	return errors.New("not implemented")
}

// Flush purges all entries in the cache
func (m *inMemoryCache) Flush() error {
	m.pool.Purge()

	return nil
}

func (m *inMemoryCache) lurker() {
	for range time.Tick(lurkerPeriod) {
		for _, key := range m.pool.Keys() {
			item, ok := m.pool.Peek(key)
			if ok && item.valid.Before(time.Now()) {
				m.pool.Remove(key)
				break
			}
		}
	}
}
