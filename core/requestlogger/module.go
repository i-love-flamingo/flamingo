package requestlogger

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/opencensus/request"
	"flamingo.me/flamingo/v3/framework/web"
	"fmt"
	"go.opencensus.io/stats/view"
)

type (
	// Module for core/requestlogger
	Module struct{}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).To(logger{})

	if err := opencensus.View("flamingo/request/http_response_count", request.HTTPResponseCount, view.Count(), request.KeyHTTPStatus); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}

// DefaultConfig configures module's default configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{}
}
