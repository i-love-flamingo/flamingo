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
	// identifier is the fake identifier implementation
	identifier struct {
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
	_ auth.RequestIdentifier = (*identifier)(nil)
	_ auth.WebCallbacker     = (*identifier)(nil)
	_ auth.WebLogouter       = (*identifier)(nil)

	errMissingUsername           = errors.New("missing username")
	errInvalidUser               = errors.New("invalid user")
	errMissingPassword           = errors.New("missing password")
	errPasswordMismatch          = errors.New("password mismatch")
	errIdentityNotSavedInSession = errors.New("identity not saved in session")
	errSessionDataInvalid        = errors.New("session data not properly decoded")
)

func identityProviderFactory(cfg config.Map) (auth.RequestIdentifier, error) {
	var fakeConfig fakeConfig
	if err := cfg.MapInto(&fakeConfig); err != nil {
		return nil, err
	}

	return &identifier{broker: fakeConfig.Broker, config: fakeConfig}, nil
}

// Inject injects module dependencies
func (i *identifier) Inject(
	reverseRouter web.ReverseRouter,
	responder *web.Responder,
	eventRouter flamingo.EventRouter,
) *identifier {
	i.reverseRouter = reverseRouter
	i.responder = responder
	i.eventRouter = eventRouter

	return i
}

// Broker returns the broker id from the config
func (i *identifier) Broker() string {
	return i.broker
}

// Authenticate action, fake
func (i *identifier) Authenticate(_ context.Context, r *web.Request) web.Result {
	return i.prepareFormResponse(nil, r)
}

func (i *identifier) prepareFormResponse(formError error, r *web.Request) web.Result {
	callbackURL, err := i.reverseRouter.Absolute(r, "core.auth.callback", map[string]string{"broker": i.broker})
	if err != nil {
		return i.responder.ServerError(err)
	}

	loginTemplate := defaultLoginTemplate
	if i.config.LoginTemplate != "" {
		loginTemplate = i.config.LoginTemplate
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
			FormURL:    callbackURL.String(),
			Message:    errMsg,
			UsernameID: i.config.UsernameFieldID,
			PasswordID: i.config.PasswordFieldID,
		})
	if err != nil {
		return i.responder.ServerError(err)
	}

	response := i.responder.HTTP(http.StatusOK, body)
	response.Header.Set("Content-Type", "text/html; charset=utf-8")
	return response
}

func (i *identifier) handlePostValues(r *web.Request) error {
	username, err := r.Form1(i.config.UsernameFieldID)
	if err != nil {
		return errMissingUsername
	}

	userCfg, found := i.config.UserConfig[username]
	if !found {
		return errInvalidUser
	}

	if i.config.ValidatePassword {
		password, err := r.Form1(i.config.PasswordFieldID)
		if err != nil {
			return errMissingPassword
		}

		if userCfg.Password != password {
			return errPasswordMismatch
		}
	}

	r.Session().Store(fmt.Sprintf(userDataSessionKey, i.broker), UserSessionData{Subject: username})

	return nil
}

// Identify action, fake
func (i *identifier) Identify(ctx context.Context, request *web.Request) (auth.Identity, error) {
	userSessionData, ok := request.Session().Load(fmt.Sprintf(userDataSessionKey, i.broker))
	if !ok {
		return nil, errIdentityNotSavedInSession
	}

	if usd, ok := userSessionData.(UserSessionData); ok {
		return &identity{
			subject: usd.Subject,
			broker:  i.broker,
		}, nil
	}

	return nil, errSessionDataInvalid
}

// Callback from fake idp
func (i *identifier) Callback(ctx context.Context, request *web.Request, returnTo func(*web.Request) *url.URL) web.Result {
	_, err := i.Identify(ctx, request)
	if err == nil {
		return i.responder.URLRedirect(returnTo(request))
	}

	i.Logout(ctx, request)

	if request.Request().Method == http.MethodPost {
		err = i.handlePostValues(request)
		if err == nil {
			identity, _ := i.Identify(ctx, request)
			i.eventRouter.Dispatch(ctx, &auth.WebLoginEvent{Request: request, Broker: i.broker, Identity: identity})
			return i.responder.URLRedirect(returnTo(request))
		}
	}

	return i.prepareFormResponse(err, request)
}

// Logout logs out
func (i *identifier) Logout(_ context.Context, request *web.Request) {
	request.Session().Delete(fmt.Sprintf(userDataSessionKey, i.broker))
}
