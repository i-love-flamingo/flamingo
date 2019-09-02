package request

import (
	"fmt"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

type (
	// Module for record request status codes
	Module struct {
	}
)

var (
	// HTTPResponseCount counts different HTTP responses
	HTTPResponseCount = stats.Int64("flamingo/request/http_response_count", "Count of http responses by status code", stats.UnitDimensionless)

	// KeyHTTPStatus defines response http status code
	KeyHTTPStatus, _ = tag.NewKey("status_code")
)

func init() {
	if err := opencensus.View("flamingo/request/http_response_count", HTTPResponseCount, view.Count(), KeyHTTPStatus); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(web.Filter)).To(new(metricsFilter))
}

// DefaultConfig configures module's default configuration
func (m *Module) DefaultConfig() config.Map {
	return config.Map{}
}
