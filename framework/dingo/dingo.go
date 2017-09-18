package dingo

import (
	"fmt"
	"log"
	"reflect"
)

const (
	// INIT state
	INIT = iota
	// DEFAULT state
	DEFAULT
)

type (
	// Injector defines bindings and multibindings
	Injector struct {
		bindings      map[reflect.Type][]*Binding
		multibindings map[reflect.Type][]*Binding
		interceptor   map[reflect.Type][]reflect.Type
		overrides     []*override
		parent        *Injector
		scopes        map[reflect.Type]Scope
		stage         uint
		delayed       []interface{}
	}

	// Module is provided by packages to generate the DI tree
	Module interface {
		Configure(injector *Injector)
	}

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
		interceptor:   make(map[reflect.Type][]reflect.Type),
		scopes:        make(map[reflect.Type]Scope),
		stage:         DEFAULT,
	}

	injector.Bind(Injector{}).ToInstance(injector)

	injector.BindScope(Singleton)
	injector.BindScope(ChildSingleton)

	injector.InitModules(modules...)

	return injector
}

// Child derives a child injector with a new ChildSingletonScope
func (injector *Injector) Child() *Injector {
	var newInjector = NewInjector()
	newInjector.parent = injector
	newInjector.Bind(Injector{}).ToInstance(newInjector)
	newInjector.BindScope(new(ChildSingletonScope))
	newInjector.multibindings = injector.multibindings

	return newInjector
}

// InitModules initializes the injector with the given modules
func (injector *Injector) InitModules(modules ...Module) {
	injector.stage = INIT

	for _, module := range modules {
		injector.requestInjection(module)
		module.Configure(injector)
	}

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
		panic("cannot override unknown binding " + override.typ.String() + " (annotated with " + override.annotatedWith + ")")
	}

	for typ, bindings := range injector.bindings {
		known := make(map[string]*Binding)
		for _, binding := range bindings {
			if known, ok := known[binding.annotatedWith]; ok && !known.equal(binding) {
				panic(fmt.Sprintf("already known binding for %s with annotation '%s'", typ, binding.annotatedWith))
			}
			known[binding.annotatedWith] = binding
		}
	}

	injector.stage = DEFAULT

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

// getInstance creates the new instance of of, returns a reflect.value
func (injector *Injector) getInstance(of interface{}, annotatedWith string) reflect.Value {
	oftype := reflect.TypeOf(of)

	if oft, ok := of.(reflect.Type); ok {
		oftype = oft
	} else {
		for oftype.Kind() == reflect.Ptr {
			oftype = oftype.Elem()
		}
	}

	return injector.resolveType(oftype, annotatedWith)
}

func (injector *Injector) findBinding(t reflect.Type, annotation string) *Binding {
	if len(injector.bindings[t]) > 0 {
		binding := injector.lookupBinding(t, annotation)
		if binding != nil {
			return binding
		}
	}

	if injector.parent != nil {
		return injector.parent.findBinding(t, annotation)
	}

	return nil
}

// resolveType resolves a requested type, with annotation
func (injector *Injector) resolveType(t reflect.Type, annotation string) reflect.Value {
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
		final = injector.internalResolveType(t, annotation)
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

// internalResolveType resolves a type request with the current injector
func (injector *Injector) internalResolveType(t reflect.Type, annotation string) reflect.Value {
	if binding := injector.findBinding(t, annotation); binding != nil {
		if binding.instance != nil {
			return binding.instance.ivalue
		}

		if binding.provider != nil {
			result := binding.provider.Create(injector)
			if result.Kind() == reflect.Slice {
				result = injector.internalResolveType(result.Type(), "")
			} else {
				injector.requestInjection(result.Interface())
			}
			return result
		}

		if binding.to != nil {
			if binding.to == t {
				panic("circular from " + t.String() + " to " + binding.to.String() + " (annotated with: " + binding.annotatedWith + ")")
			}
			return injector.resolveType(binding.to, "")
		}
	}

	// This for an injection request on a provider, such as `func() MyInstance`
	if t.Kind() == reflect.Func && t.NumOut() == 1 {
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			// create a new type
			res := reflect.New(t.Out(0))
			// dereference possible interface pointer
			if res.Kind() == reflect.Ptr && (res.Elem().Kind() == reflect.Interface || res.Elem().Kind() == reflect.Ptr) {
				res = res.Elem()
			}

			if res.Kind() == reflect.Slice {
				return []reflect.Value{injector.internalResolveType(t.Out(0), annotation)}
			}
			// set to actual value
			res.Set(injector.getInstance(t.Out(0), annotation))
			// return
			return []reflect.Value{res}
		})
	}

	// This is the injection request for multibindings
	if t.Kind() == reflect.Slice {
		if bindings, ok := injector.multibindings[t.Elem()]; ok {
			n := reflect.MakeSlice(t, 0, len(bindings))
			for _, binding := range bindings {
				if binding.annotatedWith == annotation {
					//n = reflect.Append(n, injector.getInstance(binding.to))
					n = reflect.Append(n, injector.intercept(injector.getInstance(binding.to, annotation), t.Elem()))
				}
			}
			return n
		}
	}

	if t.Kind() == reflect.Interface {
		panic("Can not instantiate interface " + t.String())
	}

	if annotation != "" {
		panic("Can not automatically create an annotated injection " + t.String() + " with annotation " + annotation)
	}

	n := reflect.New(t)
	injector.requestInjection(n.Interface())
	return n
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
			log.Printf("%s: %s\n%s\n", current.Type().PkgPath(), current.Type().Name(), current.String())
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
			injectlist = append(injectlist, current.Elem())
			continue

		// inject into struct fields
		case reflect.Struct:
			for fieldIndex := 0; fieldIndex < ctype.NumField(); fieldIndex++ {
				if tag, ok := ctype.Field(fieldIndex).Tag.Lookup("inject"); ok {
					field := current.Field(fieldIndex)
					instance := injector.resolveType(field.Type(), tag)
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

	if injector.parent != nil {
		fmt.Printf("\n-----[ parent @ %p ]-----\n", injector.parent)
		injector.parent.Debug()
	}
}
