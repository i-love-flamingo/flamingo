package dingo

import (
	"reflect"
	"sync"
)

type (
	Scope interface {
		resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value
	}

	baseScope struct {
		sync.Mutex
		instances map[reflect.Type]reflect.Value
	}
)

var Singleton = new(baseScope)

func (s *baseScope) resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value {
	if found, ok := s.instances[t]; ok {
		return found
	}

	s.Lock()
	defer s.Unlock()

	// someone already built our instance while we were waiting for the lock...
	if found, ok := s.instances[t]; ok {
		return found
	}

	if s.instances == nil {
		s.instances = make(map[reflect.Type]reflect.Value)
	}

	s.instances[t] = unscoped(t, "")
	return s.instances[t]
}
