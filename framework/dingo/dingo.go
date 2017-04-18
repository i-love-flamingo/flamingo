package dingo

// TODO: Implement child/parent

import (
	"fmt"
	"reflect"
)

type (
	// Injector defines bindings and multibindings
	Injector struct {
		bindings      map[reflect.Type][]*Binding
		multibindings map[reflect.Type][]*Binding
		//parent        *Injector
		scopes map[Scope]struct{}
	}

	// Module is provided by packages to generate the DI tree
	Module interface {
		Configure(injector *Injector)
	}
)

// NewInjector builds up a new Injector out of a list of Modules
func NewInjector(modules ...Module) *Injector {
	injector := &Injector{
		bindings:      make(map[reflect.Type][]*Binding),
		multibindings: make(map[reflect.Type][]*Binding),
		scopes:        make(map[Scope]struct{}),
	}

	injector.Bind(Injector{}).ToInstance(injector)

	injector.BindScope(Singleton)

	injector.InitModules(modules...)

	return injector
}

// InitModules initializes the injector with the given modules
func (injector *Injector) InitModules(modules ...Module) {
	for _, module := range modules {
		injector.RequestInjection(module)
		module.Configure(injector)
	}

	for typ, bindings := range injector.bindings {
		known := make(map[string]struct{})
		for _, binding := range bindings {
			if _, ok := known[binding.annotatedWith]; ok {
				panic(fmt.Sprintf("already known binding for %s with annotation %s", typ, binding.annotatedWith))
			}
			known[binding.annotatedWith] = struct{}{}
		}
	}

	// build eager singletons
	for _, bindings := range injector.bindings {
		for _, binding := range bindings {
			if binding.eager {
				injector.getInstance(binding.typeof)
			}
		}
	}
}

// GetInstance creates a new instance of what was requested
func (injector *Injector) GetInstance(of interface{}) interface{} {
	return injector.getInstance(of).Interface()
}

// getInstance creates the new instance of of, returns a reflect.value
func (injector *Injector) getInstance(of interface{}) reflect.Value {
	oftype := reflect.TypeOf(of)

	if oft, ok := of.(reflect.Type); ok {
		oftype = oft
	} else {
		for oftype.Kind() == reflect.Ptr {
			oftype = oftype.Elem()
		}
	}

	return injector.resolveType(oftype, "")
}

// resolveType resolves a requested type, with annotation
func (injector *Injector) resolveType(t reflect.Type, annotation string) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if len(injector.bindings[t]) > 0 {
		binding := injector.lookupBinding(t, annotation)

		if binding.scope != nil {
			if _, ok := injector.scopes[binding.scope]; ok {
				return binding.scope.resolveType(t, injector.internalResolveType)
			} else {
				panic(fmt.Sprintf("unknown scope %s", binding.scope))
			}
		}
	}

	return injector.internalResolveType(t, annotation)
}

// internalResolveType resolves a type request with the current injector
func (injector *Injector) internalResolveType(t reflect.Type, annotation string) reflect.Value {
	if len(injector.bindings[t]) > 0 {
		binding := injector.lookupBinding(t, annotation)

		if binding.instance != nil {
			return binding.instance.ivalue
		}

		if binding.provider != nil {
			result := binding.provider.Create(injector)
			if result.Kind() == reflect.Slice {
				result = injector.internalResolveType(result.Type(), "")
			} else {
				injector.RequestInjection(result.Interface())
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
				return []reflect.Value{injector.internalResolveType(t.Out(0), "")}
			} else {
				// set to actual value
				res.Set(injector.getInstance(t.Out(0)))
				// return
				return []reflect.Value{res}
			}
		})
	}

	// This is the injection request for multibindings
	if t.Kind() == reflect.Slice {
		if bindings, ok := injector.multibindings[t.Elem()]; ok {
			n := reflect.MakeSlice(t, 0, len(bindings))
			for _, binding := range bindings {
				if binding.annotatedWith == annotation {
					n = reflect.Append(n, injector.getInstance(binding.to))
				}
			}
			return n
		}
	}

	//if injector.parent != nil {
	//	return injector.parent.resolveType(t, annotation)
	//}

	if t.Kind() == reflect.Interface {
		panic("Can not instantiate interface " + t.String())
	}

	n := reflect.New(t)
	injector.RequestInjection(n.Interface())
	return n
}

// lookupBinding search a binding with the corresponding annotation
func (injector *Injector) lookupBinding(t reflect.Type, annotation string) *Binding {
	for _, binding := range injector.bindings[t] {
		if binding.annotatedWith == annotation {
			return binding
		}
	}

	//for _, binding := range injector.bindings[t] {
	//	if binding.annotatedWith == "" {
	//		return binding
	//	}
	//}

	panic("Can not find binding with annotation '" + annotation + "' for " + fmt.Sprintf("%T", t))
	//return injector.bindings[t][0]
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

// BindScope binds a scope to be aware of
func (injector *Injector) BindScope(s Scope) {
	injector.scopes[s] = struct{}{}
}

// Bind creates a new binding for an abstract type / interface
// Use the syntax
//	injector.Bind((*Interface)(nil))
// To specify the interface (cast it to a pointer to a nil of the type Interface)
func (injector *Injector) Bind(what interface{}) *Binding {
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
	if bindings, ok := injector.bindings[bindtype]; ok && len(bindings) > 0 {
		for i, binding := range bindings {
			if binding.annotatedWith == annotatedWith {
				binding := new(Binding)
				injector.bindings[bindtype][i] = binding
				binding.typeof = bindtype
				return binding
			}
		}
	}
	panic("cannot override unknown binding (annotated with " + annotatedWith + ")")
}

// RequestInjection requests the object to have all fields annotated with `inject` to be filled
func (injector *Injector) RequestInjection(object interface{}) {
	if _, ok := object.(reflect.Value); !ok {
		object = reflect.ValueOf(object)
	}
	var injectlist = []reflect.Value{object.(reflect.Value)}
	var i int

	for {
		if i >= len(injectlist) {
			break
		}

		current := injectlist[i]
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
	for vtype, bindings := range injector.bindings {
		fmt.Printf("   %s     >    ", vtype)
		for _, binding := range bindings {
			fmt.Printf("%s |", binding.to)
		}
		fmt.Println()
	}

	for vtype, bindings := range injector.multibindings {
		fmt.Printf("   %s     >    ", vtype)
		for _, binding := range bindings {
			fmt.Printf("%s |", binding.to)
		}
		fmt.Println()
	}
}
