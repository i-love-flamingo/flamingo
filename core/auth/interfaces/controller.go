package interfaces

import (
	"context"
	"errors"
	"net/url"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	Authcontroller struct {
		responder *web.Responder
		authservice Authservice
		router         web.ReverseRouter
	}
)


// Inject for Authcontroller
func (c *Authcontroller) Inject(responder *web.Responder, authservice Authservice,router web.ReverseRouter,) {
	c.responder = responder
	c.authservice = authservice
	c.router = router
}

//AuthAction - default auth action - starting the authorsation with the registered Authservice
func (c *Authcontroller) AuthAction(ctx context.Context, r *web.Request) web.Result {
	redirecturl, ok := r.Params["redirecturl"]
	if !ok || redirecturl == "" {
		redirecturl = r.Request().Referer()
	}

	if refURL, err := url.Parse(redirecturl); err != nil || refURL.Host != r.Request().Host {
		u, _ := c.router.Absolute(r, "", nil)
		redirecturl = u.String()
	}

	url, err := url.Parse(redirecturl)

	if err != nil {
		return c.responder.ServerError(errors.New("wrong redirect url given"))
	}
	result, err := c.authservice.Authenticate(ctx,url)
	if err != nil {
		return c.responder.ServerError(errors.New("wrong redirect url given"))
	}
	return result
}