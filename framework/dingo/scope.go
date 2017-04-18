package dingo

import (
	"reflect"
	"sync"
)

type (
	// Scope defines a scope's behaviour
	Scope interface {
		resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value
	}

	// SingletonScope is our Scope to handle Singletons
	SingletonScope struct {
		sync.Mutex
		instanceLock map[reflect.Type]*sync.Mutex
		instances    map[reflect.Type]reflect.Value
	}
)

// Singleton is the default SingletonScope for dingo
var Singleton = new(SingletonScope)

// resolve a request in this scope
func (s *SingletonScope) resolveType(t reflect.Type, unscoped func(t reflect.Type, annotation string) reflect.Value) reflect.Value {
	// we got one :)
	if found, ok := s.instances[t]; ok {
		return found
	}

	// without an existing instance we need to create one

	// Lock ourselve
	s.Lock()

	// someone already built our instance while we were waiting for the lock
	if found, ok := s.instances[t]; ok {
		s.Unlock()
		return found
	}

	// If instanceLock is empty, create it now
	if s.instanceLock == nil {
		s.instanceLock = make(map[reflect.Type]*sync.Mutex)
	}

	// check for the concrete instanceLock
	if s.instanceLock[t] == nil {
		s.instanceLock[t] = new(sync.Mutex)
	}

	// acquire the instance-type's lock
	s.instanceLock[t].Lock()
	defer s.instanceLock[t].Unlock()

	// someone already built our instance while we were waiting/setup the locks
	if found, ok := s.instances[t]; ok {
		return found
	}

	if s.instances == nil {
		s.instances = make(map[reflect.Type]reflect.Value)
	}

	// release our main lock so we won't lock ourselves when trying to create a singleton
	// with a singleton dependency
	s.Unlock()

	// save our new generated singleton
	s.instances[t] = unscoped(t, "")

	// return the new singleton
	return s.instances[t]
}
