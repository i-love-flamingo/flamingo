// Package healthcheck provides a healthcheck endpoint under the default route /status/healthcheck
// Usage:
// Register your own Status via Dingo:
// injector.BindMap(new(healthcheck.Status), "yourServiceName").To(yourServiceNameApi.Status{})
package healthcheck

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/interfaces/controllers"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module entry point for the flamingo healthcheck module
type Module struct {
	controller *controllers.Healthcheck
	checkPath  string
	pingPath   string
}

// Inject dependencies
func (m *Module) Inject(
	controller *controllers.Healthcheck,
	config *struct {
		CheckPath string `inject:"config:core.healthcheck.checkPath"`
		PingPath  string `inject:"config:core.healthcheck.pingPath"`
	},
) {
	m.controller = controller
	m.checkPath = config.CheckPath
	m.pingPath = config.PingPath
}

type routes struct {
	controller *controllers.Ping
}

func (r *routes) Inject(controller *controllers.Ping) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleAny("core.healthcheck.ping", web.WrapHTTPHandler(r.controller))
	registry.MustRoute("/health/ping", "core.healthcheck.ping")
}

// Configure dependency injection
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMap((*domain.Handler)(nil), m.pingPath).To(&controllers.Ping{})
	injector.BindMap((*domain.Handler)(nil), m.checkPath).To(&controllers.Healthcheck{})

	web.BindRoutes(injector, new(routes))
	injector.BindMulti((*prefixrouter.OptionalHandler)(nil)).AnnotatedWith("fallback").To(controllers.Ping{})
}

// CueConfig schema and configuration
func (m *Module) CueConfig() string {
	return `
core healthcheck: {
	checkAuth: bool | *false
	checkPath: string | *"/status/healthcheck"
	pingPath: string | *"/status/ping"
}
`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"healthcheck.checkAuth": "core.healthcheck.checkAuth",
		"healthcheck.checkPath": "core.healthcheck.checkPath",
		"healthcheck.pingPath":  "core.healthcheck.pingPath",
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(systemendpoint.Module),
	}
}
