package domain

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

var (
	// HTTPResponseCount counts different HTTP responses
	HTTPResponseCount = stats.Int64("flamingo/requestlogger_http_response_count", "Count of http responses by status code", stats.UnitDimensionless)

	// KeyHTTPStatus defines response http status code
	KeyHTTPStatus, _ = tag.NewKey("status_code")

	// KeyArea identifies the current application area
	KeyArea, _ = tag.NewKey("area")
)
