package oauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofrs/uuid"
	"golang.org/x/oauth2"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// OpenIDIdentity is an extension of Identity which provides an IDToken on top of OAuth2
	OpenIDIdentity interface {
		auth.Identity
		Identity
		IDToken() *oidc.IDToken
		IDTokenClaims(into interface{}) error
	}

	// CallbackErrorHandler can be used to handle errors in the Callback e.g. to cover prompt=none cases
	CallbackErrorHandler interface {
		Handle(ctx context.Context, broker string, request *web.Request, originalReturnTo func(request *web.Request) *url.URL, errString string, errDetails string) web.Result
	}

	oidcIdentity struct {
		broker            string
		subject           string
		token             token
		verifier          *oidc.IDTokenVerifier
		rawIDToken        string
		idTokenClaims     []byte
		accessTokenClaims []byte
	}

	openIDIdentifier struct {
		broker                   string
		oauth2Config             *oauth2.Config
		responder                *web.Responder
		provider                 *oidc.Provider
		reverseRouter            web.ReverseRouter
		authcodeOptions          []AuthCodeOptioner
		authCodeOptionerProvider authCodeOptionerProvider
		eventRouter              flamingo.EventRouter
		oidcConfig               oidcConfig
		verifierConfigurator     []func(*oidc.Config)
		callbackErrorHandler     CallbackErrorHandler
		stateTimeout             *time.Duration
	}

	sessionData struct {
		Subject           string
		Token             *oauth2.Token
		RawIDToken        string
		IDTokenClaims     []byte
		AccessTokenClaims []byte
	}

	oidcConfig struct {
		Broker             string   `json:"broker"`
		Endpoint           string   `json:"endpoint"`
		ClientID           string   `json:"clientID"`
		ClientSecret       string   `json:"clientSecret"`
		Scopes             []string `json:"scopes"`
		EnableOfflineToken bool     `json:"enableOfflineToken"`
		Claimset           struct {
			IDToken  []string `json:"idToken"`
			UserInfo []string `json:"userInfo"`
		} `json:"requestClaims"`
		Claims struct {
			IDToken     map[string]string `json:"idToken"`
			AccessToken map[string]string `json:"accessToken"`
		} `json:"claims"`
		EnableEndSessionEndpoint bool   `json:"enableEndSessionEndpoint"`
		OverrideIssuerURL        string `json:"overrideIssuerURL"`
		StateLifeTime            string `json:"stateLifeTime"`
	}
)

func init() {
	gob.Register(sessionData{})
}

var (
	_ OpenIDIdentity = new(oidcIdentity)

	_ auth.RequestIdentifier     = new(openIDIdentifier)
	_ auth.WebAuthenticater      = new(openIDIdentifier)
	_ auth.WebCallbacker         = new(openIDIdentifier)
	_ auth.WebIdentityRefresher  = new(openIDIdentifier)
	_ auth.WebLogoutWithRedirect = new(openIDIdentifier)

	// OpenIDTypeChecker checks the Identity for OpenID Identity
	OpenIDTypeChecker = func(identity auth.Identity) bool {
		_, ok := identity.(OpenIDIdentity)

		return ok
	}
)

func oidcFactory(cfg config.Map) (auth.RequestIdentifier, error) {
	var oidcConfig oidcConfig

	if err := cfg.MapInto(&oidcConfig); err != nil {
		return nil, err
	}

	ctx := context.Background()
	if oidcConfig.OverrideIssuerURL != "" {
		ctx = oidc.InsecureIssuerURLContext(ctx, oidcConfig.OverrideIssuerURL)
	}

	provider, err := oidc.NewProvider(ctx, oidcConfig.Endpoint)
	if err != nil {
		return nil, err
	}

	scopes := append([]string{oidc.ScopeOpenID}, oidcConfig.Scopes...)
	if oidcConfig.EnableOfflineToken {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
	}

	var authCodeOptions []AuthCodeOptioner
	if claimset := getClaimset(oidcConfig); claimset.HasClaims() {
		authCodeOption, err := claimset.AuthCodeOption()
		if err != nil {
			return nil, err
		}
		authCodeOptions = append(authCodeOptions, oauth2AuthCodeOption{authCodeOption: authCodeOption})
	}

	stateTimeout := defaultStateTimeout

	if oidcConfig.StateLifeTime != "" {
		if duration, err := time.ParseDuration(oidcConfig.StateLifeTime); err == nil {
			stateTimeout = duration
		}
	}

	return &openIDIdentifier{
		oauth2Config: &oauth2.Config{
			ClientID:     oidcConfig.ClientID,
			ClientSecret: oidcConfig.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  "", // filled on request
			Scopes:       scopes,
		},
		broker:          oidcConfig.Broker,
		provider:        provider,
		authcodeOptions: authCodeOptions,
		oidcConfig:      oidcConfig,
		stateTimeout:    &stateTimeout,
	}, nil
}

