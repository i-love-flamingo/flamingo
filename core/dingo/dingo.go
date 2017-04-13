package dingo

import (
	"fmt"
	"log"
	"reflect"
)

type (
	Injector struct {
		bindings      map[reflect.Type][]*Binding
		multibindings map[reflect.Type][]*Binding
		parent        *Injector
		scopes        map[Scope]struct{}
	}

	Module interface {
		Configure(injector *Injector)
	}
)

func NewInjector(modules ...Module) *Injector {
	injector := &Injector{
		bindings:      make(map[reflect.Type][]*Binding),
		multibindings: make(map[reflect.Type][]*Binding),
		scopes:        make(map[Scope]struct{}),
	}

	injector.Bind(injector).ToInstance(injector)

	injector.BindScope(Singleton)

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

	return injector
}

// GetInstance always creates a pointer
func (injector *Injector) GetInstance(of interface{}) interface{} {
	oftype := reflect.TypeOf(of)

	if oft, ok := of.(reflect.Type); ok {
		oftype = oft
	} else {
		for oftype.Kind() == reflect.Ptr {
			oftype = oftype.Elem()
		}
	}

	log.Printf("resolving %s %#v\n", oftype, of)

	obj := injector.resolveType(oftype, "").Interface()

	injector.RequestInjection(obj)

	log.Printf("to %s %#v\n", reflect.TypeOf(obj), obj)

	return obj
}

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

func (injector *Injector) internalResolveType(t reflect.Type, annotation string) reflect.Value {
	if len(injector.bindings[t]) > 0 {
		binding := injector.lookupBinding(t, annotation)

		if binding.instance != nil {
			log.Println("instance", binding.instance.ivalue)
			return binding.instance.ivalue
		}

		if binding.provider != nil {
			log.Println("provider", binding.provider.fnctype.String(), binding.provider)
			result := binding.provider.Create(injector)
			if result.Kind() == reflect.Slice {
				result = injector.internalResolveType(result.Type(), "")
			} else {
				injector.RequestInjection(result.Interface())
			}
			return result
		}

		if binding.to != nil {
			log.Println("resolving to...")
			return injector.resolveType(binding.to, "")
		}
	}

	if injector.parent != nil {
		return injector.parent.resolveType(t, annotation)
	}

	if t.Kind() == reflect.Func {
		log.Printf("building %s\n", t)
		return reflect.MakeFunc(t, func(args []reflect.Value) (results []reflect.Value) {
			// create a new type
			res := reflect.New(t.Out(0))
			// dereference possible interface pointer
			if res.Kind() == reflect.Ptr && res.Elem().Kind() == reflect.Interface {
				res = res.Elem()
			}

			if res.Elem().Kind() == reflect.Slice {
				return []reflect.Value{injector.internalResolveType(t.Out(0), "")}
			} else {
				// set to actual value
				res.Set(reflect.ValueOf(injector.GetInstance(t.Out(0))))
				// return
				return []reflect.Value{res}
			}
		})
	}

	if t.Kind() == reflect.Slice {
		log.Println(injector.multibindings[t.Elem()])
		if bindings, ok := injector.multibindings[t.Elem()]; ok {
			n := reflect.MakeSlice(t, len(bindings), len(bindings))
			for _, binding := range bindings {
				n = reflect.Append(n, reflect.New(binding.to))
			}
			return n
		}
	}

	if t.Kind() == reflect.Interface {
		panic("Can not instantiate interface " + t.String())
	}

	return reflect.New(t)
}

func (injector *Injector) lookupBinding(t reflect.Type, annotation string) *Binding {
	for _, binding := range injector.bindings[t] {
		if binding.annotatedWith == annotation {
			return binding
		}
	}

	for _, binding := range injector.bindings[t] {
		if binding.annotatedWith == "" {
			return binding
		}
	}

	return injector.bindings[t][0]
}

func (injector *Injector) BindMulti(what interface{}) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	binding := new(Binding)
	binding.typeof = bindtype
	injector.multibindings[bindtype] = append(injector.bindings[bindtype], binding)
	return binding
}

func (injector *Injector) BindScope(s Scope) {
	injector.scopes[s] = struct{}{}
}

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

func (injector *Injector) RequestInjection(object interface{}) {
	var injectlist = []reflect.Value{reflect.ValueOf(object)}
	var i int

	for {
		if i >= len(injectlist) {
			break
		}

		current := injectlist[i]
		ctype := current.Type()
		log.Println("field:", ctype, current)

		i++

		switch ctype.Kind() {
		case reflect.Ptr:
			injectlist = append(injectlist, current.Elem())
			continue

		case reflect.Struct:
			for fieldIndex := 0; fieldIndex < ctype.NumField(); fieldIndex++ {
				if tag, ok := ctype.Field(fieldIndex).Tag.Lookup("inject"); ok {
					field := current.Field(fieldIndex)
					instance := injector.resolveType(field.Type(), tag)
					if instance.Kind() == reflect.Ptr {
						if instance.Elem().Kind() == reflect.Func || instance.Elem().Kind() == reflect.Interface {
							instance = instance.Elem()
						}
					}
					log.Printf("setting %v of %T to %v\n", current.Type().Field(fieldIndex).Name, current.Interface(), instance)
					field.Set(instance)
					injectlist = append(injectlist, instance)
				}
			}

		case reflect.Func:
			continue

		case reflect.Interface:
			continue

		default:
			//panic("Can't inject into " + fmt.Sprintf("%#v", current))
			continue
		}
	}
}
