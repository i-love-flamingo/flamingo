package service_container

import (
	"log"
	"reflect"
	"runtime"
)

type (
	// ServiceContainer is a basic flamingo helper
	// to register default Routes, packages, etc.
	ServiceContainer struct {
		unnamed []*Object
		named   map[string]*Object
		Routes  map[string]string
		Handler map[string]interface{}
		di      *Graph
	}

	// RegisterFunc defines a callback used by packages to bootstrap themselves
	RegisterFunc func(r *ServiceContainer)

	// PostInjecter defines the PostInject() function which is called when the DI has finished
	PostInjecter interface {
		PostInject(g *Graph)
	}
)

// MarshalText serialization of RegisterFunc Names
func (r RegisterFunc) MarshalText() (text []byte, err error) {
	return []byte(runtime.FuncForPC(reflect.ValueOf(r).Pointer()).Name()), nil
}

// New creates a new ServiceContainer
func New() *ServiceContainer {
	return &ServiceContainer{
		Routes:  make(map[string]string),
		named:   make(map[string]*Object),
		Handler: make(map[string]interface{}),
	}
}

// WalkRegisterFuncs calls the provided RegisterFunc callbacks
func (r *ServiceContainer) WalkRegisterFuncs(rfs ...RegisterFunc) *ServiceContainer {
	for _, rf := range rfs {
		rf(r)
	}
	return r
}

// Handle registers Handler on ServiceContainer
func (r *ServiceContainer) Handle(name string, handler interface{}) {
	r.Handler[name] = handler
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		r.Register(handler)
	}
}

// Route adds a route
func (r *ServiceContainer) Route(path, name string) *ServiceContainer {
	r.Routes[path] = name
	return r
}

// Register registers any object for DI
func (r *ServiceContainer) Register(o interface{}, tags ...string) *ServiceContainer {
	r.Remove(o)

	object := &Object{
		Value: o,
		Tags:  tags,
	}
	r.unnamed = append(r.unnamed, object)
	return r
}

// RegisterNamed registers any object for DI with a given name
func (r *ServiceContainer) RegisterNamed(name string, o interface{}, tags ...string) *ServiceContainer {
	object := &Object{
		Value: o,
		Name:  name,
		Tags:  tags,
	}
	r.named[name] = object
	return r
}

// Remove removes an already registered object of the same type
func (r *ServiceContainer) Remove(i interface{}) {
	for k, o := range r.unnamed {
		if reflect.TypeOf(o.Value).String() == reflect.TypeOf(i).String() {
			r.unnamed = append(r.unnamed[:k], r.unnamed[k+1:]...)
		}
	}
}

// sl is a private logger to show DI logs
type sl struct{}

// Debugf DI logger
func (s sl) Debugf(a string, b ...interface{}) {
	log.Printf(a, b...)
}

// DI returns the injection graph, not populated
func (r *ServiceContainer) DI() *Graph {
	if r.di != nil {
		return r.di
	}

	r.di = new(Graph)
	//r.di.Logger = sl{}

	r.Register(r)

	for _, o := range r.unnamed {
		err := r.di.Provide(o)
		if err != nil {
			panic(err)
		}
	}

	for _, o := range r.named {
		err := r.di.Provide(o)
		if err != nil {
			panic(err)
		}
	}

	return r.di
}

// Resolve populates the injection graph
func (r *ServiceContainer) Resolve() {
	di := r.DI()
	err := di.Populate()
	if err != nil {
		panic(err)
	}

	for _, o := range di.Objects() {
		if pi, ok := o.Value.(PostInjecter); ok && !o.PostInjected {
			pi.PostInject(di)
			o.PostInjected = true
		}
	}
}

// InjectInto injects resolves the current tree into the new object, but does not pollute the original tree
// to prevent memory leaks and a growing tree
func (r *ServiceContainer) InjectInto(object interface{}) {
	var di = new(Graph)
	//di.Logger = sl{}

	for _, o := range r.unnamed {
		err := di.Provide(o)
		if err != nil {
			panic(err)
		}
	}

	for _, o := range r.named {
		err := di.Provide(o)
		if err != nil {
			panic(err)
		}
	}

	err := di.Provide(&Object{Value: object})
	if err != nil {
		panic(err)
	}

	err = di.Populate()
	if err != nil {
		panic(err)
	}

	for _, o := range di.Objects() {
		if pi, ok := o.Value.(PostInjecter); ok && !o.PostInjected {
			pi.PostInject(di)
			o.PostInjected = true
		}
	}
}

// GetByTag returns all registered objects with the given tag
func (r *ServiceContainer) GetByTag(tag string) []interface{} {
	return r.DI().GetByTag(tag)
}