func getClaimset(oidcConfig oidcConfig) *ClaimSet {
	var claimSet *ClaimSet

	claimSet = createClaimSetFromMapping(TopLevelClaimIDToken, oidcConfig.Claimset.IDToken, claimSet)
	claimSet = createClaimSetFromMapping(TopLevelClaimUserInfo, oidcConfig.Claimset.UserInfo, claimSet)

	return claimSet
}

func createClaimSetFromMapping(topLevelName string, mapping []string, claimSet *ClaimSet) *ClaimSet {
	for _, name := range mapping {
		if name == "" {
			continue
		}
		if claimSet == nil {
			claimSet = &ClaimSet{}
		}
		claimSet.AddVoluntaryClaim(topLevelName, name)
	}

	return claimSet
}

func (i *openIDIdentifier) sessionCode(s string) string {
	return "core.auth.oidc." + i.broker + "." + s
}

// Inject dependencies
func (i *openIDIdentifier) Inject(
	responder *web.Responder,
	reverseRouter web.ReverseRouter,
	eventRouter flamingo.EventRouter,
	authCodeOptionerProvider authCodeOptionerProvider,
	optionals *struct {
		CallbackErrorHandler CallbackErrorHandler `inject:",optional"`
	},
) {
	i.responder = responder
	i.reverseRouter = reverseRouter
	i.eventRouter = eventRouter
	i.authCodeOptionerProvider = authCodeOptionerProvider

	if optionals != nil && optionals.CallbackErrorHandler != nil {
		i.callbackErrorHandler = optionals.CallbackErrorHandler
	}
}

// Identify an incoming request
func (i *openIDIdentifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	sessionCode := i.sessionCode("sessiondata")

	sessiondata, err := i.sessionData(request, sessionCode)
	if err != nil {
		return nil, err
	}

	verifierConfig := &oidc.Config{ClientID: i.oauth2Config.ClientID}
	for _, configurator := range i.verifierConfigurator {
		configurator(verifierConfig)
	}

	authConfig, err := i.config(request)
	if err != nil {
		return nil, fmt.Errorf("failed to identify a reuqest: %w", err)
	}

	identity := &oidcIdentity{
		token:             token{tokenSource: authConfig.TokenSource(ctx, sessiondata.Token)},
		broker:            i.broker,
		subject:           sessiondata.Subject,
		verifier:          i.provider.Verifier(verifierConfig),
		rawIDToken:        sessiondata.RawIDToken,
		idTokenClaims:     sessiondata.IDTokenClaims,
		accessTokenClaims: sessiondata.AccessTokenClaims,
	}

	token, idtoken, err := identity.tokens(ctx)
	if err != nil {
		request.Session().Delete(sessionCode)
		return nil, err
	}

	request.Session().Store(sessionCode, sessionData{
		Token:             token,
		Subject:           idtoken.Subject,
		RawIDToken:        identity.rawIDToken,
		IDTokenClaims:     sessiondata.IDTokenClaims,
		AccessTokenClaims: sessiondata.AccessTokenClaims,
	})

	return identity, nil
}

func (i *openIDIdentifier) sessionData(request *web.Request, sessionCode string) (sessionData, error) {
	data, ok := request.Session().Load(sessionCode)
	if !ok {
		return sessionData{}, errors.New("no sessiondata")
	}

	sessiondata, ok := data.(sessionData)
	if !ok {
		request.Session().Delete(sessionCode)
		return sessionData{}, errors.New("broken sessiondata")
	}

	return sessiondata, nil
}

// RefreshIdentity by invalidating the access token from the token stored in the session data
// which will cause a refresh request the next time an identity is requested
func (i *openIDIdentifier) RefreshIdentity(_ context.Context, request *web.Request) error {
	sessionCode := i.sessionCode("sessiondata")

	sessiondata, err := i.sessionData(request, sessionCode)
	if err != nil {
		return err
	}

	sessiondata.Token.AccessToken = ""
	request.Session().Store(sessionCode, sessiondata)

	return nil
}

