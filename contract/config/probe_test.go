package config_test

import (
	"context"
	"encoding/json"
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/web"
)

// probeRoutes serves the values of T at "/probe" as JSON.
// The area that serves the request fills T through the "inject" tags of its fields.
type probeRoutes[T any] struct {
	values *T
}

// Inject the config values.
func (r *probeRoutes[T]) Inject(values *T) {
	r.values = values
}

// Routes of the probe.
func (r *probeRoutes[T]) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/probe", "probe")
	registry.HandleGet("probe", func(_ context.Context, _ *web.Request) web.Result {
		return new(web.Responder).Data(r.values)
	})
}

// probeModule binds the probe routes for T. Use it to read the values of a child area.
func probeModule[T any]() dingo.Module {
	return dingo.ModuleFunc(func(injector *dingo.Injector) {
		web.BindRoutes(injector, new(probeRoutes[T]))
	})
}

// rootValues gives the values that dingo injects in the root area.
// It asks the root injector, because root has no path prefix in these tests.
func rootValues[T any](tb testing.TB, app *harness.App) T {
	tb.Helper()

	instance, err := app.Injector(tb).GetInstance(new(T))
	if err != nil {
		tb.Fatalf("inject root values: %v", err)
	}

	values, ok := instance.(*T)
	if !ok {
		tb.Fatalf("inject root values: got %T", instance)
	}

	return *values
}

// areaValues reads the values that the probe of the area at the path prefix sees.
func areaValues[T any](tb testing.TB, app *harness.App, prefix string) T {
	tb.Helper()

	body := app.Get(tb, prefix+"/probe").Body

	var values T

	if err := json.Unmarshal([]byte(body), &values); err != nil {
		tb.Fatalf("probe %q: decode %q: %v", prefix, body, err)
	}

	return values
}
