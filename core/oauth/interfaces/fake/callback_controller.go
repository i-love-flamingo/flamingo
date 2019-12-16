package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/oauth/application/fake"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CallbackController fake controller
	CallbackController struct {
		responder *web.Responder

		mappingService *domain.UserMappingService

		userData config.Map
	}
)

// Inject dependencies
func (c *CallbackController) Inject(
	responder *web.Responder,
	mappingService *domain.UserMappingService,
	config *struct {
		UserData config.Map `inject:"config:core.oauth.fakeUserData"`
	},
) {
	c.responder = responder
	c.mappingService = mappingService
	c.userData = config.UserData
}

// Get http action
func (c *CallbackController) Get(_ context.Context, request *web.Request) web.Result {
	user := c.mappingService.MapToUser(c.userData, request.Session())
	if user == nil {
		user = domain.Guest
	}

	group, err := request.Query1("group")
	if err == nil {
		user.Groups = append(user.Groups, group)
	}

	request.Session().Store(fake.UserSessionKey, user)

	value, _ := request.Session().Load("auth.redirect")
	redirectURL, ok := value.(string)
	if !ok || redirectURL == "" {
		return c.responder.RouteRedirect("home", nil)
	}
	url, err := url.Parse(redirectURL)
	if err != nil {
		return c.responder.RouteRedirect("home", nil)
	}

	return c.responder.URLRedirect(url)
}
