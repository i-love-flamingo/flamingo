package interfaces

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/responder"
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
		tokenExtras    config.Slice
		userService    application.UserServiceInterface
	}
)

// Inject CallbackController dependencies
func (cc *CallbackController) Inject(
	redirectAware responder.RedirectAware,
	errorAware responder.ErrorAware,
	authManager *application.AuthManager,
	logger flamingo.Logger,
	eventPublisher *application.EventPublisher,
	userService application.UserServiceInterface,
	cfg *struct {
		TokenExtras config.Slice `inject:"config:auth.tokenExtras"`
	},
) {
	cc.RedirectAware = redirectAware
	cc.ErrorAware = errorAware
	cc.authManager = authManager
	cc.logger = logger
	cc.eventPublisher = eventPublisher
	cc.tokenExtras = cfg.TokenExtras
	cc.userService = userService
}

// Get handler for callbacks
func (cc *CallbackController) Get(c context.Context, request *web.Request) web.Response {
	// Verify state and errors.
	defer request.Session().Delete(application.KeyAuthstate)

	if request.Session().Try(application.KeyAuthstate) != request.MustQuery1("state") {
		cc.logger.Error("Invalid State", request.Session().Try(application.KeyAuthstate), request.MustQuery1("state"))
		return cc.Error(c, errors.New("Invalid State"))
	}

	code, cOk := request.Query1("code")
	errCode, eOk := request.Query1("error")

	if !cOk && !eOk {
		err := errors.New("missing both code and error get parameter")
		cc.logger.Error("core.auth.callback Missing parameter", err)
		return cc.Error(c, errors.WithStack(err))
	} else if cOk {
		oauth2Token, err := cc.authManager.OAuth2Config(c).Exchange(c, code)
		if err != nil {
			cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
			return cc.Error(c, errors.WithStack(err))
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

		request.Session().Store(application.KeyToken, oauth2Token)
		rawToken, err := cc.authManager.ExtractRawIDToken(oauth2Token)
		request.Session().Store(application.KeyRawIDToken, rawToken)
		request.Session().Store(application.KeyTokenExtras, tokenExtras)
		if err != nil {
			cc.logger.Error("core.auth.callback Error ExtractRawIDToken", err)
			return cc.Error(c, errors.WithStack(err))
		}

		err = cc.userService.InitUser(c, request.Session().G())
		if err != nil {
			cc.logger.Error("core.auth.callback Error init user", err)
			return cc.Error(c, errors.WithStack(err))
		}

		cc.eventPublisher.PublishLoginEvent(c, &domain.LoginEvent{Session: request.Session().G()})
		cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
		request.Session().AddFlash("successful logged in", "info")
	} else if eOk {
		cc.logger.Error("core.auth.callback Error parameter", errCode)
	}

	if redirect, ok := request.Session().Load("auth.redirect"); ok {
		request.Session().Delete("auth.redirect")
		return cc.RedirectURL(redirect.(string))
	}
	return cc.Redirect("home", nil)
}
