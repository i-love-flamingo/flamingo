package dingo

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

const (
	// INIT state
	INIT = iota
	// DEFAULT state
	DEFAULT
)

type (
	// Injector defines bindings and multibindings
	// it is possible to have a parent-injector, which can be asked if no resolution is available
	Injector struct {
		bindings      map[reflect.Type][]*Binding          // list of available bindings for a concrete type
		multibindings map[reflect.Type][]*Binding          // list of multi-bindings for a concrete type
		mapbindings   map[reflect.Type]map[string]*Binding // list of map-bindings for a concrete type
		interceptor   map[reflect.Type][]reflect.Type      // list of interceptors for a type
		overrides     []*override                          // list of overrides for a binding
		parent        *Injector                            // parent injector reference
		scopes        map[reflect.Type]Scope               // scope-bindings
		stage         uint                                 // current stage
		delayed       []interface{}                        // delayed bindings
	}

	// Module is provided by packages to generate the DI tree
	Module interface {
		Configure(injector *Injector)
	}

	// Depender defines a dependency-aware module
	Depender interface {
		Depends() []Module
	}

	// overrides are evaluated lazy, so they are scheduled here
	override struct {
		typ           reflect.Type
		annotatedWith string
		binding       *Binding
	}
)

// NewInjector builds up a new Injector out of a list of Modules
func NewInjector(modules ...Module) *Injector {
	injector := &Injector{
		bindings:      make(map[reflect.Type][]*Binding),
		multibindings: make(map[reflect.Type][]*Binding),
		mapbindings:   make(map[reflect.Type]map[string]*Binding),
		interceptor:   make(map[reflect.Type][]reflect.Type),
		scopes:        make(map[reflect.Type]Scope),
		stage:         DEFAULT,
	}

	// bind current injector
	injector.Bind(Injector{}).ToInstance(injector)

	// bind default scopes
	injector.BindScope(Singleton)
	injector.BindScope(ChildSingleton)

	// init current modules
	injector.InitModules(modules...)

	return injector
}

// Child derives a child injector with a new ChildSingletonScope
func (injector *Injector) Child() *Injector {
	newInjector := NewInjector()
	newInjector.parent = injector
	newInjector.Bind(Injector{}).ToInstance(newInjector)
	newInjector.BindScope(new(ChildSingletonScope)) // bind a new child-singleton

	return newInjector
}

// InitModules initializes the injector with the given modules
func (injector *Injector) InitModules(modules ...Module) {
	injector.stage = INIT

	// todo better dependency resolution
	newModules := make([]Module, 0, len(modules))
	for _, module := range modules {
		if d, ok := module.(Depender); ok {
			newModules = append(newModules, d.Depends()...)
		}
		newModules = append(newModules, module)
	}
	modules = newModules

	known := make(map[reflect.Type]struct{}, len(modules))
	for _, module := range modules {
		if _, ok := known[reflect.TypeOf(module)]; ok {
			continue
		}
		injector.requestInjection(module)
		module.Configure(injector)
		known[reflect.TypeOf(module)] = struct{}{}
	}

	// evaluate overrides when modules were loaded
	for _, override := range injector.overrides {
		bindtype := override.typ
		if bindtype.Kind() == reflect.Ptr {
			bindtype = bindtype.Elem()
		}
		if bindings, ok := injector.bindings[bindtype]; ok && len(bindings) > 0 {
			for i, binding := range bindings {
				if binding.annotatedWith == override.annotatedWith {
					injector.bindings[bindtype][i] = override.binding
				}
			}
			continue
		}
		panic("cannot override unknown binding " + override.typ.String() + " (annotated with " + override.annotatedWith + ")") // todo ok?
	}

	// make sure there are no duplicated bindings
	for typ, bindings := range injector.bindings {
		known := make(map[string]*Binding)
		for _, binding := range bindings {
			if known, ok := known[binding.annotatedWith]; ok && !known.equal(binding) {
				panic(fmt.Sprintf("already known binding for %s with annotation '%s' | Known binding: %#v%#v Try %#v%#v", typ, binding.annotatedWith, known.to.PkgPath(), known.to.Name(), binding.to.PkgPath(), binding.to.Name()))
			}
			known[binding.annotatedWith] = binding
		}
	}

	injector.stage = DEFAULT

	// continue with delayed injections
	for _, object := range injector.delayed {
		injector.requestInjection(object)
	}

	injector.delayed = nil

	// build eager singletons
	for _, bindings := range injector.bindings {
		for _, binding := range bindings {
			if binding.eager {
				injector.getInstance(binding.typeof, binding.annotatedWith)
			}
		}
	}
}

