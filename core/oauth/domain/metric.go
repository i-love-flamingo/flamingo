package domain

import (
	"fmt"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	// LoginFailCount counts the failed login attempts
	LoginFailCount = stats.Int64("flamingo/oauth_login_fail_count", "Count of failed login attempts", stats.UnitDimensionless)
)

func init() {
	if err := opencensus.View("flamingo/oauth_login_fail_count", LoginFailCount, view.Count()); err != nil {
		panic(fmt.Sprintf("failed to register opencensus view: %s", err))
	}
}
