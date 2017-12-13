package cache

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/groupcache/singleflight"
)

type (
	// HTTPLoader returns a response. it will be cached unless there is an error. this means 400/500 responses are cached too!
	HTTPLoader func() (*http.Response, error)

	// HTTPFrontend stores and caches http responses
	HTTPFrontend struct {
		singleflight.Group
		Backend Backend
	}

	nopCloser struct {
		io.Reader
	}

	cachedResponse struct {
		orig *http.Response
		body []byte
	}
)

// Close the nopCloser to implement io.Closer
func (nopCloser) Close() error { return nil }

func copyResponse(response cachedResponse, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}

	newResponse := *response.orig

	buf := make([]byte, len(response.body))
	copy(buf, response.body)
	newResponse.Body = nopCloser{bytes.NewBuffer(buf)}

	return &newResponse, nil
}

// Get a http response, with tags and a loader
// the tags will be used when the entry is stored
func (hf *HTTPFrontend) Get(key string, lifetime, gracetime time.Duration, loader HTTPLoader, tags ...string) (*http.Response, error) {
	log.Println("Trying to cache")

	if entry, ok := hf.Backend.Get(key); ok {
		if entry.Lifetime.After(time.Now()) {
			return copyResponse(entry.Data.(cachedResponse), nil)
		}

		if entry.Lifetime.Add(gracetime).After(time.Now()) {
			go hf.load(key, lifetime, gracetime, loader, tags...)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}
	}

	return copyResponse(hf.load(key, lifetime, gracetime, loader, tags...))
}

func (hf *HTTPFrontend) load(key string, lifetime, gracetime time.Duration, loader HTTPLoader, tags ...string) (cachedResponse, error) {
	data, err := hf.Do(key, func() (interface{}, error) {
		return loader()
	})

	log.Println("loaded cache", data, err)

	if err != nil {
		return cachedResponse{}, err
	}

	response := data.(*http.Response)
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	cached := cachedResponse{
		orig: response,
		body: body,
	}

	hf.Backend.Set(key, &CacheEntry{
		Data:      cached,
		Lifetime:  time.Now().Add(lifetime),
		Gracetime: time.Now().Add(lifetime + gracetime),
		Tags:      tags,
	})

	return cached, nil
}
