package pug_template

import (
	di "flamingo/framework/dependencyinjection"
	"flamingo/framework/router"
	"flamingo/core/pug_template/pugast"
	"flamingo/core/pug_template/template_functions"
	"flamingo/core/template"
	"net/http"
)

// Register Services for pug_template package
func Register(c *di.Container) {
	basedir, debug := c.GetParameter("pug_template.basedir").(string), c.GetParameter("pug_template.debug").(bool)

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
