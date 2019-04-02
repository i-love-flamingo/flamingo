package application

import (
	"context"
	"encoding/gob"
	"net/http"

	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	// keyToken defines where the authentication token is saved
	keyToken = "auth.token"

	// keyRawIDToken defines where the raw ID token is saved
	keyRawIDToken = "auth.rawidtoken"

	// keyAuthstate defines the current internal authentication state
	keyAuthstate = "auth.state"

	// keyToken defines where the authentication token extras are saved
	keyTokenExtras = "auth.token.extras"
)

func init() {
	gob.Register(&oauth2.Token{})
	gob.Register(&oidc.IDToken{})
	gob.Register(&domain.TokenExtras{})
}

type (
	// AuthManager handles authentication related operations
	AuthManager struct {
		server              string
		secret              string
		clientID            string
		disableOfflineToken bool
		scopes              config.Slice
		idTokenMapping      config.Slice
		userInfoMapping     config.Slice
		logger              flamingo.Logger
		router              *web.Router
		openIDProvider      *oidc.Provider
	}
)

// Inject authManager dependencies
func (am *AuthManager) Inject(logger flamingo.Logger, router *web.Router, openIDProvider *oidc.Provider, config *struct {
	Server              string       `inject:"config:oauth.server"`
	Secret              string       `inject:"config:oauth.secret"`
	ClientID            string       `inject:"config:oauth.clientid"`
	DisableOfflineToken bool         `inject:"config:oauth.disableOfflineToken"`
	Scopes              config.Slice `inject:"config:oauth.scopes"`
	IDTokenMapping      config.Slice `inject:"config:oauth.claims.idToken"`
	UserInfoMapping     config.Slice `inject:"config:oauth.claims.userInfo"`
}) {
	am.logger = logger
	am.router = router
	am.server = config.Server
	am.secret = config.Secret
	am.clientID = config.ClientID
	am.disableOfflineToken = config.DisableOfflineToken
	am.scopes = config.Scopes
	am.idTokenMapping = config.IDTokenMapping
	am.userInfoMapping = config.UserInfoMapping

	var err error
	am.openIDProvider, err = oidc.NewProvider(context.Background(), config.Server)
	if err != nil {
		panic(err)
	}
}

// Auth tries to retrieve the authentication context for a active session
func (am *AuthManager) Auth(c context.Context, session *web.Session) (domain.Auth, error) {
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
	return am.openIDProvider
}

// OAuth2Config is lazy setup oauth2config
func (am *AuthManager) OAuth2Config(ctx context.Context, req *web.Request) *oauth2.Config {
	var redirectUrl string
	if req != nil {
		callbackURL, _ := am.router.Absolute(req, "auth.callback", nil)
		redirectUrl = callbackURL.String()
	}

	var scopes []string
	err := am.scopes.MapInto(&scopes)
	if err != nil {
		am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("could not parse scopes from config", am.scopes, err)
	}

	scopes = append([]string{oidc.ScopeOpenID}, scopes...)
	if !am.disableOfflineToken {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     am.clientID,
		ClientSecret: am.secret,
		RedirectURL:  redirectUrl,

		Endpoint: am.OpenIDProvider().Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,

		ClaimSet: am.getClaimsRequestParameter(),
	}

	am.logger.WithField(flamingo.LogKeyCategory, "auth").Debug("am.oauth2Config", oauth2Config)
	return oauth2Config
}

// Verifier creates an OID verifier
func (am *AuthManager) Verifier() *oidc.IDTokenVerifier {
	return am.OpenIDProvider().Verifier(&oidc.Config{ClientID: am.clientID})
}

// OAuth2Token retrieves the oauth2 token from the session
func (am *AuthManager) OAuth2Token(session *web.Session) (*oauth2.Token, error) {
	if _, ok := session.Load(keyToken); !ok {
		return nil, errors.New("no token")
	}

	value, _ := session.Load(keyToken)
	oauth2Token, ok := value.(*oauth2.Token)
	if !ok {
		return nil, errors.Errorf("invalid token %#v", value)
	}

	return oauth2Token, nil
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) IDToken(c context.Context, session *web.Session) (*oidc.IDToken, error) {
	token, _, err := am.getIDToken(c, session)
	return token, err
}

// GetRawIDToken gets the raw IDToken from session
func (am *AuthManager) GetRawIDToken(c context.Context, session *web.Session) (string, error) {
	_, raw, err := am.getIDToken(c, session)
	return raw, err
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) getIDToken(c context.Context, session *web.Session) (*oidc.IDToken, string, error) {
	if session == nil {
		return nil, "", errors.New("no session configured")
	}

	if token, ok := session.Load(keyRawIDToken); ok {
		idtoken, err := am.Verifier().Verify(c, token.(string))
		if err == nil {
			return idtoken, token.(string), nil
		}
	}

	token, raw, err := am.getNewIDToken(c, session)
	if err != nil {
		return nil, "", err
	}

	session.Store(keyRawIDToken, raw)

	return token, raw, nil
}

// IDToken retrieves and validates the ID Token from the session
func (am *AuthManager) getNewIDToken(c context.Context, session *web.Session) (*oidc.IDToken, string, error) {
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
	err := configuration.MapInto(&mapping)
	if err != nil {
		am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("could not map configuration", err)
	}

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
func (am *AuthManager) TokenSource(c context.Context, session *web.Session) (oauth2.TokenSource, error) {
	oauth2Token, err := am.OAuth2Token(session)
	if err != nil {
		return nil, err
	}

	return am.OAuth2Config(c, nil).TokenSource(c, oauth2Token), nil
}

// HTTPClient to retrieve a client with automatic tokensource
func (am *AuthManager) HTTPClient(c context.Context, session *web.Session) (*http.Client, error) {
	ts, err := am.TokenSource(c, session)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(c, ts), nil
}

// StoreTokenDetails stores all token related data into session
func (am *AuthManager) StoreTokenDetails(session *web.Session, oauth2Token *oauth2.Token, rawToken string, tokenExtras *domain.TokenExtras) {
	session.Store(keyToken, oauth2Token)
	session.Store(keyRawIDToken, rawToken)
	session.Store(keyTokenExtras, tokenExtras)
}

// DeleteTokenDetails deletes all token related data from session
func (am *AuthManager) DeleteTokenDetails(session *web.Session) {
	session.Delete(keyToken)
	session.Delete(keyRawIDToken)
	session.Delete(keyTokenExtras)
}

// StoreAuthState stores auth state into session, used to connect passed state id in auth callback with the one stored in session
func (am *AuthManager) StoreAuthState(session *web.Session, state string) {
	session.Store(keyAuthstate, state)
}

// LoadAuthState loads auth state from session
func (am *AuthManager) LoadAuthState(session *web.Session) (string, bool) {
	value, _ := session.Load(keyAuthstate)
	state, ok := value.(string)
	return state, ok
}

// DeleteAuthState deletes auth state from session
func (am *AuthManager) DeleteAuthState(session *web.Session) {
	session.Delete(keyAuthstate)
}
