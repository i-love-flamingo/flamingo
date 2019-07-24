package gotemplate

import (
	"context"
	"html/template"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// urlFunc allows templates to access the routers `URL` helper method
	urlFunc struct {
		router urlRouter
	}

	dataFunc struct {
		router urlRouter
	}

	// getFunc allows templates to access the router's `get` method
	getFunc struct {
		router urlRouter
	}

	// plainHTMLFunc returns the given string as plain template HTML
	plainHTMLFunc struct{}

	// plainJSFunc returns the given string as plain template JS
	plainJSFunc struct{}
)

var (
	_ flamingo.TemplateFunc = new(urlFunc)
	_ flamingo.TemplateFunc = new(getFunc)
	_ flamingo.TemplateFunc = new(dataFunc)
	_ flamingo.TemplateFunc = new(plainHTMLFunc)
	_ flamingo.TemplateFunc = new(plainJSFunc)
)

func (g *getFunc) Inject(router urlRouter) *getFunc {
	g.router = router
	return g
}

// TemplateFunc as implementation of get method
func (g *getFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...string) interface{} {
		var p = make(map[interface{}]interface{})
		for i := 0; i < len(params); i += 2 {
			p[params[i]] = params[i+1]
		}
		return g.router.Data(ctx, what, p)
	}
}

func (d *dataFunc) Inject(router urlRouter) *dataFunc {
	d.router = router
	return d
}

// Func as implementation of get method
func (d *dataFunc) Func(ctx context.Context) interface{} {
	return func(what string, params ...string) interface{} {
		var p = make(map[interface{}]interface{})
		for i := 0; i < len(params); i += 2 {
			p[params[i]] = params[i+1]
		}
		return d.router.Data(ctx, what, p)
	}
}

func (u *urlFunc) Inject(router urlRouter) *urlFunc {
	u.router = router
	return u
}

// Func as implementation of url method
func (u *urlFunc) Func(context.Context) interface{} {
	return func(where string, params ...string) template.URL {
		var p = make(map[string]string)
		for i := 0; i < len(params); i += 2 {
			p[params[i]] = params[i+1]
		}
		url, _ := u.router.Relative(where, p)
		return template.URL(url.String())
	}
}

// Func returns the given string as plain template HTML
func (s *plainHTMLFunc) Func(_ context.Context) interface{} {
	return func(in string) template.HTML {
		return template.HTML(in)
	}
}

// Func returns the given string as plain template JS
func (s *plainJSFunc) Func(_ context.Context) interface{} {
	return func(in string) template.JS {
		return template.JS(in)
	}
}
