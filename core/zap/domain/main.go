package domain

import (
	"fmt"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"flamingo.me/flamingo/v3/framework/opencensus"
)

var (
	// ErrorCount counts logged errors
	ErrorCount = stats.Int64("flamingo/error_count", "Count of logged errors", stats.UnitDimensionless)

	// KeyArea identifies the current application area
	KeyArea, _ = tag.NewKey("area")
)

func init() {
	if err := opencensus.View("flamingo/error_count", ErrorCount, view.Count()); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}
