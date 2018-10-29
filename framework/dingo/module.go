package dingo

import (
	"flamingo.me/dingo"
)

type (
	// Module is provided by packages to generate the DI tree
	// deprecated: use flamingo.me/dingo
	Module = dingo.Module

	// Depender defines a dependency-aware module
	// deprecated: use flamingo.me/dingo
	Depender = dingo.Depender
)

// TryModule tests if modules are properly bound
// deprecated: use flamingo.me/dingo
func TryModule(modules ...Module) (resultingError error) {
	return dingo.TryModule(modules...)
}
