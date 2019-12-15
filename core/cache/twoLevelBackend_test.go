package cache_test

import (
	"flamingo.me/flamingo/v3/framework/flamingo"
	"testing"

	"flamingo.me/flamingo/v3/core/cache"
)

func Test_RunDefaultBackendTestCase_TwoLevelBackend(t *testing.T) {
	f := cache.TwoLevelBackendFactory{}
	c := cache.TwoLevelBackendConfig{
		FirstLevel:  cache.NewInMemoryCache(),
		SecondLevel: cache.NewFileBackend("", "mutlilevelBackendTest"),
	}

	backend := f.Inject(flamingo.NullLogger{}).SetConfig(c).Build()

	testcase := cache.NewBackendTestCase(t, backend, true)
	testcase.RunTests()
}
