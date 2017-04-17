package dingo

import (
	"log"
	"reflect"
	"sync"
)

type (
	Scope interface {
		resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value
	}

	SingletonScope struct {
		sync.Mutex
		instanceLock map[reflect.Type]*sync.Mutex
		instances    map[reflect.Type]reflect.Value
	}
)

var Singleton = new(SingletonScope)

func (s *SingletonScope) resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value {
	if found, ok := s.instances[t]; ok {
		return found
	}

	s.Lock()

	// someone already built our instance while we were waiting for the lock...
	if found, ok := s.instances[t]; ok {
		s.Unlock()
		return found
	}

	if s.instanceLock == nil {
		s.instanceLock = make(map[reflect.Type]*sync.Mutex)
	}

	if s.instanceLock[t] == nil {
		s.instanceLock[t] = new(sync.Mutex)
	}

	s.instanceLock[t].Lock()
	defer s.instanceLock[t].Unlock()

	// someone already built our instance while we were waiting for the lock...
	if found, ok := s.instances[t]; ok {
		return found
	}

	s.Unlock()

	if s.instances == nil {
		s.instances = make(map[reflect.Type]reflect.Value)
	}

	log.Println("singleton creates unscoped ", t)
	s.instances[t] = unscoped(t, "")
	return s.instances[t]
}
