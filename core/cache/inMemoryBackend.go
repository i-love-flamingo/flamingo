package cache

import (
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const lurkerPeriod = 1 * time.Minute

type (
	inMemoryBackend struct {
		cacheMetrics Metrics
		pool         *lru.TwoQueueCache
	}

	//InMemoryBackendConfig config
	InMemoryBackendConfig struct {
		Size int
	}

	//InMemoryBackendFactory factory
	InMemoryBackendFactory struct {
		config       InMemoryBackendConfig
		frontendName string
	}

	inMemoryCacheEntry struct {
		valid time.Time
		data  interface{}
	}
)

// NewInMemoryCache creates a new lru TwoQueue backed cache backend
// Depricated - use the InMemoryBackendFactory or general CacheFactory
func NewInMemoryCache() Backend {
	f := InMemoryBackendFactory{}
	return f.SetConfig(InMemoryBackendConfig{
		Size: 100}).SetFrontendName("default").Build()
}

//SetConfig for factory
func (f *InMemoryBackendFactory) SetConfig(config InMemoryBackendConfig) *InMemoryBackendFactory {
	f.config = config
	return f
}

//SetFrontendName used in Metrics
func (f *InMemoryBackendFactory) SetFrontendName(frontendName string) *InMemoryBackendFactory {
	f.frontendName = frontendName
	return f
}

//Build factory func
func (f *InMemoryBackendFactory) Build() Backend {
	cache, _ := lru.New2Q(f.config.Size)

	m := &inMemoryBackend{
		pool:         cache,
		cacheMetrics: NewCacheMetrics("inMemory", f.frontendName),
	}
	go m.lurker()
	return m
}

// Get tries to get an object from cache
func (m *inMemoryBackend) SetSize(size int) error {
	cache, err := lru.New2Q(100)
	if err != nil {
		return err
	}
	m.pool = cache
	return nil
}

// Get tries to get an object from cache
func (m *inMemoryBackend) Get(key string) (*Entry, bool) {
	entry, ok := m.pool.Get(key)
	if !ok {
		m.cacheMetrics.countMiss()
		return nil, ok
	}
	m.cacheMetrics.countHit()
	return entry.(inMemoryCacheEntry).data.(*Entry), ok
}

// Set a cache entry with a key
func (m *inMemoryBackend) Set(key string, entry *Entry) error {
	m.pool.Add(key, inMemoryCacheEntry{
		data:  entry,
		valid: entry.Meta.gracetime,
	})

	return nil
}

// Purge a cache key
func (m *inMemoryBackend) Purge(key string) error {
	m.pool.Remove(key)

	return nil
}

// Flush purges all entries in the cache
func (m *inMemoryBackend) Flush() error {
	m.pool.Purge()

	return nil
}

func (m *inMemoryBackend) lurker() {
	for range time.Tick(lurkerPeriod) {
		for _, key := range m.pool.Keys() {
			item, ok := m.pool.Peek(key)
			if ok && item.(inMemoryCacheEntry).valid.Before(time.Now()) {
				m.pool.Remove(key)
				break
			}
		}
	}
}
