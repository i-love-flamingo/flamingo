package cache

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

// Backend is an autogenerated mock type for the Backend type
type MockBackend struct {
	mock.Mock
}

// Flush provides a mock function with given fields:
func (_m *MockBackend) Flush() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *MockBackend) Get(key string) (*Entry, bool) {
	ret := _m.Called(key)

	var r0 *Entry
	if rf, ok := ret.Get(0).(func(string) *Entry); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Entry)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// Purge provides a mock function with given fields: key
func (_m *MockBackend) Purge(key string) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PurgeTags provides a mock function with given fields: tags
func (_m *MockBackend) PurgeTags(tags []string) error {
	ret := _m.Called(tags)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Set provides a mock function with given fields: key, entry
func (_m *MockBackend) Set(key string, entry *Entry) error {
	ret := _m.Called(key, entry)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *Entry) error); ok {
		r0 = rf(key, entry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func createLoader(statusCode int, body string, err error) func(ctx context.Context) (*http.Response, *Meta, error) {
	return func(ctx context.Context) (*http.Response, *Meta, error) {
		return createResponse(statusCode, body), nil, err
	}
}

func createResponse(statusCode int, body string) *http.Response {
	response, _ := http.ReadResponse(
		bufio.NewReader(
			strings.NewReader(
				fmt.Sprintf("HTTP/1.1 %d OK\nContent-Type: text/html\n\n%s", statusCode, body),
			),
		),
		nil,
	)
	return response
}

func loaderWithWatingTime(ctx context.Context) (*http.Response, *Meta, error) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write([]byte("Test 123"))
	}))

	defer server.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, nil, err
}

func TestHTTPFrontend_Get(t *testing.T) {
	// wait channel top check async cache setting
	cacheSetComplete := make(chan struct{}, 1)
	inFuture := time.Now().Add(time.Hour)
	inPast := time.Now().Add(-time.Hour)

	type args struct {
		key    string
		loader HTTPLoader
	}
	tests := []struct {
		name string
		args args
		// Entry returned by Backend Get
		cacheEntry *Entry
		// expected response from HTTPFrontend
		want    *http.Response
		wantErr bool
		// wanted data to be passed to Backend Set
		wantedCachedData []byte
		// want Backend Set to be called
		wantSet bool
		ctx     context.Context
	}{
		{
			name: "empty cache",
			args: args{
				key:    "test",
				loader: createLoader(200, "body", nil),
			},
			cacheEntry:       nil,
			want:             createResponse(200, "body"),
			wantedCachedData: []byte("body"),
			wantSet:          true,
			wantErr:          false,
		},
		{
			name: "empty cache, error on loader",
			args: args{
				key:    "test",
				loader: createLoader(200, "body", errors.New("test error on loader")),
			},
			cacheEntry: nil,
			want:       nil,
			// even in error case of the loader, the result is expected to be cached (easy circuit breaker)
			wantedCachedData: nil,
			wantSet:          true,
			wantErr:          true,
		},
		{
			name: "data from cache in lifetime",
			args: args{
				key:    "test",
				loader: nil,
			},
			cacheEntry: &Entry{
				Meta: Meta{
					lifetime:  inFuture,
					gracetime: inFuture,
				},
				Data: cachedResponse{
					orig: createResponse(200, "foo"),
					body: []byte("foo"),
				},
			},
			want:             createResponse(200, "foo"),
			wantedCachedData: nil,
			wantSet:          false,
			wantErr:          false,
		},
		{
			name: "from cache out of lifetime but in gracetime",
			args: args{
				key:    "test",
				loader: createLoader(200, "body", nil),
			},
			cacheEntry: &Entry{
				Meta: Meta{
					lifetime:  inPast,
					gracetime: inFuture,
				},
				Data: cachedResponse{
					orig: createResponse(200, "foo"),
					body: []byte("foo"),
				},
			},
			// we expect the cached value as result, but a the actual value from loader to be cached (async)
			want:             createResponse(200, "foo"),
			wantedCachedData: []byte("body"),
			wantSet:          true,
			wantErr:          false,
		},
		{
			name: "from cache out of lifetime and out of gracetime",
			args: args{
				key:    "test",
				loader: createLoader(200, "body", nil),
			},
			cacheEntry: &Entry{
				Meta: Meta{
					lifetime:  inPast,
					gracetime: inPast,
				},
				Data: cachedResponse{
					orig: createResponse(200, "foo"),
					body: []byte("foo"),
				},
			},
			// cached value is invalid, so the actual value from loader is expected as response and to be cached
			want:             createResponse(200, "body"),
			wantedCachedData: []byte("body"),
			wantSet:          true,
			wantErr:          false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backendMock := &MockBackend{}
			backendMock.On("Get", tt.args.key).Return(tt.cacheEntry, tt.cacheEntry != nil).Once()

			setCall := backendMock.On(
				"Set",
				tt.args.key,
				mock.MatchedBy(func(e *Entry) bool {
					return assert.Equal(t, e.Data.(cachedResponse).body, tt.wantedCachedData)
				}),
			).Run(func(args mock.Arguments) {
				cacheSetComplete <- struct{}{}
			}).Return(nil)
			if tt.wantSet {
				setCall.Once()
			} else {
				setCall.Maybe()
				// no cache set expected, so we complete directly
				cacheSetComplete <- struct{}{}
			}

			hf := new(HTTPFrontend).Inject(
				backendMock,
				&flamingo.NullLogger{},
			)

			got, err := hf.Get(context.Background(), tt.args.key, tt.args.loader)
			// wait for event. async Backend Set completion before asserting expectations
			<-cacheSetComplete
			backendMock.AssertExpectations(t)

			assert.Equalf(t, err != nil, tt.wantErr, "Get() error = %v, wantErr %v", err, tt.wantErr)
			if tt.wantErr {
				assert.Nil(t, got, "response is expected to be nil in error case")
				return
			}
			require.NotNil(t, got, "result of Get() is nil")
			assert.Equal(t, got.Header, tt.want.Header)
			gotBody, _ := io.ReadAll(got.Body)
			wantBody, _ := io.ReadAll(tt.want.Body)
			assert.Equal(t, string(wantBody), string(gotBody))
		})
	}
}