// GetInstance creates a new instance of what was requested
func (injector *Injector) GetInstance(of interface{}) interface{} {
	return injector.getInstance(of, "").Interface()
}

// GetAnnotatedInstance creates a new instance of what was requested with the given annotation
func (injector *Injector) GetAnnotatedInstance(of interface{}, annotatedWith string) interface{} {
	return injector.getInstance(of, annotatedWith).Interface()
}

// getInstance creates the new instance of typ, returns a reflect.value
func (injector *Injector) getInstance(typ interface{}, annotatedWith string) reflect.Value {
	oftype := reflect.TypeOf(typ)

	if oft, ok := typ.(reflect.Type); ok {
		oftype = oft
	} else {
		for oftype.Kind() == reflect.Ptr {
			oftype = oftype.Elem()
		}
	}

	return injector.resolveType(oftype, annotatedWith, false)
}

func (injector *Injector) findBinding(t reflect.Type, annotation string) *Binding {
	if len(injector.bindings[t]) > 0 {
		binding := injector.lookupBinding(t, annotation)
		if binding != nil {
			return binding
		}
	}

	// inject one key of a map-binding
	if len(annotation) > 4 && annotation[:4] == "map:" {
		return injector.mapbindings[t][annotation[4:]]
	}

	// ask parent
	if injector.parent != nil {
		return injector.parent.findBinding(t, annotation)
	}

	return nil
}

// resolveType resolves a requested type, with annotation
func (injector *Injector) resolveType(t reflect.Type, annotation string, optional bool) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var final reflect.Value

	if binding := injector.findBinding(t, annotation); binding != nil {
		if binding.scope != nil {
			if scope, ok := injector.scopes[reflect.TypeOf(binding.scope)]; ok {
				final = scope.ResolveType(t, annotation, injector.internalResolveType)
				if !final.IsValid() {
					panic(fmt.Sprintf("%T did no resolve %s", scope, t))
				}
			} else {
				panic(fmt.Sprintf("unknown scope %T for %s", binding.scope, t))
			}
		}
	}

	if !final.IsValid() {
		final = injector.internalResolveType(t, annotation, optional)
	}

	if !final.IsValid() {
		panic("can not resolve " + t.String())
	}

	final = injector.intercept(final, t)

	return final
}

func (injector *Injector) intercept(final reflect.Value, t reflect.Type) reflect.Value {
	for _, interceptor := range injector.interceptor[t] {
		of := final
		final = reflect.New(interceptor)
		injector.requestInjection(final.Interface())
		final.Elem().Field(0).Set(of)
	}
	if injector.parent != nil {
		return injector.parent.intercept(final, t)
	}
	return final
}

func (injector *Injector) resolveBinding(binding *Binding, t reflect.Type, optional bool) (reflect.Value, error) {
	if binding.instance != nil {
		return binding.instance.ivalue, nil
	}

	if binding.provider != nil {
		result := binding.provider.Create(injector)
		if result.Kind() == reflect.Slice {
			result = injector.internalResolveType(result.Type(), "", optional)
		} else {
			injector.requestInjection(result.Interface())
		}
		return result, nil
	}

	if binding.to != nil {
		if binding.to == t {
			panic("circular from " + t.String() + " to " + binding.to.String() + " (annotated with: " + binding.annotatedWith + ")")
		}
		return injector.resolveType(binding.to, "", optional), nil
	}

	return reflect.Value{}, fmt.Errorf("binding is not bound: %v for %s", binding, t.String())
}

