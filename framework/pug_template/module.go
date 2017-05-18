package pug_template

import (
	"flamingo/framework/dingo"
	"flamingo/framework/pug_template/pugast"
	"flamingo/framework/pug_template/template_functions"
	"flamingo/framework/router"
	"flamingo/framework/template"
	"net/http"
)

type Module struct {
	RouterRegistry *router.RouterRegistry `inject:""`
	Basedir        string                 `inject:"config:pug_template.basedir"`
}

func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(m.Basedir))))
	m.RouterRegistry.Route("/static/{n:.*}", "_static")

	m.RouterRegistry.Handle("_pugtpl_debug", new(DebugController))
	m.RouterRegistry.Route("/_pugtpl/debug", "_pugtpl_debug")

	// We bind the Template Engine to the ChildSingleton level (in case there is different config handling
	// We use the provider to make sure both are always the same injected type
	injector.Bind(pugast.PugTemplateEngine{}).In(dingo.ChildSingleton)
	injector.Bind((*template.Engine)(nil)).
		In(dingo.ChildSingleton).
		ToProvider(func(t *pugast.PugTemplateEngine, i *dingo.Injector) template.Engine { return (template.Engine)(t) })

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.MathLib{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.DebugFunc{})
}
