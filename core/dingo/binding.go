package dingo

import (
	"fmt"
	"reflect"
)

type (
	Binding struct {
		typeof reflect.Type

		to       reflect.Type
		instance *Instance
		provider *Provider

		// todo
		eager         bool
		annotatedWith string
		scope         Scope
	}

	Instance struct {
		itype  reflect.Type
		ivalue reflect.Value
	}

	Provider struct {
		fnctype reflect.Type
		fnc     reflect.Value
	}
)

func (b *Binding) To(what interface{}) *Binding {
	to := reflect.TypeOf(what)

	for to.Kind() == reflect.Ptr {
		to = to.Elem()
	}

	if !to.AssignableTo(b.typeof) && !reflect.PtrTo(to).AssignableTo(b.typeof) {
		panic(fmt.Sprintf("%s not assignable to %s", to, b.typeof))
	}

	b.to = to

	return b
}

func (b *Binding) ToInstance(instance interface{}) *Binding {
	b.instance = &Instance{
		itype:  reflect.TypeOf(instance),
		ivalue: reflect.ValueOf(instance),
	}
	return b
}

func (b *Binding) ToProvider(p interface{}) *Binding {
	provider := &Provider{
		fnc: reflect.ValueOf(p),
	}
	provider.fnctype = provider.fnc.Type().Out(0)
	if provider.fnctype != b.typeof && provider.fnctype != reflect.PtrTo(b.typeof) {
		panic(fmt.Sprintf("wrong provider type %s for %s", provider.fnctype, b.typeof))
	}
	b.provider = provider
	return b
}

func (b *Binding) AnnotatedWith(annotation string) *Binding {
	b.annotatedWith = annotation
	return b
}

func (b *Binding) In(scope Scope) *Binding {
	b.scope = scope
	return b
}

func (b *Binding) AsEagerSingleton() *Binding {
	b.In(Singleton)
	b.eager = true
	return b
}

func (p *Provider) Create(injector *Injector) reflect.Value {
	in := make([]reflect.Value, p.fnc.Type().NumIn())
	for i := 0; i < p.fnc.Type().NumIn(); i++ {
		in[i] = reflect.ValueOf(injector.GetInstance(p.fnc.Type().In(i)))
	}
	res := p.fnc.Call(in)[0]
	injector.RequestInjection(res)
	return res
}
