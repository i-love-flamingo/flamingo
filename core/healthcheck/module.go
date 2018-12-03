/*
Flamingo Package that provides a healthcheck endpoint under the default route /status/healthcheck
Usage:
 Register your own Status via Dingo:
 injector.BindMap((*healthcheck.Status)(nil), "yourServiceName").To(yourServiceNameApi.Status{})

*/
package healthcheck

import (
	"flamingo.me/flamingo/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/core/healthcheck/interfaces/controllers"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
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
	path        string
	healthcheck *controllers.Healthcheck
}

func (r *routes) Inject(healthcheck *controllers.Healthcheck, cfg *struct {
	Path string `inject:"config:healthcheck.path"`
}) {
	r.healthcheck = healthcheck
	r.path = cfg.Path
}

func (r *routes) Routes(registry *router.Registry) {
	registry.HandleGet("healthcheck", r.healthcheck.Get)
	registry.Route(r.path, "healthcheck")
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"healthcheck": config.Map{
			"checkSession": false,
			"checkAuth":    false,
			"path":         "/status/healthcheck",
		},
	}
}
