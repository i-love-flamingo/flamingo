package gotemplate

import (
	"flamingo.me/dingo"

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
	flamingo.BindTemplateFunc(injector, "plainHtml", new(plainHTMLFunc))
	flamingo.BindTemplateFunc(injector, "plainJs", new(plainJSFunc))
}

// CueConfig definition
func (m *Module) CueConfig() string {
	return `
// general config
core gotemplate engine: {
	templates basepath: string | *"templates"
	layout dir: string | *""
}
`
}

// FlamingoLegacyConfigAlias mapping
func (m *Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{
		"gotemplates.engine.templates.basepath": "core.gotemplate.engine.templates.basepath",
		"gotemplates.engine.layout.dir":         "core.gotemplate.engine.layout.dir",
	}
}
