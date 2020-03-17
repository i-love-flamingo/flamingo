package application

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"

	"flamingo.me/flamingo/v3/core/auth/oauth"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
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
	gob.Register(oauth2.Token{})
	gob.Register(domain.TokenExtras{})
}

type (
	// AuthManager handles authentication related operations
	// Deprecated: use core/auth instead
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
		tokenExtras         config.Slice
		AuthCodeOptions     []oauth2.AuthCodeOption
	}

	loggingRoundTripper struct {
		originalTransport http.RoundTripper
	}
)

// RoundTrip implements RoundTripper interface and adds logging
func (f *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("No request given")
	}
	b, err := httputil.DumpRequest(req, true)
	log.Println()
	log.Println("############### OAUTH REQUEST:")
	log.Printf("%v  %v ", string(b), err)
	res, err := f.originalTransport.RoundTrip(req)
	b, _ = httputil.DumpResponse(res, true)
	log.Println("############### OAUTH RESPONSE:")
	log.Printf("%v  %v ", string(b), err)
	log.Println("############### OAUTH Call Stack:")
	log.Println()
	debug.PrintStack()
	return res, err
}

// Inject authManager dependencies
func (am *AuthManager) Inject(logger flamingo.Logger, router *web.Router, config *struct {
	Server              string       `inject:"config:core.oauth.server"`
	Secret              string       `inject:"config:core.oauth.secret"`
	ClientID            string       `inject:"config:core.oauth.clientid"`
	DisableOfflineToken bool         `inject:"config:core.oauth.disableOfflineToken"`
	Scopes              config.Slice `inject:"config:core.oauth.scopes"`
	IDTokenMapping      config.Slice `inject:"config:core.oauth.claims.idToken"`
	UserInfoMapping     config.Slice `inject:"config:core.oauth.claims.userInfo"`
	TokenExtras         config.Slice `inject:"config:core.oauth.tokenExtras"`
	DebugMode           bool         `inject:"config:flamingo.debug.mode"`
	Enabled             bool         `inject:"config:core.oauth.enabled"`
}) {
	am.logger = logger.WithField(flamingo.LogKeyModule, "oauth")
	am.router = router
	if config != nil {
		am.server = config.Server
		am.secret = config.Secret
		am.clientID = config.ClientID
		am.disableOfflineToken = config.DisableOfflineToken
		am.scopes = config.Scopes
		am.idTokenMapping = config.IDTokenMapping
		am.userInfoMapping = config.UserInfoMapping
		am.tokenExtras = config.TokenExtras
		if !config.Enabled {
			am.logger.Info("OIDC is disabled. Modules depending on OAuth features cannot work properly")
			return
		}

		var err error
		am.openIDProvider, err = oidc.NewProvider(context.Background(), config.Server)
		if err != nil {
			if config.DebugMode {
				am.logger.Error(err)
			} else {
				//panic on err since we really expect a valid authmanager state and application is in a failed state otherwise
				panic(err)
			}
		}
	}
}

