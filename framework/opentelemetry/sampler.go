package opentelemetry

import (
	"net/http"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type ConfiguredURLPrefixSampler struct {
	Whitelist        config.Slice
	Blacklist        config.Slice
	AllowParentTrace bool
}

// Inject dependencies
func (c *ConfiguredURLPrefixSampler) Inject(
	cfg *struct {
		Whitelist        config.Slice `inject:"config:flamingo.opentelemetry.tracing.sampler.whitelist,optional"`
		Blacklist        config.Slice `inject:"config:flamingo.opentelemetry.tracing.sampler.blacklist,optional"`
		AllowParentTrace bool         `inject:"config:flamingo.opentelemetry.tracing.sampler.allowParentTrace,optional"`
	},
) *ConfiguredURLPrefixSampler {
	if cfg != nil {
		c.Whitelist = cfg.Whitelist
		c.Blacklist = cfg.Blacklist
		c.AllowParentTrace = cfg.AllowParentTrace
	}
	return c
}

func (c *ConfiguredURLPrefixSampler) GetFilterOption() otelhttp.Filter {
	var allowed, blocked []string
	_ = c.Whitelist.MapInto(&allowed)
	_ = c.Blacklist.MapInto(&blocked)

	return URLPrefixSampler(allowed, blocked, c.AllowParentTrace)
}

func URLPrefixSampler(allowed, blocked []string, allowParentTrace bool) otelhttp.Filter {
	return func(request *http.Request) bool {
		path := request.URL.Path
		isParentSampled := trace.SpanContextFromContext(request.Context()).IsSampled()
		// empty allowed means all
		sample := len(allowed) == 0
		// check allowed if len is > 0, and decide if we should sample
		for _, p := range allowed {
			if strings.HasPrefix(path, p) {
				sample = true
				break
			}
		}

		// we do not sample, unless the parent is sampled
		if !sample {
			return !allowParentTrace && isParentSampled
		}

		// check sampling decision against blocked
		for _, p := range blocked {
			if strings.HasPrefix(path, p) {
				sample = false
				break
			}
		}

		// we sample, or the parent sampled
		return (!allowParentTrace && isParentSampled) || sample
	}
}
