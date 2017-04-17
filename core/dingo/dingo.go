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

	for _, bindings := range injector.bindings {
		for _, binding := range bindings {
			if binding.eager {
				injector.GetInstance(binding.typeof)
			}
		}
	}

	return injector
}

func (injector *Injector) Child() *Injector {
	child := &Injector{
		bindings:      make(map[reflect.Type][]*Binding),
		multibindings: make(map[reflect.Type][]*Binding),
		scopes:        make(map[Scope]struct{}),
		parent:        injector,
	}

	return child
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

	obj := injector.resolveType(oftype, "").Interface()

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
				panic("circular for " + t.String() + "::" + binding.annotatedWith)
			}
			return injector.resolveType(binding.to, "")
		}
	}

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
				res.Set(reflect.ValueOf(injector.GetInstance(t.Out(0))))
				// return
				return []reflect.Value{res}
			}
		})
	}

	if t.Kind() == reflect.Slice {
		log.Println(injector.multibindings[t.Elem()])
		if bindings, ok := injector.multibindings[t.Elem()]; ok {
			n := reflect.MakeSlice(t, 0, len(bindings))
			for _, binding := range bindings {
				n = reflect.Append(n, reflect.ValueOf(injector.GetInstance(binding.to)))
			}
			return n
		}
	}

	if injector.parent != nil {
		return injector.parent.resolveType(t, annotation)
	}

	if t.Kind() == reflect.Interface {
		panic("Can not instantiate interface " + t.String())
	}

	n := reflect.New(t)
	injector.RequestInjection(n.Interface())
	return n
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

	panic("Can not find binding with annotation or empty for " + fmt.Sprintf("%T", t) + " " + annotation)
	//return injector.bindings[t][0]
}

func (injector *Injector) BindMulti(what interface{}) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	binding := &Binding{
		typeof: bindtype,
	}
	imb := injector.multibindings[bindtype]
	imb = append(imb, binding)
	injector.multibindings[bindtype] = imb
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

func (injector *Injector) Override(what interface{}) *Binding {
	bindtype := reflect.TypeOf(what)
	if bindtype.Kind() == reflect.Ptr {
		bindtype = bindtype.Elem()
	}
	if bindings, ok := injector.bindings[bindtype]; ok && len(bindings) > 0 {
		binding := new(Binding)
		injector.bindings[bindtype][0] = binding
		binding.typeof = bindtype
		return binding
	}
	panic("cannot override unknown binding")
}

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
		case reflect.Ptr:
			injectlist = append(injectlist, current.Elem())
			continue

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
					//log.Printf("setting %v of %T to %v\n", current.Type().Field(fieldIndex).Name, current.Interface(), instance)
					field.Set(instance)
				}
			}

		case reflect.Func:
			continue

		case reflect.Interface:
			continue

		case reflect.Slice:
			continue

		default:
			//panic("Can't inject into " + fmt.Sprintf("%#v", current))
			continue
		}
	}
}