// Auth tries to retrieve the authentication context for a active session - this is used to pass Authentication to services
//	- if the stored token for the Auth is not valid anymore it will refresh the token before
func (am *AuthManager) Auth(c context.Context, session *web.Session) (domain.Auth, error) {
	c = am.OAuthCtx(c)
	currentToken, err := am.OAuth2Token(session)
	if err != nil {
		am.logger.WithContext(c).Debug(err)
		return domain.Auth{}, err
	}
	if !currentToken.Valid() {
		err := am.refreshTokenAndUpdateStore(c, session)
		if err != nil {
			am.logger.WithContext(c).Error(err)
			return domain.Auth{}, err
		}
	}
	ts, err := am.TokenSource(c, session)
	if err != nil {
		am.logger.WithContext(c).Error(err)
		return domain.Auth{}, err
	}
	idToken, err := am.IDToken(c, session)
	if err != nil {
		am.logger.WithContext(c).Error(err)
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

//OAuthCtx - returns ctx that should be used to pass to oauth2 lib - it enables logging for Debug reasons
func (am *AuthManager) OAuthCtx(ctx context.Context) context.Context {
	if os.Getenv("OAUTHDEBUG") == "1" {
		oauthHTTPClient := &http.Client{
			Transport: &loggingRoundTripper{
				originalTransport: http.DefaultTransport,
			},
		}
		return context.WithValue(ctx, oauth2.HTTPClient, oauthHTTPClient)
	}
	return ctx
}

// OAuth2Config is lazy setup oauth2config
func (am *AuthManager) OAuth2Config(_ context.Context, req *web.Request) *oauth2.Config {

	var redirectURL string
	if req != nil {
		callbackURL, _ := am.router.Absolute(req, "auth.callback", nil)
		redirectURL = callbackURL.String()
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

	claimset := am.getClaimsRequestParameter()
	if claimset.HasClaims() {
		authCodeOption, err := am.getClaimsRequestParameter().AuthCodeOption()
		if err != nil {
			am.logger.WithField(flamingo.LogKeyCategory, "auth").Error("could not map configuration", err)
		} else {
			am.AuthCodeOptions = append(am.AuthCodeOptions, authCodeOption)
		}
	}

	oauth2Config := &oauth2.Config{
		ClientID:     am.clientID,
		ClientSecret: am.secret,
		RedirectURL:  redirectURL,

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: scopes,
	}
	if am.OpenIDProvider() != nil {
		oauth2Config.Endpoint = am.OpenIDProvider().Endpoint()
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
	oauth2Token, ok := value.(oauth2.Token)
	if !ok {
		am.DeleteTokenDetails(session)
		return nil, fmt.Errorf("invalid token in session %#v", value)
	}

	return &oauth2Token, nil
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
	if am.OpenIDProvider() == nil {
		return nil, "", errors.New("oauth not available")
	}
	c = am.OAuthCtx(c)
	if session == nil {
		return nil, "", errors.New("no session configured")
	}

	if token, ok := session.Load(keyRawIDToken); ok {
		idtoken, err := am.Verifier().Verify(c, token.(string))
		if err == nil {
			return idtoken, token.(string), nil
		}
		am.logger.WithContext(c).Debug("keyRawIDToken not verified (anymore)")
		err = am.refreshTokenAndUpdateStore(c, session)
		if err != nil {
			return nil, "", err
		}
		token, ok = session.Load(keyRawIDToken)
		if !ok {
			return nil, "", errors.New("no token after refreshToken")
		}
		idtoken, err = am.Verifier().Verify(c, token.(string))
		if err != nil {
			return nil, "", errors.New("no verified id token after refreshToken")
		}
		return idtoken, token.(string), nil

	}

	return nil, "", errors.New("no id token in session")

}

// refreshTokenAndUpdateStore
func (am *AuthManager) refreshTokenAndUpdateStore(c context.Context, session *web.Session) error {
	c = am.OAuthCtx(c)
	tokenSource, err := am.TokenSource(c, session)
	if err != nil {
		return fmt.Errorf("error retrieving tokenSource: %w", err)
	}

	token, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("error retrieving token: %w", err)
	}

	err = am.StoreTokenDetails(c, session, token)
	if err != nil {
		return fmt.Errorf("error storing token: %w", err)
	}
	return nil
}

func (am *AuthManager) getClaimsRequestParameter() *oauth.ClaimSet {
	var claimSet *oauth.ClaimSet

	claimSet = am.createClaimSetFromMapping(oauth.TopLevelClaimIDToken, am.idTokenMapping, claimSet)
	claimSet = am.createClaimSetFromMapping(oauth.TopLevelClaimUserInfo, am.userInfoMapping, claimSet)

	return claimSet
}

func (am *AuthManager) createClaimSetFromMapping(topLevelName string, configuration config.Slice, claimSet *oauth.ClaimSet) *oauth.ClaimSet {
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
			claimSet = &oauth.ClaimSet{}
		}
		claimSet.AddVoluntaryClaim(topLevelName, name)
	}

	return claimSet
}

// AccessToken - used to get access token
func (am *AuthManager) AccessToken(ctx context.Context, session *web.Session) (string, error) {
	auth, err := am.Auth(ctx, session)
	if err != nil {
		return "", err
	}
	token, err := auth.TokenSource.Token()
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// ExtractRawIDToken from the provided (fresh) oatuh2token
func (am *AuthManager) ExtractRawIDToken(oauth2Token *oauth2.Token) (string, error) {
	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("no id token %T / %v", oauth2Token.Extra("id_token"), oauth2Token.Extra("id_token"))
	}

	return rawIDToken, nil
}

// TokenSource - return oauth2.TokenSource initialized with the Refreshtoken stored in the
// to be used in situations where you need it
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
func (am *AuthManager) StoreTokenDetails(ctx context.Context, session *web.Session, oauth2Token *oauth2.Token) error {
	if oauth2Token == nil {
		err := errors.New("StoreTokenDetails got no token")
		am.logger.WithContext(ctx).Error(err)
		return err
	}
	if oauth2Token.AccessToken == "" {
		err := errors.New("StoreTokenDetails got token without accesstoken")
		am.logger.WithContext(ctx).Error(err)
		return err
	}
	if !oauth2Token.Valid() {
		am.logger.WithContext(ctx).Warn("StoreTokenDetails got already invalid token")
	}
	rawToken, err := am.ExtractRawIDToken(oauth2Token)
	if err != nil {
		am.logger.Error("core.auth.callback Error ExtractRawIDToken", err)
		return err
	}

	var extras []string
	err = am.tokenExtras.MapInto(&extras)
	if err != nil {
		return err
	}
	tokenExtras := domain.TokenExtras{}
	for _, extra := range extras {
		value := oauth2Token.Extra(extra)
		parsed, ok := value.(string)
		if !ok {
			am.logger.Error("core.auth.callback invalid type for extras", value)
			continue
		}
		tokenExtras.Add(extra, parsed)
	}

	var token oauth2.Token
	token = *oauth2Token
	session.Store(keyToken, token)
	session.Store(keyRawIDToken, rawToken)
	session.Store(keyTokenExtras, tokenExtras)
	return nil
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
