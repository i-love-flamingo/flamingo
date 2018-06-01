package cache

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/groupcache/singleflight"
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	// HTTPLoader returns a response. it will be cached unless there is an error. this means 400/500 responses are cached too!
	HTTPLoader func() (*http.Response, *Meta, error)

	// HTTPFrontend stores and caches http responses
	HTTPFrontend struct {
		singleflight.Group
		Backend Backend         `inject:""`
		Logger  flamingo.Logger `inject:""`
	}

	nopCloser struct {
		io.Reader
	}

	cachedResponse struct {
		orig *http.Response
		body []byte
	}
)

func GetHTTPFrontendCacheWithNullBackend() *HTTPFrontend {
	return &HTTPFrontend{
		Backend: &NullBackend{},
		Logger:  flamingo.NullLogger{},
	}
}

// Close the nopCloser to implement io.Closer
func (nopCloser) Close() error { return nil }

func copyResponse(response cachedResponse, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	var newResponse http.Response
	if response.orig != nil {
		newResponse = *response.orig
	}

	buf := make([]byte, len(response.body))
	copy(buf, response.body)
	newResponse.Body = nopCloser{bytes.NewBuffer(buf)}

	return &newResponse, nil
}

// Get a http response, with tags and a loader
// the tags will be used when the entry is stored
func (hf *HTTPFrontend) Get(key string, loader HTTPLoader) (*http.Response, error) {
	if hf.Backend == nil {
		return nil, errors.New("NO Backend in Cache")
	}
	if entry, ok := hf.Backend.Get(key); ok {
		if entry.Meta.lifetime.After(time.Now()) {
			hf.Logger.WithField("category", "httpFrontendCache").Debugf("Serving from cache %v", key)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}

		if entry.Meta.gracetime.After(time.Now()) {
			go hf.load(key, loader)
			hf.Logger.WithField("category", "httpFrontendCache").Debugf("Gracetime! Serving from cache %v", key)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}
	}
	hf.Logger.WithField("category", "httpFrontendCache").Debugf("No cache entry for %v", key)
	return copyResponse(hf.load(key, loader))
}

func (hf *HTTPFrontend) load(key string, loader HTTPLoader) (cachedResponse, error) {
	data, err := hf.Do(key, func() (res interface{}, resultErr error) {
		defer func() {
			if err := recover(); err != nil {
				//resultErr = errors.WithStack(fmt.Errorf("%#v", err))
				if err2, ok := err.(error); ok {
					resultErr = errors.WithStack(err2) //fmt.Errorf("%#v", err)
				} else {
					resultErr = errors.WithStack(fmt.Errorf("HTTPFrontend.load exception: %#v", err))
				}
			}
		}()

		data, meta, err := loader()
		if meta == nil {
			meta = &Meta{
				Lifetime:  30 * time.Second,
				Gracetime: 10 * time.Minute,
			}
		}
		if err != nil {
			return loaderResponse{nil, meta}, err
		}

		response := data
		body, _ := ioutil.ReadAll(response.Body)

		response.Body.Close()

		cached := cachedResponse{
			orig: response,
			body: body,
		}

		return loaderResponse{cached, meta}, err
	})

	//if err != nil {
	//	if hf.Logger != nil {
	//		hf.Logger.Error("cache load failed: ", err)
	//	}
	//	return cachedResponse{}, err
	//}

	if data == nil {
		data = loaderResponse{
			cachedResponse{
				orig: new(http.Response),
				body: []byte{},
			},
			&Meta{
				Lifetime:  30 * time.Second,
				Gracetime: 10 * time.Minute,
			},
		}
	}

	loadedData := data.(loaderResponse).data
	var cached cachedResponse
	if loadedData != nil {
		cached = loadedData.(cachedResponse)
	}

	hf.Logger.WithField("category", "httpFrontendCache").Debugf("Store in Cache %v / Meta: %#v", key, data.(loaderResponse).meta)
	hf.Backend.Set(key, &Entry{
		Data: cached,
		Meta: Meta{
			lifetime:  time.Now().Add(data.(loaderResponse).meta.Lifetime),
			gracetime: time.Now().Add(data.(loaderResponse).meta.Lifetime + data.(loaderResponse).meta.Gracetime),
			Tags:      data.(loaderResponse).meta.Tags,
		},
	})

	return cached, err
}
