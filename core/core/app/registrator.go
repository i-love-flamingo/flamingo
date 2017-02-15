package app

import (
	"log"

	"github.com/facebookgo/inject"
)

type (
	// Registrator is a basic flamingo helper
	// to register default routes, packages, etc.
	Registrator struct {
		objects []interface{}

		handlers map[string]interface{}

		routes map[string]string
	}

	// RegisterFunc defines a callback used by packages to bootstrap themselves
	RegisterFunc func(r *Registrator)
)

// NewRegistrator creates a new Registrator
func NewRegistrator() *Registrator {
	return &Registrator{
		handlers: make(map[string]interface{}),
		routes:   make(map[string]string),
	}
}

// Register calls the provided RegisterFunc callbacks
func (r *Registrator) Register(rfs ...RegisterFunc) {
	for _, rf := range rfs {
		rf(r)
	}
}

// Route adds a route
func (r *Registrator) Route(path, name string) {
	r.routes[path] = name
}

// Handle adds a handler
func (r *Registrator) Handle(name string, handler interface{}) {
	r.handlers[name] = handler
}

// Object registers any object for DI
func (r *Registrator) Object(i ...interface{}) {
	r.objects = append(r.objects, i...)
}

// sl is a private logger to show DI logs
type sl struct{}

// Debugf DI logger
func (_ sl) Debugf(a string, b ...interface{}) {
	log.Printf(a, b...)
}

// DI returns the injection graph, not populated
func (r *Registrator) DI() inject.Graph {
	var di inject.Graph

	di.Logger = sl{}

	for _, o := range r.objects {
		di.Provide(&inject.Object{Value: o})
	}
	for _, h := range r.handlers {
		di.Provide(&inject.Object{Value: h})
	}

	return di
}

// Resolve populates the injection graph
func (r *Registrator) Resolve() {
	r.DI().Populate()
}
