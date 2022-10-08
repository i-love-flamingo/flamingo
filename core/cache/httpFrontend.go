package cache

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/golang/groupcache/singleflight"
	"go.opencensus.io/trace"
)

type (
	// HTTPLoader returns a response. it will be cached unless there is an error. this means 400/500 responses are cached too!
	HTTPLoader func(context.Context) (*http.Response, *Meta, error)

	// HTTPFrontend stores and caches http responses
	HTTPFrontend struct {
		singleflight.Group
		backend Backend
		logger  flamingo.Logger
	}

	nopCloser struct {
		io.Reader
	}

	cachedResponse struct {
		orig *http.Response
		body []byte
	}
)

// Inject HTTPFrontend dependencies
func (hf *HTTPFrontend) Inject(backend Backend, logger flamingo.Logger) *HTTPFrontend {
	hf.backend = backend
	hf.logger = logger

	return hf
}

// GetHTTPFrontendCacheWithNullBackend helper for tests
func GetHTTPFrontendCacheWithNullBackend() *HTTPFrontend {
	return &HTTPFrontend{
		backend: &NullBackend{},
		logger:  flamingo.NullLogger{},
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
func (hf *HTTPFrontend) Get(ctx context.Context, key string, loader HTTPLoader) (*http.Response, error) {
	if hf.backend == nil {
		return nil, errors.New("NO backend in Cache")
	}

	ctx, span := trace.StartSpan(ctx, "flamingo/cache/httpFrontend/Get")
	span.Annotate(nil, key)
	defer span.End()

	if entry, ok := hf.backend.Get(key); ok {
		if entry.Meta.lifetime.After(time.Now()) {
			hf.logger.WithField("category", "httpFrontendCache").Debug("Serving from cache", key)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}

		if entry.Meta.gracetime.After(time.Now()) {
			go hf.load(context.Background(), key, loader, true)
			hf.logger.WithField("category", "httpFrontendCache").Debug("Gracetime! Serving from cache", key)
			return copyResponse(entry.Data.(cachedResponse), nil)
		}
	}
	hf.logger.WithField("category", "httpFrontendCache").Debug("No cache entry for", key)

	return copyResponse(hf.load(ctx, key, loader, false))
}

func (hf *HTTPFrontend) load(ctx context.Context, key string, loader HTTPLoader, keepExistingEntry bool) (cachedResponse, error) {
	ctx, span := trace.StartSpan(ctx, "flamingo/cache/httpFrontend/load")
	span.Annotate(nil, key)
	defer span.End()

	data, err := hf.Do(key, func() (res interface{}, resultErr error) {
		ctx, fetchRoutineSpan := trace.StartSpan(context.Background(), "flamingo/cache/httpFrontend/fetchRoutine")
		fetchRoutineSpan.Annotate(nil, key)
		defer fetchRoutineSpan.End()

		defer func() {
			if err := recover(); err != nil {
				if err2, ok := err.(error); ok {
					resultErr = fmt.Errorf("httpfrontend load: %w", err2)
				} else {
					resultErr = fmt.Errorf("httpfrontend load: %v", err2)
				}
			}
		}()

		data, meta, err := loader(ctx)
		if meta == nil {
			meta = &Meta{
				Lifetime:  30 * time.Second,
				Gracetime: 10 * time.Minute,
			}
		}
		if err != nil {
			return loaderResponse{nil, meta, fetchRoutineSpan.SpanContext()}, err
		}

		response := data
		body, _ := io.ReadAll(response.Body)

		response.Body.Close()

		cached := cachedResponse{
			orig: response,
			body: body,
		}

		return loaderResponse{cached, meta, fetchRoutineSpan.SpanContext()}, err
	})

	keepExistingEntry = keepExistingEntry && (err != nil || data == nil)

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
			trace.SpanContext{},
		}
	}

	loadedData := data.(loaderResponse).data
	var cached cachedResponse
	if loadedData != nil {
		cached = loadedData.(cachedResponse)
	}

	if keepExistingEntry {
		hf.logger.WithContext(ctx).WithField("category", "httpFrontendCache").Debug("No store/overwrite in cache because we couldn't fetch new data", key)
	} else {
		hf.logger.WithContext(ctx).WithField("category", "httpFrontendCache").Debug("Store in Cache", key, data.(loaderResponse).meta)
		hf.backend.Set(key, &Entry{
			Data: cached,
			Meta: Meta{
				lifetime:  time.Now().Add(data.(loaderResponse).meta.Lifetime),
				gracetime: time.Now().Add(data.(loaderResponse).meta.Lifetime + data.(loaderResponse).meta.Gracetime),
				Tags:      data.(loaderResponse).meta.Tags,
			},
		})
	}

	span.AddAttributes(trace.StringAttribute("parenttrace", data.(loaderResponse).span.TraceID.String()))
	span.AddAttributes(trace.StringAttribute("parentspan", data.(loaderResponse).span.SpanID.String()))
	//span.AddLink(trace.Link{
	//	SpanID:  data.(loaderResponse).span.SpanID,
	//	TraceID: data.(loaderResponse).span.TraceID,
	//	Type:    trace.LinkTypeChild,
	//})

	return cached, err
}
