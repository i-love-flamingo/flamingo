package gotemplate

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// Module for gotemplate engine
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind(new(flamingo.TemplateEngine)).In(dingo.ChildSingleton).To(engine{})
	injector.Bind(new(urlRouter)).To(web.Router{})

	flamingo.BindTemplateFunc(injector, "url", new(urlFunc))
	flamingo.BindTemplateFunc(injector, "get", new(getFunc))
	flamingo.BindTemplateFunc(injector, "data", new(dataFunc))
}

// DefaultConfig for gotemplate module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"gotemplates.engine.templates.basepath": "templates",
		"gotemplates.engine.layout.dir":         "",
	}
}
