package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// CallbackControllerInterface is the callback HTTP action provider
	CallbackControllerInterface interface {
		Get(context.Context, *web.Request) web.Result
	}

	// CallbackController handles the oauth2.0 callback
	CallbackController struct {
		responder      *web.Responder
		authManager    *application.AuthManager
		logger         flamingo.Logger
		eventPublisher *application.EventPublisher
		userService    application.UserServiceInterface
	}
)

// Inject CallbackController dependencies
func (cc *CallbackController) Inject(
	responder *web.Responder,
	authManager *application.AuthManager,
	logger flamingo.Logger,
	eventPublisher *application.EventPublisher,
	userService application.UserServiceInterface,
) {
	cc.responder = responder
	cc.authManager = authManager
	cc.logger = logger
	cc.eventPublisher = eventPublisher

	cc.userService = userService
}

// Get handler for callbacks
func (cc *CallbackController) Get(ctx context.Context, request *web.Request) web.Result {
	// Verify state and errors.
	defer cc.authManager.DeleteAuthState(request.Session())

	if state, ok := cc.authManager.LoadAuthState(request.Session()); !ok || state != request.Request().URL.Query().Get("state") {
		cc.logger.Error("Invalid State", state, request.Request().URL.Query().Get("state"))
		return cc.responder.ServerError(errors.New("Invalid State"))
	}

	code := request.Request().URL.Query().Get("code")
	errCode := request.Request().URL.Query().Get("error")

	if code == "" && errCode == "" {
		err := errors.New("missing both code and error get parameter")
		cc.logger.Error("core.auth.callback Missing parameter", err)
		return cc.responder.ServerError(errors.WithStack(err))
	} else if code != "" {
		oauth2Token, err := cc.authManager.OAuth2Config(ctx, request).Exchange(cc.authManager.OAuthCtx(ctx), code)
		if err != nil {
			cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
			return cc.responder.ServerError(errors.WithStack(err))
		}

		err = cc.authManager.StoreTokenDetails(request.Session(), oauth2Token)
		if err != nil {
			cc.logger.Error("core.auth.callback Error", err)
			return cc.responder.ServerError(errors.WithStack(err))
		}
		cc.eventPublisher.PublishLoginEvent(ctx, &domain.LoginEvent{Session: request.Session()})
		cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
		cc.logger.Debugf("Token expiry %v", oauth2Token.Expiry)
		request.Session().AddFlash("successful logged in")
	} else if errCode != "" {
		cc.logger.Error("core.auth.callback Error parameter", errCode)
	}

	if redirect, ok := request.Session().Load("auth.redirect"); ok {
		request.Session().Delete("auth.redirect")
		redirectURL, _ := url.Parse(redirect.(string))
		return cc.responder.URLRedirect(redirectURL)
	}
	return cc.responder.RouteRedirect("home", nil)
}