// internalResolveType resolves a type request with the current injector
func (injector *Injector) internalResolveType(t reflect.Type, annotation string, optional bool) reflect.Value {
	if binding := injector.findBinding(t, annotation); binding != nil {
		r, err := injector.resolveBinding(binding, t, optional)
		if err == nil {
			return r
		}

		if annotation != "" {
			return injector.resolveType(binding.typeof, "", false)
		}
	}

	// This for an injection request on a provider, such as `func() MyInstance`
	if t.Kind() == reflect.Func && t.NumOut() == 1 && strings.HasSuffix(t.Name(), "Provider") {
		return injector.createProvider(t, annotation, optional)
	}

	// This is the injection request for multibindings
	if t.Kind() == reflect.Slice {
		return injector.resolveMultibinding(t, annotation, optional)
	}

	// Map Binding injection
	if t.Kind() == reflect.Map && t.Key().Kind() == reflect.String {
		return injector.resolveMapbinding(t, annotation, optional)
	}

	if annotation != "" && !optional {
		panic("Can not automatically create an annotated injection " + t.String() + " with annotation " + annotation)
	}

	if t.Kind() == reflect.Interface && !optional {
		panic("Can not instantiate interface " + t.String())
	}

	if t.Kind() == reflect.Func && !optional {
		panic("Can not create a new function " + t.String() + " (Do you want a provider? Then suffix type with Provider)")
	}

	n := reflect.New(t)
	injector.requestInjection(n.Interface())
	return n
}

func (injector *Injector) createProvider(t reflect.Type, annotation string, optional bool) reflect.Value {
	return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
		// create a new type
		res := reflect.New(t.Out(0))
		// dereference possible interface pointer
		if res.Kind() == reflect.Ptr && (res.Elem().Kind() == reflect.Interface || res.Elem().Kind() == reflect.Ptr) {
			res = res.Elem()
		}

		// multibindings
		if res.Elem().Kind() == reflect.Slice {
			return []reflect.Value{injector.internalResolveType(t.Out(0), annotation, optional)}
		}

		// mapbindings
		if res.Elem().Kind() == reflect.Map && res.Elem().Type().Key().Kind() == reflect.String {
			return []reflect.Value{injector.internalResolveType(t.Out(0), annotation, optional)}
		}

		// set to actual value
		res.Set(injector.getInstance(t.Out(0), annotation))
		// return
		return []reflect.Value{res}
	})
}

func (injector *Injector) createProviderForBinding(t reflect.Type, binding *Binding, annotation string, optional bool) reflect.Value {
	return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
		// create a new type
		res := reflect.New(binding.typeof)
		// dereference possible interface pointer
		if res.Kind() == reflect.Ptr && (res.Elem().Kind() == reflect.Interface || res.Elem().Kind() == reflect.Ptr) {
			res = res.Elem()
		}

		if r, err := injector.resolveBinding(binding, t, optional); err == nil {
			res.Set(r)
			return []reflect.Value{res}
		}

		// set to actual value
		res.Set(injector.getInstance(binding.typeof, annotation))
		// return
		return []reflect.Value{res}
	})
}

func (injector *Injector) joinMultibindings(t reflect.Type, annotation string) []*Binding {
	var parent []*Binding
	if injector.parent != nil {
		parent = injector.parent.joinMultibindings(t, annotation)
	}

	bindings := make([]*Binding, len(parent)+len(injector.multibindings[t]))
	copy(bindings, parent)
	c := len(parent)
	for _, b := range injector.multibindings[t] {
		if b.annotatedWith == annotation {
			bindings[c] = b
			c++
		}
	}
	return bindings[:c]
}

func (injector *Injector) resolveMultibinding(t reflect.Type, annotation string, optional bool) reflect.Value {
	targetType := t.Elem()
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	providerType := targetType
	provider := strings.HasSuffix(targetType.Name(), "Provider") && targetType.Kind() == reflect.Func

	if provider {
		targetType = targetType.Out(0)
	}

	if bindings := injector.joinMultibindings(targetType, annotation); len(bindings) > 0 {
		n := reflect.MakeSlice(t, 0, len(bindings))
		for _, binding := range bindings {
			if provider {
				n = reflect.Append(n, injector.createProviderForBinding(providerType, binding, annotation, false))
				continue
			}

			r, err := injector.resolveBinding(binding, t, optional)
			if err != nil {
				panic(err)
			}
			n = reflect.Append(n, r)
		}
		return n
	}

	return reflect.MakeSlice(t, 0, 0)
}

