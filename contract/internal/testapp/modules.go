// Package testapp is the application that the process tests start, with its config.
package testapp

import (
	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/basket"
	"flamingo.me/flamingo/v3/contract/internal/testapp/modules/ops"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
)

// Main starts the application with the prefix router and the areas "en" and "de".
func Main() {
	flamingo.App(
		[]dingo.Module{
			new(zap.Module),
			new(prefixrouter.Module),
			new(systemendpoint.Module),
			new(opencensus.Module),
			new(healthcheck.Module),
			new(flamingoFramework.SessionModule),
			new(ops.Module),
		},
		flamingo.ChildAreas(
			config.NewArea("en", []dingo.Module{new(basket.Module)}),
			config.NewArea("de", []dingo.Module{new(basket.Module)}),
		),
	)
}

// PlainServe starts an application without prefixrouter, so flamingo's own serve command runs the server.
func PlainServe() {
	flamingo.App([]dingo.Module{
		new(zap.Module),
		new(systemendpoint.Module),
		new(ops.Module),
	})
}
