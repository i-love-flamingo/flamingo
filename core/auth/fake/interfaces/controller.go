package interfaces

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"flamingo.me/flamingo/v3/core/auth/fake/domain"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// IdpController is the fake identity provider controller
	IdpController struct {
		responder     *web.Responder
		reverseRouter web.ReverseRouter
		config        *fakeConfig
	}

	viewData struct {
		FormURL    string
		Message    string
		UsernameID string
		PasswordID string
		OtpID      string
	}
)

const (
	errMissingUsername  = "missing username"
	errInvalidUser      = "invalid user"
	errMissingPassword  = "missing password"
	errFakeConfigFaulty = "fake auth config error"
	errPasswordMismatch = "password mismatch"
	errOtpMismatch      = "otp mismatch"
	errMissingOtp       = "otp missing"

	defaultUserNameFieldID = "username"
	defaultPasswordFieldID = "password"
	defaultOtpFieldID      = "otp"

	userDataSessionKey = "core.auth.fake.%s.data"
)

const defaultIDPTemplate = `
<body>
  <h1>Login!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
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
func (c *IdpController) Inject(
	responder *web.Responder,
	reverseRouter web.ReverseRouter,
) *IdpController {
	c.responder = responder
	c.reverseRouter = reverseRouter

	return c
}

// Auth action to simulate OIDC / Oauth Login Page
func (c *IdpController) Auth(ctx context.Context, r *web.Request) web.Result {
	broker, ok := r.Params["broker"]
	if !ok || broker == "" {
		return c.responder.ServerError(errors.New("broker not known"))
	}

	config, found := identifierConfig[broker]
	if !found {
		return c.responder.ServerError(errors.New("broker not known"))
	}

	c.config = &config

	formError := errors.New("")

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

	var idpTemplate string
	if c.config.LoginTemplate != "" {
		idpTemplate = c.config.LoginTemplate
	} else {
		idpTemplate = defaultIDPTemplate
	}

	t := template.New("fake")
	t, err = t.Parse(idpTemplate)
	if err != nil {
		return c.responder.ServerError(err)
	}

	var body = new(bytes.Buffer)

	err = t.Execute(
		body,
		viewData{
			FormURL:    formURL.String(),
			Message:    formError.Error(),
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

func (c *IdpController) handlePostValues(r *web.Request, values map[string][]string, broker string) error {
	usernameVal, ok := values[c.config.UsernameFieldID]
	if !ok {
		return errors.New(errMissingUsername)
	}

	user := usernameVal[0]

	userCfg, found := c.config.UserConfig[user]
	if !found {
		return errors.New(errInvalidUser)
	}

	if c.config.ValidatePassword {
		passwordVal, ok := values[c.config.PasswordFieldID]
		if !ok {
			return errors.New(errMissingPassword)
		}

		expectedPassword := passwordVal[0]
		userPassword := userCfg.Password
		if expectedPassword != userPassword {
			return errors.New(errPasswordMismatch)
		}
	}

	if c.config.ValidateOtp {
		otpVal, ok := values[c.config.OtpFieldID]
		if !ok {
			return errors.New(errMissingOtp)
		}

		expectedOtp := otpVal[0]
		userOtp := userCfg.Otp
		if expectedOtp != userOtp {
			return errors.New(errOtpMismatch)
		}
	}

	sessionData := domain.UserSessionData{Subject: user}
	r.Session().Store(fmt.Sprintf(userDataSessionKey, broker), sessionData)

	return nil
}
