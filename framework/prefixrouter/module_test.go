package prefixrouter_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
)

type testingNullLogger struct{}

func (m *testingNullLogger) Configure(injector *dingo.Injector) {
	injector.Bind(new(flamingo.Logger)).To(flamingo.NullLogger{})
}

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(framework.InitModule), new(testingNullLogger), new(prefixrouter.Module)); err != nil {
		t.Error(err)
	}
}
