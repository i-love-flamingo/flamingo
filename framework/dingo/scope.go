package dingo

import (
	"flamingo.me/dingo"
)

type (
	// Scope defines a scope's behaviour
	// deprecated: use flamingo.me/dingo
	Scope = dingo.Scope

	// SingletonScope is our Scope to handle Singletons
	// deprecated: use flamingo.me/dingo
	SingletonScope = dingo.SingletonScope

	// ChildSingletonScope manages child-specific singleton
	// deprecated: use flamingo.me/dingo
	ChildSingletonScope = dingo.ChildSingletonScope
)

var (
	// Singleton is the default SingletonScope for dingo
	// deprecated: use flamingo.me/dingo
	Singleton Scope = NewSingletonScope()

	// ChildSingleton is a per-child singleton
	// deprecated: use flamingo.me/dingo
	ChildSingleton Scope = NewChildSingletonScope()
)

// NewSingletonScope creates a new singleton scope
// deprecated: use flamingo.me/dingo
func NewSingletonScope() *SingletonScope {
	return dingo.NewSingletonScope()
}

// NewChildSingletonScope creates a new child singleton scope
// deprecated: use flamingo.me/dingo
func NewChildSingletonScope() *ChildSingletonScope {
	return dingo.NewChildSingletonScope()
}
