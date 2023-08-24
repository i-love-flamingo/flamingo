package cache_test

import (
	"testing"
	"time"

	"flamingo.me/flamingo/v3/core/cache"
	"github.com/stretchr/testify/assert"
)

func Test_inMemoryCache_Flush(t *testing.T) {
	t.Parallel()

	inMemoryCache := cache.NewInMemoryCache()

	err := inMemoryCache.Set("foo", &cache.Entry{
		Meta: cache.Meta{
			Tags:      []string{"bar"},
			Lifetime:  5 * time.Second,
			Gracetime: 10 * time.Second,
		},
		Data: "test",
	})

	assert.NoError(t, err)

	err = inMemoryCache.Set("baz", &cache.Entry{
		Meta: cache.Meta{
			Tags:      []string{"bar"},
			Lifetime:  5 * time.Second,
			Gracetime: 10 * time.Second,
		},
		Data: "test",
	})

	assert.NoError(t, err)

	_, found := inMemoryCache.Get("foo")
	assert.True(t, found)

	assert.NoError(t, inMemoryCache.Flush())

	_, found = inMemoryCache.Get("foo")
	assert.False(t, found)

	_, found = inMemoryCache.Get("baz")
	assert.False(t, found)
}

func Test_inMemoryCache_Purge(t *testing.T) {
	t.Parallel()

	inMemoryCache := cache.NewInMemoryCache()

	err := inMemoryCache.Set("foo", &cache.Entry{
		Meta: cache.Meta{
			Tags:      []string{"bar"},
			Lifetime:  5 * time.Second,
			Gracetime: 10 * time.Second,
		},
		Data: "test",
	})

	assert.NoError(t, err)

	_, found := inMemoryCache.Get("foo")
	assert.True(t, found)

	assert.NoError(t, inMemoryCache.Purge("foo"))

	_, found = inMemoryCache.Get("foo")
	assert.False(t, found)
}

func Test_inMemoryCache_SetGet(t *testing.T) {
	t.Parallel()

	inMemoryCache := cache.NewInMemoryCache()

	err := inMemoryCache.Set("foo", &cache.Entry{
		Meta: cache.Meta{
			Tags:      []string{"bar"},
			Lifetime:  5 * time.Second,
			Gracetime: 10 * time.Second,
		},
		Data: "test",
	})

	assert.NoError(t, err)

	entry, found := inMemoryCache.Get("foo")
	assert.True(t, found)
	assert.Equal(t, "test", entry.Data)
	assert.Equal(t, cache.Meta{
		Tags:      []string{"bar"},
		Lifetime:  5 * time.Second,
		Gracetime: 10 * time.Second,
	}, entry.Meta)
}
