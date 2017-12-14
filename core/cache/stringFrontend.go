package cache

import (
	"time"

	"github.com/golang/groupcache/singleflight"
)

type (
	StringLoader   func() (string, *Meta, error)
	StringFrontend struct {
		singleflight.Group
		Backend Backend `inject:""`
	}
)

// Get string cache
func (sf *StringFrontend) Get(key string, loader StringLoader) (string, error) {
	if entry, ok := sf.Backend.Get(key); ok {
		if entry.Meta.lifetime.After(time.Now()) {
			return entry.Data.(string), nil
		}

		if entry.Meta.gracetime.After(time.Now()) {
			go sf.load(key, loader)
			return entry.Data.(string), nil
		}
	}

	return sf.load(key, loader)
}

func (sf *StringFrontend) load(key string, loader StringLoader) (string, error) {
	data, err := sf.Do(key, func() (interface{}, error) {
		data, meta, err := loader()
		if meta == nil {
			meta = &Meta{
				Lifetime:  30 * time.Second,
				Gracetime: 10 * time.Minute,
			}
		}
		return loaderResponse{data, meta}, err
	})

	if err != nil {
		return "", err
	}

	sf.Backend.Set(key, &Entry{
		Data: data.(loaderResponse).data,
		Meta: Meta{
			lifetime:  time.Now().Add(data.(loaderResponse).meta.Lifetime),
			gracetime: time.Now().Add(data.(loaderResponse).meta.Lifetime + data.(loaderResponse).meta.Gracetime),
			Tags:      data.(loaderResponse).meta.Tags,
		},
	})

	return data.(string), nil
}
