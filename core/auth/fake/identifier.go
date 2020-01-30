package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	Identifier struct {
		broker string
	}

	Authenticator struct {
		responder *web.Responder
	}
)

const fakeAuthURL string = "/fake/auth"

var (
	_ auth.RequestIdentifier = (*Identifier)(nil)
)

func (i *Identifier) Broker() string {
	return i.broker
}

func (a *Authenticator) Authenticate(ctx context.Context, _ *web.Request) web.Result {
	fakeAuthUrl, _ := url.Parse(fakeAuthURL)
	return a.responder.URLRedirect(fakeAuthUrl)
}

func (i *Identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {

	fakeSubject := "" // TODO

	return &Identity{
		subject: fakeSubject,
		broker:  i.broker,
	}, nil
}
