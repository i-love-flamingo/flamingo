package internalauth

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/internalauth/application"
	"flamingo.me/flamingo/v3/core/internalauth/domain"
)

// InternalAuth module for backend oauth usage
type InternalAuth struct{}

// Configure the DI
func (m *InternalAuth) Configure(injector *dingo.Injector) {
	injector.Bind(new(domain.InternalAuthService)).To(application.OauthService{})
}