func (i *oidcIdentity) tokens(ctx context.Context) (*oauth2.Token, *oidc.IDToken, error) {
	token, err := i.token.tokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	if idtoken, ok := token.Extra("id_token").(string); ok {
		i.rawIDToken = idtoken
	}

	idToken, err := i.verifier.Verify(ctx, i.rawIDToken)
	if err != nil {
		return nil, nil, err
	}

	return token, idToken, nil
}

// Broker information
func (i *oidcIdentity) Broker() string {
	return i.broker
}

// Subject getter
func (i *oidcIdentity) Subject() string {
	return i.subject
}

// IDToken getter
func (i *oidcIdentity) IDToken() *oidc.IDToken {
	_, idtoken, _ := i.tokens(context.Background()) // TODO background ok?
	return idtoken
}

// IDTokenClaims mapper
func (i *oidcIdentity) IDTokenClaims(into interface{}) error {
	return json.Unmarshal(i.idTokenClaims, into)
}

// AccessTokenClaims mapper
func (i *oidcIdentity) AccessTokenClaims(into interface{}) error {
	return json.Unmarshal(i.accessTokenClaims, into)
}

// TokenSource getter
func (i *oidcIdentity) TokenSource() oauth2.TokenSource {
	return i.token.TokenSource()
}

// String returns a readable token
func (i *oidcIdentity) String() string {
	return fmt.Sprintf("%s, (%s) expiry: %s", i.subject, string(i.idTokenClaims), i.IDToken().Expiry)
}

// RawIDToken returns the actual raw ID token for further processing
func (i *oidcIdentity) RawIDToken() string {
	return i.rawIDToken
}

// Broker getter
func (i *openIDIdentifier) Broker() string {
	return i.broker
}

func (i *openIDIdentifier) config(request *web.Request) (*oauth2.Config, error) {
	oauth2Config := *i.oauth2Config

	u, err := i.reverseRouter.Absolute(request, "core.auth.callback", map[string]string{"broker": i.broker})
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth config: %w", err)
	}

	oauth2Config.RedirectURL = u.String()

	return &oauth2Config, nil
}

// StateEntry stores entries of recent states during login (oidc states)
type StateEntry struct {
	State string
	TS    time.Time
}

const defaultStateTimeout = time.Minute * 30

func init() {
	gob.Register([]StateEntry(nil))
	gob.Register(StateEntry{})
}

const sessionStatesKey = "states"

var now = time.Now

func (i *openIDIdentifier) validateSessionCode(request *web.Request, code string) bool {
	stateTimeout := defaultStateTimeout

	if i.stateTimeout != nil {
		stateTimeout = *i.stateTimeout
	}

	sessionStates, ok := request.Session().Load(i.sessionCode(sessionStatesKey))
	if !ok {
		return false
	}
	states, ok := sessionStates.([]StateEntry)
	if !ok {
		return false
	}
	newStates := make([]StateEntry, 0, len(states))
	validated := false
	for _, state := range states {
		if state.TS.Add(stateTimeout).Before(now()) {
			continue
		}
		if state.State == code {
			validated = true
			continue
		}
		newStates = append(newStates, state)
	}
	request.Session().Store(i.sessionCode(sessionStatesKey), newStates)
	return validated
}

func (i *openIDIdentifier) createSessionCode(request *web.Request, code string) {
	sessionStates, ok := request.Session().Load(i.sessionCode(sessionStatesKey))
	if !ok {
		sessionStates = []StateEntry{}
	}
	states := sessionStates.([]StateEntry)
	states = append(states, StateEntry{
		State: code,
		TS:    now(),
	})
	request.Session().Store(i.sessionCode(sessionStatesKey), states)
}

// Authenticate a user
func (i *openIDIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	state := uuid.Must(uuid.NewV4()).String()
	i.createSessionCode(request, state)
	options := make([]oauth2.AuthCodeOption, 0, len(i.authcodeOptions))
	for _, o := range i.authcodeOptions {
		options = append(options, o.Options(ctx, i.Broker(), request)...)
	}
	for _, o := range i.authCodeOptionerProvider() {
		options = append(options, o.Options(ctx, i.Broker(), request)...)
	}

	authConfig, err := i.config(request)
	if err != nil {
		return i.responder.ServerError(err)
	}

	u, err := url.Parse(authConfig.AuthCodeURL(state, options...))
	if err != nil {
		return i.responder.ServerError(err)
	}

	return i.responder.URLRedirect(u)
}

