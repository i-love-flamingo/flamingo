package systemendpoint_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(systemendpoint.Module)); err != nil {
		t.Error(err)
	}
}
