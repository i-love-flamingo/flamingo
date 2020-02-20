package fake

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// controller is the fake identity provider controller
	controller struct {
		responder     *web.Responder
		reverseRouter web.ReverseRouter
		config        *fakeConfig
	}

	viewData struct {
		FormURL    string
		Message    error
		UsernameID string
		PasswordID string
		OtpID      string
	}
)

var (
	errMissingUsername  = errors.New("missing username")
	errInvalidUser      = errors.New("invalid user")
	errMissingPassword  = errors.New("missing password")
	errPasswordMismatch = errors.New("password mismatch")
	errOtpMismatch      = errors.New("otp mismatch")
	errMissingOtp       = errors.New("otp missing")
)

const (
	defaultUserNameFieldID = "username"
	defaultPasswordFieldID = "password"
	defaultOtpFieldID      = "otp"

	userDataSessionKey = "core.auth.fake.%s.data"
)

const defaultLoginTemplate = `
<body>
  <h1>Login!</h1>
  <form name="fake-login-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
	<label for="{{.UsernameID}}">Username</label>
	<input type="text" name="{{.UsernameID}}" id="{{.UsernameID}}">
	<label for="{{.PasswordID}}">Password</label>
    <input type="password" name="{{.PasswordID}}" id="{{.PasswordID}}">
	<label for="{{.OtpID}}">2 Factor OTP</label>
    <input type="text" name="{{.OtpID}}" id="{{.OtpID}}">
	<button type="submit" id="submit">Fake Login</button>
  </form>
</body>
`

// Inject injects module dependencies
func (c *controller) Inject(
	responder *web.Responder,
	reverseRouter web.ReverseRouter,
) *controller {
	c.responder = responder
	c.reverseRouter = reverseRouter

	return c
}

// Auth action to simulate OIDC / Oauth Login Page
func (c *controller) Auth(ctx context.Context, r *web.Request) web.Result {
	broker, ok := r.Params["broker"]
	if !ok || broker == "" {
		return c.responder.ServerError(errors.New("broker not known"))
	}

	config, found := identifierConfig[broker]
	if !found {
		return c.responder.ServerError(errors.New("broker not known"))
	}

	c.config = &config

	var formError error

	postValues, err := r.FormAll()
	if err == nil {
		delete(postValues, "broker")
		if len(postValues) > 0 {
			formError = c.handlePostValues(r, postValues, broker)

			if formError == nil {
				return c.responder.RouteRedirect("core.auth.callback", map[string]string{"broker": broker})
			}
		}
	}

	// get formURL to callback with broker filled in
	formURL, err := c.reverseRouter.Absolute(r, "core.auth.fake.auth", map[string]string{"broker": broker})
	if err != nil {
		return c.responder.ServerError(err)
	}

	var loginTemplate string
	if c.config.LoginTemplate != "" {
		loginTemplate = c.config.LoginTemplate
	} else {
		loginTemplate = defaultLoginTemplate
	}

	t := template.New("fake")
	t, err = t.Parse(loginTemplate)
	if err != nil {
		return c.responder.ServerError(err)
	}

	var body = new(bytes.Buffer)

	err = t.Execute(
		body,
		viewData{
			FormURL:    formURL.String(),
			Message:    formError,
			UsernameID: c.config.UsernameFieldID,
			PasswordID: c.config.PasswordFieldID,
			OtpID:      c.config.OtpFieldID,
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

func (c *controller) handlePostValues(r *web.Request, values map[string][]string, broker string) error {
	usernameVal, ok := values[c.config.UsernameFieldID]
	if !ok {
		return errMissingUsername
	}

	user := usernameVal[0]

	userCfg, found := c.config.UserConfig[user]
	if !found {
		return errInvalidUser
	}

	if c.config.ValidatePassword {
		passwordVal, ok := values[c.config.PasswordFieldID]
		if !ok {
			return errMissingPassword
		}

		expectedPassword := passwordVal[0]
		userPassword := userCfg.Password
		if expectedPassword != userPassword {
			return errPasswordMismatch
		}
	}

	if c.config.ValidateOtp {
		otpVal, ok := values[c.config.OtpFieldID]
		if !ok {
			return errMissingOtp
		}

		expectedOtp := otpVal[0]
		userOtp := userCfg.Otp
		if expectedOtp != userOtp {
			return errOtpMismatch
		}
	}

	sessionData := UserSessionData{Subject: user}
	r.Session().Store(fmt.Sprintf(userDataSessionKey, broker), sessionData)

	return nil
}
