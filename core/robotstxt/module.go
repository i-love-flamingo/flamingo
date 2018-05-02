package robotstxt

import (
	"net/http"

	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
)

type (
	Module struct {
		DefaultMux *http.ServeMux `inject:",optional"`
		Filepath   string         `inject:"config:robotstxt.filepath"`
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

// DefaultConfig for setting pug-related config options
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"robotstxt.filepath": "frontend/robots.txt",
	}
}
