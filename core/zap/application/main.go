package application

import (
	"go.opencensus.io/stats"
)

var (
	// ErrorCount counts logged errors
	ErrorCount = stats.Int64("flamingo/zap/errors", "Count of logged errors", stats.UnitDimensionless)
)
