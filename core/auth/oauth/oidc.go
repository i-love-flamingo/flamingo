package oauth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type (
	// OpenIDIdentity is an extension of Identity which provides an IDToken on top of OAuth2
	OpenIDIdentity interface {
		auth.Identity
		Identity
		IDToken() *oidc.IDToken
		IDTokenClaims(into interface{}) error
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
	}

	sessionData struct {
		Subject           string
		Token             *oauth2.Token
		RawIDToken        string
		IDTokenClaims     []byte
		AccessTokenClaims []byte
	}
)

func init() {
	gob.Register(sessionData{})
}

var _ OpenIDIdentity = new(oidcIdentity)

type oidcConfig struct {
	Broker              string   `json:"broker"`
	Endpoint            string   `json:"endpoint"`
	ClientID            string   `json:"clientID"`
	ClientSecret        string   `json:"clientSecret"`
	Scopes              []string `json:"scopes"`
	EnabledOfflineToken bool     `json:"enabledOfflineToken"`
	Claimset            struct {
		IDToken  []string `json:"idToken"`
		UserInfo []string `json:"userInfo"`
	} `json:"requestClaims"`
	Claims struct {
		IDToken     map[string]string `json:"idToken"`
		AccessToken map[string]string `json:"accessToken"`
	} `json:"claims"`
}

func oidcFactory(cfg config.Map) (auth.RequestIdentifier, error) {
	var oidcConfig oidcConfig

	if err := cfg.MapInto(&oidcConfig); err != nil {
		return nil, err
	}

	provider, err := oidc.NewProvider(context.Background(), oidcConfig.Endpoint)
	if err != nil {
		return nil, err
	}

	var authCodeOptions []AuthCodeOptioner

	scopes := append([]string{oidc.ScopeOpenID}, oidcConfig.Scopes...)
	if oidcConfig.EnabledOfflineToken {
		scopes = append(scopes, oidc.ScopeOfflineAccess)
	}

	if claimset := getClaimset(oidcConfig); claimset.HasClaims() {
		authCodeOption, err := claimset.AuthCodeOption()
		if err != nil {
			return nil, err
		}
		authCodeOptions = append(authCodeOptions, oauth2AuthCodeOption{authCodeOption: authCodeOption})
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
) {
	i.responder = responder
	i.reverseRouter = reverseRouter
	i.eventRouter = eventRouter
	i.authCodeOptionerProvider = authCodeOptionerProvider
}

// Identify an incoming request
func (i *openIDIdentifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	sessionCode := i.sessionCode("sessiondata")

	data, ok := request.Session().Load(sessionCode)
	if !ok {
		return nil, errors.New("no sessiondata")
	}

	sessiondata, ok := data.(sessionData)
	if !ok {
		request.Session().Delete(sessionCode)
		return nil, errors.New("no sessiondata")
	}

	identity := &oidcIdentity{
		token:             token{tokenSource: i.config(request).TokenSource(ctx, sessiondata.Token)},
		broker:            i.broker,
		subject:           sessiondata.Subject,
		verifier:          i.provider.Verifier(&oidc.Config{ClientID: i.oauth2Config.ClientID}),
		rawIDToken:        sessiondata.RawIDToken,
		idTokenClaims:     sessiondata.IDTokenClaims,
		accessTokenClaims: sessiondata.AccessTokenClaims,
	}

	token, idtoken, err := identity.tokens(ctx)
	if err != nil {
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
	_, idtoken, _ := i.tokens(context.Background())
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

// Broker getter
func (i *openIDIdentifier) Broker() string {
	return i.broker
}

func (i *openIDIdentifier) config(request *web.Request) *oauth2.Config {
	oauth2Config := *i.oauth2Config
	u, _ := i.reverseRouter.Absolute(request, "core.auth.callback", map[string]string{"broker": i.broker})
	oauth2Config.RedirectURL = u.String()
	return &oauth2Config
}

// Authenticate a user
func (i *openIDIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	state := uuid.NewV4().String()
	request.Session().Store(i.sessionCode("state"), state)
	options := make([]oauth2.AuthCodeOption, 0, len(i.authcodeOptions))
	for _, o := range i.authcodeOptions {
		options = append(options, o.Options(ctx, i.Broker(), request)...)
	}
	for _, o := range i.authCodeOptionerProvider() {
		options = append(options, o.Options(ctx, i.Broker(), request)...)
	}
	u, err := url.Parse(i.config(request).AuthCodeURL(state, options...))
	if err != nil {
		return i.responder.ServerError(err)
	}

	return i.responder.URLRedirect(u)
}

// Callback for OIDC code exchange
func (i *openIDIdentifier) Callback(ctx context.Context, request *web.Request, returnTo func(request *web.Request) *url.URL) web.Result {
	errString, err := request.Query1("error")
	if err == nil {
		errDetails, _ := request.Query1("error_description")
		return i.responder.ServerError(fmt.Errorf("OpenID Connect error: %q (%q)", errString, errDetails))
	}

	state, ok := request.Session().Load(i.sessionCode("state"))
	if !ok {
		return i.responder.ServerError(errors.New("no state in session"))
	}
	if queryState, err := request.Query1("state"); err != nil || queryState != state.(string) {
		return i.responder.ServerError(errors.New("state mismatch"))
	}
	request.Session().Delete(i.sessionCode("state"))

	code, err := request.Query1("code")
	if err != nil {
		return i.responder.ServerError(err)
	}

	oauth2Token, err := i.config(request).Exchange(ctx, code)
	if err != nil {
		return i.responder.ServerError(err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return i.responder.ServerError(errors.New("claim id_token missing"))
	}

	verifier := i.provider.Verifier(&oidc.Config{ClientID: i.oauth2Config.ClientID})

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
	for k, v := range i.oidcConfig.Claims.AccessToken {
		accessTokenClaims[k] = oauth2Token.Extra(v)
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
		return i.responder.ServerError(err)
	}

	i.eventRouter.Dispatch(ctx, &auth.WebLoginEvent{Broker: i.broker, Request: request, Identity: identity})

	return i.responder.URLRedirect(returnTo(request))
}

// Logout based on a request
func (i *openIDIdentifier) Logout(ctx context.Context, request *web.Request) {
	request.Session().Delete(i.sessionCode("sessiondata"))
}

// OpenIDConnectProvder getter for openID Connect Provider
func (i *openIDIdentifier) OpenIDConnectProvider() *oidc.Provider {
	return i.provider
}
