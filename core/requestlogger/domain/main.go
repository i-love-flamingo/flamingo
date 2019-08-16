package domain

import (
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"flamingo.me/flamingo/v3/framework/opencensus"
)

var (
	// ResponseHTTPStatusCount counts different HTTP response type
	ResponseHTTPStatusCount = stats.Int64("flamingo/requestlogger_response_http_status_count", "Count of specific (http) response type", stats.UnitDimensionless)

	// KeyHTTPStatus defines response http status code
	KeyHTTPStatus, _ = tag.NewKey("http-status")

	// KeyArea identifies the current application area
	KeyArea, _ = tag.NewKey("area")
)

func init() {
	if err := opencensus.View("flamingo/requestlogger_response_http_status_count", ResponseHTTPStatusCount, view.Count(), KeyHTTPStatus); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}
