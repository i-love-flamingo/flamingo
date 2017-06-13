package gotemplate

import (
	"flamingo/core/gotemplate/framework"
	"flamingo/framework/dingo"
	"flamingo/framework/template"
)

// Module for core/go_template
type Module struct{}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.Bind((*template.Engine)(nil)).To(new(framework.Engine))
}
