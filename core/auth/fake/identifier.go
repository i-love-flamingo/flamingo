package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	Identifier struct {
		responder *web.Responder
		broker    string
	}
)

const fakeAuthURL string = "/fake/auth"

var (
	_ auth.RequestIdentifier = (*Identifier)(nil)
)

func (i *Identifier) Broker() string {
	return i.broker
}

func (i *Identifier) Authenticate(ctx context.Context, _ *web.Request) web.Result {
	authUrl, _ := url.Parse(fakeAuthURL)
	authUrl.Query().Add("broker", i.Broker())
	return i.responder.URLRedirect(authUrl)
}

func (i *Identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	fakeSubject := "" // TODO

	return &Identity{
		subject: fakeSubject,
		broker:  i.broker,
	}, nil
}
