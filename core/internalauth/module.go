package internalauth

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type (
	// InternalAuth
	InternalAuth struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the DI
func (m *InternalAuth) Configure(injector *dingo.Injector) {
}
