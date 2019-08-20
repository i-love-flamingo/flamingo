package requestlogger

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/requestlogger/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"fmt"
	"go.opencensus.io/stats/view"
)

type (
	// Module for core/requestlogger
	Module struct {
		trackResponseCount bool
	}
)

// Inject module dependencies
func (m *Module) Inject(cfg *struct {
	TrackResponseCount bool `inject:"config:requestlogger.metrics.responseCountTracking.enabled"`
}) *Module {
	if cfg != nil {
		m.trackResponseCount = cfg.TrackResponseCount
	}

	return m
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).To(logger{})

	if m.trackResponseCount {
		if err := opencensus.View("flamingo/requestlogger_http_response_count", domain.HTTPResponseCount, view.Count(), domain.KeyHTTPStatus); err != nil {
			panic(fmt.Sprintf("failed to register opencensus view: %s", err))
		}
	}
}

// DefaultConfig configures module's default configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"requestlogger.metrics.responseCountTracking.enabled": false,
	}
}
