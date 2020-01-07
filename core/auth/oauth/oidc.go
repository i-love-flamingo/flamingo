package oauth

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/url"

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
		IDToken() oidc.IDToken
	}

	oidcIdentity struct {
		broker  string
		subject string
		token   token
		idToken *oidc.IDToken
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

func oidcFactory(cfg config.Map) auth.Identifier {
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
			Scopes:       []string{oidc.ScopeOpenID},
			ClaimSet:     nil,
		},
		broker:   cfg["broker"].(string),
		provider: provider,
	}
}

func (i *openIDIdentifier) Inject(responder *web.Responder, reverseRouter web.ReverseRouter) {
	i.responder = responder
	i.reverseRouter = reverseRouter
}

func (i *openIDIdentifier) Identify(ctx context.Context, request *web.Request) auth.Identity {
	sessionCode := "core.auth.oidc." + i.broker + ".sessiondata"

	data, ok := request.Session().Load(sessionCode)
	if !ok {
		return nil
	}

	sessiondata, ok := data.(sessionData)
	if !ok {
		request.Session().Delete(sessionCode)
		return nil
	}

	verifier := i.provider.Verifier(&oidc.Config{ClientID: i.oauth2Config.ClientID})
	idToken, err := verifier.Verify(ctx, sessiondata.RawIDToken)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &oidcIdentity{
		token: token{
			config: i.config(request),
			token:  sessiondata.Token,
		},
		broker:  i.broker,
		subject: sessiondata.Subject,
		idToken: idToken,
	}
}

func (i *oidcIdentity) Broker() string {
	return i.broker
}

func (i *oidcIdentity) Subject() string {
	return i.subject
}

func (i *oidcIdentity) IDToken() *oidc.IDToken {
	return i.idToken
}

func (i *openIDIdentifier) Broker() string {
	return i.broker
}

func (i *openIDIdentifier) config(request *web.Request) *oauth2.Config {
	oauth2Config := *i.oauth2Config
	u, _ := i.reverseRouter.Absolute(request, "core.auth.callback", map[string]string{"broker": i.broker})
	oauth2Config.RedirectURL = u.String()
	return &oauth2Config
}

func (i *openIDIdentifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
	u, err := url.Parse(i.config(request).AuthCodeURL("state", oauth2.AccessTypeOffline))
	if err != nil {
		return i.responder.ServerError(err)
	}

	return i.responder.URLRedirect(u)
}

func (i *openIDIdentifier) Callback(ctx context.Context, request *web.Request, returnTo func(request *web.Request) *url.URL) web.Result {
	errString, err := request.Query1("error")
	if err == nil {
		errDetails, _ := request.Query1("error_description")
		return i.responder.ServerError(fmt.Errorf("OpenID Connect error: %q (%q)", errString, errDetails))
	}

	code, err := request.Query1("code")
	if err != nil {
		return i.responder.ServerError(err)
	}

	// Verify state and errors.
	// TODO verify state
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

	sessionCode := "core.auth.oidc." + i.broker + ".sessiondata"
	request.Session().Store(sessionCode, sessionData{
		Token:      oauth2Token,
		Subject:    idToken.Subject,
		RawIDToken: rawIDToken,
	})

	return i.responder.URLRedirect(returnTo(request))
}