func (injector *Injector) joinMapbindings(t reflect.Type, annotation string) map[string]*Binding {
	var parent map[string]*Binding
	if injector.parent != nil {
		parent = injector.parent.joinMapbindings(t, annotation)
	}

	bindings := make(map[string]*Binding, len(parent)+len(injector.multibindings[t]))
	for k, v := range parent {
		bindings[k] = v
	}
	for k, v := range injector.mapbindings[t] {
		if v.annotatedWith == annotation {
			bindings[k] = v
		}
	}
	return bindings
}

func (injector *Injector) resolveMapbinding(t reflect.Type, annotation string, optional bool) reflect.Value {
	targetType := t.Elem()
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	providerType := targetType
	provider := strings.HasSuffix(targetType.Name(), "Provider") && targetType.Kind() == reflect.Func

	if provider {
		targetType = targetType.Out(0)
	}

	if bindings := injector.joinMapbindings(targetType, annotation); len(bindings) > 0 {
		n := reflect.MakeMapWithSize(t, len(bindings))
		for key, binding := range bindings {
			if provider {
				n.SetMapIndex(reflect.ValueOf(key), injector.createProviderForBinding(providerType, binding, annotation, false))
				continue
			}

			r, err := injector.resolveBinding(binding, t, optional)
			if err != nil {
				panic(err)
			}
			n.SetMapIndex(reflect.ValueOf(key), r)
		}
		return n
	}

	return reflect.MakeMap(t)
}

// lookupBinding search a binding with the corresponding annotation
func (injector *Injector) lookupBinding(t reflect.Type, annotation string) *Binding {
	for _, binding := range injector.bindings[t] {
		if binding.annotatedWith == annotation {
			return binding
		}
	}

	return nil
}

// BindMulti binds multiple concrete types to the same abstract type / interface
func (injector *Injector) BindMulti(what interface{}) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	binding := new(Binding)
	binding.typeof = bindtype
	imb := injector.multibindings[bindtype]
	imb = append(imb, binding)
	injector.multibindings[bindtype] = imb
	return binding
}

// BindMap does a registry-like map-based binding, like BindMulti
func (injector *Injector) BindMap(what interface{}, key string) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	binding := new(Binding)
	binding.typeof = bindtype
	bindingMap := injector.mapbindings[bindtype]
	if bindingMap == nil {
		bindingMap = make(map[string]*Binding)
	}
	bindingMap[key] = binding
	injector.mapbindings[bindtype] = bindingMap

	return binding
}

// BindInterceptor intercepts to interface with interceptor
func (injector *Injector) BindInterceptor(to, interceptor interface{}) {
	totype := reflect.TypeOf(to)
	if totype.Kind() == reflect.Ptr {
		totype = totype.Elem()
	}
	if totype.Kind() != reflect.Interface {
		panic("can only intercept interfaces " + fmt.Sprintf("%v", to))
	}
	m := injector.interceptor[totype]
	m = append(m, reflect.TypeOf(interceptor))
	injector.interceptor[totype] = m
}

// BindScope binds a scope to be aware of
func (injector *Injector) BindScope(s Scope) {
	injector.scopes[reflect.TypeOf(s)] = s
}

// Bind creates a new binding for an abstract type / interface
// Use the syntax
//	injector.Bind((*Interface)(nil))
// To specify the interface (cast it to a pointer to a nil of the type Interface)
func (injector *Injector) Bind(what interface{}) *Binding {
	if what == nil {
		panic("Cannot bind nil")
	}
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	binding := new(Binding)
	binding.typeof = bindtype
	injector.bindings[bindtype] = append(injector.bindings[bindtype], binding)
	return binding
}

// Override a binding
func (injector *Injector) Override(what interface{}, annotatedWith string) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	var binding = new(Binding)
	binding.typeof = bindtype

	injector.overrides = append(injector.overrides, &override{typ: bindtype, annotatedWith: annotatedWith, binding: binding})

	return binding
}

// RequestInjection requests the object to have all fields annotated with `inject` to be filled
func (injector *Injector) RequestInjection(object interface{}) {
	if injector.stage == INIT {
		injector.delayed = append(injector.delayed, object)
	} else {
		injector.requestInjection(object)
	}
}

