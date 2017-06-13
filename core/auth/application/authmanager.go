package application

import (
	"context"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	// KeyToken defines where the authentication token is saved
	KeyToken = "auth.token"

	// KeyRawIDToken defines where the raw ID token is saved
	KeyRawIDToken = "auth.rawidtoken"

	// KeyAuthstate defines the current internal authentication state
	KeyAuthstate = "auth.state"
)

type (
	// AuthManager handles authentication related operations
	AuthManager struct {
		Server    string `inject:"config:auth.server"`
		Secret    string `inject:"config:auth.secret"`
		ClientID  string `inject:"config:auth.clientid"`
		MyHost    string `inject:"config:auth.myhost"`
		NotBefore time.Time

		Router *router.Router `inject:""`

		openIDProvider *oidc.Provider
		oauth2Config   *oauth2.Config
	}
)

// OpenIDProvider is a lazy initialized OID provider
func (authmanager *AuthManager) OpenIDProvider() *oidc.Provider {
	if authmanager.openIDProvider == nil {
		var err error
		authmanager.openIDProvider, err = oidc.NewProvider(context.Background(), authmanager.Server)
		if err != nil {
			panic(err)
		}
	}
	return authmanager.openIDProvider
}

// OAuth2Config is lazy setup oauth2config
func (authmanager *AuthManager) OAuth2Config() *oauth2.Config {
	if authmanager.oauth2Config != nil {
		return authmanager.oauth2Config
	}

	callbackURL := authmanager.Router.URL("auth.callback", nil)
	myhost, _ := url.Parse(authmanager.MyHost)
	callbackURL.Host = myhost.Host
	callbackURL.Scheme = myhost.Scheme
	authmanager.oauth2Config = &oauth2.Config{
		ClientID:     authmanager.ClientID,
		ClientSecret: authmanager.Secret,
		RedirectURL:  callbackURL.String(),

		// Discovery returns the OAuth2 endpoints.
		Endpoint: authmanager.OpenIDProvider().Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return authmanager.oauth2Config
}

// Verifier creates an OID verifier
func (authmanager *AuthManager) Verifier() *oidc.IDTokenVerifier {
	return authmanager.OpenIDProvider().Verifier(&oidc.Config{ClientID: authmanager.ClientID})
}

// OAuth2Token retrieves the oauth2 token from the session
func (authmanager *AuthManager) OAuth2Token(c web.Context) (*oauth2.Token, error) {
	if _, ok := c.Session().Values[KeyToken]; !ok {
		return nil, errors.New("no token")
	}

	oauth2Token, ok := c.Session().Values[KeyToken].(*oauth2.Token)
	if !ok {
		return nil, errors.Errorf("invalid token %T %v", c.Session().Values[KeyToken], c.Session().Values[KeyToken])
	}

	return oauth2Token, nil
}

// IDToken retrieves and validates the ID Token from the session
func (authmanager *AuthManager) IDToken(c web.Context) (*oidc.IDToken, error) {
	if _, ok := c.Session().Values[KeyRawIDToken]; !ok {
		return nil, errors.New("no id token")
	}

	rawIDToken, ok := c.Session().Values[KeyRawIDToken].(string)
	if !ok {
		return nil, errors.Errorf("invalid id token %T %v", c.Session().Values[KeyRawIDToken], c.Session().Values[KeyRawIDToken])
	}

	// Parse and verify ID Token payload.
	idToken, err := authmanager.Verifier().Verify(c, rawIDToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if !authmanager.NotBefore.IsZero() && idToken.IssuedAt.Before(authmanager.NotBefore) {
		return nil, errors.New("issued before allowed")
	}

	return idToken, nil
}

// ExtractRawIDToken from the provided (fresh) oatuh2token
func (authmanager *AuthManager) ExtractRawIDToken(oauth2Token *oauth2.Token) (string, error) {
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", errors.Errorf("no id token %T %v", oauth2Token.Extra("id_token"), oauth2Token.Extra("id_token"))
	}

	return rawIDToken, nil
}

// TokenSource to be used in situations where you need it
func (authmanager *AuthManager) TokenSource(c web.Context) (oauth2.TokenSource, error) {
	oauth2Token, err := authmanager.OAuth2Token(c)
	if err != nil {
		return nil, err
	}

	return authmanager.OAuth2Config().TokenSource(c, oauth2Token), nil
}

// HTTPClient to retrieve a client with automatic tokensource
func (authmanager *AuthManager) HTTPClient(c web.Context) (*http.Client, error) {
	ts, err := authmanager.TokenSource(c)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(c, ts), nil
}
