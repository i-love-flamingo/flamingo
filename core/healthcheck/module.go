/*
Flamingo Package that provides a healthcheck endpoint under the default route /status/healthcheck
Usage:
 Register your own Status via Dingo:
 injector.BindMap((*healthcheck.Status)(nil), "yourServiceName").To(yourServiceNameApi.Status{})

*/
package healthcheck

import (
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/core/healthcheck/interfaces/controllers"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
)

type Module struct {
	CheckSession    bool   `inject:"config:healthcheck.checkSession"`
	CheckAuthServer bool   `inject:"config:healthcheck.checkAuth"`
	SessionBackend  string `inject:"config:session.backend"`
}

func (m *Module) Configure(injector *dingo.Injector) {
	router.Bind(injector, new(routes))

	if m.CheckSession {
		switch m.SessionBackend {
		case "redis":
			injector.BindMap((*healthcheck.Status)(nil), "session").To(healthcheck.RedisSession{})
		case "file":
			injector.BindMap((*healthcheck.Status)(nil), "session").To(healthcheck.FileSession{})
		default:
			injector.BindMap((*healthcheck.Status)(nil), "session").To(healthcheck.Nil{})
		}
	}
	if m.CheckAuthServer {
		injector.BindMap((*healthcheck.Status)(nil), "auth").To(healthcheck.Auth{})
	}
}

type routes struct {
	checkPath   string
	pingPath    string
	healthcheck *controllers.Healthcheck
}

func (r *routes) Inject(healthcheck *controllers.Healthcheck, cfg *struct {
	CheckPath string `inject:"config:healthcheck.checkPath"`
	PingPath  string `inject:"config:healthcheck.pingPath"`
}) {
	r.healthcheck = healthcheck
	r.checkPath = cfg.CheckPath
	r.pingPath = cfg.PingPath
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleGet("health.check", r.healthcheck.Healthcheck)
	registry.Route(r.checkPath, "health.check")

	registry.HandleGet("health.ping", r.healthcheck.Ping)
	registry.Route(r.pingPath, "health.ping")
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"healthcheck": config.Map{
			"checkSession": false,
			"checkAuth":    false,
			"checkPath":    "/status/healthcheck",
			"pingPath":     "/status/ping",
		},
	}
}
