package prefixrouter_test

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"testing"
)

type testingNullLogger struct{}

func (m *testingNullLogger) Configure(injector *dingo.Injector) {
	injector.Bind(new(flamingo.Logger)).To(flamingo.NullLogger{})
}

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(testingNullLogger), new(prefixrouter.Module)); err != nil {
		t.Error(err)
	}
}
