package internalauth

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/core/internalauth/domain"
	"flamingo.me/flamingo/core/internalauth/application"
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
