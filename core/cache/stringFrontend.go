package cache

import (
	"time"

	"github.com/golang/groupcache/singleflight"
	"go.opencensus.io/trace"
)

type (
	// StringLoader is used to load strings for singleflight cache loads
	StringLoader func() (string, *Meta, error)

	// StringFrontend manages cache entries as strings
	StringFrontend struct {
		singleflight.Group
		backend Backend
	}
)

// Inject StringFrontend dependencies
func (sf *StringFrontend) Inject(backend Backend) {
	sf.backend = backend
}

// Get and load string cache entries
func (sf *StringFrontend) Get(key string, loader StringLoader) (string, error) {
	if entry, ok := sf.backend.Get(key); ok {
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
		return loaderResponse{data, meta, trace.SpanContext{}}, err
	})

	if err != nil {
		return "", err
	}

	sf.backend.Set(key, &Entry{
		Data: data.(loaderResponse).data,
		Meta: Meta{
			lifetime:  time.Now().Add(data.(loaderResponse).meta.Lifetime),
			gracetime: time.Now().Add(data.(loaderResponse).meta.Lifetime + data.(loaderResponse).meta.Gracetime),
			Tags:      data.(loaderResponse).meta.Tags,
		},
	})

	return data.(string), nil
}
