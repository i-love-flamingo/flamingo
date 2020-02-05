package interfaces

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"flamingo.me/flamingo/v3/core/auth/fake/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// IdpController is the fake identity provider controller
	IdpController struct {
		responder        *web.Responder
		reverseRouter    web.ReverseRouter
		template         string
		userConfig       config.Map
		validatePassword bool
		validateOtp      bool
		usernameFieldID  string
		passwordFieldID  string
		otpFieldID       string
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
	defaultOtpFieldID      = "m2fa-otp"

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
	cfg *struct {
		Template         string     `inject:"config:core.auth.fake.loginTemplate,optional"`
		UserConfig       config.Map `inject:"config:core.auth.fake.userConfig"`
		ValidatePassword bool       `inject:"config:core.auth.fake.validatePassword,optional"`
		ValidateOtp      bool       `inject:"config:core.auth.fake.validateOtp,optional"`
		UsernameFieldID  string     `inject:"config:core.auth.fake.usernameFieldId,optional"`
		PasswordFieldID  string     `inject:"config:core.auth.fake.passwordFieldId,optional"`
		OtpFieldID       string     `inject:"config:core.auth.fake.otpFieldId,optional"`
	},
) *IdpController {
	c.responder = responder
	c.reverseRouter = reverseRouter

	if cfg != nil {
		c.template = cfg.Template
		c.userConfig = cfg.UserConfig
		c.validatePassword = cfg.ValidatePassword
		c.validateOtp = cfg.ValidateOtp
		c.usernameFieldID = cfg.UsernameFieldID
		c.passwordFieldID = cfg.PasswordFieldID
		c.otpFieldID = cfg.OtpFieldID
	}

	if c.usernameFieldID == "" {
		c.usernameFieldID = defaultUserNameFieldID
	}
	if c.passwordFieldID == "" {
		c.passwordFieldID = defaultPasswordFieldID
	}
	if c.otpFieldID == "" {
		c.otpFieldID = defaultOtpFieldID
	}

	return c
}

// Auth action to simulate OIDC / Oauth Login Page
func (c *IdpController) Auth(ctx context.Context, r *web.Request) web.Result {
	broker, ok := r.Params["broker"]
	if !ok || broker == "" {
		return c.responder.ServerError(errors.New("broker nor known"))
	}

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

	if c.template != "" {
		return c.responder.Render(c.template, viewData{
			FormURL:    formURL.String(),
			Message:    formError.Error(),
			UsernameID: c.usernameFieldID,
			PasswordID: c.passwordFieldID,
			OtpID:      c.otpFieldID,
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
			FormURL:    formURL.String(),
			Message:    formError.Error(),
			UsernameID: c.usernameFieldID,
			PasswordID: c.passwordFieldID,
			OtpID:      c.otpFieldID,
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
	usernameVal, ok := values[c.usernameFieldID]
	if !ok {
		return errors.New(errMissingUsername)
	}

	user := usernameVal[0]

	userCfgRaw, found := c.userConfig.Get(user)
	if !found {
		return errors.New(errInvalidUser)
	}

	userCfg := userCfgRaw.(config.Map)

	if c.validatePassword {
		passwordVal, ok := values[c.passwordFieldID]
		if !ok {
			return errors.New(errMissingPassword)
		}

		expectedPassword := passwordVal[0]
		userPasswordRaw, found := userCfg.Get("password")
		if !found {
			return errors.New(errFakeConfigFaulty)
		}

		userPassword := userPasswordRaw.(string)
		if expectedPassword != userPassword {
			return errors.New(errPasswordMismatch)
		}
	}

	if c.validateOtp {
		otpVal, ok := values[c.otpFieldID]
		if !ok {
			return errors.New(errMissingOtp)
		}

		expectedOtp := otpVal[0]
		userOtpRaw, found := userCfg.Get("otp")
		if !found {
			return errors.New(errFakeConfigFaulty)
		}

		userOtp := userOtpRaw.(string)
		if expectedOtp != userOtp {
			return errors.New(errOtpMismatch)
		}
	}

	sessionData := domain.UserSessionData{Subject: user}
	r.Session().Store(fmt.Sprintf(userDataSessionKey, broker), sessionData)

	return nil
}
