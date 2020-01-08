package oauth

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"

	uuid "github.com/satori/go.uuid"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type (
	// OpenIDIdentity is an extension of Identity which provides an IDToken on top of OAuth2
	OpenIDIdentity interface {
		auth.Identity
		TokenSourcer
		IDToken() *oidc.IDToken
	}

	oidcIdentity struct {
		broker     string
		subject    string
		token      token
		verifier   *oidc.IDTokenVerifier
		rawIDToken string
	}

	openIDIdentifier struct {
		broker        string
		oauth2Config  *oauth2.Config
		responder     *web.Responder
		provider      *oidc.Provider
		reverseRouter web.ReverseRouter
	}

	sessionData struct {
		Subject    string
		Token      *oauth2.Token
		RawIDToken string
	}
)

func init() {
	gob.Register(sessionData{})
}

var _ OpenIDIdentity = new(oidcIdentity)

func oidcFactory(cfg config.Map) auth.RequestIdentifier {
	provider, err := oidc.NewProvider(context.Background(), cfg["endpoint"].(string))
	if err != nil {
		panic(err)
	}

	return &openIDIdentifier{
		oauth2Config: &oauth2.Config{
			ClientID:     cfg["clientID"].(string),
			ClientSecret: cfg["clientSecret"].(string),
			Endpoint:     provider.Endpoint(),
			RedirectURL:  "",
			Scopes:       []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess},
			ClaimSet:     nil,
		},
		broker:   cfg["broker"].(string),
		provider: provider,
	}
}

func (i *openIDIdentifier) sessionCode(s string) string {
	return "core.auth.oidc." + i.broker + "." + s
}

// Inject dependencies
func (i *openIDIdentifier) Inject(responder *web.Responder, reverseRouter web.ReverseRouter) {
	i.responder = responder
	i.reverseRouter = reverseRouter
}

// Identify an incoming request
func (i *openIDIdentifier) Identify(ctx context.Context, request *web.Request) auth.Identity {
	sessionCode := i.sessionCode("sessiondata")

	data, ok := request.Session().Load(sessionCode)
	if !ok {
		return nil
	}

	sessiondata, ok := data.(sessionData)
	if !ok {
		request.Session().Delete(sessionCode)
		return nil
	}

	identity := &oidcIdentity{
		token:      token{tokenSource: i.config(request).TokenSource(ctx, sessiondata.Token)},
		broker:     i.broker,
		subject:    sessiondata.Subject,
		verifier:   i.provider.Verifier(&oidc.Config{ClientID: i.oauth2Config.ClientID}),
		rawIDToken: sessiondata.RawIDToken,
	}

	token, idtoken := identity.tokens(ctx)

	request.Session().Store(sessionCode, sessionData{
		Token:      token,
		Subject:    idtoken.Subject,
		RawIDToken: identity.rawIDToken,
	})

	return identity
}

func (i *oidcIdentity) tokens(ctx context.Context) (*oauth2.Token, *oidc.IDToken) {
	token, err := i.token.tokenSource.Token()

	if err != nil {
		panic(err)
	}

	if idtoken, ok := token.Extra("id_token").(string); ok {
		i.rawIDToken = idtoken
	}

	idToken, err := i.verifier.Verify(ctx, i.rawIDToken)

	return token, idToken
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
	_, idtoken := i.tokens(context.Background())
	return idtoken
}

// TokenSource getter
func (i *oidcIdentity) TokenSource() oauth2.TokenSource {
	return i.token.TokenSource()
}

// String returns a readable token
func (i *oidcIdentity) String() string {
	return fmt.Sprintf("%s, expiry: %s", i.subject, i.IDToken().Expiry)
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
	u, err := url.Parse(i.config(request).AuthCodeURL(state, oauth2.AccessTypeOffline))
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

	// Extract custom claims
	// TODO
	var claims struct {
		Email    string `json:"email"`
		Verified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return i.responder.ServerError(err)
	}

	sessionCode := i.sessionCode("sessiondata")
	request.Session().Store(sessionCode, sessionData{
		Token:      oauth2Token,
		Subject:    idToken.Subject,
		RawIDToken: rawIDToken,
	})

	return i.responder.URLRedirect(returnTo(request))
}

// Logout based on a request
func (i *openIDIdentifier) Logout(ctx context.Context, request *web.Request) {
	request.Session().Delete(i.sessionCode("sessiondata"))
}
