package custom

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

// OidcModule configures an "overridden" openid connect module
type OidcModule struct{}

// Configure dependency injection
func (*OidcModule) Configure(injector *dingo.Injector) {
	injector.BindMap(new(auth.RequestIdentifierFactory), "customOidcBroker").ToProvider(func(identifierConfig *struct {
		Oidc auth.RequestIdentifierFactory `inject:"map:oidc"`
	}) auth.RequestIdentifierFactory {
		return func(cfg config.Map) (auth.RequestIdentifier, error) {
			identifier, err := identifierConfig.Oidc(cfg["oidc"].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			if err := injector.RequestInjection(identifier); err != nil {
				return nil, err
			}
			return &oidcBroker{oidcBroker: identifier}, nil
		}
	})
}

// CueConfig for oidcBroker
func (*OidcModule) CueConfig() string {
	return `
customOidcBroker :: {
	typ: "customOidcBroker"
	broker: "customOidcBroker"
	oidc: {
		core.auth.oidc
		broker: "customOidcBroker"
		clientID: "customOidcBroker"
		clientSecret: "customOidcBroker"
		"endpoint": "http://127.0.0.1:3351/dex"
	}
}
`
}

type oidcBroker struct {
	oidcBroker auth.RequestIdentifier
	responder  *web.Responder
}

// Inject dependencies
func (b *oidcBroker) Inject(responder *web.Responder) {
	b.responder = responder
}

// Broker getter
func (*oidcBroker) Broker() string {
	return "customOidcBroker"
}

// Identify request
func (b *oidcBroker) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	if val, ok := request.Session().Load("customOidcBroker.loggedIn"); !ok || !val.(bool) {
		return nil, errors.New("no customOidcBroker loggedIn flag")
	}
	return b.oidcBroker.Identify(ctx, request)
}

// Authenticate request
func (b *oidcBroker) Authenticate(ctx context.Context, request *web.Request) web.Result {
	return b.oidcBroker.(auth.WebAuthenticater).Authenticate(ctx, request)
}

// Callback handler for request
func (b *oidcBroker) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	if identity, err := b.oidcBroker.Identify(ctx, request); err != nil && identity == nil {
		// okURL marks a final callback, e.g. everything went fine
		okURL := new(url.URL)
		if result, ok := b.oidcBroker.(auth.WebCallbacker).Callback(ctx, request, func(request *web.Request) *url.URL {
			return okURL
		}).(*web.URLRedirectResponse); ok && result.URL != okURL {
			// upstream callback failed
			return result
		}
	}

	if q, err := request.Query1("oidcbroker"); err == nil && q != "" {
		request.Session().Store("customOidcBroker.loggedIn", true)
	}

	if val, ok := request.Session().Load("customOidcBroker.loggedIn"); !ok || !val.(bool) {
		return b.responder.HTTP(http.StatusOK, strings.NewReader(`<a href="?oidcbroker=1">Confirm Login</a>`))
	}

	// we have finished
	return b.responder.URLRedirect(returnTo(request))
}

// Logout handler for request
func (b *oidcBroker) Logout(ctx context.Context, request *web.Request) *url.URL {
	return b.oidcBroker.(auth.WebLogoutWithRedirect).Logout(ctx, request)
}
