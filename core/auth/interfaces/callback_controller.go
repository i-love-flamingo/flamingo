package interfaces

import (
	"context"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
	"github.com/pkg/errors"
)

type (
	CallbackControllerInterface interface {
		Get(context.Context, *web.Request) web.Response
	}

	// CallbackController handles the oauth2.0 callback
	CallbackController struct {
		responder.RedirectAware
		responder.ErrorAware
		authManager    *application.AuthManager
		logger         flamingo.Logger
		eventPublisher *application.EventPublisher
	}
)

// Inject CallbackController dependencies
func (cc *CallbackController) Inject(
	redirectAware responder.RedirectAware,
	errorAware responder.ErrorAware,
	authManager *application.AuthManager,
	logger flamingo.Logger,
	eventPublisher *application.EventPublisher,
) {
	cc.RedirectAware = redirectAware
	cc.ErrorAware = errorAware
	cc.authManager = authManager
	cc.logger = logger
	cc.eventPublisher = eventPublisher
}

// Get handler for callbacks
func (cc *CallbackController) Get(c context.Context, request *web.Request) web.Response {
	// Verify state and errors.
	defer request.Session().Delete(application.KeyAuthstate)

	if request.Session().Try(application.KeyAuthstate) != request.MustQuery1("state") {
		cc.logger.Error("Invalid State", request.Session().Try(application.KeyAuthstate), request.MustQuery1("state"))
		return cc.Error(c, errors.New("Invalid State"))
	}

	oauth2Token, err := cc.authManager.OAuth2Config(c).Exchange(c, request.MustQuery1("code"))
	if err != nil {
		cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
		return cc.Error(c, errors.WithStack(err))
	}

	request.Session().Store(application.KeyToken, oauth2Token)
	rawToken, err := cc.authManager.ExtractRawIDToken(oauth2Token)
	request.Session().Store(application.KeyRawIDToken, rawToken)
	if err != nil {
		cc.logger.Error("core.auth.callback Error ExtractRawIDToken", err)
		return cc.Error(c, errors.WithStack(err))
	}
	cc.eventPublisher.PublishLoginEvent(c, &domain.LoginEvent{Session: request.Session().G()})
	cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
	request.Session().AddFlash("successful logged in", "info")

	if redirect, ok := request.Session().Load("auth.redirect"); ok {
		request.Session().Delete("auth.redirect")
		return cc.RedirectURL(redirect.(string))
	}
	return cc.Redirect("home", nil)
}
