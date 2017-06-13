// Package config provides supporting code for multi-tenant setups
package config

import (
	"flamingo/framework/dingo"
	"os"
	"regexp"
)

type (
	// Area defines a configuration area for multi-site setups
	// it is initialized by project main package and partly loaded by config files
	Area struct {
		Name    string
		BaseURL string

		Parent   *Area `json:"-"`
		Childs   []*Area
		Modules  []dingo.Module
		Injector *dingo.Injector `json:"-"`

		Routes        []Route                `yaml:"routes"`
		Configuration map[string]interface{} `yaml:"config" json:"config"`
	}

	// DefaultConfigModule is used to get a module's default configuration
	DefaultConfigModule interface {
		DefaultConfig() map[string]interface{}
	}

	// Route defines the yaml structure for a route, consisting of a path and a controller, as well as optional args
	Route struct {
		Path       string
		Controller string
		Name       string
	}
)

// NewArea creates a new Area with childs
func NewArea(name string, modules []dingo.Module, baseURL string, childs ...*Area) *Area {
	ctx := &Area{
		Name:    name,
		Modules: modules,
		Childs:  childs,
		BaseURL: baseURL,
	}

	for _, c := range childs {
		c.Parent = ctx
	}

	return ctx
}

// GetFlatContexts returns a map of context-relative-name->*Area, which has been flatted to inherit all parent's
// tree settings such as DI & co, and filtered to only list tree nodes specified by Contexts of area.
func (area *Area) GetFlatContexts() []*Area {
	var result []*Area
	flat := area.Flat()

	for relativeContextKey, context := range flat {
		if context.BaseURL == "" {
			continue
		}
		result = append(result, &Area{
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
func (area *Area) GetInitializedInjector() *dingo.Injector {
	var injector *dingo.Injector
	if area.Parent != nil {
		injector = area.Parent.Injector.Child()
	} else {
		injector = dingo.NewInjector()
	}
	injector.Bind(Area{}).ToInstance(area)

	for _, module := range area.Modules {
		if cfgmodule, ok := module.(DefaultConfigModule); ok {
			for k, v := range cfgmodule.DefaultConfig() {
				if _, ok := area.Configuration[k]; !ok {
					area.Configuration[k] = v
				}
			}
		}
	}

	var regex = regexp.MustCompile(`%%ENV:([^%]+)%%`)
	for k, v := range area.Configuration {
		if val, ok := v.(string); ok {
			v = regex.ReplaceAllStringFunc(val, func(a string) string { return os.Getenv(regex.FindStringSubmatch(a)[1]) })
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}

	injector.InitModules(area.Modules...)

	return injector
}

// Flat returns a map of name->*Area of contexts, were all values have been inherited (yet overriden) of the parent context tree.
func (area *Area) Flat() map[string]*Area {
	res := make(map[string]*Area)
	res[area.Name] = area

	area.Injector = area.GetInitializedInjector()

	for _, child := range area.Childs {
		for cn, flatchild := range child.Flat() {
			res[area.Name+`/`+cn] = MergeFrom(*flatchild, *area)
		}
	}

	return res
}

// MergeFrom merges two Contexts into a new one
// We do not merge config, as we use the DI to handle it
func MergeFrom(baseContext, incomingContext Area) *Area {
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
func (area *Area) Config(key string) interface{} {
	if config, ok := area.Configuration[key]; ok {
		return config
	}

	if area.Parent != nil {
		return area.Parent.Config(key)
	}

	return nil
}
