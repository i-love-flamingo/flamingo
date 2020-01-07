package oauth

import (
	"encoding/gob"

	"golang.org/x/oauth2"
)

type (
	// TokenSourcer defines a TokenSource which is can be used to get an AccessToken vor OAuth2 flows
	TokenSourcer interface {
		TokenSource() oauth2.TokenSource
	}

	token struct {
		tokenSource oauth2.TokenSource
	}

	//identifier struct {
	//	config    *oauth2.Config
	//	broker    string
	//	responder *web.Responder
	//}
)

func init() {
	gob.Register(oauth2.Token{})
}

func (i token) TokenSource() oauth2.TokenSource {
	return i.tokenSource
}

//
//func oauth2Factory(cfg config.Map) auth.Identifier {
//	return &identifier{
//		config: &oauth2.Config{
//			ClientID:     cfg["clientID"].(string),
//			ClientSecret: cfg["clientSecret"].(string),
//			Endpoint:     oauth2.Endpoint{},
//			RedirectURL:  "",
//			Scopes:       nil,
//			ClaimSet:     nil,
//		},
//		broker: cfg["broker"].(string),
//	}
//}
//
//func (i *identifier) Inject(responder *web.Responder) {
//	i.responder = responder
//}
//
//func (i *identifier) Identify(ctx context.Context, request *web.Request) auth.Identity {
//	sessionCode := "core.auth.oauth." + i.broker + ".sessiondata"
//
//	data, ok := request.Session().Load(sessionCode)
//	if !ok {
//		return nil
//	}
//
//	sessiondata, ok := data.(sessionData)
//	if !ok {
//		request.Session().Delete(sessionCode)
//	}
//
//	return token{
//		accessToken:  sessiondata.AccessToken,
//		refreshToken: sessiondata.RefreshToken,
//		roles:        sessiondata.Roles,
//		subject:      sessiondata.Subject,
//		broker:       i.broker,
//	}
//}
//
//func (i *identifier) Broker() string {
//	return i.broker
//}
//
//func (i *identifier) Authenticate(ctx context.Context, request *web.Request) web.Result {
//	u, err := url.Parse(i.config.AuthCodeURL("state", oauth2.AccessTypeOffline))
//	if err != nil {
//		return i.responder.ServerError(err)
//	}
//
//	return i.responder.URLRedirect(u)
//}
//
//func (i *identifier) Callback(ctx context.Context, request *web.Request) web.Result {
//	code, err := request.Query1("code")
//	if err != nil {
//		return i.responder.ServerError(err)
//	}
//
//	token, err := i.config.Exchange(ctx, code)
//	if err != nil {
//		return i.responder.ServerError(err)
//	}
//
//	request.Session().Store("core.oauth."+i.broker+".sessiondata", sessionData{
//		AccessToken:  token.AccessToken,
//		RefreshToken: token.RefreshToken,
//		Roles:        nil,
//		Subject:      token.Extra("email").(string),
//	})
//
//	return nil
//}
