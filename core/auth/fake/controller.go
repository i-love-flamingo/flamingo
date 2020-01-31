package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	idpController struct {
		responder     *web.Responder
		reverseRouter web.ReverseRouter
		template      string
	}

	viewData struct {
		FormURL string
	}
)

const defaultIDPContentHtml = `
<body>
  <h1>Login!</h1>
  <form action="{{.FormURL}}">
  </form>
</body>
`

// Inject injects module dependencies
func (c *idpController) Inject(
	responder *web.Responder,
	reverseRouter web.ReverseRouter,
	cfg *struct {
		Template string `inject:"config:auth.fake.loginTemplate"`
	},
) *idpController {
	c.responder = responder
	c.reverseRouter = reverseRouter

	if cfg != nil {
		c.template = cfg.Template
	}

	return c
}

// Auth action to simulate OIDC / Oauth Login Page
func (c *idpController) Auth(_ context.Context, r *web.Request) web.Result {
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