// Callback for OIDC code exchange
func (i *openIDIdentifier) Callback(ctx context.Context, request *web.Request, returnTo func(request *web.Request) *url.URL) web.Result {
	if errString, err := request.Query1("error"); err == nil {
		errDetails, _ := request.Query1("error_description")
		if i.callbackErrorHandler != nil {
			if result := i.callbackErrorHandler.Handle(ctx, i.broker, request, returnTo, errString, errDetails); result != nil {
				return result
			}
		}

		return i.responder.ServerError(fmt.Errorf("OpenID Connect error: %q (%q)", errString, errDetails))
	}

	queryState, err := request.Query1("state")
	if err != nil {
		return i.responder.ServerError(errors.New("no state in request"))
	}
	if !i.validateSessionCode(request, queryState) {
		return i.responder.ServerError(errors.New("state mismatch"))
	}

	code, err := request.Query1("code")
	if err != nil {
		return i.responder.ServerError(err)
	}

	oauthConfig, err := i.config(request)
	if err != nil {
		return i.responder.ServerError(err)
	}

	options := make([]oauth2.AuthCodeOption, 0)

	if i.authCodeOptionerProvider != nil {
		for _, o := range i.authCodeOptionerProvider() {
			options = append(options, o.Options(ctx, i.Broker(), request)...)
		}
	}

	oauth2Token, err := oauthConfig.Exchange(ctx, code, options...)
	if err != nil {
		return i.responder.ServerError(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return i.responder.ServerError(errors.New("claim id_token missing"))
	}

	verifierConfig := &oidc.Config{ClientID: i.oauth2Config.ClientID}
	for _, configurator := range i.verifierConfigurator {
		configurator(verifierConfig)
	}
	verifier := i.provider.Verifier(verifierConfig)

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return i.responder.ServerError(err)
	}

	var (
		idTokenClaims     = make(map[string]interface{})
		tempIDTokenClaims = make(map[string]interface{})
		accessTokenClaims = make(map[string]interface{})
	)

	if err := idToken.Claims(&tempIDTokenClaims); err != nil {
		return i.responder.ServerError(err)
	}
	for k, v := range i.oidcConfig.Claims.IDToken {
		idTokenClaims[k] = tempIDTokenClaims[v]
	}

	accessTokenParts := strings.Split(oauth2Token.AccessToken, ".")
	if len(accessTokenParts) >= 2 {
		decoded, _ := base64.RawURLEncoding.DecodeString(accessTokenParts[1])
		var cs map[string]interface{}
		_ = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(&cs)
		if cs != nil {
			for k, v := range i.oidcConfig.Claims.AccessToken {
				accessTokenClaims[k] = cs[v]
			}
		}
	}
	for k, v := range i.oidcConfig.Claims.AccessToken {
		if oauth2Token.Extra(v) != nil {
			accessTokenClaims[k] = oauth2Token.Extra(v)
		}
	}

	itc, _ := json.Marshal(idTokenClaims)
	atc, _ := json.Marshal(accessTokenClaims)

	sessionCode := i.sessionCode("sessiondata")
	request.Session().Store(sessionCode, sessionData{
		Token:             oauth2Token,
		Subject:           idToken.Subject,
		RawIDToken:        rawIDToken,
		IDTokenClaims:     itc,
		AccessTokenClaims: atc,
	})

	identity, err := i.Identify(ctx, request)
	if err != nil {
		i.Logout(ctx, request)
		return i.responder.ServerError(err)
	}

	i.eventRouter.Dispatch(ctx, &auth.WebLoginEvent{Broker: i.broker, Request: request, Identity: identity})

	return i.responder.URLRedirect(returnTo(request))
}

// Logout based on a request
func (i *openIDIdentifier) Logout(ctx context.Context, request *web.Request) *url.URL {
	identity, err := i.Identify(ctx, request)
	request.Session().Delete(i.sessionCode("sessiondata"))

	// return if we are not logged in
	if identity == nil || err != nil || !i.oidcConfig.EnableEndSessionEndpoint {
		return nil
	}

	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}
	// we can ignore errors here as we are just fine handling default values
	if err := i.OpenIDConnectProvider().Claims(&claims); err != nil {
		return nil
	}
	if claims.EndSessionEndpoint == "" {
		return nil
	}
	returnURL, err := url.Parse(claims.EndSessionEndpoint)
	if err != nil {
		return nil
	}
	query := returnURL.Query()
	query.Set("id_token_hint", identity.(*oidcIdentity).rawIDToken)
	returnURL.RawQuery = query.Encode()
	return returnURL
}

// OpenIDConnectProvider getter for openID Connect Provider
func (i *openIDIdentifier) OpenIDConnectProvider() *oidc.Provider {
	return i.provider
}
