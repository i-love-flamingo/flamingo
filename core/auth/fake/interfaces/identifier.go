package interfaces

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/fake/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Identifier is the fake Identifier implementation
	Identifier struct {
		responder   *web.Responder
		broker      string
		eventRouter flamingo.EventRouter
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
	_ = request

	sess := web.SessionFromContext(ctx)
	userSubject, ok := sess.Load(fmt.Sprintf(userNameSessionKey, i.broker))
	if !ok {
		return nil, errors.New("identity not saved in session")
	}

	return domain.NewIdentity(userSubject.(string), i.broker), nil
}

// Callback from fake idp
func (i *Identifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	identity, err := i.Identify(ctx, request)
	if err != nil {
		i.Logout(ctx, request)

		return i.responder.ServerError(err)
	}

	i.eventRouter.Dispatch(ctx, &auth.WebLoginEvent{Broker: i.broker, Request: request, Identity: identity})

	return i.responder.URLRedirect(returnTo(request))
}

// Logout logs out
func (i *Identifier) Logout(ctx context.Context, request *web.Request) {
	_ = request

	sess := web.SessionFromContext(ctx)
	sess.Delete(fmt.Sprintf(userNameSessionKey, i.broker))
}
