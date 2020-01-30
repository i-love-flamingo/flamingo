package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Identifier is the fake Identifier implementation
	Identifier struct {
		responder *web.Responder
		broker    string
	}
)

const fakeAuthURL string = "/fake/auth"

var (
	_ auth.RequestIdentifier = (*Identifier)(nil)
)

// Broker returns the broker id from the config
func (i *Identifier) Broker() string {
	return i.broker
}

// Authenticate action, fake
func (i *Identifier) Authenticate(ctx context.Context, _ *web.Request) web.Result {
	authURL, _ := url.Parse(fakeAuthURL)
	authURL.Query().Add("broker", i.Broker())

	return i.responder.URLRedirect(authURL)
}

// Identify action, fake
func (i *Identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	fakeSubject := "" // TODO

	return &Identity{
		subject: fakeSubject,
		broker:  i.broker,
	}, nil
}
