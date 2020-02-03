package interfaces

import (
	"bytes"
	"context"
	"errors"
	"flamingo.me/flamingo/v3/framework/config"
	"html/template"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
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
		FormURL string
		Message string
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
)

const defaultIDPTemplate = `
<body>
  <h1>Login!</h1>
  <form name="fake-idp-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
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
func (c *IdpController) Inject(
	responder *web.Responder,
	reverseRouter web.ReverseRouter,
	cfg *struct {
		Template         string     `inject:"config:auth.fake.loginTemplate,optional"`
		UserConfig       config.Map `inject:"config:auth.fake.userConfig"`
		ValidatePassword bool       `inject:"config:auth.fake.validatePassword,optional"`
		ValidateOtp      bool       `inject:"config:auth.fake.validateOtp,optional"`
		UsernameFieldID  string     `inject:"config:auth.fake.usernameFieldId,optional"`
		PasswordFieldID  string     `inject:"config:auth.fake.passwordFieldId,optional"`
		OtpFieldID       string     `inject:"config:auth.fake.otpFieldId,optional"`
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
	broker, err := r.Query1("broker")
	if err != nil || broker == "" {
		return c.responder.ServerError(err)
	}

	formError := errors.New("")

	postValues, err := r.FormAll()
	if err == nil {
		delete(postValues, "broker")
		if len(postValues) > 0 {
			formError = c.handlePostValues(ctx, postValues, broker)

			if formError == nil {
				return c.responder.RouteRedirect("core.auth.callback(broker)", map[string]string{"broker": broker})
			}
		}
	}

	// get formURL to callback with broker filled in
	formURL, err := c.reverseRouter.Absolute(r, "core.auth.callback(broker)", map[string]string{"broker": broker})
	if err != nil {
		return c.responder.ServerError(err)
	}

	if c.template != "" {
		return c.responder.Render(c.template, viewData{
			FormURL: formURL.String(),
			Message: formError.Error(),
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
			Message: formError.Error(),
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

func (c *IdpController) handlePostValues(ctx context.Context, values map[string][]string, broker string) error {
	// TODO: make this configurable
	// TODO: make sure password and otp are only verified when configured

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

	return nil
}
