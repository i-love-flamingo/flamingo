package internalauth

import (
	"flamingo.me/flamingo/v3/core/internalauth/application"
	"flamingo.me/flamingo/v3/core/internalauth/domain"
	"flamingo.me/flamingo/v3/framework/dingo"
)

// InternalAuth
type InternalAuth struct{}

// Configure the DI
func (m *InternalAuth) Configure(injector *dingo.Injector) {
	injector.Bind((*domain.InternalAuthService)(nil)).To(application.OauthService{})
}
