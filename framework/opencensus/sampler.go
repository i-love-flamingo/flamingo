package opencensus

import (
	"net/http"
	"path"
	"strings"

	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/config"
)

// URLPrefixSampler creates a sampling getter for ochttp.Server.
//
// If the whitelist is empty it is treated as allowed, otherwise checked first.
// If the blacklist is set it will disable sampling again.
// If ignoreParentDecision is set we allow the decision to be taken from incoming tracing,
// otherwise we enforce our decision
func URLPrefixSampler(allowed, blocked []string, ignoreParentDecision bool) func(*http.Request) trace.StartOptions {
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
					return trace.SamplingDecision{Sample: !ignoreParentDecision && p.ParentContext.IsSampled()}
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
				return trace.SamplingDecision{Sample: (!ignoreParentDecision && p.ParentContext.IsSampled()) || sample}
			},
		}
	}
}

// ConfiguredURLPrefixSampler constructs the prefix GetStartOptions getter with the default opencensus configuration
type ConfiguredURLPrefixSampler struct {
	Whitelist        config.Slice
	Blacklist        config.Slice
	AllowParentTrace bool
	area             *config.Area
}

// Inject dependencies
func (c *ConfiguredURLPrefixSampler) Inject(
	area *config.Area,
	cfg *struct {
		Whitelist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.whitelist,optional"`
		Blacklist        config.Slice `inject:"config:flamingo.opencensus.tracing.sampler.blacklist,optional"`
		AllowParentTrace bool         `inject:"config:flamingo.opencensus.tracing.sampler.allowParentTrace,optional"`
	},
) *ConfiguredURLPrefixSampler {
	c.area = area
	if cfg != nil {
		c.Whitelist = cfg.Whitelist
		c.Blacklist = cfg.Blacklist
		c.AllowParentTrace = cfg.AllowParentTrace
	}
	return c
}

// GetStartOptions constructor for ochttp.Server
func (c *ConfiguredURLPrefixSampler) GetStartOptions() func(*http.Request) trace.StartOptions {
	areas, _ := c.area.GetFlatContexts()
	prefixes := make([]string, 0, len(areas))

	for _, area := range areas {
		pathValue, pathSet := area.Configuration.Get("flamingo.router.path")
		hostValue, hostSet := area.Configuration.Get("flamingo.router.host")

		prefix := "/"
		if pathSet {
			prefix = path.Join("/", pathValue.(string), "/")
		}
		if hostSet && hostValue != "" {
			prefix = hostValue.(string) + prefix
		}

		prefixes = append(prefixes, prefix)
	}

	allowed := make([]string, 0, len(c.Whitelist)+len(prefixes)+1)
	blocked := make([]string, 0, len(c.Blacklist)+len(prefixes)+1)

	_ = c.Whitelist.MapInto(&allowed)
	_ = c.Blacklist.MapInto(&blocked)

	// prefixed routes
	for _, prefix := range prefixes {
		for _, p := range c.Whitelist {
			allowed = append(allowed, prefix+p.(string))
		}
		for _, p := range c.Blacklist {
			blocked = append(blocked, prefix+p.(string))
		}
	}

	return URLPrefixSampler(allowed, blocked, c.AllowParentTrace)
}
