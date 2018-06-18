package gotemplate

import (
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/template"
)

// Module for gotemplate engine
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*template.Engine)(nil)).In(dingo.ChildSingleton).To(engine{})
	injector.BindMulti((*template.Function)(nil)).To(urlFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(getFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(dataFunc{})
}

// DefaultConfig for gotemplate module
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"gotemplates.engine.templates.basepath": "templates",
		"gotemplates.engine.layout.dir":         "",
	}
}
