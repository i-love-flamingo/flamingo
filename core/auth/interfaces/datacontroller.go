package interfaces

import (
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	// UserController uc
	UserController struct {
		userService *application.UserService
	}
)

// Inject UserController dependencies
func (u *UserController) Inject(service *application.UserService) {
	u.userService = service
}

// Data controller to return userinfo
func (u *UserController) Data(c web.Context) interface{} {
	return u.userService.GetUser(c, c.Session())
}
