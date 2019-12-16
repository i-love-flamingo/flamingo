package robotstxt

import (
	"net/http"

	"flamingo.me/dingo"
)

type (
	// Module for robotstxt
	Module struct {
		DefaultMux *http.ServeMux `inject:",optional"`
		Filepath   string         `inject:"config:core.robotstxt.filepath"`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	if m.DefaultMux != nil {
		m.DefaultMux.HandleFunc("/robots.txt", func(rw http.ResponseWriter, req *http.Request) {
			http.ServeFile(rw, req, m.Filepath)
		})
	}
}

// CueConfig schema
func (*Module) CueConfig() string {
	return `
core: robotstxt: filepath: string | *"frontend/robots.txt"
`
}

// FlamingoLegacyConfigAlias mapping
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"robotstxt.filepath": "core.robotstxt.filepath"}
}
