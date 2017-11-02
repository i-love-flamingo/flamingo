package gotemplate

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"

	"go.aoe.com/flamingo/framework/router"
	flamingotemplate "go.aoe.com/flamingo/framework/template"
	"go.aoe.com/flamingo/framework/web"
)

type (
	engine struct {
		Glob              string                             `inject:"config:gotemplates.engine.glob"`
		TemplateFunctions *flamingotemplate.FunctionRegistry `inject:""`
	}

	// urlFunc allows templates to access the routers `URL` helper method
	urlFunc struct {
		Router *router.Router `inject:""`
	}

	// getFunc allows templates to access the router's `get` method
	getFunc struct {
		Router *router.Router `inject:""`
	}
)

func (e *engine) Render(context web.Context, name string, data interface{}) (io.Reader, error) {
	done := context.Profile("template engine", "load templates")

	functionsMap := template.FuncMap{
		"Upper": strings.ToUpper,
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"map": func(p ...interface{}) map[string]interface{} {
			res := make(map[string]interface{})
			for i := 0; i < len(p); i += 2 {
				res[fmt.Sprint(p[i])] = p[i+1]
			}
			return res
		},
	}

	funcs := e.TemplateFunctions.Populate()
	for k, f := range e.TemplateFunctions.ContextAware {
		funcs[k] = f(context)
	}

	tpl := template.Must(template.New("").Funcs(functionsMap).Funcs(funcs).ParseGlob(e.Glob))

	done()

	defer context.Profile("template engine", "render template "+name)()
	buf := &bytes.Buffer{}
	err := tpl.ExecuteTemplate(buf, name + ".html", data)

	return buf, err
}

// Name alias for use in template
func (g getFunc) Name() string {
	return "get"
}

// Func as implementation of get method
func (g *getFunc) Func(ctx web.Context) interface{} {
	return func(what string, params ...map[string]interface{}) interface{} {
		var p = make(map[interface{}]interface{})
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return g.Router.Get(what, ctx, p)
	}
}

// Name alias for use in template
func (u urlFunc) Name() string {
	return "url"
}

// Func as implementation of url method
func (u *urlFunc) Func() interface{} {
	return func(where string, params ...map[string]interface{}) template.URL {
		var p = make(map[string]string)
		if len(params) == 1 {
			for k, v := range params[0] {
				p[k] = fmt.Sprint(v)
			}
		}
		return template.URL(u.Router.URL(where, p).String())
	}
}
