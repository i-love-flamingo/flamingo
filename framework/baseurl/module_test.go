package baseurl_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/baseurl"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(baseurl.Module)); err != nil {
		t.Error(err)
	}
}
