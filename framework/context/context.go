// Package context provides supporting code for multi-tenant setups
package context

import (
	"flamingo/framework/dingo"
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

		// Basti: Why json - ?
		// Also: what is interface{}? why not string?
		// Also - Should we better use composite for the stuff that is Unmarschalled from yaml?
		Routes        []Route                `yaml:"routes"`
		Configuration map[string]interface{} `yaml:"config" json:"config"`
	}

	// RoutingConfig Value struct are representing the informations required by routing
	RoutingConfig struct {
		Name     string
		BaseURL  string
		Routes   []Route
		Injector *dingo.Injector
	}

	// Route defines the yaml structure for a route, consisting of a path and a controller, as well as optional args
	Route struct {
		Path       string
		Controller string
		Args       map[string]string
	}
)

// This is the RootContext - its expected that this is set from project package
var RootContext *Context

// New returns Context Pointers with RegisterFuncs.

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
func (ctx *Context) GetRoutingConfigs() []*RoutingConfig {
	var result []*RoutingConfig
	flat := ctx.Flat()

	for _, context := range flat {
		if context.BaseURL == "" {
			continue
		}
		result = append(result, &RoutingConfig{
			Name:     context.Name,
			BaseURL:  context.BaseURL,
			Routes:   context.Routes,
			Injector: context.GetInitializedInjector(),
		})

	}
	return result
}

// Returns initialized container - based on the configuration
func (ctx *Context) GetInitializedInjector() *dingo.Injector {
	injector := dingo.NewInjector()
	for k, v := range ctx.Configuration {
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}
	injector.InitModules(ctx.Modules...)
	return injector
}

// Flat returns a map of name->*Context of contexts, were all values have been inherited (yet overriden) of the parent context tree.
func (ctx *Context) Flat() map[string]*Context {
	res := make(map[string]*Context)
	res[ctx.Name] = ctx

	for _, child := range ctx.Childs {
		for cn, flatchild := range child.Flat() {
			res[ctx.Name+`/`+cn] = MergeFrom(*flatchild, *ctx)
		}
	}

	return res
}

// MergeFrom merges two Contexts into a new one
func MergeFrom(baseContext, incomingContext Context) *Context {
	if baseContext.Configuration == nil {
		baseContext.Configuration = make(map[string]interface{})
	}

	for k, v := range incomingContext.Configuration {
		if _, ok := baseContext.Configuration[k]; !ok {
			baseContext.Configuration[k] = v
		}
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

	baseContext.Modules = append(incomingContext.Modules, baseContext.Modules...)

	return &baseContext
}
