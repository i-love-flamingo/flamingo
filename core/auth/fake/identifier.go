package fake

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"text/template"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Identifier is the fake Identifier implementation
	Identifier struct {
		responder     *web.Responder
		broker        string
		reverseRouter web.ReverseRouter
		eventRouter   flamingo.EventRouter
		config        fakeConfig
	}

	fakeConfig struct {
		Broker           string                `json:"broker"`
		LoginTemplate    string                `json:"loginTemplate"`
		ValidatePassword bool                  `json:"validatePassword"`
		UsernameFieldID  string                `json:"usernameFieldId"`
		PasswordFieldID  string                `json:"passwordFieldId"`
		UserConfig       map[string]userConfig `json:"userConfig"`
	}

	userConfig struct {
		Password string `json:"password"`
	}

	viewData struct {
		FormURL    string
		Message    string
		UsernameID string
		PasswordID string
	}
)

const (
	defaultLoginTemplate = `
<body>
  <h1>Login!</h1>
  <form name="fake-login-form" action="{{.FormURL}}" method="post">
	<div>{{.Message}}</div>
	<label for="{{.UsernameID}}">Username</label>
	<input type="text" name="{{.UsernameID}}" id="{{.UsernameID}}">
	<label for="{{.PasswordID}}">Password</label>
    <input type="password" name="{{.PasswordID}}" id="{{.PasswordID}}">
	<button type="submit" id="submit">Fake Login</button>
  </form>
</body>
`

	defaultUserNameFieldID = "username"
	defaultPasswordFieldID = "password"

	userDataSessionKey = "core.auth.fake.%s.data"
)

var (
	_ auth.RequestIdentifier = (*Identifier)(nil)
	_ auth.WebCallbacker     = (*Identifier)(nil)
	_ auth.WebLogouter       = (*Identifier)(nil)

	errMissingUsername  = errors.New("missing username")
	errInvalidUser      = errors.New("invalid user")
	errMissingPassword  = errors.New("missing password")
	errPasswordMismatch = errors.New("password mismatch")
)

// IdentityProviderFactory -
func IdentityProviderFactory(cfg config.Map) (auth.RequestIdentifier, error) {
	var fakeConfig fakeConfig
	if err := cfg.MapInto(&fakeConfig); err != nil {
		return nil, err
	}

	return &Identifier{broker: fakeConfig.Broker, config: fakeConfig}, nil
}

// Inject injects module dependencies
func (i *Identifier) Inject(
	reverseRouter web.ReverseRouter,
	responder *web.Responder,
	eventRouter flamingo.EventRouter,
) *Identifier {
	i.reverseRouter = reverseRouter
	i.responder = responder
	i.eventRouter = eventRouter

	return i
}

// Broker returns the broker id from the config
func (i *Identifier) Broker() string {
	return i.broker
}

// Authenticate action, fake
func (i *Identifier) Authenticate(_ context.Context, r *web.Request) web.Result {
	var formError error

	postValues, err := r.FormAll()
	if err == nil {
		delete(postValues, "broker")
		if len(postValues) > 0 {
			formError = i.handlePostValues(r, postValues, i.broker)

			if formError == nil {
				return i.responder.RouteRedirect("core.auth.callback", map[string]string{"broker": i.broker})
			}
		}
	}

	// get formURL to callback with broker filled in
	formURL, err := i.reverseRouter.Absolute(r, "core.auth.login", map[string]string{"broker": i.broker})
	if err != nil {
		return i.responder.ServerError(err)
	}

	var loginTemplate string
	if i.config.LoginTemplate != "" {
		loginTemplate = i.config.LoginTemplate
	} else {
		loginTemplate = defaultLoginTemplate
	}

	t := template.New("fake")
	t, err = t.Parse(loginTemplate)
	if err != nil {
		return i.responder.ServerError(err)
	}

	var body = new(bytes.Buffer)
	var errMsg string

	if formError != nil {
		errMsg = formError.Error()
	}

	err = t.Execute(
		body,
		viewData{
			FormURL:    formURL.String(),
			Message:    errMsg,
			UsernameID: i.config.UsernameFieldID,
			PasswordID: i.config.PasswordFieldID,
		})
	if err != nil {
		return i.responder.ServerError(err)
	}

	return &web.Response{
		Header: http.Header{"ContentType": []string{"text/html; charset=utf-8"}},
		Status: http.StatusOK,
		Body:   body,
	}
}

func (i *Identifier) handlePostValues(r *web.Request, values map[string][]string, broker string) error {
	usernameVal, ok := values[i.config.UsernameFieldID]
	if !ok {
		return errMissingUsername
	}

	user := usernameVal[0]

	userCfg, found := i.config.UserConfig[user]
	if !found {
		return errInvalidUser
	}

	if i.config.ValidatePassword {
		passwordVal, ok := values[i.config.PasswordFieldID]
		if !ok {
			return errMissingPassword
		}

		expectedPassword := passwordVal[0]
		userPassword := userCfg.Password
		if expectedPassword != userPassword {
			return errPasswordMismatch
		}
	}

	sessionData := UserSessionData{Subject: user}
	r.Session().Store(fmt.Sprintf(userDataSessionKey, broker), sessionData)

	return nil
}

// Identify action, fake
func (i *Identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	userSessionData, ok := request.Session().Load(fmt.Sprintf(userDataSessionKey, i.broker))
	if !ok {
		return nil, errors.New("identity not saved in session")
	}

	if usd, ok := userSessionData.(UserSessionData); ok {
		return &identity{
			subject: usd.Subject,
			broker:  i.broker,
		}, nil
	}

	return nil, errors.New("session data not properly decoded")
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
func (i *Identifier) Logout(_ context.Context, request *web.Request) {
	request.Session().Delete(fmt.Sprintf(userDataSessionKey, i.broker))
}
