package opencensus

import (
	"net/http"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"go.opencensus.io/trace"
)

// URLPrefixSampler creates a sampling getter for ochttp.Server.
//
// If the whitelist is empty it is treated as allowed, otherwise checked first.
// If the blacklist is set it will disable sampling again.
// If takeParentDecision is set we allow the decision to be taken from incoming tracing,
// otherwise we enforce our decision
func URLPrefixSampler(allowed, blocked []string, allowParentTrace bool) func(*http.Request) trace.StartOptions {
	return func(request *http.Request) trace.StartOptions {

		path := request.URL.Path

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
			return trace.StartOptions{
				Sampler: func(p trace.SamplingParameters) trace.SamplingDecision {
					return trace.SamplingDecision{Sample: !allowParentTrace && p.ParentContext.IsSampled()}
				},
			}
		}

		// check sampling decision against blocked
		for _, p := range blocked {
			if strings.HasPrefix(path, p) {
				sample = false
				break
			}
		}

		// we sample, or the parent sampled
		return trace.StartOptions{
			Sampler: func(p trace.SamplingParameters) trace.SamplingDecision {
				return trace.SamplingDecision{Sample: (!allowParentTrace && p.ParentContext.IsSampled()) || sample}
			},
		}
	}
}

// ConfiguredURLPrefixSampler constructs the prefix GetStartOptions getter with the default opencensus configuration
type ConfiguredURLPrefixSampler struct {
	Whitelist        config.Slice
	Blacklist        config.Slice
	AllowParentTrace bool
}

// Inject dependencies
func (c *ConfiguredURLPrefixSampler) Inject(
	cfg *struct {
		Whitelist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.whitelist,optional"`
		Blacklist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.blacklist,optional"`
		AllowParentTrace bool         `inject:"config:flamingo.opencensus.tracing.sampler.allowParentTrace,optional"`
	},
) *ConfiguredURLPrefixSampler {
	if cfg != nil {
		c.Whitelist = cfg.Whitelist
		c.Blacklist = cfg.Blacklist
		c.AllowParentTrace = cfg.AllowParentTrace
	}
	return c
}

// GetStartOptions constructor for ochttp.Server
func (c *ConfiguredURLPrefixSampler) GetStartOptions() func(*http.Request) trace.StartOptions {
	var allowed, blocked []string
	_ = c.Whitelist.MapInto(&allowed)
	_ = c.Blacklist.MapInto(&blocked)

	return URLPrefixSampler(allowed, blocked, c.AllowParentTrace)
}
