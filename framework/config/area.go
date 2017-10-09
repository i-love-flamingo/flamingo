// Package config provides supporting code for multi-tenant setups
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/framework/dingo"
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

		Routes        []Route `yaml:"routes"`
		Configuration Map     `yaml:"-" json:"-"`
		LoadedConfig  Map     `yaml:"config" json:"config"`
	}

	// Map contains configuration
	Map map[string]interface{}

	// DefaultConfigModule is used to get a module's default configuration
	DefaultConfigModule interface {
		DefaultConfig() Map
	}

	// OverrideConfigModule allows to override config dynamically
	OverrideConfigModule interface {
		OverrideConfig(current Map) Map
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

// Add to the map (deep merge)
func (m Map) Add(cfg Map) {
	for k, v := range cfg {
		if vv, ok := v.(map[string]interface{}); ok {
			v = Map(vv)
		}

		if strings.Index(k, ".") > -1 {
			k, sub := strings.SplitN(k, ".", 2)[0], strings.SplitN(k, ".", 2)[1]
			if mm, ok := m[k]; !ok {
				m[k] = make(Map)
				m[k].(Map).Add(Map{sub: v})
			} else if mm, ok := mm.(Map); ok {
				mm.Add(Map{sub: v})
			} else {
				panic(fmt.Sprintf("Config conflict! %q.%q: %v into %v", k, sub, v, m[k]))
			}
		} else {
			_, mapleft := m[k].(Map)
			_, mapright := v.(Map)
			if mapleft && mapright {
				m[k].(Map).Add(v.(Map))
			} else if mapleft && !mapright {
				panic(fmt.Sprintf("Config conflict! %q:%v into %v", k, v, m[k]))
			} else if mapright {
				m[k] = make(Map)
				m[k].(Map).Add(v.(Map))
			} else {
				m[k] = v
			}
		}
	}
}

// Flat map
func (m Map) Flat() Map {
	res := make(Map)

	for k, v := range m {
		res[k] = v

		if v, ok := v.(Map); ok {
			for sk, sv := range v.Flat() {
				res[k+"."+sk] = sv
			}
		}
	}

	return res
}

// MarshalTo tries to marshal the configuration map into a given interface
func (m Map) MarshalTo(out interface{}) error {

	jsonBytes, err := json.Marshal(m)

	if err != nil {
		return errors.Wrap(err, "Problem with marshaling map")
	}

	err = json.Unmarshal(jsonBytes, &out)
	if err == nil {
		return errors.Wrap(err, fmt.Sprintf("Problem with unmarshaling into given structure %t", out))
	}

	return nil
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

	area.Configuration = make(Map)
	for _, module := range area.Modules {
		if cfgmodule, ok := module.(DefaultConfigModule); ok {
			area.Configuration.Add(cfgmodule.DefaultConfig())
		}
	}

	area.Configuration.Add(area.LoadedConfig)

	for _, module := range area.Modules {
		if cfgmodule, ok := module.(OverrideConfigModule); ok {
			area.Configuration.Add(cfgmodule.OverrideConfig(area.Configuration))
		}
	}

	regex := regexp.MustCompile(`%%ENV:([^%]+)%%`)
	for k, v := range area.Configuration.Flat() {
		if val, ok := v.(string); ok {
			v = regex.ReplaceAllStringFunc(val, func(a string) string { return os.Getenv(regex.FindStringSubmatch(a)[1]) })
		}
		if v == nil {
			log.Printf("Warning: %s has nil value Configured!", k)
			continue
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
		baseContext.Configuration = make(Map)
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
	if config, ok := area.Configuration.Flat()[key]; ok {
		return config
	}

	if area.Parent != nil {
		return area.Parent.Config(key)
	}

	return nil
}
