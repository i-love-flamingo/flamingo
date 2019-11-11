package cache_test

// @TODO: write unit-tests for all exported methods

import (
	"testing"

	"flamingo.me/flamingo/v3/core/cache"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

func Test_RunDefaultBackendTestCase_MultiLevelBackend(t *testing.T) {
	backend1 := cache.NewInMemoryCache()
	backend2 := cache.NewFileBackend("")
	backend := cache.NewMultiLevelBackend(
		cache.MultiLevelBackendOptions{
			Backends: []cache.Backend{backend1, backend2},
		},
		flamingo.NullLogger{},
	)

	testcase := cache.NewBackendTestCase(t, backend, true)
	testcase.RunTests()
}
