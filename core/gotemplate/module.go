package gotemplate

import (
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/template"
)

// Module for gotemplate engine
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*template.Engine)(nil)).In(dingo.ChildSingleton).To(engine{})

	template.BindFunc(injector, "url", new(urlFunc))
	template.BindCtxFunc(injector, "get", new(getFunc))
	template.BindCtxFunc(injector, "data", new(dataFunc))
}

// DefaultConfig for gotemplate module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"gotemplates.engine.templates.basepath": "templates",
		"gotemplates.engine.layout.dir":         "",
	}
}
