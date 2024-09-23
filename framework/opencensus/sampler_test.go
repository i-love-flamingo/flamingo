package opencensus_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/opencensus"
)

func TestURLPrefixSampler(t *testing.T) {
	sampler := new(opencensus.ConfiguredURLPrefixSampler)
	sampler.Inject(
		&config.Area{},
		&struct {
			Whitelist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.whitelist,optional"`
			Blacklist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.blacklist,optional"`
			AllowParentTrace bool         `inject:"config:flamingo.opencensus.tracing.sampler.allowParentTrace,optional"`
		}{})
	assert.NotNil(t, sampler.GetStartOptions())

	assert.True(t, opencensus.URLPrefixSampler(nil, nil, false)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{}).Sample)
	assert.True(t, opencensus.URLPrefixSampler(nil, nil, true)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{}).Sample)

	assert.False(t, opencensus.URLPrefixSampler([]string{"/test"}, nil, true)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{}).Sample)
	assert.False(t, opencensus.URLPrefixSampler([]string{"/test"}, nil, true)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{ParentContext: trace.SpanContext{TraceOptions: trace.TraceOptions(1)}}).Sample)

	assert.False(t, opencensus.URLPrefixSampler([]string{"/test"}, nil, true)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{}).Sample)
	assert.True(t, opencensus.URLPrefixSampler([]string{"/test"}, nil, false)(httptest.NewRequest(http.MethodGet, "/", nil)).Sampler(trace.SamplingParameters{ParentContext: trace.SpanContext{TraceOptions: trace.TraceOptions(1)}}).Sample)

	assert.True(t, opencensus.URLPrefixSampler([]string{"/test"}, nil, false)(httptest.NewRequest(http.MethodGet, "/test", nil)).Sampler(trace.SamplingParameters{}).Sample)

	assert.True(t, opencensus.URLPrefixSampler([]string{"/test"}, []string{"/test/foo"}, false)(httptest.NewRequest(http.MethodGet, "/test", nil)).Sampler(trace.SamplingParameters{}).Sample)
	assert.False(t, opencensus.URLPrefixSampler([]string{"/test"}, []string{"/test/foo"}, false)(httptest.NewRequest(http.MethodGet, "/test/foo", nil)).Sampler(trace.SamplingParameters{}).Sample)
}
