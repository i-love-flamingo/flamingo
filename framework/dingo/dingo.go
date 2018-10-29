package dingo

import (
	"flamingo.me/dingo"
)

const (
	// INIT state
	INIT = dingo.INIT
	// DEFAULT state
	DEFAULT = dingo.INIT
)

// EnableCircularTracing activates dingo's trace feature to find circular dependencies
// this is super expensive (memory wise), so it should only be used for debugging purposes
// deprecated: use flamingo.me/dingo
func EnableCircularTracing() {
	dingo.EnableCircularTracing()
}

type (
	// Injector defines bindings and multibindings
	// it is possible to have a parent-injector, which can be asked if no resolution is available
	// deprecated: use flamingo.me/dingo
	Injector = dingo.Injector
)

// NewInjector builds up a new Injector out of a list of Modules
// deprecated: use flamingo.me/dingo
func NewInjector(modules ...Module) *Injector {
	return dingo.NewInjector(modules...)
}
