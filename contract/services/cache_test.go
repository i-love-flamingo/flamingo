package services_test

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/core/cache"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

const cacheReaders = 10

func TestHTTPFrontend_LoadsConcurrentMissOnce(t *testing.T) {
	t.Parallel()

	// created outside the bubble, because its cleanup goroutine never ends
	backend := cache.NewInMemoryCache()

	synctest.Test(t, func(t *testing.T) {
		frontend := new(cache.HTTPFrontend).Inject(backend, flamingo.NullLogger{})

		var (
			loads   atomic.Int32
			release = make(chan struct{})
			bodies  = make([]string, cacheReaders)
			wg      sync.WaitGroup
		)

		var loader cache.HTTPLoader = func(context.Context) (*http.Response, *cache.Meta, error) {
			// each call gives a body with its own number
			body := "brand-" + strconv.Itoa(int(loads.Add(1)))

			<-release

			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))},
				&cache.Meta{Lifetime: time.Hour}, nil
		}

		for reader := range cacheReaders {
			wg.Add(1)

			go func() {
				defer wg.Done()

				resp, err := frontend.Get(t.Context(), "brand", loader)
				if !assert.NoError(t, err) {
					return
				}

				defer func() { _ = resp.Body.Close() }()

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				bodies[reader] = string(body)
			}()
		}

		synctest.Wait() // every reader is blocked, in the loader or behind it
		close(release)
		wg.Wait()

		require.EqualValues(t, 1, loads.Load(), "loader calls")

		for reader, body := range bodies {
			assert.Equal(t, "brand-1", body, "reader %d", reader)
		}
	})
}
