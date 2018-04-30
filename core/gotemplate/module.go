package gotemplate

import (
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/template"
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
		"gotemplates.engine.glob": "templates/*",
	}
}
