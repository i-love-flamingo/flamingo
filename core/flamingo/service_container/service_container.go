package service_container

import (
	"log"
	"reflect"
	"runtime"

	"github.com/facebookgo/inject"
)

type (
	// ServiceContainer is a basic flamingo helper
	// to register default Routes, packages, etc.
	ServiceContainer struct {
		unnamed []*inject.Object
		named   map[string]*inject.Object
		tags    map[string][]*inject.Object
		Routes  map[string]string
		Handler map[string]interface{}
		di      *inject.Graph
	}

	// RegisterFunc defines a callback used by packages to bootstrap themselves
	RegisterFunc func(r *ServiceContainer)

	// PostInjecter defines the PostInject() function which is called when the DI has finished
	PostInjecter interface {
		PostInject()
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
		named:   make(map[string]*inject.Object),
		tags:    make(map[string][]*inject.Object),
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
	r.Register(handler)
}

// Route adds a route
func (r *ServiceContainer) Route(path, name string) *ServiceContainer {
	r.Routes[path] = name
	return r
}

// Register registers any object for DI
func (r *ServiceContainer) Register(o interface{}, tags ...string) *ServiceContainer {
	r.Remove(o)

	object := &inject.Object{
		Value: o,
	}
	r.unnamed = append(r.unnamed, object)
	for _, tag := range tags {
		r.tags[tag] = append(r.tags[tag], object)
	}
	return r
}

// RegisterNamed registers any object for DI with a given name
func (r *ServiceContainer) RegisterNamed(name string, o interface{}, tags ...string) *ServiceContainer {
	object := &inject.Object{
		Value: o,
		Name:  name,
	}
	r.named[name] = object
	for _, tag := range tags {
		r.tags[tag] = append(r.tags[tag], object)
	}
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
func (_ sl) Debugf(a string, b ...interface{}) {
	log.Printf(a, b...)
}

// DI returns the injection graph, not populated
func (r *ServiceContainer) DI() *inject.Graph {
	if r.di != nil {
		return r.di
	}

	r.di = new(inject.Graph)

	r.Register(r)

	for _, o := range r.unnamed {
		r.di.Provide(o)
	}

	for _, o := range r.named {
		r.di.Provide(o)
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
		if o, ok := o.Value.(PostInjecter); ok {
			o.PostInject()
		}
	}
}

// Create a new object of type object and inject into it
func (r *ServiceContainer) Create(object interface{}) interface{} {
	var t = reflect.TypeOf(object)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var o = reflect.New(t).Interface()
	r.InjectInto(o)
	return o
}

// InjectInto injects resolves the current tree into the new object, but does not pollute the original tree
// to prevent memory leaks and a growing tree
func (r *ServiceContainer) InjectInto(object interface{}) {
	var di inject.Graph
	//di.Logger = sl{}
	di.Provide(r.DI().Objects()...)
	di.Provide(&inject.Object{Value: object})
	err := di.Populate()
	if err != nil {
		panic(err)
	}
	if object, ok := object.(PostInjecter); ok {
		object.PostInject()
	}
}

// GetByTag returns all registered objects with the given tag
func (r *ServiceContainer) GetByTag(tag string) (res []interface{}) {
	for _, o := range r.tags[tag] {
		res = append(res, o.Value)
	}
	return
}