func (injector *Injector) requestInjection(object interface{}) {
	if _, ok := object.(reflect.Value); !ok {
		object = reflect.ValueOf(object)
	}
	var injectlist = []reflect.Value{object.(reflect.Value)}
	var i int
	var current reflect.Value

	defer func() {
		if e := recover(); e != nil {
			log.Printf("%s: %s: injecting into %s", current.Type().PkgPath(), current.Type().Name(), current.String())
			panic(e)
		}
	}()

	for {
		if i >= len(injectlist) {
			break
		}

		current = injectlist[i]
		ctype := current.Type()

		i++

		switch ctype.Kind() {
		// dereference pointer
		case reflect.Ptr:
			if setup := current.MethodByName("Inject"); setup.IsValid() {
				args := make([]reflect.Value, setup.Type().NumIn())
				for i := range args {
					args[i] = injector.getInstance(setup.Type().In(i), "")
				}
				setup.Call(args)
			}
			injectlist = append(injectlist, current.Elem())
			continue

		// inject into struct fields
		case reflect.Struct:
			for fieldIndex := 0; fieldIndex < ctype.NumField(); fieldIndex++ {
				if tag, ok := ctype.Field(fieldIndex).Tag.Lookup("inject"); ok {
					field := current.Field(fieldIndex)

					if field.Kind() == reflect.Struct {
						panic(fmt.Sprintf("Can not inject into struct %#v of %#v", field, current))
					}

					var optional bool
					for _, option := range strings.Split(tag, ",") {
						switch strings.TrimSpace(option) {
						case "optional":
							optional = true
						}
					}
					tag = strings.Split(tag, ",")[0]

					instance := injector.resolveType(field.Type(), tag, optional)
					if instance.Kind() == reflect.Ptr {
						if instance.Elem().Kind() == reflect.Func || instance.Elem().Kind() == reflect.Interface || instance.Elem().Kind() == reflect.Slice {
							instance = instance.Elem()
						}
					}
					if field.Kind() != reflect.Ptr && field.Kind() != reflect.Interface && instance.Kind() == reflect.Ptr {
						field.Set(instance.Elem())
					} else {
						field.Set(instance)
					}
				}
			}

		default:
			continue
		}
	}
}

// Debug Output
func (injector *Injector) Debug() {
	fmt.Print("My scopes: ")
	for _, scope := range injector.scopes {
		fmt.Printf("%T@%p ", scope, scope)
	}
	fmt.Println()
	fmt.Println()

	for vtype, bindings := range injector.bindings {
		fmt.Printf("\t%30s  >  ", vtype)
		for _, binding := range bindings {
			if binding.annotatedWith != "" {
				fmt.Printf(" (%s)", binding.annotatedWith)
			}
			if binding.instance != nil {
				fmt.Printf(" %s |", binding.instance.ivalue.String())
			} else if binding.provider != nil {
				fmt.Printf(" %s |", binding.provider.fnc.String())
			} else if binding.to != nil {
				fmt.Printf(" %s |", binding.to)
			}
		}
		fmt.Println()
	}

	for vtype, bindings := range injector.multibindings {
		fmt.Printf("\t%30s  >  ", vtype)
		for _, binding := range bindings {
			if binding.annotatedWith != "" {
				fmt.Printf(" (%s)", binding.annotatedWith)
			}
			if binding.instance != nil {
				fmt.Printf(" %s |", binding.instance.ivalue.String())
			} else if binding.provider != nil {
				fmt.Printf(" %s |", binding.provider.fnc.String())
			} else if binding.to != nil {
				fmt.Printf(" %s |", binding.to)
			}
		}
		fmt.Println()
	}

	for vtype, bindings := range injector.mapbindings {
		fmt.Printf("\t%30s  >  ", vtype)
		for key, binding := range bindings {
			fmt.Printf("%s:", key)
			if binding.annotatedWith != "" {
				fmt.Printf(" (%s)", binding.annotatedWith)
			}
			if binding.instance != nil {
				fmt.Printf(" %s |", binding.instance.ivalue.String())
			} else if binding.provider != nil {
				fmt.Printf(" %s |", binding.provider.fnc.String())
			} else if binding.to != nil {
				fmt.Printf(" %s |", binding.to)
			}
		}
		fmt.Println()
	}

	if injector.parent != nil {
		fmt.Printf("\n-----[ parent @ %p ]-----\n", injector.parent)
		injector.parent.Debug()
	}
}
