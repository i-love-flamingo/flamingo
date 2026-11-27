package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	"go.uber.org/zap/zaptest/observer"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/contract/internal/harness"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// brandModule binds two routes that log the message "brand".
// The route "/brand/traced" logs inside a span and gives the IDs of the span.
type brandModule struct{}

func (*brandModule) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(brandRoutes))
}

type brandRoutes struct {
	responder *web.Responder
	logger    flamingo.Logger
}

func (r *brandRoutes) Inject(responder *web.Responder, logger flamingo.Logger) {
	r.responder = responder
	r.logger = logger
}

func (r *brandRoutes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/brand", "brand")
	registry.HandleGet("brand", func(ctx context.Context, _ *web.Request) web.Result {
		r.logger.WithContext(ctx).Info("brand")

		return r.responder.Data(map[string]string{})
	})

	registry.MustRoute("/brand/traced", "brand.traced")
	registry.HandleGet("brand.traced", func(ctx context.Context, _ *web.Request) web.Result {
		ctx, span := trace.StartSpan(ctx, "brand")
		defer span.End()

		r.logger.WithContext(ctx).Info("brand")

		return r.responder.Data(map[string]string{
			"trace": span.SpanContext().TraceID.String(),
			"span":  span.SpanContext().SpanID.String(),
		})
	})
}

// brandEntry gives the log entry with the message "brand" and its fields.
func brandEntry(tb testing.TB, app *harness.App) map[string]any {
	tb.Helper()

	var found []observer.LoggedEntry

	for _, entry := range app.Logs() {
		if entry.Message == "brand" {
			found = append(found, entry)
		}
	}

	require.Len(tb, found, 1, "log entries with the message \"brand\"")

	return found[0].ContextMap()
}

// Log searches of applications depend on these field names.
func TestLogger_WithContextAddsRequestAndTrace(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
	t.Run("request", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		app := harness.Boot(t, []dingo.Module{new(brandModule)}, harness.Config{})

		app.Get(t, "/brand")

		fields := brandEntry(t, app)

		assert.Equal(t, "GET", fields["method"])
		assert.Equal(t, "/brand", fields["path"])
	})

	t.Run("span", func(t *testing.T) { //nolint:paralleltest // boots an app, dingo scopes are process-global
		app := harness.Boot(t, []dingo.Module{new(brandModule)}, harness.Config{})

		body := app.Get(t, "/brand/traced").Body

		fields := brandEntry(t, app)

		// the router starts a span for each request, the entry has the span of the route
		assert.Equal(t, "/brand/traced", fields["path"])
		assert.Contains(t, body, `"trace":"`+fieldString(t, fields, "traceID")+`"`)
		assert.Contains(t, body, `"span":"`+fieldString(t, fields, "spanID")+`"`)
	})
}

func fieldString(tb testing.TB, fields map[string]any, key string) string {
	tb.Helper()

	value, ok := fields[key].(string)
	require.True(tb, ok, "field %q is a string, fields: %v", key, fields)
	require.NotEmpty(tb, value, "field %q", key)

	return value
}
