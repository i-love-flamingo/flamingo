package interfaces

import (
	"context"
	"errors"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// LegacyIdentity is an oauth.OIDCIdentifier for old oauth module
type LegacyIdentity struct {
	auth       domain.Auth
	rawIDToken string
}

var _ oauth.OpenIDIdentity = new(LegacyIdentity)

// Subject for this identity
func (identity *LegacyIdentity) Subject() string {
	return identity.auth.IDToken.Subject
}

// Broker code, hardcoded
func (identity *LegacyIdentity) Broker() string {
	return "flamingo.core.oauth"
}

// TokenSource returns the oauth2 token source
func (identity *LegacyIdentity) TokenSource() oauth2.TokenSource {
	return identity.auth.TokenSource
}

// AccessTokenClaims is not supported with the old module
func (identity *LegacyIdentity) AccessTokenClaims(into interface{}) error {
	return errors.New("flamingo.core.oauth does not support AccessTokenClaims")
}

// IDToken getter
func (identity *LegacyIdentity) IDToken() *oidc.IDToken {
	return identity.auth.IDToken
}

// IDTokenClaims mapper
func (identity *LegacyIdentity) IDTokenClaims(into interface{}) error {
	return identity.auth.IDToken.Claims(into)
}

// RawIDToken returns the raw JWT id token
func (identity *LegacyIdentity) RawIDToken() string {
	return identity.rawIDToken
}

// LegacyIdentifier bridges core/oauth and core/auth/oauth together
type LegacyIdentifier struct {
	authmanager        *application.AuthManager
	responder          *web.Responder
	callbackController CallbackControllerInterface
	loginController    LoginControllerInterface
	logoutController   LogoutControllerInterface
}

// Inject dependencies
func (identifier *LegacyIdentifier) Inject(
	authmanager *application.AuthManager,
	responder *web.Responder,
	callbackController CallbackControllerInterface,
	loginController LoginControllerInterface,
	logoutController LogoutControllerInterface,
) *LegacyIdentifier {
	identifier.authmanager = authmanager
	identifier.responder = responder
	identifier.loginController = loginController
	identifier.callbackController = callbackController
	identifier.logoutController = logoutController
	return identifier
}

// Broker hardcoded to flamingo.core.oauth
func (*LegacyIdentifier) Broker() string {
	return "flamingo.core.oauth"
}

// Identify an incoming request with the authmanager
func (identifier *LegacyIdentifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	authData, err := identifier.authmanager.Auth(ctx, request.Session())
	if err != nil {
		return nil, err
	}

	rawIDToken, _ := identifier.authmanager.GetRawIDToken(ctx, request.Session())
	return &LegacyIdentity{auth: authData, rawIDToken: rawIDToken}, nil
}

// Authenticate an incoming request with the logincontroller
func (identifier *LegacyIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	return identifier.loginController.Get(ctx, request)
}

// Callback for the current request
func (identifier *LegacyIdentifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	request.Session().Store("auth.redirect", returnTo(request).String())
	return identifier.callbackController.Get(ctx, request)
}

// Logout using the legacy logout controller
func (identifier *LegacyIdentifier) Logout(ctx context.Context, request *web.Request) *url.URL {
	resp := identifier.logoutController.Get(ctx, request)
	if ur, ok := resp.(*web.URLRedirectResponse); ok {
		return ur.URL
	}
	return nil
}
