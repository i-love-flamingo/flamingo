package flamingo

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang/groupcache/singleflight"
	"github.com/hashicorp/golang-lru"
)

const (
	lifetime     = 30 * time.Second
	gracetime    = 10 * time.Minute
	lurkerPeriod = 10 * time.Second
	cachesize    = 100
)

type (
	Cache interface {
		Get(ctx context.Context, key string, loader func() (interface{}, error)) (interface{}, bool)
	}

	HttpCache struct {
		Cache `inject:""`
	}

	NullCache struct{}

	inMemoryCacheEntry struct {
		valid time.Time
		data  interface{}
	}

	inMemoryCache struct {
		singleflight.Group
		sync.RWMutex
		pool *lru.TwoQueueCache
	}

	nopCloser struct {
		io.Reader
	}

	cachedResponse struct {
		orig *http.Response
		body []byte
	}
)

func (nopCloser) Close() error { return nil }

func (*NullCache) Get(ctx context.Context, key string, loader func() (interface{}, error)) (interface{}, bool) {
	r, err := loader()
	return r, err == nil
}

func (hc *HttpCache) Get(ctx context.Context, key string, loader func() (*http.Response, error)) (*http.Response, bool) {
	cached, ok := hc.Cache.Get(ctx, key, func() (interface{}, error) {
		resp, err := loader()
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		resp.Body.Close()

		return cachedResponse{
			orig: resp,
			body: body,
		}, nil
	})

	response := &(*cached.(cachedResponse).orig)
	response.Body = nopCloser{bytes.NewBuffer(cached.(cachedResponse).body)}

	return response, ok
}

func NewInMemoryCache() *inMemoryCache {
	cache, _ := lru.New2Q(cachesize)

	m := &inMemoryCache{
		pool: cache,
	}
	go m.lurker()
	return m
}

func (m *inMemoryCache) Get(ctx context.Context, key string, loader func() (interface{}, error)) (interface{}, bool) {
	m.RLock()

	if data, ok := m.pool.Get(key); ok {
		if time.Now().Before(data.(inMemoryCacheEntry).valid) {
			m.RUnlock()
			log.Println("got cache for", key)
			return data.(inMemoryCacheEntry).data, true
		}

		// gracetime?
		if time.Now().Before(data.(inMemoryCacheEntry).valid.Add(gracetime)) {
			m.RUnlock()

			log.Println("scheduling recache for", key)
			// schedule a refresh
			go m.load(context.Background(), key, loader)

			return data.(inMemoryCacheEntry).data, true
		}
	}

	m.RUnlock()

	return m.load(context.Background(), key, loader)
}

func (m *inMemoryCache) load(ctx context.Context, key string, loader func() (interface{}, error)) (interface{}, bool) {
	data, err := m.Do(key, loader)

	if err != nil {
		log.Println("error loading cache", err)
		return nil, false
	}

	m.Lock()
	defer m.Unlock()

	if entry, ok := m.pool.Get(key); ok {
		if time.Now().Before(entry.(inMemoryCacheEntry).valid) {
			log.Println("entry already valid", key, entry.(inMemoryCacheEntry).valid)
			return entry.(inMemoryCacheEntry).data, true
		}
	}

	log.Println("saved cache for", key)

	m.pool.Add(key, inMemoryCacheEntry{
		valid: time.Now().Add(lifetime),
		data:  data,
	})

	return data, true
}

func (m *inMemoryCache) lurker() {
	for tick := range time.Tick(lurkerPeriod) {
		log.Println("cache lurker", tick)
		m.Lock()

		for _, key := range m.pool.Keys() {
			item, ok := m.pool.Peek(key)
			if ok && time.Now().After(item.(inMemoryCacheEntry).valid.Add(gracetime)) {
				m.pool.Remove(key)
				log.Println("cache lurker", "cleared", key)
				break
			}
		}

		m.Unlock()
	}
}
