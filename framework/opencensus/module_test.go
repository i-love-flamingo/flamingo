package opencensus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{
		"flamingo.opencensus.jaeger.enable": true,
		"flamingo.opencensus.zipkin.enable": true,
	}, new(Module)); err != nil {
		t.Error(err)
	}

	assert.NotNil(t, new(Module).FlamingoLegacyConfigAlias())
}

func TestView(t *testing.T) {
	assert.NoError(t, View("test", stats.Int64("testm", "testd", "tests"), view.Count()))
	assert.Error(t, View("test", stats.Int64("testm2", "testd2", "tests2"), view.Count()))
}

type testRoundTripper struct{}

func (testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{Header: req.Header}, nil
}

func TestCorrelationIDInjector(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.Background()
	ctx, span := trace.StartSpan(ctx, "test")
	req = req.WithContext(ctx)
	resp, err := (&correlationIDInjector{next: testRoundTripper{}}).RoundTrip(req)
	assert.NoError(t, err)
	assert.Equal(t, span.SpanContext().TraceID.String(), resp.Header.Get("X-Correlation-ID"))
}

func TestLocalAddr(t *testing.T) {
	//TODO: how do we test this?
	localAddr()
}
