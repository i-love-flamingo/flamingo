package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	idpController struct {
		template  string
		responder *web.Responder
	}

	viewData struct {
		FormURL string
	}
)

func (c *idpController) FakeAuth(_ context.Context, request *web.Request) web.Result {
	return c.responder.Render(c.template, viewData{
		FormURL: "",
	})
}
