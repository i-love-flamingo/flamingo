package interfaces

import (
	"context"
	"errors"
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
	_ auth.WebLogouter       = (*Identifier)(nil)
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
	sess := web.SessionFromContext(ctx)
	userSubject, ok := sess.Load(userNameSessionKey)
	if !ok {
		return nil, errors.New("identity not saved in session")
	}

	return domain.NewIdentity(userSubject.(string), i.broker), nil
}

// Callback from fake idp
func (i *Identifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {

}

// Logout logs out
func (i *Identifier) Logout(ctx context.Context, request *web.Request) {
	sess := web.SessionFromContext(ctx)
	sess.Delete(userNameSessionKey)
}
