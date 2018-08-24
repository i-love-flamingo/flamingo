package application

import (
	"context"
	"encoding/gob"
	"net/http"
	"net/url"

	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
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

func init() {
	gob.Register(&oauth2.Token{})
	gob.Register(&oidc.IDToken{})
}

type (
	// authManager handles authentication related operations
	AuthManager struct {
		server              string
		secret              string
		clientID            string
		myHost              string
		disableOfflineToken bool
		logger              flamingo.Logger
		router              *router.Router

		openIDProvider *oidc.Provider
		oauth2Config   *oauth2.Config
	}
)

// Inject authManager dependencies
func (am *AuthManager) Inject(logger flamingo.Logger, router *router.Router, config *struct {
	Server              string `inject:"config:auth.server"`
	Secret              string `inject:"config:auth.secret"`
	ClientID            string `inject:"config:auth.clientid"`
	MyHost              string `inject:"config:auth.myhost"`
	DisableOfflineToken bool   `inject:"config:auth.disableOfflineToken"`
}) {
	am.logger = logger
	am.router = router
	am.server = config.Server
	am.secret = config.Secret
	am.clientID = config.ClientID
	am.myHost = config.MyHost
	am.disableOfflineToken = config.DisableOfflineToken
}

// Auth tries to retrieve the authentication context for a active session
func (am *AuthManager) Auth(c context.Context, session *sessions.Session) (domain.Auth, error) {
	ts, err := am.TokenSource(c, session)
	if err != nil {
		return domain.Auth{}, err
	}
	idToken, err := am.IDToken(c, session)
	if err != nil {
		return domain.Auth{}, err
	}

	return domain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}, nil
}

// OpenIDProvider is a lazy initialized OID provider
func (am *AuthManager) OpenIDProvider() *oidc.Provider {
	if am.openIDProvider == nil {
		var err error
		am.openIDProvider, err = oidc.NewProvider(context.Background(), am.server)
		if err != nil {
			panic(err)
		}
	}
	return am.openIDProvider
}

// OAuth2Config is lazy setup oauth2config
func (am *AuthManager) OAuth2Config() *oauth2.Config {
	if am.oauth2Config != nil {
		return am.oauth2Config
	}

	callbackURL := am.router.URL("auth.callback", nil)

	am.logger.WithField(flamingo.LogKeyCategory, "auth").Debug("am Callback", am, callbackURL)

	myhost, err := url.Parse(am.myHost)
	if err != nil {
		am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("Url parse failed", am.myHost, err)
	}
	callbackURL.Host = myhost.Host
	callbackURL.Scheme = myhost.Scheme
	scopes := []string{oidc.ScopeOpenID, "profile", "email"}
	if !am.disableOfflineToken {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
	}

	am.oauth2Config = &oauth2.Config{
		ClientID:     am.clientID,
		ClientSecret: am.secret,
		RedirectURL:  callbackURL.String(),

		// Discovery returns the OAuth2 endpoints.
		// It might panic here if Endpoint cannot be discovered
		Endpoint: am.OpenIDProvider().Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,
	}
	am.logger.WithField(flamingo.LogKeyCategory, "auth").Debug("am.oauth2Config", am.oauth2Config)
	return am.oauth2Config
}

// Verifier creates an OID verifier
func (am *AuthManager) Verifier() *oidc.IDTokenVerifier {
	return am.OpenIDProvider().Verifier(&oidc.Config{ClientID: am.clientID})
}

// OAuth2Token retrieves the oauth2 token from the session
func (am *AuthManager) OAuth2Token(session *sessions.Session) (*oauth2.Token, error) {
	if _, ok := session.Values[KeyToken]; !ok {
		return nil, errors.New("no token")
	}

	oauth2Token, ok := session.Values[KeyToken].(*oauth2.Token)
	if !ok {
		return nil, errors.Errorf("invalid token %T %v", session.Values[KeyToken], session.Values[KeyToken])
	}

	return oauth2Token, nil
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) IDToken(c context.Context, session *sessions.Session) (*oidc.IDToken, error) {
	token, _, err := am.getIDToken(c, session)
	return token, err
}

// GetRawIDToken gets the raw IDToken from session
func (am *AuthManager) GetRawIDToken(c context.Context, session *sessions.Session) (string, error) {
	_, raw, err := am.getIDToken(c, session)
	return raw, err
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) getIDToken(c context.Context, session *sessions.Session) (*oidc.IDToken, string, error) {
	if session == nil {
		return nil, "", errors.New("no session configured")
	}

	if token, ok := session.Values[KeyRawIDToken]; ok {
		idtoken, err := am.Verifier().Verify(c, token.(string))
		if err == nil {
			return idtoken, token.(string), nil
		}
	}

	token, raw, err := am.getNewIdToken(c, session)
	if err != nil {
		return nil, "", err
	}

	session.Values[KeyRawIDToken] = raw

	return token, raw, nil
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) getNewIdToken(c context.Context, session *sessions.Session) (*oidc.IDToken, string, error) {
	tokenSource, err := am.TokenSource(c, session)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	raw, err := am.ExtractRawIDToken(token)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	idtoken, err := am.Verifier().Verify(c, raw)

	if idtoken == nil {
		return nil, "", errors.New("idtoken nil")
	}

	return idtoken, raw, err
}

// ExtractRawIDToken from the provided (fresh) oatuh2token
func (am *AuthManager) ExtractRawIDToken(oauth2Token *oauth2.Token) (string, error) {
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", errors.Errorf("no id token %T %v", oauth2Token.Extra("id_token"), oauth2Token.Extra("id_token"))
	}

	return rawIDToken, nil
}

// TokenSource to be used in situations where you need it
func (am *AuthManager) TokenSource(c context.Context, session *sessions.Session) (oauth2.TokenSource, error) {
	oauth2Token, err := am.OAuth2Token(session)
	if err != nil {
		return nil, err
	}

	return am.OAuth2Config().TokenSource(c, oauth2Token), nil
}

// HTTPClient to retrieve a client with automatic tokensource
func (am *AuthManager) HTTPClient(c context.Context, session *sessions.Session) (*http.Client, error) {
	ts, err := am.TokenSource(c, session)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(c, ts), nil
}
