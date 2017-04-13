package pug_template

import (
	"flamingo/core/dingo"
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/router"
	"flamingo/core/packages/pug_template/pugast"
	"flamingo/core/packages/pug_template/template_functions"
	"flamingo/core/template"
	"net/http"
)

type Module struct {
	Router *router.Router `inject:""`
}

func (module *Module) Configure(injector *dingo.Injector) {
	module.Router.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/dist"))))
	module.Router.Route("/static/{n:.*}", "_static")

	module.Router.Handle("_pugtpl_debug", new(DebugController))
	module.Router.Route("/_pugtpl/debug", "_pugtpl_debug")

	injector.Bind((*template.Engine)(nil)).ToProvider(func() template.Engine { return pugast.NewPugTemplateEngine("frontend/dist", true) })

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.DebugFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.MathLib{})
}

// Register Services for pug_template package
func Register(c *di.Container) {
	basedir, debug := c.GetParameter("pug_template.basedir").(string), c.GetParameter("pug_template.debug").(bool)

	c.Register(pugast.NewPugTemplateEngine(basedir, debug))
	c.Register(new(template.FunctionRegistry))
	c.Register(new(template_functions.AssetFunc), "template.func")
	c.Register(new(template_functions.DebugFunc), "template.func")
	c.Register(new(template_functions.MathLib), "template.func")
}
