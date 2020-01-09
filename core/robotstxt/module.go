package robotstxt

import (
	"net/http"

	"flamingo.me/dingo"
)

type (
	// Module for robotstxt
	Module struct {
		DefaultMux           *http.ServeMux `inject:",optional"`
		RobotsTxtFilepath    string
		SecurityTxtActivated bool
		SecurityTxtFilepath  string
		HumansTxtActivated   bool
		HumansTxtFilepath    string
	}
)

// Inject dependencies
func (m *Module) Inject(
	config *struct {
		RobotsTxtFilepath    string `inject:"config:core.robotstxt.filepath"`
		SecurityTxtActivated bool   `inject:"config:core.securitytxt.enabled,optional"`
		SecurityTxtFilepath  string `inject:"config:core.securitytxt.filepath"`
		HumansTxtActivated   bool   `inject:"config:core.humanstxt.enabled,optional"`
		HumansTxtFilepath    string `inject:"config:core.humanstxt.filepath"`
	},
) {
	m.RobotsTxtFilepath = config.RobotsTxtFilepath
	m.SecurityTxtActivated = config.SecurityTxtActivated
	m.SecurityTxtFilepath = config.SecurityTxtFilepath
	m.HumansTxtActivated = config.HumansTxtActivated
	m.HumansTxtFilepath = config.HumansTxtFilepath
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	if m.DefaultMux != nil {
		m.DefaultMux.HandleFunc("/robots.txt", func(rw http.ResponseWriter, req *http.Request) {
			http.ServeFile(rw, req, m.RobotsTxtFilepath)
		})

		if m.SecurityTxtActivated {
			// https://securitytxt.org/
			m.DefaultMux.HandleFunc("/.well-known/security.txt", func(rw http.ResponseWriter, req *http.Request) {
				http.ServeFile(rw, req, m.SecurityTxtFilepath)
			})
		}

		if m.HumansTxtActivated {
			// http://humanstxt.org/
			m.DefaultMux.HandleFunc("/humans.txt", func(rw http.ResponseWriter, req *http.Request) {
				http.ServeFile(rw, req, m.HumansTxtFilepath)
			})
		}
	}
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
