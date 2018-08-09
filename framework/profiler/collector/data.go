package collector

import "flamingo.me/flamingo/framework/web"

type (
	// DataCollector for external collections
	// deprecated: use opencensus
	DataCollector interface {
		Collect(ctx web.Context) string
	}
)
