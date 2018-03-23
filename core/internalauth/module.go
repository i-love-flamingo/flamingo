package internalauth

import (
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/core/internalauth/domain"
	"go.aoe.com/flamingo/core/internalauth/application"
)

type (
	// InternalAuth
	InternalAuth struct {
		RouterRegistry *router.Registry `inject:""`
	}
)

// Configure the DI
func (m *InternalAuth) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.InternalAuthService)(nil)).To(application.OauthService{})
}
