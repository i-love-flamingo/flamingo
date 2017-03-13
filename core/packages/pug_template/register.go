package pug_template

import (
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/router"
	"flamingo/core/packages/pug_template/pugast"
	"flamingo/core/packages/pug_template/template_functions"
	"flamingo/core/template"
	"net/http"
)

// Register Services for pug_template package
func Register(basedir string, debug bool) di.RegisterFunc {
	return func(c *di.Container) {
		c.Register(func(r *router.Router) {
			r.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(basedir))))
			r.Route("/static/{n:.*}", "_static")

			r.Handle("_pugtpl_debug", new(DebugController))
			r.Route("/_pugtpl/debug", "_pugtpl_debug")
		}, router.RouterRegister)

		c.Register(pugast.NewPugTemplateEngine(basedir, debug))
		c.Register(new(template.FunctionRegistry))
		c.Register(new(template_functions.AssetFunc), "template.func")
		c.Register(new(template_functions.DebugFunc), "template.func")
		c.Register(new(template_functions.MathLib), "template.func")
	}
}
