package dingo

import (
	"errors"
	"fmt"
)

type (
	// Module is provided by packages to generate the DI tree
	Module interface {
		Configure(injector *Injector)
	}

	// Depender defines a dependency-aware module
	Depender interface {
		Depends() []Module
	}
)

// TryModule tests if modules are properly bound
func TryModule(modules ...Module) (resultingError error) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				resultingError = err
				return
			}
			resultingError = errors.New(fmt.Sprint(err))
		}
	}()

	injector := NewInjector()
	injector.buildEagerSingletons = false
	injector.InitModules(modules...)
	return nil
}
