package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/auth/domain"
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
		tokenExtras    config.Slice
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
	cfg *struct {
		TokenExtras config.Slice `inject:"config:auth.tokenExtras"`
	},
) {
	cc.responder = responder
	cc.authManager = authManager
	cc.logger = logger
	cc.eventPublisher = eventPublisher
	cc.tokenExtras = cfg.TokenExtras
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
		oauth2Token, err := cc.authManager.OAuth2Config(ctx).Exchange(ctx, code)
		if err != nil {
			cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
			return cc.responder.ServerError(errors.WithStack(err))
		}

		var extras []string
		err = cc.tokenExtras.MapInto(&extras)
		if err != nil {
			panic(err)
		}
		tokenExtras := &domain.TokenExtras{}
		for _, extra := range extras {
			value := oauth2Token.Extra(extra)
			parsed, ok := value.(string)
			if !ok {
				cc.logger.Error("core.auth.callback invalid type for extras", value)
				continue
			}
			tokenExtras.Add(extra, parsed)
		}

		rawToken, err := cc.authManager.ExtractRawIDToken(oauth2Token)
		if err != nil {
			cc.logger.Error("core.auth.callback Error ExtractRawIDToken", err)
			return cc.responder.ServerError(errors.WithStack(err))
		}
		cc.authManager.StoreTokenDetails(request.Session(), oauth2Token, rawToken, tokenExtras)

		cc.eventPublisher.PublishLoginEvent(ctx, &domain.LoginEvent{Session: request.Session()})
		cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
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
