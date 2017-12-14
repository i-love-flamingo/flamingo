package cache

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/groupcache/singleflight"
)

type (
	// HTTPLoader returns a response. it will be cached unless there is an error. this means 400/500 responses are cached too!
	HTTPLoader func() (*http.Response, *Meta, error)

	// HTTPFrontend stores and caches http responses
	HTTPFrontend struct {
		singleflight.Group
		Backend Backend `inject:""`
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
func (hf *HTTPFrontend) Get(key string, loader HTTPLoader) (*http.Response, error) {
	if entry, ok := hf.Backend.Get(key); ok {
		if entry.Meta.lifetime.After(time.Now()) {
			return copyResponse(entry.Data.(cachedResponse), nil)
		}

		if entry.Meta.gracetime.After(time.Now()) {
			go hf.load(key, loader)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}
	}

	return copyResponse(hf.load(key, loader))
}

func (hf *HTTPFrontend) load(key string, loader HTTPLoader) (cachedResponse, error) {
	data, err := hf.Do(key, func() (interface{}, error) {
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
		return cachedResponse{}, err
	}

	response := data.(loaderResponse).data.(*http.Response)
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	cached := cachedResponse{
		orig: response,
		body: body,
	}

	hf.Backend.Set(key, &Entry{
		Data: cached,
		Meta: Meta{
			lifetime:  time.Now().Add(data.(loaderResponse).meta.Lifetime),
			gracetime: time.Now().Add(data.(loaderResponse).meta.Lifetime + data.(loaderResponse).meta.Gracetime),
			Tags:      data.(loaderResponse).meta.Tags,
		},
	})

	return cached, nil
}
