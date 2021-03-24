package robotstxt

import (
	"net/http"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/robotstxt/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module for robotstxt
	Module struct {
		defaultMux *http.ServeMux
		filepath   string
	}

	routes struct {
		securityTxtActivated bool
		humansTxtActivated   bool
		files                *interfaces.FileController
	}
)

// Inject dependencies
func (m *Module) Inject(
	optionals *struct {
		DefaultMux *http.ServeMux `inject:",optional"`
	},
	cfg *struct {
		Filepath string `inject:"config:core.robotstxt.filepath"`
	},
) *Module {
	if optionals != nil {
		m.defaultMux = optionals.DefaultMux
	}
	if cfg != nil {
		m.filepath = cfg.Filepath
	}
	return m
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	if m.defaultMux != nil {
		m.defaultMux.HandleFunc("/robots.txt", func(rw http.ResponseWriter, req *http.Request) {
			http.ServeFile(rw, req, m.filepath)
		})
	}
	web.BindRoutes(injector, new(routes))
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: robotstxt: filepath: string | *"frontend/robots.txt"
core: securitytxt: {
	enabled:  bool | *false
	filepath: string | *"frontend/security.txt"
}
core: humanstxt: {
	enabled:  bool | *false
	filepath: string | *"frontend/humans.txt"
}
`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"robotstxt.filepath": "core.robotstxt.filepath"}
}

// Inject routes dependencies
func (r *routes) Inject(
	files *interfaces.FileController,
	cfg *struct {
		SecurityTxtActivated bool `inject:"config:core.securitytxt.enabled,optional"`
		HumansTxtActivated   bool `inject:"config:core.humanstxt.enabled,optional"`
	},
) {
	if cfg != nil {
		r.securityTxtActivated = cfg.SecurityTxtActivated
		r.humansTxtActivated = cfg.HumansTxtActivated
	}
	r.files = files
}

// Routes module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("robotstxt.robotstxt", r.files.GetRobotsTxt)
	registry.MustRoute("/robots.txt", "robotstxt.robotstxt")

	if r.securityTxtActivated {
		registry.HandleGet("robotstxt.securitytxt", r.files.GetSecurityTxt)
		registry.MustRoute("/.well-known/security.txt", "robotstxt.securitytxt")
	}

	if r.humansTxtActivated {
		registry.HandleGet("robotstxt.humanstxt", r.files.GetHumansTxt)
		registry.MustRoute("/humans.txt", "robotstxt.humanstxt")
	}
}
