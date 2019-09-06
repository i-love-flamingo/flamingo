package application

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"

	"flamingo.me/flamingo/v3/framework/opencensus"
)

var (
	// ErrorCount counts logged errors
	ErrorCount = stats.Int64("flamingo/zap/errors", "Count of logged errors", stats.UnitDimensionless)
)

func init() {
	if err := opencensus.View("flamingo/zap/errors", ErrorCount, view.Count()); err != nil {
		panic(err)
	}
}
