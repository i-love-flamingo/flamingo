package opencensus

/*
Package opencensus contains support for OpenCensus usage.

Enabling the module will run a prometheus exporter on port 13210.

Using metrics in your module:

	package controller

	import (
		"flamingo.me/flamingo/v3/framework/opencensus"
		"go.opencensus.io/stats"
		"go.opencensus.io/stats/view"
	)

	var rt = stats.Int64("flamingo/package/mystat", "my stat records 5 milliseconds per call", stats.UnitMilliseconds)

	// register view in opencensus
	func init() {
		opencensus.View("flamingo/package/mystat/sum", rt, view.Sum())
		opencensus.View("flamingo/package/mystat/count", rt, view.Count())
	}

    // measure mystat
	func (c *Controller) Get(ctx context.Context, r *web.Request) web.Response {
		// ...

		// record 5ms per call
		stats.Record(ctx, rt.M(5))

		// ...
	}
*/
