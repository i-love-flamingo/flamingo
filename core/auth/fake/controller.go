package fake

import (
	"bytes"
	"context"
	"html/template"
	"net/http"

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

const defaultIDPTemplate = `
<body>
  <h1>Login!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
	<label for="username">Username</label>   
	<input type="text" name="username" id="username">
	<label for="password">Password</label>
    <input type="password" name="password" id="password">
	<label for="m2fa-otp">2 Factor OTP</label>    
    <input type="text" name="m2fa-otp" id="m2fa-otp">
	<button type="submit" id="submit">Fake Login</button> 
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

	// get formURL to callback with broker filled in
	formURL, err := c.reverseRouter.Absolute(r, "core.auth.callback(broker)", map[string]string{"broker": broker})
	if err != nil {
		return c.responder.ServerError(err)
	}

	if c.template != "" {
		return c.responder.Render(c.template, viewData{
			FormURL: formURL.String(),
		})
	}

	// no custom template specified, use fallback template

	t := template.New("fake")

	t, err = t.Parse(defaultIDPTemplate)
	if err != nil {
		return c.responder.ServerError(err)
	}

	var body = new(bytes.Buffer)

	err = t.Execute(
		body,
		viewData{
			FormURL: formURL.String(),
		})
	if err != nil {
		return c.responder.ServerError(err)
	}

	return &web.Response{
		Header: http.Header{"ContentType": []string{"text/html; charset=utf-8"}},
		Status: http.StatusOK,
		Body:   body,
	}
}
