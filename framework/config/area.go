// Package config provides supporting code for multi-tenant setups
package config

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"flamingo.me/dingo"
	"github.com/pkg/errors"
)

type (
	// Area defines a configuration area for multi-site setups
	// it is initialized by project main package and partly loaded by config files
	Area struct {
		Name string

		Parent   *Area
		Childs   []*Area
		Modules  []dingo.Module
		Injector *dingo.Injector

		Routes        []Route
		Configuration Map
		LoadedConfig  Map
	}

	// Map contains configuration
	Map map[string]interface{}

	// Slice contains a list of configuration options
	Slice []interface{}

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

// NewArea creates a new Area with optional childs
func NewArea(name string, modules []dingo.Module, childs ...*Area) *Area {
	ctx := &Area{
		Name:          name,
		Modules:       modules,
		Childs:        childs,
		Configuration: make(Map),
	}

	for _, c := range childs {
		c.Parent = ctx
	}

	return ctx
}

// GetFlatContexts returns a map of context-relative-name->*Area, which has been flatted to inherit all parent's
// tree settings such as DI & co, and filtered to only list tree nodes specified by Contexts of area.
func (area *Area) GetFlatContexts() ([]*Area, error) {
	var result []*Area
	flat, err := area.Flat()
	if err != nil {
		return nil, err
	}

	for relativeContextKey, context := range flat {
		result = append(result, &Area{
			Name:          relativeContextKey,
			Routes:        context.Routes,
			Injector:      context.Injector,
			Configuration: context.Configuration,
		})

	}
	return result, nil
}

// Add to the map (deep merge)
func (m Map) Add(cfg Map) error {
	// so we can not deep merge if we have `.` in our own keys, we need to ensure our keys are clean first
	for k, v := range m {
		var toClean Map
		if strings.Contains(k, ".") {
			if toClean == nil {
				toClean = make(Map)
			}
			toClean[k] = v
			delete(m, k)
		}
		if toClean != nil {
			if err := m.Add(toClean); err != nil {
				return err
			}
		}
	}

	for k, v := range cfg {
		if vv, ok := v.(map[string]interface{}); ok {
			v = Map(vv)
		} else if vv, ok := v.([]interface{}); ok {
			v = Slice(vv)
		}

		if strings.Index(k, ".") > -1 {
			k, sub := strings.SplitN(k, ".", 2)[0], strings.SplitN(k, ".", 2)[1]
			if mm, ok := m[k]; !ok {
				m[k] = make(Map)
				if err := m[k].(Map).Add(Map{sub: v}); err != nil {
					return err
				}
			} else if mm, ok := mm.(Map); ok {
				if err := mm.Add(Map{sub: v}); err != nil {
					return err
				}
			} else {
				return errors.Errorf("config conflict: %q.%q: %v into %v", k, sub, v, m[k])
			}
		} else {
			_, mapleft := m[k].(Map)
			_, mapright := v.(Map)
			// if left side already is a map and will be assigned to nil in config_dev
			if mapleft && v == nil {
				m[k] = nil
			} else if mapleft && mapright {
				if err := m[k].(Map).Add(v.(Map)); err != nil {
					return err
				}
			} else if mapleft && !mapright {
				return errors.Errorf("config conflict: %q:%v into %v", k, v, m[k])
			} else if mapright {
				m[k] = make(Map)
				if err := m[k].(Map).Add(v.(Map)); err != nil {
					return err
				}
			} else {
				// convert non-float64 to float64
				switch vv := v.(type) {
				case int:
					v = float64(vv)
				case int8:
					v = float64(vv)
				case int16:
					v = float64(vv)
				case int32:
					v = float64(vv)
				case int64:
					v = float64(vv)
				case uint:
					v = float64(vv)
				case uint8:
					v = float64(vv)
				case uint16:
					v = float64(vv)
				case uint32:
					v = float64(vv)
				case uint64:
					v = float64(vv)
				case float32:
					v = float64(vv)
				}
				m[k] = v
			}
		}
	}
	return nil
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

// MapInto tries to map the configuration map into a given interface
func (m Map) MapInto(out interface{}) error {
	jsonBytes, err := json.Marshal(m)

	if err != nil {
		return errors.Wrap(err, "Problem with marshaling map")
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Problem with unmarshaling into given structure %T", out))
	}

	return nil
}

// MapInto tries to map the configuration map into a given interface
func (s Slice) MapInto(out interface{}) error {
	jsonBytes, err := json.Marshal(s)

	if err != nil {
		return errors.Wrap(err, "Problem with marshaling map")
	}

	err = json.Unmarshal(jsonBytes, &out)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Problem with unmarshaling into given structure %T", out))
	}

	return nil
}

// Get a value by it's path
func (m Map) Get(key string) (interface{}, bool) {
	keyParts := strings.SplitN(key, ".", 2)
	val, ok := m[keyParts[0]]
	if len(keyParts) == 2 {
		mm, ok := val.(Map)
		if ok {
			return mm.Get(keyParts[1])
		}
		return mm, false
	}
	return val, ok
}

// GetInitializedInjector returns initialized container based on the configuration
// we derive our injector from our parent
func (area *Area) GetInitializedInjector() (*dingo.Injector, error) {
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
			if err := area.Configuration.Add(cfgmodule.DefaultConfig()); err != nil {
				return nil, err
			}
		}
	}

	if err := area.Configuration.Add(Map{"area": area.Name}); err != nil {
		return nil, err
	}
	if err := area.Configuration.Add(area.LoadedConfig); err != nil {
		return nil, err
	}

	for _, module := range area.Modules {
		if cfgmodule, ok := module.(OverrideConfigModule); ok {
			if err := area.Configuration.Add(cfgmodule.OverrideConfig(area.Configuration)); err != nil {
				return nil, err
			}
		}
	}

	for k, v := range area.Configuration.Flat() {
		if v == nil {
			continue
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}

	if config, ok := area.Configuration.Get("flamingo.modules.disabled"); ok {
		for _, disabled := range config.(Slice) {
			for i, module := range area.Modules {
				tm := reflect.TypeOf(module).Elem()
				if tm.PkgPath()+"."+tm.Name() == disabled.(string) {
					area.Modules = append(area.Modules[:i], area.Modules[i+1:]...)
				}
			}
		}
	}

	injector.InitModules(area.Modules...)

	return injector, nil
}

// Flat returns a map of name->*Area of contexts, were all values have been inherited (yet overriden) of the parent context tree.
func (area *Area) Flat() (map[string]*Area, error) {
	res := make(map[string]*Area)
	res[area.Name] = area
	var err error
	area.Injector, err = area.GetInitializedInjector()
	if err != nil {
		return nil, err
	}

	for _, child := range area.Childs {
		flat, err := child.Flat()
		if err != nil {
			return nil, err
		}
		for cn, flatchild := range flat {
			res[area.Name+`/`+cn] = MergeFrom(*flatchild, *area)
		}
	}

	return res, nil
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

// Config get a config value (recursive thru all parents if possible)
func (area *Area) Config(key string) (interface{}, bool) {
	if config, ok := area.Configuration.Get(key); ok {
		return config, true
	}

	if area.Parent != nil {
		return area.Parent.Config(key)
	}

	return nil, false
}

// HasConfigKey checks recursive if the config has a given key
func (area *Area) HasConfigKey(key string) bool {
	if _, ok := area.Configuration.Get(key); ok {
		return true
	}

	if area.Parent != nil {
		return area.Parent.HasConfigKey(key)
	}

	return false
}
