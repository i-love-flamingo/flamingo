package interfaces

import (
	"flamingo/core/auth/application"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"

	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type (
	// LoginController handles the login redirect
	LoginController struct {
		*responder.RedirectAware `inject:""`
		AuthManager              *application.AuthManager `inject:""`
	}

	// LogoutController handles the logout
	LogoutController struct {
		*responder.RedirectAware `inject:""`
	}

	// CallbackController handles the oauth2.0 callback
	CallbackController struct {
		*responder.RedirectAware `inject:""`
		*responder.ErrorAware    `inject:""`
		AuthManager              *application.AuthManager `inject:""`
	}
)

// Get handler for logins (redirect)
func (l *LoginController) Get(c web.Context) web.Response {
	state := uuid.NewV4().String()
	c.Session().Values["auth.state"] = state

	return l.RedirectUrl(l.AuthManager.OAuth2Config().AuthCodeURL(state))
}

// Get handler for logout
func (l *LogoutController) Get(c web.Context) web.Response {
	delete(c.Session().Values, application.KEY_AUTHSTATE)
	delete(c.Session().Values, application.KEY_RAWIDTOKEN)
	delete(c.Session().Values, application.KEY_TOKEN)

	return l.Redirect("home")
}

// Get handler for callbacks
func (cc *CallbackController) Get(c web.Context) web.Response {
	// Verify state and errors.
	defer delete(c.Session().Values, application.KEY_AUTHSTATE)

	if c.Session().Values[application.KEY_AUTHSTATE] != c.MustQuery1("state") {
		return cc.Error(c, errors.New("Invalid State"))
	}

	finish := c.Profile("auth.callback", "code: "+c.MustQuery1("code"))
	oauth2Token, err := cc.AuthManager.OAuth2Config().Exchange(c, c.MustQuery1("code"))
	finish()
	if err != nil {
		return cc.Error(c, errors.WithStack(err))
	}

	c.Session().Values[application.KEY_TOKEN] = oauth2Token
	c.Session().Values[application.KEY_RAWIDTOKEN], err = cc.AuthManager.ExtractRawIdToken(oauth2Token)
	if err != nil {
		return cc.Error(c, errors.WithStack(err))
	}

	return cc.Redirect("home")
}
