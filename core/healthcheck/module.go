// Package healthcheck provides a healthcheck endpoint under the default route /status/healthcheck
// Usage:
// Register your own Status via Dingo:
// injector.BindMap(new(healthcheck.Status), "yourServiceName").To(yourServiceNameApi.Status{})
package healthcheck

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/core/healthcheck/interfaces/controllers"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/systemendpoint"
	"flamingo.me/flamingo/v3/framework/systemendpoint/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module entry point for the flamingo healthcheck module
type Module struct {
	controller      *controllers.Healthcheck
	checkSession    bool
	checkAuthServer bool
	checkPath       string
	pingPath        string
	sessionBackend  string
}

// Inject dependencies
func (m *Module) Inject(
	controller *controllers.Healthcheck,
	config *struct {
		CheckSession    bool   `inject:"config:healthcheck.checkSession"`
		CheckAuthServer bool   `inject:"config:healthcheck.checkAuth"`
		CheckPath       string `inject:"config:healthcheck.checkPath"`
		PingPath        string `inject:"config:healthcheck.pingPath"`
		SessionBackend  string `inject:"config:session.backend"`
	},
) {
	m.controller = controller
	m.checkSession = config.CheckSession
	m.checkAuthServer = config.CheckAuthServer
	m.checkPath = config.CheckPath
	m.pingPath = config.PingPath
	m.sessionBackend = config.SessionBackend
}

type routes struct {
	controller *controllers.Ping
}

func (r *routes) Inject(controller *controllers.Ping) {
	r.controller = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleAny("core.healtchcheck.ping", web.WrapHTTPHandler(r.controller))
	registry.Route("/status/ping", "core.healthcheck.ping")
}

// Configure dependency injection
func (m *Module) Configure(injector *dingo.Injector) {
	if m.checkSession {
		switch m.sessionBackend {
		case "redis":
			injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.RedisSession{})
		case "file":
			injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.FileSession{})
		default:
			injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.Nil{})
		}
	}
	if m.checkAuthServer {
		injector.BindMap((*healthcheck.Status)(nil), "auth").To(healthcheck.Auth{})
	}

	injector.BindMap((*domain.Handler)(nil), m.pingPath).To(&controllers.Ping{})
	injector.BindMap((*domain.Handler)(nil), m.checkPath).To(&controllers.Healthcheck{})

	web.BindRoutes(injector, new(routes))
}

// DefaultConfig for healthcheck module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"healthcheck": config.Map{
			"checkSession": true,
			"checkAuth":    false,
			"checkPath":    "/status/healthcheck",
			"pingPath":     "/status/ping",
		},
	}
}

// Depends on other modules
func (m *Module) Depends() []dingo.Module {
	return []dingo.Module{
		new(systemendpoint.Module),
	}
}
