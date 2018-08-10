package dingo

import (
	"fmt"
	"reflect"
)

type (
	// Binding defines a type mapped to a more concrete type
	Binding struct {
		typeof reflect.Type

		to       reflect.Type
		instance *Instance
		provider *Provider

		eager         bool
		annotatedWith string
		scope         Scope
	}

	// Instance holds quick-references to type and value
	Instance struct {
		itype  reflect.Type
		ivalue reflect.Value
	}

	// Provider holds the provider function
	Provider struct {
		fnctype reflect.Type
		fnc     reflect.Value
		binding *Binding
	}
)

// To binds a concrete type to a binding
func (b *Binding) To(what interface{}) *Binding {
	to := reflect.TypeOf(what)

	for to.Kind() == reflect.Ptr {
		to = to.Elem()
	}

	if !to.AssignableTo(b.typeof) && !reflect.PtrTo(to).AssignableTo(b.typeof) {
		panic(fmt.Sprintf("%s#%s not assignable to %s#%s", to.PkgPath(), to.Name(), b.typeof.PkgPath(), b.typeof.Name()))
	}

	b.to = to

	return b
}

// ToInstance binds an instance to a binding
func (b *Binding) ToInstance(instance interface{}) *Binding {
	b.instance = &Instance{
		itype:  reflect.TypeOf(instance),
		ivalue: reflect.ValueOf(instance),
	}
	if !b.instance.itype.AssignableTo(b.typeof) && !b.instance.itype.AssignableTo(reflect.PtrTo(b.typeof)) {
		panic(fmt.Sprintf("%s#%s not assignable to %s#%s", b.instance.itype.PkgPath(), b.instance.itype.Name(), b.typeof.PkgPath(), b.typeof.Name()))
	}
	return b
}

// ToProvider binds a provider to an instance. The provider's arguments are automatically injected
func (b *Binding) ToProvider(p interface{}) *Binding {
	provider := &Provider{
		fnc:     reflect.ValueOf(p),
		binding: b,
	}
	provider.fnctype = provider.fnc.Type().Out(0)
	if !provider.fnctype.AssignableTo(b.typeof) && !provider.fnctype.AssignableTo(reflect.PtrTo(b.typeof)) {
		panic(fmt.Sprintf("provider returns %q which is not assignable to %q", provider.fnctype, b.typeof))
	}
	b.provider = provider
	return b
}

// AnnotatedWith sets the binding's annotation
func (b *Binding) AnnotatedWith(annotation string) *Binding {
	b.annotatedWith = annotation
	return b
}

// In set's the bindings scope
func (b *Binding) In(scope Scope) *Binding {
	b.scope = scope
	return b
}

// AsEagerSingleton set's the binding to singleton and requests eager initialization
func (b *Binding) AsEagerSingleton() *Binding {
	b.In(Singleton)
	b.eager = true
	return b
}

func (b *Binding) equal(to *Binding) bool {
	return reflect.DeepEqual(b, to)
}

// Create creates a new instance by the provider and requests injection, all provider arguments are automatically filled
func (p *Provider) Create(injector *Injector) reflect.Value {
	in := make([]reflect.Value, p.fnc.Type().NumIn())
	for i := 0; i < p.fnc.Type().NumIn(); i++ {
		in[i] = injector.getInstance(p.fnc.Type().In(i), "")
	}
	res := p.fnc.Call(in)[0]
	injector.requestInjection(res)
	return res
}
