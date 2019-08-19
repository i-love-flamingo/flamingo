package domain

import (
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"flamingo.me/flamingo/v3/framework/opencensus"
)

var (
	// HTTPResponseCount counts different HTTP responses
	HTTPResponseCount = stats.Int64("flamingo/requestlogger_http_response_count", "Count of http responses by status code", stats.UnitDimensionless)

	// KeyHTTPStatus defines response http status code
	KeyHTTPStatus, _ = tag.NewKey("status_code")

	// KeyArea identifies the current application area
	KeyArea, _ = tag.NewKey("area")
)

func init() {
	if err := opencensus.View("flamingo/requestlogger_http_response_count", HTTPResponseCount, view.Count(), KeyHTTPStatus); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}
