package application

import (
	"context"
	"encoding/gob"
	"net/http"
	"net/url"

	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
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

	// KeyToken defines where the authentication token extras are saved
	KeyTokenExtras = "auth.token.extras"
)

func init() {
	gob.Register(&oauth2.Token{})
	gob.Register(&oidc.IDToken{})
	gob.Register(&domain.TokenExtras{})
}

type (
	// authManager handles authentication related operations
	AuthManager struct {
		server              string
		secret              string
		clientID            string
		myHost              string
		allowHostFromReq    bool
		disableOfflineToken bool
		scopes              config.Slice
		idTokenMapping      config.Slice
		userInfoMapping     config.Slice
		logger              flamingo.Logger
		router              *router.Router

		openIDProvider *oidc.Provider
		oauth2Config   map[string]*oauth2.Config
	}
)

// Inject authManager dependencies
func (am *AuthManager) Inject(logger flamingo.Logger, router *router.Router, config *struct {
	Server              string       `inject:"config:auth.server"`
	Secret              string       `inject:"config:auth.secret"`
	ClientID            string       `inject:"config:auth.clientid"`
	MyHost              string       `inject:"config:auth.myhost"`
	AllowHostFromReq    bool         `inject:"config:auth.allowHostFromReq,optional"`
	DisableOfflineToken bool         `inject:"config:auth.disableOfflineToken"`
	Scopes              config.Slice `inject:"config:auth.scopes"`
	IdTokenMapping      config.Slice `inject:"config:auth.claims.idToken"`
	UserInfoMapping     config.Slice `inject:"config:auth.claims.userInfo"`
}) {
	am.logger = logger
	am.router = router
	am.server = config.Server
	am.secret = config.Secret
	am.clientID = config.ClientID
	am.myHost = config.MyHost
	am.allowHostFromReq = config.AllowHostFromReq
	am.disableOfflineToken = config.DisableOfflineToken
	am.scopes = config.Scopes
	am.idTokenMapping = config.IdTokenMapping
	am.userInfoMapping = config.UserInfoMapping
	am.oauth2Config = make(map[string]*oauth2.Config)
}

func (am *AuthManager) URL(ctx context.Context, path string) (*url.URL, error) {
	ubase := *am.router.Base()
	u := &ubase
	if path != "" {
		u.Path = path
	}

	myhost, err := url.Parse(am.myHost)
	if err != nil {
		return nil, err
	}

	u.Host = myhost.Host
	if r, ok := web.FromContext(ctx); ok && am.allowHostFromReq {
		u.Host = r.Request().Host
	}
	u.Scheme = myhost.Scheme

	return u, nil
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
func (am *AuthManager) OAuth2Config(ctx context.Context) *oauth2.Config {
	callbackURL, err := am.URL(ctx, am.router.URL("auth.callback", nil).Path)
	if err != nil {
		am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("could not get url", err)
	}

	if cfg, ok := am.oauth2Config[callbackURL.String()]; ok {
		return cfg
	}

	am.logger.WithField(flamingo.LogKeyCategory, "auth").Debug("am Callback", am, callbackURL)

	var scopes []string
	err = am.scopes.MapInto(&scopes)
	if err != nil {
		am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("could not parse scopes from config", am.scopes, err)
	}

	scopes = append([]string{oidc.ScopeOpenID}, scopes...)
	if !am.disableOfflineToken {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
	}

	am.oauth2Config[callbackURL.String()] = &oauth2.Config{
		ClientID:     am.clientID,
		ClientSecret: am.secret,
		RedirectURL:  callbackURL.String(),

		// Discovery returns the OAuth2 endpoints.
		// It might panic here if Endpoint cannot be discovered
		Endpoint: am.OpenIDProvider().Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,

		ClaimSet: am.getClaimsRequestParameter(),
	}

	am.logger.WithField(flamingo.LogKeyCategory, "auth").Debug("am.oauth2Config", am.oauth2Config)
	return am.oauth2Config[callbackURL.String()]
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

func (am *AuthManager) getClaimsRequestParameter() *oauth2.ClaimSet {
	var claimSet *oauth2.ClaimSet

	claimSet = am.createClaimSetFromMapping(oauth2.IdTokenClaim, am.idTokenMapping, claimSet)
	claimSet = am.createClaimSetFromMapping(oauth2.UserInfoClaim, am.userInfoMapping, claimSet)

	return claimSet
}

func (am *AuthManager) createClaimSetFromMapping(topLevelName string, configuration config.Slice, claimSet *oauth2.ClaimSet) *oauth2.ClaimSet {
	var mapping []string
	configuration.MapInto(&mapping)

	for _, name := range mapping {
		if name == "" {
			continue
		}
		if claimSet == nil {
			claimSet = &oauth2.ClaimSet{}
		}
		claimSet.AddVoluntaryClaim(topLevelName, name)
	}

	return claimSet
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

	return am.OAuth2Config(c).TokenSource(c, oauth2Token), nil
}

// HTTPClient to retrieve a client with automatic tokensource
func (am *AuthManager) HTTPClient(c context.Context, session *sessions.Session) (*http.Client, error) {
	ts, err := am.TokenSource(c, session)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(c, ts), nil
}
