package domain

import (
	"go.opencensus.io/stats"
)

var (
	// LoginFailedCount counts the failed login attempts
	LoginFailedCount = stats.Int64("flamingo/oauth_login_failed_count", "Count of failed login attempts", stats.UnitDimensionless)
	// LoginSucceededCount counts the successful login attempts
	LoginSucceededCount = stats.Int64("flamingo/oauth_login_succeeded_count", "Count of succeeded login attempts", stats.UnitDimensionless)
)
