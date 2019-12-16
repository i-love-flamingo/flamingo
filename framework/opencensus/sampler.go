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
func URLPrefixSampler(whitelist, blacklist []string, allowParentTrace bool) func(*http.Request) trace.StartOptions {
	return func(request *http.Request) trace.StartOptions {

		path := request.URL.Path

		// empty whitelist means all
		sample := len(whitelist) == 0
		// check whitelist if len is > 0, and decide if we should sample
		for _, p := range whitelist {
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

		// check sampling decision against blacklist
		for _, p := range blacklist {
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
	Whitelist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.whitelist,optional"`
	Blacklist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.blacklist,optional"`
	AllowParentTrace bool         `inject:"config:flamingo.opencensus.tracing.sampler.allowParentTrace,optional"`
}

// GetStartOptions constructor for ochttp.Server
func (c *ConfiguredURLPrefixSampler) GetStartOptions() func(*http.Request) trace.StartOptions {
	var whitelist, blacklist []string
	c.Whitelist.MapInto(&whitelist)
	c.Blacklist.MapInto(&blacklist)

	return URLPrefixSampler(whitelist, blacklist, c.AllowParentTrace)
}
