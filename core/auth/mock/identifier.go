package mock

import (
	"context"
	"errors"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	IdentifyMethod     func(*Identifier, context.Context, *web.Request) (auth.Identity, error)
	AuthenticateMethod func(*Identifier, context.Context, *web.Request) web.Result
	CallbackMethod     func(*Identifier, context.Context, *web.Request, func(*web.Request) *url.URL) web.Result
	LogoutMethod       func(*Identifier, context.Context, *web.Request)

	// Identifier mocks request identification
	Identifier struct {
		broker             string
		authenticateMethod AuthenticateMethod
		callbackMethod     CallbackMethod
		logoutMethod       LogoutMethod
		identifyMethod     IdentifyMethod
	}
)

var (
	_ auth.RequestIdentifier = new(Identifier)
	_ auth.Identity          = new(Identity)
	_ oauth.OpenIDIdentity   = new(OIDCIdentity)
)

// Broker getter
func (m *Identifier) Broker() string {
	return m.broker
}

// SetBroker identity for the identifier
func (m *Identifier) SetBroker(broker string) {
	m.broker = broker
}

// Identify an incoming request
func (m *Identifier) Identify(ctx context.Context, r *web.Request) (auth.Identity, error) {
	if m.identifyMethod != nil {
		return m.identifyMethod(m, ctx, r)
	}
	return nil, errors.New("can not identify")
}

// SetIdentifyMethod for mock identifier
func (m *Identifier) SetIdentifyMethod(method IdentifyMethod) *Identifier {
	m.identifyMethod = method
	return m
}

// Identify an incoming request
func (m *Identifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	if m.authenticateMethod != nil {
		return m.authenticateMethod(m, ctx, request)
	}
	return nil
}

// Identify an incoming request
func (m *Identifier) SetAuthenticateMethod(method AuthenticateMethod) *Identifier {
	m.authenticateMethod = method
	return m
}

// Identify an incoming request
func (m *Identifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	if m.callbackMethod != nil {
		return m.callbackMethod(m, ctx, request, returnTo)
	}
	return nil
}

// Identify an incoming request
func (m *Identifier) SetCallbackMethod(method CallbackMethod) *Identifier {
	m.callbackMethod = method
	return m
}

// Identify an incoming request
func (m *Identifier) Logout(ctx context.Context, request *web.Request) {
	if m.logoutMethod != nil {
		m.logoutMethod(m, ctx, request)
	}
}

// Identify an incoming request
func (m *Identifier) SetLogoutMethod(method LogoutMethod) *Identifier {
	m.logoutMethod = method
	return m
}
