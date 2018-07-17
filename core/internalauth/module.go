package internalauth

import (
	"flamingo.me/flamingo/core/internalauth/application"
	"flamingo.me/flamingo/core/internalauth/domain"
	"flamingo.me/flamingo/framework/dingo"
)

// InternalAuth
type InternalAuth struct{}

// Configure the DI
func (m *InternalAuth) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.InternalAuthService)(nil)).To(application.OauthService{})
}
