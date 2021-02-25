// Package config provides supporting code for multi-tenant setups
package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"flamingo.me/dingo"
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
		LoadedConfig  Map // Deprecated: empty and should not be used anymore

		cueBuildInstance *build.Instance
		cueInstance      *cue.Instance
		cueConfig        *ast.File
		defaultConfig    Map
		loadedConfig     Map
	}

	// DefaultConfigModule is used to get a module's default configuration
	// Deprecated: use CueConfigModule instead
	DefaultConfigModule interface {
		DefaultConfig() Map
	}

	// CueConfigModule provides a cue schema with default configuration which is used to validate and set config
	CueConfigModule interface {
		CueConfig() string
	}

	// OverrideConfigModule allows to override config dynamically
	OverrideConfigModule interface {
		OverrideConfig(current Map) Map
	}

	flamingoLegacyConfigAlias interface {
		FlamingoLegacyConfigAlias() map[string]string
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
// Deprecated: just do it yourself if necessary, with Flat()
func (area *Area) GetFlatContexts() ([]*Area, error) {
	var result []*Area
	flat, err := area.Flat()
	if err != nil {
		return nil, err
	}

	for relativeContextKey, context := range flat {
		newArea := *context
		newArea.Name = relativeContextKey
		result = append(result, &newArea)
	}
	return result, nil
}

var typeOfModuleFunc = reflect.TypeOf(dingo.ModuleFunc(nil))

// resolveDependencies tries to get a complete list of all modules, including all dependencies
// known can be empty initially, and will then be used for subsequent recursive calls
func resolveDependencies(modules []dingo.Module, known map[interface{}]struct{}) []dingo.Module {
	final := make([]dingo.Module, 0, len(modules))

	if known == nil {
		known = make(map[interface{}]struct{})
	}

	for _, module := range modules {
		var identity interface{} = reflect.TypeOf(module)
		if identity == typeOfModuleFunc {
			identity = reflect.ValueOf(module)
		}
		if _, ok := known[identity]; ok {
			continue
		}
		known[identity] = struct{}{}
		if depender, ok := module.(dingo.Depender); ok {
			final = append(final, resolveDependencies(depender.Depends(), known)...)
		}
		final = append(final, module)
	}

	return final
}

func moduleName(m dingo.Module) string {
	tm := reflect.TypeOf(m)
	for tm.Kind() == reflect.Ptr {
		tm = tm.Elem()
	}
	return tm.PkgPath() + "." + tm.Name()
}

func cueError(err error) error {
	if p, ok := err.(errors.Error); ok {
		return fmt.Errorf("%s: %w", p.Position(), err)
	}
	return err
}

func (area *Area) loadCueConfig() error {
	area.Modules = resolveDependencies(area.Modules, nil)

	if err := area.cueBuildInstance.AddFile("flamingo.modules.disabled", "flamingo?: modules?: disabled?: [...string]"); err != nil {
		return cueError(err)
	}

	for _, module := range area.Modules {
		if cuemodule, ok := module.(CueConfigModule); ok {
			if err := area.cueBuildInstance.AddFile(moduleName(module), cuemodule.CueConfig()); err != nil {
				return fmt.Errorf("loading config for %s failed: %w", moduleName(module), cueError(err))
			}
		}
	}

	envFile := "flamingo: { os: { env: { \n[string]: string\n"
	for _, v := range os.Environ() {
		v := strings.SplitN(v, "=", 2)
		envFile += fmt.Sprintf("\"%s\": \"%s\"\n", esc(v[0]), esc(v[1]))
	}
	envFile += "} } }\n"

	if err := area.cueBuildInstance.AddFile("flamingo-os-env-file", envFile); err != nil {
		return fmt.Errorf("%s: %w", area.Name, cueError(err))
	}

	return nil
}

func (area *Area) loadDefaultConfig() error {
	area.Modules = resolveDependencies(area.Modules, nil)
	area.defaultConfig = make(Map)

	for _, module := range area.Modules {
		if cfgmodule, ok := module.(DefaultConfigModule); ok {
			if err := area.defaultConfig.Add(cfgmodule.DefaultConfig()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (area *Area) checkLegacyConfig(warn bool) {
	for _, module := range area.Modules {
		if cfgmodule, ok := module.(flamingoLegacyConfigAlias); ok {
			for old, new := range cfgmodule.FlamingoLegacyConfigAlias() {
				if oldval, ok := area.Configuration.Get(old); ok {
					if warn {
						log.Printf("WARNING: legacy config %q set, migrate to %q", old, new)
					}
					if newval, ok := area.Configuration.Get(new); !ok {
						if err := area.Configuration.Add(Map{new: oldval}); err != nil {
							log.Fatal(err)
						}
					} else if ok && !reflect.DeepEqual(oldval, newval) {
						// don't warn on complex/map type
						if _, ok := newval.(Map); !ok {
							log.Fatalf("ERROR: legacy config mismatch for new %q=%q and old %q=%q", new, newval, old, oldval)
						}
					}
				}
				if newval, ok := area.Configuration.Get(new); ok {
					if err := area.Configuration.Add(Map{old: newval}); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}
}

func esc(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", `\r`, -1)
	return s
}

func (area *Area) loadConfig(legacy, logLegacy bool) error {
	area.cueBuildInstance = build.NewContext().NewInstance(area.Name, nil)

	if err := area.loadCueConfig(); err != nil {
		return err
	}

	if err := area.loadDefaultConfig(); err != nil {
		return err
	}

	area.Configuration = Map{"area": area.Name}

	if err := area.Configuration.Add(area.defaultConfig); err != nil {
		return err
	}
	if err := area.Configuration.Add(area.loadedConfig); err != nil {
		return err
	}

	for _, module := range area.Modules {
		if cfgmodule, ok := module.(OverrideConfigModule); ok {
			if err := area.Configuration.Add(cfgmodule.OverrideConfig(area.Configuration)); err != nil {
				return err
			}
		}
	}

	if legacy {
		area.checkLegacyConfig(logLegacy)
	}

	// TODO workaround for issue https://github.com/cuelang/cue/issues/220
	// we mark every nil-value as `*null | _`, which includes everything but defaults to null
	purgeNil := ""
	for k, v := range area.Configuration.Flat() {
		if v == nil {
			purgeNil += `"` + strings.Replace(k, `.`, `": "`, -1) + `": *null | _` + "\n"
		}
	}

	if err := area.cueBuildInstance.AddFile("flamingo.config.purgenil", purgeNil); err != nil {
		return cueError(err)
	}
	if area.cueConfig != nil {
		if err := area.cueBuildInstance.AddSyntax(area.cueConfig); err != nil {
			return cueError(err)
		}
	}

	var err error
	area.cueInstance, err = new(cue.Runtime).Build(area.cueBuildInstance)
	if err != nil {
		return cueError(err)
	}

	area.cueInstance, err = area.cueInstance.Fill(area.Configuration)
	if err != nil {
		return fmt.Errorf("%s: %w", area.Name, cueError(err))
	}

	m := make(Map)
	if err := area.cueInstance.Value().Decode(&m); err != nil {
		return fmt.Errorf("%s: %w", area.Name, cueError(err))
	}
	if err := area.Configuration.Add(m); err != nil {
		return fmt.Errorf("%s: %w", area.Name, err)
	}

	if legacy {
		area.checkLegacyConfig(false)
	}

	return nil
}

// GetInitializedInjector returns initialized container based on the configuration
// we derive our injector from our parent
func (area *Area) GetInitializedInjector() (*dingo.Injector, error) {
	if area.Injector != nil {
		return area.Injector, nil
	}

	if area.Parent != nil {
		parent, err := area.Parent.GetInitializedInjector()
		if err != nil {
			return nil, err
		}
		if area.Injector, err = parent.Child(); err != nil {
			return nil, err
		}
	} else {
		var err error
		if area.Injector, err = dingo.NewInjector(); err != nil {
			return nil, err
		}
	}
	area.Injector.SetBuildEagerSingletons(false)
	area.Injector.Bind(Area{}).ToInstance(area)

	if err := area.loadConfig(true, true); err != nil {
		return nil, err
	}

	for k, v := range area.Configuration.Flat() {
		if v == nil {
			continue
		}
		area.Injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
		if vf, ok := v.(float64); ok && vf == float64(int64(vf)) {
			area.Injector.Bind(new(int64)).AnnotatedWith("config:" + k).ToInstance(int64(vf))
			area.Injector.Bind(new(int)).AnnotatedWith("config:" + k).ToInstance(int(int64(vf)))
		}
	}

	if config, ok := area.Configuration.Get("flamingo.modules.disabled"); ok {
		for _, disabled := range config.(Slice) {
			area.Modules = disableModule(area.Modules, disabled.(string))
		}
	}

	if err := area.Injector.InitModules(area.Modules...); err != nil {
		return nil, err
	}
	return area.Injector, nil // area.Injector.BuildEagerSingletons(false)
}

func disableModule(input []dingo.Module, disabled string) []dingo.Module {
	for i, module := range input {
		tm := reflect.TypeOf(module).Elem()
		if tm.PkgPath()+"."+tm.Name() == disabled {
			return append(input[:i], input[i+1:]...)
		}
	}
	return input
}

// Flat returns a map of name->*Area of contexts, were all values have been inherited (yet overridden) of the parent context tree.
func (area *Area) Flat() (map[string]*Area, error) {
	res := make(map[string]*Area, 1+len(area.Childs))
	res[area.Name] = area

	for _, child := range area.Childs {
		flat, err := child.Flat()
		if err != nil {
			return nil, err
		}
		for cn, flatchild := range flat {
			res[area.Name+`/`+cn] = MergeFrom(*flatchild, *area)
			_ = res[area.Name+`/`+cn].loadConfig(true, false) // we load the config as far as possible
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
