package domain

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

var (
	// ErrorCount counts logged errors
	ErrorCount = stats.Int64("flamingo/error_count", "Count of logged errors", stats.UnitDimensionless)

	// KeyArea identifies the current application area
	KeyArea, _ = tag.NewKey("area")
)
