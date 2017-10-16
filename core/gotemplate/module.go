package gotemplate

import (
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/template"
)

type Module struct{}

func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*template.Engine)(nil)).To(engine{})
	injector.BindMulti((*template.Function)(nil)).To(urlFunc{})
	injector.BindMulti((*template.ContextFunction)(nil)).To(getFunc{})
}

func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"gotemplates.engine.glob": "templates/*",
	}
}
