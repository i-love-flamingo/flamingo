package opencensus

import (
	"net/http"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"go.opencensus.io/trace"
)

// URLPrefixSampler creates a sampling getter for ochttp.Server.
//
// If the include list is empty it is treated as allowed, otherwise checked first.
// If the exclude list is set it will disable sampling again.
// If takeParentDecision is set we allow the decision to be taken from incoming tracing,
// otherwise we enforce our decision
func URLPrefixSampler(include, exclude []string, allowParentTrace bool) func(*http.Request) trace.StartOptions {
	return func(request *http.Request) trace.StartOptions {

		path := request.URL.Path

		// empty include means all
		sample := len(include) == 0
		// check include if len is > 0, and decide if we should sample
		for _, p := range include {
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

		// check sampling decision against exclude
		for _, p := range exclude {
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
	Include          config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.include,optional"`
	Exclude          config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.exclude,optional"`
	AllowParentTrace bool         `inject:"config:flamingo.opencensus.tracing.sampler.allowParentTrace,optional"`
}

// GetStartOptions constructor for ochttp.Server
func (c *ConfiguredURLPrefixSampler) GetStartOptions() func(*http.Request) trace.StartOptions {
	var include, exclude []string
	c.Include.MapInto(&include)
	c.Exclude.MapInto(&exclude)

	return URLPrefixSampler(include, exclude, c.AllowParentTrace)
}
