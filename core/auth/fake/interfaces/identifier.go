package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/fake/domain"
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
	_ auth.WebCallbacker     = (*Identifier)(nil)
)

// Broker returns the broker id from the config
func (i *Identifier) Broker() string {
	return i.broker
}

// Authenticate action, fake
func (i *Identifier) Authenticate(_ context.Context, _ *web.Request) web.Result {
	authURL, _ := url.Parse(fakeAuthURL)
	urlValues := url.Values{}
	urlValues.Add("broker", i.Broker())
	authURL.RawQuery = urlValues.Encode()

	return i.responder.URLRedirect(authURL)
}

// Identify action, fake
func (i *Identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	fakeSubject := "" // TODO

	return domain.NewIdentity(fakeSubject, i.broker), nil
}

func (i *Identifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	panic("not implemtddededdd")
}
