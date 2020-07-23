package robotstxt

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/robotstxt/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module for robotstxt
	Module struct{}

	routes struct {
		securityTxtActivated bool
		humansTxtActivated   bool
		files                interfaces.FileControllerInterface
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(interfaces.FileControllerInterface)).To(interfaces.DefaultFileControllerInterface{})
	web.BindRoutes(injector, new(routes))
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: robotstxt: filepath: string | *"frontend/robots.txt"
core: securitytxt: enabled:  bool | *false
core: securitytxt: filepath: string | *"frontend/security.txt"
core: humanstxt: enabled:  bool | *false
core: humanstxt: filepath: string | *"frontend/humans.txt"
`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"robotstxt.filepath": "core.robotstxt.filepath"}
}

// Inject routes dependencies
func (r *routes) Inject(
	config *struct {
		SecurityTxtActivated bool `inject:"config:core.securitytxt.enabled,optional"`
		HumansTxtActivated   bool `inject:"config:core.humanstxt.enabled,optional"`
	},
	files interfaces.FileControllerInterface,
) {
	r.securityTxtActivated = config.SecurityTxtActivated
	r.humansTxtActivated = config.HumansTxtActivated
	r.files = files
}

// Routes module
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("robotstxt.robotstxt", r.files.GetRobotsTxt)
	_, err := registry.Route("/robots.txt", "robotstxt.robotstxt")
	if err != nil {
		panic(err)
	}

	if r.securityTxtActivated {
		registry.HandleGet("robotstxt.securitytxt", r.files.GetSecurityTxt)
		_, err := registry.Route("/.well-known/security.txt", "robotstxt.securitytxt")
		if err != nil {
			panic(err)
		}
	}

	if r.humansTxtActivated {
		registry.HandleGet("robotstxt.humanstxt", r.files.GetHumansTxt)
		_, err := registry.Route("/humans.txt", "robotstxt.humanstxt")
		if err != nil {
			panic(err)
		}
	}
}
