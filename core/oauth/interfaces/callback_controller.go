package interfaces

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
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

var (
	// loginFailedCount counts the failed login attempts
	loginFailedCount = stats.Int64("flamingo/oauth/login_failed", "Count of failed login attempts", stats.UnitDimensionless)
	// loginSucceededCount counts the successful login attempts
	loginSucceededCount = stats.Int64("flamingo/oauth/login_succeeded", "Count of succeeded login attempts", stats.UnitDimensionless)
)

func init() {
	if err := opencensus.View("flamingo/oauth/login_failed", loginFailedCount, view.Count()); err != nil {
		panic(err)
	}
	if err := opencensus.View("flamingo/oauth/login_succeeded", loginSucceededCount, view.Count()); err != nil {
		panic(err)
	}
}

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
		if state != "" {
			cc.logger.Error(fmt.Sprintf("Invalid State - expected: %v  got: %v", state, request.Request().URL.Query().Get("state")))
		}

		stats.Record(ctx, loginFailedCount.M(1))
		return cc.responder.ServerError(errors.New("invalid state"))
	}

	// cleanup previous token information
	cc.authManager.DeleteTokenDetails(request.Session())

	code := request.Request().URL.Query().Get("code")
	errCode := request.Request().URL.Query().Get("error")

	if code == "" && errCode == "" {
		err := errors.New("missing both code and error get parameter")
		cc.logger.Error("core.auth.callback Missing parameter", err)
		stats.Record(ctx, loginFailedCount.M(1))
		return cc.responder.ServerError(err)
	} else if code != "" {
		oauth2Token, err := cc.authManager.OAuth2Config(ctx, request).Exchange(cc.authManager.OAuthCtx(ctx), code)
		if err != nil {
			cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
			stats.Record(ctx, loginFailedCount.M(1))
			return cc.responder.ServerError(fmt.Errorf("core.auth.callback error in OAuth2Config Exchange: %w", err))
		}

		err = cc.authManager.StoreTokenDetails(ctx, request.Session(), oauth2Token)
		if err != nil {
			cc.logger.Error("core.auth.callback Error", err)
			stats.Record(ctx, loginFailedCount.M(1))
			return cc.responder.ServerError(fmt.Errorf("core.auth.StoreTokenDetails error %w", err))
		}
		cc.eventPublisher.PublishLoginEvent(ctx, &domain.LoginEvent{Session: request.Session()})
		cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
		cc.logger.Debugf("Token expiry %v", oauth2Token.Expiry)
		request.Session().AddFlash("successful logged in")
		stats.Record(ctx, loginSucceededCount.M(1))
	} else if errCode != "" {
		cc.logger.Error("core.auth.callback Error parameter", errCode)
		stats.Record(ctx, loginFailedCount.M(1))
	}

	if redirect, ok := request.Session().Load("auth.redirect"); ok {
		request.Session().Delete("auth.redirect")
		redirectURL, _ := url.Parse(redirect.(string))
		return cc.responder.URLRedirect(redirectURL)
	}
	return cc.responder.RouteRedirect("home", nil)
}
