package cache

import (
	"time"

	"github.com/golang/groupcache/singleflight"
)

type (
	StringLoader   func() (string, error)
	StringFrontend struct {
		singleflight.Group
		Backend Backend
	}
)

func (sf *StringFrontend) Get(key string, lifetime, gracetime time.Duration, loader StringLoader, tags ...string) (string, error) {
	if entry, ok := sf.Backend.Get(key); ok {
		if entry.Lifetime.After(time.Now()) {
			return entry.Data.(string), nil
		}

		if entry.Lifetime.Add(gracetime).After(time.Now()) {
			go sf.load(key, lifetime, gracetime, loader, tags...)
			return entry.Data.(string), nil
		}
	}

	return sf.load(key, lifetime, gracetime, loader, tags...)
}

func (sf *StringFrontend) load(key string, lifetime, gracetime time.Duration, loader StringLoader, tags ...string) (string, error) {
	data, err := sf.Do(key, func() (interface{}, error) {
		return loader()
	})

	if err != nil {
		return "", err
	}

	sf.Backend.Set(key, &CacheEntry{
		Data:      data,
		Lifetime:  time.Now().Add(lifetime),
		Gracetime: time.Now().Add(lifetime + gracetime),
		Tags:      tags,
	})

	return data.(string), nil
}
