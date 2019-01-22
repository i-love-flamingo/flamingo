package fake

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/application/fake"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/responder"
)

type (
	CallbackController struct {
		responder.RedirectAware

		mappingService *domain.UserMappingService

		userData config.Map
	}
)

func (c *CallbackController) Inject(
	redirectAware responder.RedirectAware,
	mappingService *domain.UserMappingService,
	config *struct {
		UserData config.Map `inject:"config:auth.fakeUserData"`
	},
) {
	c.RedirectAware = redirectAware
	c.mappingService = mappingService
	c.userData = config.UserData
}

func (c *CallbackController) Get(_ context.Context, request *web.Request) web.Response {
	user := c.mappingService.MapToUser(c.userData, request.Session())
	if user == nil {
		user = domain.Guest
	}
	request.Session().Store(fake.UserSessionKey, user)

	value, _ := request.Session().Load("auth.redirect")
	redirectUrl, ok := value.(string)
	if !ok || redirectUrl == "" {
		return c.Redirect("home", nil)
	}

	return c.RedirectURL(redirectUrl)
}