//nolint:bodyclose // response might be nil so we cannot close the body
func TestContextDeadlineExceeded(t *testing.T) {
	t.Parallel()

	t.Run("exceeded, throw error", func(t *testing.T) {
		t.Parallel()

		entry := &Entry{
			Meta: Meta{
				lifetime:  time.Now().Add(-24 * time.Hour),
				gracetime: time.Now().Add(-24 * time.Hour),
			},
			Data: nil,
		}

		backendMock := &MockBackend{}
		backendMock.On("Get", "test").Return(entry, true)
		backendMock.On("Set", "test", mock.Anything).Return(func(string, *Entry) error {
			return nil
		})

		contextWithDeadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
		t.Cleanup(cancel)

		hf := new(HTTPFrontend).Inject(
			backendMock,
			&flamingo.NullLogger{},
		)

		got, err := hf.Get(contextWithDeadline, "test", loaderWithWatingTime)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Nil(t, got)
	})

	t.Run("did not exceed, no error", func(t *testing.T) {
		t.Parallel()

		entry := &Entry{
			Meta: Meta{
				lifetime:  time.Now().Add(-24 * time.Hour),
				gracetime: time.Now().Add(-24 * time.Hour),
			},
			Data: nil,
		}

		backendMock := &MockBackend{}
		backendMock.On("Get", "test").Return(entry, true)
		backendMock.On("Set", "test", mock.Anything).Return(func(string, *Entry) error {
			return nil
		})

		contextWithDeadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(7*time.Second))
		t.Cleanup(cancel)

		hf := new(HTTPFrontend).Inject(
			backendMock,
			&flamingo.NullLogger{},
		)

		got, err := hf.Get(contextWithDeadline, "test", loaderWithWatingTime)

		assert.NoError(t, err)
		assert.Equal(t, got.StatusCode, http.StatusOK)
	})
}
