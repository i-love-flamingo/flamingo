// Package context provides supporting code for multi-tenant setups
package context

import (
	"flamingo/framework/dingo"
	"os"
	"strings"
)

type (
	// Context defines a configuration context for multi-site setups
	// it is initialized by project main package and partly loaded by config files
	Context struct {
		Name    string
		BaseURL string

		Parent   *Context `json:"-"`
		Childs   []*Context
		Modules  []dingo.Module
		Injector *dingo.Injector `json:"-"`

		Routes        []Route                `yaml:"routes"`
		Configuration map[string]interface{} `yaml:"config" json:"config"`
	}

	// Route defines the yaml structure for a route, consisting of a path and a controller, as well as optional args
	Route struct {
		Path       string
		Controller string
		Name       string
	}
)

// New creates a new Context with childs
func New(name string, modules []dingo.Module, baseURl string, childs ...*Context) *Context {
	ctx := &Context{
		Name:    name,
		Modules: modules,
		Childs:  childs,
		BaseURL: baseURl,
	}

	for _, c := range childs {
		c.Parent = ctx
	}

	return ctx
}

// GetFlatContexts returns a map of context-relative-name->*Context, which has been flatted to inherit all parent's
// tree settings such as DI & co, and filtered to only list tree nodes specified by Contexts of ctx.
func (ctx *Context) GetFlatContexts() []*Context {
	var result []*Context
	flat := ctx.Flat()

	for relativeContextKey, context := range flat {
		if context.BaseURL == "" {
			continue
		}
		result = append(result, &Context{
			Name:     relativeContextKey,
			BaseURL:  context.BaseURL,
			Routes:   context.Routes,
			Injector: context.Injector,
		})

	}
	return result
}

// GetInitializedInjector returns initialized container based on the configuration
// we derive our injector from our parent
func (ctx *Context) GetInitializedInjector() *dingo.Injector {
	var injector *dingo.Injector
	if ctx.Parent != nil {
		injector = ctx.Parent.Injector.Child()
	} else {
		injector = dingo.NewInjector()
	}
	injector.Bind(Context{}).ToInstance(ctx)

	for k, v := range ctx.Configuration {
		if val, ok := v.(string); ok && strings.HasPrefix(val, "%%ENV:") && strings.HasSuffix(val, "%%") {
			v = os.Getenv(val[6 : len(val)-2])
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}

	injector.InitModules(ctx.Modules...)

	return injector
}

// Flat returns a map of name->*Context of contexts, were all values have been inherited (yet overriden) of the parent context tree.
func (ctx *Context) Flat() map[string]*Context {
	res := make(map[string]*Context)
	res[ctx.Name] = ctx

	ctx.Injector = ctx.GetInitializedInjector()

	for _, child := range ctx.Childs {
		for cn, flatchild := range child.Flat() {
			res[ctx.Name+`/`+cn] = MergeFrom(*flatchild, *ctx)
		}
	}

	return res
}

// MergeFrom merges two Contexts into a new one
// We do not merge config, as we use the DI to handle it
func MergeFrom(baseContext, incomingContext Context) *Context {
	if baseContext.Configuration == nil {
		baseContext.Configuration = make(map[string]interface{})
	}

	knownhandler := make(map[string]bool)
	for _, route := range baseContext.Routes {
		knownhandler[route.Controller] = true
	}

	for _, route := range incomingContext.Routes {
		if !knownhandler[route.Controller] {
			baseContext.Routes = append(baseContext.Routes, route)
		}
	}

	return &baseContext
}

// Config get a config value recursive
func (ctx *Context) Config(key string) interface{} {
	if config, ok := ctx.Configuration[key]; ok {
		return config
	}

	if ctx.Parent != nil {
		return ctx.Parent.Config(key)
	}

	return nil
}
