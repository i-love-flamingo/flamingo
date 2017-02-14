package app

import (
	"log"

	"github.com/facebookgo/inject"
)

type (
	Registrator struct {
		objects []interface{}

		handlers map[string]interface{}

		routes map[string]string
	}

	RegisterFunc func(r *Registrator)
)

func NewRegistrator() *Registrator {
	return &Registrator{
		handlers: make(map[string]interface{}),
		routes:   make(map[string]string),
	}
}

func (r *Registrator) Register(rfs ...RegisterFunc) {
	for _, rf := range rfs {
		rf(r)
	}
}

func (r *Registrator) Route(path, name string) {
	r.routes[path] = name
}

func (r *Registrator) Handle(name string, handler interface{}) {
	r.handlers[name] = handler
}

func (r *Registrator) Object(i ...interface{}) {
	r.objects = append(r.objects, i...)
}

type sl struct{}

func (_ sl) Debugf(a string, b ...interface{}) {
	log.Printf(a, b...)
}

func (r *Registrator) Resolve() {
	var di inject.Graph

	di.Logger = sl{}

	for _, o := range r.objects {
		di.Provide(&inject.Object{Value: o})
	}
	for _, h := range r.handlers {
		di.Provide(&inject.Object{Value: h})
	}

	di.Populate()
}
