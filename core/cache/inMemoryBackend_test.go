package cache_test

import (
	"testing"

	"flamingo.me/flamingo/v3/core/cache"
)

func Test_RunDefaultBackendTestCase_InMemoryBackend(t *testing.T) {
	backend := cache.NewInMemoryCache()

	testcase := cache.NewBackendTestCase(t, backend, true)
	testcase.RunTests()
}
