package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	idpController struct {
		template      string
		responder     *web.Responder
		reverseRouter web.ReverseRouter
	}

	viewData struct {
		FormURL string
	}
)

// FakeAuth action to simulate OIDC / Oauth Login Page
func (c *idpController) FakeAuth(_ context.Context, r *web.Request) web.Result {
	broker, err := r.Query1("broker")
	if err != nil || broker == "" {
		return c.responder.ServerError(err)
	}

	formURL, err := c.reverseRouter.Absolute(r, "core.auth.callback(broker)", map[string]string{"broker": broker})
	if err != nil {
		return c.responder.ServerError(err)
	}

	return c.responder.Render(c.template, viewData{
		FormURL: formURL.String(),
	})
}
