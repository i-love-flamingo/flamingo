// Package context provides supporting code for multi-tenant setups
package context

import (
	"flamingo/core/dingo"
	"fmt"
)

type (
	// Context defines a configuration context for multi-site setups
	Context struct {
		Name    string
		BaseURL string

		Parent  *Context `json:"-"`
		Childs  []*Context
		Modules []dingo.Module
		//ServiceContainer *di.Container `json:"-"`
		Injector *dingo.Injector `json:"-"`

		Routes        []Route                `yaml:"routes"`
		Configuration map[string]interface{} `yaml:"config" json:"config"`
		Contexts      map[string]string      `yaml:"contexts"`
	}

	// Route defines the yaml structure for a route, consisting of a path and a controller, as well as optional args
	Route struct {
		Path       string
		Controller string
		Args       map[string]string
	}
)

// New returns Context Pointers with RegisterFuncs.
func New(name string, modules []dingo.Module, childs ...*Context) *Context {
	ctx := &Context{
		Name:    name,
		Modules: modules,
		Childs:  childs,
	}

	for _, c := range childs {
		c.Parent = ctx
	}

	return ctx
}

// GetFlatContexts returns a map of context-relative-name->*Context, which has been flatted to inherit all parent's
// tree settings such as DI & co, and filtered to only list tree nodes specified by Contexts of ctx.
func (ctx *Context) GetFlatContexts() map[string]*Context {
	result := make(map[string]*Context)
	flat := ctx.Flat()
	for baseurl, name := range ctx.Contexts {
		result[name] = flat[ctx.Name+`/`+name]
		result[name].BaseURL = baseurl
		result[name].Childs = nil
		result[name].Contexts = nil
		result[name].Name = name
		result[name].Injector = dingo.NewInjector(result[name].Modules...)
		for k, v := range result[name].Configuration {
			//result[name].ServiceContainer.SetParameter(k, v)
			result[name].Injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
		}
	}

	fmt.Println(result)

	return result
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

	//baseContext.RegisterFuncs = append(incomingContext.RegisterFuncs, baseContext.RegisterFuncs...)
	baseContext.Modules = append(incomingContext.Modules, baseContext.Modules...)

	return &baseContext
}
