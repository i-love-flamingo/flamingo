package pug_template

import (
	"flamingo/core/dingo"
	"flamingo/core/pug_template/pugast"
	"flamingo/core/pug_template/template_functions"
	"flamingo/core/template"
	"flamingo/framework/router"
	"net/http"
)

type Module struct {
	RouterRegistry *router.RouterRegistry `inject:""`
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/dist"))))
	module.RouterRegistry.Route("/static/{n:.*}", "_static")

	module.RouterRegistry.Handle("_pugtpl_debug", new(DebugController))
	module.RouterRegistry.Route("/_pugtpl/debug", "_pugtpl_debug")

	injector.Bind((*pugast.PugTemplateEngine)(nil)).AsEagerSingleton().ToProvider(func() *pugast.PugTemplateEngine { return pugast.NewPugTemplateEngine("frontend/dist", true) })
	injector.Bind((*template.Engine)(nil)).AsEagerSingleton().To(pugast.PugTemplateEngine{})

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.MathLib{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.DebugFunc{})
}
