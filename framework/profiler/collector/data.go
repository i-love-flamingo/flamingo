package collector

import "go.aoe.com/flamingo/framework/web"

type (
	// DataCollector for external collections
	DataCollector interface {
		Collect(ctx web.Context) string
	}
)
