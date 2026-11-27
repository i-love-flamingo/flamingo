// Package account is a test module with routes that show the login state.
// The app that loads it must also load auth.WebModule, oauth.Module and security.Module.
package account

import (
	"context"
	"errors"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/core/security/interface/middleware"
	"flamingo.me/flamingo/v3/framework/web"
)

// refreshBroker is the broker that the route "account.refresh" refreshes.
const refreshBroker = "shop"

var (
	errNoRefresher = errors.New("no broker " + refreshBroker + " that can refresh an identity")
	errNotOpenID   = errors.New("the identity is not an OpenID identity")
)

type (
	// Module of the account routes.
	Module struct{}

	routes struct {
		responder   *web.Responder
		security    application.SecurityService
		middleware  *middleware.SecurityMiddleware
		identity    *auth.WebIdentityService
		identifiers []auth.RequestIdentifier
	}
)

// Configure DI.
func (*Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

// Inject dependencies.
func (r *routes) Inject(
	responder *web.Responder,
	security application.SecurityService,
	securityMiddleware *middleware.SecurityMiddleware,
	identity *auth.WebIdentityService,
	identifiers []auth.RequestIdentifier,
) {
	r.responder = responder
	r.security = security
	r.middleware = securityMiddleware
	r.identity = identity
	r.identifiers = identifiers
}

// Routes of the module.
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/account/status", "account.status")
	registry.HandleGet("account.status", func(ctx context.Context, req *web.Request) web.Result {
		return r.responder.Data(map[string]bool{"loggedIn": r.security.IsLoggedIn(ctx, req.Session())})
	})

	registry.MustRoute("/account", "account.page")
	registry.HandleGet("account.page", r.middleware.HandleIfLoggedIn(func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]bool{"ok": true})
	}))

	registry.MustRoute("/account/login", "auth.login(redirecturl?)")
	registry.HandleGet("auth.login", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]string{})
	})

	registry.MustRoute("/account/claims", "account.claims")
	registry.HandleGet("account.claims", r.claims)

	registry.MustRoute("/account/refresh", "account.refresh")
	registry.HandleGet("account.refresh", r.refresh)
}

// claims gives the subject, the claims that the broker config maps to the keys "userID", "mail" and
// "tokenSubject", and whether the identity has an access token.
func (r *routes) claims(ctx context.Context, req *web.Request) web.Result {
	identity, err := r.identity.IdentifyAs(ctx, req, oauth.OpenIDTypeChecker)
	if err != nil {
		return r.responder.ForbiddenWithContext(ctx, err)
	}

	openID, ok := identity.(oauth.OpenIDIdentity)
	if !ok {
		return r.responder.ServerErrorWithContext(ctx, errNotOpenID)
	}

	var idToken struct {
		UserID string `json:"userID"`
		Mail   string `json:"mail"`
	}

	if err := openID.IDTokenClaims(&idToken); err != nil {
		return r.responder.ServerErrorWithContext(ctx, err)
	}

	var accessToken struct {
		TokenSubject string `json:"tokenSubject"`
	}

	if err := openID.AccessTokenClaims(&accessToken); err != nil {
		return r.responder.ServerErrorWithContext(ctx, err)
	}

	token, err := openID.TokenSource().Token()
	if err != nil {
		return r.responder.ServerErrorWithContext(ctx, err)
	}

	return r.responder.Data(map[string]any{
		"sub":            openID.Subject(),
		"userID":         idToken.UserID,
		"email":          idToken.Mail,
		"tokenSubject":   accessToken.TokenSubject,
		"hasAccessToken": token.AccessToken != "",
	})
}

// refresh makes the broker "shop" refresh the token at the next identification.
func (r *routes) refresh(ctx context.Context, req *web.Request) web.Result {
	for _, identifier := range r.identifiers {
		refresher, ok := identifier.(auth.WebIdentityRefresher)
		if !ok || identifier.Broker() != refreshBroker {
			continue
		}

		if err := refresher.RefreshIdentity(ctx, req); err != nil {
			return r.responder.ServerErrorWithContext(ctx, err)
		}

		return r.responder.Data(map[string]bool{"refreshed": true})
	}

	return r.responder.ServerErrorWithContext(ctx, errNoRefresher)
}
